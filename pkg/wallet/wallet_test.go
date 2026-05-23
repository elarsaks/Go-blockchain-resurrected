package wallet

import (
	"encoding/json"
	"testing"
)

func TestNewTransactionMapsFields(t *testing.T) {
	transaction := NewTransaction("hello", "recipient", "sender", nil, nil, 2)

	if transaction.message != "hello" {
		t.Fatalf("message = %q, want hello", transaction.message)
	}
	if transaction.recipientBlockchainAddress != "recipient" {
		t.Fatalf("recipient = %q, want recipient", transaction.recipientBlockchainAddress)
	}
	if transaction.senderBlockchainAddress != "sender" {
		t.Fatalf("sender = %q, want sender", transaction.senderBlockchainAddress)
	}
}

func TestTransactionMarshalJSONUsesCorrectAddressFields(t *testing.T) {
	transaction := NewTransaction("hello", "recipient", "sender", nil, nil, 2)

	got, err := json.Marshal(transaction)
	if err != nil {
		t.Fatalf("Marshal transaction: %v", err)
	}

	const want = `{"message":"hello","recipientBlockchainAddress":"recipient","senderBlockchainAddress":"sender","value":2}`
	if string(got) != want {
		t.Fatalf("transaction JSON = %s, want %s", got, want)
	}
}
