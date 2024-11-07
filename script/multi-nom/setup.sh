#!/bin/bash
# set -xeu
# # set price
# onomyd tx oracle set-price nomUSD 1 --home=$HOME/.onomyd/validator1  --from validator1 --keyring-backend test --fees 20stake --chain-id testing-1 -y
# onomyd tx oracle set-price fet 1.3 --home=$HOME/.onomyd/validator2  --from validator2 --keyring-backend test --fees 20stake --chain-id testing-1 -y
# onomyd tx oracle set-price atom 8.0 --home=$HOME/.onomyd/validator3  --from validator3 --keyring-backend test --fees 20stake --chain-id testing-1 -y
# onomyd tx oracle set-price nomEUR 1 --home=$HOME/.onomyd/validator1  --from validator1 --keyring-backend test --fees 20stake --chain-id testing-1 -y
sleep 7
# onomyd q oracle  price atom nomEUR
onomyd tx gov submit-proposal ./script/multi-nom/proposal-vault-1.json --home=$HOME/.onomyd/validator1  --from validator1 --keyring-backend test --fees 20stake --chain-id testing-1 -y
# onomyd tx gov submit-proposal ./script/multi-nom/proposal-vault-2.json --home=$HOME/.onomyd/validator2  --from validator2 --keyring-backend test --fees 20stake --chain-id testing-1 -y
# vote∑
sleep 7
onomyd tx gov vote 1 yes  --from validator1 --keyring-backend test --home ~/.onomyd/validator1 --chain-id testing-1 -y --fees 20stake
onomyd tx gov vote 1 yes  --from validator2 --keyring-backend test --home ~/.onomyd/validator2 --chain-id testing-1 -y --fees 20stake
onomyd tx gov vote 1 yes  --from validator3 --keyring-backend test --home ~/.onomyd/validator3 --chain-id testing-1 -y --fees 20stake
# sleep 7
# onomyd tx gov vote 2 yes  --from validator1 --keyring-backend test --home ~/.onomyd/validator1 --chain-id testing-1 -y --fees 20stake
# onomyd tx gov vote 2 yes  --from validator2 --keyring-backend test --home ~/.onomyd/validator2 --chain-id testing-1 -y --fees 20stake
# onomyd tx gov vote 2 yes  --from validator3 --keyring-backend test --home ~/.onomyd/validator3 --chain-id testing-1 -y --fees 20stake

# wait voting_perio=15s
echo "========sleep=========="
sleep 7
test1="today auto lazy finger shoulder abstract oppose south sunny glass similar great feature rhythm raise evil owner orange auction absurd half mail ice glory"
echo $test1 | onomyd keys add test1 --recover --keyring-backend=test --home=$HOME/.onomyd/validator1

test2="convince ocean tower relax toward cattle sort wonder cross enhance pull describe typical total link home glare polar clip trim fish divorce arrest fall"
echo $test2 | onomyd keys add test2 --recover --keyring-backend=test --home=$HOME/.onomyd/validator1

test3="famous aware lens hair relax cancel glove gloom enforce shoulder spread valley any uncover slush gain dawn slim pipe kidney come exit bench bomb"
echo $test3 | onomyd keys add test3 --recover --keyring-backend=test --home=$HOME/.onomyd/validator1

onomyd tx bank multi-send $( onomyd keys show validator1 --home=$HOME/.onomyd/validator1  --keyring-backend test -a) $( onomyd keys show test1 --home=$HOME/.onomyd/validator1  --keyring-backend test -a) $( onomyd keys show test2 --home=$HOME/.onomyd/validator1  --keyring-backend test -a) $( onomyd keys show test3 --home=$HOME/.onomyd/validator1  --keyring-backend test -a) 1000000000nomUSD,10000stake,100000000atom,1000000000nomEUR,10000000000000000000nomJPY --home=$HOME/.onomyd/validator1 --keyring-backend test --fees 250stake --chain-id testing-1 -y --gas 2036038
sleep 8
onomyd q gov proposals