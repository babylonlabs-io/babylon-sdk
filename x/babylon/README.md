# Babylon Cosmos Module

The `babylon` module provides the following functionalities:

- Exposing the `BeginBlock` and `EndBlock` interfaces to the Cosmos BSN
  smart contracts.
- Handling rewards generation for BSNs, i.e., allocating a part of the protocol
  revenue to BTC and BABY stakers.
- Providing wrapper functions for simplifying the instantiation of Cosmos BSN
  contracts.
- Exposing GRPC queries, i.e. whitelisting some queries for the Cosmos BSN smart
  contracts to use.

## Code Reference

### Begin and End Block Messages

See the `x/babylon/keeper/wasm.go` file for the `SendBeginBlockMsg` and
`SendEndBlockMsg` functions, which are used to send the `BeginBlock` and
`EndBlock` messages to the Cosmos BSN smart contracts.

### Rewards Generation

See the `x/babylon/keeper/mint_rewards.go` file for the `MintBlockRewards`
function, which is used to generate the rewards for the block, and send them to
the `btc-finality` smart contract.

### Contract Instantiation

See the `x/babylon/keeper/msg_server.go` file for the
`InstantiateBabylonContracts` function, which is used to instantiate the
Cosmos BSN smart contracts with the necessary parameters.