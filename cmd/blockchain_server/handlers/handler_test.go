package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/elarsaks/Go-blockchain/pkg/block"
	"github.com/elarsaks/Go-blockchain/pkg/wallet"
)

type testBlockchainServer struct {
	bc     *block.Blockchain
	wallet *wallet.Wallet
}

func (s testBlockchainServer) Port() uint16 {
	return 5001
}

func (s testBlockchainServer) GetWallet() *wallet.Wallet {
	return s.wallet
}

func (s testBlockchainServer) GetBlockchain() *block.Blockchain {
	return s.bc
}

func newTestBlockchainHandler(t *testing.T) (*BlockchainServerHandler, *block.Blockchain, *wallet.Wallet) {
	t.Helper()

	minerWallet, err := wallet.NewWalletWithError()
	if err != nil {
		t.Fatalf("create miner wallet: %v", err)
	}

	bc := block.NewBlockchain(minerWallet.BlockchainAddress(), 5001)
	handler := NewBlockchainServerHandler(testBlockchainServer{
		bc:     bc,
		wallet: minerWallet,
	})

	return handler, bc, minerWallet
}

func TestGetChainReturnsGenesisChain(t *testing.T) {
	handler, _, _ := newTestBlockchainHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/chain", nil)
	rec := httptest.NewRecorder()

	handler.GetChain(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", rec.Code, rec.Body.String())
	}
	if got := rec.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("Content-Type = %q, want application/json", got)
	}

	var response struct {
		Chain []json.RawMessage `json:"chain"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(response.Chain) != 1 {
		t.Fatalf("chain length = %d, want 1", len(response.Chain))
	}
}

func TestGetTransactionsReturnsPoolAndLength(t *testing.T) {
	handler, bc, _ := newTestBlockchainHandler(t)
	if ok, err := bc.AddTransaction(block.MINING_SENDER, "recipient", "reward", 1, nil, nil); err != nil || !ok {
		t.Fatalf("seed transaction pool: ok=%v err=%v", ok, err)
	}

	req := httptest.NewRequest(http.MethodGet, "/transactions", nil)
	rec := httptest.NewRecorder()

	handler.Transactions(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var response struct {
		Transactions []json.RawMessage `json:"transactions"`
		Length       int               `json:"length"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response.Length != 1 {
		t.Fatalf("length = %d, want 1", response.Length)
	}
	if len(response.Transactions) != 1 {
		t.Fatalf("transactions length = %d, want 1", len(response.Transactions))
	}
}

func TestPostTransactionRejectsMissingFieldsWithoutMutatingPool(t *testing.T) {
	handler, bc, _ := newTestBlockchainHandler(t)

	req := httptest.NewRequest(http.MethodPost, "/transactions", strings.NewReader(`{}`))
	rec := httptest.NewRecorder()

	handler.Transactions(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected current status 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var response map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response["message"] != "fail" {
		t.Fatalf("message = %q, want fail", response["message"])
	}
	if got := len(bc.TransactionPool()); got != 0 {
		t.Fatalf("transaction pool length = %d, want 0", got)
	}
}

func TestMinerWalletReturnsServerWallet(t *testing.T) {
	handler, _, minerWallet := newTestBlockchainHandler(t)

	req := httptest.NewRequest(http.MethodPost, "/miner/wallet", nil)
	rec := httptest.NewRecorder()

	handler.MinerWallet(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var response struct {
		BlockchainAddress string `json:"blockchainAddress"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response.BlockchainAddress != minerWallet.BlockchainAddress() {
		t.Fatalf("blockchainAddress = %q, want %q", response.BlockchainAddress, minerWallet.BlockchainAddress())
	}
}
