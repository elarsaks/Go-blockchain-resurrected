package handlers

import (
	"log"
	"net/http"
)

// TODO: Describe the purpose of this function

// Get the gateway of the BlockchainServer
func (h *BlockchainServerHandler) GetChain(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		bc := h.server.GetBlockchain()
		m, _ := bc.MarshalJSON()
		_, _ = w.Write(m)
	default:
		log.Printf("ERROR: Invalid HTTP Method")
		w.WriteHeader(http.StatusMethodNotAllowed)

	}
}
