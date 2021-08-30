#!/bin/bash
set -eu

echo "building environment"

# Initial dir
CURRENT_WORKING_DIR=$(pwd)
# Name of the network to bootstrap
CHAINID="onomy"
# Name of the onomy artifact
ONOMY=onomyd
# The name of the onomy node
ONOMY_NODE_NAME="onomy"
# The address to run onomy node
ONOMY_HOST="0.0.0.0"
# Home folder for onomy config
ONOMY_HOME="$CURRENT_WORKING_DIR/$CHAINID/$ONOMY_NODE_NAME"
# Home flag for home folder
ONOMY_HOME_FLAG="--home $ONOMY_HOME"
# Config directories for onomy node
ONOMY_HOME_CONFIG="$ONOMY_HOME/config"
# Config file for onomy node
ONOMY_NODE_CONFIG="$ONOMY_HOME_CONFIG/config.toml"
# Genesis file for onomy node
ONOMY_NODE_GENESIS="$ONOMY_HOME_CONFIG/genesis.json"
# App config file for onomy node
ONOMY_APP_CONFIG="$ONOMY_HOME_CONFIG/app.toml"
# Keyring flag
ONOMY_KEYRING_FLAG="--keyring-backend test"
# Chain ID flag
ONOMY_CHAINID_FLAG="--chain-id $CHAINID"
# The name of the onomy validator
ONOMY_VALIDATOR_NAME=val
# The name of the onomy orchestrator/validator
ONOMY_ORCHESTRATOR_NAME=orch
# Onomy chain demons
STAKE_DENOM="nom"
# The ethereum address on validator/orchestrator
ETH_ORCHESTRATOR_VALIDATOR_ADDRESS=0x2d9480eBA3A001033a0B8c3Df26039FD3433D55d

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
ONOMY_GENESIS_COINS="1000000000000$STAKE_DENOM"

# Initialize the home directory and add some keys
echo "Init test chain"
$ONOMY $ONOMY_HOME_FLAG $ONOMY_CHAINID_FLAG init $ONOMY_NODE_NAME
echo "Set stake/mint demon to $STAKE_DENOM"
fsed "s#\"stake\"#\"$STAKE_DENOM\"#g" $ONOMY_NODE_GENESIS
echo "Add validator key"
$ONOMY $ONOMY_HOME_FLAG keys add $ONOMY_VALIDATOR_NAME $ONOMY_KEYRING_FLAG --output json | jq . >> $ONOMY_HOME/validator_key.json
echo "Adding validator addresses to genesis files"
$ONOMY $ONOMY_HOME_FLAG add-genesis-account "$($ONOMY $ONOMY_HOME_FLAG keys show $ONOMY_VALIDATOR_NAME -a $ONOMY_KEYRING_FLAG)" $ONOMY_GENESIS_COINS
echo "Generating orchestrator keys"
$ONOMY $ONOMY_HOME_FLAG keys add $ONOMY_ORCHESTRATOR_NAME $ONOMY_KEYRING_FLAG --output json | jq . >> $ONOMY_HOME/orchestrator_key.json
echo "Adding orchestrator addresses to genesis files"
$ONOMY $ONOMY_HOME_FLAG add-genesis-account "$($ONOMY $ONOMY_HOME_FLAG keys show $ONOMY_ORCHESTRATOR_NAME -a $ONOMY_KEYRING_FLAG)" $ONOMY_GENESIS_COINS

echo "Creating gentxs"
$ONOMY $ONOMY_HOME_FLAG gentx --ip $ONOMY_HOST $ONOMY_VALIDATOR_NAME 100000000000$STAKE_DENOM "$ETH_ORCHESTRATOR_VALIDATOR_ADDRESS" "$(jq -r .address $ONOMY_HOME/orchestrator_key.json)" $ONOMY_KEYRING_FLAG $ONOMY_CHAINID_FLAG

echo "Collecting gentxs in $ONOMY_NODE_NAME"
$ONOMY $ONOMY_HOME_FLAG collect-gentxs

echo "Exposing ports and APIs of the $ONOMY_NODE_NAME"

# Change ports
fsed "s#\"tcp://127.0.0.1:26656\"#\"tcp://$ONOMY_HOST:26656\"#g" $ONOMY_NODE_CONFIG
fsed "s#\"tcp://127.0.0.1:26657\"#\"tcp://$ONOMY_HOST:26657\"#g" $ONOMY_NODE_CONFIG
fsed 's#addr_book_strict = true#addr_book_strict = false#g' $ONOMY_NODE_CONFIG
fsed 's#external_address = ""#external_address = "tcp://'$ONOMY_HOST:26656'"#g' $ONOMY_NODE_CONFIG
fsed 's#enable = false#enable = true#g' $ONOMY_APP_CONFIG
fsed 's#swagger = false#swagger = true#g' $ONOMY_APP_CONFIG
