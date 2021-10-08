# Get NOM via the Testnet Faucet

In order to test the Onomy faucet and receive faucet testnet NOM tokens, you'll first need to setup a wallet/account.

## What do I need?

You'll need a Linux server with any modern Linux distribution, 2gb of RAM and at least 20gb storage. The requirements are very minimal.

These docs will outline how to setup a wallet/account and test the Onomy faucet. You must first have a [full node running](setting-up-a-fullnode-manual.md).


### Generate your key

Be sure to back up the phrase you will get! You’ll need it in a bit. If you don't back up the phrase here, just follow the steps again to generate a new key.

Note 'account_name' is just the name of your key here. You can pick anything you like, just make sure to remember it later.


```
onomyd --home $HOME/.onomy keys add {account_name} --keyring-backend test --output json | jq . >> $HOME/.onomy/validator_key.json
```
Now you will be able to access your wallet details in `$HOME/.onomy/validator_key.json`. Here, $HOME is path to your home directory. You will be able to see wallet name, wallet type, wallet address, wallet public key and your wallet mnemonic
### Request tokens from the faucet

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

Copy your address from the 'address' field and paste it into the command below in the place of $ONOMY_VALIDATOR_ADDRESS:

```
curl -X POST http://testnet1.onomy.io:8000/ -H  "accept: application/json" -H  "Content-Type: application/json" -d "{  \"address\": \"$ONOMY_VALIDATOR_ADDRESS\",  \"coins\": [    \"10nom\"  ]}"
```

This will provide you 10 NOM from the faucet.

You can check your account balances by using this command:
```
onomyd --home $HOME/.onomy query bank balances $ACCOUNT_ADDRESS 

```
