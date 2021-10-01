# Onomy setup full node

Setting up a node requires first installing [onomyd](onomy-testnet-docs/install-onomyd.md).

### Config setup
You'll want to set a few variables before continuing. 

Choose your `chainId` you would like to join. The current testnet is `onomy-testnet1`.
```
CHAIN_ID=<chainId>

# example
CHAIN_ID=onomy-testnet1
```

Choose you `moniker-name`. This can be anything of your choosing and will be used to identify you in the explorer.
```
MONIKER_NAME=<moniker-name>

# example
MONIKER_NAME="Powerman-5000"
```

### Init the config files
```
cd $HOME
onomyd init $MONIKER_NAME --chain-id $CHAIN_ID --home $HOME/$CHAIN_ID/onomy
```

### Copy the genesis file

```
rm $HOME/$CHAIN_ID/onomy/config/genesis.json
wget http://147.182.190.16:26657/genesis? -O $HOME/raw.json
jq .result.genesis $HOME/raw.json >> $HOME/$CHAIN_ID/onomy/config/genesis.json
rm -rf $HOME/raw.json
```

### Add seed node

Change the seed field in `$HOME/$CHAIN_ID/onomy/config/config.toml` to contain the following:

```
seeds = "5e0f5b9d54d3e038623ddb77c0b91b559ff13495@147.182.190.16:26656"
```

## Increasing the default open files limit
If we don't raise this value nodes will crash once the network grows large enough
```
sudo su -c "echo 'fs.file-max = 65536' >> /etc/sysctl.conf"
sysctl -p

sudo su -c "echo '* hard nofile 94000' >> /etc/security/limits.conf"
sudo su -c "echo '* soft nofile 94000' >> /etc/security/limits.conf"

sudo su -c "echo 'session required pam_limits.so' >> /etc/pam.d/common-session"
```
For this to take effect you'll need to (A) reboot (B) close and re-open all ssh sessions.
To check if this has worked run
```
ulimit -n
```
If you see 1024 then you need to reboot

### Syncing your node
To sync your node, you will use systemd, which manages the `onomy` daemon and automatically restarts it in case of failure. To use systemd, you will create a service file. Be sure to replace `<your_user>` with the user on your server, and `<chain_id>` with your chosen chain id.:

```bash:
sudo tee /etc/systemd/system/onomyd.service > /dev/null <<'EOF'
[Unit]
Description=Onomy daemon
After=network-online.target
[Service]
User=<your_user>
ExecStart=/home/<your_user>/go/bin/onomyd start --home /home/<your_user>/<chain_id>/onomy
Restart=on-failure
RestartSec=3
LimitNOFILE=4096
[Install]
WantedBy=multi-user.target
EOF
```

To start syncing:

```
# Start the node
sudo systemctl enable onomyd
sudo systemctl start onomyd
```

To check on the status of syncing:

```
onomyd status --output json | jq '.sync_info'
```

This will give output that looks like the following:
```
"sync_info": {
"latest_block_hash": "7BF95EED4EB50073F28CF833119FDB8C7DFE0562F611DF194CF4123A9C1F4640",
"latest_app_hash": "7C0C89EC4E903BAC730D9B3BB369D870371C6B7EAD0CCB5080B5F9D3782E3559",
"latest_block_height": "668538",
"latest_block_time": "2020-10-31T17:50:56.800119764Z",
"earliest_block_hash": "E7CAD87A4FDC47DFDE3D4E7C24D80D4C95517E8A6526E2D4BB4D6BC095404113",
"earliest_app_hash": "",
"earliest_block_height": "1",
"earliest_block_time": "2020-09-15T14:02:31Z",
"catching_up": false
}
```
The main thing to watch is that the block height is increasing. Once you are caught up with the chain, `catching_up` will become false. At that point, you can start using your node to create a validator

To check the logs of the node:

```
sudo journalctl -u onomyd -f
```

Your node is now fully set up! 
From here you can:
- request tokens from the faucet
- become a validator