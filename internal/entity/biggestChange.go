package entity

// @Description Наибольшее изменение .
type BiggestChange struct {
	Address       string `json:"address"`
	Amount        string `json:"amount"`
	LastBlock     string `json:"lastBlock"`
	CountOfBlocks int64  `json:"countOfBlocks"`
	IsRecieved    bool   `json:"isRecieved"`
}
