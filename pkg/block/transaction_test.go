package block

import (
	"encoding/json"
	"testing"

	walletpkg "github.com/elarsaks/Go-blockchain/pkg/wallet"
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

func TestAddTransactionRejectsNegativeSystemValue(t *testing.T) {
	bc := NewBlockchain("miner", 5001)

	ok, err := bc.AddTransaction(MINING_SENDER, "recipient", "reward", -1, nil, nil)
	if err == nil {
		t.Fatal("expected AddTransaction to reject negative system value")
	}
	if ok {
		t.Fatal("AddTransaction returned true for negative system value")
	}
	if got := len(bc.TransactionPool()); got != 0 {
		t.Fatalf("transaction pool length = %d, want 0", got)
	}
}

func TestTransactionRequestValidateRejectsEmptyAndInvalidValues(t *testing.T) {
	message := ""
	recipient := "recipient"
	sender := "sender"
	publicKey := "public-key"
	signature := "signature"
	value := float32(1)

	request := &TransactionRequest{
		Message:                    &message,
		RecipientBlockchainAddress: &recipient,
		SenderBlockchainAddress:    &sender,
		SenderPublicKey:            &publicKey,
		Signature:                  &signature,
		Value:                      &value,
	}

	if request.Validate() {
		t.Fatal("expected empty message to be invalid")
	}

	message = "payment"
	value = 0
	if request.Validate() {
		t.Fatal("expected zero value to be invalid")
	}

	value = 1
	recipient = sender
	if request.Validate() {
		t.Fatal("expected self-transfer to be invalid")
	}
}

func TestVerifyTransactionSignatureAcceptsWalletSignature(t *testing.T) {
	senderWallet, err := walletpkg.NewWalletWithError()
	if err != nil {
		t.Fatalf("NewWalletWithError returned error: %v", err)
	}
	walletTransaction := walletpkg.NewTransaction(
		"payment",
		"recipient",
		senderWallet.BlockchainAddress(),
		senderWallet.PrivateKey(),
		senderWallet.PublicKey(),
		1,
	)
	signature, err := walletTransaction.GenerateSignatureWithError()
	if err != nil {
		t.Fatalf("GenerateSignatureWithError returned error: %v", err)
	}

	bc := NewBlockchain("miner", 5001)
	blockTransaction := NewTransaction(senderWallet.BlockchainAddress(), "recipient", "payment", 1)

	if !bc.VerifyTransactionSignature(senderWallet.PublicKey(), signature, blockTransaction) {
		t.Fatal("expected wallet signature to verify")
	}
}

func TestVerifyTransactionSignatureRejectsTamperedTransaction(t *testing.T) {
	senderWallet, err := walletpkg.NewWalletWithError()
	if err != nil {
		t.Fatalf("NewWalletWithError returned error: %v", err)
	}
	walletTransaction := walletpkg.NewTransaction(
		"payment",
		"recipient",
		senderWallet.BlockchainAddress(),
		senderWallet.PrivateKey(),
		senderWallet.PublicKey(),
		1,
	)
	signature, err := walletTransaction.GenerateSignatureWithError()
	if err != nil {
		t.Fatalf("GenerateSignatureWithError returned error: %v", err)
	}

	bc := NewBlockchain("miner", 5001)
	tamperedTransaction := NewTransaction(senderWallet.BlockchainAddress(), "attacker", "payment", 1)

	if bc.VerifyTransactionSignature(senderWallet.PublicKey(), signature, tamperedTransaction) {
		t.Fatal("expected tampered transaction to fail signature verification")
	}
}
