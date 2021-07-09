#!/bin/bash
set -eu

echo "running ethereum node"

# Initial dir
CURRENT_WORKING_DIR=$(pwd)
# Name of the network to bootstrap
CHAINID="testchain"
# Name of the gravity artifact
GRAVITY=gravity
# The name of the gravity node
GRAVITY_NODE_NAME="gravity"
# The address to run gravity node
#GRAVITY_HOST="0.0.0.0"
GIT_HUB_USER=$1
GIT_HUB_PASS=$2
GIT_HUB_EMAIL=$3
GIT_HUB_BRANCH=$4
GRAVITY_HOST=$5
# Home folder for gravity config
GRAVITY_HOME="$CURRENT_WORKING_DIR/$CHAINID/$GRAVITY_NODE_NAME"
# Home flag for home folder
GRAVITY_HOME_FLAG="--home $GRAVITY_HOME"
# Prefix of cosmos addresses
GRAVITY_ADDRESS_PREFIX=cosmos
# Gravity chain demons
STAKE_DENOM="stake"

ETH_MINER_PRIVATE_KEY="0xb1bab011e03a9862664706fc3bbaa1b16651528e5f0e7fbfcbfdd8be302a13e7"
ETH_MINER_PUBLIC_KEY="0xBf660843528035a5A4921534E156a27e64B231fE"
# The host of ethereum node
#ETH_HOST="0.0.0.0"
ETH_HOST=$6


echo "Applying contracts"
GRAVITY_DIR=/go/src/github.com/onomyprotocol/gravity-bridge/
cd $GRAVITY_DIR/solidity
ls
echo "running contract-deployer.ts: GRAVITY_HOST: $GRAVITY_HOST ETH_HOST: $ETH_HOST"

npx ts-node \
    contract-deployer.ts \
    --cosmos-node="http://$GRAVITY_HOST:26657" \
    --eth-node="http://$ETH_HOST:8545" \
    --eth-privkey="$ETH_MINER_PRIVATE_KEY" \
    --contract=artifacts/contracts/Gravity.sol/Gravity.json \
    --test-mode=true | grep "Gravity deployed at Address" | grep -Eow '0x[0-9a-fA-F]{40}' >> $GRAVITY_HOME/eth_contract_address

CONTRACT_ADDRESS=$(cat $GRAVITY_HOME/eth_contract_address)
echo "Contract address: $CONTRACT_ADDRESS"

##----------------------------- commit save Contract address-----
cd /root/onomy/
sh deploy/master-cosmos-orchestrator-node/scripts/store-ethereum-contract-info.sh $GIT_HUB_USER $GIT_HUB_PASS $GIT_HUB_EMAIL $GIT_HUB_BRANCH

# return back to home
cd $CURRENT_WORKING_DIR
sleep 10