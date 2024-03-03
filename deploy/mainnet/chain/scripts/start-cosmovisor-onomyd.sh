#!/bin/bash
set -eu

echo "Starting onomy node"

export PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/games:/usr/local/games:/snap/bin:/home/ubuntu/.onomy/bin:/home/ubuntu/.onomy/cosmovisor/genesis/bin
export DAEMON_HOME=/home/ubuntu/.onomy/
export DAEMON_NAME=onomyd
export DAEMON_RESTART_AFTER_UPGRADE=true

if [ "$(ulimit -n)" -lt 65535 ]; then
    echo "Fail ulimit: $(ulimit -n) < 65535"
    exit 1
fi

cosmovisor start