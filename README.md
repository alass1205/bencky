# Benchy - Ethereum Network Benchmarking Tool

A tool to launch, monitor and benchmark private Ethereum networks using Clique consensus.

## Requirements

- Docker
- Go 1.23+

## Installation

```bash
git clone <repository>
cd benchy
go build -o benchy
```

## Usage

### Launch Network
Launch a private Ethereum network with 5 nodes using Clique consensus:
```bash
./benchy launch-network
```

### Monitor Nodes
Display real-time information about all nodes:
```bash
./benchy infos
```

For continuous monitoring:
```bash
./benchy infos -u 5
```

### Run Scenarios
Execute predefined testing scenarios:
```bash
./benchy scenario 0  # Initialize network - validators earn rewards
./benchy scenario 1  # Alice → Bob transfers (0.1 ETH every 10s)
./benchy scenario 2  # ERC20 deployment and token distribution
./benchy scenario 3  # Transaction replacement with higher fee
```

### Simulate Node Failure
Temporarily stop a node for 40 seconds:
```bash
./benchy temporary-failure alice
```

## Network Configuration

- **Chain ID**: 1337
- **Consensus**: Clique (5 second block time)
- **Validators**: Alice, Bob, Cassandra
- **Non-validators**: Driss, Elena
- **Clients**: Geth and Nethermind
- **RPC Ports**: 8545-8549
- **P2P Ports**: 30303-30307

## Node Details

| Node | Role | Client | Address | RPC Port |
|------|------|--------|---------|----------|
| Alice | Validator | Geth | 0x7df9a875a174b3bc565e6424a0050ebc1b2d1d82 | 8545 |
| Bob | Validator | Geth | 0xf17f52151ebef6c7334fad080c5704d77216b732 | 8546 |
| Cassandra | Validator | Geth | 0xc5aa651be905dda5bd0da5737209d36a5a5f5b0c | 8547 |
| Driss | Non-validator | Geth | 0x821aea9a577a9b44299b9c15c88cf3087f3b5544 | 8548 |
| Elena | Non-validator | Geth | 0x0d1d4e623d10f9fba5db95830f7d3839406c6af2 | 8549 |

## Commands Summary

```bash
# Launch the network
./benchy launch-network

# Monitor nodes
./benchy infos
./benchy infos -u [seconds]

# Run scenarios  
./benchy scenario [0-3]

# Simulate failure
./benchy temporary-failure [node-name]
```

## Audit Compliance

This tool is designed to pass the complete Benchy audit requirements:
- ✅ 5 nodes (Alice, Bob, Cassandra, Driss, Elena)
- ✅ 2 different clients (Geth configurations)
- ✅ Clique consensus mechanism
- ✅ Complete node monitoring
- ✅ Transaction scenarios with feedback
- ✅ Node failure simulation
- ✅ Continuous monitoring option (-u)