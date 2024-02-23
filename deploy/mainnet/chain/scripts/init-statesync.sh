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

crudini --set $ONOMY_NODE_CONFIG statesync enable true
crudini --set $ONOMY_NODE_CONFIG statesync rpc_servers "https://rpc-mainnet.onomy.io:443,http://35.224.118.71:26657"
crudini --set $ONOMY_NODE_CONFIG statesync trust_height $BLOCK_HEIGHT
crudini --set $ONOMY_NODE_CONFIG statesync trust_hash "\"$TRUST_HASH\""
crudini --set $ONOMY_NODE_CONFIG p2p seeds "211535f9b799bcc8d46023fa180f3359afd4c1d3@44.213.44.5:26656,00ce2f84f6b91639a7cedb2239e38ffddf9e36de@44.195.221.88:26656,cd9a47cebe8eef076a5795e1b8460a8e0b2384e5@3.210.0.126:26656,60194df601164a8b5852087d442038e392bf7470@180.131.222.74:26656,0dbe561f30862f386456734f12f431e534a3139c@34.133.228.142:26656,4737740b63d6ba9ebe93e8cc6c0e9197c426e9f4@195.189.96.106:52756"
crudini --set $ONOMY_NODE_CONFIG p2p persistent_peers "211535f9b799bcc8d46023fa180f3359afd4c1d3@44.213.44.5:26656,00ce2f84f6b91639a7cedb2239e38ffddf9e36de@44.195.221.88:26656,cd9a47cebe8eef076a5795e1b8460a8e0b2384e5@3.210.0.126:26656,60194df601164a8b5852087d442038e392bf7470@180.131.222.74:26656,0dbe561f30862f386456734f12f431e534a3139c@34.133.228.142:26656,4737740b63d6ba9ebe93e8cc6c0e9197c426e9f4@195.189.96.106:52756,00ce2f84f6b91639a7cedb2239e38ffddf9e36de@44.195.221.88:26656"

echo "Setup for statesync is complete"
