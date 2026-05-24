package handlers

import (
	"encoding/json"
	"net/http"
	"time"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type WalletServer interface {
	Port() uint16
	Gateway() string
	MinerGateway(minerID string) (string, error)
}

type WalletServerHandler struct {
	server WalletServer
	client HTTPClient
}

func NewWalletServerHandler(s WalletServer) *WalletServerHandler {
	return NewWalletServerHandlerWithClient(s, &http.Client{Timeout: 10 * time.Second})
}

func NewWalletServerHandlerWithClient(s WalletServer, client HTTPClient) *WalletServerHandler {
	return &WalletServerHandler{server: s, client: client}
}

func (h *WalletServerHandler) gatewayForRequest(req *http.Request) (string, error) {
	minerID := req.URL.Query().Get("miner_id")
	if minerID == "" {
		return h.server.Gateway(), nil
	}
	return h.server.MinerGateway(minerID)
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
