package block

import "testing"

func TestChainReturnsSnapshot(t *testing.T) {
	bc := NewBlockchain("miner", 5001)

	chain := bc.Chain()
	chain[0] = nil

	if bc.LastBlock() == nil {
		t.Fatal("mutating returned chain slice changed blockchain state")
	}
}

func TestTransactionPoolReturnsSnapshot(t *testing.T) {
	bc := NewBlockchain("miner", 5001)
	_, err := bc.AddTransaction(MINING_SENDER, "recipient", "reward", 1, nil, nil)
	if err != nil {
		t.Fatalf("AddTransaction returned error: %v", err)
	}

	pool := bc.TransactionPool()
	pool[0] = nil

	if bc.TransactionPool()[0] == nil {
		t.Fatal("mutating returned transaction pool slice changed blockchain state")
	}
}

func TestBlockTransactionsReturnsSnapshot(t *testing.T) {
	block := NewBlock(1, [32]byte{}, []*Transaction{
		NewTransaction("sender", "recipient", "message", 1),
	})

	transactions := block.Transactions()
	transactions[0] = nil

	if block.Transactions()[0] == nil {
		t.Fatal("mutating returned transactions slice changed block state")
	}
}
