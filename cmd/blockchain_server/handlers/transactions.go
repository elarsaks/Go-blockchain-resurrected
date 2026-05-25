package handlers

import (
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"

	"github.com/elarsaks/Go-blockchain/pkg/block"
	"github.com/elarsaks/Go-blockchain/pkg/utils"
)

func (h *BlockchainServerHandler) HandleGetTransaction(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	bc := h.server.GetBlockchain()

	transactions := bc.TransactionPool()

	_ = json.NewEncoder(w).Encode(struct {
		Transactions []*block.Transaction `json:"transactions"`
		Length       int                  `json:"length"`
	}{
		Transactions: transactions,
		Length:       len(transactions),
	})
}

func (h *BlockchainServerHandler) HandlePostTransaction(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(req.Body)
	var t block.TransactionRequest
	err := decoder.Decode(&t)

	if err != nil {
		log.Printf("ERROR: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(utils.JsonStatus("fail"))
		return
	}

	if !t.Validate() {
		log.Println("ERROR: missing field(s)")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(utils.JsonStatus("fail"))
		return
	}

	if !isHexString(*t.SenderPublicKey, 128) || !isHexString(*t.Signature, 128) {
		log.Println("ERROR: malformed public key or signature")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(utils.JsonStatus("fail"))
		return
	}

	publicKey := utils.PublicKeyFromString(*t.SenderPublicKey)
	signature := utils.SignatureFromString(*t.Signature)

	bc := h.server.GetBlockchain()

	isCreated, err := bc.CreateTransaction(*t.SenderBlockchainAddress,
		*t.RecipientBlockchainAddress, *t.Message, *t.Value, publicKey, signature)

	var m []byte
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errMsg := struct {
			Status  string `json:"status"`
			Message string `json:"message"`
		}{
			Status:  "fail",
			Message: err.Error(),
		}

		m, _ = json.Marshal(errMsg)
	} else if !isCreated {
		w.WriteHeader(http.StatusBadRequest)
		m = utils.JsonStatus("fail")
	} else {
		w.WriteHeader(http.StatusCreated)
		m = utils.JsonStatus("success")
	}

	_, _ = w.Write(m)
}

func (h *BlockchainServerHandler) HandlePutTransaction(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(req.Body)
	var t block.TransactionRequest
	err := decoder.Decode(&t)

	if err != nil {
		log.Printf("ERROR: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(utils.JsonStatus("fail"))
		return
	}

	if !t.Validate() {
		log.Println("ERROR: missing field(s)")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(utils.JsonStatus("fail"))
		return
	}

	if !isHexString(*t.SenderPublicKey, 128) || !isHexString(*t.Signature, 128) {
		log.Println("ERROR: malformed public key or signature")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(utils.JsonStatus("fail"))
		return
	}

	publicKey := utils.PublicKeyFromString(*t.SenderPublicKey)
	signature := utils.SignatureFromString(*t.Signature)

	bc := h.server.GetBlockchain()

	isUpdated, err := bc.AddTransaction(*t.SenderBlockchainAddress,
		*t.RecipientBlockchainAddress,
		*t.Message,
		*t.Value,
		publicKey,
		signature)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "fail", "error": err.Error()})
		return
	}

	var m []byte
	if !isUpdated {
		w.WriteHeader(http.StatusBadRequest)
		m = utils.JsonStatus("fail")
	} else {
		m = utils.JsonStatus("success")
	}

	_, _ = w.Write(m)
}

func (h *BlockchainServerHandler) HandleDeleteTransaction(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	bc := h.server.GetBlockchain()

	bc.ClearTransactionPool()

	_, _ = w.Write(utils.JsonStatus("success"))
}

func (h *BlockchainServerHandler) Transactions(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		h.HandleGetTransaction(w, req)
	case http.MethodPost:
		h.HandlePostTransaction(w, req)
	case http.MethodPut:
		h.HandlePutTransaction(w, req)
	case http.MethodDelete:
		h.HandleDeleteTransaction(w, req)
	default:
		log.Println("ERROR: Invalid HTTP Method")
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func isHexString(s string, expectedLength int) bool {
	if len(s) != expectedLength {
		return false
	}
	_, err := hex.DecodeString(s)
	return err == nil
}
