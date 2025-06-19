# Babylon SDK

Cosmos module for Cosmos BSN (Bitcoin Supercharged Network) Staking Integration.

This module provides the necessary functionality to integrate Cosmos BSN chains
with the Babylon BTC Staking protocol. It allows Cosmos BSN chains to stake BTC
and receive economic security from Bitcoin, while also providing the necessary
functionality to interact with the Babylon Genesis chain.

Please see the [Cosmos BSN contracts](https://github.com/babylonlabs-io/cosmos-bsn-contracts)
repository for more documentation and CosmWasm contracts.

The code is forked from https://github.com/osmosis-labs/mesh-security-sdk.

## Project Structure

* `x/babylon` - Module code that is to be imported by BSNs.
* `demo/app` - Example application and CLI that is using the babylon module.
* `tests/e2e` - End-to-end tests with the demo app and Cosmos BSN smart contracts.

## High Level Overview

Babylon SDK provides a thin layer that BSNs need to integrate, in the form of a
Cosmos module.

## Code Reference

### Babylon Module

See the [Babylon module README](x/babylon/README.md) for more information.

### GRPC Queries

See the `demo/app/wasm/grpc_whitelist.go` file for the `WhitelistedGrpcQuery`
function, which defines the list of whitelisted GRPC queries.