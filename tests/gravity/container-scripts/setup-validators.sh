#!/bin/bash
set -eux
# your onomyd binary name
BIN=onomyd

CHAIN_ID="gravity-test"

NODES=$1

ALLOCATION="10000000000000000000000anom,10000000000footoken,10000000000ibc/nometadatatoken"

# first we start a genesis.json with validator 1
# validator 1 will also collect the gentx's once gnerated
STARTING_VALIDATOR=1
STARTING_VALIDATOR_HOME="--home /validator$STARTING_VALIDATOR"
$BIN init $STARTING_VALIDATOR_HOME --chain-id=$CHAIN_ID validator1

## Modify generated genesis.json to our liking by editing fields using jq
## we could keep a hardcoded genesis file around but that would prevent us from
## testing the generated one with the default values provided by the module.

# add in denom metadata for both native tokens
jq '.app_state.bank.denom_metadata += [{"name": "Foo Token", "symbol": "FOO", "base": "footoken", display: "mfootoken", "description": "A non-staking test token", "denom_units": [{"denom": "footoken", "exponent": 0}, {"denom": "mfootoken", "exponent": 6}]},{"name": "NOM", "symbol": "NOM", "base": "anom", display: "nom", "description": "Nom token", "denom_units": [{"denom": "anom", "exponent": 0}, {"denom": "nom", "exponent": 18}]}]' /validator$STARTING_VALIDATOR/config/genesis.json > /metadata-genesis.json

# a 60 second voting period to allow us to pass governance proposals in the tests
jq '.app_state.gov.voting_params.voting_period = "60s"' /metadata-genesis.json > /community-pool-genesis.json

# Add some funds to the community pool to test Airdrops, note that the gravity address here is the first 20 bytes
# of the sha256 hash of 'distribution' to create the address of the module
# To get from code: app.AccountKeeper.GetModuleAddress(distrtypes.ModuleName).String() // onomy1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8a7s2c6
jq '.app_state.distribution.fee_pool.community_pool = [{"denom": "anom", "amount": "10000000000.0"}]' /community-pool-genesis.json > /community-pool2-genesis.json
jq '.app_state.auth.accounts += [{"@type": "/cosmos.auth.v1beta1.ModuleAccount", "base_account": { "account_number": "0", "address": "onomy1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8a7s2c6","pub_key": null,"sequence": "0"},"name": "distribution","permissions": ["basic"]}]' /community-pool2-genesis.json > /community-pool3-genesis.json
jq '.app_state.bank.balances += [{"address": "onomy1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8a7s2c6", "coins": [{"amount": "10000000000", "denom": "anom"}]}]' /community-pool3-genesis.json > /community-pool4-genesis.json
jq '.app_state.bank.supply += [{"denom": "anom", "amount": "10000000000"}]' /community-pool4-genesis.json > /edited-genesis.json

# rename base denom to anom
sed -i 's/stake/anom/g' /edited-genesis.json
mv /edited-genesis.json /genesis.json

echo "The test genesis is ready:"
cat /genesis.json

# Sets up an arbitrary number of validators on a single machine by manipulating
# the --home parameter on onomyd
for i in $(seq 1 $NODES);
do
ONOMY_HOME="--home /validator$i"
ARGS="$ONOMY_HOME --keyring-backend test"

# Generate a validator key, orchestrator key, and eth key for each validator
$BIN keys add $ARGS validator$i 2>> /validator-phrases
$BIN keys add $ARGS orchestrator$i 2>> /orchestrator-phrases
$BIN eth_keys add >> /validator-eth-keys

VALIDATOR_KEY=$($BIN keys show validator$i -a $ARGS)
ORCHESTRATOR_KEY=$($BIN keys show orchestrator$i -a $ARGS)
# move the genesis in
mkdir -p /validator$i/config/
mv /genesis.json /validator$i/config/genesis.json
$BIN add-genesis-account $ARGS $VALIDATOR_KEY $ALLOCATION
$BIN add-genesis-account $ARGS $ORCHESTRATOR_KEY $ALLOCATION
# move the genesis back out
mv /validator$i/config/genesis.json /genesis.json
done


for i in $(seq 1 $NODES);
do
cp /genesis.json /validator$i/config/genesis.json
ONOMY_HOME="--home /validator$i"
ARGS="$ONOMY_HOME --keyring-backend test"
ORCHESTRATOR_KEY=$($BIN keys show orchestrator$i -a $ARGS)
ETHEREUM_KEY=$(grep address /validator-eth-keys | sed -n "$i"p | sed 's/.*://')
# the /8 containing 7.7.7.7 is assigned to the DOD and never routable on the public internet
# we're using it in private to prevent onomy from blacklisting it as unroutable
# and allow local pex
$BIN gravity gentx $ARGS $ONOMY_HOME --moniker validator$i --chain-id=$CHAIN_ID --ip 7.7.7.$i validator$i 500000000000000000000anom $ETHEREUM_KEY $ORCHESTRATOR_KEY
# obviously we don't need to copy validator1's gentx to itself
if [ $i -gt 1 ]; then
cp /validator$i/config/gentx/* /validator1/config/gentx/
fi
done


$BIN gravity collect-gentxs $STARTING_VALIDATOR_HOME
GENTXS=$(ls /validator1/config/gentx | wc -l)
cp /validator1/config/genesis.json /genesis.json
echo "Collected $GENTXS gentx"

# put the now final genesis.json into the correct folders
for i in $(seq 1 $NODES);
do
cp /genesis.json /validator$i/config/genesis.json
done
