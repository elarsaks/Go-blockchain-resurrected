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
	err := json.Unmarshal([]byte(`{"nonce":1,"previousHash":"","transactions":[]}`), &block)
	if err == nil {
		t.Fatal("expected missing timestamp to return an error")
	}
}

func TestTransactionUnmarshalJSONRejectsMissingFields(t *testing.T) {
	var transaction Transaction
	err := json.Unmarshal([]byte(`{"message":"hello"}`), &transaction)
	if err == nil {
		t.Fatal("expected missing transaction fields to return an error")
	}
}
