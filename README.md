# x402-demo

Production-grade Go implementation of the [x402 payment protocol](https://x402.org/) — HTTP-native micropayments over any EVM blockchain.

Three independently deployable components demonstrate the full x402 payment lifecycle: a **Facilitator** that verifies and settles EIP-3009 payments on-chain, a **Resource Server** that gates API access behind HTTP 402, and a **Client CLI** that signs and submits payments automatically.

Tested and verified on **Base Sepolia** with real USDC transfers.

## How It Works

```
                          x402 Payment Flow
                          ─────────────────

   ┌────────────┐                ┌──────────────────┐                ┌────────────────────┐
   │ Client CLI │  1. GET /api   │  Resource Server  │                │ Facilitator Server │
   │            │ ──────────────>│                   │                │                    │
   │            │                │  "No payment      │                │                    │
   │            │  2. HTTP 402   │   header found"   │                │                    │
   │            │ <──────────────│                   │                │                    │
   │            │  + PAYMENT-    │                   │                │                    │
   │            │    REQUIRED    │                   │                │                    │
   │            │                │                   │                │                    │
   │  Signs     │                │                   │                │                    │
   │  EIP-3009  │  3. GET /api   │                   │  4. POST       │                    │
   │  payload   │ ──────────────>│  Parses header,   │ ──/verify────> │  Recovers signer,  │
   │            │  + PAYMENT-    │  forwards to      │                │  checks balance,   │
   │            │    SIGNATURE   │  facilitator      │  5. POST       │  simulates call    │  ┌───────────┐
   │            │                │                   │ ──/settle────> │                    │  │ EVM Chain │
   │            │                │                   │                │  Builds EIP-1559   │  │           │
   │            │                │                   │                │  tx, calls         │──│ transfer  │
   │            │                │                   │  6. tx hash    │  transferWith-     │  │ WithAuth  │
   │            │  7. HTTP 200   │  Returns API data │ <─────────────│  Authorization     │  │ (USDC)    │
   │            │ <──────────────│  + PAYMENT-       │                │                    │  │           │
   │            │  + response    │    RESPONSE       │                │                    │  └───────────┘
   └────────────┘                └──────────────────┘                └────────────────────┘
```

| Step | What happens |
|------|-------------|
| 1 | Client sends a normal HTTP request to a protected endpoint |
| 2 | Resource Server responds with **HTTP 402** and a `PAYMENT-REQUIRED` header containing accepted schemes, network, amount, and recipient address |
| 3 | Client creates an EIP-3009 `transferWithAuthorization` signature (EIP-712 typed data) and retries with a `PAYMENT-SIGNATURE` header |
| 4-5 | Resource Server delegates verification and on-chain settlement to the Facilitator via `/verify` and `/settle` |
| 6 | Facilitator submits the `transferWithAuthorization` transaction, pays gas, waits for confirmation, and returns the tx hash |
| 7 | Client receives the API response along with a `PAYMENT-RESPONSE` header containing the settlement transaction hash |

## Wallet Roles

The system uses three wallets with distinct roles. Notably, `PAY_TO_ADDRESS` does **not** need a private key:

```
Client Wallet                          PAY_TO Address
(signs EIP-3009)                       (receives USDC)
      │                                      ▲
      │  0.1 USDC (transferWithAuth)         │
      └──────────────────────────────────────┘
                        │
               Facilitator Wallet
               (pays gas only, never touches USDC)
```

| Wallet | Private Key? | Holds | Role |
|--------|:---:|-------|------|
| **Facilitator** | Yes | ETH (gas) | Submits `transferWithAuthorization` tx on-chain |
| **Client** | Yes | USDC | Signs EIP-3009 payment authorizations |
| **PAY_TO** | **No** | Receives USDC | Any EVM address — MetaMask, exchange, multisig, etc. |

## Prerequisites

- **Go 1.24+**
- Two wallets with private keys:
  - **Facilitator wallet** — needs ETH for gas (Base Sepolia: get from [Alchemy Faucet](https://www.alchemy.com/faucets/base-sepolia))
  - **Client wallet** — needs USDC (Base Sepolia: get from [Circle Faucet](https://faucet.circle.com))
- A recipient address (`PAY_TO_ADDRESS`) — any EVM address, no private key needed

## Quick Start

### 1. Clone and build

```bash
git clone https://github.com/example/x402-demo.git
cd x402-demo
make build
```

### 2. Configure environment

```bash
cp .env.example .env
```

Edit `.env` with your values:

```bash
# Facilitator wallet (pays gas for on-chain settlement)
FACILITATOR_PRIVATE_KEY=0x...

# Client wallet (holds USDC, signs payment authorizations)
CLIENT_PRIVATE_KEY=0x...

# Recipient address (receives USDC payments — no private key needed)
PAY_TO_ADDRESS=0x...

# Chain configuration (defaults to Base Sepolia)
NETWORK=eip155:84532
RPC_URL=https://sepolia.base.org
USDC_ADDRESS=0x036CbD53842c5426634e7929541eC2318f3dCF7e
```

### 3. Start the servers

```bash
# Terminal 1 — Facilitator (port 4022)
make run-facilitator

# Terminal 2 — Resource Server (port 4021)
make run-resource
```

### 4. Make a paid API call

```bash
# Terminal 3 — Client
make run-client
```

Expected output:

```
→ GET http://localhost:4021/weather
← 200 OK

💰 Payment Settlement:
   Success:     true
   Transaction: 0x99e49093d0bb2805b2e1097a6c71336c73f5871a4e51ec2dacc733f51faedc24
   Network:     eip155:84532
   Payer:       0x47322Ca28a85B12a7EA64a251Cd8b9Ea1fac037b

Response:
{
  "city": "New York",
  "condition": "Windy",
  "humidity": 86,
  "temperature": 25,
  "timestamp": "2026-03-19T01:16:19Z"
}
```

### 5. Check balances

```bash
go run ./cmd/balance
```

```
==========================================
  Base Sepolia — After Payment
==========================================

Facilitator:   0x23fbdE5A14dFB508502f5A2622f66c0D3B0ab37A
  ETH:  0.000299
  USDC: 0.100000

Client:        0x47322Ca28a85B12a7EA64a251Cd8b9Ea1fac037b
  ETH:  0.000000
  USDC: 19.800000

PAY_TO (you):  0xDBCbC75772954F82d436700cDC4B7c8F434e07F5
  ETH:  0.000000
  USDC: 0.100000
```

### 6. Verify endpoints

```bash
curl http://localhost:4022/health      # Facilitator health
curl http://localhost:4021/health      # Resource Server health
curl -v http://localhost:4021/weather  # 402 without payment
curl http://localhost:4022/supported   # Facilitator capabilities
```

## Protected Endpoints

The Resource Server exposes three payment-gated demo APIs:

| Endpoint | Price (USDC) | Description |
|----------|-------------|-------------|
| `GET /weather` | $0.10 | Random city weather data |
| `GET /joke` | $0.10 | Programming jokes |
| `GET /premium-data` | $0.10 | Mock analytics report |
| `GET /health` | Free | Server health check |

Prices are defined in `internal/server/routes.go` and can be changed per-endpoint.

## Client CLI Flags

```bash
go run ./cmd/client [flags]

  -endpoint string   API path to request (default: $ENDPOINT_PATH or /weather)
  -url string        Resource server URL (default: $RESOURCE_URL)
  -v                 Verbose output (debug logging)
```

Examples:

```bash
# Request a joke
go run ./cmd/client -endpoint /joke

# Request premium data from a specific server
go run ./cmd/client -url http://api.example.com -endpoint /premium-data -v
```

## Configuration Reference

All configuration is via environment variables. Copy `.env.example` to `.env` for local development — the app loads `.env` automatically via [godotenv](https://github.com/joho/godotenv).

| Variable | Used by | Default | Description |
|----------|---------|---------|-------------|
| `FACILITATOR_PRIVATE_KEY` | facilitator | *required* | Hex-encoded private key (with or without `0x` prefix). This wallet pays gas for on-chain settlement. |
| `CLIENT_PRIVATE_KEY` | client | *required* | Hex-encoded private key. This wallet holds USDC and signs EIP-3009 authorizations. |
| `RPC_URL` | facilitator, client | `https://sepolia.base.org` | JSON-RPC endpoint for the target EVM chain. |
| `NETWORK` | all | `eip155:84532` | CAIP-2 network identifier. Determines which chain the system operates on. |
| `USDC_ADDRESS` | resource, client | `0x036CbD53842c5426634e7929541eC2318f3dCF7e` | EIP-3009 compatible token contract address. |
| `FACILITATOR_URL` | resource | *required* | Base URL of the Facilitator Server (e.g., `http://localhost:4022`). |
| `PAY_TO_ADDRESS` | resource | *required* | Ethereum address that receives payments. No private key needed. |
| `FACILITATOR_PORT` | facilitator | `4022` | HTTP listen port. |
| `RESOURCE_PORT` | resource | `4021` | HTTP listen port. |
| `RESOURCE_URL` | client | *required* | Base URL of the Resource Server (e.g., `http://localhost:4021`). |
| `ENDPOINT_PATH` | client | `/weather` | Default API endpoint to call. |
| `LOG_LEVEL` | all | `info` | Log verbosity: `debug`, `info`, `warn`, `error`. |

## Switching Chains

The system is chain-agnostic. To target a different EVM network, change three environment variables:

```bash
# Example: Base Mainnet
NETWORK=eip155:8453
RPC_URL=https://mainnet.base.org
USDC_ADDRESS=0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913
```

### Chain Compatibility

| Chain | EIP-3009 | Status |
|-------|:---:|--------|
| Base Sepolia | Yes | **Verified working** — tested with real txs |
| Base Mainnet | Yes | SDK built-in support |
| Polygon | Yes | SDK built-in support |
| Arbitrum | Yes | SDK built-in support |
| Chiliz Mainnet | **No** | Bridged USDC (ChainPort) — no `transferWithAuthorization` in bytecode |
| Chiliz Spicy | **No** | No USDC deployed |

> **Chiliz Note:** The USDC on Chiliz (`0xa37936F56249965d407E39347528a1A91eB1cbef`) is bridged via Chainport. It is a basic ERC-20 (1,798 bytes) named `"Bridged USDC (ChainPort)"` and does not implement EIP-3009. To use x402 on Chiliz, you would need to deploy a custom EIP-3009 compatible token or use the Permit2 transfer method.

## Docker

```bash
# Start Facilitator + Resource Server
docker compose up --build

# Client (run locally against Docker services)
RESOURCE_URL=http://localhost:4021 make run-client

# Shut down
docker compose down
```

## Project Structure

```
x402-demo/
├── cmd/
│   ├── facilitator/main.go    Facilitator HTTP server entrypoint
│   ├── resource/main.go       Resource HTTP server entrypoint
│   ├── client/main.go         Client CLI entrypoint
│   └── balance/main.go        Wallet balance checker utility
├── internal/
│   ├── config/                Environment variable loading + validation
│   ├── facilserver/           Facilitator HTTP handlers (/verify, /settle, /supported)
│   ├── server/                Resource Server route definitions + API handlers
│   └── signer/                FacilitatorEvmSigner implementation (EIP-1559, EIP-712)
├── pkg/health/                Shared health check response type
├── test/                      Integration tests and fixtures
├── .env.example               Environment variable template
├── Dockerfile                 Multi-stage build (facilitator / resource / client targets)
├── docker-compose.yml         Facilitator + Resource orchestration
└── Makefile                   Build, test, run targets
```

## Technology Stack

| Component | Technology |
|-----------|-----------|
| Language | Go 1.24+ |
| x402 SDK | [coinbase/x402/go](https://github.com/coinbase/x402) (v2.6.0, V2 protocol) |
| EVM Client | [go-ethereum](https://github.com/ethereum/go-ethereum) v1.17 |
| HTTP Framework | [Gin](https://github.com/gin-gonic/gin) v1.12 |
| Payment Scheme | EIP-3009 `transferWithAuthorization` (exact scheme) |
| Signature Standard | EIP-712 Typed Structured Data |
| Transaction Format | EIP-1559 (dynamic fee) |
| Logging | `log/slog` (structured JSON) |
| Config | Environment variables via [godotenv](https://github.com/joho/godotenv) |

## Facilitator API Reference

The Facilitator Server exposes endpoints consumed by Resource Servers (not end users):

### `POST /verify`

Validates a payment payload without executing on-chain. Checks signature recovery, payer balance, timestamp validity, and nonce uniqueness.

```json
// Request
{
  "x402Version": 2,
  "paymentPayload": { "..." },
  "paymentRequirements": { "..." }
}

// Response (200 OK)
{
  "isValid": true,
  "payer": "0x47322Ca28a85B12a7EA64a251Cd8b9Ea1fac037b"
}
```

### `POST /settle`

Executes the payment on-chain by calling `transferWithAuthorization` on the token contract. The Facilitator wallet pays gas.

```json
// Response (200 OK)
{
  "success": true,
  "transaction": "0x99e49093d0bb2805b2e1097a6c71336c73f5871a4e51ec2dacc733f51faedc24",
  "network": "eip155:84532",
  "payer": "0x47322Ca28a85B12a7EA64a251Cd8b9Ea1fac037b"
}
```

### `GET /supported`

Returns the schemes, networks, and protocol versions this Facilitator supports.

```json
{
  "kinds": [
    { "x402Version": 2, "scheme": "exact", "network": "eip155:84532" }
  ],
  "extensions": [],
  "signers": {
    "eip155:*": ["0x23fbdE5A14dFB508502f5A2622f66c0D3B0ab37A"]
  }
}
```

### `GET /health`

```json
{
  "status": "ok",
  "service": "facilitator",
  "network": "eip155:84532",
  "address": "0x23fbdE5A14dFB508502f5A2622f66c0D3B0ab37A"
}
```

## Protocol Details

This implementation uses the **x402 V2 protocol** with the following HTTP headers:

| Header | Direction | Purpose |
|--------|-----------|---------|
| `PAYMENT-REQUIRED` | Server -> Client | Base64-encoded payment requirements (in 402 response) |
| `PAYMENT-SIGNATURE` | Client -> Server | Base64-encoded signed payment payload |
| `PAYMENT-RESPONSE` | Server -> Client | Base64-encoded settlement result (tx hash) |

The payment scheme is `exact` using **EIP-3009** `transferWithAuthorization`:

1. Client signs an EIP-712 typed data message authorizing a token transfer
2. Facilitator calls `transferWithAuthorization(from, to, value, validAfter, validBefore, nonce, v, r, s)` on the token contract
3. Nonces are random 32-byte values (not sequential), allowing concurrent authorizations
4. The Facilitator wallet pays gas; the Client wallet only signs — USDC goes directly from Client to `PAY_TO_ADDRESS`

## Development

```bash
make build           # Compile all binaries
make test            # Run unit tests
make test-integration # Run integration tests (requires testnet)
make lint            # Run golangci-lint
make clean           # Remove compiled binaries
```

Run a single test:

```bash
go test ./internal/config -run TestLoadFacilitator -v
```

## Verified Transactions

The following transactions were executed during testing on Base Sepolia:

| Tx Hash | From | To | Amount |
|---------|------|-----|--------|
| [`0x99e4...dc24`](https://sepolia.basescan.org/tx/0x99e49093d0bb2805b2e1097a6c71336c73f5871a4e51ec2dacc733f51faedc24) | `0x4732...037b` | `0x23fb...b37A` | 0.1 USDC |
| [`0x6d3a...2445`](https://sepolia.basescan.org/tx/0x6d3a230de24f0650703fc87fd9b3f0cb19cc914e6530aca4512d5956f4fb2445) | `0x4732...037b` | `0xDBCb...07F5` | 0.1 USDC |

## Further Reading

- [x402 Protocol Specification](https://x402.org/)
- [x402 Documentation](https://docs.x402.org/)
- [Coinbase x402 Go SDK](https://pkg.go.dev/github.com/coinbase/x402/go)
- [EIP-3009: Transfer With Authorization](https://eips.ethereum.org/EIPS/eip-3009)
- [EIP-712: Typed Structured Data Hashing and Signing](https://eips.ethereum.org/EIPS/eip-712)
- [Coinbase x402 GitHub Repository](https://github.com/coinbase/x402)

## License

MIT
