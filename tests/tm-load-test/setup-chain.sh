#!/bin/bash
set -eux

CONNECTIONS=$1
TIME=$2
RATE=$3
SIZE=$4
FILE_NAME=$5 || "log.txt"

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
result=$( docker images -q cosmoschain )

sudo bash $DIR/buid-container.sh	
# Remove existing container instance
set +e
docker rm -f cosmos_test_instance
set -e

pushd $DIR/../


# Run new test container instance
docker run -d --name cosmos_test_instance -v "$(pwd)"/target:/root/logs -p 9090:9090 -p 26657:26657 -p 1317:1317 -it cosmoschain /bin/bash -c "sleep infinity | bash /root/home/master-validator/init.sh | bash /root/home/run-load-test/master-slave.sh $CONNECTIONS $TIME $RATE $SIZE $FILE_NAME"
sleep $TIME