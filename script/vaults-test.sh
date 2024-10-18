#!/bin/bash
set -xeu
# set price
sleep 7
onomyd tx oracle set-price nomUSD 1 --home=$HOME/.onomyd/validator1  --from validator1 --keyring-backend test --fees 20stake --chain-id testing-1 -y
onomyd tx oracle set-price fet 1.3 --home=$HOME/.onomyd/validator2  --from validator2 --keyring-backend test --fees 20stake --chain-id testing-1 -y
onomyd tx oracle set-price atom 8.0 --home=$HOME/.onomyd/validator3  --from validator3 --keyring-backend test --fees 20stake --chain-id testing-1 -y
sleep 7
onomyd q oracle  get-price fet
onomyd tx gov submit-proposal ./script/proposal-1.json --home=$HOME/.onomyd/validator1  --from validator1 --keyring-backend test --fees 20stake --chain-id testing-1 -y
# # vote
sleep 7
onomyd tx gov vote 1 yes  --from validator1 --keyring-backend test --home ~/.onomyd/validator1 --chain-id testing-1 -y --fees 20stake
onomyd tx gov vote 1 yes  --from validator2 --keyring-backend test --home ~/.onomyd/validator2 --chain-id testing-1 -y --fees 20stake
onomyd tx gov vote 1 yes  --from validator3 --keyring-backend test --home ~/.onomyd/validator3 --chain-id testing-1 -y --fees 20stake

# wait voting_perio=15s
echo "========sleep=========="
sleep 7
test1="today auto lazy finger shoulder abstract oppose south sunny glass similar great feature rhythm raise evil owner orange auction absurd half mail ice glory"
echo $test1 | onomyd keys add test1 --recover --keyring-backend=test --home=$HOME/.onomyd/validator1

test2="convince ocean tower relax toward cattle sort wonder cross enhance pull describe typical total link home glare polar clip trim fish divorce arrest fall"
echo $test2 | onomyd keys add test2 --recover --keyring-backend=test --home=$HOME/.onomyd/validator1

test3="famous aware lens hair relax cancel glove gloom enforce shoulder spread valley any uncover slush gain dawn slim pipe kidney come exit bench bomb"
echo $test3 | onomyd keys add test3 --recover --keyring-backend=test --home=$HOME/.onomyd/validator1

onomyd tx bank multi-send $( onomyd keys show validator1 --home=$HOME/.onomyd/validator1  --keyring-backend test -a) $( onomyd keys show test1 --home=$HOME/.onomyd/validator1  --keyring-backend test -a) $( onomyd keys show test2 --home=$HOME/.onomyd/validator1  --keyring-backend test -a) $( onomyd keys show test3 --home=$HOME/.onomyd/validator1  --keyring-backend test -a) 1000000000nomUSD,10000stake,100000000atom --home=$HOME/.onomyd/validator1 --keyring-backend test --fees 20stake --chain-id testing-1 -y
sleep 8
onomyd q gov proposals
onomyd tx vaults create-vault 12500000atom 50000000nomUSD --from test2 --home=$HOME/.onomyd/validator1 --keyring-backend test --fees 20stake --chain-id testing-1 -y
# onomyd tx vaults create-vault 10000000atom 20000000nomUSD --from validator2 --home=$HOME/.onomyd/validator2 --keyring-backend test --fees 20stake --chain-id testing-1 -y

# sleep 7

# onomyd tx vaults mint 0 20000000nomUSD --from validator1 --home=$HOME/.onomyd/validator1 --keyring-backend test --fees 20stake --chain-id testing-1 -y

# sleep 7 

# onomyd tx vaults deposit 0 20000000atom --from validator1 --home=$HOME/.onomyd/validator1 --keyring-backend test --fees 20stake --chain-id testing-1 -y

# sleep 7

# onomyd tx vaults withdraw 0 1000atom --from validator1 --home=$HOME/.onomyd/validator1 --keyring-backend test --fees 20stake --chain-id testing-1 -y

# sleep 7

# onomyd tx vaults repay 0 40000000nomUSD --from validator1 --home=$HOME/.onomyd/validator1 --keyring-backend test --fees 20stake --chain-id testing-1 -y

sleep 7

onomyd q bank balances $(onomyd keys show test2 -a --keyring-backend test --home $HOME/.onomyd/validator1)

onomyd tx oracle set-price atom 5 --home=$HOME/.onomyd/validator1  --from validator1 --keyring-backend test --fees 20stake --chain-id testing-1 -y

sleep 31
onomyd q bank balances $(onomyd keys show test1 --home=$HOME/.onomyd/validator1  --keyring-backend test -a)

onomyd tx auction bid 0 15000000nomUSD 6.0 --from test1 --home=$HOME/.onomyd/validator1 --keyring-backend test --fees 20stake --chain-id testing-1 -y

sleep 7

echo "I will buy all the remaining"
onomyd tx auction bid 0 60000000nomUSD 5.0 --from test3 --home=$HOME/.onomyd/validator1 --keyring-backend test --fees 20stake --chain-id testing-1 -y

echo "wating long time, query auction ratecurrent = 1.1...liquidate"
# onomyd tx aution 