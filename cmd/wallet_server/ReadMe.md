# WALLET - Wallet Server

## About
The Wallet Server is a web server that handles the connection between clients and the blockchain network.

## Dependenies
- Golang

## Installation

If you haven't already installed the project from the parent folder, follow these steps to set up the Wallet Server:

1. Navigate to the parent folder of the project in your terminal.

2. Run the following command to download the necessary dependencies:

```bash
go mod tidy
```

## Running
To run it directly via Golang, execute the following command from the repository root:
```bash
PORT=5000 MINER_HOST=127.0.0.1 go run ./cmd/wallet_server
```


