# Setting up your Validator

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
```

To check on the status of your validator:
```bash:
onomyd status --output json | jq '.validator_info'
```

After you have completed this guide, your validator should be up and ready to receive delegations.