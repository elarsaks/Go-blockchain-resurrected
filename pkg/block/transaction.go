package block

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/elarsaks/Go-blockchain/pkg/utils"
)

// --- Types ---
// Transaction represents a single transaction in the blockchain.
type Transaction struct {
	message                    string
	recipientBlockchainAddress string
	senderBlockchainAddress    string
	value                      float32
}

// TransactionRequest represents a request to create a new transaction.
type TransactionRequest struct {
	Message                    *string  `json:"message"`
	RecipientBlockchainAddress *string  `json:"recipientBlockchainAddress"`
	SenderBlockchainAddress    *string  `json:"senderBlockchainAddress"`
	SenderPublicKey            *string  `json:"senderPublicKey"`
	Signature                  *string  `json:"signature"`
	Value                      *float32 `json:"value"`
}

// AmountResponse represents the response with the amount in a transaction.
type BalanceResponse struct {
	Balance float32 `json:"balance"`
	Error   string  `json:"error"`
}

func (bc *Blockchain) AddTransaction(sender string,
	recipient string,
	message string,
	value float32,
	senderPublicKey *ecdsa.PublicKey,
	s *utils.Signature) (bool, error) {
	bc.mux.Lock()
	defer bc.mux.Unlock()

	return bc.addTransactionLocked(sender, recipient, message, value, senderPublicKey, s)
}

func (bc *Blockchain) addTransactionLocked(sender string,
	recipient string,
	message string,
	value float32,
	senderPublicKey *ecdsa.PublicKey,
	s *utils.Signature) (bool, error) {
	if strings.TrimSpace(sender) == "" {
		return false, fmt.Errorf("ERROR: sender blockchain address is required")
	}
	if strings.TrimSpace(recipient) == "" {
		return false, fmt.Errorf("ERROR: recipient blockchain address is required")
	}
	if strings.TrimSpace(message) == "" {
		return false, fmt.Errorf("ERROR: transaction message is required")
	}
	if sender != MINING_SENDER && value <= 0 {
		return false, fmt.Errorf("ERROR: transaction value must be positive")
	}
	if sender != MINING_SENDER && sender == recipient {
		return false, fmt.Errorf("ERROR: sender and recipient must be different")
	}

	// Create a new transaction
	t := NewTransaction(sender, recipient, message, value)

	// If the sender is the mining address, add the transaction to the pool and return true
	if sender == MINING_SENDER {
		if value < 0 {
			return false, fmt.Errorf("ERROR: system transaction value must not be negative")
		}
		bc.transactionPool = append(bc.transactionPool, t)
		return true, nil
	}

	if senderPublicKey == nil || s == nil {
		return false, fmt.Errorf("ERROR: sender public key and signature are required")
	}

	// If the transaction signature is not verified, return false and an error
	if !bc.VerifyTransactionSignature(senderPublicKey, s, t) {
		return false, fmt.Errorf("ERROR: Verify Transaction")
	}

	// Calculate the total balance of the sender
	balance, err := bc.calculateTotalBalanceLocked(sender)
	if err != nil {
		// If there is an error calculating the balance, return false and the error
		return false, fmt.Errorf("ERROR: CalculateTotalAmount: %v", err)
	}

	// If the sender's balance is less than the value of the transaction, return false and an error
	if balance < value {
		return false, fmt.Errorf("ERROR: Not enough balance in a wallet")
	}

	// Add the transaction to the transaction pool
	bc.transactionPool = append(bc.transactionPool, t)

	// Return true and no error
	return true, nil
}

// Empty the transaction pool the Blockchain
func (bc *Blockchain) ClearTransactionPool() {
	bc.mux.Lock()
	defer bc.mux.Unlock()

	bc.transactionPool = bc.transactionPool[:0]
}

// Create a new transaction
func (bc *Blockchain) CreateTransaction(sender string, recipient string, message string, value float32,
	senderPublicKey *ecdsa.PublicKey, s *utils.Signature) (bool, error) {

	isTransacted, err := bc.AddTransaction(sender, recipient, message, value, senderPublicKey, s)

	// If there was an error while adding the transaction, log the error and return it
	if err != nil {

		log.Printf("ERROR: %v", err)
		return false, err
	}

	// If the transaction was added successfully, broadcast it to the network
	if isTransacted {
		// Reverse engineer this part of the code
		for _, n := range bc.Neighbors() {
			publicKeyStr := fmt.Sprintf("%064x%064x", senderPublicKey.X.Bytes(),
				senderPublicKey.Y.Bytes())
			signatureStr := s.String()
			bt := &TransactionRequest{
				Message:                    &message,
				RecipientBlockchainAddress: &recipient,
				SenderBlockchainAddress:    &sender,
				SenderPublicKey:            &publicKeyStr,
				Signature:                  &signatureStr,
				Value:                      &value,
			}
			m, err := json.Marshal(bt)
			if err != nil {
				log.Printf("ERROR: %v", err)
				return false, err
			}
			endpoint := peerEndpoint(n, "/transactions")
			if err := doPeerRequest("PUT", endpoint, bytes.NewReader(m)); err != nil {
				log.Printf("ERROR: %v", err)
				continue
			}
		}
	}

	return isTransacted, nil
}

