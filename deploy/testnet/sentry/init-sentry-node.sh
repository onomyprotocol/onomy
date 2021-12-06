#!/bin/bash
set -eu

echo "Initializing full node"

# Initial dir
ONOMY_HOME=$HOME/.onomy

# Name of the network to bootstrap
CHAINID="onomy-testnet"
# Name of the onomy artifact
ONOMY=onomyd

# The address to run onomy node
ONOMY_HOST="0.0.0.0"
# Config directories for onomy node
ONOMY_HOME_CONFIG="$ONOMY_HOME/config"
# Config file for onomy node
ONOMY_NODE_CONFIG="$ONOMY_HOME_CONFIG/config.toml"
# Chain ID flag
ONOMY_CHAINID_FLAG="--chain-id $CHAINID"

read -r -p "Enter a name for your node [onomy-sentry]:" ONOMY_NODE_NAME
ONOMY_NODE_NAME=${ONOMY_NODE_NAME:-onomy-sentry}

ONOMY_SEEDS=
while [[ $ONOMY_SEEDS = "" ]]; do
   read -r -p "Enter seeds peers, id@ip:port,id2@ip2:port : " ONOMY_SEEDS
done

ONOMY_VALIDATOR_ID=
while [[ $ONOMY_VALIDATOR_ID = "" ]]; do
   read -r -p "Enter node id of an existing validator:" ONOMY_VALIDATOR_ID
done

ONOMY_VALIDATOR_IP=
while [[ $ONOMY_VALIDATOR_IP = "" ]]; do
   read -r -p "Enter Hostname/IP Address of the validator node:" ONOMY_VALIDATOR_IP
done

ONOMY_VALIDATOR_PEER="$ONOMY_VALIDATOR_ID@$ONOMY_VALIDATOR_IP:26656"

default_ip=$(hostname -I | awk '{print $1}')
read -r -p "Enter your ip address [$default_ip]:" ip
ip=${ip:-$default_ip}

# create home directory
mkdir -p $ONOMY_HOME

# make a file to store validator information
echo '{
  "node_name": "",
  "validator_name": "",
  "chain_id": "",
  "orchestrator_name": ""
}' > $ONOMY_HOME/node_info.json

# ------------------ Init onomy ------------------

echo "Creating $ONOMY_NODE_NAME node with chain-id=$CHAINID..."
echo "Initializing genesis files"

# Initialize the home directory and add some keys
echo "Init test chain"
$ONOMY $ONOMY_CHAINID_FLAG init $ONOMY_NODE_NAME

#copy master genesis file from the seed
seed_id=$(sed 's/.*@\(.*\):.*/\1/' <<< "$ONOMY_SEEDS")
wget $seed_id:26657/genesis? -O $ONOMY_HOME/raw_genesis.json
rm $ONOMY_HOME_CONFIG/genesis.json
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
fsed 's#seeds = ""#seeds = "'$ONOMY_SEEDS'"#g' $ONOMY_NODE_CONFIG

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