package types

var (
	OK     = "OK"
	FAILED = "FAILED"
)

type Response struct {
	Status string      `json:"status"`
	Error  string      `json:"error"`
	Data   interface{} `json:"data"`
}
