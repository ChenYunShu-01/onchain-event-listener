package types

import (
	"math/big"
)

type L2DepositRequest struct {
	Type     string   `json:"type"`
	VaultId  *big.Int `json:"vault_id"`
	Amount   string   `json:"amount"`
	TokenId  string   `json:"token_id"`
	StarkKey string   `json:"stark_key"`
}
