package handlers

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/elarsaks/Go-blockchain/pkg/block"
	"github.com/elarsaks/Go-blockchain/pkg/utils"
	"github.com/elarsaks/Go-blockchain/pkg/wallet"
)

// Create a new transaction
func (h *WalletServerHandler) CreateTransaction(w http.ResponseWriter, req *http.Request) {
	//* NOTE: We are not just passing request to miner, because we need to sign the transaction
	// Switching on the HTTP method
	switch req.Method {
	case http.MethodPost:

		// Decoding the body of the request into a TransactionRequest object
		decoder := json.NewDecoder(req.Body)
		var t wallet.TransactionRequest
		err := decoder.Decode(&t)

		// If there was an error decoding the request, log the error and send a fail response
		if err != nil {
			log.Printf("ERROR: %v", err)
			writeJSONError(w, http.StatusBadRequest, "invalid transaction request")
			return
		}

		// Validate the transaction request, send a fail response if validation fails
		if !t.Validate() {
			log.Println("ERROR: missing field(s)")
			writeJSONError(w, http.StatusBadRequest, "missing required transaction field")
			return
		}

		gateway, err := h.gatewayForRequest(req)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		if !isValidPublicKeyString(*t.SenderPublicKey) || !isValidPrivateKeyString(*t.SenderPrivateKey) {
			writeJSONError(w, http.StatusBadRequest, "invalid sender key")
			return
		}

		// Convert the sender's public and private keys from strings to their appropriate types
		publicKey := utils.PublicKeyFromString(*t.SenderPublicKey)
		privateKey := utils.PrivateKeyFromString(*t.SenderPrivateKey, publicKey)

		// Parse the value from the request, handle error if the value is not a valid float
		value, err := strconv.ParseFloat(*t.Value, 32)
		if err != nil || value <= 0 {
			log.Println("ERROR: parse error")
			writeJSONError(w, http.StatusBadRequest, "transaction value must be positive")
			return
		}
		value32 := float32(value)

		// Create a new Transaction object
		transaction := wallet.NewTransaction(
			*t.Message,
			*t.RecipientBlockchainAddress,
			*t.SenderBlockchainAddress,
			privateKey,
			publicKey,
			value32)

		// Generate a signature for the transaction
		signature, err := transaction.GenerateSignatureWithError()
		if err != nil {
			log.Printf("ERROR: Failed to sign transaction: %v", err)
			writeJSONError(w, http.StatusBadRequest, "failed to sign transaction")
			return
		}
		signatureStr := signature.String()

		// Create a new TransactionRequest object that will be sent to the miner
		bt := &block.TransactionRequest{
			Message:                    t.Message,
			RecipientBlockchainAddress: t.RecipientBlockchainAddress,
			SenderBlockchainAddress:    t.SenderBlockchainAddress,
			SenderPublicKey:            t.SenderPublicKey,
			Signature:                  &signatureStr,
			Value:                      &value32,
		}

		// Serialize the TransactionRequest object into JSON
		m, err := json.Marshal(bt)
		if err != nil {
			log.Printf("ERROR: Failed to encode transaction: %v", err)
			writeJSONError(w, http.StatusInternalServerError, "failed to encode transaction")
			return
		}
		buf := bytes.NewBuffer(m)

		// Make a POST request to the miner's API to create a new transaction
		minerReq, err := http.NewRequest(http.MethodPost, gateway+"/transactions", buf)
		if err != nil {
			log.Printf("ERROR: Failed to build POST request: %v", err)
			writeJSONError(w, http.StatusInternalServerError, "failed to create transaction")
			return
		}
		minerReq.Header.Set("Content-Type", "application/json")

		resp, err := h.client.Do(minerReq)

		// Check if there was an error while making the POST request
		if err != nil {
			// Log the error message
			log.Printf("ERROR: Failed to make POST request: %v", err)

			// Pass the error message to the client
			writeJSONError(w, http.StatusBadGateway, "failed to reach miner")
			return
		}
		defer resp.Body.Close()

		// Read the response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			// Log the error message
			log.Printf("ERROR: Failed to read response body: %v", err)

			// Pass the error message to the client
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Check the response status code and send a success response if it was 201
		if resp.StatusCode == 201 {
			writeJSON(w, http.StatusOK, map[string]string{"message": "success"})
			return
		}

		// If the status code was not 201, send the response body (which contains the error message) to the client
		w.WriteHeader(resp.StatusCode)
		io.WriteString(w, string(body))

	// If the HTTP method is not POST, send a 400 response and log an error message
	default:
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		log.Println("ERROR: Invalid HTTP Method")
	}
}

func isValidPublicKeyString(value string) bool {
	return len(value) == 128 && isHexString(value)
}

func isValidPrivateKeyString(value string) bool {
	return value != "" && isHexString(value)
}

func isHexString(value string) bool {
	if len(value)%2 != 0 {
		value = "0" + value
	}
	_, err := hex.DecodeString(value)
	return err == nil
}
