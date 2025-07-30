#!/bin/bash

display_usage() {
	echo "Missing parameters. Please check if all parameters were specified."
	echo "Usage: setup-bcd.sh [CHAIN_ID] [CHAIN_DIR] [RPC_PORT] [P2P_PORT] [PROFILING_PORT] [GRPC_PORT] [BABYLON_CONTRACT_CODE_FILE] [BTC_LC_CONTRACT_CODE_FILE] [BTCSTAKING_CONTRACT_CODE_FILE] [BTCFINALITY_CONTRACT_CODE_FILE]"
	echo "Example: setup-bcd.sh test-chain-id ./data 26657 26656 6060 9090 ./babylon_contract.wasm ./btc_light_client.wasm ./btc_staking.wasm ./btc_finality.wasm"
	exit 1
}

BINARY=bcd
DENOM=stake
BASEDENOM=ustake
KEYRING=--keyring-backend="test"
SILENT=1

redirect() {
	if [ "$SILENT" -eq 1 ]; then
		"$@" >/dev/null 2>&1
	else
		"$@"
	fi
}

if [ "$#" -lt "9" ]; then
	display_usage
	exit 1
fi

CHAINID=$1
CHAINDIR=$2
RPCPORT=$3
P2PPORT=$4
PROFPORT=$5
GRPCPORT=$6
BABYLON_CONTRACT_CODE_FILE=$7
BTC_LC_CONTRACT_CODE_FILE=$8
BTCSTAKING_CONTRACT_CODE_FILE=$9
BTCFINALITY_CONTRACT_CODE_FILE=${10}

# ensure the binary exists
if ! command -v $BINARY &>/dev/null; then
	echo "$BINARY could not be found"
	exit
fi

# Delete chain data from old runs
echo "Deleting $CHAINDIR/$CHAINID folders..."
rm -rf $CHAINDIR/$CHAINID &>/dev/null
rm $CHAINDIR/$CHAINID.log &>/dev/null

echo "Creating $BINARY instance: home=$CHAINDIR | chain-id=$CHAINID | p2p=:$P2PPORT | rpc=:$RPCPORT | profiling=:$PROFPORT | grpc=:$GRPCPORT"

# Add dir for chain, exit if error
if ! mkdir -p $CHAINDIR/$CHAINID 2>/dev/null; then
	echo "Failed to create chain folder. Aborting..."
	exit 1
fi
# Build genesis file incl account for passed address
coins="100000000000$DENOM,100000000000$BASEDENOM"
delegate="50000000000$DENOM"

redirect $BINARY --home $CHAINDIR/$CHAINID --chain-id $CHAINID init $CHAINID
$BINARY --home $CHAINDIR/$CHAINID keys add validator $KEYRING --output json >$CHAINDIR/$CHAINID/validator_seed.json 2>&1
$BINARY --home $CHAINDIR/$CHAINID keys add user $KEYRING --output json >$CHAINDIR/$CHAINID/key_seed.json 2>&1
redirect $BINARY --home $CHAINDIR/$CHAINID genesis add-genesis-account $($BINARY --home $CHAINDIR/$CHAINID keys $KEYRING show user -a) $coins
redirect $BINARY --home $CHAINDIR/$CHAINID genesis add-genesis-account $($BINARY --home $CHAINDIR/$CHAINID keys $KEYRING show validator -a) $coins
redirect $BINARY --home $CHAINDIR/$CHAINID genesis gentx validator $delegate $KEYRING --chain-id $CHAINID
redirect $BINARY --home $CHAINDIR/$CHAINID genesis collect-gentxs

