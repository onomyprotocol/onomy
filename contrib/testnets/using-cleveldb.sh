#!/bin/bash

make install GAIA_BUILD_OPTIONS="cleveldb"

onomyd init "t6" --home ./t6 --chain-id t6

onomyd unsafe-reset-all --home ./t6

mkdir -p ./t6/data/snapshots/metadata.db

onomyd keys add validator --keyring-backend test --home ./t6

onomyd add-genesis-account $(onomyd keys show validator -a --keyring-backend test --home ./t6) 100000000stake --keyring-backend test --home ./t6

onomyd gentx validator 100000000stake --keyring-backend test --home ./t6 --chain-id t6

onomyd collect-gentxs --home ./t6

onomyd start --db_backend cleveldb --home ./t6
