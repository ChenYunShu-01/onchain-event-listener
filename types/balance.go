package types

type BalanceResponse struct {
	Address  string  `json:"address"`
	Balance  float64 `json:"balance"`
	Currency string  `json:"currency"`
	Symbol   string  `json:"symbol"`
	Type     string  `json:"type"`
	Website  string  `json:"website"`
	Logo     string  `json:"logo"`
}
