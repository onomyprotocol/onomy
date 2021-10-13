#!/bin/bash
set -eu

echo "building environment"
# Initial dir
ONOMY_HOME=$HOME/.onomy

# Name of the network to bootstrap
#echo "Enter chain-id"
#read chainid
CHAINID="onomy-testnet1"
# Name of the gravity artifact
GRAVITY=onomyd
# The name of the gravity node
GRAVITY_NODE_NAME="onomy"
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
# The name of the gravity validator
GRAVITY_VALIDATOR_NAME=validator1
# The name of the gravity orchestrator/validator
GRAVITY_ORCHESTRATOR_NAME=orch
# Gravity chain demons
STAKE_DENOM="nom"
#NORMAL_DENOM="samoleans"
NORMAL_DENOM="footoken"
# The port of the onomy gRPC
GRAVITY_GRPC_PORT="9090"

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
# ------------------ Init gravity ------------------

echo "Creating $GRAVITY_NODE_NAME validator with chain-id=$CHAINID..."
echo "Initializing genesis files"
# Build genesis file incl account for passed address
GRAVITY_GENESIS_COINS="1000000000000$STAKE_DENOM,1000000000000$NORMAL_DENOM"

# Initialize the home directory and add some keys
echo "Init test chain"
$GRAVITY $GRAVITY_HOME_FLAG $GRAVITY_CHAINID_FLAG init $GRAVITY_NODE_NAME

echo "Set stake/mint demon to $STAKE_DENOM"
fsed "s#\"stake\"#\"$STAKE_DENOM\"#g" $GRAVITY_HOME_CONFIG/genesis.json

# add in denom metadata for both native tokens
jq '.app_state.bank.denom_metadata += [{"base": "footoken", display: "mfootoken", "description": "A non-staking test token", "denom_units": [{"denom": "footoken", "exponent": 0}, {"denom": "mfootoken", "exponent": 6}]}, {"base": "nom", display: "mnom", "description": "A staking test token", "denom_units": [{"denom": "nom", "exponent": 0}, {"denom": "mnom", "exponent": 6}]}]' $GRAVITY_HOME_CONFIG/genesis.json > $ONOMY_HOME/metadata-genesis.json

# a 60 second voting period to allow us to pass governance proposals in the tests
jq '.app_state.gov.voting_params.voting_period = "60s"' $ONOMY_HOME/metadata-genesis.json > $ONOMY_HOME/edited-genesis.json
mv $ONOMY_HOME/edited-genesis.json $ONOMY_HOME/genesis.json
mv $ONOMY_HOME/genesis.json $GRAVITY_HOME_CONFIG/genesis.json

echo "Add validator key"
$GRAVITY $GRAVITY_HOME_FLAG keys add $GRAVITY_VALIDATOR_NAME $GRAVITY_KEYRING_FLAG --output json | jq . >> $GRAVITY_HOME/validator_key.json
jq .mnemonic $GRAVITY_HOME/validator_key.json | sed 's#\"##g' > $ONOMY_HOME/validator-phrases

echo "Generating orchestrator keys"
$GRAVITY $GRAVITY_HOME_FLAG keys add --output=json $GRAVITY_ORCHESTRATOR_NAME $GRAVITY_KEYRING_FLAG | jq . >> $GRAVITY_HOME/orchestrator_key.json
jq .mnemonic $GRAVITY_HOME/orchestrator_key.json | sed 's#\"##g' > $ONOMY_HOME/orchestrator-phrases

echo "Adding validator addresses to genesis files"
$GRAVITY $GRAVITY_HOME_FLAG add-genesis-account "$($GRAVITY $GRAVITY_HOME_FLAG keys show $GRAVITY_VALIDATOR_NAME -a $GRAVITY_KEYRING_FLAG)" $GRAVITY_GENESIS_COINS
echo "Adding orchestrator addresses to genesis files"
$GRAVITY $GRAVITY_HOME_FLAG add-genesis-account "$($GRAVITY $GRAVITY_HOME_FLAG keys show $GRAVITY_ORCHESTRATOR_NAME -a $GRAVITY_KEYRING_FLAG)" $GRAVITY_GENESIS_COINS

echo "Adding faucet account addresses to genesis files"
$GRAVITY $GRAVITY_HOME_FLAG keys add --output=json faucet_account1 $GRAVITY_KEYRING_FLAG | jq . >> $GRAVITY_HOME/faucet_account1.json
$GRAVITY $GRAVITY_HOME_FLAG keys add --output=json faucet_account2 $GRAVITY_KEYRING_FLAG | jq . >> $GRAVITY_HOME/faucet_account2.json
$GRAVITY $GRAVITY_HOME_FLAG keys add --output=json faucet_account3 $GRAVITY_KEYRING_FLAG | jq . >> $GRAVITY_HOME/faucet_account3.json
$GRAVITY $GRAVITY_HOME_FLAG keys add --output=json faucet_account4 $GRAVITY_KEYRING_FLAG | jq . >> $GRAVITY_HOME/faucet_account4.json
$GRAVITY $GRAVITY_HOME_FLAG keys add --output=json faucet_account5 $GRAVITY_KEYRING_FLAG | jq . >> $GRAVITY_HOME/faucet_account5.json
$GRAVITY $GRAVITY_HOME_FLAG add-genesis-account "$($GRAVITY $GRAVITY_HOME_FLAG keys show faucet_account1 -a $GRAVITY_KEYRING_FLAG)" $GRAVITY_GENESIS_COINS
$GRAVITY $GRAVITY_HOME_FLAG add-genesis-account "$($GRAVITY $GRAVITY_HOME_FLAG keys show faucet_account2 -a $GRAVITY_KEYRING_FLAG)" $GRAVITY_GENESIS_COINS
$GRAVITY $GRAVITY_HOME_FLAG add-genesis-account "$($GRAVITY $GRAVITY_HOME_FLAG keys show faucet_account3 -a $GRAVITY_KEYRING_FLAG)" $GRAVITY_GENESIS_COINS
$GRAVITY $GRAVITY_HOME_FLAG add-genesis-account "$($GRAVITY $GRAVITY_HOME_FLAG keys show faucet_account4 -a $GRAVITY_KEYRING_FLAG)" $GRAVITY_GENESIS_COINS
$GRAVITY $GRAVITY_HOME_FLAG add-genesis-account "$($GRAVITY $GRAVITY_HOME_FLAG keys show faucet_account5 -a $GRAVITY_KEYRING_FLAG)" $GRAVITY_GENESIS_COINS
#echo "Adding orchestrator keys to genesis"
#GRAVITY_ORCHESTRATOR_KEY="$(jq .address $GRAVITY_HOME/orchestrator_key.json)"

