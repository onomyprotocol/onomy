#!/bin/bash
set -e

echo "Starting eth orchestrator"

# The prefix for cosmos addresses
ONOMY_ADDRESS_PREFIX="onomy"
# The address to run onomy node
ONOMY_HOST="0.0.0.0"
# The port of the onomy gRPC
ONOMY_GRPC_PORT="9191"

# Onomy chain denom
ONOMY_STAKE_DENOM="anom"
# The json file with all bridges addresses
ONOMY_CONTRACT_ADDRESSES_PATH=$CHAIN_SCRIPTS"assets/bridge/addresses.json"

ONOMY_ETH_BRIDGE_CONTRACT_ADDRESS=$(jq -r .ethereum $ONOMY_CONTRACT_ADDRESSES_PATH)
if [[ -z "${ONOMY_ETH_BRIDGE_CONTRACT_ADDRESS}" ]]; then
  echo "Fail: ONOMY_ETH_BRIDGE_CONTRACT_ADDRESS is empty, check the $ONOMY_CONTRACT_ADDRESSES_PATH file"
  exit 1
fi

ONOMY_ETH_ORCHESTRATOR_ETH_PRIVATE_KEY=$(pass keyring-onomy/eth-orchestrator-eth-private-key)
if [[ -z "${ONOMY_ETH_ORCHESTRATOR_ETH_PRIVATE_KEY}" ]]; then
  echo "Fail: check if key exists in pass: keyring-onomy/eth-orchestrator-eth-private-key"
  exit 1
fi

ONOMY_ETH_ORCHESTRATOR_NAME=eth-orchestrator
ONOMY_ETH_ORCHESTRATOR_MNEMONIC=$(pass keyring-onomy/$ONOMY_ETH_ORCHESTRATOR_NAME-mnemonic)
if [[ -z "${ONOMY_ETH_ORCHESTRATOR_MNEMONIC}" ]]; then
  echo "Fail: check if key exists in pass: keyring-onomy/$ONOMY_ETH_ORCHESTRATOR_NAME-mnemonic"
  exit 1
fi

if [[ -z "${ETH_RPC_ADDRESS}" ]]; then
  echo "Fail: ETH_RPC_ADDRESS is not set"
  exit 1
fi

echo "ETH_RPC_ADDRESS: $ETH_RPC_ADDRESS, ONOMY_ETH_BRIDGE_CONTRACT_ADDRESS: $ONOMY_ETH_BRIDGE_CONTRACT_ADDRESS"

#-------------------- Run orchestrator --------------------

while ! timeout 1 bash -c "</dev/tcp/$ONOMY_HOST/$ONOMY_GRPC_PORT"; do
  sleep 1
done

echo "Starting orchestrator"

gbt -a $ONOMY_ADDRESS_PREFIX orchestrator \
             --cosmos-phrase="$ONOMY_ETH_ORCHESTRATOR_MNEMONIC" \
             --ethereum-key="$ONOMY_ETH_ORCHESTRATOR_ETH_PRIVATE_KEY" \
             --cosmos-grpc="http://$ONOMY_HOST:$ONOMY_GRPC_PORT/" \
             --ethereum-rpc="$ETH_RPC_ADDRESS" \
             --fees="1$ONOMY_STAKE_DENOM" \
             --gravity-contract-address="$ONOMY_ETH_BRIDGE_CONTRACT_ADDRESS"