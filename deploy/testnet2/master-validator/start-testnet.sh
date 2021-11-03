#!/bin/bash

CURRENT_WORKING_DIR=$(pwd)
CHAINID="onomy-testnet2"

CHAINDIR="$CURRENT_WORKING_DIR/testdata"
ONOMY=onomyd
ONOMY_HOST="0.0.0.0"


home_dir="$CHAINDIR/$CHAINID"
n0name="onomy0"
n1name="onomy1"
n2name="onomy2"

# Folders for nodes
n0dir="$home_dir/$n0name"
n1dir="$home_dir/$n1name"
n2dir="$home_dir/$n2name"

# Home flag for folder
home0="--home $n0dir"
home1="--home $n1dir"
home2="--home $n2dir"

# Config directories for nodes
n0cfgDir="$n0dir/config"
n1cfgDir="$n1dir/config"
n2cfgDir="$n2dir/config"

# Config files for nodes
n0cfg="$n0cfgDir/config.toml"
n1cfg="$n1cfgDir/config.toml"
n2cfg="$n2cfgDir/config.toml"

# App config files for nodes
n0appCfg="$n0cfgDir/app.toml"
n1appCfg="$n1cfgDir/app.toml"
n2appCfg="$n2cfgDir/app.toml"

# Common flags
kbt="--keyring-backend test"
cid="--chain-id $CHAINID"

echo "Initializing genesis files"
STAKE_DENOM="nom"
NORMAL_DENOM="footoken"

coins="1000000000000$STAKE_DENOM,1000000000000$NORMAL_DENOM"

# Switch sed command in the case of linux
fsed() {
  if [ `uname` = 'Linux' ]; then
    sed -i "$@"
  else
    sed -i '' "$@"
  fi
}

# Initialize the 3 home directories and add some keys
$ONOMY $home0 $cid init n0 &>/dev/null
$ONOMY $home0 keys add val $kbt --output json | jq . >> $n0dir/validator_key.json
$ONOMY $home0 keys add orch $kbt --output json | jq . >> $n0dir/orchestrator_key.json
$ONOMY $home0 keys add faucet_account1 $kbt --output json | jq . >> $n0dir/faucet_account1.json
$ONOMY $home0 keys add faucet_account2 $kbt --output json | jq . >> $n0dir/faucet_account2.json
$ONOMY $home0 keys add faucet_account3 $kbt --output json | jq . >> $n0dir/faucet_account3.json
$ONOMY $home0 keys add faucet_account4 $kbt --output json | jq . >> $n0dir/faucet_account4.json
$ONOMY $home0 keys add faucet_account5 $kbt --output json | jq . >> $n0dir/faucet_account5.json

$ONOMY $home1 $cid init n1 &>/dev/null
$ONOMY $home1 keys add val $kbt --output json | jq . >> $n1dir/validator_key.json
$ONOMY $home1 keys add orch $kbt --output json | jq . >> $n1dir/orchestrator_key.json

$ONOMY $home2 $cid init n2 &>/dev/null
$ONOMY $home2 keys add val $kbt --output json | jq . >> $n2dir/validator_key.json
$ONOMY $home2 keys add orch $kbt --output json | jq . >> $n2dir/orchestrator_key.json

echo "Set stake/mint demon to $STAKE_DENOM"
fsed "s#\"stake\"#\"$STAKE_DENOM\"#g" $n0cfgDir/genesis.json

# add in denom metadata for both native tokens
jq '.app_state.bank.denom_metadata += [{"base": "footoken", display: "mfootoken", "description": "A non-staking test token", "denom_units": [{"denom": "footoken", "exponent": 0}, {"denom": "mfootoken", "exponent": 6}]}, {"base": "nom", display: "mnom", "description": "A staking test token", "denom_units": [{"denom": "nom", "exponent": 0}, {"denom": "mnom", "exponent": 6}]}]' $n0cfgDir/genesis.json > $CURRENT_WORKING_DIR/metadata-genesis.json

# a 60 second voting period to allow us to pass governance proposals in the tests
jq '.app_state.gov.voting_params.voting_period = "60s"' $CURRENT_WORKING_DIR/metadata-genesis.json > $CURRENT_WORKING_DIR/edited-genesis.json
mv $CURRENT_WORKING_DIR/edited-genesis.json $CURRENT_WORKING_DIR/genesis.json
mv $CURRENT_WORKING_DIR/genesis.json $n0cfgDir/genesis.json

echo "Adding validator addresses to genesis files"
$ONOMY $home0 add-genesis-account "$($ONOMY $home0 keys show val -a $kbt)" $coins &>/dev/null
$ONOMY $home0 add-genesis-account "$($ONOMY $home1 keys show val -a $kbt)" $coins &>/dev/null
$ONOMY $home0 add-genesis-account "$($ONOMY $home2 keys show val -a $kbt)" $coins &>/dev/null

