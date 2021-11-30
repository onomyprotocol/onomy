#!/bin/bash
set -eu

echo "Initializing full node"

# Initial dir
ONOMY_HOME=$HOME/.onomy

# Name of the network to bootstrap
#echo "Enter chain-id"
#read chainid
CHAINID="onomy-testnet1"
# Name of the onomy artifact
ONOMY=onomyd

read -r -p "Enter a name for your node [onomy]:" ONOMY_NODE_NAME
ONOMY_NODE_NAME=${ONOMY_NODE_NAME:-onomy}

# The address to run onomy node
ONOMY_HOST="0.0.0.0"
# Config directories for onomy node
ONOMY_HOME_CONFIG="$ONOMY_HOME/config"
# Config file for onomy node
ONOMY_NODE_CONFIG="$ONOMY_HOME_CONFIG/config.toml"
# App config file for onomy node
ONOMY_APP_CONFIG="$ONOMY_HOME_CONFIG/app.toml"
# Chain ID flag
ONOMY_CHAINID_FLAG="--chain-id $CHAINID"
# Seed node

SHOULD_USE_STATE_SYNC="true"

STAKE_DENOM="anom"


read -r -p "Enter Hostname/IP Address to fetch genesis file [testnet1.onomy.io]:" ONOMY_SEED_IP
ONOMY_SEED_IP=${ONOMY_SEED_IP:-"testnet1.onomy.io"}

read -r -p "Enter seeds of nodes that have enabled snapshot [node_Id1@IP_node_Id1:26656,node_Id2@IP_node_Id2:26656]:" ONOMY_SEED
ONOMY_SEED=$ONOMY_SEED

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


# STATE SYNC
if [ "$SHOULD_USE_STATE_SYNC" = "true" ]; then

    read -r -p "Enter state sync RPC servers that have enabled snapshot [http://node1_IP:26657,http://node2_IP:26657]:" STATE_SYNC_RPC_SERVERS
    STATE_SYNC_RPC_SERVERS=${STATE_SYNC_RPC_SERVERS}
    IFS=', ' read -r -a STATE_SYNC_RPC_SERVERS_ARRAY <<< ${STATE_SYNC_RPC_SERVERS}

    TEMP_STATE_HEIGHT=$(curl -s ${STATE_SYNC_RPC_SERVERS_ARRAY[0]}/commit | jq -r ".result.signed_header.header.height")
    STATE_SYNC_HEIGHT=$(((($TEMP_STATE_HEIGHT/500)-1)*500))
    STATE_SYNC_HASH=$(curl -s ${STATE_SYNC_RPC_SERVERS_ARRAY[0]}/commit?height=${STATE_SYNC_HEIGHT} | jq '.result.signed_header.commit.block_id.hash')

    HASHES_MATCH=true
    for SERVER in "${STATE_SYNC_RPC_SERVERS_ARRAY[@]}"
    do
        TEMP_HASH=$(curl -s ${SERVER}/commit?height=${STATE_SYNC_HEIGHT} | jq '.result.signed_header.commit.block_id.hash')
        if [ "$STATE_SYNC_HASH" != "$TEMP_HASH" ]; then
            HASHES_MATCH=false
            break
        fi
    done

    if [ "$HASHES_MATCH" = "true" ]; then
        echo "Hashed matched"
        fsed -i "s/enable = false/enable = true/g" $ONOMY_NODE_CONFIG
        fsed -i "s~rpc_servers = \".*\"~rpc_servers = \"${STATE_SYNC_RPC_SERVERS}\"~g" $ONOMY_NODE_CONFIG
        fsed -i "s/trust_height = 0/trust_height = ${STATE_SYNC_HEIGHT}/g" $ONOMY_NODE_CONFIG
        fsed -i "s/trust_hash = \".*\"/trust_hash = ${STATE_SYNC_HASH}/g" $ONOMY_NODE_CONFIG
    else
        echo "Hashed from different peers don't match. State sync is OFF"
    fi
fi

# Change ports
fsed "s#\"tcp://127.0.0.1:26656\"#\"tcp://$ONOMY_HOST:26656\"#g" $ONOMY_NODE_CONFIG
fsed "s#\"tcp://127.0.0.1:26657\"#\"tcp://$ONOMY_HOST:26657\"#g" $ONOMY_NODE_CONFIG
fsed 's#addr_book_strict = true#addr_book_strict = false#g' $ONOMY_NODE_CONFIG
fsed 's#external_address = ""#external_address = "tcp://'$ONOMY_HOST:26656'"#g' $ONOMY_NODE_CONFIG
fsed 's#seeds = ""#seeds = "'$ONOMY_SEED'"#g' $ONOMY_NODE_CONFIG

fsed 's#minimum-gas-prices = ""#minimum-gas-prices = "0'$STAKE_DENOM'"#g' $ONOMY_APP_CONFIG
fsed 's#enable = false#enable = true#g' $ONOMY_APP_CONFIG
fsed 's#swagger = false#swagger = true#g' $ONOMY_APP_CONFIG
fsed 's#"chain_id": ""#"chain_id": "'$CHAINID'"#g'  $ONOMY_HOME/node_info.json
fsed 's#"node_name": ""#"node_name": "'$ONOMY_NODE_NAME'"#g'  $ONOMY_HOME/node_info.json

echo "The initialisation of $ONOMY_NODE_NAME is done"