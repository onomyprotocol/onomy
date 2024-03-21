#!/bin/bash

ONOMY_CHAIN_ID=onomy-mainnet-1
OSMO_CHAIN_ID=osmosis-1

sudo apt update && sudo apt upgrade -y
sudo apt install curl tar nano build-essential pkg-config libssl-dev tmux vim -y
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
# Press enter for default installation
source "$HOME/.cargo/env"
cargo install ibc-relayer-cli --bin hermes --locked
export PATH=$HOME/.cargo/bin:$PATH
hermes version
echo "export PATH=\$HOME/.cargo/bin:\$PATH" >> ~/.profile

mkdir -p $HOME/.hermes
# TODO: change URL once merged
curl https://raw.githubusercontent.com/onomyprotocol/onomy/scripts/deploy/mainnet/chain/scripts/assets/hermes/config.toml > $HOME/.hermes/config.toml
hermes keys add --chain $ONOMY_CHAIN_ID --mnemonic-file mnemonic-onomy.txt
hermes keys add --chain $OSMO_CHAIN_ID --mnemonic-file mnemonic-osmosis.txt
rm mnemonic-onomy.txt mnemonic-osmosis.txt
hermes keys balance --chain $ONOMY_CHAIN_ID
hermes keys balance --chain $OSMO_CHAIN_ID

echo "
[Unit]
Description=hermes
After=network.target

[Service]
User=ubuntu
ExecStart=$HOME/.cargo/bin/hermes start

Restart=always
RestartSec=3

[Install]
WantedBy=multi-user.target
" | sudo tee /etc/systemd/system/hermes.service

sudo systemctl daemon-reload
sudo systemctl start hermes.service
sudo systemctl enable hermes.service
sudo systemctl status hermes.service --no-pager
hermes health-check
sudo journalctl -u hermes.service -n 10 --no-pager

