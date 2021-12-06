# Steps to run the sentry node

## Install dependencies from source code

```
./bin.sh
```

## Get the seed

You can use default or get in from the master node

## Init sentry node

```
./init-sentry-node.sh
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

Add start-node.sh to your ~/.bashrc in order to start automatically after the OS restart.