#!/bin/bash
set -eu

# The address to run onomy node
# The node is running on the host machine be the call to it we expect from the container.
# The hist to make the test pass on mac and linux
ONOMY_HOST="host.docker.internal"
if ! ping -c 1 $ONOMY_HOST &> /dev/null
then
  ONOMY_HOST="0.0.0.0"
fi

# The port of the onomy gRPC
ONOMY_GRPC_PORT="9090"

# The prefix for cosmos addresses
ONOMY_ADDRESS_PREFIX="onomy"

# read the mnemonic from param.
ONOMY_ORCHESTRATOR_MNEMONIC="$1"

# onomy stake denom
STAKE_DENOM=nom

# The URL of the running mock eth node.
ETH_ADDRESS="http://0.0.0.0:8545/"

# The ETH key used for orchestrator signing of the transactions.
ETH_ORCHESTRATOR_PRIVATE_KEY=c40f62e75a11789dbaf6ba82233ce8a52c20efb434281ae6977bb0b3a69bf709

# read already deployed address from the file
GRAVITY_CONTRACT_ADDRESS=$(cat gravity_contract_address)

echo "Starting orchestrator"

gbt --address-prefix="$ONOMY_ADDRESS_PREFIX" orchestrator \
             --cosmos-phrase="$ONOMY_ORCHESTRATOR_MNEMONIC" \
             --ethereum-key="$ETH_ORCHESTRATOR_PRIVATE_KEY" \
             --cosmos-grpc="http://$ONOMY_HOST:$ONOMY_GRPC_PORT/" \
             --ethereum-rpc="$ETH_ADDRESS" \
             --fees="1$STAKE_DENOM" \
             --gravity-contract-address="$GRAVITY_CONTRACT_ADDRESS"