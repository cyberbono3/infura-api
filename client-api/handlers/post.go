package handlers

import (
	"net/http"

	"github.com/cyberbono3/infura-coding/client-api/data"
)

func (ch *ClientHandler) GetTransaction(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	td := r.Context().Value(KeyTransaction{}).(*data.TransactionData)

	ch.l.Debug("Fetching transaction data: %#v\n", td.GetBlockNumber())

	tx, err := ch.c.EthGetTransactionByBlockNumberAndIndex(td.GetBlockNumber(), td.GetTransactionIndex())
	if err != nil {
		ch.l.Error("fetching transaction", "error", err)
		rw.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	err = data.ToJSON(tx, rw)
	if err != nil {
		ch.l.Error("[ERROR] serializing transaction", "error", err)
	}

}

func (ch *ClientHandler) GetBlock(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	bd := r.Context().Value(KeyBlock{}).(*data.BlockData)

	blockNumber := bd.GetBlockNumber()

	block, err := ch.c.EthGetBlockByNumber(blockNumber, bd.GetShowTransFlag())
	if err != nil {
		ch.l.Error("Unable to fetch block", "error", err)

		rw.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	err = data.ToJSON(block, rw)
	if err != nil {
		ch.l.Error("Unable to serializing block", err)
	}
}
