package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/elarsaks/Go-blockchain/pkg/wallet"
)

// Get User wallet
func (h *WalletServerHandler) GetUserWallet(w http.ResponseWriter, req *http.Request) {

	if req.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	gateway, err := h.gatewayForRequest(req)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	userWallet, err := wallet.NewWalletWithError()
	if err != nil {
		log.Println("ERROR: Failed to create wallet:", err)
		writeJSONError(w, http.StatusInternalServerError, "failed to create wallet")
		return
	}

	// Create a payload containing the userWallet's blockchain address
	payload := struct {
		BlockchainAddress string `json:"blockchainAddress"`
	}{
		BlockchainAddress: userWallet.BlockchainAddress(),
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Println("ERROR: Failed to marshal payload:", err)
		writeJSONError(w, http.StatusInternalServerError, "failed to encode wallet registration")
		return
	}

	// Register the userWallet on the blockchain
	minerReq, err := http.NewRequest(http.MethodPost, gateway+"/wallet/register", bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Printf("ERROR: Failed to build wallet registration request: %v", err)
		writeJSONError(w, http.StatusInternalServerError, "failed to register wallet")
		return
	}
	minerReq.Header.Set("Content-Type", "application/json")

	resp, err := h.client.Do(minerReq)
	if err != nil {
		log.Printf("ERROR: Failed to register wallet: %v", err)
		writeJSONError(w, http.StatusBadGateway, "failed to register wallet")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("ERROR: Failed to register wallet")
		writeJSONError(w, http.StatusBadGateway, "failed to register wallet")
		return
	}

	writeJSON(w, http.StatusOK, userWallet)
}
