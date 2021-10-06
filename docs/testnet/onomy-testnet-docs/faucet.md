# Get NOM via the Testnet Faucet

In order to test the Onomy faucet and receive faucet testnet NOM tokens, you'll first need to setup a wallet/account.

## What do I need?

You'll need a Linux server with any modern Linux distribution, 2gb of RAM and at least 20gb storage. The requirements are very minimal.

These docs will outline how to setup a wallet/account and test the Onomy faucet. You must first have a [full node running](setting-up-a-fullnode-manual.md).


### Generate your key

Be sure to back up the phrase you will get! You’ll need it in a bit. If you don't back up the phrase here, just follow the steps again to generate a new key.

Note 'account_name' is just the name of your key here. You can pick anything you like, just make sure to remember it later.

You'll be prompted to create a password. We suggest you pick something short since you'll be typing it a lot.

```
onomyd --home $HOME/.onomy keys add {account_name} --keyring-backend test --output json | jq . >> $HOME/.onomy/validator_key.json
```

### Request tokens from the faucet

First, list all of your keys:

```
onomyd --home $HOME/.onomy keys list --keyring-backend test
```

You'll see an output like this:

```
- name: keyname
  type: local
  address: cosmos1youraddresswillgohere
  pubkey: cosmospub1yourpublickleywillgohere
  mnemonic: ""
  threshold: 0
  pubkeys: []

```

Copy your address from the 'address' field and paste it into the command below in the place of $ONOMY_VALIDATOR_ADDRESS:

```
curl -X POST http://testnet1.onomy.io:8000/ -H  "accept: application/json" -H  "Content-Type: application/json" -d "{  \"address\": \"$ONOMY_VALIDATOR_ADDRESS\",  \"coins\": [    \"10nom\"  ]}"
```

This will provide you 10 NOM from the faucet storage.

You can check your account balances by using this command line:
```
onomyd --home $HOME/.onomy query bank balances $ACCOUNT_ADDRESS 

```
