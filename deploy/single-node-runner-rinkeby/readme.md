# Single node runner rinkeby

Table of Contents
=================

* [Description](#Description)
* [Entrypoints](#Entrypoints)
* [Build](#Build)
* [Run](#Run)
* [Accounts](#Accounts)
* [Client](#Client)
* [Fauset](#Fauset)

## Description

The "onomy-single-node-runner-rinkeby" is a docker image that contains already prebuilt onomy home files, and a constant
Gravity.sol contract address deployed to the Rinkeby. The container runs the onomy test net with one validator and one
bridge orchestrator. This is a minimum that makes the bridge works.

## Entrypoints

0.0.0.0 should be change to a host of a deployed container

- onomy swagger: [http://0.0.0.0:1317/](http://0.0.0.0:1317/)
- onomy rpc: [http://0.0.0.0:1317/](http://0.0.0.0:1317/)
- onomy grpc: [http://0.0.0.0:9090/](http://0.0.0.0:9090/)
- onomy fauset (nom): [http://0.0.0.0:8080/](http://0.0.0.0:8080/)
- ethereum rpc: [http://0.0.0.0:8545/](http://0.0.0.0:8545/)
  | [geth API docs](https://geth.ethereum.org/docs/rpc/server)

## Build

### Build locally

  ```
  docker build -t onomy/onomy-single-node-runner-rinkeby:local  -f Dockerfile ../../
  ```

### Steps to rebuild image from scratch

- Prepare onomy (generate genesys and etc.)

Run init.sh inside the container with binaries and copy $ONOMY_HOME folder outside. This home contains testnet chain
settings and genesys file of the chain.

- Deploy Ethereum contact. The instruction is in cosmos-gravity-bridge repo.

***The contract used for the runner***.

```
Gravity deployed at Address - 0x39cea5A03Bd5266D8e6D9A942b258923df99926D
```

## Run

To run use docker-compose-remote.yml to run on the remote node, and docker-compose.yml for the local.

## Accounts

### Ethereum rinkeby accounts

- root validator (contract deployer)

```
name: test-chain-root
address: 0x97D5F5D4fDf83b9D2Cb342A09b8DF297167a73d0
private key: e0b21b1d80e53f38734a3ed395796956b50c637916ddbb6cedb096b848053d2d
```

- orchestrator/validator

```
name: test-chain-orchestrator-validator:  
address: 0x2d9480eBA3A001033a0B8c3Df26039FD3433D55d
private key: c40f62e75a11789dbaf6ba82233ce8a52c20efb434281ae6977bb0b3a69bf709
```

### onomy test chain accounts

- orchestrator

```
{
  "name": "orch",
  "type": "local",
  "address": "onomy1yf0xqjuhk44krpdgm3mrlemlvad3dr9xucetad",
  "pubkey": "onomypub1addwnpepqgmdjw54zcx5j0jk62mhwp6gp686rp08c3e6fegyth0m7gvqsm7y7vd4sw6",
  "mnemonic": "require pitch mansion frequent bean agent swing say trick remain sausage clever blind axis dove spell leg float wrist tackle rather million theory dolphin"
}
```

- validator

```
{
  "name": "val",
  "type": "local",
  "address": "onomy17vfz8e0ecvpj2emff5q96awcp78tq34tg8fjnm",
  "pubkey": "onomypub1addwnpepqfk93v68tqavpze4007zswrmv3vhp0mdwxfq4zz58tfk0ryx7h8yxsjsf3a",
  "mnemonic": "enough settle civil month timber wrap genre coconut stick neutral frozen load dutch venue install depth poet lyrics skull orchard rail trip worry robust"
}
```

## Client

The ETC20 coins used in that part are FAU. You can substitute "erc20-address" to any other token.

### Inside the single-runner-container (all tools already installed)

- Connect to the container and go to orchestrator/client folder

  ```
  docker exec -it onomy-single-node-runner-rinkeby bash
  ```

- Mint some FAU tokens for the 0x97D5F5D4fDf83b9D2Cb342A09b8DF297167a73d0 on the [page](https://erc20faucet.com/)

- Get cosmos user's balance:
  ```
  onomyd query bank balances onomy17vfz8e0ecvpj2emff5q96awcp78tq34tg8fjnm
  ```

- Send from eth to cosmos (from 0x97D5F5D4fDf83b9D2Cb342A09b8DF297167a73d0 to
  onomy17vfz8e0ecvpj2emff5q96awcp78tq34tg8fjnm)
  ```
  gbt -a onomy client eth-to-cosmos \
          --ethereum-key="c40f62e75a11789dbaf6ba82233ce8a52c20efb434281ae6977bb0b3a69bf709" \
          --ethereum-rpc="$ETH_RPC_ADDRESS" \
          --gravity-contract-address=0x39cea5A03Bd5266D8e6D9A942b258923df99926D \
          --token-contract-address=0xFab46E002BbF0b4509813474841E0716E6730136 \
          --amount=10 \
          --destination=onomy17vfz8e0ecvpj2emff5q96awcp78tq34tg8fjnm
  ```

- Now check users balances on both sides

- Send wnom from eth to cosmos (from 0x97D5F5D4fDf83b9D2Cb342A09b8DF297167a73d0 to
  onomy17vfz8e0ecvpj2emff5q96awcp78tq34tg8fjnm)
  ```
  gbt -a onomy client eth-to-cosmos \
        --ethereum-key="c40f62e75a11789dbaf6ba82233ce8a52c20efb434281ae6977bb0b3a69bf709" \
        --ethereum-rpc="$ETH_RPC_ADDRESS" \
        --gravity-contract-address=0x39cea5A03Bd5266D8e6D9A942b258923df99926D \
        --token-contract-address=0xe7c0fd1f0A3f600C1799CD8d335D31efBE90592C \
        --amount=1 \
        --destination=onomy17vfz8e0ecvpj2emff5q96awcp78tq34tg8fjnm
  ```
- Now check users balances on both sides the amon should be increased

- Send from cosmos to eth (from onomy17vfz8e0ecvpj2emff5q96awcp78tq34tg8fjnm to
  0x2d9480eBA3A001033a0B8c3Df26039FD3433D55d /different eth address)

  ```
  gbt -a onomy client cosmos-to-eth --cosmos-phrase="enough settle civil month timber wrap genre coconut stick neutral frozen load dutch venue install depth poet lyrics skull orchard rail trip worry robust" \
                  --cosmos-grpc="http://0.0.0.0:9090" \
                  --fees=1gravity0xFab46E002BbF0b4509813474841E0716E6730136 \
                  --amount=555555gravity0xFab46E002BbF0b4509813474841E0716E6730136 \
                  --eth-destination=0x2d9480eBA3A001033a0B8c3Df26039FD3433D55d

- Now check users balances on both sides one more time

### Using linux or docker/linux to connect

- Configure your environment Run docker with ubuntu and connect to the container

  ```
  docker run --name onomy-single-node-runner-rinkeby-client  -v `pwd`/client/linux:/root/home/client -w /root/home/client -it ubuntu 
  ```

- Add permissions

  ``` 
  chmod 777 ./gbt
  ```

- Install curl (optional if already installed)

  ```
  apt-get update && apt-get install curl -yq
  ```

Now you can use ./gbt to run the commands

# Fauset

[docs description](https://github.com/tendermint/faucet)

* install

```
curl https://get.starport.network/faucet! | bash
```

* run

```
faucet -cli-name=onomyd -keyring-backend=test -node=tcp://0.0.0.0:26657 -mnemonic="enough settle civil month timber wrap genre coconut stick neutral frozen load dutch venue install depth poet lyrics skull orchard rail trip worry robust"
```

* fauset API

[http://0.0.0.0:8000/](http://0.0.0.0:8000/)

* CURL example

Send 10nom to onomy1yf0xqjuhk44krpdgm3mrlemlvad3dr9xucetad

```
curl -X POST "http://0.0.0.0:8000/" -H  "accept: application/json" -H  "Content-Type: application/json" -d "{  \"address\": \"onomy1yf0xqjuhk44krpdgm3mrlemlvad3dr9xucetad\",  \"coins\": [    \"10nom\"  ]}"
```