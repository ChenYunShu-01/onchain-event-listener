package types

import "gorm.io/gorm"

// VaultRequest is the request body for the /api/v1/vault endpoint.
type NonceRequest struct {
	StarkKey string `form:"stark_key"`
}

func (r *NonceRequest) Verify() error {

	if r.StarkKey == "" {
		return ErrStarkKeyMissing
	}

	return nil
}

type NonceResponse struct {
	Nonce int64 `json:"nonce"`
}

// TODO: extend vault model support all kinds of token
type NonceModel struct {
	gorm.Model
	StarkKey string `json:"stark_key"`
	Nonce    int64  `json:"nonce"`
}
