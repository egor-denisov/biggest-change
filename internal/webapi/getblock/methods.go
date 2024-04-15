package webapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/egor-denisov/biggest-change/internal/entity"
)

// Checking validity of url.
func isValidUrl(url string) bool {
	return strings.HasPrefix(url, "https://go.getblock.io/")
}

// Trying making retry requests .
func (w *StatsOfChangingWebAPI) retryRequest(
	ctx context.Context,
	body *bytes.Buffer,
	response interface{},
) (err error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, w.url, body)
	if err != nil {
		return fmt.Errorf("StatsOfChangingWebAPI - retryRequest: %w", err)
	}

	request.Header.Set("Content-Type", "application/json")

	for i := 0; i < w.maxRetries; i++ {
		w.limiter.WaitForAvailability()

		// If context is done, return error and decrease limiter counter
		select {
		case <-ctx.Done():
			w.limiter.Rollback()
			return ctx.Err()
		default:
		}

		resp, err := w.client.Do(request)
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
		time.Sleep(w.timeBetweenRetries)
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
	ctx context.Context,
	blockNumber *big.Int,
) ([]*entity.Transaction, error) {
	body, err := getBlockByNumberBuildRequestBody(blockNumber)
	if err != nil {
		return nil,
			fmt.Errorf("StatsOfChangingWebAPI - getTransactionsByBlockNumber - getBlockByNumberBuildRequestBody: %w", err)
	}

	response := getBlockByNumberResponse{}

	if err := w.retryRequest(ctx, body, &response); err != nil {
		return nil,
			fmt.Errorf("StatsOfChangingWebAPI - getTransactionsByBlockNumber - w.retryRequest: %w", err)
	}

	res := make([]*entity.Transaction, len(response.Result.Transactions))
	// Сonverting the values from hex to *big.Int
	for i, t := range response.Result.Transactions {
		gas, err := hex2int(t.Gas)
		if err != nil {
			return nil,
				fmt.Errorf("StatsOfChangingWebAPI - getTransactionsByBlockNumber - hex2int: %w", err)
		}

		gasPrice, err := hex2int(t.GasPrice)
		if err != nil {
			return nil,
				fmt.Errorf("StatsOfChangingWebAPI - getTransactionsByBlockNumber - hex2int: %w", err)
		}

		value, err := hex2int(t.Value)
		if err != nil {
			return nil,
				fmt.Errorf("StatsOfChangingWebAPI - getTransactionsByBlockNumber - hex2int: %w", err)
		}

		res[i] = &entity.Transaction{
			From:     t.From,
			To:       t.To,
			Gas:      gas,
			GasPrice: gasPrice,
			Value:    value,
		}
	}

	return res, nil
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
	ctx context.Context,
) (*big.Int, error) {
	body, err := blockNumberBuildRequestBody()
	if err != nil {
		return nil,
			fmt.Errorf("StatsOfChangingWebAPI - getCurrentBlockNumber - blockNumberBuildRequestBody: %w", err)
	}

	response := blockNumberResponse{}

	if err := w.retryRequest(ctx, body, &response); err != nil {
		return nil,
			fmt.Errorf("StatsOfChangingWebAPI - getCurrentBlockNumber - w.retryRequest: %w", err)
	}

	res, err := hex2int(response.Result)
	if err != nil {
		return nil,
			fmt.Errorf("StatsOfChangingWebAPI - getCurrentBlockNumber - hex2int: %w", err)
	}

	return res, nil
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
