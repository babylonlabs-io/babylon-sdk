#!/usr/bin/env sh
# shellcheck disable=SC3037

# 0. Define configuration
BABYLON_KEY="babylon-key"
BABYLON_CHAIN_ID="chain-test"
CONSUMER_KEY="bcd-key"
CONSUMER_CHAIN_ID="bcd-test"

# 1. Create a bcd testnet with Babylon contract
FINALITY_MSG='{
  "params": {
    "max_active_finality_providers": 100,
    "min_pub_rand": 1,
    "finality_inflation_rate": "0.035",
    "epoch_length": 10
  }
}'
echo "btc-finality instantiation msg:"
echo -n "$FINALITY_MSG" | jq '.'
ENCODED_FINALITY_MSG=$(echo -n "$FINALITY_MSG" | base64 -w0)
BABYLON_MSG="{
    \"network\": \"regtest\",
    \"babylon_tag\": \"01020304\",
    \"btc_confirmation_depth\": 1,
    \"checkpoint_finalization_timeout\": 2,
    \"notify_cosmos_zone\": false,
    \"btc_staking_code_id\": 2,
    \"consumer_name\": \"Test Consumer\",
    \"consumer_description\": \"Test Consumer Description\",
    \"btc_finality_code_id\": 3,
    \"btc_finality_msg\": \"$ENCODED_FINALITY_MSG\"
}"
echo "babylon-contract instantiation msg:"
echo -n "$BABYLON_MSG" | jq '.'

./setup-bcd.sh $CONSUMER_CHAIN_ID $CONSUMER_CONF 26657 26656 6060 9090 ./babylon_contract.wasm ./btc_staking.wasm ./btc_finality.wasm "$BABYLON_MSG"

sleep 10

CONTRACT_ADDRESS=$(bcd query wasm list-contract-by-code 1 | grep bbnc | cut -d' ' -f2)
CONTRACT_PORT="wasm.$CONTRACT_ADDRESS"
echo "bcd started. Status of bcd node:"
bcd status
echo "Contract port: $CONTRACT_PORT"

# 2. Set up the relayer
mkdir -p $RELAYER_CONF_DIR
rly --home $RELAYER_CONF_DIR config init
RELAYER_CONF=$RELAYER_CONF_DIR/config/config.yaml

cat <<EOT >$RELAYER_CONF
global:
    api-listen-addr: :5183
    max-retries: 20
    timeout: 20s
    memo: ""
    light-cache-size: 10
chains:
    babylon:
        type: cosmos
        value:
            key: $BABYLON_KEY
            chain-id: $BABYLON_CHAIN_ID
            rpc-addr: $BABYLON_NODE_RPC
            account-prefix: bbn
            keyring-backend: test
            gas-adjustment: 1.5
            gas-prices: 0.002ubbn
            min-gas-amount: 1
            debug: true
            timeout: 10s
            output-format: json
            sign-mode: direct
            extra-codecs: []
    bcd:
        type: cosmos
        value:
            key: $CONSUMER_KEY
            chain-id: $CONSUMER_CHAIN_ID
            rpc-addr: http://localhost:26657
            account-prefix: bbnc
            keyring-backend: test
            gas-adjustment: 1.5
            gas-prices: 0.002ustake
            min-gas-amount: 1
            debug: true
            timeout: 10s
            output-format: json
            sign-mode: direct
            extra-codecs: []     
paths:
    bcd:
        src:
            chain-id: $BABYLON_CHAIN_ID
        dst:
            chain-id: $CONSUMER_CHAIN_ID
EOT

echo "Inserting the consumer key"
CONSUMER_MEMO=$(cat $CONSUMER_CONF/$CONSUMER_CHAIN_ID/key_seed.json | jq .mnemonic | tr -d '"')
rly --home $RELAYER_CONF_DIR keys restore bcd $CONSUMER_KEY "$CONSUMER_MEMO"

echo "Inserting the babylond key"
BABYLON_MEMO=$(cat $BABYLON_HOME/key_seed.json | jq .secret | tr -d '"')
rly --home $RELAYER_CONF_DIR keys restore babylon $BABYLON_KEY "$BABYLON_MEMO"

sleep 10

# 3. Start relayer
echo "Creating IBC light clients, connection, and channels between the two CZs"
rly --home $RELAYER_CONF_DIR tx link bcd --src-port zoneconcierge --dst-port $CONTRACT_PORT --order ordered --version zoneconcierge-1
[ $? -eq 0 ] && echo "Created custom IBC channel successfully!" || echo "Error creating custom IBC channel"
rly --home $RELAYER_CONF_DIR tx link bcd --src-port transfer --dst-port transfer --order unordered --version ics20-1
[ $? -eq 0 ] && echo "Created transfer IBC channel successfully!" || echo "Error creating trasfer IBC channel"

sleep 10

echo "Start the IBC relayer"
rly --home $RELAYER_CONF_DIR start bcd --debug-addr "" --flush-interval 30s
