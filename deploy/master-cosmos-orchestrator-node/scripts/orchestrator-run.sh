#!/bin/bash
set -eu

echo "running gravity-bridge"

# Initial dir
CURRENT_WORKING_DIR=$(pwd)
# Name of the network to bootstrap
CHAINID="testchain"
# Name of the gravity artifact
GRAVITY=gravity
# The name of the gravity node
GRAVITY_NODE_NAME="gravity"
# The address to run gravity node
#GRAVITY_HOST="0.0.0.0"
GRAVITY_HOST=$1
# Home folder for gravity config
GRAVITY_HOME="$CURRENT_WORKING_DIR/$CHAINID/$GRAVITY_NODE_NAME"
# Home flag for home folder
GRAVITY_HOME_FLAG="--home $GRAVITY_HOME"
# Prefix of cosmos addresses
GRAVITY_ADDRESS_PREFIX=cosmos
# Gravity chain demons
STAKE_DENOM="stake"

ETH_MINER_PUBLIC_KEY="0xBf660843528035a5A4921534E156a27e64B231fE"
# The host of ethereum node
ETH_HOST="0.0.0.0"
ETH_HOST=$2

CONTRACT_ADDRESS=$(cat $GRAVITY_HOME/eth_contract_address)
echo "Contract address: $CONTRACT_ADDRESS"

echo "Gathering keys for orchestrator"
COSMOS_GRPC="http://$GRAVITY_HOST:9090/"
COSMOS_PHRASE=$(jq -r .mnemonic $GRAVITY_HOME/orchestrator_key.json)
ETH_RPC=http://$ETH_HOST:8545
ETH_PRIVATE_KEY=$(jq -r .private_key $GRAVITY_HOME/eth_key.json)
echo "Run orchestrator"

orchestrator --cosmos-phrase="$COSMOS_PHRASE" \
             --ethereum-key="$ETH_PRIVATE_KEY" \
             --cosmos-grpc="$COSMOS_GRPC" \
             --ethereum-rpc="$ETH_RPC" \
             --fees="$STAKE_DENOM" \
             --contract-address="$CONTRACT_ADDRESS"\
             --address-prefix="$GRAVITY_ADDRESS_PREFIX"