package data

type BlockData struct {
	BlockNumber   int  `json:"block" validate:"required,gt=0"`
	ShowTransFlag bool `json:"show"`
}

func (bd *BlockData) GetBlockNumber() int {
	return bd.BlockNumber

}
func (bd *BlockData) GetShowTransFlag() bool {
	return bd.ShowTransFlag
}
