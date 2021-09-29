# Onomy testnet faucet

In order to use the faucet and receive testnet tokens you'll first need to setup a wallet

## What do I need?

A Linux server with any modern Linux distribution, 2gb of ram and at least 20gb storage. Requirements are very minimal.

### Download/install Onomy chain binaries
To download binary follow these commands
```
cd $HOME
mkdir binaries
cd binaries
wget https://github.com/onomyprotocol/onomy/releases/download/v0.0.1/onomyd
wget https://github.com/onomyprotocol/onomy/releases/download/v0.0.1/gbt
wget https://github.com/onomyprotocol/onomy/releases/download/v0.0.1/geth
cd ..
chmod -R +x binaries
export PATH=$PATH:$HOME/binaries/
```
or If you have Fedora (Fedora 34) or Redhat (Red Hat Enterprise Linux 8.4 (Ootpa))
and you want to make binaries yourself, then follow these steps
```
sudo yum install -y git
git clone -b dev https://github.com/onomyprotocol/onomy.git
cd onomy/deploy/testnet
bash bin.sh
```
The second way may be unsafe because it used the latest version of the artifacts.
### Init the config files

```
cd $HOME
onomyd --home $HOME/onomy-testnet1/onomy init {validator moniker} --chain-id onomy-testnet1
```
### Generate your key

Be sure to back up the phrase you get! You’ll need it in a bit. If you don't back up the phrase here just follow the steps again to generate a new key.

Note 'myvalidatorkeyname' is just the name of your key here, you can pick anything you like, just remember it later.

You'll be prompted to create a password, I suggest you pick something short since you'll be typing it a lot

```
onomyd --home $HOME/onomy-testnet1/onomy keys add {myvalidatorkeyname} --keyring-backend test --output json | jq . >> $HOME/onomy-testnet1/onomy/validator_key.json
```

### Copy the genesis file

```
rm $HOME/onomy-testnet1/onomy/config/genesis.json
wget http://147.182.190.16:26657/genesis? -O $HOME/raw.json
jq .result.genesis $HOME/raw.json >> $HOME/onomy-testnet1/onomy/config/genesis.json
rm -rf $HOME/raw.json
```

### Add seed node

Change the seed field in $HOME/onomy-testnet1/onomy/config/config.toml to contain the following:

```

replace seeds = "" to seeds = "5e0f5b9d54d3e038623ddb77c0b91b559ff13495@147.182.190.16:26656"

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

### Start your full node and wait for it to sync
```
onomyd --home $HOME/onomy-testnet1/onomy start
```
To check if the node is synced or not, do:
```
onomyd status
```
When the value of catching_up is false, your node is fully sync'd with the network.
```
"sync_info": {
"latest_block_hash": "7BF95EED4EB50073F28CF833119FDB8C7DFE0562F611DF194CF4123A9C1F4640",
"latest_app_hash": "7C0C89EC4E903BAC730D9B3BB369D870371C6B7EAD0CCB5080B5F9D3782E3559",
"latest_block_height": "668538",
"latest_block_time": "2020-10-31T17:50:56.800119764Z",
"earliest_block_hash": "E7CAD87A4FDC47DFDE3D4E7C24D80D4C95517E8A6526E2D4BB4D6BC095404113",
"earliest_app_hash": "",
"earliest_block_height": "1",
"earliest_block_time": "2020-09-15T14:02:31Z",
"catching_up": false
}
```

### Request tokens from the faucet

First list all of your keys

```
onomyd --home $HOME/onomy-testnet1/onomy keys list --keyring-backend test
```

You'll see an output like this

```
- name: keyname
  type: local
  address: cosmos1youraddresswillgohere
  pubkey: cosmospub1yourpublickleywillgohere
  mnemonic: ""
  threshold: 0
  pubkeys: []

```

Copy your address from the 'address' field and paste it into the command below in place of $ONOMY_VALIDATOR_ADDRESS

```
curl -X POST http://147.182.190.16:8000/ -H  "accept: application/json" -H  "Content-Type: application/json" -d "{  \"address\": \"$ONOMY_VALIDATOR_ADDRESS\",  \"coins\": [    \"10nom\"  ]}"
```

This will provide you 10nom from the faucet storage.
You can check balances of account using this command as given below
```
onomyd --home $HOME/onomy-testnet1/onomy query bank balances $ACCOUNT_ADDRESS 

```

