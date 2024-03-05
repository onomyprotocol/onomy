#!/bin/bash
set -eu

echo "Starting node exporter"

export PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/games:/usr/local/games:/snap/bin:/home/ubuntu/.onomy/bin:/home/ubuntu/.onomy/cosmovisor/genesis/bin

ONOMY_NODE_EXPORTER_PORT=9100

if [ "$(ulimit -n)" -lt 65535 ]; then
    echo "Fail ulimit: $(ulimit -n) < 65535"
    exit 1
fi

node_exporter --web.listen-address=:$ONOMY_NODE_EXPORTER_PORT