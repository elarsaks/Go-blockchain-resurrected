package block

import "testing"

func TestFilterOutSelfPort(t *testing.T) {
	neighbors := []string{
		"http://miner-1:5001",
		"http://miner-2:5002",
		"127.0.0.1:5001",
	}

	filtered := filterOutSelfPort(neighbors, "5001")

	if len(filtered) != 1 {
		t.Fatalf("filtered length = %d, want 1", len(filtered))
	}
	if filtered[0] != "http://miner-2:5002" {
		t.Fatalf("filtered[0] = %q, want http://miner-2:5002", filtered[0])
	}
}
