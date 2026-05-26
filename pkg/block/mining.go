package block

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

// ==============================
// Blockchain Proof and Mining Methods
// ==============================

// Mining creates a new block and adds it to the blockchain.
func (bc *Blockchain) Mining() bool {
	// Lock the blockchain while mining
	bc.mux.Lock()

	// Log out blockchain
	// bc.Print() // TODO: Remove debug

	//* DEBUG #Consensus Wallet registration mining should be done some where else
	// Don't mine when there is no transaction and blockchain already has few blocks
	if len(bc.transactionPool) == 0 {
		bc.mux.Unlock()
		return false
	}

	// Add a mining reward transaction
	_, err := bc.addTransactionLocked(MINING_SENDER, bc.blockchainAddress, "MINING REWARD", MINING_REWARD, nil, nil)

	// If an error occurred adding the transaction, log the error and return false
	if err != nil {
		bc.mux.Unlock()
		log.Printf("ERROR: %v", err)
		return false
	}

	// Find a new proof of work and create a new block
	transactions := bc.copyTransactionPoolLocked()
	previousHash := bc.chain[len(bc.chain)-1].Hash()
	nonce := bc.proofOfWork(transactions, previousHash)
	bc.createBlockLocked(nonce, previousHash, transactions)
	bc.mux.Unlock()

	// Log a successful mining operation
	// #debug
	log.Println("action=mining, status=success")

	bc.clearNeighborTransactionPools()

	// Send a consensus request to each neighbor
	for _, n := range bc.Neighbors() {

		fmt.Println("Send consensus to neigbour ", n)

		endpoint := peerEndpoint(n, "/consensus")
		if err := doPeerRequest(http.MethodPut, endpoint, nil); err != nil {
			log.Printf("ERROR: %v", err)
		}
	}

	// Return true indicating the mining operation was successful
	return true
}

// StartMining initiates the mining process.
func (bc *Blockchain) StartMining(ctx context.Context) {
	bc.Mining()
	// Schedule the next mining operation to occur after MINING_TIMER_SEC seconds.
	go func() {
		ticker := time.NewTicker(time.Second * MINING_TIMER_SEC)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				bc.Mining()
			}
		}
	}()
}

// ValidProof validates the proof of work.
func (bc *Blockchain) ValidProof(nonce int, previousHash [32]byte, transactions []*Transaction, difficulty int) bool {
	if difficulty < 0 {
		return false
	}
	zeros := strings.Repeat("0", difficulty)
	guessBlock := Block{0, nonce, previousHash, transactions}
	guessHashStr := fmt.Sprintf("%x", guessBlock.Hash())
	if difficulty > len(guessHashStr) {
		return false
	}

	return guessHashStr[:difficulty] == zeros
}

// ProofOfWork finds the proof of work.
func (bc *Blockchain) ProofOfWork() int {
	bc.mux.RLock()
	defer bc.mux.RUnlock()

	transactions := bc.copyTransactionPoolLocked()
	previousHash := bc.chain[len(bc.chain)-1].Hash()

	return bc.proofOfWork(transactions, previousHash)
}

func (bc *Blockchain) proofOfWork(transactions []*Transaction, previousHash [32]byte) int {
	nonce := 0
	for !bc.ValidProof(nonce, previousHash, transactions, MINING_DIFFICULTY) {
		nonce += 1
	}
	return nonce
}

// ValidChain validates the chain.
func (bc *Blockchain) ValidChain(chain []*Block) bool {
	if len(chain) == 0 {
		return false
	}

	//* DEGUG #Consensus
	preBlock := chain[0]
	currentIndex := 1
	for currentIndex < len(chain) {
		b := chain[currentIndex]
		if b.previousHash != preBlock.Hash() {
			return false
		}

		if !bc.ValidProof(b.Nonce(), b.PreviousHash(), b.Transactions(), MINING_DIFFICULTY) {
			return false
		}

		preBlock = b
		currentIndex += 1
	}
	return true
}

