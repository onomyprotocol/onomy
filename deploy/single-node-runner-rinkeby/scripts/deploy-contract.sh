#!/bin/bash
set -eu

echo "running onomy with bridge"

# Name of the onomy artifact
ONOMY=onomyd
# The name of the onomy node
ONOMY_NODE_NAME="onomy"
# Home folder for onomy config
ONOMY_HOME="$HOME/.onomy"
# The address to run onomy node
ONOMY_HOST="0.0.0.0"
# The port of the onomy gRPC
ONOMY_GRPC_PORT="9090"
# The ETH key used for orchestrator signing of the transactions
ETH_ORCHESTRATOR_PRIVATE_KEY=c40f62e75a11789dbaf6ba82233ce8a52c20efb434281ae6977bb0b3a69bf709

# ------------------ Run onomy ------------------

echo "Starting $ONOMY_NODE_NAME"
$ONOMY start --pruning=nothing &

echo "Waiting $ONOMY_NODE_NAME to launch gRPC $ONOMY_GRPC_PORT..."

while ! timeout 1 bash -c "</dev/tcp/$ONOMY_HOST/$ONOMY_GRPC_PORT"; do
  sleep 1
done

echo "$ONOMY_NODE_NAME launched"

#-------------------- Deploy the contract --------------------

echo "Deploying Gravity contract, using: $ETH_RPC_ADDRESS"

contract-deployer \
--cosmos-node="http://$ONOMY_HOST:26657" \
--eth-node="$ETH_RPC_ADDRESS" \
--eth-privkey="$ETH_ORCHESTRATOR_PRIVATE_KEY" \
--contract=/root/home/solidity/Gravity.json \
--test-mode=false | grep "Gravity deployed at Address" | grep -Eow '0x[0-9a-fA-F]{40}' >> $ONOMY_HOME/eth_contract_address

CONTRACT_ADDRESS=$(cat $ONOMY_HOME/eth_contract_address)

echo "Contract address: $CONTRACT_ADDRESS"


