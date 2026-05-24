package block

import "testing"

func TestCalculateTotalBalanceTracksReceivedAndSpentTransactions(t *testing.T) {
	bc := NewBlockchain("miner", 5001)

	bc.transactionPool = []*Transaction{
		NewTransaction(MINING_SENDER, "alice", "initial balance", 10),
	}
	bc.CreateBlock(1, bc.LastBlock().Hash())

	bc.transactionPool = []*Transaction{
		NewTransaction("alice", "bob", "payment", 3),
	}
	bc.CreateBlock(2, bc.LastBlock().Hash())

	aliceBalance, err := bc.CalculateTotalBalance("alice")
	if err != nil {
		t.Fatalf("CalculateTotalBalance(alice) returned error: %v", err)
	}
	if aliceBalance != 7 {
		t.Fatalf("alice balance = %v, want 7", aliceBalance)
	}

	bobBalance, err := bc.CalculateTotalBalance("bob")
	if err != nil {
		t.Fatalf("CalculateTotalBalance(bob) returned error: %v", err)
	}
	if bobBalance != 3 {
		t.Fatalf("bob balance = %v, want 3", bobBalance)
	}
}

func TestNewBlockchainRegistersMinerAddressWithZeroBalance(t *testing.T) {
	bc := NewBlockchain("miner", 5001)

	balance, err := bc.CalculateTotalBalance("miner")
	if err != nil {
		t.Fatalf("CalculateTotalBalance(miner) returned error: %v", err)
	}
	if balance != 0 {
		t.Fatalf("miner balance = %v, want 0", balance)
	}
}

func TestCalculateTotalBalanceTreatsLocalMinerAddressAsKnownAfterChainReplacement(t *testing.T) {
	bc := NewBlockchain("miner", 5001)
	bc.chain = []*Block{NewBlock(0, [32]byte{}, nil)}

	balance, err := bc.CalculateTotalBalance("miner")
	if err != nil {
		t.Fatalf("CalculateTotalBalance(miner) returned error: %v", err)
	}
	if balance != 0 {
		t.Fatalf("miner balance = %v, want 0", balance)
	}
}

func TestCalculateTotalBalanceRejectsUnknownAddress(t *testing.T) {
	bc := NewBlockchain("miner", 5001)

	if balance, err := bc.CalculateTotalBalance("missing"); err == nil {
		t.Fatalf("CalculateTotalBalance returned balance %v, want error", balance)
	}
}
