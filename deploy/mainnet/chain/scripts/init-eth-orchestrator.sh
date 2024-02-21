#!/bin/bash
set -eu

echo "Initialising eth orchestrator"

# The prefix for cosmos addresses
ONOMY_ADDRESS_PREFIX="onomy"
# The address to run onomy node
ONOMY_HOST="0.0.0.0"
# The port of the onomy gRPC
ONOMY_GRPC_PORT="9191"
# The path to the gbt settings
ONOMY_GBT_CONFIG_HOME=$HOME/.gbt
# Onomy chain denom
ONOMY_STAKE_DENOM="anom"

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

ONOMY_VALIDATOR_NAME=validator
ONOMY_VALIDATOR_MNEMONIC=$(pass keyring-onomy/$ONOMY_VALIDATOR_NAME-mnemonic)
if [[ -z "${ONOMY_VALIDATOR_MNEMONIC}" ]]; then
  echo "Fail: check if key exists in pass: keyring-onomy/$ONOMY_ETH_ORCHESTRATOR_NAME-mnemonic"
  exit 1
fi

mkdir -p $ONOMY_GBT_CONFIG_HOME
cat assets/bridge/eth-orchestrator-config.toml > $ONOMY_GBT_CONFIG_HOME/config.toml

echo "The orchestrator config is initialised with the config:"
echo -e "\n$(cat $ONOMY_GBT_CONFIG_HOME/config.toml)\n"

echo "Registering orchestrator address in onomy chain"

gbt -a $ONOMY_ADDRESS_PREFIX keys register-orchestrator-address \
      --cosmos-phrase="$ONOMY_ETH_ORCHESTRATOR_MNEMONIC" \
      --validator-phrase="$ONOMY_VALIDATOR_MNEMONIC" \
      --ethereum-key="$ONOMY_ETH_ORCHESTRATOR_ETH_PRIVATE_KEY" \
      --cosmos-grpc="http://$ONOMY_HOST:$ONOMY_GRPC_PORT/" \
      --fees="1$ONOMY_STAKE_DENOM"
