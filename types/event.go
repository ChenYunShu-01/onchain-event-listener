package types

import (
	"gorm.io/gorm"
)

type EventName = string

const (
	LogDeposit            = EventName("LogDeposit")
	LogNftDeposit         = EventName("LogNftDeposit")
	LogDepositWithTokenId = EventName("LogDepositWithTokenId")
)

type EventLog struct {
	gorm.Model
	EventName   string `json:"event_name"`
	Data        string `json:"data"`
	BlockNumber uint64 `json:"block_number"`
	LogIndex    uint   `json:"log_index"`
}

func (event *EventLog) IsNil() bool {
	return event.EventName == ""
}

func (event *EventLog) NewerThan(log *EventLog) bool {
	if event.BlockNumber > log.BlockNumber {
		return true
	}
	if event.BlockNumber == log.BlockNumber {
		if event.LogIndex > log.LogIndex {
			return true
		}
		return false
	}
	return false

}
