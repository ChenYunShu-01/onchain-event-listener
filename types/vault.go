package types

import "gorm.io/gorm"

// VaultRequest is the request body for the /api/v1/vault endpoint.
type VaultRequest struct {
	Address  string `form:"address"`
	TokenID  string `form:"token_id"`
	StarkKey string `form:"stark_key"`
}

func (v *VaultRequest) Verify() error {
	if v.Address == "" {
		return ErrAddressMissing
	}

	if v.StarkKey == "" {
		return ErrStarkKeyMissing
	}

	// if contracts.IsERC20(m.ContractAddress) && m.Amount <= 0 {
	// 	return ErrAmountInvalid
	// }

	return nil
}

type VaultResponse struct {
	VaultId int64 `json:"vault_id"`
}

// TODO: extend vault model support all kinds of token
type VaultModel struct {
	gorm.Model
	Address  string    `json:"address"`
	Type     TokenType `json:"type"`
	StarkKey string    `json:"stark_key"`
	TokenID  string    `json:"token_id"`
}
