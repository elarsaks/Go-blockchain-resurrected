package handlers

import (
	"log"
	"net/http"

	"github.com/elarsaks/Go-blockchain/pkg/utils"
)

// TODO: Desciption

// Mine the Block in the BlockchainServer
func (h *BlockchainServerHandler) Mine(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		bc := h.server.GetBlockchain()
		isMined := bc.Mining()

		var m []byte
		if !isMined {
			w.WriteHeader(http.StatusBadRequest)
			m = utils.JsonStatus("fail")
		} else {
			m = utils.JsonStatus("success")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(m)
	default:
		log.Println("ERROR: Invalid HTTP Method")
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
