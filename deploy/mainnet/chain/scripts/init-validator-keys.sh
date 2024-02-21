#!/bin/bash
set -eu

echo "Initializing validator keys"

# Keyring flag
ONOMY_KEYRING_FLAG="--keyring-backend pass"
# The name of the onomy validator
ONOMY_VALIDATOR_NAME="validator"
# The name of the onomy eth-orchestrator/validator
ONOMY_ETH_ORCHESTRATOR_NAME="eth-orchestrator"

# -----------------Adding Validator---------------------

if onomyd keys show $ONOMY_VALIDATOR_NAME -a $ONOMY_KEYRING_FLAG &> /dev/null; then
    echo "The $ONOMY_VALIDATOR_NAME key is already set"
else
    onomyd keys add $ONOMY_VALIDATOR_NAME $ONOMY_KEYRING_FLAG --output json | jq -r '.mnemonic' | pass insert -e keyring-onomy/$ONOMY_VALIDATOR_NAME-mnemonic
fi

if onomyd keys show $ONOMY_ETH_ORCHESTRATOR_NAME -a $ONOMY_KEYRING_FLAG &> /dev/null; then
    echo "The $ONOMY_ETH_ORCHESTRATOR_NAME key is already set"
else
    onomyd keys add $ONOMY_ETH_ORCHESTRATOR_NAME $ONOMY_KEYRING_FLAG --output json | jq -r '.mnemonic' | pass insert -e keyring-onomy/$ONOMY_ETH_ORCHESTRATOR_NAME-mnemonic
fi

# Show validator-info

echo "The validator address is: $(onomyd keys show $ONOMY_VALIDATOR_NAME -a $ONOMY_KEYRING_FLAG)"
echo "The eth-orchestrator address is: $(onomyd keys show $ONOMY_ETH_ORCHESTRATOR_NAME -a $ONOMY_KEYRING_FLAG)"
