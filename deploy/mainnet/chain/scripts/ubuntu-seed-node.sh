#!/bin/bash

sudo apt update && sudo apt upgrade -y
sudo apt install -y tmux vim crudini jq git liblz4-tool
git clone https://github.com/onomyprotocol/onomy ~/onomy

## Unnecessary after merge of scripts branch to mainnet
cd ~/onomy
git checkout scripts

cd ~/onomy/deploy/mainnet/chain/scripts/
source bin-ubuntu.sh
./init-seed-node.sh
./init-statesync.sh
./expose-metrics.sh

## ulimits
echo 'fs.file-max = 65536' | sudo tee -a /etc/sysctl.conf
sudo sysctl -p
echo '* hard nofile 94000' | sudo tee -a /etc/security/limits.conf
echo '* soft nofile 94000' | sudo tee -a /etc/security/limits.conf
# The following line is necessary on RHEL/Oracle but not Ubuntu
# echo 'session required /lib/security/pam_limits.so' | sudo tee -a /etc/pam.d/login

./autostake.sh

# Service setup
sudo ./add-service.sh cosmovisor-onomyd ${PWD}/start-cosmovisor-onomyd.sh
sudo ./add-service.sh node-exporter ${PWD}/start-node-exporter.sh

# Oracle Cloud Ubuntu Firewall Config
IPTABLES_CONFIG=/etc/iptables/rules.v4
if test -f "$IPTABLES_CONFIG"; then
  # Oracle Cloud Ubuntu Firewall Config
  sudo sed -i 's/22 -j ACCEPT/&\n-A INPUT -p tcp -m state --state NEW -m tcp --dport 26656 -j ACCEPT/' $IPTABLES_CONFIG
  sudo sed -i 's/22 -j ACCEPT/&\n-A INPUT -p tcp -m state --state NEW -m tcp --dport 26657 -j ACCEPT/' $IPTABLES_CONFIG
  sudo iptables-restore < $IPTABLES_CONFIG
fi

sudo chown -R ubuntu:ubuntu ~/.onomy/config/

echo "Completed node setup"
# Get seed id to share
echo "seed=$(onomyd tendermint show-node-id)@$(hostname -I | awk '{print $1}'):26656"
cosmovisor status
