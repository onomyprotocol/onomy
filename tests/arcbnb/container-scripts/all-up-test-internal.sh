#!/bin/bash
# the script run inside the container for all-up-test.sh
NODES=$1
TEST_TYPE=$2
ALCHEMY_ID=$3
set -eux

bash /onomy/tests/arcbnb/container-scripts/setup-validators.sh $NODES
bash /onomy/tests/arcbnb/container-scripts/run-testnet.sh $NODES $TEST_TYPE $ALCHEMY_ID &

# deploy the ethereum contracts
DEPLOY_CONTRACTS=1 RUST_BACKTRACE=full CHAIN_BINARY=onomyd ADDRESS_PREFIX=onomy RUST_LOG=INFO test-runner
bash /onomy/tests/arcbnb/container-scripts/integration-tests.sh $NODES $TEST_TYPE