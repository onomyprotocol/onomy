# Become a Validator on the Onomy Network

## What Do I Need?
Minimum Requirements:
- Any modern Linux distribution 
- A quad-core CPU
- 16 GiB RAM
- 320gb of storage space

You also need to have [full node running](setting-up-a-fullnode.md) before trying to set up a validator.


## How to become validator
In order to become validator, you need to follow steps below
1. [Set up Validator](#validator)
2. [Set up Gravity Bridge](#gravityBridge)
3. [Set up Go Ethereum (GEth)](#GEth)
4. [Start Orchestrator](#orchestrator)


## <a name="validator"></a> 1. Set up Validator

In order to set up your node as a validator, you first need to have a [full-node running](setting-up-a-fullnode.md). Once you have set up a full node and it has synced with the blockchain, you can create a validator.

There are a lot of ways to setup your validator in order to secure it. This document will guide you through setting up sentry node architecture for your validator as well. Setting up a sentry node is optional and it is used only to increase security of the validator node.

First, you can convert your full node to validator by follwoing steps below:

#### Generate your key

Use the following command to generate your keys. Your keys will be stored in `$HOME/.onomy/validator_key.json`.
```
onomyd --home $HOME/.onomy keys add {account_name} --keyring-backend test --output json | jq . >> $HOME/.onomy/validator_key.json
```

You will be able to see wallet name, wallet type, wallet address, wallet public key and your wallet mnemonic

#### Request tokens from the faucet

First, list all of your keys:
```
onomyd --home $HOME/.onomy keys list --keyring-backend test
```
You'll see an output like this:
```
- name: keyname
  type: local
  address: onomy1somerandomtext
  pubkey: onomypub1somemorerandomtext
  mnemonic: ""
  threshold: 0
  pubkeys: []

```
Copy your address from the 'address' field and paste it into the command below in the place of $ONOMY_VALIDATOR_ADDRESS:

```
curl -X POST http://testnet1.onomy.io:8000/ -H  "accept: application/json" -H  "Content-Type: application/json" -d "{  \"address\": \"$ONOMY_VALIDATOR_ADDRESS\",  \"coins\": [    \"100000000nom\"  ]}"
```
This will provide you 100000000nom from the faucet.

#### Send your validator setup transaction, but make sure your node is fully synced before:

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

### \[Optional\]create a sentry structure
In order to increase security of your validator, you can add a few sentry nodes to your validator setup. The validator will only talk to the sentry nodes on a private network. and your sentry nodes in turn will communicate with rest of the network. Diagram below shows the sentry architecture.

![Sentry Architecture](sentry_architecture.png)


In order to set up the sentry nodes, we will first setup a few [full nodes](setting-up-a-fullnode.md). Once you initialize the full node, either using the script or manually, do not start it. We will need to change some parameters in`$ONOMY_HOME/config/config.toml`

change the following parameters in config file of sentry nodes:
```
pex: true
persistent-peers: nodeid@ip:port (for validator node)
private-peer-ids: node id of validator
unconditional-peer-ids: node id of validator and optionally sentries
addr-book-strict: false
```

change the following parameters in config file of validator node:
```
pex: false
persistent-peers: nodeid@ip:port (for all the sentry nodes)
unconditional-peer-ids: optionally sentries
addr-book-strict: false
double-sign-check-height: 10
```

After these changes, restart both validator and sentry nodes and your sentry structure should be ready.

Also dont forget to make necessary changes to firewall in your validator node to stop unauthorized access.

#### Confirm that you are validating

You can check your validator status by using the following command
```
onomyd query staking validator "$($ONOMY keys show $ONOMY_VALIDATOR_NAME --bech val --address --keyring-backend test)"
```

If you see `status: BOND_STATUS_BONDED`, you are validating.

## <a name="gravityBridge"></a> 2. Setup Gravity bridge

You are now validating on the Onomy Network. However, as a validator, you also need to run the Gravity bridge components or you will be slashed and removed from the validator set after about 16 hours.

### Register your delegate keys

Delegate keys allow the for the validator private keys to be kept in secure storage while the Orchestrator can use its own delegated keys for Gravity functions. `The delegate keys registration tool will generate Ethereum and Cosmos keys for you if you don't provide any. These will be saved in your local config for later use`.

\*\*If you have set a minimum fee value in your `$HOME/.onomy/config/app.toml`, modify the `--fees` parameter to match that value!

Here, you will need validator phrase which you have saved while creating the validator key:
```
cat $HOME/.onomy/validator_key.json
```
```
gbt -a onomy init

gbt -a onomy keys register-orchestrator-address --validator-phrase "the phrase you saved earlier" --fees=0nom

```
### Fund your delegate keys

Both your Ethereum delegate key and your Cosmos delegate key will need some tokens to pay for gas. On the Onomy Network side you were sent some 'footoken' along with your ANOM. We're essentially using ANOM as a gas token for this testnet.

In a production network, only relayers would need Ethereum to fund relaying, but for this testnet, all validators run relayers by default, allowing us to more easily simulate a lively economy of many relayers.

You should have received 100000000 Onomy ANOM tokens.

To get the address for your validator key, you can run the command below, where 'myvalidatorkeyname' is whatever you named your key in the 'generate your key' step:

```
onomyd --home $HOME/.onomy keys show {myvalidatorkeyname} --keyring-backend test
```
Basically, there are two ways to fund your delegate cosmos address:

1. Transfer from an account
```
onomyd --home $HOME/.onomy tx bank send <your validator address> <your delegate cosmos address> 1000000nom --chain-id=onomy-testnet1 --keyring-backend test
```
2. Using faucet command, from the Onomy-side faucet:
```
curl -X POST http://testnet1.onomy.io:8000/ -H  "accept: application/json" -H  "Content-Type: application/json" -d "{  \"address\": \"<your delegate cosmos address>\",  \"coins\": [    \"100000000nom\"  ]}"
```
 Now, we need some Rinkeby ETH in the Ethereum delegate key. You can get some RInkeby ETH from official Rinkeby faucet:
```
https://www.rinkeby.io/#faucet
```

## <a name="GEth"></a> 3. Setup Geth on the Rinkeby testnet

We will be using Geth Ethereum light clients for this task. For production Gravity, we suggest that you point your Orchestrator at a Geth light client and then configure your light client to peer with the full nodes that you control. This provides higher reliability as light clients are very quick to start/stop and resync, allowing you to rebuild an Ethereum full node without having days of Orchestrator downtime, for example.

Geth full nodes do not serve light clients by default, as light clients do not trust full nodes, but if there are no full nodes to request proofs from, they cannot operate. Therefore, we are collecting the largest possible list of Geth full nodes from our community that will serve light clients.

If you have more than 40gb of free storage, an SSD, and extra memory/CPU power, please run a full node and share the node url. If you do not, please use the the light client instructions.

_Please only run one or the other of the below instructions in new terminal; running both will not work_

### Light client instructions

```
geth --rinkeby --syncmode "light"  --http --http.port "8545"
```

### Fullnode instructions

```
geth --rinkeby --syncmode "full"  --http --http.port "8545"
```

You'll see this url, please note your IP and share this node's url and your ip in chat to be added to the light client nodes list

```
INFO [06-10|14:11:03.104] Started P2P networking self=enode://71b8bb569dad23b16822a249582501aef5ed51adf384f424a060aec4151b7b5c4d8a1503c7f3113ef69e24e1944640fc2b422764cf25dbf9db91f34e94bf4571@127.0.0.1:30303
```

Finally, you'll need to wait for several hours until your node is synced. You cannot continue with the instructions before your node is fully synced.

## <a name="orchestrator"></a> 4. Start your Orchestrator

Now that the setup is complete, you can start your Orchestrator. Use the Cosmos mnemonic and Ethereum private key generated in the 'register delegate keys' step. You should setup your Orchestrator in systemd or elsewhere to keep it running and restart it when it crashes.

If your Orchestrator goes down for more than 16 hours during the testnet, you will be slashed and booted from the active validator set.

Since you'll be running this a lot, we suggest putting the command into a script. The next version of the orchestrator will use a config file for these values and have encrypted key storage.

\*\*If you have set a minimum fee value in your `$HOME/.onomy/config/app.toml` modify the `--fees` parameter to match that value!

Shown below is the command to start orchestrator

```
gbt --address-prefix="onomy" orchestrator \
        --cosmos-phrase="<registered delegate consmos phrase>" \
        --cosmos-grpc="http://0.0.0.0:9090" \
        --ethereum-key="<registered delegate ethereum private key>" \
        --ethereum-rpc="http://0.0.0.0:8545" \
        --fees="1nom" \
        --gravity-contract-address="0xB4BAd4Cef22a4EAeF67434644ebaB4cEC54Db37A"
```

## Next Step

Now that the validator node and gravity bridge setup is successfully done, you must [test gravity bridge](testing-gravity.md).
