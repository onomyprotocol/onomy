#!/bin/bash
set -eu

echo "Starting node exporter"

ONOMY_NODE_EXPORTER_PORT=9100

if [ "$(ulimit -n)" -lt 65535 ]; then
    echo "Fail ulimit: $(ulimit -n) < 65535"
    exit 1
fi

node_exporter --web.listen-address=:$ONOMY_NODE_EXPORTER_PORT