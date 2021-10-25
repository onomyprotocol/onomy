#!/bin/bash
set -eu

echo "Initializing master node"
# Initial dir
ONOMY_HOME=$HOME/.onomy
# Name of the network to bootstrap
CHAINID="onomy-testnet1"
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
# The name of the onomy orchestrator/validator
ONOMY_ORCHESTRATOR_NAME=orch
# Onomy chain demons
STAKE_DENOM="nom"
#NORMAL_DENOM="samoleans"
NORMAL_DENOM="footoken"
# The port of the onomy gRPC
ONOMY_GRPC_PORT="9090"

mkdir -p $ONOMY_HOME
# The host of ethereum node
ETH_HOST="0.0.0.0"
echo '{
        "validator_name": "",
        "chain_id": "",
        "orchestrator_name": ""
}' > $ONOMY_HOME/val_info.json

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
ONOMY_GENESIS_COINS="1000000000000$STAKE_DENOM,1000000000000$NORMAL_DENOM"

# Initialize the home directory and add some keys
echo "Init test chain"
$ONOMY $ONOMY_CHAINID_FLAG init $ONOMY_NODE_NAME

echo "Set stake/mint demon to $STAKE_DENOM"
fsed "s#\"stake\"#\"$STAKE_DENOM\"#g" $ONOMY_HOME_CONFIG/genesis.json

# add in denom metadata for both native tokens
jq '.app_state.bank.denom_metadata += [{"base": "footoken", display: "mfootoken", "description": "A non-staking test token", "denom_units": [{"denom": "footoken", "exponent": 0}, {"denom": "mfootoken", "exponent": 6}]}, {"base": "nom", display: "mnom", "description": "A staking test token", "denom_units": [{"denom": "nom", "exponent": 0}, {"denom": "mnom", "exponent": 6}]}]' $ONOMY_HOME_CONFIG/genesis.json > $ONOMY_HOME/metadata-genesis.json

# a 60 second voting period to allow us to pass governance proposals in the tests
jq '.app_state.gov.voting_params.voting_period = "60s"' $ONOMY_HOME/metadata-genesis.json > $ONOMY_HOME/edited-genesis.json
mv $ONOMY_HOME/edited-genesis.json $ONOMY_HOME/genesis.json
mv $ONOMY_HOME/genesis.json $ONOMY_HOME_CONFIG/genesis.json

echo "Add validator key"
$ONOMY keys add $ONOMY_VALIDATOR_NAME $ONOMY_KEYRING_FLAG --output json | jq . >> $ONOMY_HOME/validator_key.json
jq -r .mnemonic $ONOMY_HOME/validator_key.json > $ONOMY_HOME/validator-phrases

echo "Generating orchestrator keys"
$ONOMY keys add --output=json $ONOMY_ORCHESTRATOR_NAME $ONOMY_KEYRING_FLAG | jq . >> $ONOMY_HOME/orchestrator_key.json
jq -r .mnemonic $ONOMY_HOME/orchestrator_key.json > $ONOMY_HOME/orchestrator-phrases

echo "Adding validator addresses to genesis files"
$ONOMY add-genesis-account "$($ONOMY keys show $ONOMY_VALIDATOR_NAME -a $ONOMY_KEYRING_FLAG)" $ONOMY_GENESIS_COINS
echo "Adding orchestrator addresses to genesis files"
$ONOMY add-genesis-account "$($ONOMY keys show $ONOMY_ORCHESTRATOR_NAME -a $ONOMY_KEYRING_FLAG)" $ONOMY_GENESIS_COINS

echo "Adding faucet account addresses to genesis files"
$ONOMY keys add --output=json faucet_account1 $ONOMY_KEYRING_FLAG | jq . >> $ONOMY_HOME/faucet_account1.json
$ONOMY keys add --output=json faucet_account2 $ONOMY_KEYRING_FLAG | jq . >> $ONOMY_HOME/faucet_account2.json
$ONOMY keys add --output=json faucet_account3 $ONOMY_KEYRING_FLAG | jq . >> $ONOMY_HOME/faucet_account3.json
$ONOMY keys add --output=json faucet_account4 $ONOMY_KEYRING_FLAG | jq . >> $ONOMY_HOME/faucet_account4.json
$ONOMY keys add --output=json faucet_account5 $ONOMY_KEYRING_FLAG | jq . >> $ONOMY_HOME/faucet_account5.json
$ONOMY add-genesis-account "$($ONOMY keys show faucet_account1 -a $ONOMY_KEYRING_FLAG)" $ONOMY_GENESIS_COINS
$ONOMY add-genesis-account "$($ONOMY keys show faucet_account2 -a $ONOMY_KEYRING_FLAG)" $ONOMY_GENESIS_COINS
$ONOMY add-genesis-account "$($ONOMY keys show faucet_account3 -a $ONOMY_KEYRING_FLAG)" $ONOMY_GENESIS_COINS
$ONOMY add-genesis-account "$($ONOMY keys show faucet_account4 -a $ONOMY_KEYRING_FLAG)" $ONOMY_GENESIS_COINS
$ONOMY add-genesis-account "$($ONOMY keys show faucet_account5 -a $ONOMY_KEYRING_FLAG)" $ONOMY_GENESIS_COINS

echo "Generating ethereum keys"
$ONOMY eth_keys add --output=json | jq . >> $ONOMY_HOME/eth_key.json
echo "private: $(jq -r .private_key $ONOMY_HOME/eth_key.json)" > $ONOMY_HOME/validator-eth-keys
echo "public: $(jq -r .public_key $ONOMY_HOME/eth_key.json)" >> $ONOMY_HOME/validator-eth-keys
echo "address: $(jq -r .address $ONOMY_HOME/eth_key.json)" >> $ONOMY_HOME/validator-eth-keys

echo "Creating gentxs"
$ONOMY gentx --ip $ONOMY_HOST $ONOMY_VALIDATOR_NAME 1000000000000$STAKE_DENOM "$(jq -r .address $ONOMY_HOME/eth_key.json)" "$(jq -r .address $ONOMY_HOME/orchestrator_key.json)" $ONOMY_KEYRING_FLAG $ONOMY_CHAINID_FLAG

echo "Collecting gentxs in $ONOMY_NODE_NAME"
$ONOMY collect-gentxs

echo "Exposing ports and APIs of the $ONOMY_NODE_NAME"

# Change ports
fsed "s#\"tcp://127.0.0.1:26656\"#\"tcp://$ONOMY_HOST:26656\"#g" $ONOMY_NODE_CONFIG
fsed "s#\"tcp://127.0.0.1:26657\"#\"tcp://$ONOMY_HOST:26657\"#g" $ONOMY_NODE_CONFIG
fsed 's#addr_book_strict = true#addr_book_strict = false#g' $ONOMY_NODE_CONFIG
fsed 's#external_address = ""#external_address = "tcp://'$ONOMY_HOST:26656'"#g' $ONOMY_NODE_CONFIG
fsed 's#enable = false#enable = true#g' $ONOMY_APP_CONFIG
fsed 's#swagger = false#swagger = true#g' $ONOMY_APP_CONFIG

# Save validator-info
fsed 's#"validator_name": ""#"validator_name": "'$ONOMY_VALIDATOR_NAME'"#g'  $ONOMY_HOME/val_info.json
fsed 's#"chain_id": ""#"chain_id": "'$CHAINID'"#g'  $ONOMY_HOME/val_info.json