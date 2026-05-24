package main

import "testing"

func TestMinerGatewayUsesDockerMinerServiceNames(t *testing.T) {
	t.Setenv("MINER_HOST", "miner-1")

	server := &WalletServer{}
	tests := map[string]string{
		"1": "http://miner-1:5001",
		"2": "http://miner-2:5002",
		"3": "http://miner-3:5003",
	}

	for minerID, want := range tests {
		got, err := server.MinerGateway(minerID)
		if err != nil {
			t.Fatalf("MinerGateway(%q) returned error: %v", minerID, err)
		}
		if got != want {
			t.Fatalf("MinerGateway(%q) = %q, want %q", minerID, got, want)
		}
	}
}

func TestMinerGatewayUsesSingleHostForLocalPorts(t *testing.T) {
	t.Setenv("MINER_HOST", "127.0.0.1")

	server := &WalletServer{}
	tests := map[string]string{
		"1": "http://127.0.0.1:5001",
		"2": "http://127.0.0.1:5002",
		"3": "http://127.0.0.1:5003",
	}

	for minerID, want := range tests {
		got, err := server.MinerGateway(minerID)
		if err != nil {
			t.Fatalf("MinerGateway(%q) returned error: %v", minerID, err)
		}
		if got != want {
			t.Fatalf("MinerGateway(%q) = %q, want %q", minerID, got, want)
		}
	}
}
