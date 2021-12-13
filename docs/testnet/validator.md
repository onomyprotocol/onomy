# Validator on the Onomy Network

1. [Become a Validator on the Onomy Network](#validator)
2. [Jailing/Unjailing](#jailing)

## <a name="validator"></a> Become a Validator on the Onomy Network

## What Do I Need?

System requirements:

- Any modern Linux distribution (RHEL 8 or Fedora 36 are preferred)
- A quad-core CPU
- 16 GiB RAM
- 320gb of storage space

## How to run a validator node

There are two ways for setting up the node, one is the standard method where everything is done by running a
pre-configured script while second one is for more advanced users who might want to customize their installation:

1. [Standard Method](#standardMethod)
2. [Advanced/Manual Method](#advancedMethod)

## <a name="standardMethod"></a> 1. Standard Method

### Initiate chain

Clone the latest release of github repo

```
git clone https://github.com/onomyprotocol/onomy.git
```

Pass through the [validator](../../deploy/testnet/validator.md) instruction

## <a name="advancedMethod"></a> 2. Advanced/Manual Method

In order to become validator, you need to use or Advance with the steps below.

1. [Set up Validator](#validator)
2. [Set up Gravity Bridge](#gravityBridge)
3. [Set up Go Ethereum (Geth)](#Geth)
4. [Start Orchestrator](#orchestrator)

## <a name="validator"></a> 1. Set up Validator

In order to set up your node as a validator, you first need to have a [full-node running](full.md). Once you have set up
a full node and it has synced with the blockchain, you can create a validator.

There are a lot of ways to setup your validator in order to secure it. This document will guide you through setting up
sentry node architecture for your validator as well. Setting up a sentry node is optional and it is used only to
increase security of the validator node.

First, you can convert your full node to validator by follwoing steps below:

#### Generate your key if you already not have them

Use the following command to generate your keys. Your keys will be stored in `$HOME/.onomy/validator_key.json`.

```
onomyd --home $HOME/.onomy keys add {account_name} --keyring-backend test --output json | jq . >> $HOME/.onomy/validator_key.json
```

You will be able to see wallet name, wallet type, wallet address, wallet public key and your wallet mnemonic

#### Request tokens (for validator) from the master account by "text" request to onomy team.

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

The 'address' field aif your validator address.

#### Send your validator setup transaction, but make sure your node is fully synced before:

```
onomyd --home $HOME/.onomy tx staking create-validator \
 --amount=10000000000000000000anom \
 --pubkey=$(onomyd --home $HOME/.onomy tendermint show-validator) \
 --moniker="put your validator name here" \
 --chain-id=onomy-testnet \
 --commission-rate="0.10" \
 --commission-max-rate="0.20" \
 --commission-max-change-rate="0.01" \
 --min-self-delegation="1" \
 --gas="auto" \
 --gas-adjustment=1.5 \
 --gas-prices="1anom" \
 --from="put your validator key name here" \
 --keyring-backend test
```

### Optionally create a sentry structure

In order to increase security of your validator, you can add a few sentry nodes to your validator setup. The validator
will only talk to the sentry nodes on a private network. and your sentry nodes in turn will communicate with rest of the
network. Diagram below shows the sentry architecture.

![Sentry Architecture](sentry_architecture.png)

In order to set up the sentry nodes, we will first setup a few [full nodes](full.md). Once you initialize the full node,
either using the script or manually, do not start it. We will need to change some parameters
in`$ONOMY_HOME/config/config.toml`

change the following parameters in config file of sentry nodes:

```
pex: true
persistent_peers: nodeid@ip:port (for validator node)
private_peer_ids: node id of validator
unconditional_peer_ids: node id of validator and optionally sentries
addr_book_strict: false
```

change the following parameters in config file of validator node:

```
pex: false
persistent_peers: nodeid@ip:port (for all the sentry nodes)
unconditional_peer_ids: optionally sentries
addr_book_strict: false
double_sign_check_height: 10
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

You are now validating on the Onomy Network. However, as a validator, you also need to run the Gravity bridge components
or you will be slashed and removed from the validator set after about 16 hours.

### Register your delegate keys

Delegate keys allow the for the validator private keys to be kept in secure storage while the Orchestrator can use its
own delegated keys for Gravity functions. The delegate keys registration tool will generate Ethereum and Cosmos keys for
you if you don't provide any. These will be saved in your local config for later use .

Here, you will need validator phrase which you have saved while creating the validator key:

```
cat $HOME/.onomy/validator_key.json
```

```
gbt -a onomy init

gbt -a onomy keys register-orchestrator-address --validator-phrase "the phrase you saved earlier" --fees=0anom

```

### Fund your delegate keys

Both your Ethereum delegate key, and your Cosmos delegate key will need some tokens to pay for gas. On the Onomy Network
side you were sent some 'footoken' along with your ANOM. We're essentially using ANOM as a gas token for this testnet.

To get the address for your validator key, you can run the command below, where 'myvalidatorkeyname' is whatever you
named your key in the 'generate your key' step:

```
onomyd --home $HOME/.onomy keys show {myvalidatorkeyname} --keyring-backend test
```

Delegate from an account

```
onomyd --home $HOME/.onomy tx bank send <your validator address> <your delegate cosmos address> 1000000anom --chain-id=onomy-testnet --keyring-backend test
```

Now, we need some Rinkeby ETH in the Ethereum delegate key. You can get some Rinkeby ETH from official Rinkeby faucet:

```
https://www.rinkeby.io/#faucet
```

## <a name="Geth"></a> 3. Setup Geth on the Rinkeby testnet

We will be using Geth Ethereum clients for this task. For production Gravity, we suggest that you point your
Orchestrator at a Geth client and then configure your client to peer with the full nodes that you control. This provides
higher reliability as clients are very quick to start/stop and resync, allowing you to rebuild an Ethereum full node
without having days of Orchestrator downtime, for example.

Geth full nodes do not serve clients by default, as clients do not trust full nodes, but if there are no full nodes to
request proofs from, they cannot operate. Therefore, we are collecting the largest possible list of Geth full nodes from
our community that will serve clients.

If you have more than 40gb of free storage, an SSD, and extra memory/CPU power, please run a full node and share the
node url. If you do not, please use the the client instructions.

_Please only run one or the other of the below instructions in new terminal; running both will not work_

```
geth --rinkeby --syncmode "full"  --http --http.port "8545"
```

You'll see this url, please note your IP and share this node's url and your ip in chat to be added to the client nodes
list

```
INFO [06-10|14:11:03.104] Started P2P networking self=enode://71b8bb569dad23b16822a249582501aef5ed51adf384f424a060aec4151b7b5c4d8a1503c7f3113ef69e24e1944640fc2b422764cf25dbf9db91f34e94bf4571@127.0.0.1:30303
```

Finally, you'll need to wait for several hours until your node is synced. You cannot continue with the instructions
before your node is fully synced.

## <a name="orchestrator"></a> 4. Start your Orchestrator

Now that the setup is complete, you can start your Orchestrator. Use the Cosmos mnemonic and Ethereum private key
generated in the 'register delegate keys' step. You should setup your Orchestrator in systemd or elsewhere to keep it
running and restart it when it crashes.

If your Orchestrator goes down for more than 16 hours during the testnet, you will be slashed and booted from the active
validator set.

Since you'll be running this a lot, we suggest putting the command into a script. The next version of the orchestrator
will use a config file for these values and have encrypted key storage.

If you have set a minimum fee value in your `$HOME/.onomy/config/app.toml` modify the `--fees` parameter to match that
value!

Shown below is the command to start orchestrator

```
gbt --address-prefix="onomy" orchestrator \
        --cosmos-phrase="<registered delegate consmos phrase>" \
        --cosmos-grpc="http://0.0.0.0:9090" \
        --ethereum-key="<registered delegate ethereum private key>" \
        --ethereum-rpc="http://0.0.0.0:8545" \
        --fees="1anom" \
        --gravity-contract-address="0x119999cf67269C29CEC39337106e00bBbcd68bf9"
```

## Next Step

Now that the validator node and gravity bridge setup is successfully done, you
can [test gravity bridge](testing-gravity.md).

# <a name="jailing"></a> Jailing/Unjailing

## What is Jailing

When a validator disconnects from the network due to connection loss or server fail or it double signs, it needs to be
eliminated from the validator list. It is known as 'jailing'. A validator is jailed if it fails to validate at least 50%
of the last 100 blocks.

When jailed due to downtime, the validator's total stake is slashed by 1%. and if jailed due to double signing,
validaor's total stake is slashed by 5%.

Once jailed, validators can be unjailed again after 10 minutes. These configurations can be found in the genesis file
under the slashing section

```
"slashing": {
      "params": {
        "signed_blocks_window": "100",
        "min_signed_per_window": "0.500000000000000000",
        "downtime_jail_duration": "600s",
        "slash_fraction_double_sign": "0.050000000000000000",
        "slash_fraction_downtime": "0.010000000000000000"
      },
      "signing_infos": [],
      "missed_blocks": []
    }
```

## Unjailing validator

In order to unjail the validator, you may run the following command once 10 minutes have passed

```
onomyd tx slashing unjail --from <validator-name> --chain-id=onomy-testnet --gas auto --gas-adjustment 1.5 --gas-prices 1anom --keyring-backend test
```



