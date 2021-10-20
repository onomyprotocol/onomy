# Testing the Onomy Network through Gravity

Now that we've made it this far, it's time to actually play around with the bridge.

This first command will send some ERC20 tokens to an address of your choice on the Onomy
Network. 

Notice that the Ethereum key is pre-filled. This address has both some test ETH and a large balance of ERC20 tokens from the contracts listed here:

```
0xFab46E002BbF0b4509813474841E0716E6730136 - FooToken (FAU)
0xd92e713d051c37ebb2561803a3b5fbabc4962431 - Test USDT (TUSDT)
```

Note that the 'amount' field for this command is now in whole coins rather than wei like in the previous testnets:

```
gbt -a onomy client eth-to-cosmos \
        --ethereum-key "0x07a78996067f8049a73cefbe18ed114a47fc0136e8cf57fb153917b79cdc0066" \
        --gravity-contract-address "0xB4BAd4Cef22a4EAeF67434644ebaB4cEC54Db37A" \
        --token-contract-address "0xFab46E002BbF0b4509813474841E0716E6730136" \
        --amount=100 \
        --destination "any Cosmos address, I suggest your delegate Cosmos address"
```

You should see a message like this on your Orchestrator. The details, of course, will be different, but this means that your Orchestrator has observed the event on Ethereum and has sent the details to the Cosmos chain!

```
[2021-02-13T12:35:54Z INFO  orchestrator::ethereum_event_watcher] Oracle observed deposit with sender 0xBf660843528035a5A4921534E156a27e64B231fE, destination cosmos1xpfu40gseet70wfeazds773v05pjx3dwe7e03f, amount
999999984306749440, and event nonce 3
```

Once the event has been observed, we can check our balance on the Cosmos side. We will see some peggy <ERC20 address> tokens in our balance. We have a lot of code in production right now, so the module renaming from 'Peggy' to 'Gravity' has been put on hold until production is complete.

```
onomyd --home $HOME/.onomy query bank balances <any cosmos address>
```

Now that we have some tokens on the Onomy Network, we can try sending them back to Ethereum. Remember to use the Cosmos phrase for the address you actually sent the tokens to. Alternately you can send Cosmos native tokens with this command.

The denom of a bridged token will be:

```
gravity0xFab46E002BbF0b4509813474841E0716E6730136
```

```
gbt -a onomy client cosmos-to-eth \
        --cosmos-phrase "the phrase containing the Gravity bridged tokens (delegate keys mnemonic)" \
        --fees 1nom \
        --amount 1000gravity0xFab46E002BbF0b4509813474841E0716E6730136 \
        --eth-destination "any eth address, try your delegate eth address"
```

It will take a moment or two for EtherScan to catch up, but once it has, you'll see the new ERC20 token balance reflected on EtherScan.
