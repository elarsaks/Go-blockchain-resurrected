# Backend API

The backend has two HTTP API layers:

- wallet server: dashboard-facing gateway on port `5000`
- blockchain server: miner-node API on ports `5001`, `5002`, and `5003`

The React dashboard should use the wallet server. The wallet server forwards selected requests to a miner.

## Wallet Server API

Base URL in Docker/local dashboard mode:

```text
http://localhost:5000
```

### `GET /`

Returns a small JSON map of available routes.

### `POST /user/wallet`

Creates a new user wallet, registers that wallet address on the selected miner, and returns wallet details.

Response:

```json
{
  "privateKey": "...",
  "publicKey": "...",
  "blockchainAddress": "..."
}
```

Notes:

- The wallet server creates the wallet with `pkg/wallet.NewWallet`.
- It registers the wallet by posting the address to the selected miner's `/wallet/register`.
- The miner records a zero-value registration transaction and starts mining.

### `POST /miner/wallet?miner_id=1`

Selects a miner gateway and returns that miner's wallet details.

Query parameters:

| Name | Required | Description |
| --- | --- | --- |
| `miner_id` | no | `"1"`, `"2"`, or `"3"`. Defaults to miner 1 when not recognized. |

Response:

```json
{
  "privateKey": "...",
  "publicKey": "...",
  "blockchainAddress": "..."
}
```

Side effect:

- Updates the wallet server's active gateway.

### `GET /wallet/balance?blockchainAddress=...`

Fetches the selected miner's balance calculation for a blockchain address.

Response:

```json
{
  "balance": 1,
  "error": ""
}
```

If the miner cannot find the address in the chain, `error` contains a message.

### `POST /transaction`

Creates and signs a transaction, then forwards it to the selected miner.

Request:

```json
{
  "message": "USER TRANSACTION",
  "recipientBlockchainAddress": "...",
  "senderBlockchainAddress": "...",
  "senderPrivateKey": "...",
  "senderPublicKey": "...",
  "value": "1"
}
```

Response:

```json
{
  "message": "success"
}
```

Notes:

- The dashboard sends value as a string.
- The wallet server parses the value into `float32`.
- The wallet server signs the transaction with the sender private key.
- The miner verifies the signature before accepting the transaction.

### `GET /miner/blocks?amount=10`

Fetches recent blocks from the selected miner.

Query parameters:

| Name | Required | Description |
| --- | --- | --- |
| `amount` | yes | Positive integer. The current miner handler always returns the latest 10 blocks, but the wallet server still validates this query parameter. |

Response:

```json
[
  {
    "timestamp": 0,
    "nonce": 0,
    "previousHash": "...",
    "transactions": []
  }
]
```

## Blockchain Server API

Each miner exposes the same API. In Docker/local mode:

```text
http://localhost:5001
http://localhost:5002
http://localhost:5003
```

### `GET /chain`

Returns the miner's full blockchain.

Response:

```json
{
  "chain": [
    {
      "timestamp": 0,
      "nonce": 0,
      "previousHash": "...",
      "transactions": []
    }
  ]
}
```

### `GET /miner/blocks`

Returns the latest blocks from this miner. The current handler returns `bc.GetBlocks(10)`.

### `POST /miner/wallet`

Returns this miner's wallet.

### `GET /balance?blockchainAddress=...`

Calculates the total balance for an address by scanning all transactions in the chain.

### `POST /transactions`

Accepts a signed transaction from the wallet server.

Request:

```json
{
  "message": "...",
  "recipientBlockchainAddress": "...",
  "senderBlockchainAddress": "...",
  "senderPublicKey": "...",
  "signature": "...",
  "value": 1
}
```

Behavior:

- validates required fields
- parses sender public key and signature
- verifies the ECDSA signature against the transaction JSON
- checks sender balance
- adds the transaction to the miner's transaction pool
- broadcasts the accepted transaction to peers with `PUT /transactions`

### `PUT /transactions`

Accepts a transaction broadcast from another miner and adds it to this miner's transaction pool.

### `GET /transactions`

Returns the current transaction pool and its length.

### `DELETE /transactions`

Clears the transaction pool.

This is used after a miner creates a block and asks peer miners to clear already-mined transactions.

### `GET /mine`

Mines one block if there are pending transactions.

Response:

```json
{
  "message": "success"
}
```

If there are no pending transactions, the handler returns `400` and `{"message":"fail"}`.

### `GET /mine/start`

Starts mining using the miner's scheduled mining loop.

### `PUT /consensus`

Runs conflict resolution by fetching peer chains and replacing the local chain if a longer valid chain is found.

### `POST /wallet/register`

Registers a new wallet address by adding a zero-value blockchain transaction and starting mining.

Request:

```json
{
  "blockchainAddress": "..."
}
```
