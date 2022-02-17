#!/bin/bash
set -eu

echo "Initializing master node"
# Initial dir
ONOMY_HOME=$HOME/.onomy
# Name of the network to bootstrap
CHAINID="onomy-tm-load-test-chain"
# Name of the onomy artifact
ONOMY=onomyd
# The name of the onomy node
ONOMY_NODE_NAME="onomy"
# The address to run onomy node
ONOMY_HOST="0.0.0.0"
# Config directories for onomy node
ONOMY_HOME_CONFIG="$ONOMY_HOME/config"
# Config file for onomy node
ONOMY_NODE_CONFIG="$ONOMY_HOME_CONFIG/config.toml"
# App config file for onomy node
ONOMY_APP_CONFIG="$ONOMY_HOME_CONFIG/app.toml"
# Keyring flag
ONOMY_KEYRING_FLAG="--keyring-backend test"
# Chain ID flag
ONOMY_CHAINID_FLAG="--chain-id $CHAINID"
# The name of the onomy validator
ONOMY_VALIDATOR_NAME=validator1
# Onomy chain demons
STAKE_DENOM="anom"
#NORMAL_DENOM="samoleans"
NORMAL_DENOM="footoken"

mkdir -p $ONOMY_HOME

# Switch sed command in the case of linux
fsed() {
  if [ `uname` = 'Linux' ]; then
    sed -i "$@"
  else
    sed -i '' "$@"
  fi
}
# ------------------ Init onomy ------------------

echo "Creating $ONOMY_NODE_NAME validator with chain-id=$CHAINID..."
echo "Initializing genesis files"
# Build genesis file incl account for passed address
ONOMY_GENESIS_COINS="1000000000000000000000000$STAKE_DENOM,1000000000000000000000000$NORMAL_DENOM"

# Initialize the home directory and add some keys
echo "Init test chain"
$ONOMY $ONOMY_CHAINID_FLAG init $ONOMY_NODE_NAME

echo "Set stake/mint demon to $STAKE_DENOM"
fsed "s#\"stake\"#\"$STAKE_DENOM\"#g" $ONOMY_HOME_CONFIG/genesis.json

echo "Add validator key"
$ONOMY keys add $ONOMY_VALIDATOR_NAME $ONOMY_KEYRING_FLAG --output json | jq . >> $ONOMY_HOME/validator_key.json
jq -r .mnemonic $ONOMY_HOME/validator_key.json > $ONOMY_HOME/validator-phrases

echo "Adding validator addresses to genesis files"
$ONOMY add-genesis-account "$($ONOMY keys show $ONOMY_VALIDATOR_NAME -a $ONOMY_KEYRING_FLAG)" $ONOMY_GENESIS_COINS

echo "Generating ethereum keys"
$ONOMY eth_keys add --output=json | jq . >> $ONOMY_HOME/eth_key.json

echo "Creating gentxs"
$ONOMY gravity gentx $ONOMY_VALIDATOR_NAME 1000000000000000000000000$STAKE_DENOM "$(jq -r .address $ONOMY_HOME/eth_key.json)" "$(jq -r .address $ONOMY_HOME/validator_key.json)" --moniker validator --ip $ONOMY_HOST $ONOMY_KEYRING_FLAG $ONOMY_CHAINID_FLAG

echo "Collecting gentxs in $ONOMY_NODE_NAME"
$ONOMY gravity collect-gentxs

echo "Exposing ports and APIs of the $ONOMY_NODE_NAME"

# Change ports
fsed "s#\"tcp://127.0.0.1:26656\"#\"tcp://$ONOMY_HOST:26656\"#g" $ONOMY_NODE_CONFIG
fsed "s#\"tcp://127.0.0.1:26657\"#\"tcp://$ONOMY_HOST:26657\"#g" $ONOMY_NODE_CONFIG
fsed 's#addr_book_strict = true#addr_book_strict = false#g' $ONOMY_NODE_CONFIG
fsed 's#external_address = ""#external_address = "tcp://'$ONOMY_HOST:26656'"#g' $ONOMY_NODE_CONFIG
fsed 's#log_level = \"info\"#log_level = \"error\"#g' $ONOMY_NODE_CONFIG
fsed 's#enable = false#enable = true#g' $ONOMY_APP_CONFIG
fsed 's#swagger = false#swagger = true#g' $ONOMY_APP_CONFIG