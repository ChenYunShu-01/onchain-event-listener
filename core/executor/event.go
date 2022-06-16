package executor

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/reddio-com/red-adapter/config"
	"github.com/reddio-com/red-adapter/pkg/starkex"
	"github.com/reddio-com/starkex-contracts-source/source/deposits"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/go-logr/zapr"
	"github.com/reddio-com/red-adapter/pkg/contracts"
	rdb "github.com/reddio-com/red-adapter/pkg/db"
	event2 "github.com/reddio-com/red-adapter/pkg/event"
	"github.com/reddio-com/red-adapter/types"
	"github.com/reddio-com/starkex-utils/asset"
	starkex_types "github.com/reddio-com/starkex-utils/types"
	"go.uber.org/zap"
)

const (
	ProjectStartBlock = uint64(7019341)
)

var (
	zapLogger, _ = zap.NewDevelopment()
	loger        = zapr.NewLogger(zapLogger)
	logger       = loger.WithValues("package", "executor")
)

type Executor struct {
	db  *rdb.DB
	cfg *config.Cfg
}

func NewExecutor(cfg *config.Cfg) (*Executor, error) {
	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	executor := &Executor{
		db:  rdb.NewDB(db),
		cfg: cfg,
	}

	err = executor.db.InitL1AdapterTable()

	return executor, err
}

func (e *Executor) StartToWatchEvent() {
	client, _ := ethclient.Dial(e.cfg.Chains["goerli"].RPC)
	eventBlockGap := e.cfg.EventBlockGap
	go func() {
		err := watchEvent(contracts.Deposit, types.LogDeposit, e.db, client, eventBlockGap)
		if err != nil {
			logger.Error(err, "watch event failed")
		}
	}()
	go func() {
		err := watchEvent(contracts.Deposit, types.LogNftDeposit, e.db, client, eventBlockGap)
		if err != nil {
			logger.Error(err, "watch event failed")
		}
	}()
	go func() {
		err := watchEvent(contracts.Deposit, types.LogDepositWithTokenId, e.db, client, eventBlockGap)
		if err != nil {
			logger.Error(err, "watch event failed")
		}
	}()
	<-make(chan struct{})
}

func watchEvent(contractName contracts.ContractName, eventName types.EventName, db *rdb.DB, client *ethclient.Client, eventBlockGap uint64) error {
	var (
		latestEvent *types.EventLog
		filterQuery ethereum.FilterQuery
	)
	contractMeta := contracts.GetContractMeta(contractName)
	query := append([][]interface{}{{contractMeta.ABI.Events[eventName].ID}})
	topics, err := abi.MakeTopics(query...)
	if err != nil {
		return err
	}
	latestEvent = db.GetLastEvent(eventName)

	if latestEvent.IsNil() {
		latestEvent.BlockNumber = ProjectStartBlock
	}
	parsedBlock := latestEvent.BlockNumber
	currentBlockNumber := latestEvent.BlockNumber

	logTicker := time.Tick(1 * time.Minute)
	infoTicker := time.Tick(3 * time.Minute)

	for {
		select {
		case <-logTicker:
			currentBlockNumber, err = client.BlockNumber(context.Background())
			if err != nil {
				logger.Info("got error", "details:", err.Error())
				continue
			}

			if currentBlockNumber-parsedBlock < eventBlockGap {
				continue
			}

			endBlockNumber := currentBlockNumber - eventBlockGap
			filterQuery = ethereum.FilterQuery{
				Addresses: []common.Address{contractMeta.ContractAddress},
				Topics:    topics,
				FromBlock: new(big.Int).SetUint64(parsedBlock + 1),
				ToBlock:   new(big.Int).SetUint64(endBlockNumber),
			}
			logs, err := client.FilterLogs(context.Background(), filterQuery)
			if err != nil {
				logger.Info("got error", "details:", err.Error())
				continue
			}
			for _, log := range logs {
				currenEventLog, err := event2.FromOriginLogToEventLog(log, eventName)
				if err != nil {
					return err
				}
				if currenEventLog.NewerThan(latestEvent) {
					fmt.Println("event Name", eventName)
					db.Create(&currenEventLog)
					l2DepositRequest, err := computeL2DepositRequest(log, eventName)
					if err != nil {
						logger.Error(err, "compute l2 deposit request failed")
						return err
					}
					if l2DepositRequest.StarkKey != "0x38cae143fe6d2b8bdb7051f211744017d98f7e6a67e45a5dfc08759c119cf3c" && l2DepositRequest.StarkKey != "0x7d459f9c3ff9fda3073a4f793f809e1edcb6e4ef27a9a385f7e2b414d5d8e41" {
						continue
					}

					var txid int64

					for i := 0; i < 5; i++ {
						txid, err = starkex.Deposit(l2DepositRequest)
						if err == nil {
							break
						}

						if err != nil {
							logger.Error(err, "deposit failed, retrying")
						}
					}
					if err != nil {
						return err
					}

					logger.Info("send to starkex", "txid:", fmt.Sprint(txid), "origin hash:", log.TxHash.Hex())
					latestEvent = currenEventLog
				}
			}
			parsedBlock = endBlockNumber
		case <-infoTicker:
			logger.Info("info", "eventName:", eventName, "currentBlockNumber:", fmt.Sprint(currentBlockNumber), "parsed block", fmt.Sprint(parsedBlock))
		}
	}
}

func computeL2DepositRequest(log ethtypes.Log, eventName types.EventName) (*types.L2DepositRequest, error) {
	l2DepositRequest := &types.L2DepositRequest{}
	depositLog := &deposits.DepositsLogDeposit{}
	depositNFTLog := &deposits.DepositsLogNftDeposit{}
	contract := contracts.GetContractMeta(contracts.Deposit)
	var err error
	if eventName == types.LogDeposit {
		err = contract.ToBoundContract().UnpackLog(depositLog, types.LogDeposit, log)
		if err != nil {
			return l2DepositRequest, err
		}
		l2DepositRequest.StarkKey = fmt.Sprint("0x", depositLog.StarkKey.Text(16))
		l2DepositRequest.VaultId = depositLog.VaultId
		l2DepositRequest.Amount = depositLog.QuantizedAmount.String()
		l2DepositRequest.TokenId = fmt.Sprint("0x", depositLog.AssetType.Text(16))
	} else if eventName == types.LogNftDeposit {
		err = contract.ToBoundContract().UnpackLog(depositNFTLog, types.LogNftDeposit, log)
		if err != nil {
			return l2DepositRequest, err
		}
		l2DepositRequest.StarkKey = fmt.Sprint("0x", depositNFTLog.StarkKey.Text(16))
		l2DepositRequest.VaultId = depositNFTLog.VaultId
		l2DepositRequest.Amount = "1"
		l2DepositRequest.TokenId = fmt.Sprint("0x", asset.GetAssetIDByAssetType(starkex_types.ERC721, depositNFTLog.AssetType, depositNFTLog.TokenId).Text(16))
	}

	return l2DepositRequest, nil
}
