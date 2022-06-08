package db

import (
	"fmt"

	"github.com/reddio-com/red-adapter/types"
	"gorm.io/gorm"
)

const (
	BaseVaultID = 1222222322
)

type DB struct {
	*gorm.DB
}

func NewDB(db *gorm.DB) *DB {
	return &DB{
		DB: db,
	}
}

func (db *DB) InitL1AdapterTable() {
	db.AutoMigrate(&types.EventLog{})
}

func (db *DB) InitTable() {
	db.AutoMigrate(&types.Contracts{}, &types.VaultModel{}, &types.MintModel{}, &types.NonceModel{})
}

func (db *DB) IncreaseMintIndex(m *types.MintModel) int64 {
	db.Model(m).Where("address =?", m.Address).Update("`index`", gorm.Expr("`index` + 1"))
	db.First(m, "address = ?", m.Address)
	return m.Index
}

// TODO: handle not found error
func (db *DB) GetContractInfo(contractAddress string) (*types.Contracts, error) {
	var info types.Contracts
	result := db.First(&info, "address = ?", contractAddress)
	if result.Error != nil && (result.Error == gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("contract %s not registered", contractAddress)
	}

	if result.Error != nil {
		return nil, fmt.Errorf("error getting contract info: %w", result.Error)
	}
	return &info, nil
}

func (db *DB) RegisterContractInfo(address string, typ string, quantum int64) error {
	var contractInfo types.Contracts
	contractInfo.Address = address
	contractInfo.Type = typ
	contractInfo.Quantum = quantum
	result := db.Create(&contractInfo)
	if result.Error != nil {
		return fmt.Errorf("error registering contract info: %w", result.Error)
	}
	return nil
}

func (db *DB) GetOrCreateVault(address string, starkKey string, tokenID int64) int64 {
	var vault types.VaultModel
	result := db.First(&vault, "address = ? and stark_key = ? and token_id = ?", address, starkKey, tokenID)
	if (result.Error != nil) && (result.Error == gorm.ErrRecordNotFound) {
		vault.Address = address
		vault.StarkKey = starkKey
		vault.TokenID = fmt.Sprint(tokenID)
		db.Create(&vault)
	}
	return int64(vault.ID) + BaseVaultID
}

func (db *DB) QueryVault(address string, starkKey string, tokenID int64) (int64, error) {
	var vault types.VaultModel
	result := db.First(&vault, "address = ? and stark_key = ? and token_id = ?", address, starkKey, tokenID)
	if (result.Error != nil) && (result.Error == gorm.ErrRecordNotFound) {
		return 0, result.Error
	}
	return int64(vault.ID) + BaseVaultID, nil
}

func (db *DB) GetNonce(starkKey string) int64 {
	var nonce types.NonceModel
	result := db.First(&nonce, "stark_key = ?", starkKey)
	if (result.Error != nil) && (result.Error == gorm.ErrRecordNotFound) {
		return int64(0)
	}
	return nonce.Nonce
}

func (db *DB) IncreaseNonce(starkKey string) int64 {
	var nonce types.NonceModel

	result := db.First(&nonce, "stark_key = ?", starkKey)
	if (result.Error != nil) && (result.Error == gorm.ErrRecordNotFound) {
		nonce.Nonce = 0
		nonce.StarkKey = starkKey
		db.Create(&nonce)
	}

	db.Model(&nonce).Where("stark_key =?", starkKey).Update("`nonce`", gorm.Expr("`nonce` + 1"))
	db.First(&nonce, "stark_key = ?", starkKey)
	return nonce.Nonce
}

func (db *DB) GetLastEvent(eventName types.EventName) *types.EventLog {
	latestEvent := new(types.EventLog)
	db.Where("event_name = ?", eventName).Last(&latestEvent)
	return latestEvent
}
