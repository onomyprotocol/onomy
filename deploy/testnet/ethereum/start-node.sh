#!/bin/bash
set -eu

echo "Starting ethereum full node"

# Name of the onomy artifact
ONOMY_HOME="$HOME/.onomy"
mkdir -p $ONOMY_HOME/logs
GETH_LOG_FILE=$ONOMY_HOME/logs/geth.log

if [ "$(ulimit -n)" -lt 65535 ]; then
    echo "Fail ulimit: $(ulimit -n) < 65535"
    exit 1
fi

geth --rinkeby --http --http.addr 0.0.0.0 --http.api eth,net,web3,personal,txpool,admin --syncmode full \
 --pprof --pprof.addr 0.0.0.0 --pprof.port 6060 --metrics &>> $GETH_LOG_FILE &

echo "Ethereum rinkeby is started, check the logs file $GETH_LOG_FILE"

GETH_SYNCING_COMMAND="geth --rinkeby attach --exec \"eth.syncing\""
echo "You can check the sync status by command: $GETH_SYNCING_COMMAND"