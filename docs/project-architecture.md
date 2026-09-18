# Architecture

## Purpose

This repository is a local blockchain playground. It demonstrates the moving pieces of a blockchain-style system:

- wallet/key generation
- signed transactions
- miner nodes
- proof-of-work mining
- peer transaction broadcast
- longest-valid-chain conflict resolution
- a browser dashboard for interaction

The architecture is intentionally small and local. It is optimized for learning and experimentation rather than production safety.

## Runtime Components

| Component | Path | Default port | Responsibility |
| --- | --- | ---: | --- |
| React dashboard | `cmd/react_dashboard` | `3000` | Browser UI for creating wallets, selecting miners, sending transactions, viewing balances, and showing recent blocks. |
| Wallet server | `cmd/wallet_server` | `5000` | API gateway used by the dashboard. Generates user wallets, signs outgoing transactions, selects a miner gateway, and proxies reads/writes to miners. |
| Miner node 1 | `cmd/blockchain_server` | `5001` | Owns a blockchain instance, miner wallet, transaction pool, mining loop, and peer calls. |
| Miner node 2 | `cmd/blockchain_server` | `5002` | Same binary/configuration as miner 1 with a different port. |
| Miner node 3 | `cmd/blockchain_server` | `5003` | Same binary/configuration as miner 1 with a different port. |

## Docker Compose Topology

`docker-compose.yml` starts five services:

- `react-dashboard`
- `wallet-server`
- `miner-1`
- `miner-2`
- `miner-3`

The dashboard receives `VITE_GATEWAY_API_URL=http://localhost:5000`, so browser calls go to the wallet server on the host.

The wallet server receives `MINER_HOST=miner-1`. Its gateway selector builds miner URLs from that host and ports `5001`, `5002`, and `5003`.

The miner containers receive:

- `PORT=5001`, `PORT=5002`, or `PORT=5003`
- `MINER_HOST=miner`

When `MINER_HOST` is set on a miner, the miner creates static peers:

- `http://miner-1:5001`
- `http://miner-2:5002`
- `http://miner-3:5003`

It then filters out the peer with its own port.

## Code Organization

```text
cmd/
  blockchain_server/     Miner node HTTP server
  wallet_server/         Wallet API gateway HTTP server
  react_dashboard/       React/Vite dashboard

pkg/
  block/                 Blockchain domain model and mining logic
  wallet/                ECDSA wallet, address, and transaction signing logic
  utils/                 JSON status helpers, CORS, key/signature parsing, neighbor discovery

scripts/
  e2e-transfer.sh        Docker-based full-stack transfer test

docs/
  *.md                   Architecture and operation documentation
```

## Package Boundaries

### `pkg/block`

Owns the blockchain model:

- `Block`
- `Blockchain`
- `Transaction`
- transaction pool
- proof of work
- mining reward
- chain validation
- conflict resolution
- miner neighbor synchronization
- balance calculation

This package currently also performs HTTP calls to peer miners. That means the domain logic and network transport are coupled.

### `pkg/wallet`

Owns wallet and sender-side transaction behavior:

- creates an ECDSA private/public key pair
- derives a Base58 blockchain address
- serializes wallet data for API responses
- creates a transaction value object
- signs transaction JSON with ECDSA

### `pkg/utils`

Contains small shared helpers:

- `JsonStatus`
- CORS middleware
- hex string to ECDSA key/signature conversion
- local IPv4 neighbor discovery

## Server Responsibilities

### Wallet Server

The wallet server is the dashboard-facing API. It does not maintain blockchain state itself. Its most important responsibilities are:

- generate user wallets
- register user wallets on a miner
- fetch miner wallets
- choose the active miner gateway
- sign user-submitted transactions
- forward transactions to miners
- proxy balance and block reads

### Blockchain Server

Each blockchain server process represents one miner node. It:

- creates a miner wallet lazily on first blockchain access
- creates an in-memory blockchain
- starts neighbor synchronization
- starts a periodic mining loop
- exposes blockchain and transaction HTTP endpoints
- mines pending transactions into blocks
- adds a miner reward transaction before mining
- asks peer miners to resolve conflicts after mining

## State Model

State is in memory only. Restarting containers resets wallets, transaction pools, chains, and miner state.

| State | Owner | Persistence |
| --- | --- | --- |
| User wallet returned to dashboard | Browser/runtime response | Not persisted by backend |
| Miner wallet | Miner server process | In memory |
| Blockchain chain | Miner server process | In memory |
| Transaction pool | Miner server process | In memory |
| Selected wallet gateway | Wallet server process | In memory |

## Current Design Caveats

- Chain data is not persisted to disk.
- Wallet private keys are returned to the dashboard and are not stored securely.
- Balances use `float32`, which is not appropriate for real money-like values.
- Mining and HTTP peer calls live in the same package.
- Some package methods log or print directly rather than returning structured errors.
- Peer synchronization is simple and assumes local Docker service names.
- The blockchain is educational and should not be treated as production-grade cryptographic infrastructure.

## Distributed Architecture Target

The planned distributed version keeps the current learning scope, but moves runtime ownership into Kubernetes and adds an event/observability layer.

Planned runtime ownership:

| Component | Kubernetes shape | Notes |
| --- | --- | --- |
| React dashboard | Deployment + Service | Browser entrypoint, optionally exposed through Ingress. |
| Wallet server | Deployment + Service | API gateway for dashboard commands and reads. |
| Miner nodes | StatefulSet + headless Service | Stable miner identities for peer discovery and consensus demos. |
| Kafka | Helm chart or operator-managed cluster | Event stream for transaction, block, and consensus activity. |
| Prometheus | Helm chart or operator-managed deployment | Scrapes service and Kafka metrics. |
| Grafana | Deployment or chart-managed service | Dashboards for chain height, mining, transaction pool, latency, and errors. |

The first Kubernetes version should preserve existing HTTP behavior. Kafka should start as an event stream for auditability and monitoring instead of becoming the source of truth immediately. This keeps the migration small and makes each phase easier to test.

See [Distributed Systems Plan](distributed-systems-plan.md) for rollout phases and open decisions.
