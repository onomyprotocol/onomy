#!/bin/bash
set -eu

echo "Deploying gravity contract"

ONOMY_HOME=$HOME/.onomy
GRAVITY_SRC=$ONOMY_HOME/src/gravity

# The address to run onomy node
ONOMY_HOST="0.0.0.0"

#-------------------- Deploy the contract --------------------

if [[ -z "${ETH_RPC_ADDRESS}" ]]; then
  echo "Fail: ETH_RPC_ADDRESS is not set"
  exit
fi

if [[ -z "${ETH_PRIVATE_KEY}" ]]; then
  echo "Fail: ETH_PRIVATE_KEY is not set"
  exit
fi

echo "ETH_RPC_ADDRESS: $ETH_RPC_ADDRESS and ETH_PRIVATE_KEY: $ETH_PRIVATE_KEY"

cd $GRAVITY_SRC/solidity

deploy_response=$(./contract-deployer \
--cosmos-node="http://$ONOMY_HOST:26657" \
--eth-node="$ETH_RPC_ADDRESS" \
--eth-privkey="$ETH_PRIVATE_KEY" \
--contract=artifacts/contracts/Gravity.sol/Gravity.json \
--test-mode=false)

CONTRACT_ADDRESS=$(echo "$deploy_response" | grep "Gravity deployed at Address"  | grep -Eow '0x[0-9a-fA-F]{40}')

if [[ -z "${CONTRACT_ADDRESS}" ]]; then
  echo "Something went wrong: $deploy_response"
  exit
fi

echo $CONTRACT_ADDRESS >> $ONOMY_HOME/eth_contract_address

echo "Contract deployed successfully, address: $CONTRACT_ADDRESS, file path: $ONOMY_HOME/eth_contract_address"

cd $ONOMY_HOME