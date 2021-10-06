# Onomy testnet faucet

In order to test the Onomy faucet and receive Onomy faucet tokens you'll first need to setup a wallet/account

## What do I need?

A Linux server with any modern Linux distribution, 2gb of ram and at least 20gb storage. Requirements are very minimal.

The following will outline how to setup wallet/account and test the Onomy faucet. You must first have a [full node running](setting-up-a-fullnode-manual.md).


### Generate your key


Note 'account_name' is just the name of your key here, you can pick anything you like, just remember it later.

```
onomyd --home $HOME/.onomy keys add {account_name} --keyring-backend test --output json | jq . >> $HOME/.onomy/validator_key.json
```

### Request tokens from the faucet

First list all of your keys

```
onomyd --home $HOME/.onomy keys list --keyring-backend test
```

You'll see an output like this

```
- name: keyname
  type: local
  address: onomy1youraddresswillgohere
  pubkey: onomypub1yourpublickleywillgohere
  mnemonic: ""
  threshold: 0
  pubkeys: []

```

Copy your address from the 'address' field and paste it into the command below in place of $ONOMY_VALIDATOR_ADDRESS

```
curl -X POST http://testnet1.onomy.io:8000/ -H  "accept: application/json" -H  "Content-Type: application/json" -d "{  \"address\": \"$ONOMY_VALIDATOR_ADDRESS\",  \"coins\": [    \"10nom\"  ]}"
```

This will provide you 10nom from the faucet storage.
You can check balances of account using this command as given below
```
onomyd --home $HOME/.onomy query bank balances $ACCOUNT_ADDRESS 

```