# Set proper defaults and change ports
echo "Change settings in config.toml and genesis.json files..."
# Use temporary files to avoid permission issues with mounted volumes
cp $CHAINDIR/$CHAINID/config/config.toml /tmp/config.toml.tmp
sed 's#"tcp://127.0.0.1:26657"#"tcp://0.0.0.0:'"$RPCPORT"'"#g' /tmp/config.toml.tmp > /tmp/config1.tmp && mv /tmp/config1.tmp /tmp/config.toml.tmp
sed 's#"tcp://0.0.0.0:26656"#"tcp://0.0.0.0:'"$P2PPORT"'"#g' /tmp/config.toml.tmp > /tmp/config1.tmp && mv /tmp/config1.tmp /tmp/config.toml.tmp
sed 's#"localhost:6060"#"localhost:'"$PROFPORT"'"#g' /tmp/config.toml.tmp > /tmp/config1.tmp && mv /tmp/config1.tmp /tmp/config.toml.tmp
sed 's/timeout_commit = "5s"/timeout_commit = "1s"/g' /tmp/config.toml.tmp > /tmp/config1.tmp && mv /tmp/config1.tmp /tmp/config.toml.tmp
sed 's/max_body_bytes = 1000000/max_body_bytes = 1000000000/g' /tmp/config.toml.tmp > /tmp/config1.tmp && mv /tmp/config1.tmp /tmp/config.toml.tmp
sed 's/timeout_propose = "3s"/timeout_propose = "1s"/g' /tmp/config.toml.tmp > /tmp/config1.tmp && mv /tmp/config1.tmp /tmp/config.toml.tmp
sed 's/index_all_keys = false/index_all_keys = true/g' /tmp/config.toml.tmp > /tmp/config1.tmp && mv /tmp/config1.tmp /tmp/config.toml.tmp
cp /tmp/config.toml.tmp $CHAINDIR/$CHAINID/config/config.toml

cp $CHAINDIR/$CHAINID/config/app.toml /tmp/app.toml.tmp
sed 's/minimum-gas-prices = ""/minimum-gas-prices = "0.00001ustake"/g' /tmp/app.toml.tmp > /tmp/app1.tmp && mv /tmp/app1.tmp /tmp/app.toml.tmp
sed 's#"tcp://0.0.0.0:1317"#"tcp://0.0.0.0:1318"#g' /tmp/app.toml.tmp > /tmp/app1.tmp && mv /tmp/app1.tmp /tmp/app.toml.tmp # ensure port is not conflicted with Babylon
cp /tmp/app.toml.tmp $CHAINDIR/$CHAINID/config/app.toml

cp $CHAINDIR/$CHAINID/config/genesis.json /tmp/genesis.json.tmp
sed 's/"bond_denom": "stake"/"bond_denom": "'"$DENOM"'"/g' /tmp/genesis.json.tmp > /tmp/genesis1.tmp && mv /tmp/genesis1.tmp /tmp/genesis.json.tmp
cp /tmp/genesis.json.tmp $CHAINDIR/$CHAINID/config/genesis.json

# Clean up temporary files
rm -f /tmp/config.toml.tmp /tmp/config1.tmp /tmp/app.toml.tmp /tmp/app1.tmp /tmp/genesis.json.tmp /tmp/genesis1.tmp

# sed -i '' 's#index-events = \[\]#index-events = \["message.action","send_packet.packet_src_channel","send_packet.packet_sequence"\]#g' $CHAINDIR/$CHAINID/config/app.toml

# Modify governance parameters for faster testing
echo "Updating governance parameters for faster testing..."

# Use temporary files to avoid permission issues with mounted volumes
GENESIS_TEMP="/tmp/genesis.json.tmp"
GENESIS_WORK="/tmp/genesis_work.tmp"
cp $CHAINDIR/$CHAINID/config/genesis.json "$GENESIS_TEMP"

# Apply governance parameter modifications
sed 's/"voting_period": "[^"]*"/"voting_period": "60s"/g' "$GENESIS_TEMP" > "$GENESIS_WORK" && mv "$GENESIS_WORK" "$GENESIS_TEMP"
sed 's/"amount": "10000000"/"amount": "1000000"/g' "$GENESIS_TEMP" > "$GENESIS_WORK" && mv "$GENESIS_WORK" "$GENESIS_TEMP"
sed 's/"max_deposit_period": "[^"]*"/"max_deposit_period": "30s"/g' "$GENESIS_TEMP" > "$GENESIS_WORK" && mv "$GENESIS_WORK" "$GENESIS_TEMP"

