#!/bin/bash
set -eu

echo "Initialising gbt"

ONOMY_HOME=$HOME/.onomy
# The prefix for cosmos addresses
ONOMY_ADDRESS_PREFIX="onomy"
# The address to run onomy node
ONOMY_HOST="0.0.0.0"
# The port of the onomy gRPC
ONOMY_GRPC_PORT="9191"
# The mnemonic of validator
ONOMY_VALIDATOR_MNEMONIC=$(jq -r .mnemonic $ONOMY_HOME/validator_key.json)
# The mnemonic of orchestrator
ONOMY_ORCHESTRATOR_MNEMONIC=$(jq -r .mnemonic $ONOMY_HOME/orchestrator_key.json)

# Onomy chain demons
STAKE_DENOM="anom"

if [[ -z "${ONOMY_VALIDATOR_MNEMONIC}" ]]; then
  echo "Fail: ONOMY_VALIDATOR_MNEMONIC is empty, check the file: $ONOMY_HOME/validator_key.json"
  exit 1
fi

if [[ -z "${ONOMY_ORCHESTRATOR_MNEMONIC}" ]]; then
  echo "Fail: ONOMY_ORCHESTRATOR_MNEMONIC is empty, check the file: $ONOMY_HOME/orchestrator_key.json"
  exit 1
fi

gbt init

echo "Registering orchestrator address"

gbt -a $ONOMY_ADDRESS_PREFIX keys register-orchestrator-address \
      --cosmos-phrase="$ONOMY_ORCHESTRATOR_MNEMONIC" \
      --validator-phrase="$ONOMY_VALIDATOR_MNEMONIC" \
      --ethereum-key="$ETH_ORCHESTRATOR_VALIDATOR_PRIVATE_KEY" \
      --cosmos-grpc="http://$ONOMY_HOST:$ONOMY_GRPC_PORT/" \
      --fees="1$STAKE_DENOM"

gbt -a $ONOMY_ADDRESS_PREFIX keys set-orchestrator-key --phrase="$ONOMY_ORCHESTRATOR_MNEMONIC"

echo "gbt is initialised"