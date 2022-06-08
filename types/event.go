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
	if event.EventName == "" {
		return true
	}
	return false
}

func (event *EventLog) NewerThan(log *EventLog) bool {
	if event.BlockNumber > log.BlockNumber {
		return true
	} else if event.BlockNumber == log.BlockNumber {
		if event.LogIndex > log.LogIndex {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}