# Copy the modified file back and clean up
cp "$GENESIS_TEMP" $CHAINDIR/$CHAINID/config/genesis.json
rm -f "$GENESIS_TEMP" "$GENESIS_WORK"
# Start
echo "Starting $BINARY..."
$BINARY --home $CHAINDIR/$CHAINID start --pruning=nothing --grpc-web.enable=false --grpc.address="0.0.0.0:$GRPCPORT" --log_level trace --trace --log_format 'plain' --log_no_color 2>&1 | tee $CHAINDIR/$CHAINID.log &
sleep 10

# Set up relayer and create IBC channels BEFORE uploading contracts
echo "Setting up relayer and creating IBC channels..."
BABYLON_KEY="babylon-key"
BABYLON_CHAIN_ID="chain-test"
CONSUMER_KEY="bcd-key"
RELAYER_CONF_DIR="/data/relayer"
BABYLON_HOME="/data/node1/babylond"
BABYLON_NODE_RPC="http://babylondnode1:26657"

mkdir -p $RELAYER_CONF_DIR
rly --home $RELAYER_CONF_DIR config init
RELAYER_CONF=$RELAYER_CONF_DIR/config/config.yaml

cat <<EOT >$RELAYER_CONF
global:
    api-listen-addr: :5183
    max-retries: 20
    timeout: 30s
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
            timeout: 30s
            output-format: json
            sign-mode: direct
            extra-codecs: []
    bcd:
        type: cosmos
        value:
            key: $CONSUMER_KEY
            chain-id: $CHAINID
            rpc-addr: http://localhost:$RPCPORT
            account-prefix: bbnc
            keyring-backend: test
            gas-adjustment: 1.5
            gas-prices: 0.002ustake
            min-gas-amount: 1
            debug: true
            timeout: 30s
            output-format: json
            sign-mode: direct
            extra-codecs: []
paths:
    bcd:
        src:
            chain-id: $BABYLON_CHAIN_ID
        dst:
            chain-id: $CHAINID
EOT

echo "Inserting relayer keys..."
if [ ! -f "$CHAINDIR/$CHAINID/key_seed.json" ]; then
    echo "ERROR: Consumer key file not found!"
    exit 1
fi

CONSUMER_MEMO=$(cat $CHAINDIR/$CHAINID/key_seed.json | jq .mnemonic | tr -d '"')
if ! rly --home $RELAYER_CONF_DIR keys list bcd 2>/dev/null | grep -q "$CONSUMER_KEY"; then
    rly --home $RELAYER_CONF_DIR keys restore bcd $CONSUMER_KEY "$CONSUMER_MEMO"
fi

if [ ! -f "$BABYLON_HOME/key_seed.json" ]; then
    echo "ERROR: Babylon key file not found!"
    exit 1
fi

BABYLON_MEMO=$(cat $BABYLON_HOME/key_seed.json | jq .secret | tr -d '"')
if ! rly --home $RELAYER_CONF_DIR keys list babylon 2>/dev/null | grep -q "$BABYLON_KEY"; then
    rly --home $RELAYER_CONF_DIR keys restore babylon $BABYLON_KEY "$BABYLON_MEMO"
fi

sleep 5

# Create IBC infrastructure
echo "Creating IBC clients..."
rly --home $RELAYER_CONF_DIR tx clients bcd
sleep 10

echo "Creating IBC connection..."
rly --home $RELAYER_CONF_DIR tx connection bcd
sleep 5

echo "Creating IBC transfer channel..."
rly --home $RELAYER_CONF_DIR tx channel bcd --src-port transfer --dst-port transfer --order unordered --version ics20-1
sleep 3

# Upload contract code and capture code IDs
echo "Storing Babylon contract code..."
$BINARY --home $CHAINDIR/$CHAINID tx wasm store "$BABYLON_CONTRACT_CODE_FILE" $KEYRING --from user --chain-id $CHAINID --gas 200000000 --gas-prices 0.01$BASEDENOM --node http://localhost:$RPCPORT -y
sleep 3
BABYLON_CODE_ID=$($BINARY --home $CHAINDIR/$CHAINID query wasm list-code --output json | jq -r '.code_infos[-1].code_id')
echo "BABYLON_CODE_ID: $BABYLON_CODE_ID"

