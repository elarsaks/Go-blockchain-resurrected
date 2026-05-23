package utils

import "testing"

func TestFindNeighborsRejectsNonIPv4Host(t *testing.T) {
	if neighbors := FindNeighbors("localhost", 5001, 0, 1, 5001, 5002); neighbors != nil {
		t.Fatalf("FindNeighbors returned %v, want nil", neighbors)
	}
}
