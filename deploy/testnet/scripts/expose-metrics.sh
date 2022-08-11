#!/bin/bash
set -eu

echo "Enable the node's monitoring"

# Onomy home dir
ONOMY_HOME=$HOME/.onomy
# Config directories for onomy node
ONOMY_HOME_CONFIG="$ONOMY_HOME/config"
# Config file for onomy node
ONOMY_NODE_CONFIG="$ONOMY_HOME_CONFIG/config.toml"

crudini --set $ONOMY_NODE_CONFIG instrumentation prometheus true

echo "The prometheus is enabled"