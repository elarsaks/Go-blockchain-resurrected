package handlers

import (
	"log"
	"net/http"

	"github.com/elarsaks/Go-blockchain/pkg/utils"
)

// TODO: Desciption

// Start the mining process in the BlockchainServer
func (h *BlockchainServerHandler) StartMine(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		bc := h.server.GetBlockchain()
		bc.Mining()

		m := utils.JsonStatus("success")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(m)
	default:
		log.Println("ERROR: Invalid HTTP Method")
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
