#!/bin/bash
#Setup Validator
set -eu

echo "Initializing validator"

# Initial dir
ONOMY_HOME=$HOME/.onomy
# Name of the network to bootstrap
CHAINID=$(jq -r .chain_id $ONOMY_HOME/node_info.json)
# Name of the onomy artifact
ONOMY=onomyd
# The name of the validator node
ONOMY_VALIDATOR_NAME=$(jq -r .validator_name $ONOMY_HOME/node_info.json)
# Keyring flag
ONOMY_KEYRING_FLAG="--keyring-backend test"
# Gravity chain demons
STAKE_DENOM="anom"

echo -e "\n Creating validator:"

# Store the public key of validator
PUB_KEY=$($ONOMY tendermint show-validator)

# Do the create validator transaction
$ONOMY tx staking create-validator \
--amount=$NOM_STAKE_AMOUNT$STAKE_DENOM \
--pubkey="$PUB_KEY" \
--moniker=\"$ONOMY_VALIDATOR_NAME\" \
--chain-id=$CHAINID \
--commission-rate="0.10" \
--commission-max-rate="0.20" \
--commission-max-change-rate="0.01" \
--min-self-delegation="10" \
--gas="auto" \
--gas-adjustment=1.5 \
--gas-prices="1$STAKE_DENOM" \
--from=$ONOMY_VALIDATOR_NAME \
$ONOMY_KEYRING_FLAG -y

sleep 2
echo -e "\n Validator information:"

$ONOMY query staking validator "$($ONOMY keys show $ONOMY_VALIDATOR_NAME --bech val --address $ONOMY_KEYRING_FLAG)"

echo -e "\nYour node is validating if [status: BOND_STATUS_BONDED]."
echo -e  "User keys are located in the file $ONOMY_HOME/validator_key.json"
echo -e  "Private validator keys in the file $ONOMY_HOME/config/priv_validator_key.json"