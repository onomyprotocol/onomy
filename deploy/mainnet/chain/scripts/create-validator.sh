#!/bin/bash
set -eu

echo "Creating validator"

# Onomy home dir
ONOMY_HOME=$HOME/.onomy
# Name of the network to bootstrap
CHAINID="onomy-mainnet-1"
# The name of the validator node
ONOMY_VALIDATOR_NAME=validator
# Keyring flag
ONOMY_KEYRING_FLAG="--keyring-backend pass"
# Onomy chain denom
ONOMY_STAKE_DENOM="anom"

while read -r -p "Provide the validator moniker (name):" ONOMY_MONIKER; do
    if [ -z $ONOMY_MONIKER ]; then
        echo "The moniker name is required."
        continue
    fi
    break
done

read -r -p "Provide the validator web site (optional):" ONOMY_VALIDATOR_WEB_SITE;
read -r -p "Provide the validator identity (https://keybase.io/ identity key) (optional):" ONOMY_VALIDATOR_IDENTITY;
read -r -p "Provide the validator commission rate [0.10]:" ONOMY_VALIDATOR_COMMISSION_RATE;
ONOMY_VALIDATOR_COMMISSION_RATE=${ONOMY_VALIDATOR_COMMISSION_RATE:-"0.10"}
read -r -p "Provide the validator commission max rate [0.20]:" ONOMY_VALIDATOR_COMMISSION_MAX_RATE;
ONOMY_VALIDATOR_COMMISSION_MAX_RATE=${ONOMY_VALIDATOR_COMMISSION_MAX_RATE:-"0.20"}
read -r -p "Provide the validator commission max change rate [0.01]:" ONOMY_VALIDATOR_COMMISSION_MAX_CHANGE_RATE;
ONOMY_VALIDATOR_COMMISSION_MAX_CHANGE_RATE=${ONOMY_VALIDATOR_COMMISSION_MAX_CHANGE_RATE:-"0.01"}
read -r -p "Provide the validator min self delegation amount of aNOMs it should be more or equal 225000000000000000000000 [225000000000000000000000]:" ONOMY_MIN_DELEGATION_AMOUNT;
ONOMY_MIN_DELEGATION_AMOUNT=${ONOMY_MIN_DELEGATION_AMOUNT:-"225000000000000000000000"}
read -r -p "Provide the validator delegation amount of aNOMs, it should be more or equal min self delegation [225000000000000000000000]:" ONOMY_DELEGATION_AMOUNT;
ONOMY_DELEGATION_AMOUNT=${ONOMY_DELEGATION_AMOUNT:-"225000000000000000000000"}

# Store the public key of validator
PUB_KEY=$(onomyd tendermint show-validator)

# Do the create validator transaction
onomyd tx staking create-validator \
--amount=$ONOMY_DELEGATION_AMOUNT$ONOMY_STAKE_DENOM \
--pubkey="$PUB_KEY" \
--moniker="$ONOMY_MONIKER" \
--website="$ONOMY_VALIDATOR_WEB_SITE" \
--identity="$ONOMY_VALIDATOR_IDENTITY" \
--commission-rate="$ONOMY_VALIDATOR_COMMISSION_RATE" \
--commission-max-rate="$ONOMY_VALIDATOR_COMMISSION_MAX_RATE" \
--commission-max-change-rate="$ONOMY_VALIDATOR_COMMISSION_MAX_CHANGE_RATE" \
--min-self-delegation=$ONOMY_MIN_DELEGATION_AMOUNT \
--gas="auto" \
--gas-adjustment=1.5 \
--chain-id="$CHAINID" \
--from=$ONOMY_VALIDATOR_NAME \
$ONOMY_KEYRING_FLAG -y

sleep 5
echo -e "\n Validator information:"

onomyd query staking validator "$(onomyd keys show $ONOMY_VALIDATOR_NAME --bech val --address $ONOMY_KEYRING_FLAG)"

echo -e "Your node is validating if [status: BOND_STATUS_BONDED]."
echo -e  "Private validator keys in the file $ONOMY_HOME/config/priv_validator_key.json"
