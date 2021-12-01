# Steps to run the master node

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

## Start the node
Before run the script please set up "ulimit > 65536":

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

Also you cat check all current validators now.

```
onomyd q staking validators
```

## Send som tokens from you validator to your orchestrator

```
onomyd tx bank send {validator-address} {orchestrator-address} 5000000000000000000anom --chain-id=onomy-testnet --keyring-backend test
```
Check the orchestrator balance now
```
onomyd q bank balances {orchestrator-address}
```

## Init gbt
Before run the script please set evn variable:

* $ETH_ORCHESTRATOR_VALIDATOR_PRIVATE_KEY - the Ethereum private key which will be use for the orchestrator

```
./init-gbt.sh
```

And then check that your Ethereum address is in the list of curren calset
```
onomyd q gravity current-valset
```

## Run orchestrator
Before run the script please set evn variable:

* ETH_RPC_ADDRESS - the RPC address of the Ethereum node
* ETH_ORCHESTRATOR_VALIDATOR_PRIVATE_KEY - the Ethereum private key which will be use for the orchestrator

```
./start-orchestrator.sh
```

# Run inside the container with master validator

## Set up the master node from the instruction in the master-validator.

## Commit master node
```
docker commit onomy-testnet-master onomy-testnet-master-working
```
## Create network to conntct container
```
docker network create --driver bridge testnet
```
## Run the master (from the testnet folder)
```
docker run -dit --name onomy-testnet-master-working -v `pwd`/master-validator:/root/master-validator -w /root/master-validator --network testnet onomy-testnet-master-working sleep 10000000000
```
## Login to onomy-testnet-master-working
```
docker exec -it onomy-testnet-master-working bash
```
## Run the validator (from the testnet folder)
```
docker run -dit --name onomy-testnet-peer-working -v `pwd`/peer-validator:/root/peer-validator -w /root/peer-validator --network testnet fedora:35 sleep 10000000000
```
## Login to onomy-testnet-peer-working
```
docker exec -it onomy-testnet-peer-working bash
```
## Ping master node
```
yum install iputils
```
Ping and capture the output
```
ping onomy-testnet-master-working
```

Then the node is ready to start the main steps.
