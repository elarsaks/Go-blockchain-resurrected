# Request Flows

## Startup

```mermaid
sequenceDiagram
    participant Docker as Docker Compose
    participant W as Wallet Server
    participant M1 as Miner 1
    participant M2 as Miner 2
    participant M3 as Miner 3

    Docker->>W: Start PORT=5000, MINER_HOST=miner-1
    Docker->>M1: Start PORT=5001, MINER_HOST=miner
    Docker->>M2: Start PORT=5002, MINER_HOST=miner
    Docker->>M3: Start PORT=5003, MINER_HOST=miner

    M1->>M1: Create miner wallet lazily
    M1->>M1: Create blockchain and genesis block
    M1->>M1: Set neighbors and start mining timer

    M2->>M2: Same initialization
    M3->>M3: Same initialization
```

Each miner uses the same binary. Port and Docker service name determine which node it behaves as.

## Miner Wallet Selection

```mermaid
sequenceDiagram
    participant UI as React Dashboard
    participant W as Wallet Server
    participant M as Selected Miner

    UI->>W: POST /miner/wallet?miner_id=1
    W->>W: Set gateway to miner 1
    W->>M: POST /miner/wallet
    M-->>W: Miner wallet JSON
    W-->>UI: Miner wallet JSON
```

The dashboard includes the selected `miner_id` on later block, balance, and transaction requests.

## User Wallet Creation

```mermaid
sequenceDiagram
    participant UI as React Dashboard
    participant W as Wallet Server
    participant M as Selected Miner

    UI->>W: POST /user/wallet
    W->>W: Generate ECDSA wallet
    W->>M: POST /wallet/register
    M->>M: Add zero-value registration transaction
    M->>M: Start mining
    M-->>W: Registration success
    W-->>UI: User wallet JSON
```

The returned user wallet includes private key, public key, and blockchain address.

## Transaction Submission

```mermaid
sequenceDiagram
    participant UI as React Dashboard
    participant W as Wallet Server
    participant M as Selected Miner
    participant P as Peer Miners

    UI->>W: POST /transaction
    W->>W: Parse value
    W->>W: Recreate sender keys
    W->>W: Sign wallet transaction
    W->>M: POST /transactions
    M->>M: Validate required fields
    M->>M: Verify signature
    M->>M: Check sender balance
    M->>M: Add to transaction pool
    M->>P: PUT /transactions
    P->>P: Add broadcast transaction to peer pools
    M-->>W: Created
    W-->>UI: Success/failure JSON
```

The transaction is not spendable by the recipient until a miner includes it in a block.

## Mining

```mermaid
sequenceDiagram
    participant M as Miner
    participant P as Peer Miners

    M->>M: Timer calls StartMining
    M->>M: Mining checks transaction pool
    M->>M: Add mining reward transaction
    M->>M: ProofOfWork finds nonce
    M->>M: CreateBlock
    M->>P: DELETE /transactions
    M->>P: PUT /consensus
    P->>P: ResolveConflicts
```

Mining returns early if the transaction pool is empty.

## Balance Read

```mermaid
sequenceDiagram
    participant UI as React Dashboard
    participant W as Wallet Server
    participant M as Selected Miner

    UI->>W: GET /wallet/balance?blockchainAddress=...
    W->>M: GET /balance?blockchainAddress=...
    M->>M: Scan all block transactions
    M-->>W: BalanceResponse
    W-->>UI: BalanceResponse
```

The balance is calculated from chain history only. Pending transaction-pool entries do not count.

## Recent Block Read

```mermaid
sequenceDiagram
    participant UI as React Dashboard
    participant W as Wallet Server
    participant M as Selected Miner

    UI->>W: GET /miner/blocks?amount=10
    W->>M: GET /miner/blocks?amount=10
    M->>M: Get latest blocks
    M-->>W: Block list JSON
    W-->>UI: Block list JSON
```

The dashboard uses this flow to render recent blockchain activity.

## Planned Kubernetes Startup

```mermaid
sequenceDiagram
    participant K8s as Kubernetes
    participant W as wallet-server Deployment
    participant M as miner StatefulSet
    participant K as Kafka
    participant P as Prometheus

    K8s->>K: Start Kafka broker(s)
    K8s->>M: Start miner-0, miner-1, miner-2
    K8s->>W: Start wallet-server
    M->>M: Create or restore blockchain state
    M->>M: Discover peer miner DNS names
    W->>K: Publish service startup event
    M->>K: Publish miner startup event
    P->>W: Scrape metrics endpoint
    P->>M: Scrape metrics endpoints
```

In the first distributed version, Kubernetes changes service discovery and runtime management, but the user-facing command path can stay HTTP-based.

## Planned Event Publishing

```mermaid
sequenceDiagram
    participant UI as React Dashboard
    participant W as Wallet Server
    participant M as Selected Miner
    participant K as Kafka
    participant G as Grafana / Consumers

    UI->>W: POST /transaction
    W->>M: POST /transactions
    M->>M: Validate and add to transaction pool
    M->>K: Publish transactions.accepted
    M-->>W: Created
    W->>K: Publish transactions.created
    W-->>UI: Success JSON
    G->>K: Consume event stream or read derived metrics
```

Kafka should initially describe successful state changes. It should not replace transaction validation, mining, or consensus until the current behavior is stable in Kubernetes.
