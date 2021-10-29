#!/bin/bash
set -eu

echo "running onomy with bridge"

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
ONOMY_HOME="$HOME/.onomy"
# Onomy mnemonic used for orchestrator signing of the transactions (orchestrator_key.json file)
ONOMY_ORCHESTRATOR_MNEMONIC=$(jq -r .mnemonic $ONOMY_HOME/orchestrator_key.json)
# Onomy mnemonic used for fauset from (validator_key.json file)
ONOMY_VALIDATOR_MNEMONIC=$(jq -r .mnemonic $ONOMY_HOME/validator_key.json)

# Onomy chain demons
STAKE_DENOM="anom"

# The ETH key used for orchestrator signing of the transactions
ETH_ORCHESTRATOR_PRIVATE_KEY=c40f62e75a11789dbaf6ba82233ce8a52c20efb434281ae6977bb0b3a69bf709

ETH_GRAVITY_CONTRACT_ADDRESS=0x0F23c3f0C77582a5dB7fB3D61097B619982fb32f

# ------------------ Run onomy ------------------

echo "Starting $ONOMY_NODE_NAME"
$ONOMY start --pruning=nothing &

echo "Waiting $ONOMY_NODE_NAME to launch gRPC $ONOMY_GRPC_PORT..."

while ! timeout 1 bash -c "</dev/tcp/$ONOMY_HOST/$ONOMY_GRPC_PORT"; do
  sleep 1
done

echo "$ONOMY_NODE_NAME launched"

# ------------------ Run fauset ------------------

echo "Starting fauset based on validator account"
faucet -cli-name=$ONOMY -mnemonic="$ONOMY_VALIDATOR_MNEMONIC" &

#-------------------- Run orchestrator --------------------

echo "Starting orchestrator"

gbt --address-prefix="$ONOMY_ADDRESS_PREFIX" orchestrator \
             --cosmos-phrase="$ONOMY_ORCHESTRATOR_MNEMONIC" \
             --ethereum-key="$ETH_ORCHESTRATOR_PRIVATE_KEY" \
             --cosmos-grpc="http://$ONOMY_HOST:$ONOMY_GRPC_PORT/" \
             --ethereum-rpc="$ETH_RPC_ADDRESS" \
             --fees="1$STAKE_DENOM" \
             --gravity-contract-address="$ETH_GRAVITY_CONTRACT_ADDRESS"