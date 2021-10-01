# Requesting tokens from the faucet

This assumes you've already set up a full node and are fully caught up.

### Request tokens from the faucet

First list all of your keys

```
onomyd keys list --keyring-backend test --home $HOME/$CHAIN_ID/onomy
```

You'll see an output like this

```
- name: lav.test
  type: local
  address: onomy1anlymr6kr7ns6t3c0lhjcvse7aqegffwm3vady
  pubkey: onomypub1addwnpepqwtm4uklvv3t2snr2w3sxhhyrh867dj95n5yzcme7ltyzgzwhxk0gvtsdnv
  mnemonic: ""
  threshold: 0
  pubkeys: []

```

Copy your address from the 'address' field and paste it into the command below in place of $ACCOUNT_ADDRESS

```
curl -X POST http://147.182.190.16:8000/ -H  "accept: application/json" -H  "Content-Type: application/json" -d "{  \"address\": \"$ACCOUNT_ADDRESS\",  \"coins\": [    \"10nom\"  ]}"
```

This will provide you 10nom from the faucet storage.
You can check balances of account using this command as given below
```
onomyd query bank balances $ACCOUNT_ADDRESS --home $HOME/$CHAIN_ID/onomy
```
