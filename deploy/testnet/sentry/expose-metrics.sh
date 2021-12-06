#!/bin/bash
set -eu

echo "Enable the node's monitoring"

# Initial dir
ONOMY_HOME=$HOME/.onomy
# Config directories for onomy node
ONOMY_HOME_CONFIG="$ONOMY_HOME/config"
# Config file for onomy node
ONOMY_NODE_CONFIG="$ONOMY_HOME_CONFIG/config.toml"

# Switch sed command in the case of linux
fsed() {
  if [ `uname` = 'Linux' ]; then
    sed -i "$@"
  else
    sed -i '' "$@"
  fi
}

fsed 's#prometheus = false#prometheus = true#g' $ONOMY_NODE_CONFIG

echo "The prometheus is enabled"