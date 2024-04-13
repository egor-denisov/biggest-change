package entity

import "math/big"

// @Description Транзакция .
type Transaction struct {
	From     string   `json:"from"`
	Gas      *big.Int `json:"gas"`
	GasPrice *big.Int `json:"gasPrice"`
	To       string   `json:"to"`
	Value    *big.Int `json:"value"`
}