echo "Storing BTC Light Client contract code..."
$BINARY --home $CHAINDIR/$CHAINID tx wasm store "$BTC_LC_CONTRACT_CODE_FILE" $KEYRING --from user --chain-id $CHAINID --gas 200000000 --gas-prices 0.01$BASEDENOM --node http://localhost:$RPCPORT -y
sleep 3
BTC_LC_CODE_ID=$($BINARY --home $CHAINDIR/$CHAINID query wasm list-code --output json | jq -r '.code_infos[-1].code_id')
echo "BTC_LC_CODE_ID: $BTC_LC_CODE_ID"

echo "Storing BTC Staking contract code..."
$BINARY --home $CHAINDIR/$CHAINID tx wasm store "$BTCSTAKING_CONTRACT_CODE_FILE" $KEYRING --from user --chain-id $CHAINID --gas 200000000 --gas-prices 0.01$BASEDENOM --node http://localhost:$RPCPORT -y
sleep 3
BTCSTAKING_CODE_ID=$($BINARY --home $CHAINDIR/$CHAINID query wasm list-code --output json | jq -r '.code_infos[-1].code_id')
echo "BTCSTAKING_CODE_ID: $BTCSTAKING_CODE_ID"

echo "Storing BTC Finality contract code..."
$BINARY --home $CHAINDIR/$CHAINID tx wasm store "$BTCFINALITY_CONTRACT_CODE_FILE" $KEYRING --from user --chain-id $CHAINID --gas 200000000 --gas-prices 0.01$BASEDENOM --node http://localhost:$RPCPORT -y
sleep 3
BTCFINALITY_CODE_ID=$($BINARY --home $CHAINDIR/$CHAINID query wasm list-code --output json | jq -r '.code_infos[-1].code_id')
echo "BTCFINALITY_CODE_ID: $BTCFINALITY_CODE_ID"

# Prepare init messages for the other contracts
ADMIN=$($BINARY --home $CHAINDIR/$CHAINID keys show user --keyring-backend test -a)
NETWORK="regtest"
BTC_CONFIRMATION_DEPTH=1
CHECKPOINT_FINALIZATION_TIMEOUT=2
BABYLON_TAG="01020304"
CONSUMER_NAME="test-consumer"
CONSUMER_DESCRIPTION="test-consumer-description"
ICS20_CHANNEL_ID="channel-1"

BTC_LC_INIT_MSG=$(jq -n --arg network "$NETWORK" --argjson btc_confirmation_depth $BTC_CONFIRMATION_DEPTH --argjson checkpoint_finalization_timeout $CHECKPOINT_FINALIZATION_TIMEOUT '{network: $network, btc_confirmation_depth: $btc_confirmation_depth, checkpoint_finalization_timeout: $checkpoint_finalization_timeout}')
BTCSTAKING_INIT_MSG=$(jq -n --arg admin "$ADMIN" '{admin: $admin}')
BTCFINALITY_INIT_MSG=$(jq -n --arg admin "$ADMIN" '{admin: $admin}')



# Base64 encode the init messages as required by the Babylon contract
BTC_LC_INIT_MSG_B64=$(echo -n "$BTC_LC_INIT_MSG" | base64 | tr -d '\n')
BTCSTAKING_INIT_MSG_B64=$(echo -n "$BTCSTAKING_INIT_MSG" | base64 | tr -d '\n')
BTCFINALITY_INIT_MSG_B64=$(echo -n "$BTCFINALITY_INIT_MSG" | base64 | tr -d '\n')

