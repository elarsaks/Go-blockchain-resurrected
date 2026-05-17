#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
WALLET_URL="${WALLET_URL:-http://localhost:5000}"
COMPOSE="${COMPOSE:-docker compose}"
POLL_INTERVAL_SECONDS="${POLL_INTERVAL_SECONDS:-2}"
TIMEOUT_SECONDS="${TIMEOUT_SECONDS:-90}"

cd "$ROOT_DIR"

cleanup() {
  local status=$?

  if [[ "$status" -ne 0 ]]; then
    echo "E2E transfer failed. Recent Docker Compose logs:" >&2
    $COMPOSE logs --no-color --tail=200 >&2 || true
  fi

  $COMPOSE down --remove-orphans >/dev/null || true
  exit "$status"
}
trap cleanup EXIT

request() {
  local method="$1"
  local path="$2"
  local data="${3:-}"
  local response_file
  response_file="$(mktemp)"
  local status

  if [[ -n "$data" ]]; then
    status="$(curl --silent --show-error \
      --output "$response_file" \
      --write-out "%{http_code}" \
      --request "$method" \
      --header "Content-Type: application/json" \
      --data "$data" \
      "$WALLET_URL$path")"
  else
    status="$(curl --silent --show-error \
      --output "$response_file" \
      --write-out "%{http_code}" \
      --request "$method" \
      "$WALLET_URL$path")"
  fi

  if [[ "$status" -lt 200 || "$status" -ge 300 ]]; then
    echo "HTTP $status for $method $path" >&2
    cat "$response_file" >&2
    rm -f "$response_file"
    return 1
  fi

  cat "$response_file"
  rm -f "$response_file"
}

json_get() {
  local json="$1"
  local path="$2"

  node -e "
    const value = JSON.parse(process.argv[1]);
    const path = process.argv[2].split('.');
    let current = value;
    for (const key of path) current = current?.[key];
    if (current === undefined || current === null) process.exit(1);
    process.stdout.write(String(current));
  " "$json" "$path"
}

wait_for_wallet_server() {
  local deadline=$((SECONDS + TIMEOUT_SECONDS))

  until request POST "/miner/wallet?miner_id=1" >/dev/null 2>&1; do
    if (( SECONDS >= deadline )); then
      echo "Timed out waiting for wallet server and miner to become ready." >&2
      return 1
    fi
    sleep "$POLL_INTERVAL_SECONDS"
  done
}

wait_for_balance_at_least() {
  local address="$1"
  local expected="$2"
  local label="$3"
  local deadline=$((SECONDS + TIMEOUT_SECONDS))
  local balance="0"

  while (( SECONDS < deadline )); do
    local response
    response="$(request GET "/wallet/balance?blockchainAddress=$address")"
    balance="$(json_get "$response" "balance")"

    if node -e "process.exit(Number(process.argv[1]) >= Number(process.argv[2]) ? 0 : 1)" "$balance" "$expected"; then
      echo "$label balance is $balance."
      return 0
    fi

    sleep "$POLL_INTERVAL_SECONDS"
  done

  echo "Timed out waiting for $label balance to reach $expected. Last balance: $balance." >&2
  return 1
}

echo "Building and starting Docker stack..."
$COMPOSE up --build -d

echo "Waiting for wallet server and miner wallet..."
wait_for_wallet_server

miner_wallet="$(request POST "/miner/wallet?miner_id=1")"
user_wallet="$(request POST "/user/wallet")"

miner_address="$(json_get "$miner_wallet" "blockchainAddress")"
miner_public_key="$(json_get "$miner_wallet" "publicKey")"
miner_private_key="$(json_get "$miner_wallet" "privateKey")"
user_address="$(json_get "$user_wallet" "blockchainAddress")"

echo "Miner wallet: $miner_address"
echo "User wallet:  $user_address"

echo "Waiting for miner reward to be mined..."
wait_for_balance_at_least "$miner_address" "1" "Miner"

transaction_payload="$(node -e "
  process.stdout.write(JSON.stringify({
    message: 'E2E TRANSFER',
    recipientBlockchainAddress: process.argv[1],
    senderBlockchainAddress: process.argv[2],
    senderPrivateKey: process.argv[3],
    senderPublicKey: process.argv[4],
    value: '1'
  }));
" "$user_address" "$miner_address" "$miner_private_key" "$miner_public_key")"

echo "Sending 1 coin from miner wallet to user wallet..."
request POST "/transaction" "$transaction_payload" >/dev/null

echo "Waiting for transfer to be mined..."
wait_for_balance_at_least "$user_address" "1" "User"

echo "E2E transfer passed."
