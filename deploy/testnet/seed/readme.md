# Steps to run the seed node

## Install dependencies from source code

```
./bin.sh
```

## Init seed node

```
./init-seed-node.sh
```

## Copy the genesys from the master node to the seed node

Path to the genesis is: /root/.onomy/config/genesis.json

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

Get the node id:

```
onomyd tendermint show-node-id
```

Get the node ip:

```
hostname -I | awk '{print $1}'
```

The seed peer is: id@$ip:26656


## Setup auto-start

Add start-node.sh to your crontab or /etc/init.d in order to start automatically after the OS restart.

If you used the bin.sh installation then additionally you need to add
```
export PATH=$PATH:$ONOMY_HOME/bin
```

In your start scripts (after the ONOMY_HOME initialization)
