#!/bin/bash
#Setup Validator
set -eu

echo "building environment"
# Initial dir
ONOMY_HOME=$HOME/.onomy
# Name of the network to bootstrap
CHAINID=$(jq .chain_id $ONOMY_HOME/node_info.json | sed 's#\"##g')
# Name of the gravity artifact
GRAVITY=onomyd
# The name of the gravity node
GRAVITY_NODE_NAME=$(jq .node_name $ONOMY_HOME/node_info.json | sed 's#\"##g')
# The address to run gravity node
GRAVITY_HOST="0.0.0.0"
# Home folder for gravity config
GRAVITY_HOME="$ONOMY_HOME/$CHAINID/$GRAVITY_NODE_NAME"
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
# Gravity chain demons
STAKE_DENOM="nom"
NORMAL_DENOM="footoken"

# -----------------Adding Validator---------------------
echo "Adding validator key"
read -p "Enter a name for your validator: " GRAVITY_VALIDATOR_NAME
echo $GRAVITY_HOME
$GRAVITY $GRAVITY_HOME_FLAG keys add $GRAVITY_VALIDATOR_NAME $GRAVITY_KEYRING_FLAG --output json | jq . >> $GRAVITY_HOME/validator_key.json
jq .mnemonic $GRAVITY_HOME/validator_key.json | sed 's#\"##g' >> $HOME/validator-phrases


# Save validator-info
# Switch sed command in the case of linux
fsed() {
  if [ `uname` = 'Linux' ]; then
    sed -i "$@"
  else
    sed -i '' "$@"
  fi
}

fsed 's#"validator_name": ""#"validator_name": "'$GRAVITY_VALIDATOR_NAME'"#g'  $ONOMY_HOME/node_info.json

# -------------------Get Faucet URL----------------
read -p "Please enter Faucet URL: " -i "http://testnet1.onomy.io:8000/" -e FAUCET_TOKEN_BASE_URL

# ------------------Get Tokens from Faucet------------------
ONOMY_VALIDATOR_ADDRESS=$(jq -r .address $GRAVITY_HOME/validator_key.json)
curl -X POST "$FAUCET_TOKEN_BASE_URL" -H  "accept: application/json" -H  "Content-Type: application/json" -d "{  \"address\": \"$ONOMY_VALIDATOR_ADDRESS\",  \"coins\": [    \"200000000nom\"  ]}"

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