# Build the Babylon contract instantiation message
BABYLON_INIT_MSG=$(jq -n \
  --arg network "$NETWORK" \
  --arg babylon_tag "$BABYLON_TAG" \
  --argjson btc_confirmation_depth $BTC_CONFIRMATION_DEPTH \
  --argjson checkpoint_finalization_timeout $CHECKPOINT_FINALIZATION_TIMEOUT \
  --argjson notify_cosmos_zone false \
  --argjson btc_light_client_code_id $BTC_LC_CODE_ID \
  --arg btc_light_client_msg "$BTC_LC_INIT_MSG_B64" \
  --argjson btc_staking_code_id $BTCSTAKING_CODE_ID \
  --arg btc_staking_msg "$BTCSTAKING_INIT_MSG_B64" \
  --argjson btc_finality_code_id $BTCFINALITY_CODE_ID \
  --arg btc_finality_msg "$BTCFINALITY_INIT_MSG_B64" \
  --arg consumer_name "$CONSUMER_NAME" \
  --arg consumer_description "$CONSUMER_DESCRIPTION" \
  --arg ics20_channel_id "$ICS20_CHANNEL_ID" \
  --arg admin "$ADMIN" \
  '{network: $network, babylon_tag: $babylon_tag, btc_confirmation_depth: $btc_confirmation_depth, checkpoint_finalization_timeout: $checkpoint_finalization_timeout, notify_cosmos_zone: $notify_cosmos_zone, btc_light_client_code_id: $btc_light_client_code_id, btc_light_client_msg: $btc_light_client_msg, btc_staking_code_id: $btc_staking_code_id, btc_staking_msg: $btc_staking_msg, btc_finality_code_id: $btc_finality_code_id, btc_finality_msg: $btc_finality_msg, consumer_name: $consumer_name, consumer_description: $consumer_description, ics20_channel_id: $ics20_channel_id, admin: $admin}')
echo "Babylon contract instantiation message: $BABYLON_INIT_MSG"

# Instantiate only the Babylon contract
echo "Instantiating Babylon contract with Code ID $BABYLON_CODE_ID..."
$BINARY --home $CHAINDIR/$CHAINID tx wasm instantiate $BABYLON_CODE_ID "$BABYLON_INIT_MSG" --admin $ADMIN --label "babylon" $KEYRING --from user --chain-id $CHAINID --gas 20000000 --gas-prices 0.01$BASEDENOM --node http://localhost:$RPCPORT -y
sleep 3
# Get the Babylon contract address by querying contracts by code ID
BABYLON_ADDR=$($BINARY --home $CHAINDIR/$CHAINID query wasm list-contract-by-code $BABYLON_CODE_ID --output json | jq -r '.contracts[-1]')
echo "Babylon Address: $BABYLON_ADDR"

# Query the Babylon contract's Config {} to get all contract addresses
CONFIG_QUERY='{"config":{}}'
CONFIG_RES=$($BINARY --home $CHAINDIR/$CHAINID query wasm contract-state smart $BABYLON_ADDR "$CONFIG_QUERY" --node http://localhost:$RPCPORT --output json)
BTC_LC_ADDR=$(echo $CONFIG_RES | jq -r '.data.btc_light_client')
BTC_STAKING_ADDR=$(echo $CONFIG_RES | jq -r '.data.btc_staking')
BTC_FINALITY_ADDR=$(echo $CONFIG_RES | jq -r '.data.btc_finality')

# Get the governance module account address (this is the authority)
GOV_AUTHORITY=$($BINARY --home $CHAINDIR/$CHAINID query auth module-account gov --output json | jq -r '.account.value.address')
echo "Governance authority: $GOV_AUTHORITY"

# Create proposal JSON file with correct message type
echo "Creating governance proposal JSON..."
PROPOSAL_FILE="$CHAINDIR/$CHAINID/bsn_contracts_proposal.json"
DEPOSIT_AMOUNT="1000000$DENOM"
jq -n \
  --arg gov_authority "$GOV_AUTHORITY" \
  --arg babylon_addr "$BABYLON_ADDR" \
  --arg btc_lc_addr "$BTC_LC_ADDR" \
  --arg btc_staking_addr "$BTC_STAKING_ADDR" \
  --arg btc_finality_addr "$BTC_FINALITY_ADDR" \
  --arg deposit_amount "$DEPOSIT_AMOUNT" \
  '{
    "messages": [
      {
        "@type": "/babylonlabs.babylon.v1beta1.MsgSetBSNContracts",
        "authority": $gov_authority,
        "contracts": {
          "babylon_contract": $babylon_addr,
          "btc_light_client_contract": $btc_lc_addr,
          "btc_staking_contract": $btc_staking_addr,
          "btc_finality_contract": $btc_finality_addr
        }
      }
    ],
    "metadata": "Set BSN Contracts",
    "title": "Set BSN Contracts",
    "summary": "Set contract addresses for Babylon system",
    "deposit": $deposit_amount
  }' > "$PROPOSAL_FILE"

