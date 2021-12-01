#!/bin/bash
set -eu

echo "Starting onomy full node"

# Name of the onomy artifact
ONOMY=onomyd
ONOMY_LOG_FILE=$HOME/.onomy/onomyd.log

if [ "$(ulimit -n)" -lt 65536 ]; then
    echo "Fail: $(ulimit -n) < 65536"
    exit
fi

$ONOMY start &>> $ONOMY_LOG_FILE &

echo "Onomy is started, check the logs file $ONOMY_LOG_FILE"
