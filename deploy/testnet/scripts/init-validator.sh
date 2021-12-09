#!/bin/bash
set -eu

echo "Initializing validator"

# Initial dir
ONOMY_HOME=$HOME/.onomy

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
  onomyd keys add $ONOMY_VALIDATOR_NAME $ONOMY_KEYRING_FLAG --output json | jq . >> $ONOMY_HOME/validator_key.json
done


ONOMY_ORCHESTRATOR_NAME=''

if [[ -f "$ONOMY_HOME/orchestrator_key.json" ]]
then
    echo "Orchestrator key already exist $ONOMY_HOME/orchestrator_key.json"
    ONOMY_ORCHESTRATOR_NAME=$(jq -r .name $ONOMY_HOME/orchestrator_key.json)
fi

read -r -p "Enter a name for your orchestrator [orchestrator]:" ONOMY_ORCHESTRATOR_NAME
ONOMY_ORCHESTRATOR_NAME=${ONOMY_ORCHESTRATOR_NAME:-orchestrator}
onomyd keys add $ONOMY_ORCHESTRATOR_NAME $ONOMY_KEYRING_FLAG --output json | jq . >> $ONOMY_HOME/orchestrator_key.json

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

fsed 's#"orchestrator_name": ""#"orchestrator_name": "'$ONOMY_ORCHESTRATOR_NAME'"#g'  $ONOMY_HOME/node_info.json
ONOMY_ORCHESTRATOR_ADDRESS=$(jq -r .address $ONOMY_HOME/orchestrator_key.json)
echo "The orchestrator address is: $ONOMY_ORCHESTRATOR_ADDRESS"