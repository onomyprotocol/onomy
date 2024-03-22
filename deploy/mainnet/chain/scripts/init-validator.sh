#!/bin/bash
ONOMY_HOME=$HOME/.onomy
CHAINID="onomy-mainnet-1"
ONOMY_HOME_CONFIG="$ONOMY_HOME/config"
ONOMY_APP_CONFIG="$ONOMY_HOME_CONFIG/config.toml"
ONOMY_NODE_CONFIG="$ONOMY_HOME_CONFIG/config.toml"
ONOMY_CHAINID_FLAG="--chain-id $CHAINID"
ONOMY_SENTRY_DEFAULT_IPS="a.sentry.mainnet.onomy.io,b.sentry.mainnet.onomy.io,c.sentry.mainnet.onomy.io"

read -r -p "Enter a name for your node [onomy-sentry]:" ONOMY_NODE_NAME
ONOMY_NODE_NAME=${ONOMY_NODE_NAME:-onomy-sentry}

read -r -p "Enter sentry ips [$ONOMY_SENTRY_DEFAULT_IPS]:" ONOMY_SENTRY_IPS
ONOMY_SENTRY_IPS=${ONOMY_SENTRY_IPS:-$ONOMY_SENTRY_DEFAULT_IPS}

ONOMY_SENTRYS=
ONOMY_SENTRY_RPCS=
ONOMY_SENTRY_IDS=
for sentryIP in ${ONOMY_SENTRY_IPS//,/ } ; do
  wget $sentryIP:26657/status? -O $ONOMY_HOME/sentry_status.json
  sentryID=$(jq -r .result.node_info.id $ONOMY_HOME/sentry_status.json)

  if [[ -z "${sentryID}" ]]; then
    echo "Something went wrong, can't fetch $sentryIP info: "
    cat $ONOMY_HOME/sentry_status.json
    exit 1
  fi

  rm $ONOMY_HOME/sentry_status.json

  ONOMY_SENTRYS="$ONOMY_SENTRYS$sentryID@$sentryIP:26656,"
  ONOMY_SENTRY_RPCS="${ONOMY_SENTRY_RPCS}http://$sentryIP:26657,"
  ONOMY_SENTRY_IDS="$ONOMY_SENTRY_IDS$sentryID,"
done

echo "Creating $ONOMY_NODE_NAME node with chain-id=$CHAINID..."
echo "Initializing chain"
onomyd $ONOMY_CHAINID_FLAG init $ONOMY_NODE_NAME

#copy genesis file
cp -r ../genesis/genesis-mainnet-1.json $ONOMY_HOME_CONFIG/genesis.json

echo "Updating node config"

# change config

crudini --set $ONOMY_NODE_CONFIG p2p addr_book_strict true
crudini --set $ONOMY_NODE_CONFIG p2p external_address "\"\""
#crudini --set $ONOMY_NODE_CONFIG p2p seeds "\"$ONOMY_SENTRYS\""
crudini --set $ONOMY_NODE_CONFIG p2p seeds "\"\""
crudini --set $ONOMY_NODE_CONFIG p2p persistent_peers "\"$ONOMY_SENTRYS\""
crudini --set $ONOMY_NODE_CONFIG p2p unconditional_peer_ids "\"$ONOMY_SENTRY_IDS\""
crudini --set $ONOMY_NODE_CONFIG p2p pex false
crudini --set $ONOMY_NODE_CONFIG rpc laddr "\"tcp://127.0.0.1:26657\""
crudini --set $ONOMY_NODE_CONFIG statesync enable true
crudini --set $ONOMY_NODE_CONFIG statesync rpc_servers "\"$ONOMY_SENTRY_RPCS\""
crudini --set $ONOMY_NODE_CONFIG statesync trust_height 6953812
crudini --set $ONOMY_NODE_CONFIG statesync trust_hash "\"DE3ED360800652D73A4BB15695A494F573FEE392E1C246936D40046937DA98D8\""


crudini --set $ONOMY_APP_CONFIG api enable false
crudini --set $ONOMY_APP_CONFIG rosetta enable false
crudini --set $ONOMY_APP_CONFIG grpc enable false
crudini --set $ONOMY_APP_CONFIG grpc-web enable false
# create home directory
#
echo "The initial setup of $ONOMY_NODE_NAME is done"
