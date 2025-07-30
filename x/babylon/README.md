# Babylon Module

The `babylon` module in the Babylon SDK provides BSN integration capabilities
for Cosmos SDK based chains.
This module serves as a bridge between Cosmos BSNs and
the Babylon Genesis Bitcoin staking infrastructure,
enabling seamless integration to become a BSN (Bitcoin Supercharged Network).

## Table of contents

* [Table of contents](#table-of-contents)
* [Concepts](#concepts)
  * [Cosmos BSN Integration](#cosmos-bsn-integration)
  * [Smart Contract Communication](#smart-contract-communication)
* [States](#states)
  * [Parameters](#parameters)
* [Messages](#messages)
  * [MsgSetBSNContracts](#msgsetbsncontracts)
  * [MsgUpdateParams](#msgupdateparams)
* [BeginBlocker](#beginblocker)
* [EndBlocker](#endblocker)
* [Events](#events)
* [Queries](#queries)
* [Contract Integration](#contract-integration)
  * [Out-Messages](#out-messages)

## Concepts

The Cosmos BSNs integration stack enables Cosmos SDK based chains
to integrate with Babylon's Bitcoin staking infrastructure through a set of
smart contracts that handle Bitcoin staking, finality, and
light client functionality.

### Cosmos BSN Integration Stack

The Cosmos BSN integration stack consists of several CosmWasm smart contracts
that work together to provide Bitcoin staking capabilities:

* **Babylon Contract**: The main orchestrator contract that coordinates all
  other contracts.
* **BTC Light Client Contract**: Maintains Bitcoin header information and
  validates Bitcoin transactions.
* **BTC Staking Contract**: Manages Bitcoin staking operations and delegations.
* **BTC Finality Contract**: Handles finality voting and reward distribution.

The `babylon` module in this repository provides the necessary infrastructure
to instantiate and communicate with these contracts from Cosmos layer.

### Smart Contract Communication

The module communicates with smart contracts through one main mechanism:

1. **Sudo Messages**: Messages sent from the module to smart contracts during
   block processing

This communication enables the module to:
- Send block information to contracts during `BeginBlock` and `EndBlock`
- Coordinate contract instantiation and configuration

## States

The Babylon SDK module maintains the following state information:

### Parameters

The module parameters are defined in the `Params` protobuf message and include:

```protobuf
message Params {
  // Gas limits
  uint32 max_gas_begin_blocker = 1;
}
```

The parameters are managed through the `x/babylon/keeper/params.go` file and include:

* **Gas Limits**: Maximum gas allowed for contract sudo callbacks

### Genesis State

The module's genesis state includes the following fields for contract addresses:

```protobuf
message GenesisState {
  Params params = 1;
  BSNContracts bsn_contracts = 2;
}

message BSNContracts {
  string babylon_contract = 1;
  string btc_light_client_contract = 2;
  string btc_staking_contract = 3;
  string btc_finality_contract = 4;
}
```

* **Contract Addresses**: All Cosmos BSN smart contract addresses are now grouped in a single `BSNContracts` object. These addresses are set at chain genesis and used for communication with the respective contracts during block processing and other operations.

To set contract addresses at chain start, specify them in the genesis file under the `babylon` module's state as a `bsn_contracts` object. If not set, they can be set later via the `SetBSNContracts` message.

## Messages

The `babylon` module handles the following messages:

### MsgSetBSNContracts

Sets the addresses of the Cosmos BSN smart contracts in the module state. This message can be submitted by governance or another authorized entity to update the contract addresses after genesis.

```protobuf
message MsgSetBSNContracts {
  string authority = 1;
  BSNContracts contracts = 2;
}

message BSNContracts {
  string babylon_contract = 1;
  string btc_light_client_contract = 2;
  string btc_staking_contract = 3;
  string btc_finality_contract = 4;
}
```

**Parameters:**
- `authority`: Address with authority to set contract addresses (usually x/gov)
- `contracts`: A `BSNContracts` object containing all contract addresses

All contract addresses must be valid Bech32 addresses. The module validates the entire `BSNContracts` object atomically.

### MsgUpdateParams

Updates the module parameters. Only the authority can execute this message.

```protobuf
message MsgUpdateParams {
  string authority = 1;
  Params params = 2;
}
```

**Parameters:**
- `authority`: Address with authority to update parameters
- `params`: New parameter values

## BeginBlocker

The `BeginBlocker` is executed at the beginning of each block and
sends `BeginBlock` sudo messages to the BTC staking and finality contracts
containing the current block hash and app hash.

## EndBlocker

The `EndBlocker` is executed at the end of each block and sends
`EndBlock` sudo messages to the BTC finality contract containing
the current block hash and app hash.

## Events

The module emits events for various operations:

- **Contract Communication**: Events when messages are sent to contracts
- **Parameter Updates**: Events when module parameters are updated

Event definitions are located in `x/babylon/types/events.go`.

## Queries

The module provides the following query endpoints:

### QueryParams

Retrieves the current module parameters.

```protobuf
message QueryParamsRequest {}

message QueryParamsResponse {
  Params params = 1;
}
```

**Usage:**
```bash
babylond query babylon params
```

### QueryBSNContracts

Retrieves all BSN contract addresses as a single object.

```protobuf
message QueryBSNContractsRequest {}

message QueryBSNContractsResponse {
  BSNContracts bsn_contracts = 1;
}
```

**Usage:**
```bash
babylond query babylon bsn-contracts
```

## Contract Integration

The module integrates with the Cosmos BSN contracts through outbound messages:

### Out-Messages

Messages sent from the module to smart contracts:

#### BeginBlock

Sent to contracts at the beginning of each block:

```go
type BeginBlock struct {
    HashHex    string `json:"hash_hex"`
    AppHashHex string `json:"app_hash_hex"`
}
```

#### EndBlock

Sent to contracts at the end of each block:

```go
type EndBlock struct {
    HashHex    string `json:"hash_hex"`
    AppHashHex string `json:"app_hash_hex"`
}
```

## Architecture Notes

### Security

The module includes an **IntegrityHandler** that prevents smart contracts from bypassing
virtual staking mechanisms by blocking unauthorized staking operations. This ensures
that all staking operations go through the proper BTC staking infrastructure.