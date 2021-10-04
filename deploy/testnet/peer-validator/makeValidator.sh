#!/bin/bash
set -eu

echo "building environment"
# Initial dir
CURRENT_WORKING_DIR=$HOME
# Name of the network to bootstrap
CHAINID=$(jq .chain_id $HOME/val_info.json | sed 's#\"##g')
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
GRAVITY_VALIDATOR_NAME=$(jq .validator_name $HOME/val_info.json | sed 's#\"##g')
# Gravity chain demons
STAKE_DENOM="nom"
#NORMAL_DENOM="samoleans"
NORMAL_DENOM="footoken"
echo "Please enter faucet url get faucet token for example http://domain_name:8000/"
read url
FAUCET_TOKEN_BASE_URL="$url"


# ------------------ get faucet token------------------
ONOMY_VALIDATOR_ADDRESS=$(jq -r .address $GRAVITY_HOME/validator_key.json)
curl -X POST $FAUCET_TOKEN_BASE_URL -H  "accept: application/json" -H  "Content-Type: application/json" -d "{  \"address\": \"$ONOMY_VALIDATOR_ADDRESS\",  \"coins\": [    \"200000000nom\"  ]}"

#  wait 5 sec to sync balances in the validator account
sleep 5
# Store the public key of validator
PUB_KEY=$($GRAVITY $GRAVITY_HOME_FLAG tendermint show-validator)

# Do the create validator transaction
$GRAVITY $GRAVITY_HOME_FLAG tx staking create-validator \
--amount=100000000$STAKE_DENOM \
--pubkey="$PUB_KEY" \
--moniker=\"$GRAVITY_VALIDATOR_NAME\" \
--chain-id=$CHAINID \
--commission-rate="0.10" \
--commission-max-rate="0.20" \
--commission-max-change-rate="0.01" \
--min-self-delegation="10" \
--gas="auto" \
--gas-adjustment=1.5 \
--gas-prices="1$STAKE_DENOM" \
--from=$GRAVITY_VALIDATOR_NAME \
$GRAVITY_KEYRING_FLAG -y