package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	"github.com/elarsaks/Go-blockchain/pkg/block"
)

func (h *WalletServerHandler) GetWalletBalance(w http.ResponseWriter, req *http.Request) {
	// Check if the HTTP method is GET
	if req.Method != http.MethodGet {
		log.Println("ERROR: Invalid HTTP Method")
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Extract the blockchain address from the URL query parameters
	blockchainAddress := req.URL.Query().Get("blockchainAddress")
	if blockchainAddress == "" {
		writeJSONError(w, http.StatusBadRequest, "blockchainAddress is required")
		return
	}

	gateway, err := h.gatewayForRequest(req)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Construct the endpoint URL for the blockchain API
	endpoint, err := url.Parse(gateway + "/balance")
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "invalid miner gateway")
		return
	}
	query := endpoint.Query()
	query.Set("blockchainAddress", blockchainAddress)
	endpoint.RawQuery = query.Encode()

	minerReq, err := http.NewRequest(http.MethodGet, endpoint.String(), nil)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "failed to build balance request")
		return
	}

	// Send a GET request to the blockchain API
	resp, err := h.client.Do(minerReq)
	if err != nil {
		log.Printf("ERROR: %v", err)
		writeJSONError(w, http.StatusBadGateway, "failed to reach miner")
		return
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode == http.StatusOK {
		// Decode the response JSON into the existing response struct
		br := &block.BalanceResponse{}
		err := json.NewDecoder(resp.Body).Decode(br)
		if err != nil {
			log.Printf("ERROR: %v", err)
			writeJSONError(w, http.StatusBadGateway, "invalid balance response")
			return
		}

		writeJSON(w, http.StatusOK, br)
	} else {
		// Create a new response struct for the failure case
		failureResponse := &block.BalanceResponse{
			Error: "Failed to get wallet balance",
		}
		writeJSON(w, resp.StatusCode, failureResponse)
	}
}
