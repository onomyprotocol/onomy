# Steps to run the full node

## Install dependencies from source code

```
./bin.sh
```

## Get the seed

You can use default or get in from the master node

## Init full node

```
./init-full-node.sh
```

Get the node id:

```
onomyd tendermint show-node-id
```

Get the node ip:

```
hostname -I | awk '{print $1}'
```

## Optionally expose monitoring

```
./expose-metrics.sh
```

This script will enable the prometheus metrics in your node config.

## Start the node

Before run the script please set up "ulimit > 65535":

```
./start-node.sh
```

## Setup auto-start

Add env ETH_RPC_ADDRESS, start-node.sh and start-orchestrator.sh
to your ~/.bashrc in order to start automatically after the OS restart.