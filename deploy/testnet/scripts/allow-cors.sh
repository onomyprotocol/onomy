#!/bin/bash
set -eu

echo "Enabling cors"

# Initial dir
ONOMY_HOME=$HOME/.onomy

# Config directories for onomy node
ONOMY_HOME_CONFIG="$ONOMY_HOME/config"
# Config file for onomy node
ONOMY_NODE_CONFIG="$ONOMY_HOME_CONFIG/config.toml"
# App config file for onomy node
ONOMY_APP_CONFIG="$ONOMY_HOME_CONFIG/app.toml"

# Switch sed command in the case of linux
fsed() {
  if [ `uname` = 'Linux' ]; then
    sed -i "$@"
  else
    sed -i '' "$@"
  fi
}

fsed 's#cors_allowed_origins = \[\]#cors_allowed_origins = \[\"*\"\]#g' $ONOMY_NODE_CONFIG
fsed 's#enabled-unsafe-cors = false#enabled-unsafe-cors = true#g' $ONOMY_APP_CONFIG