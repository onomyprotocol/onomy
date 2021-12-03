#!/bin/bash
set -eu

echo "Starting onomy full node"

# Name of the onomy artifact
ONOMY_HOME="$HOME/.onomy"
mkdir -p $ONOMY_HOME/logs
GETH_LOG_FILE=$ONOMY_HOME/logs/geth.log

if [ "$(ulimit -n)" -lt 65536 ]; then
    echo "Fail ulimit: $(ulimit -n) < 65536"
    exit
fi

geth --rinkeby --http --http.addr 0.0.0.0 --http.api eth,net,web3,personal,txpool,admin --syncmode full &>> $GETH_LOG_FILE &

GETH_SYNCING_COMMAND="geth --rinkeby attach --exec \"eth.syncing\""
echo "You can check the sync status by command: $GETH_SYNCING_COMMAND"