#!/bin/bash
#Setup Validator
set -eu

echo "Initializing validator"

# Initial dir
ONOMY_HOME=$HOME/.onomy
# Name of the network to bootstrap
CHAINID=$(jq -r .chain_id $ONOMY_HOME/node_info.json)
# Name of the gravity artifact
ONOMY=onomyd
# The name of the gravity node
ONOMY_NODE_NAME=$(jq -r .node_name $ONOMY_HOME/node_info.json)
# Keyring flag
ONOMY_KEYRING_FLAG="--keyring-backend test"
# Gravity chain demons
STAKE_DENOM="nom"
NOM_REQUEST_AMOUNT=11000000
NOM_STAKE_AMOUNT=10000000

# -----------------Adding Validator---------------------

ONOMY_VALIDATOR_NAME=''

if [[ -f "$ONOMY_HOME/validator_key.json" ]]
then
    echo "Validator key already exist $ONOMY_HOME/validator_key.json"
    ONOMY_VALIDATOR_NAME=$(jq -r .name $ONOMY_HOME/validator_key.json)
fi

while [[ $ONOMY_VALIDATOR_NAME == '' ]]
do
   # The name of onomy validator
  read -r -p "Enter a name for your validator [validator]:" ONOMY_VALIDATOR_NAME
  ONOMY_VALIDATOR_NAME=${ONOMY_VALIDATOR_NAME:-validator}
  $ONOMY keys add $ONOMY_VALIDATOR_NAME $ONOMY_KEYRING_FLAG --output json | jq . >> $ONOMY_HOME/validator_key.json
done

# Save validator-info
# Switch sed command in the case of linux
fsed() {
  if [ `uname` = 'Linux' ]; then
    sed -i "$@"
  else
    sed -i '' "$@"
  fi
}

fsed 's#"validator_name": ""#"validator_name": "'$ONOMY_VALIDATOR_NAME'"#g'  $ONOMY_HOME/node_info.json

# -------------------Get Faucet URL----------------

ONOMY_VALIDATOR_ADDRESS=$(jq -r .address $ONOMY_HOME/validator_key.json)

# if the validator doesn't have enough amount use faucet
amount=$(jq -r .amount <<< "$($ONOMY q bank balances --denom $STAKE_DENOM --output json $ONOMY_VALIDATOR_ADDRESS)")
if [ "$amount" -lt $NOM_REQUEST_AMOUNT  ]; then
  read -r -p "Please enter Faucet URL [http://testnet1.onomy.io:8000/]:" FAUCET_TOKEN_BASE_URL
  FAUCET_TOKEN_BASE_URL=${FAUCET_TOKEN_BASE_URL:-http://testnet1.onomy.io:8000/}

  curl -X POST "$FAUCET_TOKEN_BASE_URL" -H  "accept: application/json" -H  "Content-Type: application/json" -d "{  \"address\": \"$ONOMY_VALIDATOR_ADDRESS\",  \"coins\": [    \"$NOM_REQUEST_AMOUNT$STAKE_DENOM\"  ]}"
fi

echo -e '\nWaiting for balance synchronization'

for i in {1..20}; do
  amount=$(jq -r .amount <<< "$($ONOMY q bank balances --denom $STAKE_DENOM --output json $ONOMY_VALIDATOR_ADDRESS)")
  if [ "$amount" -lt $NOM_REQUEST_AMOUNT  ]; then
     sleep 1
     continue
  fi
  break
done

if [ "$amount" -lt $NOM_REQUEST_AMOUNT  ]; then
 echo "The $ONOMY_VALIDATOR_ADDRESS hasn't received $STAKE_DENOM, check the node synchronization"
 exit
fi

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