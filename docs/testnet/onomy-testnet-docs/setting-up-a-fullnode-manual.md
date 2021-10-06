# How to Run a Full Node (Manual Setup - Recommended)

As a Cosmos-based chain, the ONET full nodes are similar to any Cosmos full nodes. Unlike the validator flow, running a full node requires no external software. 

## What do I Need?

A Linux server with any modern Linux distribution, 2gb of RAM and at least 20gb storage. Requirements are minimal.

## Full Node Tutorial

[Pre installation](pre-installation.md) is required prior to setting up a full node.

### Init the config files

```
cd $HOME
onomyd --home $HOME/.onomy init {validator moniker} --chain-id onomy-testnet1
```
Note:- Here the default home directory path is `~/.onomy` or `$HOME/.onomy`, if you want to change it then you can provide in `onomyd` command with flag `--home` as shown in above command. If you are using the default home directory path then its optional to provide in the `onomyd` commands.

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
