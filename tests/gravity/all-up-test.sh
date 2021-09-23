#!/bin/bash
set -eux
# the directory of this script, useful for allowing this script
# to be run with any PWD
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

result=$( docker images -q onomy-base )

# builds the container containing various system deps
# also builds Peggy once in order to cache Go deps, this container
# must be rebuilt every time you run this test if you want a faster
# solution use start chains and then run tests
bash $DIR/build-container.sh

# Remove existing container instance
set +e
docker rm -f onomy_all_up_test_instance
set -e

NODES=3
set +u
TEST_TYPE=$1
set -u

echo "starting test $TEST_TYPE"
# Run new test container instance
docker run --name onomy_all_up_test_instance --cap-add=NET_ADMIN -t onomy-base /bin/bash /onomy/tests/gravity/container-scripts/all-up-test-internal.sh $NODES $TEST_TYPE
echo "test $TEST_TYPE completed successfully"