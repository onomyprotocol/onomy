# How to run a Onomy testnet full node

A Onomy chain full node is just like any other Cosmos chain full node and unlike the validator flow requires no external software

## What do I need?

A Linux server with any modern Linux distribution, 2gb of ram and at least 20gb storage. Requirements are very minimal.

### Download/install Onomy chain binaries
```
To download binary follow these commands
mkdir binaries
cd binaries
wget https://github.com/onomyprotocol/onomy/releases/download/v0.0.1/onomyd
wget https://github.com/onomyprotocol/onomy/releases/download/v0.0.1/gbt
wget https://github.com/onomyprotocol/onomy/releases/download/v0.0.1/geth
cd ..
chmod -R +x binaries
export PATH=$PATH:$HOME/binaries/


or If you have Fedora (Fedora 34) or Redhat (Red Hat Enterprise Linux 8.4 (Ootpa))
 and you want to make binaries yourself, then follow these steps

sudo yum install -y git
git clone -b dev https://github.com/onomyprotocol/onomy.git
cd onomy/deploy/testnet
bash bin.sh
```

### Init the config files

```
cd $HOME
onomyd --home $HOME/onomy/onomy init mymoniker --chain-id onomy
```

### Copy the genesis file

```
rm $HOME/onomy/onomy/config/genesis.json
wget http://147.182.128.38:26657/genesis? -O $HOME/raw.json
jq .result.genesis $HOME/raw.json >> $HOME/onomy/onomy/config/genesis.json
rm -rf $HOME/raw.json
```

### Update toml configuration files

Change in $HOME/onomy/onomy/config/config.toml to contain the following:

```

replace seeds = "" to seeds = "1302d0ed290d74d6f061fb8506e0e34f3f67f7ff@147.182.128.38:26656"
replace "tcp://127.0.0.1:26657" to "tcp://0.0.0.0:26657"
replace "tcp://127.0.0.1:26656" to "tcp://0.0.0.0:26656"
replace addr_book_strict = true to addr_book_strict = false
replace external_address = "" to external_address = "tcp://0.0.0.0:26656"
```

Change in $HOME/onomy/onomy/config/app.toml to contain the following:

```
replace enable = false to enable = true
replace swagger = false to swagger = true
```

### Start your full node in another terminal and wait for it to sync

```
onomyd --home $HOME/onomy/onomy start
```
### Check the status of the Onomy chain

You should be good to go! You can check the status of the three
Onomy chain by running.
```
curl http://localhost:26657/status
```
if catching_up is false means your node is fully synced