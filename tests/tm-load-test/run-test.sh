#!/bin/bash
set -eux

CONNECTIONS=$1
TIME=$2
RATE=$3
SIZE=$4

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
result=$( docker images -q cosmoschain )

bash $DIR/buid-container.sh
# Remove existing container instance
set +e
docker rm -f cosmos_test_instance
set -e

pushd $DIR/../

# Run new test container instance
docker run -i --name cosmos_test_instance cosmoschain /bin/bash -c "bash start-node.sh && bash master-slave.sh $CONNECTIONS $TIME $RATE $SIZE"