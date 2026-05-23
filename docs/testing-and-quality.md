# Testing and Quality

## Current Test Layers

| Layer | Location | Purpose |
| --- | --- | --- |
| Go package tests | `pkg/**` | Validate blockchain, transaction, JSON, mining, and wallet behavior. |
| React unit/component tests | `cmd/react_dashboard/src/tests` | Validate dashboard API helpers, reducers, utilities, and selected components. |
| Docker e2e test | `scripts/e2e-transfer.sh` | Validates the complete wallet/miner/dashboard-adjacent backend flow through Docker Compose. |

## Go Checks

Run all Go quality checks from the repository root:

```bash
gofmt -w $(find . -path './cmd/react_dashboard/node_modules' -prune -o -name '*.go' -print)
go list ./... | grep -v '/cmd/react_dashboard/node_modules/' | xargs go vet
go list ./... | grep -v '/cmd/react_dashboard/node_modules/' | xargs go test
```

The `node_modules` exclusion exists because one JavaScript dependency contains Go files that should not be treated as part of this module.

## React Checks

Run from `cmd/react_dashboard`:

```bash
npm run format:check
npm run lint
npx tsc --noEmit
npm run test:coverage
npm run build
```

## E2E Transfer Check

Run from the repository root:

```bash
./scripts/e2e-transfer.sh
```

This test is useful because package tests alone do not prove that the wallet server, miner server, HTTP handlers, Docker networking, mining timer, and balance polling work together.

## What Tests Should Protect

High-value backend behavior:

- transaction field mapping between wallet and block packages
- transaction JSON stability
- block JSON encode/decode compatibility
- signature generation and verification
- balance calculation
- proof-of-work difficulty validation
- chain validation
- miner wallet registration
- block endpoint decode behavior through the wallet server

High-value frontend behavior:

- API error parsing
- wallet validation
- wallet reducer state transitions
- miner selection behavior
- recent block rendering
- notification behavior

## CI Expectations

Pull requests should report at least:

- `go-quality`
- `e2e-transfer`

Dashboard-only changes should also report:

- `react-dashboard`

The branch protection rule should require `go-quality` and `e2e-transfer` so package tests and full-stack behavior both matter before merge.
