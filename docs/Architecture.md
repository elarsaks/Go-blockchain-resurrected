# Project Architecture

This project is currently maintained as a local development blockchain playground. It is not intended to describe or support a production deployment.

## Runtime Components

```mermaid
flowchart LR
    user[Browser user]
    dashboard[React Dashboard\ncmd/react_dashboard\nlocalhost:3000]
    wallet[Wallet Server\ncmd/wallet_server\nlocalhost:5000]
    miner1[Blockchain Miner 1\ncmd/blockchain_server\nlocalhost:5001]
    miner2[Blockchain Miner 2\ncmd/blockchain_server\nlocalhost:5002]
    miner3[Blockchain Miner 3\ncmd/blockchain_server\nlocalhost:5003]

    user --> dashboard
    dashboard -->|HTTP API\n/user/wallet\n/wallet/balance\n/transaction\n/miner/blocks\n/miner/wallet| wallet

    wallet -->|gateway HTTP API| miner1
    wallet -.->|selectable gateway| miner2
    wallet -.->|selectable gateway| miner3

    miner1 <-->|peer sync\n/transactions\n/consensus| miner2
    miner2 <-->|peer sync\n/transactions\n/consensus| miner3
    miner3 <-->|peer sync\n/transactions\n/consensus| miner1
```

## Code Organization

```mermaid
flowchart TB
    root[Repository root]
    compose[docker-compose.yml\nlocal multi-service runtime]
    backend[Go backend]
    frontend[React frontend]
    shared[Shared Go packages]

    root --> compose
    root --> backend
    root --> frontend
    root --> shared

    backend --> walletServer[cmd/wallet_server\nAPI gateway and wallet operations]
    backend --> blockchainServer[cmd/blockchain_server\nminer node HTTP API]

    frontend --> reactDashboard[cmd/react_dashboard\nCreate React App dashboard]

    shared --> blockPkg[pkg/block\nblocks, transactions, mining, neighbours]
    shared --> walletPkg[pkg/wallet\nwallet and key generation]
    shared --> utilsPkg[pkg/utils\nJSON, CORS, ECDSA, neighbour helpers]

    walletServer --> shared
    blockchainServer --> shared
    reactDashboard --> walletServer
```

## Local Ports

| Component | Default URL |
| --- | --- |
| React dashboard | http://localhost:3000 |
| Wallet server | http://localhost:5000 |
| Miner 1 | http://localhost:5001 |
| Miner 2 | http://localhost:5002 |
| Miner 3 | http://localhost:5003 |

## Request Flow

```mermaid
sequenceDiagram
    participant U as Browser
    participant R as React Dashboard
    participant W as Wallet Server
    participant M as Selected Miner
    participant P as Peer Miners

    U->>R: Use dashboard
    R->>W: Create wallet or transaction
    W->>M: Forward blockchain request
    M->>M: Validate and store pending transaction
    M->>M: Mine block
    M->>P: Share transactions and consensus state
    P-->>M: Return peer state
    M-->>W: Return blockchain response
    W-->>R: Return API response
    R-->>U: Render updated state
```
