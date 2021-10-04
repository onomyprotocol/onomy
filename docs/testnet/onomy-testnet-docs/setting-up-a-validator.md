# How to become a validator on the Onomy testnet!

## What do I need?

A Linux server with any modern Linux distribution, 4cores, 8gb of ram and at least 20gb of SSD storage.

The following will outline how to turn your node into a validator. You must first have a [full node running](setting-up-a-fullnode.md) in order to become a validator.

### Run script to make it as validator node but before that node should fully sync

Important Note: Here in the script file we have set up implementation with the home directory `$HOME/onomy-testnet1/onomy`. So if you have changed this path then provide the home directory path accordingly in the `onomyd` command.

```
bash peer-validator/makeValidator.sh

Note: 1. Script will ask for enter faucet url get faucet token(Please enter http://testnet1.onomy.io:8000/)
```
You can check the validators of the Onomy chain by running.
```
curl http//localhost:26657/validators
```
validator count will be increased by 1 count

### Setup Gravity bridge

You are now validating on the Onomy blockchain. But as a validator you also need to run the Gravity bridge components or you will be slashed and removed from the validator set after about 16 hours.

### Register your delegate keys

Delegate keys allow the for the validator private keys to be kept in secure storage while the Orchestrator can use it's own delegated keys for Gravity functions. `The delegate keys registration tool will generate Ethereum and Cosmos keys for you if you don't provide any. These will be saved in your local config for later use`.

\*\*If you have set a minimum fee value in your `$HOME/onomy-testnet1/onomy/config/app.toml` modify the `--fees` parameter to match that value!

Here we need validator phrase which we have saved while created validator key
```
cat $HOME/onomy-testnet1/onomy/validator_key.json
```

```

gbt -a onomy init

gbt -a onomy keys register-orchestrator-address --validator-phrase "the phrase you saved earlier" --fees=0nom

```

### Fund your delegate keys

Both your Ethereum delegate key and your Cosmos delegate key will need some tokens to pay gas. On the Onomy chain side you where sent some 'footoken' along with your nom. We're essentially using nom as a gas token for this testnet.

In a production network only relayers would need Ethereum to fund relaying, but for this testnet all validators run relayers by default, allowing us to more easily simulate a lively economy of many relayers.

You should have received 200000000 Onomy Governance Token in nom. We're going to send half of the nom to the delegate address

To get the address for your validator key you can run the below, where 'myvalidatorkeyname' is whatever you named your key in the 'generate your key' step.

```

onomyd --home $HOME/onomy-testnet1/onomy keys show {myvalidatorkeyname} --keyring-backend test

```
Basically there is two way to fund your delegate cosmos address
1. to transfer from a account
```
onomyd --home $HOME/onomy-testnet1/onomy tx bank send <your validator address> <your delegate cosmos address> 1000000nom --chain-id=onomy-testnet1 --keyring-backend test
```
2. using faucet command, from Onomy side faucet funded
```
curl -X POST http://testnet1.onomy.io:8000/ -H  "accept: application/json" -H  "Content-Type: application/json" -d "{  \"address\": \"<your delegate cosmos address>\",  \"coins\": [    \"100000000nom\"  ]}"
```

Now we need some Rinkeby Eth in the Ethereum delegate key

```
https://www.rinkeby.io/#faucet
```

### Setup Geth on the Rinkeby testnet
_Please only run one or the other of the below instructions in new terminal, both will not work_

#### Light client instructions

```
geth --rinkeby --syncmode "light"  --http --http.port "8545"
```

#### Fullnode instructions

```
geth --rinkeby --syncmode "full"  --http --http.port "8545"
```

You'll see this url, please note your ip and share both this node url and your ip in chat to add to the light client nodes list

```
INFO [06-10|14:11:03.104] Started P2P networking self=enode://71b8bb569dad23b16822a249582501aef5ed51adf384f424a060aec4151b7b5c4d8a1503c7f3113ef69e24e1944640fc2b422764cf25dbf9db91f34e94bf4571@127.0.0.1:30303
```

Finally you'll need to wait for several hours until your node is synced, you can not continue with the instructions until your node is synced.

### Deployment of the Gravity contract

Once 66% of the validator set has registered their delegate Ethereum key it is possible to deploy the Gravity Ethereum contract. Once deployed the Gravity contract address on rinkeby will be posted here

Here is the contract address! Move forward!

```

0xB4BAd4Cef22a4EAeF67434644ebaB4cEC54Db37A

```

### Start your Orchestrator

Now that the setup is complete you can start your Orchestrator. Use the Cosmos mnemonic generated in the 'register delegate keys' step and the Ethereum private key also generated in that step. You should setup your Orchestrator in systemd or elsewhere to keep it running and restart it when it crashes.

If your Orchestrator goes down for more than 16 hours during the testnet you will be slashed and booted from the active validator set.

Since you'll be running this a lot I suggest putting the command into a script, like so. The next version of the orchestrator will use a config file for these values and have encrypted key storage.

\*\*If you have set a minimum fee value in your `$HOME/onomy-testnet1/onomy/config/app.toml` modify the `--fees` parameter to match that value!

```
nano start-orchestrator.sh
```

```
#!/bin/bash
gbt --address-prefix="onomy" orchestrator \
        --cosmos-phrase="<registered delegate consmos phrase>" \
        --cosmos-grpc="http://0.0.0.0:9090" \
        --ethereum-key="<registered delegate ethereum private key>" \
        --ethereum-rpc="http://0.0.0.0:8545" \
        --fees="1nom" \
        --gravity-contract-address="0xB4BAd4Cef22a4EAeF67434644ebaB4cEC54Db37A"
```
Please run below command in new terminal.

```
bash start-orchestrator.sh
```

## Next Step

Now that the validator node and gravity bridge setup is successfully done, you must [test gravity bridge](testing-gravity.md).
