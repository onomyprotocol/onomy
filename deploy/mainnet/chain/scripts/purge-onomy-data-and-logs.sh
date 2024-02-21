#!/bin/bash
set -eu

echo "Purging onomy data and logs"

# Onomy home dir
ONOMY_HOME=$HOME/.onomy
# Onomy data
ONOMY_DATA=$ONOMY_HOME/data/
# Onomy logs
ONOMY_LOGS=$ONOMY_HOME/logs/
# Priv validator file path
ONOMY_PRIV_VALIDATOR_FILE_PATH=$ONOMY_DATA/priv_validator_state.json

rm -rf $ONOMY_LOGS/*
echo "$ONOMY_LOGS is empty"

rm -rf $ONOMY_DATA/*
echo "$ONOMY_DATA is empty"

cp assets/node/priv_validator_state.json $ONOMY_PRIV_VALIDATOR_FILE_PATH
echo -e "$ONOMY_PRIV_VALIDATOR_FILE_PATH file is created in the $ONOMY_DATA \n"
cat $ONOMY_PRIV_VALIDATOR_FILE_PATH
echo -e "\n"