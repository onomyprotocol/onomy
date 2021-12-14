# Bonding Curve Offering (BCO)

## Definition

Bonding curves are cryptoeconomic token models that automate the relationship between price and supply. The tokens in this model are referred to as Continuous Tokens because their price is continuously calculated. In continuous token models, there is no ICO or token launch. Instead of pre-selling tokens during a launch phase, the tokens are minted continuously over time via an automated market maker contract. Tokens are minted when purchased as needed, in conjunction with demand, and used within the protocol or application when required or desired.

Continuous Tokens have other properties such as instant liquidity and deterministic price. Bonding curves act as an automated market maker such that token buyers and sellers have an instant market. Additionally, bonding curve models don’t have central authorities responsible for issuing the tokens. Instead, users can buy a project’s token through a smart contract platform. The cost to buy these tokens is determined by the supply of those tokens. Unlike traditional models, the cost of these tokens increases as the supply increases. This price is determined by a pre-existing algorithm, further described below. A fee of 1% is applied per trade with the bonding curve.

## wNOM to NOM Distribution

The Onomy Protocol token, _**NOM**_, is primarily distributed through a bonding curve contract deployed on the Ethereum Network that issues wrapped tokens (_**wNOM**_) that will be exchanged 1:1 with _NOM_ on the Onomy Network. _NOM_ is the functional utility token of the Onomy Network. An interface will be provided for this exchange, commonly referenced as "bridging" _wNOM_ tokens from Ethereum and into the Onomy ecosystem for _NOM_. When bridged, _wNOM_ is burned from the bonding curve, setting a new floor price of the bonding curve, and a proportionate supply is issued on the Onomy Network.

The following is the equation governing the bonding curve:

![image](https://user-images.githubusercontent.com/76499838/145861419-06317db9-5450-495e-9a4b-289312829967.png)

100,000,000 is the supply of _wNOM_ loaded into the bonding curve. Supply is issued as it is purchased, thereby acting as a minting mechanism over time.

![image](https://user-images.githubusercontent.com/76499838/145861439-25e79b33-cfc6-4894-b973-f509cbf1e3e3.png)

## Distribution Benefits

The bonding curve has many benefits as a way to distribute _NOM_:

* **Deterministic Price:** The buy and sell prices of tokens increase and decrease with the number of tokens minted.&#x20;
* **Continuous Price:** The price of token _n_ is less than token _n+1_ and more than _n-1._
* **Instant Liquidity:** Tokens can be bought or sold instantaneously at any time, with the bonding curve acting as an automated market maker. A bonding curve contract acts as the counterparty of the transaction and always holds enough ETH in reserve to buy tokens back.&#x20;
* **Collateralization:** In conjunction with the staking inflation curve, ETH acts as reserve backing of the _NOM_.

## Staking & Bridge Incentives

The Onomy Network Staking Curve carries a tight relationship with the Bonding Curve. The staking curve acts as a principal motivator for participants who purchased _wNOM_ to bridge from Ethereum to the Onomy Network to capture the opportunity to earn staking rewards. _NOM_ holders can either choose to become validators themselves, or delegate their _NOM_ to a validator to earn staking rewards minus a validator fee.

![image](https://user-images.githubusercontent.com/76499838/145861495-9667f434-6c22-4361-a3b5-98abddcf1c9c.png)

The staking rewards incentivize _wNOM_ on the bonding curve to bridge into _NOM,_ thereby driving value into the Onomy Ecosystem. Once bridged, bonding curve participants formalize their entry into Onomy and relinquish the ability to sell back to the bonding curve. Liquidity thereby comes from the ONEX and anticipated exchange listings.

The bonding curve is perpetual until all 100M w_NOM_ are issued. Thus, an arbitrage opportunity may appear from time to time as the _NOM_ price on exchanges exceeds that of _wNOM_ on the bonding curve. A potential arbitrager will purchase _wNOM_, bridge to _NOM_, and send to exchanges. This incentive will further drive positive feedback mechanisms driving value, community, and adoption of _NOM,_ as more _wNOM_ is minted at more valuable prices and then bridged.

## Legal Disclaimer

_Nothing in this post shall constitute or be construed as an offer to sell or the solicitation of an offer to purchase any securities. Nothing in this post shall be construed as investment advice, strategy, or investment recommendations by Onomy or any of its affiliates. This post is for informational purposes only._

_This communication contains forward-looking statements that are based on Onomy’s beliefs and assumptions based on information currently available to Onomy. In some cases, you can identify forward-looking statements by the following words: “will,” “expect,” “would,” “intend,” “believe,” "anticipate," or other comparable terminology.These statements involve risks, uncertainties, assumptions, and other factors that may cause actual results or performance to be materially different. Onomy cannot assure you that the forward-looking statements will prove to be accurate. These forward-looking statements speak only as of the date hereof. We disclaim any obligation to update these forward-looking statements._

## Risk Disclosure

_The Bonding Curve Automated Market Maker is a decentralized protocol that people can use to create liquidity and trade wNOM tokens. Your use of the Bonding Curve Contract and its Platform involves various risks, including, but not limited to, losses due to the fluctuation of prices of tokens in a trading pair or smart contract vulnerabilities. You are responsible for doing your own diligence on those interfaces to understand the fees and risks they present._

_THE BONDING CURVE CONTRACT AND PLATFORM IS PROVIDED "AS IS", AT YOUR OWN RISK, AND WITHOUT WARRANTIES OF ANY KIND. Although NOM LABS developed the code, it does not provide, own, or control the Bonding Curve Contract or Platform, which is run by smart contracts deployed on the Ethereum blockchain. No developer or entity involved in creating the Bonding Curve Contract or Platform will be liable for any claims or damages whatsoever associated with your use, inability to use, or your interaction with other users of, the Bonding Curve Contract or Platform, including any direct, indirect, incidental, special, exemplary, punitive or consequential damages, or loss of profits, cryptocurrencies, tokens, or anything else of value._ 
