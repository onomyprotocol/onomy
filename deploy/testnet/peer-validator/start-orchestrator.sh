#!/bin/bash
set -eu

echo "Starting onomy orchestrator"

# The prefix for cosmos addresses
ONOMY_ADDRESS_PREFIX="onomy"
# The address to run onomy node
ONOMY_HOST="0.0.0.0"
# The port of the onomy gRPC
ONOMY_GRPC_PORT="9191"
# Home folder for onomy config
ONOMY_HOME="$HOME/.onomy"
# The mnemonic of orchestrator
ONOMY_ORCHESTRATOR_MNEMONIC=$(jq -r .mnemonic $ONOMY_HOME/orchestrator_key.json)
# Onomy chain demons
STAKE_DENOM="anom"
# The path to orchestrator logs
mkdir -p $ONOMY_HOME/logs
ORCHESTRATOR_LOG_FILE=$ONOMY_HOME/logs/orchestrator.log

if [[ -z "${ETH_ORCHESTRATOR_VALIDATOR_PRIVATE_KEY}" ]]; then
  echo "Fail: ETH_ORCHESTRATOR_VALIDATOR_PRIVATE_KEY is not set"
  exit 1
fi

# TODO Parth, once the new contract is deployed remove this and set fixed address of ETH_GRAVITY_CONTRACT_ADDRESS
if [[ -z "${ETH_GRAVITY_CONTRACT_ADDRESS}" ]]; then
  echo "Fail: ETH_GRAVITY_CONTRACT_ADDRESS is empty"
  exit 1
fi

if [[ -z "${ONOMY_ORCHESTRATOR_MNEMONIC}" ]]; then
  echo "Fail: ONOMY_ORCHESTRATOR_MNEMONIC is empty, check the file: $ONOMY_HOME/orchestrator_key.json"
  exit 1
fi

if [[ -z "${ETH_RPC_ADDRESS}" ]]; then
  echo "Fail: ETH_RPC_ADDRESS is not set"
  exit 1
fi

echo "ETH_RPC_ADDRESS: $ETH_RPC_ADDRESS, ETH_GRAVITY_CONTRACT_ADDRESS: $ETH_GRAVITY_CONTRACT_ADDRESS"

#-------------------- Run orchestrator --------------------

echo "Starting orchestrator"

gbt -a $ONOMY_ADDRESS_PREFIX orchestrator \
             --cosmos-grpc="http://$ONOMY_HOST:$ONOMY_GRPC_PORT/" \
             --ethereum-rpc="$ETH_RPC_ADDRESS" \
             --fees="1$STAKE_DENOM" \
             --gravity-contract-address="$ETH_GRAVITY_CONTRACT_ADDRESS" &>> $ORCHESTRATOR_LOG_FILE &

echo "Orchestrator is started, check the logs file $ORCHESTRATOR_LOG_FILE"

