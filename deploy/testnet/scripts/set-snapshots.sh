#!/bin/bash
set -eu

echo "Setting snapshot configuration"

# Onomy home dir
ONOMY_HOME=$HOME/.onomy
# Config directories for onomy node
ONOMY_HOME_CONFIG="$ONOMY_HOME/config"
# App config file for onomy node
ONOMY_APP_CONFIG="$ONOMY_HOME_CONFIG/app.toml"


# snapshot config
crudini --set $ONOMY_APP_CONFIG state-sync snapshot-interval 1000
crudini --set $ONOMY_APP_CONFIG state-sync snapshot-keep-recent 3

echo "The snapshot configuration is updated"