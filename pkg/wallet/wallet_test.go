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

func TestNewWalletWithErrorCreatesWallet(t *testing.T) {
	wallet, err := NewWalletWithError()
	if err != nil {
		t.Fatalf("NewWalletWithError returned error: %v", err)
	}
	if wallet.PrivateKey() == nil {
		t.Fatal("private key is nil")
	}
	if wallet.PublicKey() == nil {
		t.Fatal("public key is nil")
	}
	if wallet.BlockchainAddress() == "" {
		t.Fatal("blockchain address is empty")
	}
}

func TestGenerateSignatureWithErrorRejectsMissingPrivateKey(t *testing.T) {
	transaction := NewTransaction("hello", "recipient", "sender", nil, nil, 2)

	signature, err := transaction.GenerateSignatureWithError()
	if err == nil {
		t.Fatal("expected missing private key to return an error")
	}
	if signature != nil {
		t.Fatal("signature should be nil when signing fails")
	}
}
