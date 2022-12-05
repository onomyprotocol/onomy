#!/bin/bash
set -eu

echo "Starting onomy node"

if [ "$(ulimit -n)" -lt 65535 ]; then
    echo "Fail ulimit: $(ulimit -n) < 65535"
    exit 1
fi

cosmovisor start
