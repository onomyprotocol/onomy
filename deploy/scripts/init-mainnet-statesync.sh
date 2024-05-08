#!/bin/bash
set -uxe

git clone https://github.com/onomyprotocol/onomy && cd onomy && git checkout v1.1.4
make install

export GOPATH=~/go
export PATH=$PATH:~/go/bin

onomyd init test --chain-id onomy-mainnet-1
wget -O ~/.onomy/config/genesis.json https://raw.githubusercontent.com/onomyprotocol/onomy/main/genesis/mainnet/genesis-mainnet-1.json

LATEST_HEIGHT=$(curl -s "https://rpc-mainnet.onomy.io:443/block" | jq -r .result.block.header.height);
BLOCK_HEIGHT=$(($LATEST_HEIGHT - 2000)) 
TRUST_HASH=$(curl -s "https://rpc-mainnet.onomy.io:443/block?height=$BLOCK_HEIGHT" | jq -r .result.block_id.hash)

export ONOMYD_STATESYNC_ENABLE=true
export ONOMYD_P2P_MAX_NUM_OUTBOUND_PEERS=100
export ONOMYD_P2P_MAX_NUM_INBOUND_PEERS=100
export ONOMYD_STATESYNC_RPC_SERVERS="https://rpc-mainnet.onomy.io:443,http://35.224.118.71:26657"
export ONOMYD_STATESYNC_TRUST_HEIGHT=$BLOCK_HEIGHT
export ONOMYD_STATESYNC_TRUST_HASH=$TRUST_HASH
export ONOMYDD_P2P_SEEDS="211535f9b799bcc8d46023fa180f3359afd4c1d3@44.213.44.5:26656,00ce2f84f6b91639a7cedb2239e38ffddf9e36de@44.195.221.88:26656,cd9a47cebe8eef076a5795e1b8460a8e0b2384e5@3.210.0.126:26656,60194df601164a8b5852087d442038e392bf7470@180.131.222.74:26656,0dbe561f30862f386456734f12f431e534a3139c@34.133.228.142:26656,4737740b63d6ba9ebe93e8cc6c0e9197c426e9f4@195.189.96.106:52756"

onomyd start --x-crisis-skip-assert-invariants --p2p.persistent_peers 211535f9b799bcc8d46023fa180f3359afd4c1d3@44.213.44.5:26656,00ce2f84f6b91639a7cedb2239e38ffddf9e36de@44.195.221.88:26656,cd9a47cebe8eef076a5795e1b8460a8e0b2384e5@3.210.0.126:26656,60194df601164a8b5852087d442038e392bf7470@180.131.222.74:26656,0dbe561f30862f386456734f12f431e534a3139c@34.133.228.142:26656,4737740b63d6ba9ebe93e8cc6c0e9197c426e9f4@195.189.96.106:52756,00ce2f84f6b91639a7cedb2239e38ffddf9e36de@44.195.221.88:26656
