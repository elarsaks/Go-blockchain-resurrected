package handlers

import (
	"log"
	"net/http"
)

// TODO: Desciption
func (h *BlockchainServerHandler) MinerWallet(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		w.Header().Set("Content-Type", "application/json")
		myWallet := h.server.GetWallet()
		m, _ := myWallet.MarshalJSON()
		_, _ = w.Write(m)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		log.Println("ERROR: Invalid HTTP Method")
	}
}
