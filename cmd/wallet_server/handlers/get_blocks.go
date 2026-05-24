package handlers

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/elarsaks/Go-blockchain/pkg/block"
)

// Handler function to get requested blocks
func (h *WalletServerHandler) GetBlocks(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Get the 'amount' query parameter from the URL
	amountStr := req.URL.Query().Get("amount")
	amount, err := strconv.Atoi(amountStr)
	if err != nil || amount <= 0 {
		writeJSONError(w, http.StatusBadRequest, "invalid amount parameter")
		return
	}

	gateway, err := h.gatewayForRequest(req)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	endpoint, err := url.Parse(gateway + "/miner/blocks")
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "invalid miner gateway")
		return
	}
	query := endpoint.Query()
	query.Set("amount", strconv.Itoa(amount))
	endpoint.RawQuery = query.Encode()

	minerReq, err := http.NewRequest(http.MethodGet, endpoint.String(), nil)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "failed to build blocks request")
		return
	}

	// Make a GET request to miner-2's API to fetch blocks
	resp, err := h.client.Do(minerReq)
	if err != nil {
		writeJSONError(w, http.StatusBadGateway, "failed to reach miner")
		return
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		writeJSONError(w, http.StatusBadGateway, "error fetching blocks from miner")
		return
	}

	// Decode the JSON response into a slice of Block
	var blocks []block.Block
	if err := json.NewDecoder(resp.Body).Decode(&blocks); err != nil {
		writeJSONError(w, http.StatusBadGateway, "invalid blocks response")
		return
	}

	// Respond with JSON-encoded blocks
	writeJSON(w, http.StatusOK, blocks)
}
