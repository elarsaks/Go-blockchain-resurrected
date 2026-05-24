package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/elarsaks/Go-blockchain/pkg/block"
	"github.com/elarsaks/Go-blockchain/pkg/wallet"
)

type testWalletServer struct {
	defaultGateway string
	gateways       map[string]string
}

func (s testWalletServer) Port() uint16 {
	return 5000
}

func (s testWalletServer) Gateway() string {
	return s.defaultGateway
}

func (s testWalletServer) MinerGateway(minerID string) (string, error) {
	gateway, ok := s.gateways[minerID]
	if !ok {
		return "", errInvalidMinerID(minerID)
	}
	return gateway, nil
}

func errInvalidMinerID(minerID string) error {
	return &invalidMinerIDError{minerID: minerID}
}

type invalidMinerIDError struct {
	minerID string
}

func (e *invalidMinerIDError) Error() string {
	return "invalid miner_id " + e.minerID
}

func TestGetWalletBalanceUsesRequestedMinerGateway(t *testing.T) {
	miner := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/balance" {
			t.Fatalf("unexpected miner path: %s", req.URL.Path)
		}
		if got := req.URL.Query().Get("blockchainAddress"); got != "wallet address/1" {
			t.Fatalf("unexpected blockchainAddress: %q", got)
		}
		json.NewEncoder(w).Encode(&block.BalanceResponse{Balance: 3})
	}))
	defer miner.Close()

	handler := NewWalletServerHandler(testWalletServer{
		defaultGateway: "http://default.invalid",
		gateways:       map[string]string{"2": miner.URL},
	})

	query := url.Values{}
	query.Set("blockchainAddress", "wallet address/1")
	query.Set("miner_id", "2")
	req := httptest.NewRequest(http.MethodGet, "/wallet/balance?"+query.Encode(), nil)
	rec := httptest.NewRecorder()

	handler.GetWalletBalance(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var response block.BalanceResponse
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response.Balance != 3 {
		t.Fatalf("expected balance 3, got %v", response.Balance)
	}
}

func TestGetMinerWalletRejectsInvalidMinerID(t *testing.T) {
	handler := NewWalletServerHandler(testWalletServer{
		defaultGateway: "http://default.invalid",
		gateways:       map[string]string{"1": "http://miner.invalid"},
	})

	req := httptest.NewRequest(http.MethodPost, "/miner/wallet?miner_id=4", nil)
	rec := httptest.NewRecorder()

	handler.GetMinerWallet(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
}

func TestCreateTransactionSignsAndForwardsToRequestedMiner(t *testing.T) {
	sender, err := wallet.NewWalletWithError()
	if err != nil {
		t.Fatalf("create sender wallet: %v", err)
	}
	recipient, err := wallet.NewWalletWithError()
	if err != nil {
		t.Fatalf("create recipient wallet: %v", err)
	}

	var forwarded block.TransactionRequest
	miner := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/transactions" {
			t.Fatalf("unexpected miner path: %s", req.URL.Path)
		}
		if req.Method != http.MethodPost {
			t.Fatalf("unexpected miner method: %s", req.Method)
		}
		if err := json.NewDecoder(req.Body).Decode(&forwarded); err != nil {
			t.Fatalf("decode forwarded transaction: %v", err)
		}
		w.WriteHeader(http.StatusCreated)
	}))
	defer miner.Close()

	handler := NewWalletServerHandler(testWalletServer{
		defaultGateway: "http://default.invalid",
		gateways:       map[string]string{"3": miner.URL},
	})

	payload := `{
		"message":"USER TRANSACTION",
		"recipientBlockchainAddress":"` + recipient.BlockchainAddress() + `",
		"senderBlockchainAddress":"` + sender.BlockchainAddress() + `",
		"senderPrivateKey":"` + sender.PrivateKeyStr() + `",
		"senderPublicKey":"` + sender.PublicKeyStr() + `",
		"value":"1"
	}`
	req := httptest.NewRequest(http.MethodPost, "/transaction?miner_id=3", strings.NewReader(payload))
	rec := httptest.NewRecorder()

	handler.CreateTransaction(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", rec.Code, rec.Body.String())
	}
	if forwarded.Signature == nil || *forwarded.Signature == "" {
		t.Fatal("expected forwarded transaction to include signature")
	}
	if forwarded.Value == nil || *forwarded.Value != 1 {
		t.Fatalf("expected forwarded value 1, got %v", forwarded.Value)
	}
}

func TestCreateTransactionRejectsInvalidSenderKey(t *testing.T) {
	handler := NewWalletServerHandler(testWalletServer{
		defaultGateway: "http://default.invalid",
		gateways:       map[string]string{"1": "http://miner.invalid"},
	})

	payload := `{
		"message":"USER TRANSACTION",
		"recipientBlockchainAddress":"recipient",
		"senderBlockchainAddress":"sender",
		"senderPrivateKey":"not-hex",
		"senderPublicKey":"not-hex",
		"value":"1"
	}`
	req := httptest.NewRequest(http.MethodPost, "/transaction?miner_id=1", strings.NewReader(payload))
	rec := httptest.NewRecorder()

	handler.CreateTransaction(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d: %s", rec.Code, rec.Body.String())
	}
}
