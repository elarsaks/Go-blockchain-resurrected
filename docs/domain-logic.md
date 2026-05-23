# Domain Logic

This document explains how the blockchain, wallet, transaction, mining, and consensus logic work in the current codebase.

## Blocks

Defined in `pkg/block/block.go`.

A `Block` contains:

- `timestamp`
- `nonce`
- `previousHash`
- `transactions`

The fields are unexported, so JSON is handled through custom marshal/unmarshal methods.

The block hash is:

```text
sha256(json.Marshal(block))
```

This means the JSON representation is part of the hash contract. Changes to block JSON shape can affect proof-of-work validation.

## Blockchain

Defined in `pkg/block/blockchain.go`.

A `Blockchain` owns:

- `transactionPool`
- `chain`
- `blockchainAddress`
- `port`
- mining and neighbor mutexes
- peer neighbor URLs

`NewBlockchain` creates a genesis block by hashing an empty `Block` and then creating a block with nonce `0`.

`CreateBlock`:

1. Creates a new block from the current transaction pool.
2. Appends it to the chain.
3. Clears the local transaction pool.
4. Sends `DELETE /transactions` to neighbors to clear their pools.

`GetBlocks(amount)` returns the newest blocks in reverse order, so the newest block appears first.

## Transactions

Defined in `pkg/block/transaction.go`.

A block transaction contains:

- message
- recipient blockchain address
- sender blockchain address
- value

`AddTransaction` has two paths:

1. Mining sender path:
   - if sender is `MINING_SENDER`, the transaction is accepted without signature or balance checks
2. User transaction path:
   - verifies the transaction signature
   - calculates sender balance
   - rejects if balance is too low
   - appends the transaction to the transaction pool

`CreateTransaction` calls `AddTransaction` locally and then broadcasts the transaction to peer miners with `PUT /transactions`.

## Wallets

Defined in `pkg/wallet/wallet.go`.

`NewWallet`:

1. Generates an ECDSA P-256 private key.
2. Gets the public key from the private key.
3. Hashes the public key with SHA-256.
4. Hashes the result with RIPEMD-160.
5. Adds a version byte.
6. Computes a double-SHA-256 checksum.
7. Encodes the result with Base58.

The wallet JSON response includes:

- private key
- public key
- blockchain address

For this project, returning private keys is intentional for learning and dashboard-driven signing. It is not safe for a real wallet system.

## Signatures

Wallet-side transactions are signed in `pkg/wallet`.

The wallet transaction JSON is hashed with SHA-256 and signed with ECDSA:

```text
signature = ECDSA_sign(sha256(json.Marshal(transaction)))
```

Miner-side verification reconstructs the transaction JSON in `pkg/block`, hashes it the same way, and verifies the ECDSA signature using the sender public key.

This makes field order and JSON shape important. Sender, recipient, message, and value must map the same way on both sides.

## Mining

Defined in `pkg/block/mining.go`.

Mining is scheduled with:

```text
time.AfterFunc(time.Second * MINING_TIMER_SEC, bc.StartMining)
```

Default mining constants:

| Constant | Value |
| --- | ---: |
| `MINING_DIFFICULTY` | `3` |
| `MINING_REWARD` | `1.0` |
| `MINING_TIMER_SEC` | `20` |

`Mining`:

1. Locks the blockchain mining mutex.
2. Returns `false` if the transaction pool is empty.
3. Adds a mining reward transaction to the miner address.
4. Finds a nonce with `ProofOfWork`.
5. Creates a block.
6. Sends `PUT /consensus` to neighbors.

`ProofOfWork`:

1. Copies the current transaction pool.
2. Gets the previous block hash.
3. Starts nonce at `0`.
4. Increments nonce until `ValidProof` returns true.

`ValidProof` creates a candidate block using timestamp `0`, the candidate nonce, previous hash, and transactions. The candidate hash must start with `MINING_DIFFICULTY` zero characters.

## Balance Calculation

`CalculateTotalBalance(address)` scans every transaction in every block.

- If the address is the recipient, value is added.
- If the address is the sender, value is subtracted.
- If the address never appears, an error is returned.

Balances are currently `float32`. That is acceptable for a toy project, but real money-like values should use integer smallest units or a decimal type.

## Consensus

`ResolveConflicts` implements a simple longest-valid-chain strategy:

1. Fetch `/chain` from each neighbor.
2. Decode the peer blockchain.
3. If a peer chain is longer and `ValidChain` returns true, remember it.
4. Replace the local chain with the longest valid peer chain found.

`ValidChain` checks:

- each block's `previousHash` matches the previous block's hash
- each block satisfies proof-of-work validation

## Neighbor Discovery

Defined in `pkg/block/neighbour.go` and `pkg/utils/neighbor.go`.

When `MINER_HOST` is set, miners use fixed Docker names:

```text
http://miner-1:5001
http://miner-2:5002
http://miner-3:5003
```

When `MINER_HOST` is not set, the code tries local IPv4 neighbor discovery across a configured IP and port range.

## Known Limitations

- Chain and wallet state are memory-only.
- No durable transaction ledger exists outside miner memory.
- `float32` is used for balances.
- Some HTTP responses and errors are handled loosely.
- Some response bodies are not closed in older network paths.
- Mining, peer networking, consensus, and core domain logic are mixed in `pkg/block`.
- No real fork-choice/security model exists beyond longest valid chain.
- This is a learning project, not production blockchain infrastructure.
