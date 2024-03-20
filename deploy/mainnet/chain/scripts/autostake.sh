#!/bin/bash

# Adapted from https://autostake.com/networks/onomy/#services

ONOMY_HOME=~/.onomy

cd $ONOMY_HOME
cp data/priv_validator_state.json priv_validator_state.json.backup
rm -r data/*

# Download the snapshot
SNAP_URL="http://snapshots.autostake.com/lyIs25DaSWMSm8evWKHGQrb/onomy-mainnet-1/latest.tar.lz4"
curl $SNAP_URL | lz4 -d | tar -xvf -

# AutoStake addrbook
cd $ONOMY_HOME/config
curl http://snapshots.autostake.com/lyIs25DaSWMSm8evWKHGQrb/onomy-mainnet-1/addrbook.json --output addrbook.json

# AutoStake State-Sync RPC
STATESYNC_RPC="https://onomy-mainnet-rpc.autostake.com:443"
LATEST_HEIGHT=$(curl -s "$STATESYNC_RPC/block" | jq -r .result.block.header.height)
TRUST_HEIGHT=$(("$LATEST_HEIGHT" - 2000))
TRUST_HASH=$(curl -s "$STATESYNC_RPC/block?height=$TRUST_HEIGHT" | jq -r .result.block_id.hash)

cd $ONOMY_HOME
sed -i '/\[statesync\]/{:a;n;/enable/s/false/true/;Ta;}' config/config.toml
sed -i "s@rpc_servers = \".*\"@rpc_servers = \"$STATESYNC_RPC,$STATESYNC_RPC\"@" config/config.toml
sed -i "s@trust_height = .*@trust_height = $TRUST_HEIGHT@" config/config.toml
sed -i "s@trust_hash = \".*\"@trust_hash = \"$TRUST_HASH\"@" config/config.toml

# AutoStake Seed
cd $ONOMY_HOME
SEED="ebc272824924ea1a27ea3183dd0b9ba713494f83@onomy-mainnet-seed.autostake.com:27556"
sed -i "s/seeds = \"\"/seeds = \"$SEED\"/" "config/config.toml"

# AutoStake Peer
cd $ONOMY_HOME
PEER="ebc272824924ea1a27ea3183dd0b9ba713494f83@onomy-mainnet-peer.autostake.com:27556"
sed -i "s/persistent_peers = \"\"/persistent_peers = \"$PEER\"/" "config/config.toml"

