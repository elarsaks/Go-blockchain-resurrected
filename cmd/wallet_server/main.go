package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/elarsaks/Go-blockchain/cmd/wallet_server/handlers"
	"github.com/elarsaks/Go-blockchain/pkg/utils"
	"github.com/gorilla/mux"
)

type WalletServer struct {
	port    uint16
	gateway string
}

// Make sure WalletServer implements handlers.WalletServer
var _ handlers.WalletServer = &WalletServer{}

func (ws *WalletServer) Port() uint16 {
	return ws.port
}

func GetHost() string {
	if host := os.Getenv("MINER_HOST"); host != "" {
		return host
	}
	return "127.0.0.1"
}

func (ws *WalletServer) Gateway() string {
	return ws.gateway
}

func (ws *WalletServer) MinerGateway(minerID string) (string, error) {
	host := GetHost()

	switch minerID {
	case "1":
		return fmt.Sprintf("http://%s:5001", host), nil
	case "2":
		return fmt.Sprintf("http://%s:5002", host), nil
	case "3":
		return fmt.Sprintf("http://%s:5003", host), nil
	default:
		return "", fmt.Errorf("invalid miner_id %q", minerID)
	}
}

func NewWalletServer(port uint16, gateway string) *WalletServer {
	return &WalletServer{port, gateway}
}

// Run the WalletServer
func (ws *WalletServer) Run() {
	// Create router
	router := mux.NewRouter()
	router.Use(utils.CorsMiddleware())

	// Create an instance of WalletServerHandler
	handler := handlers.NewWalletServerHandler(ws)

	// Define routes
	router.HandleFunc("/", handler.GetApiDescription).Methods(http.MethodGet)
	router.HandleFunc("/user/wallet", handler.GetUserWallet).Methods(http.MethodPost)
	router.HandleFunc("/wallet/balance", handler.GetWalletBalance).Methods(http.MethodGet)
	router.HandleFunc("/transaction", handler.CreateTransaction).Methods(http.MethodPost)
	router.HandleFunc("/miner/blocks", handler.GetBlocks).Methods(http.MethodGet)
	router.HandleFunc("/miner/wallet", handler.GetMinerWallet).Methods(http.MethodPost)

	// Start server
	log.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(int(ws.Port())), router))
}

func init() {
	log.SetPrefix("Wallet Server: ")
}

func main() {
	// Retrieve port from environment variable
	portStr := os.Getenv("PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil || port <= 0 {
		port = 5000 // Default value
	}

	// Print gateway and port
	log.Printf("Port: %d\n", port)

	defaultGateway, err := (&WalletServer{}).MinerGateway("1")
	if err != nil {
		log.Fatal(err)
	}

	app := NewWalletServer(uint16(port), defaultGateway)
	app.Run()
}
