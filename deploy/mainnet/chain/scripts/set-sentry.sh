#!/bin/bash
set -eu

echo "Updating the validator to work with sentries"

# Onomy home dir
ONOMY_HOME=$HOME/.onomy
ONOMY_DOUBLE_SIGN_CHECK_HEIGHT=10

# Config directories for onomy node
ONOMY_HOME_CONFIG="$ONOMY_HOME/config"
# Config file for onomy node
ONOMY_NODE_CONFIG="$ONOMY_HOME_CONFIG/config.toml"

ONOMY_SENTRY_IPS=
while [[ $ONOMY_SENTRY_IPS = "" ]]; do
   read -r -p "Enter sentry peers ips ip,ip2 :" ONOMY_SENTRY_IPS
done

ONOMY_SENTRY_IDS=
ONOMY_SENTRY_PEERS=
for sentryIP in ${ONOMY_SENTRY_IPS//,/ } ; do
  wget $sentryIP:26657/status? -O $ONOMY_HOME/sentry_status.json
  sentryID=$(jq -r .result.node_info.id $ONOMY_HOME/sentry_status.json)

  if [[ -z "${sentryID}" ]]; then
    echo "Something went wrong, can't fetch $sentryIP info: "
    cat $ONOMY_HOME/sentry_status.json
    exit 1
  fi

  rm $ONOMY_HOME/sentry_status.json

  ONOMY_SENTRY_IDS="$ONOMY_SENTRY_IDS$sentryID,"
  ONOMY_SENTRY_PEERS="$ONOMY_SENTRY_PEERS$sentryID@$sentryIP:26656,"
done

crudini --set $ONOMY_NODE_CONFIG p2p pex false
# list of sentry nodes
crudini --set $ONOMY_NODE_CONFIG p2p persistent_peers "\"$ONOMY_SENTRY_PEERS\""
# private_peer_ids	validator node ID
crudini --set $ONOMY_NODE_CONFIG p2p private_peer_ids "\"$ONOMY_SENTRY_IDS"\"
# unconditional_peer_ids	validator node ID, optionally sentry node IDs
crudini --set $ONOMY_NODE_CONFIG p2p unconditional_peer_ids "\"$ONOMY_SENTRY_IDS"\"
crudini --set $ONOMY_NODE_CONFIG p2p addr_book_strict false

crudini --set $ONOMY_NODE_CONFIG consensus double_sign_check_height $ONOMY_DOUBLE_SIGN_CHECK_HEIGHT

echo "The master is updated for sentries"