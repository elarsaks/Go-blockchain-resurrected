package block

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestBlockUnmarshalJSONRejectsInvalidPreviousHash(t *testing.T) {
	data := []byte(`{
		"timestamp": 1,
		"nonce": 2,
		"previousHash": "abc",
		"transactions": []
	}`)

	var block Block
	err := json.Unmarshal(data, &block)
	if err == nil {
		t.Fatal("expected invalid previousHash to return an error")
	}
	if !strings.Contains(err.Error(), "decode previousHash") &&
		!strings.Contains(err.Error(), "previousHash must be") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBlockUnmarshalJSONRejectsMissingFields(t *testing.T) {
	var block Block
	err := json.Unmarshal([]byte(`{"nonce":1,"previousHash":""}`), &block)
	if err == nil {
		t.Fatal("expected missing timestamp to return an error")
	}
}

func TestBlockJSONRoundTripAcceptsNilTransactions(t *testing.T) {
	block := NewBlock(0, [32]byte{}, nil)

	data, err := json.Marshal(block)
	if err != nil {
		t.Fatalf("Marshal block: %v", err)
	}
	if strings.Contains(string(data), `"transactions":null`) {
		t.Fatalf("block JSON should encode empty transactions as [], got %s", data)
	}

	var decoded Block
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal block: %v", err)
	}
	if decoded.Transactions() == nil {
		t.Fatal("decoded transactions should be an empty slice, got nil")
	}
}

func TestBlockUnmarshalJSONAcceptsNullTransactions(t *testing.T) {
	data := []byte(`{
		"timestamp": 1,
		"nonce": 2,
		"previousHash": "0000000000000000000000000000000000000000000000000000000000000000",
		"transactions": null
	}`)

	var block Block
	if err := json.Unmarshal(data, &block); err != nil {
		t.Fatalf("Unmarshal block with null transactions: %v", err)
	}
	if block.Transactions() == nil {
		t.Fatal("transactions should be an empty slice, got nil")
	}
}

func TestTransactionUnmarshalJSONRejectsMissingFields(t *testing.T) {
	var transaction Transaction
	err := json.Unmarshal([]byte(`{"message":"hello"}`), &transaction)
	if err == nil {
		t.Fatal("expected missing transaction fields to return an error")
	}
}
