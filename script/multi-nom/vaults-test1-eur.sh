onomyd tx vaults create-vault 12500000atom 50000000nomEUR --from test2 --home=$HOME/.onomyd/validator1 --keyring-backend test --fees 20stake --chain-id testing-1 -y
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

onomyd tx auction bid 0 15000000nomEUR 6.0 --from test1 --home=$HOME/.onomyd/validator1 --keyring-backend test --fees 20stake --chain-id testing-1 -y

sleep 7

echo "I will buy all the remaining"
onomyd tx auction bid 0 60000000nomEUR 5.0 --from test3 --home=$HOME/.onomyd/validator1 --keyring-backend test --fees 20stake --chain-id testing-1 -y

echo "wating long time, query auction ratecurrent = 1.1...liquidate"
# onomyd tx aution 