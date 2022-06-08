package starkex

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/reddio-com/red-adapter/types"
)

var (
	DEFAULT_GATEWAY_URL = "https://gw.playground-v2.starkex.co/v2/gateway"
	ADD_TRANSACTION_URL = DEFAULT_GATEWAY_URL + "/add_transaction"

	// TODO: remove this when online
	FIRST_UNUSE_TXID_URL = "https://gw.playground-v2.starkex.co/testing/get_first_unused_tx_id"

	TRANSACTION_PENDING = "TRANSACTION_PENDING"
)

type Transaction struct {
	TX            interface{} `json:"tx"`
	TransactionID int64       `json:"tx_id"`
}

type Result struct {
	Code string `json:"code"`
}

func Deposit(r *types.L2DepositRequest) (int64, error) {
	r.Type = "DepositRequest"
	tx_id := GetTransactionID()
	t := Transaction{TX: r, TransactionID: tx_id}
	body, _ := json.Marshal(t)
	err := request(ADD_TRANSACTION_URL, body)
	return tx_id, err
}

func GetTransactionID() int64 {
	resp, err := http.Get(FIRST_UNUSE_TXID_URL)
	if err != nil {
		return -1
	}

	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return -1
	}

	n, err := strconv.ParseInt(string(b), 10, 64)
	if err != nil {
		return -1
	}

	return n

}

// TODO - add more error handling
func request(url string, body []byte) error {
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)

	if err != nil {
		fmt.Println(err)
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New(string(b))
		// return errors.New(resp.Status)
	}

	var r Result
	err = json.Unmarshal(b, &r)
	if err != nil {
		return err
	}

	if r.Code != TRANSACTION_PENDING {
		return errors.New(r.Code)
	}

	return nil
}
