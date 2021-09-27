#!/bin/bash
set -eu

echo "running onomy with bridge"

# Initial dir
CURRENT_WORKING_DIR=$(pwd)
# Name of the network to bootstrap
CHAINID="onomy"
# Name of the onomy artifact
ONOMY=onomyd
# The name of the onomy node
ONOMY_NODE_NAME="onomy"
# The prefix for cosmos addresses
ONOMY_ADDRESS_PREFIX="onomy"
# The address to run onomy node
ONOMY_HOST="0.0.0.0"
# The port of the onomy gRPC
ONOMY_GRPC_PORT="9090"
# Home folder for onomy config
ONOMY_HOME="$CURRENT_WORKING_DIR/$CHAINID/$ONOMY_NODE_NAME"
# Home flag for home folder
ONOMY_HOME_FLAG="--home $ONOMY_HOME"
# Onomy mnemonic used for orchestrator signing of the transactions (orchestrator_key.json file)
ONOMY_ORCHESTRATOR_MNEMONIC=$(jq -r .mnemonic $ONOMY_HOME/orchestrator_key.json)
# Onomy mnemonic used for fauset from (validator_key.json file)
ONOMY_VALIDATOR_MNEMONIC=$(jq -r .mnemonic $ONOMY_HOME/validator_key.json)

# Onomy chain demons
STAKE_DENOM="nom"

# The ETH key used for orchestrator signing of the transactions
ETH_ORCHESTRATOR_PRIVATE_KEY=c40f62e75a11789dbaf6ba82233ce8a52c20efb434281ae6977bb0b3a69bf709

# The host of ethereum node
ETH_HOST="0.0.0.0"
# Eth node rpc port
ETH_RPC_PORT="8545"

ETH_CONTRACT_ADDRESS=0x8A0814b7251138Dea19054425D0dfF0C497305d3

# ------------------ Run onomy ------------------

echo "Starting $ONOMY_NODE_NAME"
$ONOMY $ONOMY_HOME_FLAG start --pruning=nothing &

echo "Waiting $ONOMY_NODE_NAME to launch gRPC $ONOMY_GRPC_PORT..."

while ! timeout 1 bash -c "</dev/tcp/$ONOMY_HOST/$ONOMY_GRPC_PORT"; do
  sleep 1
done

echo "$ONOMY_NODE_NAME launched"

# ------------------ Run fauset ------------------

echo "Starting fauset based on validator account"
faucet -cli-name=$ONOMY -mnemonic="$ONOMY_VALIDATOR_MNEMONIC" &

#-------------------- Run ethereum (geth) --------------------

echo "Starting rinkeby eth node"

geth --rinkeby --syncmode "light" \
                               --http \
                               --http.port "$ETH_RPC_PORT" \
                               --http.addr "$ETH_HOST" \
                               --http.corsdomain "*" \
                               --http.vhosts "*" \
                               &

GETH_IPC_PATH="/root/.ethereum/rinkeby/geth.ipc"
GETH_CONSOLE="geth --rinkeby attach ipc:$GETH_IPC_PATH console --exec"

# 600000 sec to run light node
for i in {1..600000}; do
  sleep 1
  echo "attempt $i to start the eth node"

  netPeerCount=$($GETH_CONSOLE "net.peerCount")
  if [ "$netPeerCount" -lt 3  ]; then
     echo "net.peerCount : $netPeerCount"
     continue
  fi

  ethSyncing=$($GETH_CONSOLE "eth.syncing")
  if [ "$ethSyncing" != "false"  ]; then
     echo "eth.syncing : $ethSyncing"
     continue
  fi

  if [ $i -eq 600000 ]; then
     echo "timeout for ethereum light node exceed"
     exit
  fi

  break
done
set -e

echo "ethereum light node is ready, peers: $($GETH_CONSOLE "net.peerCount")"

#-------------------- Run orchestrator --------------------

echo "Starting orchestrator"

gbt --address-prefix="$ONOMY_ADDRESS_PREFIX" orchestrator \
             --cosmos-phrase="$ONOMY_ORCHESTRATOR_MNEMONIC" \
             --ethereum-key="$ETH_ORCHESTRATOR_PRIVATE_KEY" \
             --cosmos-grpc="http://$ONOMY_HOST:$ONOMY_GRPC_PORT/" \
             --ethereum-rpc="http://$ETH_HOST:$ETH_RPC_PORT/" \
             --fees="1$STAKE_DENOM" \
             --gravity-contract-address="$ETH_CONTRACT_ADDRESS"