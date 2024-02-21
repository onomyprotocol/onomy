#!/bin/bash
set -eu

echo "Initializing master node"
# Onomy home dir
ONOMY_HOME=$HOME/.onomy
# Name of the network to bootstrap
CHAINID="onomy-mainnet-1"
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
ONOMY_KEYRING_FLAG="--keyring-backend pass"
# Chain ID flag
ONOMY_CHAINID_FLAG="--chain-id $CHAINID"
# The name of the onomy validator
ONOMY_VALIDATOR_NAME="validator"
# The name of the onomy orchestrator/validator
ONOMY_ETH_ORCHESTRATOR_NAME="eth-orchestrator"
# Genesis params
# Onomy chain denom
ONOMY_STAKE_DENOM="anom"
# Equal to 550k noms
ONOMY_VALIDATOR_GENESIS_COINS="550000000000000000000000$ONOMY_STAKE_DENOM"
# Equal to 10k noms
ONOMY_ETH_ORCHESTRATOR_GENESIS_COINS="10000000000000000000000$ONOMY_STAKE_DENOM"
# Equal to 43.6m noms
ONOMY_DAO_AMOUNT="43678177000000000000000000"
# Equal to 500k noms
ONOMY_STAKE_COINS="500000000000000000000000$ONOMY_STAKE_DENOM"
# Max validators on the chain
ONOMY_MAX_VALIDATORS=20
# The address of BNOM ERC20 token on ethereum.
ONOMY_BNOM_ERC20_ADDRESS="0x4fdF157C6860F33232608198419c74ED446Ca577"
# The min gov deposit 500 noms
ONOMY_GOV_MIN_DEPOSIT="500000000000000000000"
# 225k noms
ONOMY_MIN_GLOBAL_SELF_DELEGATION="225000000000000000000000"
# the validator min self delegation
ONOMY_MIN_SELF_DELEGATION_FLAG="--min-self-delegation=$ONOMY_MIN_GLOBAL_SELF_DELEGATION"
# path to the genesis accounts
ONOMY_GENESIS_ACCOUNTS=$(cat ../genesis/accounts.txt)

if [[ -z "${ONOMY_ETH_ORCHESTRATOR_VALIDATOR_ADDRESS}" ]]; then
  echo "Fail: ONOMY_ETH_ORCHESTRATOR_VALIDATOR_ADDRESS is not set"
  exit 1
fi

default_ip=$(hostname -I | awk '{print $1}')
read -r -p "Enter your ip address [$default_ip]:" ip
ip=${ip:-$default_ip}

mkdir -p $ONOMY_HOME

# ------------------ Init onomy ------------------

echo "Creating $ONOMY_NODE_NAME validator with chain-id=$CHAINID..."

# Initialize the home directory and add some keys
echo "Initializing chain"
onomyd $ONOMY_CHAINID_FLAG init $ONOMY_NODE_NAME

while read line;
  do echo "Adding genesis account: '${line}'";
  acc_arr=($line); # line example: onomy104n3g00hms7kycg54jml24s84jjt9s2agq6r07 100anom
  onomyd add-genesis-account "${acc_arr[0]}"  "${acc_arr[1]}"
done <<< "$ONOMY_GENESIS_ACCOUNTS"

sed -i "s/\"stake\"/\"$ONOMY_STAKE_DENOM\"/g" $ONOMY_HOME_CONFIG/genesis.json

# add in denom metadata for both native tokens
jq '.app_state.bank.denom_metadata = [{ "name": "NOM", "symbol": "NOM", "description": "NOM coin", "denom_units": [ { "denom": "anom", "exponent": 0 }, { "denom": "nom", "exponent": 18 } ], "base": "anom", "display": "nom" }]' $ONOMY_HOME_CONFIG/genesis.json | sponge $ONOMY_HOME_CONFIG/genesis.json

# add initial DAO balance
jq ".app_state.dao.treasury_balance = [{\"denom\": \"$ONOMY_STAKE_DENOM\", \"amount\": \"$ONOMY_DAO_AMOUNT\"}]" $ONOMY_HOME_CONFIG/genesis.json | sponge $ONOMY_HOME_CONFIG/genesis.json

# disable community_tax
jq ".app_state.distribution.params.community_tax = \"0\"" $ONOMY_HOME_CONFIG/genesis.json | sponge $ONOMY_HOME_CONFIG/genesis.json

