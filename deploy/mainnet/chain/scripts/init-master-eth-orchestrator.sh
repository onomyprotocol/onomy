#!/bin/bash
set -eu

echo "Initialising eth orchestrator"

# The path to the gbt settings
ONOMY_GBT_CONFIG_HOME=$HOME/.gbt

mkdir -p $ONOMY_GBT_CONFIG_HOME
cat assets/bridge/eth-orchestrator-config.toml > $ONOMY_GBT_CONFIG_HOME/config.toml

echo "The orchestrator config is initialised with the config:"
echo -e "\n$(cat $ONOMY_GBT_CONFIG_HOME/config.toml)\n"