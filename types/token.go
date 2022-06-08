package types

type TokenType = string

const (
	ETH     = TokenType("ETH")
	ERC20   = TokenType("ERC20")
	ERC721  = TokenType("ERC721")
	ERC20M  = TokenType("ERC20M")
	ERC721M = TokenType("ERC721M")
)

type TokenRegisterRequest struct {
	Type    TokenType `json:"type"`
	Address string    `json:"address"`
}

func (t *TokenRegisterRequest) Verify() error {
	return nil
}

type TokenRegisterResponse struct {
	TxHash string `json:"tx_hash"`
}
