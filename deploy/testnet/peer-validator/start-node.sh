#!/bin/bash
set -eu

echo "Starting onomy full node"

# Name of the onomy artifact
ONOMY=onomyd
ONOMY_HOME="$HOME/.onomy"

mkdir -p $ONOMY_HOME/logs
ONOMY_LOG_FILE=$ONOMY_HOME/logs/onomyd.log

if [ "$(ulimit -n)" -lt 65535 ]; then
    echo "Fail ulimit: $(ulimit -n) < 65535"
    exit
fi

$ONOMY start &>> $ONOMY_LOG_FILE &

echo "Onomy is started, check the logs file $ONOMY_LOG_FILE"