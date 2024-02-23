#!/bin/bash
set -eu

echo "Setting statesync config"

# Onomy home dir
ONOMY_HOME=$HOME/.onomy
# Config directories for onomy node
ONOMY_HOME_CONFIG="$ONOMY_HOME/config"
# Config file for onomy node
ONOMY_NODE_CONFIG="$ONOMY_HOME_CONFIG/config.toml"
# App config file for onomy node

LATEST_HEIGHT=$(curl -s "https://rpc-mainnet.onomy.io:443/block" | jq -r .result.block.header.height);
BLOCK_HEIGHT=$(($LATEST_HEIGHT - 2000)) 
TRUST_HASH=$(curl -s "https://rpc-mainnet.onomy.io:443/block?height=$BLOCK_HEIGHT" | jq -r .result.block_id.hash)
RPC_SERVERS="https://rpc-mainnet.onomy.io:443,http://35.224.118.71:26657"

crudini --set $ONOMY_NODE_CONFIG statesync enable true
crudini --set $ONOMY_NODE_CONFIG statesync rpc_servers "\"$RPC_SERVERS\""
crudini --set $ONOMY_NODE_CONFIG statesync trust_height $BLOCK_HEIGHT
crudini --set $ONOMY_NODE_CONFIG statesync trust_hash "\"$TRUST_HASH\""

echo "Setup for statesync is complete"
