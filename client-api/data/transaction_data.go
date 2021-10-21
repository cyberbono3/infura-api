package data

type TransactionData struct {
	BlockNumber      int `json:"block" validate:"required,gt=0"`
	TransactionIndex int `json:"index" validate:"required,gt=0"`
}

func (td *TransactionData) GetBlockNumber() int {
	return td.BlockNumber
}

func (td *TransactionData) GetTransactionIndex() int {
	return td.TransactionIndex
}
