#!/bin/bash
set -eu

echo "Initializing full node"

# Onomy home dir
ONOMY_HOME=$HOME/.onomy

# The address to run onomy node
ONOMY_HOST="0.0.0.0"
# Name of the network to bootstrap
CHAINID="onomy-mainnet-1"
# Config directories for onomy node
ONOMY_HOME_CONFIG="$ONOMY_HOME/config"
# Config file for onomy node
ONOMY_NODE_CONFIG="$ONOMY_HOME_CONFIG/config.toml"
# Chain ID flag
ONOMY_CHAINID_FLAG="--chain-id $CHAINID"
# Seeds IPs
ONOMY_SEEDS_DEFAULT_IPS="35.184.202.80,35.247.95.86"



read -r -p "Enter a name for your node [onomy-sentry]:" ONOMY_NODE_NAME
ONOMY_NODE_NAME=${ONOMY_NODE_NAME:-onomy-sentry}

read -r -p "Enter seeds ips [$ONOMY_SEEDS_DEFAULT_IPS]:" ONOMY_SEEDS_IPS
ONOMY_SEEDS_IPS=${ONOMY_SEEDS_IPS:-$ONOMY_SEEDS_DEFAULT_IPS}

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

ONOMY_SEEDS=
for seedIP in ${ONOMY_SEEDS_IPS//,/ } ; do
  wget $seedIP:26657/status? -O $ONOMY_HOME/seed_status.json
  seedID=$(jq -r .result.node_info.id $ONOMY_HOME/seed_status.json)

  if [[ -z "${seedID}" ]]; then
    echo "Something went wrong, can't fetch $seedIP info: "
    cat $ONOMY_HOME/seed_status.json
    exit 1
  fi

  rm $ONOMY_HOME/seed_status.json

  ONOMY_SEEDS="$ONOMY_SEEDS$seedID@$seedIP:26656,"
done

# create home directory
mkdir -p $ONOMY_HOME

# ------------------ Init onomy ------------------

echo "Creating $ONOMY_NODE_NAME node with chain-id=$CHAINID..."

# Initialize the home directory and add some keys
echo "Initializing chain"
onomyd $ONOMY_CHAINID_FLAG init $ONOMY_NODE_NAME

#copy genesis file
cp -r ../genesis/genesis-mainnet-1.json $ONOMY_HOME_CONFIG/genesis.json

echo "Updating node config"

# change config

crudini --set $ONOMY_NODE_CONFIG p2p addr_book_strict false
crudini --set $ONOMY_NODE_CONFIG p2p external_address "\"tcp://$ip:26656\""
crudini --set $ONOMY_NODE_CONFIG p2p seeds "\"$ONOMY_SEEDS\""
crudini --set $ONOMY_NODE_CONFIG rpc laddr "\"tcp://$ONOMY_HOST:26657\""

# sentry specific config
# pex	true - by default
crudini --set $ONOMY_NODE_CONFIG p2p pex true
# persistent_peers	validator node, optionally other sentry nodes (nodeid@ip:port)
crudini --set $ONOMY_NODE_CONFIG p2p persistent_peers "\"$ONOMY_VALIDATOR_PEER\""
# private_peer_ids	validator node ID
crudini --set $ONOMY_NODE_CONFIG p2p private_peer_ids "\"$ONOMY_VALIDATOR_ID\""
# unconditional_peer_ids	validator node ID, optionally sentry node IDs
crudini --set $ONOMY_NODE_CONFIG p2p unconditional_peer_ids "\"$ONOMY_VALIDATOR_ID\""

echo "The initialisation of $ONOMY_NODE_NAME is done"
