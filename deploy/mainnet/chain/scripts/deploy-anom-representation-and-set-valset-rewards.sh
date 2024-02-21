#!/bin/bash

echo "Deploying ERC20 representation (might take few minutes)"

# The prefix for cosmos addresses
ONOMY_ADDRESS_PREFIX="onomy"
# Home folder for onomy config
ONOMY_HOME="$HOME/.onomy"
# The address to run onomy node
ONOMY_HOST="0.0.0.0"
# The port of the onomy gRPC
ONOMY_GRPC_PORT="9191"
# The json file with all bridges addresses
ONOMY_CONTRACT_ADDRESSES_PATH="assets/bridge/addresses.json"
# Onomy chain denom
ONOMY_STAKE_DENOM="anom"
# Config directories for onomy node
ONOMY_HOME_CONFIG="$ONOMY_HOME/config"

# The gravity valset reward 100noms
ONOMY_GRAVITY_VALSET_REWARD_AMOUNT="100000000000000000000"

#-------------------- Deploy the representation --------------------

if [[ -z "${ETH_RPC_ADDRESS}" ]]; then
  echo "Fail: ETH_RPC_ADDRESS is not provided"
  exit 1
fi

ONOMY_ETH_DEPLOYER_PRIVATE_KEY=$(pass keyring-onomy/eth-deployer-private-key)
if [[ -z "${ONOMY_ETH_DEPLOYER_PRIVATE_KEY}" ]]; then
  echo "Fail: check if key exists in pass: keyring-onomy/eth-deployer-private-key"
  exit 1
fi

ONOMY_ETH_BRIDGE_CONTRACT_ADDRESS=$(jq -r .ethereum $ONOMY_CONTRACT_ADDRESSES_PATH)
if [[ -z "${ONOMY_ETH_BRIDGE_CONTRACT_ADDRESS}" ]]; then
  echo "Fail: ONOMY_ETH_BRIDGE_CONTRACT_ADDRESS is empty, check the $ONOMY_CONTRACT_ADDRESSES_PATH file"
  exit 1
fi

ONOMY_STAKE_DENOM_ERC20_CONTRACT_ADDRESS=""
deploy_response=$(gbt -a $ONOMY_ADDRESS_PREFIX client deploy-erc20-representation \
                           --ethereum-key="$ONOMY_ETH_DEPLOYER_PRIVATE_KEY" \
                           --cosmos-grpc="http://$ONOMY_HOST:$ONOMY_GRPC_PORT/" \
                           --ethereum-rpc="$ETH_RPC_ADDRESS" \
                           --cosmos-denom="$ONOMY_STAKE_DENOM" \
                           --gravity-contract-address="$ONOMY_ETH_BRIDGE_CONTRACT_ADDRESS" 2>&1)

ONOMY_STAKE_DENOM_ERC20_CONTRACT_ADDRESS=$(echo "$deploy_response" | grep "ERC20 representation"  | grep -Eow '0x[0-9a-fA-F]{40}')
if [[ -z "${ONOMY_STAKE_DENOM_ERC20_CONTRACT_ADDRESS}" ]]; then
  echo -e "Something went wrong, Error: \n $deploy_response"
  echo -e "\n ------------------------------ \n If you see the \"max fee per gas less than block base fee\" error, try again."
  exit 1
fi

echo "The ERC20 of $ONOMY_STAKE_DENOM deployed successfully, address: $ONOMY_STAKE_DENOM_ERC20_CONTRACT_ADDRESS"

echo "Updating the gravity genesis params"

jq ".app_state.gravity.params.valset_reward = {\"denom\": \"$ONOMY_STAKE_DENOM\", \"amount\": \"$ONOMY_GRAVITY_VALSET_REWARD_AMOUNT\"}" $ONOMY_HOME_CONFIG/genesis.json | sponge $ONOMY_HOME_CONFIG/genesis.json
jq ".app_state.gravity.erc20_to_denoms = [{\"denom\": \"$ONOMY_STAKE_DENOM\", \"erc20\": \"$ONOMY_STAKE_DENOM_ERC20_CONTRACT_ADDRESS\"}]" $ONOMY_HOME_CONFIG/genesis.json | sponge $ONOMY_HOME_CONFIG/genesis.json

echo "Updated successfully"
