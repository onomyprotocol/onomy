#Use Ubuntu Latest

#Setting up constants
ONOMY_HOME=$HOME/.onomy
ONOMY_SRC=$ONOMY_HOME/src/onomy
COSMOVISOR_SRC=$ONOMY_HOME/src/cosmovisor

ONOMY_VERSION="v1.1.4"
NODE_EXPORTER_VERSION="0.18.1"
COSMOVISOR_VERSION="cosmovisor-v1.0.1"

mkdir -p $ONOMY_HOME
mkdir -p $ONOMY_HOME/src
mkdir -p $ONOMY_HOME/bin
mkdir -p $ONOMY_HOME/contracts
mkdir -p $ONOMY_HOME/logs
mkdir -p $ONOMY_HOME/cosmovisor/genesis/bin
mkdir -p $ONOMY_HOME/cosmovisor/upgrades/

echo "-----------installing dependencies---------------"
sudo apt install build-essential crudini jq

#echo "--------------installing rust---------------------------"
#curl https://sh.rustup.rs -sSf | bash -s -- -y
#export PATH=$HOME/.cargo/bin:$PATH
#cargo version

echo "--------------installing golang---------------------------"
curl https://dl.google.com/go/go1.22.0.linux-amd64.tar.gz --output $HOME/go.tar.gz
sudo tar -C /usr/local -xzf $HOME/go.tar.gz
rm $HOME/go.tar.gz
mkdir -p $HOME/go/bin
export PATH=$PATH:$HOME/go/bin
export PATH=$PATH:/usr/local/go/bin
export GOPATH=$HOME/go
echo "export GOPATH=$HOME/go" >> ~/.profile
go version

echo "----------------------installing onomy---------------"
git clone -b $ONOMY_VERSION https://github.com/onomyprotocol/onomy.git $ONOMY_SRC
cd $ONOMY_SRC
make build
mv onomyd $ONOMY_HOME/cosmovisor/genesis/bin/onomyd

echo "-------------------installing cosmovisor-----------------------"
git clone -b $COSMOVISOR_VERSION https://github.com/onomyprotocol/onomy-sdk $COSMOVISOR_SRC
cd $COSMOVISOR_SRC
make cosmovisor
cp cosmovisor/cosmovisor $ONOMY_HOME/bin/cosmovisor

# echo "----------------installing eth bridge gbt-------------"
# git clone -b $ETH_BRIDGE_VERSION https://github.com/onomyprotocol/arc.git $ETH_BRIDGE_SRC
# cd $ETH_BRIDGE_SRC/orchestrator
# rm -rf .cargo/config
# cargo build --release --all
# cp $ETH_BRIDGE_SRC/orchestrator/target/release/gbt $ONOMY_HOME/bin/gbt

echo "-------------------installing node_exporter-----------------------"
curl -LO "https://github.com/prometheus/node_exporter/releases/download/v$NODE_EXPORTER_VERSION/node_exporter-$NODE_EXPORTER_VERSION.linux-amd64.tar.gz"
tar -xvf "node_exporter-$NODE_EXPORTER_VERSION.linux-amd64.tar.gz"
mv "node_exporter-$NODE_EXPORTER_VERSION.linux-amd64/node_exporter" $ONOMY_HOME/bin/node_exporter
rm -r "node_exporter-$NODE_EXPORTER_VERSION.linux-amd64"
rm "node_exporter-$NODE_EXPORTER_VERSION.linux-amd64.tar.gz"

echo "-------------------adding binaries to path-----------------------"
chmod +x $ONOMY_HOME/bin/*
export PATH=$PATH:$ONOMY_HOME/bin
chmod +x $ONOMY_HOME/cosmovisor/genesis/bin/*
export PATH=$PATH:$ONOMY_HOME/cosmovisor/genesis/bin

echo "export PATH=$PATH" >> ~/.profile

# set the cosmovisor environments
echo "export DAEMON_HOME=$ONOMY_HOME/" >> ~/.profile
echo "export DAEMON_NAME=onomyd" >> ~/.profile
echo "export DAEMON_RESTART_AFTER_UPGRADE=true" >> ~/.profile
source $HOME/.profile
ulimit -S -n 65536

echo "Onomy binaries are installed successfully."
