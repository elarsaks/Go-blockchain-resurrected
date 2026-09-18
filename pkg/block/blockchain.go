package block

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
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
	mux               sync.RWMutex
	neighbors         []string
	muxNeighbors      sync.Mutex
}

var peerHTTPClient = &http.Client{Timeout: 5 * time.Second}

// NewBlockchain creates a new instance of Blockchain.
func NewBlockchain(blockchainAddress string, port uint16) *Blockchain {
	b := &Block{}
	bc := new(Blockchain)
	bc.blockchainAddress = blockchainAddress
	bc.transactionPool = []*Transaction{
		NewTransaction(MINING_SENDER, blockchainAddress, "REGISTER MINER WALLET", 0),
	}
	bc.CreateBlock(0, b.Hash())
	bc.port = port
	return bc
}

// Chain returns the chain of the Blockchain.
func (bc *Blockchain) Chain() []*Block {
	bc.mux.RLock()
	defer bc.mux.RUnlock()

	return append([]*Block(nil), bc.chain...)
}

// Run initializes and runs the Blockchain.
func (bc *Blockchain) Run(ctx context.Context) {
	bc.StartSyncNeighbors(ctx)
	bc.StartMining(ctx) // Start mining automatically
}

// CreateBlock creates a new block and appends it to the blockchain.
func (bc *Blockchain) CreateBlock(nonce int, previousHash [32]byte) *Block {
	bc.mux.Lock()
	transactions := bc.copyTransactionPoolLocked()
	block := bc.createBlockLocked(nonce, previousHash, transactions)
	bc.mux.Unlock()

	bc.clearNeighborTransactionPools()
	return block
}

func (bc *Blockchain) createBlockLocked(nonce int, previousHash [32]byte, transactions []*Transaction) *Block {
	b := NewBlock(nonce, previousHash, transactions)
	bc.chain = append(bc.chain, b)
	bc.transactionPool = []*Transaction{}
	return b
}

// LastBlock returns the last block of the Blockchain.
func (bc *Blockchain) LastBlock() *Block {
	bc.mux.RLock()
	defer bc.mux.RUnlock()

	return bc.lastBlockLocked()
}

func (bc *Blockchain) lastBlockLocked() *Block {
	return bc.chain[len(bc.chain)-1]
}

// GetBlocks returns the latest blocks of the Blockchain.
func (bc *Blockchain) GetBlocks(amount int) []*Block {
	bc.mux.RLock()
	defer bc.mux.RUnlock()

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
	bc.mux.RLock()
	defer bc.mux.RUnlock()

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

func (bc *Blockchain) clearNeighborTransactionPools() {
	for _, n := range bc.Neighbors() {
		endpoint := peerEndpoint(n, "/transactions")
		if err := doPeerRequest(http.MethodDelete, endpoint, nil); err != nil {
			log.Printf("ERROR: delete neighbor transactions: %v", err)
		}
	}
}

func doPeerRequest(method, endpoint string, body io.Reader) error {
	req, err := http.NewRequest(method, endpoint, body)
	if err != nil {
		return fmt.Errorf("create %s request: %w", method, err)
	}

	resp, err := peerHTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("send %s request: %w", method, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("%s returned status %d", endpoint, resp.StatusCode)
	}

	return nil
}
