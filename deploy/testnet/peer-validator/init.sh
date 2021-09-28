#!/bin/bash
set -eu

echo "building environment"
# Initial dir
CURRENT_WORKING_DIR=$HOME

# Name of the network to bootstrap
#echo "Enter chain-id"
#read chainid
CHAINID="onomy"
# Name of the gravity artifact
GRAVITY=onomyd
# The name of the gravity node
GRAVITY_NODE_NAME="onomy"
# The address to run gravity node
GRAVITY_HOST="0.0.0.0"
# Home folder for gravity config
GRAVITY_HOME="$CURRENT_WORKING_DIR/$CHAINID/$GRAVITY_NODE_NAME"
# Home flag for home folder
GRAVITY_HOME_FLAG="--home $GRAVITY_HOME"
# Config directories for gravity node
GRAVITY_HOME_CONFIG="$GRAVITY_HOME/config"
# Config file for gravity node
GRAVITY_NODE_CONFIG="$GRAVITY_HOME_CONFIG/config.toml"
# App config file for gravity node
GRAVITY_APP_CONFIG="$GRAVITY_HOME_CONFIG/app.toml"
# Keyring flag
GRAVITY_KEYRING_FLAG="--keyring-backend test"
# Chain ID flag
GRAVITY_CHAINID_FLAG="--chain-id $CHAINID"
# The name of the gravity validator
GRAVITY_VALIDATOR_NAME=validator
# Gravity chain demons
STAKE_DENOM="nom"
#NORMAL_DENOM="samoleans"
NORMAL_DENOM="footoken"
echo "Please enter node-id of any validator that is running in chain to add seed"
read seedline
echo "Please enter ip of validator for which you have added node-id"
read ip
SEED="$seedline@$ip:26656"
# make a file to store validator information
echo '{
        "validator_name": "",
        "chain_id": "",
        "orchestrator_name": ""
}' > $HOME/val_info.json

# ------------------ Init gravity ------------------

echo "Creating $GRAVITY_NODE_NAME validator with chain-id=$CHAINID..."
echo "Initializing genesis files"

# Initialize the home directory and add some keys
echo "Init test chain"
$GRAVITY $GRAVITY_HOME_FLAG $GRAVITY_CHAINID_FLAG init $GRAVITY_NODE_NAME


echo "Add validator key"
$GRAVITY $GRAVITY_HOME_FLAG keys add $GRAVITY_VALIDATOR_NAME $GRAVITY_KEYRING_FLAG --output json | jq . >> $GRAVITY_HOME/validator_key.json
jq .mnemonic $GRAVITY_HOME/validator_key.json | sed 's#\"##g' >> $HOME/validator-phrases


#copy master genesis file 
rm $GRAVITY_HOME_CONFIG/genesis.json
wget http://$ip:26657/genesis? -O $HOME/raw.json
jq .result.genesis $HOME/raw.json >> $GRAVITY_HOME_CONFIG/genesis.json
rm -rf $HOME/raw.json

echo "Exposing ports and APIs of the $GRAVITY_NODE_NAME"
# Switch sed command in the case of linux
fsed() {
  if [ `uname` = 'Linux' ]; then
    sed -i "$@"
  else
    sed -i '' "$@"
  fi
}

# Change ports
fsed "s#\"tcp://127.0.0.1:26656\"#\"tcp://$GRAVITY_HOST:26656\"#g" $GRAVITY_NODE_CONFIG
fsed "s#\"tcp://127.0.0.1:26657\"#\"tcp://$GRAVITY_HOST:26657\"#g" $GRAVITY_NODE_CONFIG
fsed 's#addr_book_strict = true#addr_book_strict = false#g' $GRAVITY_NODE_CONFIG
fsed 's#external_address = ""#external_address = "tcp://'$GRAVITY_HOST:26656'"#g' $GRAVITY_NODE_CONFIG
fsed 's#seeds = ""#seeds = "'$SEED'"#g' $GRAVITY_NODE_CONFIG
fsed 's#enable = false#enable = true#g' $GRAVITY_APP_CONFIG
fsed 's#swagger = false#swagger = true#g' $GRAVITY_APP_CONFIG
# Save validator-info
fsed 's#"validator_name": ""#"validator_name": "'$GRAVITY_VALIDATOR_NAME'"#g'  $HOME/val_info.json
fsed 's#"chain_id": ""#"chain_id": "'$CHAINID'"#g'  $HOME/val_info.json


$GRAVITY $GRAVITY_HOME_FLAG start &