// ResolveConflicts resolves conflicts in the blockchain by checking the chains of its neighbors.
func (bc *Blockchain) ResolveConflicts() bool {
	// Initialize variables to track the longest chain and its length
	var longestChain []*Block = nil
	bc.mux.RLock()
	maxLength := len(bc.chain)
	bc.mux.RUnlock()

	// Iterate over the neighbors to fetch their chains
	for _, n := range bc.Neighbors() {
		fmt.Println("Resolve conflict with ", n)

		// Construct the endpoint URL to fetch the chain from the neighbor
		endpoint := peerEndpoint(n, "/chain")

		// Send an HTTP GET request to the neighbor's endpoint to fetch their chain
		resp, err := peerHTTPClient.Get(endpoint)
		if err != nil {

			// Log any error that occurred while fetching the chain
			log.Printf("ERROR: Failed to fetch chain from neighbor %s: %v", n, err)
			continue // Skip to the next neighbor in case of error
		}

		// Check the response status code to see if the request was successful
		if resp.StatusCode == http.StatusOK {
			var bcResp Blockchain
			decoder := json.NewDecoder(resp.Body)

			// Decode the JSON response into a Blockchain object
			err := decoder.Decode(&bcResp)
			_ = resp.Body.Close()
			if err != nil {
				// Log any error that occurred during JSON decoding
				log.Printf("ERROR: Failed to decode JSON response from neighbor %s: %v", n, err)
				continue // Skip to the next neighbor in case of error
			}

			// Get the chain from the neighbor's Blockchain object
			chain := bcResp.Chain()

			// Check if the fetched chain is longer than the current longest chain
			// and if it is a valid chain using bc.ValidChain()
			if len(chain) > maxLength && bc.ValidChain(chain) {
				maxLength = len(chain)
				longestChain = chain
			}
		} else {
			_ = resp.Body.Close()
			// Log the status code if the request to the neighbor's endpoint was not successful
			log.Printf("WARNING: Failed to fetch chain from neighbor %s. Status code: %d", n, resp.StatusCode)
		}
	}

	// If a longer valid chain was found, replace the blockchain's chain with it
	if longestChain != nil {
		bc.mux.Lock()
		bc.chain = longestChain
		bc.mux.Unlock()
		log.Printf("INFO: Resolved conflicts. Replaced blockchain with the longest valid chain.")
		return true
	}

	// If no longer valid chain was found, log and return false
	log.Printf("INFO: No longer valid chain found among neighbors. No conflicts resolved.")
	return false
}

// ==============================
// Blockchain Wallet and Balance Methods
// ==============================

// RegisterNewWallet registers a new wallet on the blockchain.
func (bc *Blockchain) RegisterNewWallet(blockchainAddress string, message string) bool {

	// Add a transaction for the new wallet
	_, err := bc.AddTransaction(MINING_SENDER, blockchainAddress, message, 0, nil, nil)

	// If an error occurred adding the transaction, log the error and return false
	if err != nil {
		log.Printf("ERROR: %v", err)
		return false
	}

	// Mine a new block when the wallet is registered successfully
	bc.Mining()

	// Return true indicating the wallet was registered successfully
	return true
}

// CalculateTotalBalance calculates the total balance of crypto on the specific address in the Blockchain.
func (bc *Blockchain) CalculateTotalBalance(blockchainAddress string) (float32, error) {
	bc.mux.RLock()
	defer bc.mux.RUnlock()

	return bc.calculateTotalBalanceLocked(blockchainAddress)
}

func (bc *Blockchain) calculateTotalBalanceLocked(blockchainAddress string) (float32, error) {
	var totalBalance float32 = 0.0
	addressFound := false

	for _, b := range bc.chain {
		for _, t := range b.transactions {
			value := t.value

			if blockchainAddress == t.recipientBlockchainAddress {
				totalBalance += value
				addressFound = true
			}

			if blockchainAddress == t.senderBlockchainAddress {
				totalBalance -= value
				addressFound = true
			}
		}
	}

	if !addressFound {
		if blockchainAddress == bc.blockchainAddress {
			return 0.0, nil
		}
		return 0.0, fmt.Errorf("Address not found in the Blockchain")
	}

	return totalBalance, nil
}