echo "Generating ethereum keys"
$GRAVITY $GRAVITY_HOME_FLAG eth_keys add --output=json | jq . >> $GRAVITY_HOME/eth_key.json
echo "private: $(jq .private_key $GRAVITY_HOME/eth_key.json | sed 's#\"##g')" > $ONOMY_HOME/validator-eth-keys
echo "public: $(jq .public_key $GRAVITY_HOME/eth_key.json | sed 's#\"##g')" >> $ONOMY_HOME/validator-eth-keys
echo "address: $(jq .address $GRAVITY_HOME/eth_key.json | sed 's#\"##g')" >> $ONOMY_HOME/validator-eth-keys

echo "Creating gentxs"
$GRAVITY $GRAVITY_HOME_FLAG gentx --ip $GRAVITY_HOST $GRAVITY_VALIDATOR_NAME 1000000000000$STAKE_DENOM "$(jq -r .address $GRAVITY_HOME/eth_key.json)" "$(jq -r .address $GRAVITY_HOME/orchestrator_key.json)" $GRAVITY_KEYRING_FLAG $GRAVITY_CHAINID_FLAG

echo "Collecting gentxs in $GRAVITY_NODE_NAME"
$GRAVITY $GRAVITY_HOME_FLAG collect-gentxs

#jq ".app_state.bank.supply += [{\"denom\": \"$STAKE_DENOM\",\"amount\": \"2000000000000\"},{\"denom\": \"$NORMAL_DENOM\",\"amount\": \"2000000000000\"}]" $GRAVITY_HOME_CONFIG/genesis.json | sponge $GRAVITY_HOME_CONFIG/genesis.json

echo "Exposing ports and APIs of the $GRAVITY_NODE_NAME"

# Change ports
fsed "s#\"tcp://127.0.0.1:26656\"#\"tcp://$GRAVITY_HOST:26656\"#g" $GRAVITY_NODE_CONFIG
fsed "s#\"tcp://127.0.0.1:26657\"#\"tcp://$GRAVITY_HOST:26657\"#g" $GRAVITY_NODE_CONFIG
fsed 's#addr_book_strict = true#addr_book_strict = false#g' $GRAVITY_NODE_CONFIG
fsed 's#external_address = ""#external_address = "tcp://'$GRAVITY_HOST:26656'"#g' $GRAVITY_NODE_CONFIG
fsed 's#enable = false#enable = true#g' $GRAVITY_APP_CONFIG
fsed 's#swagger = false#swagger = true#g' $GRAVITY_APP_CONFIG

# Save validator-info
fsed 's#"validator_name": ""#"validator_name": "'$GRAVITY_VALIDATOR_NAME'"#g'  $ONOMY_HOME/val_info.json
fsed 's#"chain_id": ""#"chain_id": "'$CHAINID'"#g'  $ONOMY_HOME/val_info.json


#echo "Adding initial ethereum value for gravity validator"
#jq ".alloc |= . + {$(jq .address $GRAVITY_HOME/eth_key.json) : {\"balance\": \"0x1337000000000000000000\"}}" $HOME/market/deploy/redhat-testchain-deployment/assets/ETHGenesis.json | sponge $HOME/market/deploy/redhat-testchain-deployment/assets/ETHGenesis.json

$GRAVITY $GRAVITY_HOME_FLAG start --pruning=nothing &

#echo "Waiting $GRAVITY_NODE_NAME to launch gRPC $GRAVITY_GRPC_PORT..."

#while ! timeout 1 bash -c "</dev/tcp/0.0.0.0/9090"; do
  sleep 3
#done

#echo "$GRAVITY_NODE_NAME launched"

# ------------------ Run faucet ------------------
ONOMY_ORCHESTRATOR_MNEMONIC=$(jq -r .mnemonic $GRAVITY_HOME/orchestrator_key.json)

echo "delete faucet account if we have previously"
$GRAVITY keys delete faucet --keyring-backend test -y &

echo "Starting faucet based on validator account"
faucet -cli-name=$GRAVITY -keyring-backend=test -mnemonic="$ONOMY_ORCHESTRATOR_MNEMONIC" credit-amount=100000000  -max-credit=200000000 &
