package block

import (
	"encoding/json"
	"testing"
)

func TestNewTransactionMapsFields(t *testing.T) {
	transaction := NewTransaction("sender", "recipient", "hello", 1.5)

	if transaction.senderBlockchainAddress != "sender" {
		t.Fatalf("sender = %q, want sender", transaction.senderBlockchainAddress)
	}
	if transaction.recipientBlockchainAddress != "recipient" {
		t.Fatalf("recipient = %q, want recipient", transaction.recipientBlockchainAddress)
	}
	if transaction.message != "hello" {
		t.Fatalf("message = %q, want hello", transaction.message)
	}
	if transaction.value != 1.5 {
		t.Fatalf("value = %v, want 1.5", transaction.value)
	}
}

func TestAddTransactionMapsFields(t *testing.T) {
	bc := NewBlockchain("miner", 5001)

	ok, err := bc.AddTransaction(MINING_SENDER, "recipient", "reward", 1, nil, nil)
	if err != nil {
		t.Fatalf("AddTransaction returned error: %v", err)
	}
	if !ok {
		t.Fatal("AddTransaction returned false")
	}

	pool := bc.TransactionPool()
	if len(pool) != 1 {
		t.Fatalf("transaction pool length = %d, want 1", len(pool))
	}

	gotJSON, err := json.Marshal(pool[0])
	if err != nil {
		t.Fatalf("Marshal transaction: %v", err)
	}

	const want = `{"message":"reward","recipientBlockchainAddress":"recipient","senderBlockchainAddress":"THE BLOCKCHAIN","value":1}`
	if string(gotJSON) != want {
		t.Fatalf("transaction JSON = %s, want %s", gotJSON, want)
	}
}
