package webapi

type request struct {
	JSONRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	ID      string        `json:"id"`
	Params  []interface{} `json:"params"`
}

type transactionResponse struct {
	From     string `json:"from"`
	Gas      string `json:"gas"`
	GasPrice string `json:"gasPrice"`
	To       string `json:"to"`
	Value    string `json:"value"`
}

type getBlockByNumberResponse struct {
	Result struct {
		Transactions []*transactionResponse `json:"transactions"`
	} `json:"result"`
}

type blockNumberResponse struct {
	Result string `json:"result"`
}
