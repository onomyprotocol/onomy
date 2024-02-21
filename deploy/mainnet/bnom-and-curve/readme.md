The bnom deployment contains both the nom and bonding-curve-curve contract deployments:

In order to deploy:

* Pull the "solonomy" repo

```
git pull https://github.com/onomyprotocol/solonomy.git
```

* Create 2 files in the root:

'.secret'

```
your mnemonic
```

'.env'

```
ETHERSCAN_API_KEY=<YOURE-ETHERSCAN-KEY>
```

* Install yarn and nodejs

* Run build scripts

```
yarn install 
yarn compile
```

* Deploy the contracts

```
yarn deploy:mainnet
```

The script will deploy the contracts and save the addresses in the 'compiled/chain-goerli-NOMAddrs.json' file.
