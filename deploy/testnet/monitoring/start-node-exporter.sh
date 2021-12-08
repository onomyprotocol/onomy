#!/bin/bash
set -eu

echo "Starting node_exporter"

ONOMY_HOME=$HOME/.onomy
ONOMY_NODE_EXPORTER_PORT=9100

mkdir -p $ONOMY_HOME/logs
ONOMY_NODE_EXPORTER_LOG_FILE=$ONOMY_HOME/logs/node_exporte.log

if [ "$(ulimit -n)" -lt 65535 ]; then
    echo "Fail ulimit: $(ulimit -n) < 65535"
    exit 1
fi

node_exporter --web.listen-address=:$ONOMY_NODE_EXPORTER_PORT &>> $ONOMY_NODE_EXPORTER_LOG_FILE &

echo "Node_exporter is started, check the logs file $ONOMY_NODE_EXPORTER_LOG_FILE"