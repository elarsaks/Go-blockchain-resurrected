package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/elarsaks/Go-blockchain/pkg/wallet"
)

func TestLoadMinerWalletUsesPrivateKeyEnv(t *testing.T) {
	expected, err := wallet.NewWalletWithError()
	if err != nil {
		t.Fatalf("NewWalletWithError returned error: %v", err)
	}
	t.Setenv("MINER_PRIVATE_KEY", expected.PrivateKeyStr())

	got, err := loadMinerWallet()
	if err != nil {
		t.Fatalf("loadMinerWallet returned error: %v", err)
	}
	if got.BlockchainAddress() != expected.BlockchainAddress() {
		t.Fatalf("blockchain address = %q, want %q", got.BlockchainAddress(), expected.BlockchainAddress())
	}
}

func TestLoadMinerWalletPersistsGeneratedPrivateKey(t *testing.T) {
	t.Setenv("MINER_PRIVATE_KEY", "")
	keyFile := filepath.Join(t.TempDir(), "miner.key")
	t.Setenv("MINER_PRIVATE_KEY_FILE", keyFile)

	first, err := loadMinerWallet()
	if err != nil {
		t.Fatalf("first loadMinerWallet returned error: %v", err)
	}

	if _, err := os.Stat(keyFile); err != nil {
		t.Fatalf("expected key file to be written: %v", err)
	}

	second, err := loadMinerWallet()
	if err != nil {
		t.Fatalf("second loadMinerWallet returned error: %v", err)
	}

	if second.PrivateKeyStr() != first.PrivateKeyStr() {
		t.Fatalf("private key = %q, want %q", second.PrivateKeyStr(), first.PrivateKeyStr())
	}
}
