#!/bin/bash
set -eu

echo "running onomy with bridge"

# Initial dir
CURRENT_WORKING_DIR=$(pwd)
# Name of the network to bootstrap
CHAINID="onomy"
# Name of the onomy artifact
ONOMY=onomyd
# The name of the onomy node
ONOMY_NODE_NAME="onomy"
# The prefix for cosmos addresses
ONOMY_ADDRESS_PREFIX="onomy"
# The address to run onomy node
ONOMY_HOST="0.0.0.0"
# The port of the onomy gRPC
ONOMY_GRPC_PORT="9090"
# Home folder for onomy config
ONOMY_HOME="$CURRENT_WORKING_DIR/$CHAINID/$ONOMY_NODE_NAME"
# Home flag for home folder
ONOMY_HOME_FLAG="--home $ONOMY_HOME"
# Onomy mnemonic used for fauset from (validator_key.json file)
ONOMY_VALIDATOR_MNEMONIC=$(jq -r .mnemonic $ONOMY_HOME/validator_key.json)

# The ETH key used for orchestrator signing of the transactions
ETH_ORCHESTRATOR_PRIVATE_KEY=c40f62e75a11789dbaf6ba82233ce8a52c20efb434281ae6977bb0b3a69bf709

# ------------------ Run onomy ------------------

echo "Starting $ONOMY_NODE_NAME"
$ONOMY $ONOMY_HOME_FLAG start --pruning=nothing &

echo "Waiting $ONOMY_NODE_NAME to launch gRPC $ONOMY_GRPC_PORT..."

while ! timeout 1 bash -c "</dev/tcp/$ONOMY_HOST/$ONOMY_GRPC_PORT"; do
  sleep 1
done

echo "$ONOMY_NODE_NAME launched"

#-------------------- Deploy the contract --------------------

echo "Deploying Gravity contract"
cd /root/home/solidity

contract-deployer \
--cosmos-node="http://$ONOMY_HOST:26657" \
--eth-node="$ETH_ADDRESS" \
--eth-privkey="$ETH_ORCHESTRATOR_PRIVATE_KEY" \
--contract=Gravity.json \
--test-mode=false | grep "Gravity deployed at Address" | grep -Eow '0x[0-9a-fA-F]{40}' >> $GRAVITY_HOME/eth_contract_address

CONTRACT_ADDRESS=$(cat $GRAVITY_HOME/eth_contract_address)

echo "Contract address: $CONTRACT_ADDRESS"


