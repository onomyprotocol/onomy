#!/bin/bash

# Adapted from https://autostake.com/networks/onomy/#services

ONOMY_HOME=~/.onomy

cd $ONOMY_HOME
cp data/priv_validator_state.json priv_validator_state.json.backup
rm -r data/*

# Download the snapshot
SNAP_URL="http://snapshots.autostake.com/lyIs25DaSWMSm8evWKHGQrb/onomy-mainnet-1/latest.tar.lz4"
curl $SNAP_URL | lz4 -d | tar -xvf -

# Restore the priv_validator_state.json from backup
cp priv_validator_state.json.backup data/priv_validator_state.json

# AutoStake addrbook
cd $ONOMY_HOME/config/
curl http://snapshots.autostake.com/lyIs25DaSWMSm8evWKHGQrb/onomy-mainnet-1/addrbook.json --output $ONOMY_HOME/config/addrbook.json

# AutoStake Seed
cd $ONOMY_HOME
SEED="ebc272824924ea1a27ea3183dd0b9ba713494f83@onomy-mainnet-seed.autostake.com:27556"
sed -i "s/seeds = \"\"/seeds = \"$SEED\"/" "config/config.toml"

# AutoStake Peer
cd $ONOMY_HOME
PEER="ebc272824924ea1a27ea3183dd0b9ba713494f83@onomy-mainnet-peer.autostake.com:27556"
sed -i "s/persistent_peers = \"\"/persistent_peers = \"$PEER\"/" "config/config.toml"
