package types

type TransferRequest struct {
	Signature           Signature `json:"signature"`
	Nonce               int64     `json:"nonce"`
	Amount              int64     `json:"amount"`
	TokenID             string    `json:"token_id"`
	StarkKey            string    `json:"stark_key"`
	VaultID             int64     `json:"vault_id"`
	Receiver            string    `json:"receiver"`
	ReceiverVaultID     int64     `json:"receiver_vault_id"`
	ExpirationTimestamp int64     `json:"expiration_timestamp"`
}

func (r *TransferRequest) Verify() error {
	return nil
}

type L2TransferRequest struct {
	Type                string    `json:"type"`
	Signature           Signature `json:"signature"`
	Nonce               int64     `json:"nonce"`
	Amount              int64     `json:"amount"`
	TokenID             string    `json:"token"`
	StarkKey            string    `json:"sender_public_key"`
	VaultID             int64     `json:"sender_vault_id"`
	Receiver            string    `json:"receiver_public_key"`
	ReceiverVaultID     int64     `json:"receiver_vault_id"`
	ExpirationTimestamp int64     `json:"expiration_timestamp"`
}

type TransferResponse struct {
	TransactionID int64 `json:"transaction_id"`
}
