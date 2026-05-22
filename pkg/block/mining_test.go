package block

import "testing"

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
