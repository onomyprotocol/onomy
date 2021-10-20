# How to Run an Onomy Testnet Full Node

As a Cosmos-based chain, the ONET full nodes are similar to any Cosmos full nodes. Unlike the validator flow, running a full node requires no external software. 

## Getting Started?
Minimum system requirements:
- A Linux server with any modern Linux distribution, 
- A quad-core CPU
- 8 GiB of RAM
- 20gb of storage.

Make sure you have gone through [Pre installation Steps](pre-installation.md) prior to setting up a node.

### Init the config files

```
cd $HOME
onomyd --home $HOME/.onomy init {validator moniker} --chain-id onomy-testnet1
```
Here, validator moniker is a an arbitrary name for your validator

Note:- The default home directory path is `~/.onomy` or `$HOME/.onomy`, if you have to changed it then you need add `--home` flag to `onomyd` command as shown above. If you are using the default path then its optional.

We will refer to onomy home directory as $ONOMY_HOME for here on.

### Copy the genesis file
You need to get the latest genesis file from testnet and replace it with the one in your $ONOMY_HOME directory
```
rm $HOME/.onomy/config/genesis.json
wget http://testnet1.onomy.io:26657/genesis? -O $HOME/raw.json
jq .result.genesis $HOME/raw.json >> $HOME/.onomy/config/genesis.json
rm -rf $HOME/raw.json
```

### Update toml configuration files

1. Make following changes in $ONOMY_HOME/config/config.toml:

	- add *5e0f5b9d54d3e038623ddb77c0b91b559ff13495@testnet1.onomy.io:26656* to seeds field. It should look somthing like this
	`seeds = "5e0f5b9d54d3e038623ddb77c0b91b559ff13495@testnet1.onomy.io:26656"`

	- add *tcp://0.0.0.0:26656* to external_address field
	`external_address = "" to external_address = "tcp://0.0.0.0:26656"`

	- change "tcp://127.0.0.1:26657" to "tcp://0.0.0.0:26657"
	- change "tcp://127.0.0.1:26656" to "tcp://0.0.0.0:26656"
	- change addr_book_strict = true to addr_book_strict = false

2. Make following changes in $ONOMY_HOME/config/app.toml:

	- Change enable = false to enable = true
	- Change swagger = false to swagger = true

## Increasing the default open files limit
If we don't raise this value nodes will crash once the network grows large enough
```
sudo su -c "echo 'fs.file-max = 65536' >> /etc/sysctl.conf"
sysctl -p

sudo su -c "echo '* hard nofile 94000' >> /etc/security/limits.conf"
sudo su -c "echo '* soft nofile 94000' >> /etc/security/limits.conf"

sudo su -c "echo 'session required pam_limits.so' >> /etc/pam.d/common-session"
```
After making these changes, you will either need to reboot your system or close all the ssh sessions and connect again.

To check if changes took place, run 
```
ulimit -n
```
If limit is still 1024, changes are not in effect yet.

### Start your full node in another terminal and wait for it to sync

```
onomyd --home $ONOMY_HOME start
```

### Check the status of the Onomy chain

You should be good to go! You can check the status of the Onomy chain by running.
```
curl http://localhost:26657/status
```
if catching_up is false means your node is fully synced
