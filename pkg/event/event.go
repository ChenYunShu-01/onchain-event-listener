package event

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/reddio-com/red-adapter/pkg/contracts"
	types2 "github.com/reddio-com/red-adapter/types"
	"github.com/reddio-com/starkex-contracts-source/source/deposits"
)

func FromOriginLogToEventLog(log types.Log, name types2.EventName) (*types2.EventLog, error) {
	data, err := log.MarshalJSON()
	eventLog := new(types2.EventLog)
	eventLog.Data = string(data)
	eventLog.EventName = name
	eventLog.BlockNumber = log.BlockNumber
	eventLog.LogIndex = log.Index
	return eventLog, err
}

func FromEventLogToDepositNFTEvent(eventLog *types2.EventLog) (*deposits.DepositsLogNftDeposit, error) {
	depositContractMeta := contracts.GetContractMeta(contracts.Deposit)
	log := new(types.Log)
	err := log.UnmarshalJSON([]byte(eventLog.Data))
	if err != nil {
		return nil, err
	}
	nftLogNftDeposit := new(deposits.DepositsLogNftDeposit)
	err = depositContractMeta.ToBoundContract().UnpackLog(nftLogNftDeposit, types2.LogNftDeposit, *log)
	return nftLogNftDeposit, err
}

func FromEventLogToDepositEvent(eventLog *types2.EventLog) (*deposits.DepositsLogDeposit, error) {
	depositContractMeta := contracts.GetContractMeta(contracts.Deposit)
	log := new(types.Log)
	err := log.UnmarshalJSON([]byte(eventLog.Data))
	if err != nil {
		return nil, err
	}
	nftLogNftDeposit := new(deposits.DepositsLogDeposit)
	err = depositContractMeta.ToBoundContract().UnpackLog(nftLogNftDeposit, types2.LogDeposit, *log)
	return nftLogNftDeposit, err
}

func FromEventLogToDepositWithTokenIdEvent(eventLog *types2.EventLog) (*deposits.DepositsLogDepositWithTokenId, error) {
	depositContractMeta := contracts.GetContractMeta(contracts.Deposit)
	log := new(types.Log)
	err := log.UnmarshalJSON([]byte(eventLog.Data))
	if err != nil {
		return nil, err
	}
	logDepositWithTokenId := new(deposits.DepositsLogDepositWithTokenId)
	err = depositContractMeta.ToBoundContract().UnpackLog(logDepositWithTokenId, types2.LogDepositWithTokenId, *log)
	return logDepositWithTokenId, err
}
