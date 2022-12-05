#!/bin/bash
set -eu

echo "Enabling cors"

# Onomy home dir
ONOMY_HOME=$HOME/.onomy

# Config directories for onomy node
ONOMY_HOME_CONFIG="$ONOMY_HOME/config"
# Config file for onomy node
ONOMY_NODE_CONFIG="$ONOMY_HOME_CONFIG/config.toml"
# App config file for onomy node
ONOMY_APP_CONFIG="$ONOMY_HOME_CONFIG/app.toml"

crudini --set $ONOMY_NODE_CONFIG rpc cors_allowed_origins "[\"*\"]"
crudini --set $ONOMY_APP_CONFIG api enabled-unsafe-cors true