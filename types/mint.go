package types

import (
	"math/big"

	"gorm.io/gorm"
)

// MintRequest is the request body for the /api/v1/mint endpoint.
type MintRequest struct {
	Address  string `json:"address"`
	Amount   int64  `json:"amount"`
	StarkKey string `json:"stark_key"`
}

func (m *MintRequest) Verify() error {
	if m.Address == "" {
		return ErrAddressMissing
	}

	if m.StarkKey == "" {
		return ErrStarkKeyMissing
	}

	// TODO: fix rycle import
	// if contracts.IsERC20(m.ContractAddress) && m.Amount <= 0 {
	// 	return ErrAmountInvalid
	// }

	return nil
}

type L2MintData struct {
	VaultId  *big.Int `json:"vault_id"`
	Amount   string   `json:"amount"`
	TokenId  string   `json:"token_id"`
	StarkKey string   `json:"stark_key"`
	Type     string   `json:"type"`
}

type MintModel struct {
	gorm.Model
	Address string
	Index   int64
}

type MintResponse struct {
	TransactionID int64 `json:"transaction_id"`
}
