#!/bin/bash
set -eu

echo "building environment"

# Initial dir
CURRENT_WORKING_DIR=$(pwd)
# Name of the network to bootstrap
CHAINID="testchain"
# Name of the gravity artifact
GRAVITY=gravity
# The name of the gravity node
GRAVITY_NODE_NAME="gravity"
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
GRAVITY_VALIDATOR_NAME=val
# The name of the gravity orchestrator/validator
GRAVITY_ORCHESTRATOR_NAME=orch
# Gravity chain demons
STAKE_DENOM="stake"
NORMAL_DENOM="samoleans"

# This key is the private key for the public key defined in ETHGenesis.json
# where the full node / miner sends its rewards. Therefore it's always going
# to have a lot of ETH to pay for things like contract deployments
ETH_MINER_PRIVATE_KEY="0xb1bab011e03a9862664706fc3bbaa1b16651528e5f0e7fbfcbfdd8be302a13e7"
ETH_MINER_PUBLIC_KEY="0xBf660843528035a5A4921534E156a27e64B231fE"
# The host of ethereum node
ETH_HOST="0.0.0.0"

# ------------------ Init gravity ------------------

echo "Creating $GRAVITY_NODE_NAME validator with chain-id=$CHAINID..."
echo "Initializing genesis files"
# Build genesis file incl account for passed address
GRAVITY_GENESIS_COINS="100000000000$STAKE_DENOM,100000000000$NORMAL_DENOM"

# Initialize the home directory and add some keys
echo "Init test chain"
$GRAVITY $GRAVITY_HOME_FLAG $GRAVITY_CHAINID_FLAG init $GRAVITY_NODE_NAME
echo "Add validator key"
$GRAVITY $GRAVITY_HOME_FLAG keys add $GRAVITY_VALIDATOR_NAME $GRAVITY_KEYRING_FLAG --output json | jq . >> $GRAVITY_HOME/validator_key.json
echo "Adding validator addresses to genesis files"
$GRAVITY $GRAVITY_HOME_FLAG add-genesis-account "$($GRAVITY $GRAVITY_HOME_FLAG keys show $GRAVITY_VALIDATOR_NAME -a $GRAVITY_KEYRING_FLAG)" $GRAVITY_GENESIS_COINS
echo "Generating orchestrator keys"
$GRAVITY $GRAVITY_HOME_FLAG keys add --dry-run=true --output=json $GRAVITY_ORCHESTRATOR_NAME | jq . >> $GRAVITY_HOME/orchestrator_key.json

echo "Adding orchestrator keys to genesis"
GRAVITY_ORCHESTRATOR_KEY="$(jq .address $GRAVITY_HOME/orchestrator_key.json)"

jq ".app_state.auth.accounts += [{\"@type\": \"/cosmos.auth.v1beta1.BaseAccount\",\"address\": $GRAVITY_ORCHESTRATOR_KEY,\"pub_key\": null,\"account_number\": \"0\",\"sequence\": \"0\"}]" $GRAVITY_HOME_CONFIG/genesis.json | sponge $GRAVITY_HOME_CONFIG/genesis.json
jq ".app_state.bank.balances += [{\"address\": $GRAVITY_ORCHESTRATOR_KEY,\"coins\": [{\"denom\": \"$NORMAL_DENOM\",\"amount\": \"100000000000\"},{\"denom\": \"$STAKE_DENOM\",\"amount\": \"100000000000\"}]}]" $GRAVITY_HOME_CONFIG/genesis.json | sponge $GRAVITY_HOME_CONFIG/genesis.json

echo "Generating ethereum keys"
$GRAVITY $GRAVITY_HOME_FLAG eth_keys add --output=json --dry-run=true | jq . >> $GRAVITY_HOME/eth_key.json

echo "Creating gentxs"
$GRAVITY $GRAVITY_HOME_FLAG gentx --ip $GRAVITY_HOST $GRAVITY_VALIDATOR_NAME 100000000000$STAKE_DENOM "$(jq -r .address $GRAVITY_HOME/eth_key.json)" "$(jq -r .address $GRAVITY_HOME/orchestrator_key.json)" $GRAVITY_KEYRING_FLAG $GRAVITY_CHAINID_FLAG

echo "Collecting gentxs in $GRAVITY_NODE_NAME"
$GRAVITY $GRAVITY_HOME_FLAG collect-gentxs

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
fsed 's#enable = false#enable = true#g' $GRAVITY_APP_CONFIG
fsed 's#swagger = false#swagger = true#g' $GRAVITY_APP_CONFIG

#echo "Adding initial ethereum value for miner"
#jq ".alloc |= . + {\"$ETH_MINER_PUBLIC_KEY\" : {\"balance\": \"0x1337000000000000000000\"}}" assets/ETHGenesis.json | sponge assets/ETHGenesis.json

echo "Adding initial ethereum value for gravity validator"
jq ".alloc |= . + {$(jq .address $GRAVITY_HOME/eth_key.json) : {\"balance\": \"0x1337000000000000000000\"}}" assets/ETHGenesis.json | sponge assets/ETHGenesis.json


