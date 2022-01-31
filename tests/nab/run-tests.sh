#!/bin/bash
TEST_TYPE=$1
set -eu

echo "Clearing containers"
docker-compose down

bash $DIR/run-eth.sh

echo "Starting test"
docker-compose run nab /bin/bash /onomy/tests/nab/container-scripts/all-up-test-internal.sh 1 $TEST_TYPE

echo "Clearing containers"
docker-compose down