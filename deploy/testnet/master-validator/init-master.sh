#!/bin/bash
set -eu

echo "Initializing master node"
# Initial dir
ONOMY_HOME=$HOME/.onomy
# Name of the network to bootstrap
CHAINID="onomy-testnet"
# Name of the onomy artifact
ONOMY=onomyd
# The name of the onomy node
ONOMY_NODE_NAME="onomy"
# The address to run onomy node
ONOMY_HOST="0.0.0.0"
# The port of the onomy gRPC
ONOMY_GRPC_PORT="9191"
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
# Genesys params
# Onomy chain demons
STAKE_DENOM="anom"
# Equal to 10000 noms
ONOMY_STAKE_COINS="10000000000000000000000$STAKE_DENOM"
# The coins for genesys account
ONOMY_GENESIS_COINS="100000000000000000000000000$STAKE_DENOM"
# Max validators on the chain
ONOMY_MAX_VALIDATORS=10
# The address of WNOM ERC20 token on ethereum.
ONOMY_WNOM_ERC20_ADDRESS="0xe7c0fd1f0A3f600C1799CD8d335D31efBE90592C"

if [[ -z "${ETH_ORCHESTRATOR_VALIDATOR_ADDRESS}" ]]; then
  echo "Fail: ETH_ORCHESTRATOR_VALIDATOR_ADDRESS is not set"
  exit
fi

mkdir -p $ONOMY_HOME
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

# ------------------ Get IP Address --------------
default_ip=$(hostname -I | awk '{print $1}')

read -r -p "Enter your ip address [$default_ip]: "
ip=${ip:-$default_ip}

# ------------------ Init onomy ------------------

echo "Creating $ONOMY_NODE_NAME validator with chain-id=$CHAINID..."
echo "Initializing genesis files"

# Initialize the home directory and add some keys
echo "Init test chain"
$ONOMY $ONOMY_CHAINID_FLAG init $ONOMY_NODE_NAME

echo "Set stake/mint demon to $STAKE_DENOM"
fsed "s#\"stake\"#\"$STAKE_DENOM\"#g" $ONOMY_HOME_CONFIG/genesis.json

# add in denom metadata for both native tokens
jq '.app_state.bank.denom_metadata += [{"description": "nom coin","denom_units": [{"denom": "anom","exponent": 0,"aliases": []},{"denom": "mnom","exponent": 6,"aliases": []},{"denom": "nom","exponent": 18,"aliases": []}],"base": "anom","display": "nom"}]' $ONOMY_HOME_CONFIG/genesis.json | sponge $ONOMY_HOME_CONFIG/genesis.json

# add change max validators
jq ".app_state.staking.params.max_validators = $ONOMY_MAX_VALIDATORS" $ONOMY_HOME_CONFIG/genesis.json | sponge $ONOMY_HOME_CONFIG/genesis.json

# add wnom -> anom swap settings
jq ".app_state.gravity.params.erc20_to_denom_permanent_swap = {\"erc20\": \"$ONOMY_WNOM_ERC20_ADDRESS\", \"denom\": \"$STAKE_DENOM\"}" $ONOMY_HOME_CONFIG/genesis.json | sponge $ONOMY_HOME_CONFIG/genesis.json

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

echo "Adding additional accounts to genesis files"

for i in {1..5};
do
	echo "Adding additional account $i"
	$ONOMY keys add --output=json account$i $ONOMY_KEYRING_FLAG | jq . >> $ONOMY_HOME/account$i.json
	$ONOMY add-genesis-account "$($ONOMY keys show account$i -a $ONOMY_KEYRING_FLAG)" $ONOMY_GENESIS_COINS
done

echo "Creating gentxs"
$ONOMY gentx --ip $ONOMY_HOST $ONOMY_VALIDATOR_NAME $ONOMY_STAKE_COINS $ETH_ORCHESTRATOR_VALIDATOR_ADDRESS "$(jq -r .address $ONOMY_HOME/orchestrator_key.json)" $ONOMY_KEYRING_FLAG $ONOMY_CHAINID_FLAG

echo "Collecting gentxs in $ONOMY_NODE_NAME"
$ONOMY collect-gentxs

echo "Exposing ports and APIs of the $ONOMY_NODE_NAME"

# Change ports
fsed "s#\"tcp://127.0.0.1:26656\"#\"tcp://$ONOMY_HOST:26656\"#g" $ONOMY_NODE_CONFIG
fsed "s#\"tcp://127.0.0.1:26657\"#\"tcp://$ONOMY_HOST:26657\"#g" $ONOMY_NODE_CONFIG
fsed 's#addr_book_strict = true#addr_book_strict = false#g' $ONOMY_NODE_CONFIG
fsed 's#external_address = ""#external_address = "tcp://'$ip:26656'"#g' $ONOMY_NODE_CONFIG
fsed "s#0.0.0.0:9090#$ONOMY_HOST:$ONOMY_GRPC_PORT#g" $ONOMY_APP_CONFIG
fsed 's#enable = false#enable = true#g' $ONOMY_APP_CONFIG
fsed 's#swagger = false#swagger = true#g' $ONOMY_APP_CONFIG

# Save validator-info
fsed 's#"validator_name": ""#"validator_name": "'$ONOMY_VALIDATOR_NAME'"#g'  $ONOMY_HOME/val_info.json
fsed 's#"chain_id": ""#"chain_id": "'$CHAINID'"#g'  $ONOMY_HOME/val_info.json