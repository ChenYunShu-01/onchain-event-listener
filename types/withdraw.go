package types

type WithdrawRequest struct {
	Address  string `json:"address"`
	TokenID  string `json:"token_id"`
	StarkKey string `json:"stark_key"`
	Amount   int64  `json:"amount"`
}

func (w *WithdrawRequest) Verify() error {
	return nil
}

type L2WithdrawRequest struct {
	TokenID  string `json:"token_id"`
	Amount   string `json:"amount"`
	StarkKey string `json:"stark_key"`
	VaultID  string `json:"vault_id"`
	Type     string `json:"type"`
}

type WithdrawResponse struct {
	TransactionID int64 `json:"transaction_id"`
}
