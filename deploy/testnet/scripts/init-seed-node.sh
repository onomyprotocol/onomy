#!/bin/bash
set -eu

echo "Initializing seed node"

# Initial dir
ONOMY_HOME=$HOME/.onomy

CHAINID="onomy-testnet"
# The address to run onomy node
ONOMY_HOST="0.0.0.0"
# Config directories for onomy node
ONOMY_HOME_CONFIG="$ONOMY_HOME/config"
# Config file for onomy node
ONOMY_NODE_CONFIG="$ONOMY_HOME_CONFIG/config.toml"
# Chain ID flag
ONOMY_CHAINID_FLAG="--chain-id $CHAINID"

read -r -p "Enter a name for your node [onomy-seed]:" ONOMY_NODE_NAME
ONOMY_NODE_NAME=${ONOMY_NODE_NAME:-onomy-seed}

ONOMY_SEEDS_IPS=
read -r -p "Optionally enter seeds ips, ip1,ip2:" ONOMY_SEEDS_IPS

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

# make a file to store node information
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
echo "Init chain"
onomyd $ONOMY_CHAINID_FLAG init $ONOMY_NODE_NAME

# if variable not empty load genesis from the seed
if [[ -n "${ONOMY_SEEDS}" ]]; then
  #copy master genesis file from the seed
  seed_id=$(sed 's/.*@\(.*\):.*/\1/' <<< "$ONOMY_SEEDS")
  wget $seed_id:26657/genesis? -O $ONOMY_HOME/raw_genesis.json
  rm $ONOMY_HOME_CONFIG/genesis.json
  jq .result.genesis $ONOMY_HOME/raw_genesis.json >> $ONOMY_HOME_CONFIG/genesis.json
  rm $ONOMY_HOME/raw_genesis.json
fi

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