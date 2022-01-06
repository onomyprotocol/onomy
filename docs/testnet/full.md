# How to Run an Onomy Testnet Full Node

As a Cosmos-based chain, the ONET full nodes are similar to any Cosmos full nodes. Unlike the validator flow, running a
full node requires no external software.

## Getting Started

System requirements:

- Any modern Linux distribution (RHEL 8 or Fedora 36 are preferred)
- A quad-core CPU
- 16 GiB RAM
- 320gb of storage space

Make sure you have gone through [Installation Steps](installation.md) prior to setting up a node.

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

Pass through the [full node](../../deploy/testnet/full.md) instruction

### Check the status of the Onomy Network Testnet

You should be good to go! You can check the status of the Onomy chain by running:

```
curl http://<your_ip>:26657/status
```

If catching_up is false, this means that your node is fully synced. Congratulations!

## <a name="advancedMethod"></a> 2. Advanced/Manual Method

### Init the config files

```
cd $HOME
onomyd --home $HOME/.onomy init {validator moniker} --chain-id onomy-testnet1
```

Here, validator moniker is a an arbitrary name for your validator

Note:- The default home directory path is `~/.onomy` or `$HOME/.onomy`, if you have to changed it then you need
add `--home` flag to `onomyd` command as shown above. If you are using the default path then its optional.

We will refer to onomy home directory as $ONOMY_HOME for here on.

### Copy the genesis file

You need to get the latest genesis file from testnet and replace it with the one in your $ONOMY_HOME directory

```
rm $HOME/.onomy/config/genesis.json
wget http://64.9.136.120:26657/genesis? -O $HOME/raw.json
jq .result.genesis $HOME/raw.json >> $HOME/.onomy/config/genesis.json
rm -rf $HOME/raw.json
```

### Update toml configuration files

1. Make following changes in $ONOMY_HOME/config/config.toml:

   - add *dc6c6e2faa000f8126e537b622eb081faeac62d3@64.9.136.119:26656,48a0784919f1ed01a875962aafdef75c3d9cb878@64.9.136.120:26656,c63bfabfd5379f7485e45bfd548c1e39bbc22b97@145.40.81.21:26656* to seeds field. It should look somthing
     like this
     `seeds = "dc6c6e2faa000f8126e537b622eb081faeac62d3@64.9.136.119:26656,48a0784919f1ed01a875962aafdef75c3d9cb878@64.9.136.120:26656,c63bfabfd5379f7485e45bfd548c1e39bbc22b97@145.40.81.21:26656"`

   - add *tcp://<your_ip>:26656* to external_address field
     `external_address = "" to external_address = "tcp://<your_ip>:26656"`
     Where the <your_ip> is the external IP of your node.

   - change "tcp://127.0.0.1:26657" to "tcp://0.0.0.0:26657"
   - change "tcp://127.0.0.1:26656" to "tcp://0.0.0.0:26656"
   - change addr_book_strict = true to addr_book_strict = false

2. Optionally make following changes in $ONOMY_HOME/config/app.toml:

   - Change enable = false to enable = true
   - Change swagger = false to swagger = true

## Increasing the default open files limit

The example for ([Red Hat Enterprise Linux](../../deploy/testnet/set-ulimit-rhel8.md)).

If we don't raise this value nodes will crash once the network grows large enough

After making these changes, you will either need to reboot your system or close all the ssh sessions and connect again.

To check if changes took place, run

```
ulimit -n
```

If limit is still 1024, changes are not in effect yet.

### Start your full node in another terminal and wait for it to sync

```
onomyd start
```

### Check the status of the Onomy chain

You should be good to go! You can check the status of the Onomy chain by running.

```
curl http://<your_ip>:26657/status
```
Or by visiting the same URL from the browser.
If catching_up is false means your node is fully synced.
Optionally add the "onomyd start" as a linux service or add to a crontab