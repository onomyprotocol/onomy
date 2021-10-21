#!/bin/bash
#Setup Validator
set -eu

echo "Initializing validator"

# Initial dir
ONOMY_HOME=$HOME/.onomy
# Name of the network to bootstrap
CHAINID=$(jq .chain_id $ONOMY_HOME/node_info.json | sed 's#\"##g')
# Name of the gravity artifact
ONOMY=onomyd
# The name of the gravity node
ONOMY_NODE_NAME=$(jq .node_name $ONOMY_HOME/node_info.json | sed 's#\"##g')
# Keyring flag
ONOMY_KEYRING_FLAG="--keyring-backend test"
# Gravity chain demons
STAKE_DENOM="nom"
NOM_REQUEST_AMOUNT=20000000
NOM_STAKE_AMOUNT=10000000

# -----------------Adding Validator---------------------

ONOMY_VALIDATOR_NAME=''
while [[ $ONOMY_VALIDATOR_NAME == '' ]]
do
   # The name of onomy validator
  read -p "Enter a name for your validator: " ONOMY_VALIDATOR_NAME
done

if [[ -f "$ONOMY_HOME/validator_key.json" ]]
then
    echo "Validator key already exist in $ONOMY_HOME/validator_key.json"
    exit
fi

$ONOMY keys add $ONOMY_VALIDATOR_NAME $ONOMY_KEYRING_FLAG --output json | jq . >> $ONOMY_HOME/validator_key.json
jq .mnemonic $ONOMY_HOME/validator_key.json | sed 's#\"##g' >> $HOME/validator-phrases

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
read -p "Please enter Faucet URL: " -i "http://testnet1.onomy.io:8000/" -e FAUCET_TOKEN_BASE_URL

# ------------------Get Tokens from Faucet------------------

ONOMY_VALIDATOR_ADDRESS=$(jq -r .address $ONOMY_HOME/validator_key.json)
curl -X POST "$FAUCET_TOKEN_BASE_URL" -H  "accept: application/json" -H  "Content-Type: application/json" -d "{  \"address\": \"$ONOMY_VALIDATOR_ADDRESS\",  \"coins\": [    \"$NOM_REQUEST_AMOUNT$STAKE_DENOM\"  ]}"

echo -e '\nWaiting for balance synchronization'

for i in {1..600}; do
  amount=$(jq .amount <<< "$($ONOMY q bank balances --denom $STAKE_DENOM --output json $ONOMY_VALIDATOR_ADDRESS)" | sed 's#\"##g')
  if [ "$amount" -lt $NOM_REQUEST_AMOUNT  ]; then
     sleep 1
     continue
  fi
  break
done

if [ "$amount" -lt $NOM_REQUEST_AMOUNT  ]; then
 echo "The node hasn't received $STAKE_DENOM, check the node synchronization"
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