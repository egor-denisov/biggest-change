package webapi

type request struct {
	JsonRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Id      string        `json:"id"`
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
