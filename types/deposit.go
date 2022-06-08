package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type L2DepositRequest struct {
	Type     string   `json:"type"`
	VaultId  *big.Int `json:"vault_id"`
	Amount   string   `json:"amount"`
	TokenId  string   `json:"token_id"`
	StarkKey string   `json:"stark_key"`
}

type DepositLog struct {
	DepositorEthKey    common.Address
	StarkKey           *big.Int
	VaultId            *big.Int
	AssetType          *big.Int
	NonQuantizedAmount *big.Int
	QuantizedAmount    *big.Int
	Raw                types.Log // Blockchain specific contextual infos
}
