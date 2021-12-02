
# Build testnet-node
```
DOCKER_BUILDKIT=0 docker build -t onomy/testnet-node:local .
```

# Run inside the container with master validator

## Set up the master node from the instruction in the master-validator.

## Create network to connect container

```
docker network create --driver bridge testnet
```

## Run the master (from the testnet folder)

```
docker run -dit --name onomy-testnet-master -v `pwd`:/root/testnet --network testnet onomy/testnet-node:local sleep 100000000
```

## Login to onomy-testnet-master-working

```
docker exec -it onomy-testnet-master bash
```

## Run the validator

```
docker run -dit --name onomy-testnet-validator -v `pwd`:/root/testnet --network testnet onomy/testnet-node:local sleep 100000000
```

## Login to onomy-testnet-validator

```
docker exec -it onomy-testnet-validator bash
```

## Run the seed1

```
docker run -dit --name onomy-testnet-seed1 -v `pwd`:/root/testnet --network testnet onomy/testnet-node:local sleep 100000000
```

## Login to onomy-testnet-seed1

```
docker exec -it onomy-testnet-seed1 bash
```


## Run the sentry1

```
docker run -dit --name onomy-testnet-sentry1 -v `pwd`:/root/testnet --network testnet onomy/testnet-node:local sleep 100000000
```

## Login to onomy-testnet-sentry1

```
docker exec -it onomy-testnet-sentry1 bash
```

## Run the sentry2

```
docker run -dit --name onomy-testnet-sentry2 -v `pwd`:/root/testnet --network testnet onomy/testnet-node:local sleep 100000000
```

## Login to onomy-testnet-sentry2

```
docker exec -it onomy-testnet-sentry2 bash
```

## Run the sentry3

```
docker run -dit --name onomy-testnet-sentry3 -v `pwd`:/root/testnet --network testnet onomy/testnet-node:local sleep 100000000
```

## Login to onomy-testnet-sentry3

```
docker exec -it onomy-testnet-sentry3 bash
```

Now you are ready to set up any config

# Utils 

## Ping master node

```
yum install iputils
```

Ping and capture the output

```
ping onomy-testnet-master
```
