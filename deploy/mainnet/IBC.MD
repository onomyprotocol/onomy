# IBC

The doc describes the IBC integration for the mainnet

## "onomy-mainnet-1" -> "osmosis-1":

* Relayer address: onomy1w7ndwly5x9zz349g7ql4qy9qpzr5qhwf3mlf23
* ClientId: 07-tendermint-0 
* ConnectionId: connection-1
* PortId: transfer
* ChannelId: channel-0
* NOM denom on the osmosis: IBC/B9606D347599F0F2FDF82BA3EE339000673B7D274EA50F59494DC51EFCD42163

## "osmosis-1" -> "onomy-mainnet-1":

* Relayer address: osmo1w7ndwly5x9zz349g7ql4qy9qpzr5qhwfrpc0dx
* ClientId: 07-tendermint-2405
* ConnectionId: connection-1892
* PortId: transfer
* ChannelId: channel-525


## "onomy-mainnet-1" -> "cosmoshub-4":

* Relayer address: onomy1w7ndwly5x9zz349g7ql4qy9qpzr5qhwf3mlf23
* ClientId: 07-tendermint-1
* ConnectionId: connection-2
* PortId: transfer
* ChannelId: channel-1
* NOM denom on the cosmoshub: ibc/60E492D8249D7B1E28907E6EB884E16B2E467BC68E9BA37CE862BC11DA63488C

## "cosmoshub-4" -> "onomy-mainnet-1":

* Relayer address: cosmos1w7ndwly5x9zz349g7ql4qy9qpzr5qhwft6tlm5
* ClientId: 07-tendermint-1023
* ConnectionId: connection-714
* PortId: transfer
* ChannelId: channel-492


## "onomy-mainnet-1" -> "exchain-66"(OKeX Chain):

* Relayer address: onomy1w7ndwly5x9zz349g7ql4qy9qpzr5qhwf3mlf23
* ClientId: 07-tendermint-4
* ConnectionId: connection-6
* PortId: transfer
* ChannelId: channel-5
* NOM denom on the okc: ibc/e56506a75f79e066f9e0daf6e6786e5cfca1a5bbe6761c8a06c846ebe923b6ee
* NOM ERC20 address on okc: 0xA27D78c2707dC7Ba66d8BFa5CBae13588fFE6922

## "exchain-66"(OKeX Chain) -> "onomy-mainnet-1":

* Relayer address (cosmos): ex1nx6hsyyv2g32kd29gygvykdkev49y8dzrt0sq4
* Relayer address (EVM): 0x99B578108C5222ab35454110C259B6Cb2a521da2 
* ClientId: 0t-tendermint-118
* ConnectionId: connection-82
* PortId: transfer
* ChannelId: channel-50



## "onomy-mainnet-1" -> "injective-`":

* Relayer address: onomy1w7ndwly5x9zz349g7ql4qy9qpzr5qhwf3mlf23
* ClientId: 07-tendermint-3
* ConnectionId: connection-4
* PortId: transfer
* ChannelId: channel-3
* NOM denom on the injective: 

## "injective-4" -> "onomy-mainnet-1":

* Relayer address: inj1w7ndwly5x9zz349g7ql4qy9qpzr5qhwfpnumfv
* ClientId: 07-tendermint-141
* ConnectionId: connection-114
* PortId: transfer
* ChannelId: channel-96

Examples how to send:

Open onomy node CLI and execute transaction

```bash
onomyd tx ibc-transfer transfer transfer channel-0 my-osmo-address 1000anom --from my-accout --gas auto --gas-adjustment 1.3 -y -b block --keyring-backend pass --chain-id onomy-mainnet-1
```

Open osmosis node CLI and execute transaction

```bash
osmosisd tx ibc-transfer transfer transfer channel-525 my-onomy-address 1000osmo --from my-accout --gas auto --gas-adjustment 1.3 -y -b block --keyring-backend pass --chain-id osmosis-1
```
