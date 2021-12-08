# Validator Jailing/Unjailing
## What is Jailing

When a validator disconnects from the network due to connection loss or server fail or it double signs, it needs to be eliminated from the validator list. It is known as 'jailing'. A validator is jailed if it fails to validate at least 50% of the last 100 blocks.

When jailed due to downtime, the validator's total stake is slashed by 1%. and if jailed  due to double signing, validaor's total stake is slashed by 5%.

Once jailed, validators can be unjailed again after 10 minutes. 
These configurations can be found in the genesis file under the slashing section
```
"slashing": {
      "params": {
        "signed_blocks_window": "100",
        "min_signed_per_window": "0.500000000000000000",
        "downtime_jail_duration": "600s",
        "slash_fraction_double_sign": "0.050000000000000000",
        "slash_fraction_downtime": "0.010000000000000000"
      },
      "signing_infos": [],
      "missed_blocks": []
    }
```

## Unjailing validator
In order to unjail the validator, you may run the following command once 10 minutes have passed

```
onomyd tx slashing unjail --from <validator-name> --chain-id=onomy-testnet --keyring-backend test --gas-prices=0.025
```



