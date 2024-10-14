/Users/donglieu/102024/onomy/script/multinode.sh

# set price
sleep 7

onomyd tx oracle set-price nomUSD 1 --home=$HOME/.onomyd/validator1  --from validator1 --keyring-backend test --fees 20stake --chain-id testing-1 -y
onomyd tx oracle set-price fet 1.3 --home=$HOME/.onomyd/validator2  --from validator2 --keyring-backend test --fees 20stake --chain-id testing-1 -y
onomyd tx oracle set-price atom 8.13 --home=$HOME/.onomyd/validator3  --from validator3 --keyring-backend test --fees 20stake --chain-id testing-1 -y

# submit proposal add usdt
sleep 7
onomyd q oracle  get-price fet
onomyd tx gov submit-proposal ./script/proposal-2.json --home=$HOME/.onomyd/validator1  --from validator1 --keyring-backend test --fees 20stake --chain-id testing-1 -y
onomyd tx oracle set-price usdt 1 --home=$HOME/.onomyd/validator3  --from validator3 --keyring-backend test --fees 20stake --chain-id testing-1 -y

# # vote
sleep 7
onomyd tx gov vote 1 yes  --from validator1 --keyring-backend test --home ~/.onomyd/validator1 --chain-id testing-1 -y --fees 20stake
onomyd tx gov vote 1 yes  --from validator2 --keyring-backend test --home ~/.onomyd/validator2 --chain-id testing-1 -y --fees 20stake
onomyd tx gov vote 1 yes  --from validator3 --keyring-backend test --home ~/.onomyd/validator3 --chain-id testing-1 -y --fees 20stake

# wait voting_perio=15s
sleep 15
echo "========sleep=========="

# check add usdt, balances
onomyd q psm  all-stablecoin
onomyd q bank balances $(onomyd keys show validator1 -a --keyring-backend test --home /Users/donglieu/.onomyd/validator1)

# tx swap usdt to nomUSD
echo "========swap==========="
onomyd tx psm swap-to-nomUSD 100000000000000000000000000000usdt --from validator1 --keyring-backend test --home ~/.onomyd/validator1 --chain-id testing-1 -y --fees 20stake

sleep 7
# Check account after swap
onomyd q bank balances $(onomyd keys show validator1 -a --keyring-backend test --home /Users/donglieu/.onomyd/validator1)

# tx swap nomUSD to usdt
onomyd tx psm swap-to-stablecoin usdt 1000nomUSD --from validator1 --keyring-backend test --home ~/.onomyd/validator1 --chain-id testing-1 -y --fees 20stake

sleep 7
# Check account after swap
onomyd q bank balances $(onomyd keys show validator1 -a --keyring-backend test --home /Users/donglieu/.onomyd/validator1)

test1="today auto lazy finger shoulder abstract oppose south sunny glass similar great feature rhythm raise evil owner orange auction absurd half mail ice glory"
echo $test1 | onomyd keys add test1 --recover --keyring-backend=test --home=$HOME/.onomyd/validator1

onomyd tx bank send $( onomyd keys show validator1 --home=$HOME/.onomyd/validator1  --keyring-backend test -a) $(onomyd keys show test1 --home=$HOME/.onomyd/validator1  --keyring-backend test -a) 20000000nomUSD,10000stake,10000000usdt --home=$HOME/.onomyd/validator1 --keyring-backend test --fees 20stake --chain-id testing-1 -y

# killall onomyd || true

sleep 7

onomyd tx oracle set-price usdt 1.1 --home=$HOME/.onomyd/validator3  --from validator3 --keyring-backend test --fees 20stake --chain-id testing-1 -y
# onomyd q gov proposals
# onomyd tx gov submit-legacy-proposal active-collateral "title" "description" "atom" "10" "0.1" "10000" 10000000000000000000stake --keyring-backend=test  --home=$HOME/.onomyd/validator1 --from validator1 -y --chain-id testing-1 --fees 20stake

# onomyd tx gov submit-proposal ./script/proposal-1.json --home=$HOME/.onomyd/validator1  --from validator1 --keyring-backend test --fees 20stake --chain-id testing-1 -y

# # vote
# sleep 7
# onomyd tx gov vote 2 yes  --from validator1 --keyring-backend test --home ~/.onomyd/validator1 --chain-id testing-1 -y --fees 20stake
# onomyd tx gov vote 2 yes  --from validator2 --keyring-backend test --home ~/.onomyd/validator2 --chain-id testing-1 -y --fees 20stake
# onomyd tx gov vote 2 yes  --from validator3 --keyring-backend test --home ~/.onomyd/validator3 --chain-id testing-1 -y --fees 20stake

# # wait voting_perio=15s
# echo "========sleep=========="
# sleep 15
# onomyd q gov proposals

# killall onomyd || true

# onomyd tx vaults create-vault 10000000atom 20000000nomUSD --from validator1 --home=$HOME/.onomyd/validator1 --keyring-backend test --fees 20stake --chain-id testing-1 -y

# sleep 7
# onomyd tx vaults create-vault 10000000atom 20000000nomUSD --from validator2 --home=$HOME/.onomyd/validator2 --keyring-backend test --fees 20stake --chain-id testing-1 -y

# sleep 7

# onomyd tx vaults mint 0 20000000nomUSD --from validator1 --home=$HOME/.onomyd/validator1 --keyring-backend test --fees 20stake --chain-id testing-1 -y

# sleep 7 

# onomyd tx vaults deposit 0 20000000atom --from validator1 --home=$HOME/.onomyd/validator1 --keyring-backend test --fees 20stake --chain-id testing-1 -y

# sleep 7

# onomyd tx vaults withdraw 0 1000atom --from validator1 --home=$HOME/.onomyd/validator1 --keyring-backend test --fees 20stake --chain-id testing-1 -y


# sleep 7

# onomyd tx vaults repay 0 10000000000000000nomUSD --from validator1 --home=$HOME/.onomyd/validator1 --keyring-backend test --fees 20stake --chain-id testing-1 -y

# sleep 7

# onomyd tx vaults close 0 --from validator1 --home=$HOME/.onomyd/validator1 --keyring-backend test --fees 20stake --chain-id testing-1 -y

# sleep 7

# onomyd q bank balances $(onomyd keys show validator1 -a --keyring-backend test --home /Users/donglieu/.onomyd/validator1)

# killall onomyd || true