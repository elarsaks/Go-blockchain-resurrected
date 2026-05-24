package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

// GetMinerWallet is a handler that:
// 1. Gets the 'miner' query parameter from the URL
// 2. Sets the wallet server's gateway to blockhain to the requested miner
// 3. Makes a POST request to the miner's API to fetch the wallet
// 4. Checks the response status code
// 5. Decodes the JSON response into a struct or a map
// 6. Encodes the wallet data to JSON and writes it to the response
func (h *WalletServerHandler) GetMinerWallet(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Get the 'miner' query parameter from the URL
	minerID := req.URL.Query().Get("miner_id")
	if minerID == "" {
		minerID = "1"
	}

	gateway, err := h.server.MinerGateway(minerID)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	minerReq, err := http.NewRequest(http.MethodPost, gateway+"/miner/wallet", nil)
	if err != nil {
		log.Printf("ERROR: Failed to build miner wallet request: %v", err)
		writeJSONError(w, http.StatusInternalServerError, "failed to fetch miner wallet")
		return
	}
	minerReq.Header.Set("Content-Type", "application/json")

	resp, err := h.client.Do(minerReq)

	if err != nil {
		log.Printf("ERROR: %v", err)
		writeJSONError(w, http.StatusBadGateway, "failed to reach miner")
		return
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		log.Printf("ERROR: Error fetching wallet from %s", minerID)
		writeJSONError(w, http.StatusBadGateway, "failed to fetch miner wallet")
		return
	}

	// Decode the JSON response into a struct or a map
	var walletData map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&walletData)
	if err != nil {
		log.Printf("ERROR: Error decoding wallet response: %v", err)
		writeJSONError(w, http.StatusBadGateway, "invalid miner wallet response")
		return
	}

	writeJSON(w, http.StatusOK, walletData)
}
