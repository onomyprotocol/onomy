# Delegation

## Bounding curve

In order to get the ANOMs/NOMs pass through the [Bonding Curve](bonding-curve.md) instruction.

## What is Delegation

Delegation is lending your tokens to one or more validators in order to take part in the consensus. Delegators share
both - revenue and risks with validators. Validators share their revenue with their delegators. Validators may apply a
commission, or a fee before they distribute the profits to the delegators as well. On the other hand, if a validator
misbehaves (Such as not signing enough blocks in a sliding window or double signing the blocks), validator's and their
delegator's funds might be slashed.

## Why Delegate

Not everyone is able to join the network as a validator, these people can still be part of the consensus as delegators
and get return on their investment.

## How to Delegate

In order to delegate, you need to send a delegate transaction.

```
onomyd tx staking delegate <validator-address> <amount to delegate> --from <account to delegate from> --chain-id=onomy-testnet --keyring-backend test
```