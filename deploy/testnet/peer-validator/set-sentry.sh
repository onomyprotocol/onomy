#!/bin/bash
set -eu

echo "Updating the validator to work with sentries"

# Initial dir
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
    exit
  fi

  rm $ONOMY_HOME/sentry_status.json

  ONOMY_SENTRY_IDS="$ONOMY_SENTRY_IDS$sentryID,"
  ONOMY_SENTRY_PEERS="$ONOMY_SENTRY_PEERS$sentryID@$sentryIP:26656,"
done

# Switch sed command in the case of linux
fsed() {
  if [ `uname` = 'Linux' ]; then
    sed -i "$@"
  else
    sed -i '' "$@"
  fi
}

fsed 's#pex = true""#pex = false#g' $ONOMY_NODE_CONFIG
# list of sentry nodes
fsed 's#persistent_peers = ""#persistent_peers = "'$ONOMY_SENTRY_PEERS'"#g' $ONOMY_NODE_CONFIG
# private_peer_ids	validator node ID
fsed 's#private_peer_ids = ""#private_peer_ids = "'$ONOMY_SENTRY_IDS'"#g' $ONOMY_NODE_CONFIG
# unconditional_peer_ids	validator node ID, optionally sentry node IDs
fsed 's#unconditional_peer_ids = ""#unconditional_peer_ids = "'$ONOMY_SENTRY_IDS'"#g' $ONOMY_NODE_CONFIG
fsed 's#double_sign_check_height = 0#double_sign_check_height = '$ONOMY_DOUBLE_SIGN_CHECK_HEIGHT'#g' $ONOMY_NODE_CONFIG
# addr-book-strict	false // already done by init script

echo "The master is updated for sentries"