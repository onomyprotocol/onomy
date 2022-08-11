#Setting up constants
ONOMY_HOME=$HOME/.onomy

ONOMY_VERSION="v0.0.7"
NODE_EXPORTER_VERSION="0.18.1"
COSMOVISOR_VERSION="cosmovisor-v1.0.1"

mkdir -p $ONOMY_HOME
mkdir -p $ONOMY_HOME/bin
mkdir -p $ONOMY_HOME/contracts
mkdir -p $ONOMY_HOME/logs
mkdir -p $ONOMY_HOME/cosmovisor/genesis/bin
mkdir -p $ONOMY_HOME/cosmovisor/upgrades/bin

echo "-----------installing dependencies---------------"
sudo dnf -y update
sudo dnf -y install https://dl.fedoraproject.org/pub/epel/epel-release-latest-8.noarch.rpm
sudo dnf -y install subscription-manager
sudo subscription-manager repos --enable codeready-builder-for-rhel-8-x86_64-rpms
sudo dnf makecache --refresh
sudo dnf -y --skip-broken install curl nano ca-certificates tar git jq moreutils wget hostname procps-ng pass libsecret pinentry crudini

set -eu

echo "----------------------installing onomy---------------"
curl -LO https://github.com/onomyprotocol/onomy/releases/download/$ONOMY_VERSION/onomyd
mv onomyd $ONOMY_HOME/cosmovisor/genesis/bin/onomyd

echo "----------------------installing cosmovisor---------------"
curl -LO https://github.com/onomyprotocol/onomy-sdk/releases/download/$COSMOVISOR_VERSION/cosmovisor
mv cosmovisor $ONOMY_HOME/bin/cosmovisor

echo "-------------------adding binaries to path-----------------------"
chmod +x $ONOMY_HOME/bin/*
export PATH=$PATH:$ONOMY_HOME/bin
chmod +x $ONOMY_HOME/cosmovisor/genesis/bin/*
export PATH=$PATH:$ONOMY_HOME/cosmovisor/genesis/bin

echo "export PATH=$PATH" >> ~/.bashrc

# set the cosmovisor environments
echo "export DAEMON_HOME=$ONOMY_HOME/" >> ~/.bashrc
echo "export DAEMON_NAME=onomyd" >> ~/.bashrc
echo "export DAEMON_RESTART_AFTER_UPGRADE=true" >> ~/.bashrc

echo "Onomy binaries are installed successfully."