echo "Created proposal file: $PROPOSAL_FILE"
echo -n "Proposal content: "
cat "$PROPOSAL_FILE" | jq -r '.'

# Submit governance proposal to set BSN contracts
echo "Submitting governance proposal to set BSN contracts..."
PROPOSAL_RESP=$($BINARY --home $CHAINDIR/$CHAINID tx gov submit-proposal "$PROPOSAL_FILE" $KEYRING --from user --chain-id $CHAINID --gas 2000000 --gas-prices 0.01$BASEDENOM --node http://localhost:$RPCPORT -y --output json)
echo "Proposal response: $(echo $PROPOSAL_RESP | jq -r '.')"

# Clean up the proposal file
rm -f "$PROPOSAL_FILE"

# Extract proposal ID
PROPOSAL_TX_HASH=$(echo "$PROPOSAL_RESP" | jq -r '.txhash')
echo "Proposal transaction hash: $PROPOSAL_TX_HASH"
sleep 3
PROPOSAL_TX_RESULT=$($BINARY --home $CHAINDIR/$CHAINID query tx $PROPOSAL_TX_HASH --node http://localhost:$RPCPORT --output json)
echo "Proposal transaction result: $(echo $PROPOSAL_TX_RESULT | jq -r '.')"

PROPOSAL_ID=$(echo "$PROPOSAL_TX_RESULT" | jq -r '.events[] | select(.type=="submit_proposal") | .attributes[] | select(.key=="proposal_id") | .value')
echo "Extracted proposal ID: '$PROPOSAL_ID'"

# Verify we got a valid proposal ID
if [ -z "$PROPOSAL_ID" ] || [ "$PROPOSAL_ID" = "null" ]; then
    echo "Error: Failed to get proposal ID. Checking available proposals..."
    $BINARY --home $CHAINDIR/$CHAINID query gov proposals --output json --node http://localhost:$RPCPORT
    exit 1
fi

# Vote on the proposal
echo "Voting on proposal $PROPOSAL_ID..."
$BINARY --home $CHAINDIR/$CHAINID tx gov vote "$PROPOSAL_ID" yes $KEYRING --from validator --chain-id $CHAINID --gas 200000 --gas-prices 0.01$BASEDENOM --node http://localhost:$RPCPORT -y
sleep 3

# Wait for proposal to pass
echo "Waiting for proposal to pass..."
while true; do
    PROPOSAL_STATUS=$($BINARY --home $CHAINDIR/$CHAINID query gov proposal $PROPOSAL_ID --node http://localhost:$RPCPORT --output json | jq -r '.proposal.status')
    echo "  → Current proposal status: $PROPOSAL_STATUS"

    case "$PROPOSAL_STATUS" in
        "PROPOSAL_STATUS_PASSED")
            echo "  ✅ Proposal #$PROPOSAL_ID has passed!"
            break
            ;;
        "PROPOSAL_STATUS_REJECTED")
            echo "  ❌ Proposal #$PROPOSAL_ID was rejected!"
            exit 1
            ;;
        "PROPOSAL_STATUS_FAILED")
            echo "  ❌ Proposal #$PROPOSAL_ID failed!"
            exit 1
            ;;
        *)
            echo "  → Proposal status: $PROPOSAL_STATUS, waiting..."
            sleep 5
            ;;
    esac
done

# Verify the contracts are set
echo "Verifying BSN contracts are set..."
$BINARY query babylon bsn-contracts --node http://localhost:$RPCPORT --output json | jq -r '.'

echo "BSN contracts setup completed."
