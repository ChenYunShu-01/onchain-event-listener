package executor

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/reddio-com/red-adapter/config"
	"github.com/reddio-com/red-adapter/pkg/starkex"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/reddio-com/red-adapter/pkg/contracts"
	rdb "github.com/reddio-com/red-adapter/pkg/db"
	event2 "github.com/reddio-com/red-adapter/pkg/event"
	"github.com/reddio-com/red-adapter/types"
)

const (
	ProjectStartBlock = uint64(7019341)
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

	executor.db.InitL1AdapterTable()

	return executor, nil
}

func (e *Executor) StartToWatchEvent() {
	client, _ := ethclient.Dial(e.cfg.Chains["goerli"].RPC)
	eventBlockGap := e.cfg.EventBlockGap
	go func() {
		err := watchEvent(contracts.Deposit, types.LogDeposit, e.db, client, eventBlockGap)
		if err != nil {
			fmt.Println(err)

		}
	}()
	go func() {
		err := watchEvent(contracts.Deposit, types.LogNftDeposit, e.db, client, eventBlockGap)
		if err != nil {
			fmt.Println(err)
		}
	}()
	go func() {
		err := watchEvent(contracts.Deposit, types.LogDepositWithTokenId, e.db, client, eventBlockGap)
		if err != nil {
			fmt.Println(err)
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

	logTicker := time.Tick(3 * time.Second)
	infoTicker := time.Tick(1 * time.Minute)

	for {
		select {
		case <-logTicker:
			currentBlockNumber, err = client.BlockNumber(context.Background())
			if err != nil {
				log.Println(err)
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
				log.Println(err)
				continue
			}
			for _, log := range logs {
				currenEventLog, err := event2.FromOriginLogToEventLog(log, eventName)
				if err != nil {
					return err
				}
				if currenEventLog.NewerThan(latestEvent) {
					db.Create(&currenEventLog)
					l2DepositRequest := computeL2DepositRequest(log)
					txid, err := starkex.Deposit(l2DepositRequest)
					if err != nil {
						return err
					}
					fmt.Println("txid:", txid)
					latestEvent = currenEventLog
				}
			}
			parsedBlock = endBlockNumber
		case <-infoTicker:
			fmt.Println("eventName:", eventName, "currentBlockNumber:", currentBlockNumber, "parsed block", parsedBlock)
		}
	}
}

func computeL2DepositRequest(log ethtypes.Log) *types.L2DepositRequest {
	l2DepositRequest := &types.L2DepositRequest{}
	depositLog := &types.DepositLog{}
	contract := contracts.GetContractMeta(contracts.Deposit)
	contract.ToBoundContract().UnpackLog(depositLog, types.LogDeposit, log)
	l2DepositRequest.StarkKey = fmt.Sprint("0x", depositLog.StarkKey.Text(16))
	l2DepositRequest.VaultId = depositLog.VaultId
	l2DepositRequest.Amount = depositLog.QuantizedAmount.String()
	l2DepositRequest.TokenId = fmt.Sprint("0x", depositLog.AssetType.Text(16))

	return l2DepositRequest
}
