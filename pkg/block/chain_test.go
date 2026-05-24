package block

import "testing"

func TestValidChainAcceptsValidMinedChain(t *testing.T) {
	bc := blockchainWithMinedTransaction(t)

	if !bc.ValidChain(bc.Chain()) {
		t.Fatal("expected mined chain to be valid")
	}
}

func TestValidChainRejectsTamperedPreviousHash(t *testing.T) {
	bc := blockchainWithMinedTransaction(t)
	chain := bc.Chain()
	chain[1].previousHash = [32]byte{1}

	if bc.ValidChain(chain) {
		t.Fatal("expected chain with tampered previous hash to be invalid")
	}
}

func TestValidChainRejectsTamperedTransaction(t *testing.T) {
	bc := blockchainWithMinedTransaction(t)
	chain := bc.Chain()
	chain[1].transactions[0].value = 99

	if bc.ValidChain(chain) {
		t.Fatal("expected chain with tampered transaction to be invalid")
	}
}

func blockchainWithMinedTransaction(t *testing.T) *Blockchain {
	t.Helper()

	bc := NewBlockchain("miner", 5001)
	_, err := bc.AddTransaction(MINING_SENDER, "alice", "fund alice", 1, nil, nil)
	if err != nil {
		t.Fatalf("AddTransaction returned error: %v", err)
	}

	nonce := bc.ProofOfWork()
	bc.CreateBlock(nonce, bc.LastBlock().Hash())

	if len(bc.Chain()) != 2 {
		t.Fatalf("chain length = %d, want 2", len(bc.Chain()))
	}

	return bc
}
