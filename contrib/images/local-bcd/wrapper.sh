#!/usr/bin/env sh

# 0. Define configuration
CONSUMER_KEY="bcd-key"
CONSUMER_CHAIN_ID="bcd-test"

# 1. Create a bcd testnet with Babylon contract
./setup-bcd.sh $CONSUMER_CHAIN_ID $CONSUMER_CONF 26657 26656 6060 9090 ./babylon_contract.wasm ./btc_staking.wasm '{
    "network": "regtest",
    "babylon_tag": "01020304",
    "btc_confirmation_depth": 1,
    "checkpoint_finalization_timeout": 2,
    "notify_cosmos_zone": false,
    "btc_staking_code_id": 2,
    "consumer_name": "Test Consumer",
    "consumer_description": "Test Consumer Description"
}'

sleep 10

CONTRACT_ADDRESS=$(bcd query wasm list-contract-by-code 1 | grep bbnc | cut -d' ' -f2)
CONTRACT_PORT="wasm.$CONTRACT_ADDRESS"
echo "bcd started. Status of bcd node:"
bcd status
echo "Contract port: $CONTRACT_PORT"
