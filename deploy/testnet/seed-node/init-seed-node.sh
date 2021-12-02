#!/bin/bash
set -eu

echo "Initializing seed node"

# Initial dir
ONOMY_HOME=$HOME/.onomy

# Name of the network to bootstrap
#echo "Enter chain-id"
#read chainid
CHAINID="onomy-testnet"
# Name of the onomy artifact
ONOMY=onomyd

read -r -p "Enter a name for your node [onomy-seed]:" ONOMY_NODE_NAME
ONOMY_NODE_NAME=${ONOMY_NODE_NAME:-onomy-seed}

# The address to run onomy node
ONOMY_HOST="0.0.0.0"
# Config directories for onomy node
ONOMY_HOME_CONFIG="$ONOMY_HOME/config"
# Config file for onomy node
ONOMY_NODE_CONFIG="$ONOMY_HOME_CONFIG/config.toml"
# Chain ID flag
ONOMY_CHAINID_FLAG="--chain-id $CHAINID"
# Seed node

# create home directory
mkdir -p $ONOMY_HOME

# make a file to store validator information
echo '{
  "node_name": "",
  "validator_name": "",
  "chain_id": "",
  "orchestrator_name": ""
}' > $ONOMY_HOME/node_info.json

ONOMY_SEEDS=
read -r -p "Optionally enter seed peer id@ip:port,id2@ip2:port :" ONOMY_SEEDS

# ------------------ Get IP Address --------------
default_ip=$(hostname -I | awk '{print $1}')

read -r -p "Enter your ip address [$default_ip]:" ip
ip=${ip:-$default_ip}

# ------------------ Init onomy ------------------

echo "Creating $ONOMY_NODE_NAME node with chain-id=$CHAINID..."
echo "Initializing genesis files"

# Initialize the home directory and add some keys
echo "Init test chain"
$ONOMY $ONOMY_CHAINID_FLAG init $ONOMY_NODE_NAME

echo "Exposing ports and APIs of the $ONOMY_NODE_NAME"
# Switch sed command in the case of linux
fsed() {
  if [ `uname` = 'Linux' ]; then
    sed -i "$@"
  else
    sed -i '' "$@"
  fi
}

# Change ports
fsed "s#\"tcp://127.0.0.1:26656\"#\"tcp://$ONOMY_HOST:26656\"#g" $ONOMY_NODE_CONFIG
fsed "s#\"tcp://127.0.0.1:26657\"#\"tcp://$ONOMY_HOST:26657\"#g" $ONOMY_NODE_CONFIG
fsed 's#external_address = ""#external_address = "tcp://'$ip:26656'"#g' $ONOMY_NODE_CONFIG
fsed 's#seeds = ""#seeds = "'$ONOMY_SEEDS'"#g' $ONOMY_NODE_CONFIG
fsed 's#addr_book_strict = true#addr_book_strict = false#g' $ONOMY_NODE_CONFIG

# seed config
fsed 's#seed_mode = false#seed_mode = true#g' $ONOMY_NODE_CONFIG

fsed 's#"chain_id": ""#"chain_id": "'$CHAINID'"#g'  $ONOMY_HOME/node_info.json
fsed 's#"node_name": ""#"node_name": "'$ONOMY_NODE_NAME'"#g'  $ONOMY_HOME/node_info.json

echo "The initialisation of $ONOMY_NODE_NAME is done"