// Copy the transaction pool
func (bc *Blockchain) CopyTransactionPool() []*Transaction {
	bc.mux.RLock()
	defer bc.mux.RUnlock()

	return bc.copyTransactionPoolLocked()
}

func (bc *Blockchain) copyTransactionPoolLocked() []*Transaction {
	transactions := make([]*Transaction, 0, len(bc.transactionPool))
	for _, t := range bc.transactionPool {
		transactions = append(transactions,
			NewTransaction(
				t.senderBlockchainAddress,
				t.recipientBlockchainAddress,
				t.message,
				t.value))
	}
	return transactions
}

// NewTransaction creates a new transaction.
func NewTransaction(sender string, recipient string, message string, value float32) *Transaction {
	return &Transaction{message, recipient, sender, value}
}

// Get the transaction pool the Blockchain
func (bc *Blockchain) TransactionPool() []*Transaction {
	bc.mux.RLock()
	defer bc.mux.RUnlock()

	return append([]*Transaction(nil), bc.transactionPool...)
}

// Validate checks if the transaction request is valid.
func (tr *TransactionRequest) Validate() bool {
	if tr.SenderBlockchainAddress == nil ||
		tr.RecipientBlockchainAddress == nil ||
		tr.SenderPublicKey == nil ||
		tr.Message == nil ||
		tr.Value == nil ||
		tr.Signature == nil {
		return false
	}
	return strings.TrimSpace(*tr.SenderBlockchainAddress) != "" &&
		strings.TrimSpace(*tr.RecipientBlockchainAddress) != "" &&
		strings.TrimSpace(*tr.SenderPublicKey) != "" &&
		strings.TrimSpace(*tr.Message) != "" &&
		strings.TrimSpace(*tr.Signature) != "" &&
		*tr.Value > 0 &&
		*tr.SenderBlockchainAddress != *tr.RecipientBlockchainAddress
}

// Verify the signature of the transaction
func (bc *Blockchain) VerifyTransactionSignature(
	senderPublicKey *ecdsa.PublicKey,
	s *utils.Signature,
	t *Transaction) bool {
	if senderPublicKey == nil || senderPublicKey.X == nil || senderPublicKey.Y == nil ||
		s == nil || s.R == nil || s.S == nil {
		return false
	}

	m, _ := json.Marshal(t)

	log.Println("Validate signature", string(m))

	h := sha256.Sum256([]byte(m))
	return ecdsa.Verify(senderPublicKey, h[:], s.R, s.S)
}

// Print outputs the details of the transaction.
func (t *Transaction) Print() {
	fmt.Printf("%s\n", strings.Repeat("-", 40))
	fmt.Printf(" senderBlockchainAddress      %s\n", t.senderBlockchainAddress)
	fmt.Printf(" recipientBlockchainAddress   %s\n", t.recipientBlockchainAddress)
	fmt.Printf(" message                      %s\n", t.message)
	fmt.Printf(" value                          %.1f\n", t.value)
}

// MarshalJSON implements the Marshaler interface for the AmountResponse type.
func (br *BalanceResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Balance float32 `json:"balance"`
		Error   string  `json:"error"`
	}{
		Balance: br.Balance,
		Error:   br.Error,
	})
}

// MarshalJSON implements the Marshaler interface for the Transaction type.
func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Message   string  `json:"message"`
		Recipient string  `json:"recipientBlockchainAddress"`
		Sender    string  `json:"senderBlockchainAddress"`
		Value     float32 `json:"value"`
	}{
		Message:   t.message,
		Recipient: t.recipientBlockchainAddress,
		Sender:    t.senderBlockchainAddress,
		Value:     t.value,
	})
}

// UnmarshalJSON implements the Unmarshaler interface for the Transaction type.
func (t *Transaction) UnmarshalJSON(data []byte) error {
	var v struct {
		Message   *string  `json:"message"`
		Recipient *string  `json:"recipientBlockchainAddress"`
		Sender    *string  `json:"senderBlockchainAddress"`
		Value     *float32 `json:"value"`
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	if v.Message == nil {
		return fmt.Errorf("transaction message is required")
	}
	if v.Recipient == nil {
		return fmt.Errorf("transaction recipientBlockchainAddress is required")
	}
	if v.Sender == nil {
		return fmt.Errorf("transaction senderBlockchainAddress is required")
	}
	if v.Value == nil {
		return fmt.Errorf("transaction value is required")
	}
	if *v.Value < 0 {
		return fmt.Errorf("transaction value must not be negative")
	}
	t.message = *v.Message
	t.recipientBlockchainAddress = *v.Recipient
	t.senderBlockchainAddress = *v.Sender
	t.value = *v.Value
	return nil
}
