# How to become a validator on the Onomy testnet!

The following will outline how to turn your node into a validator. You must first have a [full node running](onomy-testnet-docs/setting-up-a-full-node.md) in order to become a validator.

### Generate your key
Be sure to back up the phrase you get! You’ll need it in a bit. If you don't back up the phrase here just follow the steps again to generate a new key.

Note `myvalidatorkeyname` is just the name of your key here, you can pick anything you like, just remember it later. This is *not* the same as your `moniker-name`.

You'll be prompted to create a password, I suggest you pick something short since you'll be typing it a lot

```
onomyd keys add {myvalidatorkeyname} \
--keyring-backend test \
--home $HOME/$CHAIN_ID/onomy \
--output json | jq . >> $HOME/$CHAIN_ID/onomy/validator_key.json
```

Confirm a key was created by doing:
```bash:
onomyd keys list --keyring-backend test --home $HOME/$CHAIN_ID/onomy 
```

You should see something like:
```bash:
- name: lav.test
  type: local
  address: onomy1anlymr6kr7ns6t3c0lhjcvse7aqegffwm3vady
  pubkey: onomypub1addwnpepqwtm4uklvv3t2snr2w3sxhhyrh867dj95n5yzcme7ltyzgzwhxk0gvtsdnv
  mnemonic: ""
  threshold: 0
  pubkeys: []
```

### Create your validator
To see the options when creating a validator:

```bash:
onomyd tx staking create-validator -h
```



### Update toml configuration files

Change in $HOME/onomy-testnet1/onomy/config/config.toml to contain the following:

```

replace seeds = "" to seeds = "5e0f5b9d54d3e038623ddb77c0b91b559ff13495@147.182.190.16:26656"
replace "tcp://127.0.0.1:26657" to "tcp://0.0.0.0:26657"
replace "tcp://127.0.0.1:26656" to "tcp://0.0.0.0:26656"
replace addr_book_strict = true to addr_book_strict = false
replace external_address = "" to external_address = "tcp://0.0.0.0:26656"
```

Change in $HOME/onomy-testnet1/onomy/config/app.toml to contain the following:

```
replace enable = false to enable = true
replace swagger = false to swagger = true
```

### Create your validator
To see the options when creating a validator:

```bash:
onomyd tx staking create-validator -h
```

An example of creating a validator with 50NOM self-delegation and 10% commission:
```bash:
# Replace <key_name> with the key you created previously
onomyd tx staking create-validator \
--amount=50000000unom \
--pubkey=$(onomyd tendermint show-validator) \
--moniker="choose moniker" \
--website="optional website for your validator"
--details="optional details for your validator"
--commission-rate="0.10" \
--commission-max-rate="0.20" \
--commission-max-change-rate="0.01" \
--min-self-delegation="1" \
--from=<key_name> \
--chain-id=$CHAIN_ID \
--gas=auto
--gas-adjustment=1.4
--keyring-backend test
```

To check on the status of your validator:
```bash:
onomyd status --output json | jq '.validator_info'
```

### Confirm that you are validating

If you see one line in the response you are validating. If you don't see any output from this command you are not validating. Check that the last command ran successfully.

Be sure to replace 'my validator key name' with your actual key name. If you want to double check you can see all your keys with 'onomyd keys list --home $HOME/$CHAIN_ID/onomy --keyring-backend test'

```

onomyd query staking validator $(onomyd --home $HOME/$CHAIN_ID/onomy keys show <myvalidatorkeyname> --bech val --address --keyring-backend test) --home $HOME/$CHAIN_ID/onomy

```

## Next Steps

Now that your validator is up and running, you can set up your [gravity bridge](onomy-testnet-docs/settings-up-gravity-relay.md).