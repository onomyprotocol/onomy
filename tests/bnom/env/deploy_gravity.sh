#!/bin/bash
set -eu

# The address to run onomy node
# The node is running on the host machine be the call to it we expect from the container.
# The hist to make the test pass on mac and linux
ONOMY_HOST="host.docker.internal"
ORIGINAL_DIR=$PWD
if ! ping -c 1 $ONOMY_HOST &> /dev/null
then
  ONOMY_HOST="0.0.0.0"
fi

echo "ONOMY_HOST: $ONOMY_HOST"

# The URL of the running mock eth node.
ETH_ADDRESS="http://0.0.0.0:8545/"

# The ETH key used for orchestrator signing of the transactions
ETH_ORCHESTRATOR_PRIVATE_KEY=c40f62e75a11789dbaf6ba82233ce8a52c20efb434281ae6977bb0b3a69bf709

#-------------------- Deploy the contract --------------------

echo "Deploying Gravity contract"

cd /root/home/gravity/solidity
./contract-deployer \
--cosmos-node="http://$ONOMY_HOST:26657" \
--eth-node="$ETH_ADDRESS" \
--eth-privkey="$ETH_ORCHESTRATOR_PRIVATE_KEY" \
--contract=Gravity.json \
--test-mode=false \
--bnom-address="0xe7c0fd1f0A3f600C1799CD8d335D31efBE90592C" | grep "Gravity deployed at Address" | grep -Eow '0x[0-9a-fA-F]{40}' > gravity_contract_address

GRAVITY_CONTRACT_ADDRESS=$(cat gravity_contract_address)
cp gravity_contract_address $ORIGINAL_DIR/gravity_contract_address

echo "Contract address: $GRAVITY_CONTRACT_ADDRESS"
