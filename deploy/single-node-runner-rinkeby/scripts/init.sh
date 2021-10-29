#!/bin/bash
set -eu

echo "building environment"

# Name of the network to bootstrap
CHAINID="onomy"
# Name of the onomy artifact
ONOMY=onomyd
# The name of the onomy node
ONOMY_NODE_NAME="onomy"
# The address to run onomy node
ONOMY_HOST="0.0.0.0"
# Home folder for onomy config
ONOMY_HOME="$HOME/.onomy"
# Home flag for home folder
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
ONOMY_VALIDATOR_NAME="val"
# The name of the onomy orchestrator/validator
ONOMY_ORCHESTRATOR_NAME="orch"
# Onomy chain demons
STAKE_DENOM="anom"
# The max allowed amount of validators
ONOMY_MAX_VALIDATORS=30
# The min gov deposite 10nom
ONOMY_GOV_MIN_DEPOSIT="10000000000000000000"
# The ethereum address on validator/orchestrator
ETH_ORCHESTRATOR_VALIDATOR_ADDRESS=0x2d9480eBA3A001033a0B8c3Df26039FD3433D55d
# Equal to 100000 noms
ONOMY_GENESIS_COINS="100000000000000000000000$STAKE_DENOM"
# Equal to 10000 noms
ONOMY_STAKE_COINS="10000000000000000000000$STAKE_DENOM"

# ------------------ Init onomy ------------------

echo "Creating $ONOMY_NODE_NAME validator with chain-id=$CHAINID..."
echo "Initializing genesis files"
echo "Onomy home: $ONOMY_HOME"

# Switch sed command in the case of linux
fsed() {
  if [ `uname` = 'Linux' ]; then
    sed -i "$@"
  else
    sed -i '' "$@"
  fi
}

# Initialize the home directory and add some keys
echo "Init test chain"
$ONOMY $ONOMY_CHAINID_FLAG init $ONOMY_NODE_NAME

# modify genesys
echo "Set stake/mint demon to $STAKE_DENOM"
fsed "s#\"stake\"#\"$STAKE_DENOM\"#g" $ONOMY_NODE_GENESIS
# add in denom metadata for both native tokens
jq '.app_state.bank.denom_metadata = [{ "description": "nom coin", "denom_units": [ { "denom": "anom", "exponent": 0 }, { "denom": "mnom", "exponent": 6 }, { "denom": "nom", "exponent": 18 } ], "base": "anom", "display": "nom" }]' $ONOMY_HOME_CONFIG/genesis.json | sponge $ONOMY_HOME_CONFIG/genesis.json
# add change max validators
jq ".app_state.staking.params.max_validators = $ONOMY_MAX_VALIDATORS" $ONOMY_HOME_CONFIG/genesis.json | sponge $ONOMY_HOME_CONFIG/genesis.json
# add change min gov deposit
jq ".app_state.gov.deposit_params.min_deposit = [{\"denom\": \"$STAKE_DENOM\", \"amount\": \"$ONOMY_GOV_MIN_DEPOSIT\"}]" $ONOMY_HOME_CONFIG/genesis.json | sponge $ONOMY_HOME_CONFIG/genesis.json

echo "Add validator key"
$ONOMY keys add $ONOMY_VALIDATOR_NAME $ONOMY_KEYRING_FLAG --output json | jq . >> $ONOMY_HOME/validator_key.json
echo "Adding validator addresses to genesis files"
$ONOMY add-genesis-account "$($ONOMY keys show $ONOMY_VALIDATOR_NAME -a $ONOMY_KEYRING_FLAG)" $ONOMY_GENESIS_COINS
echo "Generating orchestrator keys"
$ONOMY keys add $ONOMY_ORCHESTRATOR_NAME $ONOMY_KEYRING_FLAG --output json | jq . >> $ONOMY_HOME/orchestrator_key.json
echo "Adding orchestrator addresses to genesis files"
$ONOMY add-genesis-account "$($ONOMY keys show $ONOMY_ORCHESTRATOR_NAME -a $ONOMY_KEYRING_FLAG)" $ONOMY_GENESIS_COINS

echo "Creating gentxs"
$ONOMY gentx --ip $ONOMY_HOST $ONOMY_VALIDATOR_NAME $ONOMY_STAKE_COINS "$ETH_ORCHESTRATOR_VALIDATOR_ADDRESS" "$(jq -r .address $ONOMY_HOME/orchestrator_key.json)" $ONOMY_KEYRING_FLAG $ONOMY_CHAINID_FLAG

echo "Collecting gentxs in $ONOMY_NODE_NAME"
$ONOMY collect-gentxs

# Change node config
fsed "s#\"tcp://127.0.0.1:26656\"#\"tcp://$ONOMY_HOST:26656\"#g" "$ONOMY_NODE_CONFIG"
fsed "s#\"tcp://127.0.0.1:26657\"#\"tcp://$ONOMY_HOST:26657\"#g" $ONOMY_NODE_CONFIG
fsed 's#addr_book_strict = true#addr_book_strict = false#g' $ONOMY_NODE_CONFIG
fsed 's#external_address = ""#external_address = "tcp://'$ONOMY_HOST:26656'"#g' $ONOMY_NODE_CONFIG
fsed 's#log_level = \"info\"#log_level = \"error\"#g' $ONOMY_NODE_CONFIG

fsed 's#enable = false#enable = true#g' $ONOMY_APP_CONFIG
fsed 's#swagger = false#swagger = true#g' $ONOMY_APP_CONFIG
fsed 's#enabled-unsafe-cors = false#enabled-unsafe-cors = true#g' $ONOMY_APP_CONFIG