package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/btcsuite/btcutil/base58"
	"github.com/elarsaks/Go-blockchain/pkg/utils"
	"golang.org/x/crypto/ripemd160"
)

type Wallet struct {
	privateKey        *ecdsa.PrivateKey
	publicKey         *ecdsa.PublicKey
	blockchainAddress string
}

type Transaction struct {
	message                    string
	recipientBlockchainAddress string
	senderBlockchainAddress    string
	senderPrivateKey           *ecdsa.PrivateKey
	senderPublicKey            *ecdsa.PublicKey
	value                      float32
}

type TransactionRequest struct {
	Message                    *string `json:"message"`
	RecipientBlockchainAddress *string `json:"recipientBlockchainAddress"`
	SenderBlockchainAddress    *string `json:"senderBlockchainAddress"`
	SenderPrivateKey           *string `json:"senderPrivateKey"`
	SenderPublicKey            *string `json:"senderPublicKey"`
	Value                      *string `json:"value"`
}

func NewWallet() *Wallet {
	w, err := NewWalletWithError()
	if err != nil {
		log.Printf("ERROR: create wallet: %v", err)
		return nil
	}
	return w
}

func NewWalletWithError() (*Wallet, error) {
	// 1. Creating ECDSA private key (32 bytes) public key (64 bytes)
	w := new(Wallet)
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("generate private key: %w", err)
	}
	w.privateKey = privateKey
	w.publicKey = &w.privateKey.PublicKey
	w.blockchainAddress = blockchainAddressFromPublicKey(w.publicKey)
	return w, nil
}

func NewWalletFromPrivateKeyHex(privateKeyHex string) (*Wallet, error) {
	privateKeyHex = strings.TrimSpace(privateKeyHex)
	if privateKeyHex == "" {
		return nil, fmt.Errorf("private key is required")
	}
	if len(privateKeyHex)%2 != 0 {
		privateKeyHex = "0" + privateKeyHex
	}

	privateKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("decode private key: %w", err)
	}

	curve := elliptic.P256()
	d := new(big.Int).SetBytes(privateKeyBytes)
	if d.Sign() <= 0 || d.Cmp(curve.Params().N) >= 0 {
		return nil, fmt.Errorf("private key is outside P-256 range")
	}

	x, y := curve.ScalarBaseMult(d.Bytes())
	if x == nil || y == nil {
		return nil, fmt.Errorf("derive public key from private key")
	}

	privateKey := &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{Curve: curve, X: x, Y: y},
		D:         d,
	}

	w := &Wallet{
		privateKey: privateKey,
		publicKey:  &privateKey.PublicKey,
	}
	w.blockchainAddress = blockchainAddressFromPublicKey(w.publicKey)
	return w, nil
}

func blockchainAddressFromPublicKey(publicKey *ecdsa.PublicKey) string {
	// 2. Perform SHA-256 hashing on the public key (32 bytes).
	h2 := sha256.New()
	h2.Write(publicKey.X.Bytes())
	h2.Write(publicKey.Y.Bytes())
	digest2 := h2.Sum(nil)
	// 3. Perform RIPEMD-160 hashing on the result of SHA-256 (20 bytes).
	h3 := ripemd160.New()
	h3.Write(digest2)
	digest3 := h3.Sum(nil)
	// 4. Add version byte in front of RIPEMD-160 hash (0x00 for Main Network).
	vd4 := make([]byte, 21)
	vd4[0] = 0x00
	copy(vd4[1:], digest3[:])
	// 5. Perform SHA-256 hash on the extended RIPEMD-160 result.
	h5 := sha256.New()
	h5.Write(vd4)
	digest5 := h5.Sum(nil)
	// 6. Perform SHA-256 hash on the result of the previous SHA-256 hash.
	h6 := sha256.New()
	h6.Write(digest5)
	digest6 := h6.Sum(nil)
	// 7. Take the first 4 bytes of the second SHA-256 hash for checksum.
	chsum := digest6[:4]
	// 8. Add the 4 checksum bytes from 7 at the end of extended RIPEMD-160 hash from 4 (25 bytes).
	dc8 := make([]byte, 25)
	copy(dc8[:21], vd4[:])
	copy(dc8[21:], chsum[:])
	// 9. Convert the result from a byte string into base58.
	return base58.Encode(dc8)
}

// PrivateKey returns the ECDSA private key of the wallet.
func (w *Wallet) PrivateKey() *ecdsa.PrivateKey {
	return w.privateKey
}

// PrivateKeyStr returns the hexadecimal representation of the private key.
func (w *Wallet) PrivateKeyStr() string {
	return fmt.Sprintf("%x", w.privateKey.D.Bytes())
}

// PublicKey returns the ECDSA public key of the wallet.
func (w *Wallet) PublicKey() *ecdsa.PublicKey {
	return w.publicKey
}

// PublicKeyStr returns the hexadecimal representation of the public key.
func (w *Wallet) PublicKeyStr() string {
	return fmt.Sprintf("%064x%064x", w.publicKey.X.Bytes(), w.publicKey.Y.Bytes())
}

// BlockchainAddress returns the blockchain address associated with the wallet.
func (w *Wallet) BlockchainAddress() string {
	return w.blockchainAddress
}

// MarshalJSON returns the JSON representation of the wallet.
func (w *Wallet) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		PrivateKey        string `json:"privateKey"`
		PublicKey         string `json:"publicKey"`
		BlockchainAddress string `json:"blockchainAddress"`
	}{
		PrivateKey:        w.PrivateKeyStr(),
		PublicKey:         w.PublicKeyStr(),
		BlockchainAddress: w.BlockchainAddress(),
	})
}

// NewTransaction creates a new transaction with the given details.
func NewTransaction(
	message string,
	recipient string,
	sender string,
	privateKey *ecdsa.PrivateKey,
	publicKey *ecdsa.PublicKey,
	value float32) *Transaction {
	return &Transaction{message, recipient, sender, privateKey, publicKey, value}
}

// GenerateSignature generates the signature for the transaction.
func (t *Transaction) GenerateSignature() *utils.Signature {
	signature, err := t.GenerateSignatureWithError()
	if err != nil {
		log.Printf("ERROR: generate signature: %v", err)
		return nil
	}
	return signature
}

func (t *Transaction) GenerateSignatureWithError() (*utils.Signature, error) {
	if t.senderPrivateKey == nil {
		return nil, fmt.Errorf("sender private key is required")
	}

	m, err := json.Marshal(t)
	if err != nil {
		return nil, fmt.Errorf("marshal transaction: %w", err)
	}

	log.Println("Generate signature", string(m))

	h := sha256.Sum256([]byte(m))
	r, s, err := ecdsa.Sign(rand.Reader, t.senderPrivateKey, h[:])
	if err != nil {
		return nil, fmt.Errorf("sign transaction: %w", err)
	}
	return &utils.Signature{R: r, S: s}, nil
}

// MarshalJSON returns the JSON representation of the transaction.
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

// Validate checks if all the required fields in the transaction request are present.
func (tr *TransactionRequest) Validate() bool {
	if tr.SenderPrivateKey == nil ||
		tr.SenderBlockchainAddress == nil ||
		tr.RecipientBlockchainAddress == nil ||
		tr.SenderPublicKey == nil ||
		tr.Message == nil ||
		tr.Value == nil {
		return false
	}
	return true
}
