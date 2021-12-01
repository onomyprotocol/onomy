#!/bin/bash
set -eu

echo "Initializing validator"

# Initial dir
ONOMY_HOME=$HOME/.onomy
# Name of the Onomy artifact
ONOMY=onomyd

# Keyring flag
ONOMY_KEYRING_FLAG="--keyring-backend test"

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
ONOMY_VALIDATOR_ADDRESS=$(jq -r .address $ONOMY_HOME/validator_key.json)

echo "The validator address is: $ONOMY_VALIDATOR_ADDRESS"
