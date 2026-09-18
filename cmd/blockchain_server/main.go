package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/elarsaks/Go-blockchain/cmd/blockchain_server/handlers"
	"github.com/elarsaks/Go-blockchain/pkg/block"
	"github.com/elarsaks/Go-blockchain/pkg/utils"
	"github.com/elarsaks/Go-blockchain/pkg/wallet"
	"github.com/gorilla/mux"
)

type BlockchainServer struct {
	port       uint16
	blockchain *block.Blockchain
	Wallet     *wallet.Wallet
	//* NOTE: In real world app we would not attach the wallet to the server
	// But for the sake of simplicity we will do it here,
	// because we dont store miners credentials in a database.
}

// Get the port of the BlockchainServer
func (bcs *BlockchainServer) Port() uint16 {
	return bcs.port
}

// GetWallet method for BlockchainServer
func (bcs *BlockchainServer) GetWallet() *wallet.Wallet {
	return bcs.Wallet
}

// Create a new instance of BlockchainServer
func NewBlockchainServer(port uint16) *BlockchainServer {
	minersWallet, err := loadMinerWallet()
	if err != nil {
		log.Printf("ERROR: load miner wallet: %v", err)
		minersWallet = wallet.NewWallet()
	}

	return &BlockchainServer{
		port:       port,
		Wallet:     minersWallet,
		blockchain: block.NewBlockchain(minersWallet.BlockchainAddress(), port),
	}
}

// Get the blockchain of the BlockchainServer
func (bcs *BlockchainServer) GetBlockchain() *block.Blockchain {
	return bcs.blockchain
}

// Run the BlockchainServer
func (bcs *BlockchainServer) Run(ctx context.Context) error {
	bcs.GetBlockchain().Run(ctx)

	// Register the miner's wallet
	// bcs.RegisterMinersWallet()

	router := mux.NewRouter()
	router.Use(utils.CorsMiddleware())
	router.PathPrefix("/").Methods(http.MethodOptions).HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	handler := handlers.NewBlockchainServerHandler(bcs)

	// Define routes
	router.HandleFunc("/chain", handler.GetChain).Methods(http.MethodGet)
	router.HandleFunc("/balance", handler.Balance).Methods(http.MethodGet)
	router.HandleFunc("/consensus", handler.Consensus).Methods(http.MethodPut)
	router.HandleFunc("/mine", handler.Mine).Methods(http.MethodGet)
	router.HandleFunc("/mine/start", handler.StartMine).Methods(http.MethodGet)
	router.HandleFunc("/miner/blocks", handler.GetBlocks).Methods(http.MethodGet)
	router.HandleFunc("/miner/wallet", handler.MinerWallet).Methods(http.MethodPost)
	router.HandleFunc("/transactions", handler.HandleGetTransaction).Methods(http.MethodGet)
	router.HandleFunc("/transactions", handler.HandlePostTransaction).Methods(http.MethodPost)
	router.HandleFunc("/transactions", handler.HandlePutTransaction).Methods(http.MethodPut)
	router.HandleFunc("/transactions", handler.HandleDeleteTransaction).Methods(http.MethodDelete)
	router.HandleFunc("/wallet/register", handler.RegisterWallet).Methods(http.MethodPost)

	// Start the server
	server := &http.Server{
		Addr:              "0.0.0.0:" + strconv.Itoa(int(bcs.Port())),
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- server.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			return err
		}
		return nil
	case err := <-errCh:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
}

// Function to initialize the logger
func init() {
	log.SetPrefix("Blockchain: ")
}

// Main function
func main() {
	// Retrieve port from environment variable
	portStr := os.Getenv("PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil || port <= 0 {
		port = 5001 // Default value
	}

	// Print port
	log.Printf("Port: %d\n", port)

	app := NewBlockchainServer(uint16(port))
	log.Printf("blockchainAddress %v", app.GetWallet().BlockchainAddress())

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := app.Run(ctx); err != nil {
		log.Fatal(err)
	}
}

func loadMinerWallet() (*wallet.Wallet, error) {
	if privateKey := os.Getenv("MINER_PRIVATE_KEY"); strings.TrimSpace(privateKey) != "" {
		return wallet.NewWalletFromPrivateKeyHex(privateKey)
	}

	keyFile := os.Getenv("MINER_PRIVATE_KEY_FILE")
	if keyFile == "" {
		var err error
		keyFile, err = defaultMinerKeyFile()
		if err != nil {
			return nil, err
		}
	}

	privateKey, err := os.ReadFile(filepath.Clean(keyFile))
	if err == nil {
		return wallet.NewWalletFromPrivateKeyHex(string(privateKey))
	}
	if !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	minerWallet, err := wallet.NewWalletWithError()
	if err != nil {
		return nil, err
	}
	if err := os.WriteFile(filepath.Clean(keyFile), []byte(minerWallet.PrivateKeyStr()), 0600); err != nil {
		return nil, err
	}
	return minerWallet, nil
}

func defaultMinerKeyFile() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	minerDir := filepath.Join(configDir, "go-blockchain")
	if err := os.MkdirAll(minerDir, 0700); err != nil {
		return "", err
	}
	return filepath.Join(minerDir, "miner_private_key"), nil
}
