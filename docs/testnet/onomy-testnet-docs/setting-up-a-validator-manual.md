# How to become a validator on the Onomy testnet!

## What do I need?

A Linux server with any modern Linux distribution, 16cores, 16gb of ram and at least 320gb of SSD storage.

Onomy chain can be run on Windows and Mac. Binaries are provided on the releases page. But validator instructions are not provided.

I also suggest an open notepad or other document to keep track of the keys you will be generating.

The following will outline how to turn your node into a validator. You must first have a [full node running](setting-up-a-fullnode-manual.md) in order to become a validator.

### Generate your key

Be sure to back up the phrase you get! You’ll need it in a bit. If you don't back up the phrase here just follow the steps again to generate a new key.

Note 'myvalidatorkeyname' is just the name of your key here, you can pick anything you like, just remember it later.

You'll be prompted to create a password, I suggest you pick something short since you'll be typing it a lot

```
cd $HOME
onomyd --home $HOME/.onomy keys add {myvalidatorkeyname} --keyring-backend test --output json | jq . >> $HOME/.onomy/validator_key.json
```

### Request some funds be sent to your address

First find your address

```
onomyd --home $HOME/.onomy keys list --keyring-backend test
```

Copy your address from the 'address' field and paste it into the command below in place of $ONOMY_VALIDATOR_ADDRESS

```
curl -X POST http://testnet1.onomy.io:8000/ -H  "accept: application/json" -H  "Content-Type: application/json" -d "{  \"address\": \"$ONOMY_VALIDATOR_ADDRESS\",  \"coins\": [    \"100000000nom\"  ]}"
```

This will provide you 100000000nom from the faucet storage.

### Send your validator setup transaction but before that node should fully sync

```

onomyd --home $HOME/.onomy tx staking create-validator \
 --amount=100000000nom \
 --pubkey=$(onomyd --home $HOME/.onomy tendermint show-validator) \
 --moniker="put your validator name here" \
 --chain-id=onomy-testnet1 \
 --commission-rate="0.10" \
 --commission-max-rate="0.20" \
 --commission-max-change-rate="0.01" \
 --min-self-delegation="1" \
 --gas="auto" \
 --gas-adjustment=1.5 \
 --gas-prices="1nom" \
 --from="put your validator key name here" \
 --keyring-backend test

```

### Confirm that you are validating

If you see one line in the response you are validating. If you don't see any output from this command you are not validating. Check that the last command ran successfully.

Be sure to replace 'my validator key name' with your actual key name. If you want to double check you can see all your keys with 'onomyd --home $HOME/.onomy keys list --keyring-backend test'

```

onomyd --home $HOME/.onomy query staking validator $(onomyd --home $HOME/.onomy keys show <myvalidatorkeyname> --bech val --address --keyring-backend test)

```

### Setup Gravity bridge

You are now validating on the Onomy blockchain. But as a validator you also need to run the Gravity bridge components or you will be slashed and removed from the validator set after about 16 hours.

### Register your delegate keys

Delegate keys allow the for the validator private keys to be kept in secure storage while the Orchestrator can use it's own delegated keys for Gravity functions. `The delegate keys registration tool will generate Ethereum and Cosmos keys for you if you don't provide any. These will be saved in your local config for later use`.

\*\*If you have set a minimum fee value in your `$HOME/.onomy/config/app.toml` modify the `--fees` parameter to match that value!

Here we need validator phrase which we have saved while created validator key
```
cat $HOME/.onomy/validator_key.json
```
```
gbt -a onomy init

gbt -a onomy keys register-orchestrator-address --validator-phrase "the phrase you saved earlier" --fees=0nom

```

### Fund your delegate keys

Both your Ethereum delegate key and your Cosmos delegate key will need some tokens to pay gas. On the Onomy chain side you where sent some 'footoken' along with your nom. We're essentially using nom as a gas token for this testnet.

In a production network only relayers would need Ethereum to fund relaying, but for this testnet all validators run relayers by default, allowing us to more easily simulate a lively economy of many relayers.

You should have received 100000000 Onomy Governance Token in nom.

To get the address for your validator key you can run the below, where 'myvalidatorkeyname' is whatever you named your key in the 'generate your key' step.

```
onomyd --home $HOME/.onomy keys show {myvalidatorkeyname} --keyring-backend test
```
Basically there is two-way to fund your delegate cosmos address
1. to transfer from a account
```
onomyd --home $HOME/.onomy tx bank send <your validator address> <your delegate cosmos address> 1000000nom --chain-id=onomy-testnet1 --keyring-backend test
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

We will be using Geth Ethereum light clients for this task. For production Gravity we suggest that you point your Orchestrator at a Geth light client and then configure your light client to peer with full nodes that you control. This provides higher reliability as light clients are very quick to start/stop and resync. Allowing you to for example rebuild an Ethereum full node without having days of Orchestrator downtime.

Geth full nodes do not serve light clients by default, light clients do not trust full nodes, but if there are no full nodes to request proofs from they can not operate. Therefore we are collecting the largest possible
list of Geth full nodes from our community that will serve light clients.

If you have more than 40gb of free storage, an SSD and extra memory/CPU power, please run a full node and share the node url. If you do not, please use the light client instructions

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

\*\*If you have set a minimum fee value in your `$HOME/.onomy/config/app.toml` modify the `--fees` parameter to match that value!

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
