#!/bin/bash
set -eu

echo "Increasing network configuration"

# Onomy home dir
ONOMY_HOME=$HOME/.onomy
# Config directories for onomy node
ONOMY_HOME_CONFIG="$ONOMY_HOME/config"
# Config file for onomy node
ONOMY_NODE_CONFIG="$ONOMY_HOME_CONFIG/config.toml"

# rpc config
crudini --set $ONOMY_NODE_CONFIG rpc grpc_max_open_connections 10000
crudini --set $ONOMY_NODE_CONFIG rpc max_open_connections 10000
crudini --set $ONOMY_NODE_CONFIG rpc max_subscription_clients 10000
crudini --set $ONOMY_NODE_CONFIG rpc max_subscriptions_per_client 20

echo "The configuration is updated"