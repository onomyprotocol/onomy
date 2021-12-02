#!/bin/bash
set -eu

echo "Initializing full node"

# Initial dir
ONOMY_HOME=$HOME/.onomy

# Name of the network to bootstrap
#echo "Enter chain-id"
#read chainid
CHAINID="onomy-testnet"
# Name of the onomy artifact
ONOMY=onomyd

read -r -p "Enter a name for your node [onomy-sentry]:" ONOMY_NODE_NAME
ONOMY_NODE_NAME=${ONOMY_NODE_NAME:-onomy-sentry}

# The address to run onomy node
ONOMY_HOST="0.0.0.0"
# Config directories for onomy node
ONOMY_HOME_CONFIG="$ONOMY_HOME/config"
# Config file for onomy node
ONOMY_NODE_CONFIG="$ONOMY_HOME_CONFIG/config.toml"
# Chain ID flag
ONOMY_CHAINID_FLAG="--chain-id $CHAINID"
# Seed node

# @TODO @Parth update it to multiple seed nodes
read -r -p "Enter node id of an existing seed [5e0f5b9d54d3e038623ddb77c0b91b559ff13495]:" ONOMY_SEED_ID
ONOMY_SEED_ID=${ONOMY_SEED_ID:-5e0f5b9d54d3e038623ddb77c0b91b559ff13495}
# @TODO @Parth update it to any of seed nodes
read -r -p "Enter Hostname/IP Address of the same node [testnet1.onomy.io]:" ONOMY_SEED_IP
ONOMY_SEED_IP=${ONOMY_SEED_IP:-"testnet1.onomy.io"}

ONOMY_SEED="$ONOMY_SEED_ID@$ONOMY_SEED_IP:26656"

ONOMY_VALIDATOR_ID=
while [[ $ONOMY_VALIDATOR_ID = "" ]]; do
   read -r -p "Enter node id of an existing validator:" ONOMY_VALIDATOR_ID
done

ONOMY_VALIDATOR_IP=
while [[ $ONOMY_VALIDATOR_IP = "" ]]; do
   read -r -p "Enter Hostname/IP Address of the validator node:" ONOMY_VALIDATOR_IP
done

ONOMY_VALIDATOR_PEER="$ONOMY_VALIDATOR_ID@$ONOMY_VALIDATOR_IP:26656"

# create home directory
mkdir -p $ONOMY_HOME

# make a file to store validator information
echo '{
  "node_name": "",
  "validator_name": "",
  "chain_id": "",
  "orchestrator_name": ""
}' > $ONOMY_HOME/node_info.json

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

#copy master genesis file
rm $ONOMY_HOME_CONFIG/genesis.json
wget $ONOMY_SEED_IP:26657/genesis? -O $ONOMY_HOME/raw_genesis.json
jq .result.genesis $ONOMY_HOME/raw_genesis.json >> $ONOMY_HOME_CONFIG/genesis.json
rm $ONOMY_HOME/raw_genesis.json

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
fsed 's#seeds = ""#seeds = "'$ONOMY_SEED'"#g' $ONOMY_NODE_CONFIG

# sentry specific config
# pex	true - by default
# persistent_peers	validator node, optionally other sentry nodes (nodeid@ip:port)
fsed 's#persistent_peers = ""#persistent_peers = "'$ONOMY_VALIDATOR_PEER'"#g' $ONOMY_NODE_CONFIG
# private_peer_ids	validator node ID
fsed 's#private_peer_ids = ""#private_peer_ids = "'$ONOMY_VALIDATOR_ID'"#g' $ONOMY_NODE_CONFIG
# unconditional_peer_ids	validator node ID, optionally sentry node IDs
fsed 's#unconditional_peer_ids = ""#unconditional_peer_ids = "'$ONOMY_VALIDATOR_ID'"#g' $ONOMY_NODE_CONFIG
# addr-book-strict	false
fsed 's#addr_book_strict = true#addr_book_strict = false#g' $ONOMY_NODE_CONFIG

fsed 's#"chain_id": ""#"chain_id": "'$CHAINID'"#g'  $ONOMY_HOME/node_info.json
fsed 's#"node_name": ""#"node_name": "'$ONOMY_NODE_NAME'"#g'  $ONOMY_HOME/node_info.json

echo "The initialisation of $ONOMY_NODE_NAME is done"