#!/bin/bash
set -eux

# the directory of this script, useful for allowing this script
# to be run with any PWD
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# Remove existing container instance
set +e
docker rm -f onomy_test_instance
set -e

NODES=3

pushd $DIR/../

# Run new test container instance
docker run --name onomy_test_instance --mount type=bind,source="$(pwd)"/,target=/onomy --cap-add=NET_ADMIN -p 9090:9090 -p 26657:26657 -p 1317:1317 -p 8545:8545 -it onomy-base /bin/bash /onomy/tests/gravity/container-scripts/reload-code.sh $NODES