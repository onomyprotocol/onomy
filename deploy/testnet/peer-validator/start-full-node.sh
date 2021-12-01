#!/bin/bash
set -eu

# Name of the onomy artifact
ONOMY=onomyd
# The name of the Onomy node
ONOMY_NODE_NAME=$(jq -r .node_name $ONOMY_HOME/node_info.json)

echo "Starting onomy full node"

$ONOMY start > $HOME/.onomy/$ONOMY_NODE_NAME.log