# Steps to run the master node

## Install dependencies from source code

```
./bin.sh
```

## Init node and genesys

Before run the script please set evn variable:

* ETH_ORCHESTRATOR_VALIDATOR_ADDRESS - this is the Ethereum public address with which the orchestrator will be running.

```
./init-master.sh
```

## Start the node

Before run the script please set up "ulimit > 65536":

```
./start-node.sh
```

## Deploy gravity contract

Before run the script please set evn variable:

* ETH_RPC_ADDRESS - the RPC address of the Ethereum node
* ETH_PRIVATE_KEY - the Ethereum private key which deploys the contract

```
./deploy-gravity.sh
```

## Run orchestrator

* ETH_RPC_ADDRESS - the RPC address of the Ethereum node
* ETH_ORCHESTRATOR_VALIDATOR_PRIVATE_KEY - the Ethereum private key which will be use for the orchestrator

```
./start-orchestrator.sh
```

# Run inside the container

```
docker run --name onomy-testnet-master -v `pwd`/master-validator:/root/master-validator -w /root/master-validator -it fedora:36
```