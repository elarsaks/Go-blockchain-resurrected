package block

import (
	"testing"

	"github.com/elarsaks/Go-blockchain/pkg/utils"
	walletpkg "github.com/elarsaks/Go-blockchain/pkg/wallet"
)

func TestValidProofChecksDifficultyPrefix(t *testing.T) {
	bc := NewBlockchain("miner-address", 5001)
	previousHash := bc.LastBlock().Hash()
	transactions := []*Transaction{
		NewTransaction("sender", "recipient", "test message", 1),
	}
	difficulty := 2

	validNonce := 0
	for !bc.ValidProof(validNonce, previousHash, transactions, difficulty) {
		validNonce++
	}

	if !bc.ValidProof(validNonce, previousHash, transactions, difficulty) {
		t.Fatalf("expected nonce %d to satisfy proof difficulty", validNonce)
	}

	invalidNonce := 0
	for bc.ValidProof(invalidNonce, previousHash, transactions, difficulty) {
		invalidNonce++
	}

	if bc.ValidProof(invalidNonce, previousHash, transactions, difficulty) {
		t.Fatalf("expected nonce %d to fail proof difficulty", invalidNonce)
	}
}

func TestValidProofRejectsInvalidDifficulty(t *testing.T) {
	bc := NewBlockchain("miner-address", 5001)

	if bc.ValidProof(0, bc.LastBlock().Hash(), nil, -1) {
		t.Fatal("expected negative difficulty to be invalid")
	}

	if bc.ValidProof(0, bc.LastBlock().Hash(), nil, 65) {
		t.Fatal("expected difficulty larger than hash length to be invalid")
	}
}

func TestMiningCreatesBlockAndClearsTransactionPoolWithoutNeighbors(t *testing.T) {
	bc := NewBlockchain("miner-address", 5001)
	bc.neighbors = nil

	_, err := bc.AddTransaction(MINING_SENDER, "alice", "initial balance", 5, nil, nil)
	if err != nil {
		t.Fatalf("AddTransaction returned error: %v", err)
	}

	if !bc.Mining() {
		t.Fatal("Mining returned false")
	}

	if got := len(bc.Chain()); got != 2 {
		t.Fatalf("chain length = %d, want 2", got)
	}
	if got := len(bc.TransactionPool()); got != 0 {
		t.Fatalf("transaction pool length = %d, want 0", got)
	}

	transactions := bc.LastBlock().Transactions()
	if got := len(transactions); got != 2 {
		t.Fatalf("mined block transaction count = %d, want 2", got)
	}
	if transactions[1].recipientBlockchainAddress != "miner-address" ||
		transactions[1].value != MINING_REWARD {
		t.Fatalf("mining reward transaction = %#v, want reward for miner-address", transactions[1])
	}
}

func TestMiningReturnsFalseWithEmptyTransactionPool(t *testing.T) {
	bc := NewBlockchain("miner-address", 5001)
	bc.neighbors = nil

	if bc.Mining() {
		t.Fatal("Mining returned true for an empty transaction pool")
	}
	if got := len(bc.Chain()); got != 1 {
		t.Fatalf("chain length = %d, want 1", got)
	}
}

func TestCreateTransactionAddsSignedTransactionWithoutNeighbors(t *testing.T) {
	bc := NewBlockchain("miner-address", 5001)
	bc.neighbors = nil

	senderWallet := walletWithBalance(t, bc, 10)
	signature := signedWalletTransaction(t, senderWallet, "recipient", "payment", 3)

	ok, err := bc.CreateTransaction(
		senderWallet.BlockchainAddress(),
		"recipient",
		"payment",
		3,
		senderWallet.PublicKey(),
		signature,
	)
	if err != nil {
		t.Fatalf("CreateTransaction returned error: %v", err)
	}
	if !ok {
		t.Fatal("CreateTransaction returned false")
	}

	pool := bc.TransactionPool()
	if got := len(pool); got != 1 {
		t.Fatalf("transaction pool length = %d, want 1", got)
	}
	transaction := pool[0]
	if transaction.senderBlockchainAddress != senderWallet.BlockchainAddress() ||
		transaction.recipientBlockchainAddress != "recipient" ||
		transaction.message != "payment" ||
		transaction.value != 3 {
		t.Fatalf("transaction = %#v, want signed payment in pool", transaction)
	}
}

func TestCreateTransactionRejectsInsufficientBalanceWithoutNeighbors(t *testing.T) {
	bc := NewBlockchain("miner-address", 5001)
	bc.neighbors = nil

	senderWallet := walletWithBalance(t, bc, 1)
	signature := signedWalletTransaction(t, senderWallet, "recipient", "payment", 3)

	ok, err := bc.CreateTransaction(
		senderWallet.BlockchainAddress(),
		"recipient",
		"payment",
		3,
		senderWallet.PublicKey(),
		signature,
	)
	if err == nil {
		t.Fatal("expected CreateTransaction to return an insufficient balance error")
	}
	if ok {
		t.Fatal("CreateTransaction returned true for insufficient balance")
	}
	if got := len(bc.TransactionPool()); got != 0 {
		t.Fatalf("transaction pool length = %d, want 0", got)
	}
}

func walletWithBalance(t *testing.T, bc *Blockchain, value float32) *walletpkg.Wallet {
	t.Helper()

	senderWallet, err := walletpkg.NewWalletWithError()
	if err != nil {
		t.Fatalf("NewWalletWithError returned error: %v", err)
	}
	_, err = bc.AddTransaction(
		MINING_SENDER,
		senderWallet.BlockchainAddress(),
		"seed balance",
		value,
		nil,
		nil,
	)
	if err != nil {
		t.Fatalf("AddTransaction returned error: %v", err)
	}
	bc.CreateBlock(1, bc.LastBlock().Hash())

	return senderWallet
}

func signedWalletTransaction(
	t *testing.T,
	senderWallet *walletpkg.Wallet,
	recipient string,
	message string,
	value float32,
) *utils.Signature {
	t.Helper()

	transaction := walletpkg.NewTransaction(
		message,
		recipient,
		senderWallet.BlockchainAddress(),
		senderWallet.PrivateKey(),
		senderWallet.PublicKey(),
		value,
	)
	signature, err := transaction.GenerateSignatureWithError()
	if err != nil {
		t.Fatalf("GenerateSignatureWithError returned error: %v", err)
	}
	return signature
}
