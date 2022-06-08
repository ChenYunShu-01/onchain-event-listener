package types

import (
	"gorm.io/gorm"
)

type Contracts struct {
	gorm.Model

	Name        string    `json:"name"`
	Address     string    `json:"address"`
	Symbol      string    `json:"symbol"`
	Decimals    int       `json:"decimals"`     // For l1
	TotalSupply int64     `json:"total_supply"` // For L1
	Type        TokenType `json:"type"`
	Quantum     int64     `json:"quantum"` // For starkex
	Count       int64     `json:"count"`   // For starkex
}
