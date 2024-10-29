#!/bin/bash
killall onomyd || true
rm -rf $HOME/.onomy/

# init chain
onomyd init test --chain-id onomyd-1

echo $(cat /Users/donglieu/script/keys/mnemonic1)| onomyd keys add val --keyring-backend test  --recover 
echo $(cat /Users/donglieu/script/keys/mnemonic2)| onomyd keys add val2 --keyring-backend test  --recover 

# onomyd genesis add-genesis-account val 1000000000000stake --keyring-backend test
# onomyd genesis add-genesis-account val2 1000000000000stake --keyring-backend test
onomyd genesis add-genesis-account val 1000000000000000000000000stake,1000000000000000000000000atom --keyring-backend test
onomyd genesis add-genesis-account val2 10000000000000000000stake,1000000000000000000000000atom --keyring-backend test


onomyd genesis gentx val  100000000000000000000000stake --keyring-backend test --chain-id onomyd-1
# # Collect genesis tx 
onomyd genesis collect-gentxs
# # Run this to ensure everything worked and that the genesis file is setup correctly
# onomyd validate
onomyd genesis validate-genesis

update_test_genesis () {
    cat $HOME/.onomy/config/genesis.json | jq "$1" > tmp.json && mv tmp.json $HOME/.onomy/config/genesis.json
}

update_test_genesis '.app_state["gov"]["params"]["voting_period"] = "15s"'
update_test_genesis '.app_state["gov"]["params"]["expedited_voting_period"] = "10s"'


sed -i -E 's|minimum-gas-prices = ""|minimum-gas-prices = "0.0001stake"|g' $HOME/.onomy/config/app.toml
# screen -S xionx -t xionx -d -m
onomyd start

# onomyvaloper1wa3u4knw74r598quvzydvca42qsmk6jrya79zd

# onomyd testnet onomy-mainnet-1  onomyvaloper1wa3u4knw74r598quvzydvca42qsmk6jrya79zd --accounts-to-fund="onomy1wa3u4knw74r598quvzydvca42qsmk6jrc6uj7m,onomy1w7f3xx7e75p4l7qdym5msqem9rd4dyc4y47xsd,onomy1g9v3zjt6rfkwm4s8sw9wu4jgz9me8pn2ygn94a" --home=$HOME/.onomy --skip-confirmation