# Single node runner rinkeby

Table of Contents
=================
*  [Description](#Description)
*  [Entrypoints](#Entrypoints)
*  [Build](#Build)
*  [Run](#Run)
*  [Accounts](#Accounts)
*  [Client](#Client)
*  [Fauset](#Fauset)

## Description

The "onomy-single-node-runner-rinkeby" is a docker image that contains already prebuilt onomy home files, and a
constant Gravity.sol contract address deployed to the Rinkeby. The container runs the onomy test net with one validator
and one bridge orchestrator. This is a minimum that makes the bridge works.

## Entrypoints 

0.0.0.0 should be change to a host of a deployed container

- onomy swagger: [http://0.0.0.0:1317/](http://0.0.0.0:1317/)
- onomy rpc: [http://0.0.0.0:1317/](http://0.0.0.0:1317/)
- onomy grpc: [http://0.0.0.0:9090/](http://0.0.0.0:9090/)
- onomy fauset (nom): [http://0.0.0.0:8080/](http://0.0.0.0:8080/)
- ethereum rpc: [http://0.0.0.0:8545/](http://0.0.0.0:8545/) | [geth API docs](https://geth.ethereum.org/docs/rpc/server)

## Build

### Build locally
  
  ```
  docker build -t onomy/onomy-single-node-runner-rinkeby:local  -f Dockerfile ../../ --no-cache
  ```

### Steps to rebuild image from scratch

- Prepare onomy (generate genesys and etc.)

Run init.sh inside the container with binaries and copy $ONOMY_HOME folder outside.
This home contains testnet chain settings and genesys file of the chain.

- Deploy Ethereum contact. The instruction is in cosmos-gravity-bridge repo.

***The contract used for the runner***.

```
Gravity deployed at Address - 0x8A0814b7251138Dea19054425D0dfF0C497305d3
```

## Run

- Run with docker (image from the dockerhub)

  ```
  docker run --name onomy-single-node-runner-rinkeby \
              -p 26656:26656 -p 26657:26657 -p 1317:1317 -p 61278:61278 -p 9090:9090 -p 8000:8000 -p 8545:8545 \
              -v /mnt/volume_nyc1_03:/root/home/onomy/onomy/data/. \
              -it --restart on-failure onomy/onomy-single-node-runner-rinkeby:latest
  ```

  **latest** here is a tag of the runner. You can get the full list on the page [tags](https://hub.docker.com/repository/docker/onomy/onomy-single-node-runner-rinkeby/tags?page=1&ordering=last_updated)
  
  The docker command uses local "/mnt/volume_nyc1_03" directory to save onomy db files, onomy_home/data/priv_validator_state.json 
  file should be there before the first run of container. 
  
  Eth Rinkeby data files are inside the container and will be re-synchronised is case of the container replacement 
  (and restored from the Rinkeby network).
  
- Run with docker compose
  ```
  docker-compose down && docker-compose up
  ```
  This docker-compose uses local image for the run.

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
"name": "orch",
"type": "local",
"address": "onomy1c3arwy47pvrl5ap0gxjt2ujxc37kyhzle4ltjf",
"pubkey": "onomypub1addwnpepqwhnxwn9l2mezqy5q34dgjqgryuywmdu8l3fjvd8f80k6j5zk3m9qpey0pz",
"mnemonic": "edge zone faint love cherry under spell alone throw ladder common cable garden enroll uncle task lounge left blush sight unknown pencil clip chunk"
```

- validator
```
"name": "val",
"type": "local",
"address": "onomy1e6xtwjw9mgmljyrqw6mlw3nrpuz3p79gct73nw",
"pubkey": "onomypub1addwnpepq28frhwps0tz8xkhus0kep2d8smfndx830hy2p2zr2e4f7q3sj94shxlvpx",
"mnemonic": "staff pool flush bundle radar craft sister local fiction clown anger friend torch toss trial drift choice fine exist weapon hamster energy marine hub"
```

## Client

The ETC20 coins used in that part are FAU. You can substitute "erc20-address" to any other token.

  ### Inside the single-runner-container (all tools already installed)
  
  - Connect to the container and go to orchestrator/client folder
    
    ```
    docker exec -it onomy-single-node-runner-rinkeby bash
    ```
    
    ```
    cd /go/src/github.com/onomyprotocol/cosmos-gravity-bridge/orchestrator/target/x86_64-unknown-linux-musl/release
    ```
  
  - Mint some FAU tokens for the 0x97D5F5D4fDf83b9D2Cb342A09b8DF297167a73d0 on the [page](https://erc20faucet.com/)

  - Get cosmos user's balance:
    ```
    onomyd query bank balances onomy1e6xtwjw9mgmljyrqw6mlw3nrpuz3p79gct73nw --chain-id onomy
    ```
  
  - Send from eth to cosmos (from 0x97D5F5D4fDf83b9D2Cb342A09b8DF297167a73d0 to onomy1e6xtwjw9mgmljyrqw6mlw3nrpuz3p79gct73nw)
    ```
    gbt -a onomy client eth-to-cosmos \
            --ethereum-key="e0b21b1d80e53f38734a3ed395796956b50c637916ddbb6cedb096b848053d2d" \
            --ethereum-rpc="http://$ETH_HOST:8545" \
            --gravity-contract-address=0x8A0814b7251138Dea19054425D0dfF0C497305d3 \
            --token-contract-address=0xFab46E002BbF0b4509813474841E0716E6730136 \
            --amount=10 \
            --destination=onomy1e6xtwjw9mgmljyrqw6mlw3nrpuz3p79gct73nw
    ```
  
  - Now check users balances on both sides
  
  - Send from cosmos to eth (from onomy1e6xtwjw9mgmljyrqw6mlw3nrpuz3p79gct73nw to 0x2d9480eBA3A001033a0B8c3Df26039FD3433D55d /different eth address)
  
    ```
    gbt -a onomy client cosmos-to-eth --cosmos-phrase="ten stereo fortune girl mean stadium boost maze immune margin rural dragon stage gadget comfort creek cupboard expect satoshi maple machine hunt abstract entry" \
                    --cosmos-grpc="http://0.0.0.0:9090" \
                    --fees=1gravity0xFab46E002BbF0b4509813474841E0716E6730136 \
                    --amount=555555gravity0xFab46E002BbF0b4509813474841E0716E6730136 \
                    --eth-destination=0x2d9480eBA3A001033a0B8c3Df26039FD3433D55d
    
  - Now check users balances on both sides one more time

  ### Using linux or docker/linux

  - Configure your environment
  Run docker with ubuntu and connect to the container
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
  - Set Eth host and onomy host (if the hosts are different, then change them from 0.0.0.0 to your hosts)
  ```
  ETH_HOST=0.0.0.0 && ONOMY_HOST=0.0.0.0
  ```
  - Mint some FAU tokens for the 0x97D5F5D4fDf83b9D2Cb342A09b8DF297167a73d0 on the [page](https://erc20faucet.com/)

  - Check balance of the cosmos user onomy1e6xtwjw9mgmljyrqw6mlw3nrpuz3p79gct73nw on the onomy side
  ```
  curl -X GET "http://$ONOMY_HOST:1317/cosmos/bank/v1beta1/balances/onomy1e6xtwjw9mgmljyrqw6mlw3nrpuz3p79gct73nw" -H "accept: application/json"
  ```

  - Send from eth to cosmos (from 0x97D5F5D4fDf83b9D2Cb342A09b8DF297167a73d0 to onomy1e6xtwjw9mgmljyrqw6mlw3nrpuz3p79gct73nw)
  ```
  ./gbt -a onomy client eth-to-cosmos \
        --ethereum-key="e0b21b1d80e53f38734a3ed395796956b50c637916ddbb6cedb096b848053d2d" \
        --ethereum-rpc="http://$ETH_HOST:8545" \
        --gravity-contract-address=0x8A0814b7251138Dea19054425D0dfF0C497305d3 \
        --token-contract-address=0xFab46E002BbF0b4509813474841E0716E6730136 \
        --amount=10 \
        --destination=onomy1e6xtwjw9mgmljyrqw6mlw3nrpuz3p79gct73nw
  ```
  
  - Check balance of the user onomy1e6xtwjw9mgmljyrqw6mlw3nrpuz3p79gct73nw on the gravity side (should be +10 gravity0xFab46E002BbF0b4509813474841E0716E6730136)
  ```
  curl -X GET "http://$ONOMY_HOST:1317/cosmos/bank/v1beta1/balances/onomy1e6xtwjw9mgmljyrqw6mlw3nrpuz3p79gct73nw" -H "accept: application/json"
  ```

  - Send from cosmos to eth (from onomy1e6xtwjw9mgmljyrqw6mlw3nrpuz3p79gct73nw to 0x2d9480eBA3A001033a0B8c3Df26039FD3433D55d /different eth address)
  ```
  ./gbt -a onomy client cosmos-to-eth --cosmos-phrase="ten stereo fortune girl mean stadium boost maze immune margin rural dragon stage gadget comfort creek cupboard expect satoshi maple machine hunt abstract entry" \
                    --cosmos-grpc="http://$ONOMY_HOST:9090" \
                    --fees=1nom \
                    --amount=1000gravity0xFab46E002BbF0b4509813474841E0716E6730136 \
                    --eth-destination=0x2d9480eBA3A001033a0B8c3Df26039FD3433D55d
  ```
  - Now check users balances on both sides (cosmos/eth) one more time

  - Termination (optional -for the docker only)
  ```
  exit
  ```
  - Restart and attach (optional - for the docker only) 
  ```
  docker start onomy-single-node-runner-rinkeby-client -a
  ```

# Fauset

[docs description](https://github.com/tendermint/faucet)

* install
```
curl https://get.starport.network/faucet! | bash
```

* run
```
faucet -cli-name=onomyd -keyring-backend=test -node=tcp://0.0.0.0:26657 -mnemonic="staff pool flush bundle radar craft sister local fiction clown anger friend torch toss trial drift choice fine exist weapon hamster energy marine hub"
```

* fauset API

[http://0.0.0.0:8000/](http://0.0.0.0:8000/)

* CURL example

Send 10nom to onomy1c3arwy47pvrl5ap0gxjt2ujxc37kyhzle4ltjf
```
curl -X POST "http://0.0.0.0:8000/" -H  "accept: application/json" -H  "Content-Type: application/json" -d "{  \"address\": \"onomy1c3arwy47pvrl5ap0gxjt2ujxc37kyhzle4ltjf\",  \"coins\": [    \"10nom\"  ]}"
```