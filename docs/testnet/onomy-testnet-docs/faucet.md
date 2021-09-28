# Onomy testnet faucet

In order to use the faucet and receive testnet tokens you'll first need to setup a wallet

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
git clone -b ONET-65 https://github.com/sunnyk56/market.git
cd market/deploy/onomy-chain
bash bin.sh
```
### Init the config files

```
cd $HOME
onomyd --home $HOME/onomy/onomy init mymoniker --chain-id onomy
```
### Generate your key

Be sure to back up the phrase you get! You’ll need it in a bit. If you don't back up the phrase here just follow the steps again to generate a new key.

Note 'myvalidatorkeyname' is just the name of your key here, you can pick anything you like, just remember it later.

You'll be prompted to create a password, I suggest you pick something short since you'll be typing it a lot

```
cd $HOME
onomyd --home $HOME/onomy/onomy init mymoniker --chain-id onomy
onomyd --home $HOME/onomy/onomy keys add myvalidatorkeyname --keyring-backend test --output json | jq . >> $HOME/onomy/onomy/validator_key.json
```

### Copy the genesis file

```
rm $HOME/onomy/onomy/config/genesis.json
wget http://147.182.128.38:26657/genesis? -O $HOME/raw.json
jq .result.genesis $HOME/raw.json >> $HOME/onomy/onomy/config/genesis.json
rm -rf $HOME/raw.json
```

### Add seed node

Change the seed field in $HOME/onomy/onomy/config/config.toml to contain the following:

```

seeds = "1302d0ed290d74d6f061fb8506e0e34f3f67f7ff@147.182.128.38:26656"

```

### Start your full node and wait for it to sync

Ask what the current blockheight is in the chat

```
onomyd --home $HOME/onomy/onomy start
```

### Request tokens from the faucet

First list all of your keys

```
onomyd --home $HOME/onomy/onomy keys list --keyring-backend test
```

You'll see an output like this

```
- name: jkilpatr
  type: local
  address: cosmos1youraddresswillgohere
  pubkey: cosmospub1yourpublickleywillgohere
  mnemonic: ""
  threshold: 0
  pubkeys: []

```

Copy your address from the 'address' field and paste it into the command below in place of $ONOMY_VALIDATOR_ADDRESS

```
curl -X POST http://147.182.128.38:8000/ -H  "accept: application/json" -H  "Content-Type: application/json" -d "{  \"address\": \"$ONOMY_VALIDATOR_ADDRESS\",  \"coins\": [    \"10nom\"  ]}"
```

This will provide you 10nom from the faucet storage.

