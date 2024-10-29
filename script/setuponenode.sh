#!/bin/bash
set -xeu
# set price
sleep 7
onomyd tx oracle set-price nomUSD 1 --home=$HOME/.onomy  --from val --keyring-backend test --fees 20stake --chain-id onomyd-1 -y
onomyd tx oracle set-price atom 8.0 --home=$HOME/.onomy  --from val2 --keyring-backend test --fees 20stake --chain-id onomyd-1 -y
sleep 7
onomyd q oracle  get-price fet
onomyd tx gov submit-proposal ./script/proposal-1.json --home=$HOME/.onomy  --from val --keyring-backend test --fees 20stake --chain-id onomyd-1 -y
# vote∑
sleep 7
onomyd tx gov vote 1 yes  --from val --keyring-backend test --home ~/.onomy --chain-id onomyd-1 -y --fees 20stake
onomyd tx gov vote 1 yes  --from val2 --keyring-backend test --home ~/.onomy --chain-id onomyd-1 -y --fees 20stake

# wait voting_perio=15s
echo "========sleep=========="
sleep 7
sleep 8
onomyd q gov proposals


onomyd tx vaults create-vault 12500000atom 50000000nomUSD --from val2 --home=$HOME/.onomy --keyring-backend test --fees 20stake --chain-id onomyd-1 -y