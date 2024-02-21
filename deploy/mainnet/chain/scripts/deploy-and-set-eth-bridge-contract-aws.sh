#!/bin/bash
set -eu

echo "Deploying eth bridge contract"

ONOMY_HOME=$HOME/.onomy

# The address to run onomy node
ONOMY_HOST="0.0.0.0"
# The address of BNOM ERC20 token on ethereum.
ONOMY_BNOM_ERC20_ADDRESS="0x4fdF157C6860F33232608198419c74ED446Ca577"
# The json file with all bridges addresses
ONOMY_CONTRACT_ADDRESSES_PATH="assets/bridge/addresses.json"

#-------------------- Deploy the contract --------------------

if [[ -z "${ETH_RPC_ADDRESS}" ]]; then
  echo "Fail: ETH_RPC_ADDRESS is not provided"
  exit 1
fi

if [[ -z "${CONTRACT_DEPLOYER_AWS_KEYS_NAME}" ]]; then
  echo "Fail: CONTRACT_DEPLOYER_AWS_KEYS_NAME is not provided"
  exit 1
fi

ONOMY_ETH_DEPLOYER_PRIVATE_KEY=$(aws secretsmanager get-secret-value --secret-id $CONTRACT_DEPLOYER_AWS_KEYS_NAME | jq --raw-output '.SecretString' | jq -r '.eth_private_key_0')
if [[ -z "${ONOMY_ETH_DEPLOYER_PRIVATE_KEY}" ]]; then
  echo "Fail: check if deployer key exists in aws store"
  exit 1
fi

echo "Deploying using ETH_RPC_ADDRESS: $ETH_RPC_ADDRESS"

deploy_response=$(contract-deployer \
--cosmos-node="http://$ONOMY_HOST:26657" \
--eth-node="$ETH_RPC_ADDRESS" \
--eth-privkey="$ONOMY_ETH_DEPLOYER_PRIVATE_KEY" \
--contract=$ONOMY_HOME/contracts/eth-bridge/Gravity.json \
--test-mode=false \
--bnom-address=$ONOMY_BNOM_ERC20_ADDRESS)

ONOMY_ETH_BRIDGE_CONTRACT_ADDRESS=$(echo "$deploy_response" | grep "Gravity deployed at Address"  | grep -Eow '0x[0-9a-fA-F]{40}')

if [[ -z "${ONOMY_ETH_BRIDGE_CONTRACT_ADDRESS}" ]]; then
  echo "Something went wrong: $deploy_response"
  exit 1
fi

echo "Contract deployed successfully, address: $ONOMY_ETH_BRIDGE_CONTRACT_ADDRESS;"

jq ".ethereum = \"$ONOMY_ETH_BRIDGE_CONTRACT_ADDRESS\"" $ONOMY_CONTRACT_ADDRESSES_PATH | sponge $ONOMY_CONTRACT_ADDRESSES_PATH