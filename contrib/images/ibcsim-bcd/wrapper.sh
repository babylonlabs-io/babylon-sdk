#!/usr/bin/env sh
# shellcheck disable=SC3037

# 0. Define configuration
BABYLON_KEY="babylon-key"
BABYLON_CHAIN_ID="chain-test"
CONSUMER_KEY="bcd-key"
CONSUMER_CHAIN_ID="bcd-test"

# 1. Create a bcd testnet with Babylon contract (includes IBC setup)
./setup-bcd.sh $CONSUMER_CHAIN_ID $CONSUMER_CONF 26657 26656 6060 9090 ./babylon_contract.wasm ./btc_light_client.wasm ./btc_staking.wasm ./btc_finality.wasm

sleep 3

echo "bcd started. Status of bcd node:"
bcd status

# 2. Wait for consumer registration and create zoneconcierge channel
echo "Waiting for consumer registration before creating zoneconcierge channel..."
echo "Sleeping for 30 seconds to allow consumer registration by test..."
sleep 30

# Create zoneconcierge channel
CONTRACT_ADDRESS=$(bcd query wasm list-contract-by-code 1 | grep bbnc | cut -d' ' -f2)
CONTRACT_PORT="wasm.$CONTRACT_ADDRESS"
echo "Contract port: $CONTRACT_PORT"

echo "Creating zoneconcierge channel..."
rly --home $RELAYER_CONF_DIR tx channel bcd --src-port zoneconcierge --dst-port $CONTRACT_PORT --order ordered --version zoneconcierge-1
[ $? -eq 0 ] && echo "  ✅ Created zoneconcierge IBC channel successfully!" || echo "  ❌ Error creating zoneconcierge IBC channel"

# 3. Start the IBC relayer
echo "Start the IBC relayer"
rly --home $RELAYER_CONF_DIR start bcd --debug-addr '' --flush-interval 30s > /data/relayer/relayer.log
