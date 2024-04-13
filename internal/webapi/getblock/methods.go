package webapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/egor-denisov/biggest-change/internal/entity"
)

var _defaultMaxRetries = 5
var _defaultTimeBetweenRetries = 500 * time.Millisecond

// Checking validity of url.
func isValidUrl(url string) bool {
	return strings.HasPrefix(url, "https://go.getblock.io/")
}

// Trying making retry requests .
func (w *StatsOfChangingWebAPI) retryRequest(request *http.Request, response interface{}) (err error) {
	for i := 0; i < _defaultMaxRetries; i++ {
		w.Limiter.WaitForAvailability()

		resp, err := w.Client.Do(request)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		err = json.NewDecoder(resp.Body).Decode(response)
		if err == nil {
			// If successful, return response
			break
		}

		// If empty body, trying again
		time.Sleep(_defaultTimeBetweenRetries)
	}

	if err == io.EOF {
		return entity.ErrTooMuchRequestToService
	}

	return err
}

// Building Request Body for eth_getBlockByNumber request.
func getBlockByNumberBuildRequestBody(blockNumber *big.Int) (*bytes.Buffer, error) {
	data := request{
		JsonRPC: "2.0",
		Method:  "eth_getBlockByNumber",
		Params:  []interface{}{int2hex(blockNumber), true},
		Id:      "getblock.io",
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(jsonData), nil
}

// Making request and getting slice of transactions.
func (w *StatsOfChangingWebAPI) getTransactionsByBlockNumber(
	blockNumber *big.Int,
	resCh chan<- []*entity.Transaction,
	errCh chan<- error,
) {
	body, err := getBlockByNumberBuildRequestBody(blockNumber)
	if err != nil {
		errCh <- err
		return
	}

	req, err := http.NewRequest(http.MethodPost, w.URL, body)
	if err != nil {
		errCh <- err
		return
	}

	req.Header.Set("Content-Type", "application/json")

	response := getBlockByNumberResponse{}

	if err := w.retryRequest(req, &response); err != nil {
		errCh <- err
		return
	}

	res := make([]*entity.Transaction, len(response.Result.Transactions))
	// Сonverting the values from hex to *big.Int
	for i, t := range response.Result.Transactions {
		gas, err := hex2int(t.Gas)
		if err != nil {
			errCh <- err
			return
		}

		gasPrice, err := hex2int(t.GasPrice)
		if err != nil {
			errCh <- err
			return
		}

		value, err := hex2int(t.Value)
		if err != nil {
			errCh <- err
			return
		}

		res[i] = &entity.Transaction{
			From:     t.From,
			To:       t.To,
			Gas:      gas,
			GasPrice: gasPrice,
			Value:    value,
		}
	}
	resCh <- res
}

// Building Request Body for eth_blockNumber request.
func blockNumberBuildRequestBody() (*bytes.Buffer, error) {
	data := request{
		JsonRPC: "2.0",
		Method:  "eth_blockNumber",
		Params:  []interface{}{},
		Id:      "getblock.io",
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(jsonData), nil
}

// Making request and getting number of last block.
func (w *StatsOfChangingWebAPI) getCurrentBlockNumber(
	resCh chan<- *big.Int,
	errCh chan<- error,
) {
	body, err := blockNumberBuildRequestBody()
	if err != nil {
		errCh <- err
		return
	}

	req, err := http.NewRequest(http.MethodPost, w.URL, body)
	if err != nil {
		errCh <- err
		return
	}

	req.Header.Set("Content-Type", "application/json")

	w.Limiter.WaitForAvailability()

	resp, err := w.Client.Do(req)
	if err != nil {
		errCh <- err
		return
	}
	defer resp.Body.Close()

	response := blockNumberResponse{}

	if err := w.retryRequest(req, &response); err != nil {
		errCh <- err
		return
	}

	res, err := hex2int(response.Result)
	if err != nil {
		errCh <- err
		return
	}

	resCh <- res
}

func hex2int(s string) (*big.Int, error) {
	i := new(big.Int)
	if s == "" {
		return i, nil
	}

	_, ok := i.SetString(s, 0)
	if !ok {
		return nil, entity.ErrStringIsNotHex
	}

	return i, nil
}

func int2hex(i *big.Int) string {
	return fmt.Sprintf("%#x", i)
}
