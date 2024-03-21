#!/bin/bash

USER=$(whoami)
GROUP=$(id -gn)

sudo apt update && sudo apt upgrade -y
sudo apt install -y tmux vim crudini jq git liblz4-tool
git clone https://github.com/onomyprotocol/onomy ~/onomy

## Unnecessary after merge of scripts branch to mainnet
cd ~/onomy
git checkout scripts
## End unnecessary

cd ~/onomy/deploy/mainnet/chain/scripts/
source bin-ubuntu.sh
./init-seed-node.sh
#./init-statesync.sh
./expose-metrics.sh
./ubuntu-ulimits.sh
./autostake.sh

# Service setup
sudo ./add-service.sh cosmovisor-onomyd ${PWD}/start-cosmovisor-onomyd.sh
sudo ./add-service.sh node-exporter ${PWD}/start-node-exporter.sh
sudo chown -R $USER:$GROUP ~/.onomy/config/

# Oracle Cloud Ubuntu Firewall Config
IPTABLES_CONFIG=/etc/iptables/rules.v4
if test -f "$IPTABLES_CONFIG"; then
  # Oracle Cloud Ubuntu Firewall Config
  sudo sed -i 's/22 -j ACCEPT/&\n-A INPUT -p tcp -m state --state NEW -m tcp --dport 26656 -j ACCEPT/' $IPTABLES_CONFIG
  sudo sed -i 's/22 -j ACCEPT/&\n-A INPUT -p tcp -m state --state NEW -m tcp --dport 26657 -j ACCEPT/' $IPTABLES_CONFIG
  sudo iptables-restore < $IPTABLES_CONFIG
fi

echo "Completed seed node setup"
# Get seed id to share
echo "seed=$(onomyd tendermint show-node-id)@$(hostname -I | awk '{print $1}'):26656"
cosmovisor status
