#!/bin/bash
set -eu

echo "Initializing full node"

# Onomy home dir
ONOMY_HOME=$HOME/.onomy

# Name of the network to bootstrap
CHAINID="onomy-testnet-1"
# The address to run onomy node
ONOMY_HOST="0.0.0.0"
# The port of the onomy gRPC
ONOMY_GRPC_PORT="9191"
# Config directories for onomy node
ONOMY_HOME_CONFIG="$ONOMY_HOME/config"
# Config file for onomy node
ONOMY_NODE_CONFIG="$ONOMY_HOME_CONFIG/config.toml"
# App config file for onomy node
ONOMY_APP_CONFIG="$ONOMY_HOME_CONFIG/app.toml"
# Chain ID flag
ONOMY_CHAINID_FLAG="--chain-id $CHAINID"
# Seeds IPs
ONOMY_SEEDS_DEFAULT_IPS="64.9.136.119,64.9.136.120,54.88.212.224"
# Statysync servers default IPs
ONOMY_STATESYNC_SERVERS_DEFAULT_IPS="3.219.52.168,44.206.144.197"

read -r -p "Enter a name for your node [onomy]:" ONOMY_NODE_NAME
ONOMY_NODE_NAME=${ONOMY_NODE_NAME:-onomy}

read -r -p "Enter seeds ips [$ONOMY_SEEDS_DEFAULT_IPS]:" ONOMY_SEEDS_IPS
ONOMY_SEEDS_IPS=${ONOMY_SEEDS_IPS:-$ONOMY_SEEDS_DEFAULT_IPS}

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
cp -r ../genesis/genesis-testnet-1.json $ONOMY_HOME_CONFIG/genesis.json

echo "Updating node config"

read -r -p "Do you want to setup state-sync?(y/n)[N]: " statesync
statesync=${statesync:-n}

statesync_nodes=
blk_height=
blk_hash=
if [[ $statesync = 'y' ]] || [[ $statesync = 'Y' ]]; then
    read -r -p "Enter IPs of statesync nodes (at least 2) [$ONOMY_STATESYNC_SERVERS_DEFAULT_IPS]:" statesync_ips
    statesync_ips=${statesync_ips:-$ONOMY_STATESYNC_SERVERS_DEFAULT_IPS}
    for statesync_ip in ${statesync_ips//,/ } ; do
      blk_details=$(curl -s http://$statesync_ip:26657/block | jq -r '.result.block.header.height + "\n" + .result.block_id.hash')
      blk_height=$(echo $blk_details | cut -d$' ' -f1)
      blk_hash=$(echo $blk_details | cut -d$' ' -f2)
      statesync_nodes="$statesync_nodes$statesync_ip:26657,"
    done

    # Change statesync settings
    crudini --set $ONOMY_NODE_CONFIG statesync enable true
    crudini --set $ONOMY_NODE_CONFIG statesync rpc_servers "\"$statesync_nodes\""
    crudini --set $ONOMY_NODE_CONFIG statesync trust_height $blk_height
    crudini --set $ONOMY_NODE_CONFIG statesync trust_hash "\"$blk_hash\""
    echo "Setup for statesync is complete"
fi

# change config
crudini --set $ONOMY_NODE_CONFIG p2p addr_book_strict false
crudini --set $ONOMY_NODE_CONFIG p2p external_address "\"tcp://$ip:26656\""
crudini --set $ONOMY_NODE_CONFIG p2p seeds "\"$ONOMY_SEEDS\""
crudini --set $ONOMY_NODE_CONFIG rpc laddr "\"tcp://$ONOMY_HOST:26657\""

crudini --set $ONOMY_APP_CONFIG grpc enable true
crudini --set $ONOMY_APP_CONFIG grpc address "\"$ONOMY_HOST:$ONOMY_GRPC_PORT\""
crudini --set $ONOMY_APP_CONFIG api enable true
crudini --set $ONOMY_APP_CONFIG api swagger true

echo "The initialisation of $ONOMY_NODE_NAME is done"
