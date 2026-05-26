package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

// TODO: Describe the purpose of this function
// Get the last 10 blocks of the BlockchainServer
func (h *BlockchainServerHandler) GetBlocks(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		bc := h.server.GetBlockchain()
		m, _ := json.Marshal(bc.GetBlocks(10))
		_, _ = w.Write(m)
	default:
		log.Printf("ERROR: Invalid HTTP Method")
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
