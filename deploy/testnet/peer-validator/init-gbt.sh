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
# The mnemonic of orchestrator
ONOMY_ORCHESTRATOR_MNEMONIC=$(jq -r .mnemonic $ONOMY_HOME/validator_key.json)
# Onomy chain demons
STAKE_DENOM="anom"

if [[ -z "${ONOMY_ORCHESTRATOR_MNEMONIC}" ]]; then
  echo "Fail: ONOMY_ORCHESTRATOR_MNEMONIC is empty, check the file: $ONOMY_HOME/validator_key.json"
  exit
fi

gbt init

gbt -a $ONOMY_ADDRESS_PREFIX keys register-orchestrator-address \
      --cosmos-phrase="$ONOMY_ORCHESTRATOR_MNEMONIC" \
      --validator-phrase="$ONOMY_ORCHESTRATOR_MNEMONIC" \
      --ethereum-key="$ETH_ORCHESTRATOR_VALIDATOR_PRIVATE_KEY" \
      --cosmos-grpc="http://$ONOMY_HOST:$ONOMY_GRPC_PORT/" \
      --fees="0$STAKE_DENOM"

echo "gbt is initialised"