# SERV Chain

A custom Proof-of-Stake blockchain built with Cosmos SDK and Tendermint Core, featuring innovative Proof-of-Service and SERV Rewards modules.

## Overview

SERV Chain combines traditional Proof-of-Stake consensus with a novel Proof-of-Service mechanism to reward validators not only for securing the network but also for providing valuable services to the ecosystem.

### Key Components

- **Tendermint Core**: Provides Byzantine Fault Tolerant consensus and P2P networking
- **Cosmos SDK**: Modular blockchain framework
- **Custom Modules**:
  - `x/servrewards`: Tracks, calculates, and distributes rewards based on system engagement
  - `x/proofofservice`: Verifies and validates service provision by nodes/operators
  - `x/noderewards` (optional): Augments standard staking rewards with performance-based logic

## Architecture

```
+-----------------------------+
|  Clients & Interfaces       | ← Wallets, CLI, REST/gRPC, Frontends
+-----------------------------+
            |
            v
+-----------------------------+
|  Full Nodes & APIs          | ← Tendermint RPC, gRPC, REST
|  (with mempool)             |
+-----------------------------+
            |
            v
+-----------------------------+
|     Tendermint Core         | ← Consensus & Networking
| - P2P Networking             |
| - BFT Consensus              |
| - ABCI (App ↔ Consensus)     |
+-----------------------------+
            |
            v
+-----------------------------+
|     Cosmos SDK App          |
| - baseapp                   |
| - Module Manager            |
| - Custom + Standard Modules |
+-----------------------------+
            |
            v
+-----------------------------+
| Persistent Storage (DB)     | ← LevelDB / RocksDB
+-----------------------------+
```

## Getting Started

### Prerequisites

- Go 1.19+
- Docker and Docker Compose (for local development)

### Installation

```bash
# Clone the repository
git clone https://github.com/your-org/serv-chain.git
cd serv-chain

# Install dependencies
go mod tidy

# Build the application
make build

# Run a local development node
make start
```

## Development

### Module Structure

The application includes standard Cosmos SDK modules (`auth`, `bank`, `staking`, etc.) and custom modules:

- `x/servrewards`: Manages reward distribution based on service provision
- `x/proofofservice`: Validates and scores service provision
- `x/noderewards`: (Optional) Enhances validator rewards based on service scores

### Commands

```bash
# Initialize a new chain
servchaind init [moniker] --chain-id serv-chain-1

# Create a key
servchaind keys add [name]

# Add genesis account
servchaind add-genesis-account [address] 10000000serv

# Create validator transaction
servchaind gentx [name] 1000000serv --chain-id serv-chain-1

# Collect genesis transactions
servchaind collect-gentxs

# Start the node
servchaind start
```

## License

[MIT License](LICENSE)