echo "Adding orchestrator addresses to genesis files"
$ONOMY $home0 add-genesis-account "$($ONOMY $home0 keys show orch -a $kbt)" $coins &>/dev/null
$ONOMY $home0 add-genesis-account "$($ONOMY $home1 keys show orch -a $kbt)" $coins &>/dev/null
$ONOMY $home0 add-genesis-account "$($ONOMY $home2 keys show orch -a $kbt)" $coins &>/dev/null

echo "Adding faucet account addresses to genesis files"
$ONOMY $home0 add-genesis-account "$($ONOMY $home0 keys show faucet_account1 -a $kbt)" $coins &>/dev/null
$ONOMY $home0 add-genesis-account "$($ONOMY $home0 keys show faucet_account2 -a $kbt)" $coins &>/dev/null
$ONOMY $home0 add-genesis-account "$($ONOMY $home0 keys show faucet_account3 -a $kbt)" $coins &>/dev/null
$ONOMY $home0 add-genesis-account "$($ONOMY $home0 keys show faucet_account4 -a $kbt)" $coins &>/dev/null
$ONOMY $home0 add-genesis-account "$($ONOMY $home0 keys show faucet_account5 -a $kbt)" $coins &>/dev/null

echo "Copying genesis file around to sign"
cp $n0cfgDir/genesis.json $n1cfgDir/genesis.json
cp $n0cfgDir/genesis.json $n2cfgDir/genesis.json

echo "Generating ethereum keys"
$ONOMY $home0 eth_keys add --output=json | jq . >> $n0dir/eth_key.json
$ONOMY $home1 eth_keys add --output=json | jq . >> $n1dir/eth_key.json
$ONOMY $home2 eth_keys add --output=json | jq . >> $n2dir/eth_key.json

echo "Creating gentxs"
$ONOMY $home0 gentx --ip $ONOMY_HOST val 1000000000000$STAKE_DENOM "$(jq -r .address $n0dir/eth_key.json)" "$(jq -r .address $n0dir/orchestrator_key.json)" $kbt $cid &>/dev/null
$ONOMY $home1 gentx --ip $ONOMY_HOST val 1000000000000$STAKE_DENOM "$(jq -r .address $n1dir/eth_key.json)" "$(jq -r .address $n1dir/orchestrator_key.json)" $kbt $cid &>/dev/null
$ONOMY $home2 gentx --ip $ONOMY_HOST val 1000000000000$STAKE_DENOM "$(jq -r .address $n2dir/eth_key.json)" "$(jq -r .address $n2dir/orchestrator_key.json)" $kbt $cid &>/dev/null

echo "Collecting gentxs in $n0name"
cp $n1cfgDir/gentx/*.json $n0cfgDir/gentx/
cp $n2cfgDir/gentx/*.json $n0cfgDir/gentx/
$ONOMY $home0 collect-gentxs &>/dev/null

echo "Distributing genesis file into $n1name, $n2name"
cp $n0cfgDir/genesis.json $n1cfgDir/genesis.json
cp $n0cfgDir/genesis.json $n2cfgDir/genesis.json

# Change ports on n0 val
fsed "s#\"tcp://127.0.0.1:26656\"#\"tcp://0.0.0.0:26656\"#g" $n0cfg
fsed "s#\"tcp://127.0.0.1:26657\"#\"tcp://0.0.0.0:26657\"#g" $n0cfg
fsed 's#addr_book_strict = true#addr_book_strict = false#g' $n0cfg
fsed 's#external_address = ""#external_address = "tcp://'$ONOMY_HOST:26656'"#g' $n0cfg
fsed 's#enable = false#enable = true#g' $n0appCfg
fsed 's#swagger = false#swagger = true#g' $n0appCfg

# Change ports on n1 val
fsed "s#\"tcp://127.0.0.1:26656\"#\"tcp://0.0.0.0:26656\"#g" $n1cfg
fsed "s#\"tcp://127.0.0.1:26657\"#\"tcp://0.0.0.0:26657\"#g" $n1cfg
fsed 's#log_level = "main:info,state:info,statesync:info,*:error"#log_level = "info"#g' $n1cfg
fsed 's#addr_book_strict = true#addr_book_strict = false#g' $n1cfg
fsed 's#external_address = ""#external_address = "tcp://'$ONOMY_HOST':26656"#g' $n1cfg
fsed 's#enable = false#enable = true#g' $n1appCfg

# Change ports on n2 val
fsed "s#\"tcp://127.0.0.1:26656\"#\"tcp://0.0.0.0:26656\"#g" $n2cfg
fsed "s#\"tcp://127.0.0.1:26657\"#\"tcp://0.0.0.0:26657\"#g" $n2cfg
fsed 's#addr_book_strict = true#addr_book_strict = false#g' $n2cfg
fsed 's#external_address = ""#external_address = "tcp://'$ONOMY_HOST':26656"#g' $n2cfg
fsed 's#log_level = "main:info,state:info,statesync:info,*:error"#log_level = "info"#g' $n2cfg
fsed 's#enable = false#enable = true#g' $n2appCfg
