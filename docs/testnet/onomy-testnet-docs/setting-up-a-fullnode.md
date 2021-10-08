# How to Run an Onomy Testnet Full Node

As a Cosmos-based chain, the ONET full nodes are similar to any Cosmos full nodes. Unlike the validator flow, running a full node requires no external software. 

## What do I need?

A Linux server with any modern Linux distribution, 4cores, 8gb of RAM and at least 20gb of SSD storage.

[Pre installation](pre-installation.md) is required prior to setting up a node.

There are two ways for setting up the node, one is the standard method where everything is done by creating a script while second one is for more advanced users who might want to customize their installation
1. [Standard Method](#standardMethod)
2. [Advanced Method](#advancedMethod)

## <a name="standardMethod"> 1. Standard Method
### Initiate chain

Clone the latest release of github repo
```
git clone -b v0.0.1 https://github.com/onomyprotocol/onomy.git
cd onomy/deploy/testnet
```

### Run the first time bootstrapping playbook and script

This script will run the needed commands to generate keys and store these in files. Be mindful of the generated keys, as you’ll need them later.

Important Note: Here in the script file, we have set up implementation with the home directory `$HOME/.onomy/onomy-testnet1/onomy`. So, if you have changed this path, then provide the home directory path accordingly in the `onomyd` command.


```
bash peer-validator/init.sh
```
Note: 1. Script will ask for validator name (Type any name for example validator1)
      2. Script will ask for node-id of any validator that is running in chain to add seed (Please enter 5e0f5b9d54d3e038623ddb77c0b91b559ff13495)
      3. Script will ask for ip of validator for which you have added node-id (Please enter testnet1.onomy.io)


Now, it's finally up and has started the sycing block

### Check the status of the Onomy Network Testnet

You should be good to go! You can check the status of the three Onomy chain by running:
```
curl http://localhost:26657/status
```
If catching_up is false, this means that your node is fully synced. Congratulations!






# How to Run a Full Node (Manual Setup - Recommended)

## <a name="advancedMethod"> 1. Advanced Method
### Init the config files

```
cd $HOME
onomyd --home $HOME/.onomy init {validator moniker} --chain-id onomy-testnet1
```
Note:- Here the default home directory path is `~/.onomy` or `$HOME/.onomy`, if you have to changed it then you need add `--home` flag to `onomyd` command as shown above. If you are using the default path then its optional.

### Copy the genesis file

```
rm $HOME/.onomy/config/genesis.json
wget http://testnet1.onomy.io:26657/genesis? -O $HOME/raw.json
jq .result.genesis $HOME/raw.json >> $HOME/.onomy/config/genesis.json
rm -rf $HOME/raw.json
```

### Update toml configuration files

Change in $HOME/.onomy/config/config.toml to contain the following:

```

replace seeds = "" to seeds = "5e0f5b9d54d3e038623ddb77c0b91b559ff13495@testnet1.onomy.io:26656"
replace "tcp://127.0.0.1:26657" to "tcp://0.0.0.0:26657"
replace "tcp://127.0.0.1:26656" to "tcp://0.0.0.0:26656"
replace addr_book_strict = true to addr_book_strict = false
replace external_address = "" to external_address = "tcp://0.0.0.0:26656"
```

Change in $HOME/.onomy/config/app.toml to contain the following:

```
replace enable = false to enable = true
replace swagger = false to swagger = true
```
## Increasing the default open files limit
If we don't raise this value nodes will crash once the network grows large enough
```
sudo su -c "echo 'fs.file-max = 65536' >> /etc/sysctl.conf"
sysctl -p

sudo su -c "echo '* hard nofile 94000' >> /etc/security/limits.conf"
sudo su -c "echo '* soft nofile 94000' >> /etc/security/limits.conf"

sudo su -c "echo 'session required pam_limits.so' >> /etc/pam.d/common-session"
```
For this to take effect you'll need to (A) reboot (B) close and re-open all ssh sessions.
To check if this has worked run
```
ulimit -n
```
If you see 1024 then you need to reboot
### Start your full node in another terminal and wait for it to sync

```
onomyd --home $HOME/.onomy start
```
### Check the status of the Onomy chain

You should be good to go! You can check the status of the three
Onomy chain by running.
```
curl http://localhost:26657/status
```
if catching_up is false means your node is fully synced
