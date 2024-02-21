#!/bin/bash
set -eu

echo "Starting geth"

if [ "$(ulimit -n)" -lt 65535 ]; then
    echo "Fail ulimit: $(ulimit -n) < 65535"
    exit 1
fi

GETH_SYNCING_COMMAND="geth --goerli attach --exec \"eth.syncing\""
echo "You can check the sync status by command: $GETH_SYNCING_COMMAND"

geth --goerli --http --http.addr 0.0.0.0 --http.api eth,net,web3,personal,txpool,admin \
 --ws --ws.origins '*' --ws.addr 0.0.0.0 --ws.api eth,net,web3,personal,txpool,admin  --syncmode full \
 --pprof --pprof.addr 0.0.0.0 --pprof.port 6060 --metrics --http.vhosts '*'