# Testing the Onomy Network through Gravity

Now that we've made it this far, it's time to actually play around with the bridge.

## Bounding curve

The first example of the Bridge usage was on the [Bonding Curve](bonding-curve.md) page, so if you haven't pass it yet,
better to do it.

## Bridging ERC20 tokens

In order to participate in the Bridging ERC20 tokens you need to run the [full node](full.md) and install the "gbt"
artifact.

### Choose any ERC20 token to bridge

You can get the ERC20 token address of the tokens you already have on your Ethereum balance or you can get/mint new
tokens for example [FAU](https://erc20faucet.com/). The address of the FAU on the goerli
is ```0xBA62BCfcAaFc6622853cca2BE6Ac7d845BC0f2Dc```.

### Choose or generate the onomy address

If you already have the onomy address, then you can use it, for the next step. If not then you can generate a new one.
Be sure that you have already set up the [pass](./pass.md).

```
onomyd keys add my-onomy-user --keyring-backend pass
```

Where the ```my-onomy-user``` is the name of the key stored in the keyring. Also, the command will print out
the ```my-onomy-user``` address and mnemonic. Copy the output for the following steps.

### Get a user's balance

Get curren balance of the user.

```
onomyd q bank balances <your_onomy_address>
```

### Send tokens from the Ethereum to onomy

```
gbt -a onomy client eth-to-cosmos \
        --ethereum-key "<your_ethereum_private_key>" \
        --ethereum-rpc "http://0.0.0.0:8545/" \
        --gravity-contract-address "0x7DD0eAe5c6dE6C5003F06C76f5E789125066217B" \
        --token-contract-address "<token_address>" \
        --amount <token_amount> \
        --destination "<your_onomy_address>"
```

Update the templated parameter to your parameters.

If you don't have the goerli Ethereum node running on ```http://0.0.0.0:8545/``` you can use Ethereum goerli RPC
provided by the [Alchemy](https://www.alchemy.com/).

If the command has passed successfully you will see the output link this:

```
[2021-12-13T11:24:53Z INFO  gbt::client::eth_to_cosmos] Sending 100000000000000000000 / 0xBA62BCfcAaFc6622853cca2BE6Ac7d845BC0f2Dc to Cosmos from 0x2d9480eBA3A001033a0B8c3Df26039FD3433D55d to onomy1f6savxl0mqt6mgg46tjvxp3py773ql45yf0ltq
[2021-12-13T11:25:22Z INFO  gbt::client::eth_to_cosmos] Send to Cosmos txid: 0xbdfee272c18417a3e2a5d86a78c336484befb1ec0ba45da37035fa3a680154bd
```

### Validate the bridging from the Ethereum to onomy

The bridging operation might take few minutes. Await 3-5 minutes and get the onomy user balances.

```
onomyd q bank balances <your_onomy_address>
```

The valid output is

```
balances:
- amount: "<token_amount>"
  denom: gravity<token_address>
pagination:
  next_key: null
  total: "0"
```

### Send tokens from the onomy to Ethereum

Once you received the ERC20 tokens you can send them back to your Ethereum account.

```
gbt -a onomy client cosmos-to-eth \
        --cosmos-grpc="http://0.0.0.0:9191" \
        --cosmos-phrase "<your_onomy_mnemonic>" \
        --fees 1<token_address> \
        --amount <token_amount><token_address> \
        --eth-destination "<your_ethtreum_address>"
```

Update the templated values with your values.

If the command has passed successfully you will see the output link this:

```
[2021-12-13T11:35:20Z INFO  gbt::client::cosmos_to_eth] Sending from Cosmos address onomy104n3g00hms7kycg54jml24s84jjt9s2agq6r07
[2021-12-13T11:35:20Z INFO  gbt::client::cosmos_to_eth] Asset gravity0xBA62BCfcAaFc6622853cca2BE6Ac7d845BC0f2Dc has ERC20 representation 0xBA62BCfcAaFc6622853cca2BE6Ac7d845BC0f2Dc
[2021-12-13T11:35:20Z INFO  gbt::client::cosmos_to_eth] Cosmos balances [Coin { amount: Uint256(1000000000000000000), denom: "anom" }, Coin { amount: Uint256(20000000000000000000), denom: "gravity0xBA62BCfcAaFc6622853cca2BE6Ac7d845BC0f2Dc" }]
[2021-12-13T11:35:20Z INFO  gbt::client::cosmos_to_eth] Locking gravity0xBA62BCfcAaFc6622853cca2BE6Ac7d845BC0f2Dc / gravity0xBA62BCfcAaFc6622853cca2BE6Ac7d845BC0f2Dc into the batch pool
[2021-12-13T11:35:22Z INFO  gbt::client::cosmos_to_eth] Send to Eth txid 8A43402DD96C0A56734A3F488CCC9564D29C18E1D4D1E6C5A0CC886AEBD4D3E4
[2021-12-13T11:35:22Z INFO  gbt::client::cosmos_to_eth] Requesting a batch to push transaction along immediately
```

Then await 3-5 minutes and get the balance on the Ethereum side.