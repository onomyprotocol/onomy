# Steps to run the peer-validator node

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

## Optionally with seeds

### Start sentry nodes based on instructions from the [sentry](../sentry/readme.md)

### Run script to set up the private connection of the validator and sentries

You will need to provide the sentry IPs.

```
./set-sentry.sh
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

## Init validator

Init the validator account to be deposited

```
./init-validator.sh
```

## Get tokens from master account

Request tokens (for validator) from the master account by "text" request to onomy.

### To do it from the master node call:

```
onomyd tx bank send account1 {validator-address} 20000000000000000000anom --chain-id=onomy-testnet --keyring-backend test
```

### Then check the balance on validator node

```
onomyd q bank balances {validator-address}
```

If the "amount" is updated you are ready to become a validator

## Create a new onomy validator

```
./create-validator.sh
```

Also you can check all current validators now.

```
onomyd q staking validators
```

## Send some tokens from you validator to your orchestrator

```
onomyd tx bank send {validator-address} {orchestrator-address} 5000000000000000000anom --chain-id=onomy-testnet --keyring-backend test
```

Check the orchestrator balance now

```
onomyd q bank balances {orchestrator-address}
```

## Init gbt

Before run the script please set env variable:

* ETH_ORCHESTRATOR_VALIDATOR_PRIVATE_KEY - the Ethereum private key which will be use for the orchestrator
* ETH_GRAVITY_CONTRACT_ADDRESS - gravity contract address (will be a constant later)

```
./init-gbt.sh
```

And then check that your Ethereum address is in the list of curren calset

```
onomyd q gravity current-valset
```

## Run orchestrator

Before run the script please set env variable:

* ETH_RPC_ADDRESS - the RPC address of the Ethereum node

```
./start-orchestrator.sh
```

## Setup auto-start

Add env ETH_RPC_ADDRESS, start-node.sh and start-orchestrator.sh
to your crontab or /etc/init.d in order to start automatically after the OS restart.

If you used the bin.sh installation then additionally you need to add
```
export PATH=$PATH:$ONOMY_HOME/bin
```

In your start scripts (after the ONOMY_HOME initialization)
