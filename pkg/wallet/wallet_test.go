package wallet

import (
	"crypto/ecdsa"
	"crypto/sha256"
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

func TestNewWalletFromPrivateKeyHexRestoresWalletIdentity(t *testing.T) {
	original, err := NewWalletWithError()
	if err != nil {
		t.Fatalf("NewWalletWithError returned error: %v", err)
	}

	restored, err := NewWalletFromPrivateKeyHex(original.PrivateKeyStr())
	if err != nil {
		t.Fatalf("NewWalletFromPrivateKeyHex returned error: %v", err)
	}

	if restored.PrivateKeyStr() != original.PrivateKeyStr() {
		t.Fatalf("private key = %q, want %q", restored.PrivateKeyStr(), original.PrivateKeyStr())
	}
	if restored.PublicKeyStr() != original.PublicKeyStr() {
		t.Fatalf("public key = %q, want %q", restored.PublicKeyStr(), original.PublicKeyStr())
	}
	if restored.BlockchainAddress() != original.BlockchainAddress() {
		t.Fatalf("blockchain address = %q, want %q", restored.BlockchainAddress(), original.BlockchainAddress())
	}
}

func TestNewWalletFromPrivateKeyHexRejectsInvalidInput(t *testing.T) {
	if restored, err := NewWalletFromPrivateKeyHex("not-hex"); err == nil {
		t.Fatalf("NewWalletFromPrivateKeyHex returned wallet %#v, want error", restored)
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

func TestGenerateSignatureWithErrorSignsTransactionJSON(t *testing.T) {
	wallet, err := NewWalletWithError()
	if err != nil {
		t.Fatalf("NewWalletWithError returned error: %v", err)
	}
	transaction := NewTransaction(
		"hello",
		"recipient",
		wallet.BlockchainAddress(),
		wallet.PrivateKey(),
		wallet.PublicKey(),
		2,
	)

	signature, err := transaction.GenerateSignatureWithError()
	if err != nil {
		t.Fatalf("GenerateSignatureWithError returned error: %v", err)
	}
	if signature == nil {
		t.Fatal("signature is nil")
	}

	message, err := json.Marshal(transaction)
	if err != nil {
		t.Fatalf("Marshal transaction: %v", err)
	}
	digest := sha256.Sum256(message)

	if !ecdsa.Verify(wallet.PublicKey(), digest[:], signature.R, signature.S) {
		t.Fatal("signature did not verify against transaction JSON")
	}
}

func TestGenerateSignatureUsesErrorReturningPath(t *testing.T) {
	wallet, err := NewWalletWithError()
	if err != nil {
		t.Fatalf("NewWalletWithError returned error: %v", err)
	}
	transaction := NewTransaction(
		"hello",
		"recipient",
		wallet.BlockchainAddress(),
		wallet.PrivateKey(),
		wallet.PublicKey(),
		2,
	)

	if signature := transaction.GenerateSignature(); signature == nil {
		t.Fatal("GenerateSignature returned nil")
	}
}
