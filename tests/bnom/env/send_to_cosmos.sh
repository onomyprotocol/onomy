#!/bin/bash
set -eu

# The prefix for cosmos addresses
ONOMY_ADDRESS_PREFIX="onomy"

# The URL of the running mock eth node.
ETH_ADDRESS="http://0.0.0.0:8545/"

# The ETH key used for orchestrator signing of the transactions.
ETH_ORCHESTRATOR_PRIVATE_KEY=c40f62e75a11789dbaf6ba82233ce8a52c20efb434281ae6977bb0b3a69bf709

# read already deployed address from the file
GRAVITY_CONTRACT_ADDRESS=$(cat gravity_contract_address)

ERC20_TOKEN_ADDRESS="$1"
ERC20_TOKEN_AMOUNT="$2"
COSMOS_DESTINATION_ADDRESS="$3"

gbt -a "$ONOMY_ADDRESS_PREFIX" client eth-to-cosmos \
        --ethereum-key "$ETH_ORCHESTRATOR_PRIVATE_KEY" \
        --ethereum-rpc="$ETH_ADDRESS" \
        --gravity-contract-address "$GRAVITY_CONTRACT_ADDRESS" \
        --token-contract-address "$ERC20_TOKEN_ADDRESS" \
        --amount="$ERC20_TOKEN_AMOUNT" \
        --destination "$COSMOS_DESTINATION_ADDRESS"

