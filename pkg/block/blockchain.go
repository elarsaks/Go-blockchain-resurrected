package block

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

// ==============================
// Blockchain Struct and Methods
// ==============================

// Blockchain represents the entire blockchain structure.
type Blockchain struct {
	transactionPool   []*Transaction
	chain             []*Block
	blockchainAddress string
	port              uint16
	mux               sync.Mutex
	neighbors         []string
	muxNeighbors      sync.Mutex
}

var peerHTTPClient = &http.Client{Timeout: 5 * time.Second}

// NewBlockchain creates a new instance of Blockchain.
func NewBlockchain(blockchainAddress string, port uint16) *Blockchain {
	b := &Block{}
	bc := new(Blockchain)
	bc.blockchainAddress = blockchainAddress
	bc.CreateBlock(0, b.Hash())
	bc.port = port
	return bc
}

// Chain returns the chain of the Blockchain.
func (bc *Blockchain) Chain() []*Block {
	return bc.chain
}

// Run initializes and runs the Blockchain.
func (bc *Blockchain) Run() {
	bc.StartSyncNeighbors()
	bc.ResolveConflicts()
	bc.StartMining() // Start mining automatically
}

// CreateBlock creates a new block and appends it to the blockchain.
func (bc *Blockchain) CreateBlock(nonce int, previousHash [32]byte) *Block {
	b := NewBlock(nonce, previousHash, bc.transactionPool)
	bc.chain = append(bc.chain, b)
	bc.transactionPool = []*Transaction{}
	for _, n := range bc.neighbors {
		endpoint := fmt.Sprintf("http://%s/transactions", n)
		req, err := http.NewRequest("DELETE", endpoint, nil)
		if err != nil {
			log.Printf("ERROR: create delete transactions request: %v", err)
			continue
		}
		resp, err := peerHTTPClient.Do(req)
		if err != nil {
			log.Printf("ERROR: delete neighbor transactions: %v", err)
			continue
		}
		_ = resp.Body.Close()
		log.Printf("%v", resp)
	}
	return b
}

// LastBlock returns the last block of the Blockchain.
func (bc *Blockchain) LastBlock() *Block {
	return bc.chain[len(bc.chain)-1]
}

// GetBlocks returns the latest blocks of the Blockchain.
func (bc *Blockchain) GetBlocks(amount int) []*Block {
	n := len(bc.chain)
	var last10Blocks []*Block
	if n > amount {
		last10Blocks = append([]*Block(nil), bc.chain[n-amount:n]...)
	} else {
		last10Blocks = append([]*Block(nil), bc.chain...)
	}

	// Reverse the slice
	for i := len(last10Blocks)/2 - 1; i >= 0; i-- {
		opp := len(last10Blocks) - 1 - i
		last10Blocks[i], last10Blocks[opp] = last10Blocks[opp], last10Blocks[i]
	}

	return last10Blocks
}

// Print displays the entire blockchain.
func (bc *Blockchain) Print() {
	for i, block := range bc.chain {
		fmt.Printf("%s Chain %d %s\n", strings.Repeat("=", 25), i,
			strings.Repeat("=", 25))
		block.Print()
	}
	fmt.Printf("%s\n", strings.Repeat("*", 25))
}

// JSON Handling for Blockchain

// MarshalJSON customizes the JSON encoding of the blockchain.
func (bc *Blockchain) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Blocks []*Block `json:"chain"`
	}{
		Blocks: bc.chain,
	})
}

// UnmarshalJSON customizes the JSON decoding of the blockchain.
func (bc *Blockchain) UnmarshalJSON(data []byte) error {
	v := &struct {
		Blocks *[]*Block `json:"chain"`
	}{
		Blocks: &bc.chain,
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	return nil
}
