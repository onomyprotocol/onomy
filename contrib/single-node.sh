#!/bin/sh

set -o errexit -o nounset

CHAINID=$1
GENACCT=$2

if [ -z "$1" ]; then
  echo "Need to input chain id..."
  exit 1
fi

if [ -z "$2" ]; then
  echo "Need to input genesis account address..."
  exit 1
fi

# Build genesis file incl account for passed address
coins="10000000000stake,100000000000samoleans"
onomyd init --chain-id $CHAINID $CHAINID
onomyd keys add validator --keyring-backend="test"
onomyd add-genesis-account $(onomyd keys show validator -a --keyring-backend="test") $coins
onomyd add-genesis-account $GENACCT $coins
onomyd gentx validator 5000000000stake --keyring-backend="test" --chain-id $CHAINID
onomyd collect-gentxs

# Set proper defaults and change ports
echo "Setting rpc listen address"
sed -i '' 's#"tcp://127.0.0.1:26657"#"tcp://0.0.0.0:26657"#g' ~/.onomy/config/config.toml
echo 2
sed -i '' 's/timeout_commit = "5s"/timeout_commit = "1s"/g' ~/.onomy/config/config.toml
sed -i '' 's/timeout_propose = "3s"/timeout_propose = "1s"/g' ~/.onomy/config/config.toml
sed -i '' 's/index_all_keys = false/index_all_keys = true/g' ~/.onomy/config/config.toml

# Start the onomy
onomyd start --pruning=nothing