# gravity slashing
jq ".app_state.gravity.params.gravity_id = \"gravity-$(cat /dev/urandom | tr -dc 'a-z0-9' | fold -w 16 | head -n 1)\"" $ONOMY_HOME_CONFIG/genesis.json | sponge $ONOMY_HOME_CONFIG/genesis.json
jq ".app_state.gravity.params.slash_fraction_valset = \"0.002\"" $ONOMY_HOME_CONFIG/genesis.json | sponge $ONOMY_HOME_CONFIG/genesis.json
jq ".app_state.gravity.params.slash_fraction_batch = \"0.002\"" $ONOMY_HOME_CONFIG/genesis.json | sponge $ONOMY_HOME_CONFIG/genesis.json
jq ".app_state.gravity.params.slash_fraction_logic_call = \"0.002\"" $ONOMY_HOME_CONFIG/genesis.json | sponge $ONOMY_HOME_CONFIG/genesis.json
jq ".app_state.gravity.params.slash_fraction_bad_eth_signature = \"0.002\"" $ONOMY_HOME_CONFIG/genesis.json | sponge $ONOMY_HOME_CONFIG/genesis.json
jq ".app_state.gravity.params.signed_valsets_window = \"500\"" $ONOMY_HOME_CONFIG/genesis.json | sponge $ONOMY_HOME_CONFIG/genesis.json
jq ".app_state.gravity.params.signed_batches_window = \"500\"" $ONOMY_HOME_CONFIG/genesis.json | sponge $ONOMY_HOME_CONFIG/genesis.json
jq ".app_state.gravity.params.signed_logic_calls_window = \"500\"" $ONOMY_HOME_CONFIG/genesis.json | sponge $ONOMY_HOME_CONFIG/genesis.json

# gravity add bnom -> anom swap settings
jq ".app_state.gravity.params.erc20_to_denom_permanent_swap = {\"erc20\": \"$ONOMY_BNOM_ERC20_ADDRESS\", \"denom\": \"$ONOMY_STAKE_DENOM\"}" $ONOMY_HOME_CONFIG/genesis.json | sponge $ONOMY_HOME_CONFIG/genesis.json
# add change max validators
jq ".app_state.staking.params.max_validators = $ONOMY_MAX_VALIDATORS" $ONOMY_HOME_CONFIG/genesis.json | sponge $ONOMY_HOME_CONFIG/genesis.json

# add change min_global_self_delegation
jq ".app_state.staking.params.min_global_self_delegation = \"$ONOMY_MIN_GLOBAL_SELF_DELEGATION\"" $ONOMY_HOME_CONFIG/genesis.json | sponge $ONOMY_HOME_CONFIG/genesis.json

# add change min gov deposit
jq ".app_state.gov.deposit_params.min_deposit = [{\"denom\": \"$ONOMY_STAKE_DENOM\", \"amount\": \"$ONOMY_GOV_MIN_DEPOSIT\"}]" $ONOMY_HOME_CONFIG/genesis.json | sponge $ONOMY_HOME_CONFIG/genesis.json

# disable distribution withdraw_addr_enabled
jq ".app_state.distribution.params.withdraw_addr_enabled = false" $ONOMY_HOME_CONFIG/genesis.json | sponge $ONOMY_HOME_CONFIG/genesis.json

echo "Adding validator key"
onomyd keys add $ONOMY_VALIDATOR_NAME $ONOMY_KEYRING_FLAG --output json | jq -r '.mnemonic' | pass insert -e keyring-onomy/$ONOMY_VALIDATOR_NAME-mnemonic
echo "Adding orchestrator keys"
onomyd keys add $ONOMY_ETH_ORCHESTRATOR_NAME $ONOMY_KEYRING_FLAG --output json | jq -r '.mnemonic' | pass insert -e keyring-onomy/$ONOMY_ETH_ORCHESTRATOR_NAME-mnemonic

echo "Adding validator addresses to genesis file"
onomyd add-genesis-account "$(onomyd keys show $ONOMY_VALIDATOR_NAME -a $ONOMY_KEYRING_FLAG)" $ONOMY_VALIDATOR_GENESIS_COINS
echo "Adding eth orchestrator addresses to genesis file"
onomyd add-genesis-account "$(onomyd keys show $ONOMY_ETH_ORCHESTRATOR_NAME -a $ONOMY_KEYRING_FLAG)" $ONOMY_ETH_ORCHESTRATOR_GENESIS_COINS

# vesting accounts
./init-master-vesting.sh

echo "Creating gentxs"
onomyd gravity gentx $ONOMY_VALIDATOR_NAME $ONOMY_STAKE_COINS "$ONOMY_ETH_ORCHESTRATOR_VALIDATOR_ADDRESS" \
 "$(onomyd keys show $ONOMY_ETH_ORCHESTRATOR_NAME -a $ONOMY_KEYRING_FLAG)" \
  $ONOMY_MIN_SELF_DELEGATION_FLAG \
  --moniker="Jedi" \
  --ip=$ONOMY_HOST \
  $ONOMY_KEYRING_FLAG $ONOMY_CHAINID_FLAG

echo "Collecting gentxs in $ONOMY_NODE_NAME"
onomyd gravity collect-gentxs

echo "Updating node config"
# change config
crudini --set $ONOMY_NODE_CONFIG p2p addr_book_strict false
crudini --set $ONOMY_NODE_CONFIG p2p external_address "\"tcp://$ip:26656\""
crudini --set $ONOMY_NODE_CONFIG rpc laddr "\"tcp://$ONOMY_HOST:26657\""

crudini --set $ONOMY_APP_CONFIG grpc enable true
crudini --set $ONOMY_APP_CONFIG grpc address "\"$ONOMY_HOST:$ONOMY_GRPC_PORT\""
crudini --set $ONOMY_APP_CONFIG api enable true
crudini --set $ONOMY_APP_CONFIG api swagger true
