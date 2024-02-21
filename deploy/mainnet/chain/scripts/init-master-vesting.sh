#!/bin/bash
set -eu
echo "Adding vesting account to the genesys file"
# Vesting accounts
ONOMY_VESTING_ACCOUNTS=$(cat ../genesis/accounts-vesting.txt)

while read line;
  do echo "Adding genesis vesting account: '${line}'";
  acc_arr=($line);
  # onomy1b40a9e9zel7hk30jm4h6n0wa35l7s9ufvnqyrt 1000anom 1640620000 1640630000 9000|100anom,900|800anom,100|100anom
  onomyd add-genesis-account "${acc_arr[0]}"  "${acc_arr[1]}" \
     --vesting-amount="${acc_arr[1]}" \
     --vesting-start-time="${acc_arr[2]}" \
     --vesting-end-time="${acc_arr[3]}" \
     --vesting-periods-amounts="${acc_arr[4]}" \
     --keyring-backend test
done <<< "$ONOMY_VESTING_ACCOUNTS"
