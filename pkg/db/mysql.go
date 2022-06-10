package db

import (
	"github.com/reddio-com/red-adapter/types"
	"gorm.io/gorm"
)

type DB struct {
	*gorm.DB
}

func NewDB(db *gorm.DB) *DB {
	return &DB{
		DB: db,
	}
}

func (db *DB) InitL1AdapterTable() error {
	err := db.AutoMigrate(&types.EventLog{})
	return err
}

func (db *DB) GetLastEvent(eventName types.EventName) *types.EventLog {
	latestEvent := new(types.EventLog)
	db.Where("event_name = ?", eventName).Last(&latestEvent)
	return latestEvent
}
