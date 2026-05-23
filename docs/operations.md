# Operations

## Run the Full Stack

Requirements:

- Docker
- Docker Compose

Run:

```bash
docker compose up --build
```

Local URLs:

| Component | URL |
| --- | --- |
| React dashboard | http://localhost:3000 |
| Wallet server | http://localhost:5000 |
| Miner 1 | http://localhost:5001 |
| Miner 2 | http://localhost:5002 |
| Miner 3 | http://localhost:5003 |

Stop:

```bash
docker compose down --remove-orphans
```

## Run the Wallet Server Locally

```bash
PORT=5000 MINER_HOST=127.0.0.1 go run ./cmd/wallet_server
```

The wallet server defaults to miner 1.

## Run a Miner Locally

```bash
PORT=5001 go run ./cmd/blockchain_server
```

For multiple local miners, start separate processes with different ports:

```bash
PORT=5001 go run ./cmd/blockchain_server
PORT=5002 go run ./cmd/blockchain_server
PORT=5003 go run ./cmd/blockchain_server
```

## Run the Dashboard Locally

```bash
cd cmd/react_dashboard
npm install
npm start
```

The dashboard uses `VITE_GATEWAY_API_URL` when provided. Otherwise it defaults to:

```text
http://localhost:5000
```

## Environment Variables

| Variable | Used by | Purpose |
| --- | --- | --- |
| `PORT` | wallet server, blockchain server | Selects the HTTP listen port. |
| `MINER_HOST` | wallet server | Host prefix used to build miner gateway URLs. |
| `MINER_HOST` | blockchain server | Enables Docker-style static neighbor URLs. |
| `VITE_GATEWAY_API_URL` | React dashboard | Dashboard API base URL. |

## Quality Checks

Go formatting:

```bash
gofmt -w $(find . -path './cmd/react_dashboard/node_modules' -prune -o -name '*.go' -print)
```

Go vet:

```bash
go list ./... | grep -v '/cmd/react_dashboard/node_modules/' | xargs go vet
```

Go tests:

```bash
go list ./... | grep -v '/cmd/react_dashboard/node_modules/' | xargs go test
```

React checks:

```bash
cd cmd/react_dashboard
npm run format:check
npm run lint
npx tsc --noEmit
npm run test:coverage
npm run build
```

## End-to-End Transfer Test

The script `scripts/e2e-transfer.sh` verifies the full Docker path:

1. Rebuilds and starts the Docker Compose stack.
2. Waits for the wallet server and miner wallet.
3. Creates a miner wallet and user wallet.
4. Waits for the miner reward.
5. Sends one coin from miner to user.
6. Waits for the transfer to be mined.
7. Stops the stack.

Run:

```bash
./scripts/e2e-transfer.sh
```

Useful knobs:

```bash
TIMEOUT_SECONDS=120 ./scripts/e2e-transfer.sh
POLL_INTERVAL_SECONDS=1 ./scripts/e2e-transfer.sh
WALLET_URL=http://localhost:5000 ./scripts/e2e-transfer.sh
```

## Debugging

Show service logs:

```bash
docker compose logs --no-color --tail=200
```

Call the wallet server:

```bash
curl http://localhost:5000/
curl -X POST 'http://localhost:5000/miner/wallet?miner_id=1'
curl 'http://localhost:5000/miner/blocks?amount=10'
```

Call a miner directly:

```bash
curl http://localhost:5001/chain
curl 'http://localhost:5001/balance?blockchainAddress=ADDRESS'
curl http://localhost:5001/transactions
```

Common symptoms:

| Symptom | Likely cause |
| --- | --- |
| Balance is zero after sending a transaction | Transaction has not been mined yet. Wait for the mining timer. |
| Wallet server returns a miner error | The selected miner gateway may not be running or may not know the wallet address yet. |
| Blocks endpoint returns decode errors | Miner block JSON and wallet-server block decoder disagree. Check block JSON tests. |
| Docker e2e times out | Mining interval, container startup time, or peer consensus may have delayed the transfer. Check Docker logs. |
