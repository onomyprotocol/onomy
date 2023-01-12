#Setting up constants
ONOMY_HOME=$HOME/.onomy
ONOMY_SRC=$ONOMY_HOME/src/onomy
COSMOVISOR_SRC=$ONOMY_HOME/src/cosmovisor

ONOMY_VERSION="v1.0.1"
COSMOVISOR_VERSION="cosmovisor-v1.0.1"

mkdir -p $ONOMY_HOME
mkdir -p $ONOMY_HOME/src
mkdir -p $ONOMY_HOME/bin
mkdir -p $ONOMY_HOME/logs
mkdir -p $ONOMY_HOME/cosmovisor/genesis/bin
mkdir -p $ONOMY_HOME/cosmovisor/upgrades/

echo "-----------installing dependencies---------------"
sudo dnf -y update
sudo dnf -y copr enable ngompa/musl-libc
sudo dnf -y install https://dl.fedoraproject.org/pub/epel/epel-release-latest-8.noarch.rpm
sudo dnf -y install subscription-manager
sudo subscription-manager config --rhsm.manage_repos=1
sudo subscription-manager repos --enable codeready-builder-for-rhel-8-x86_64-rpms
sudo dnf makecache --refresh
sudo dnf -y --skip-broken install curl nano ca-certificates tar git jq gcc-c++ gcc-toolset-9 openssl-devel musl-devel musl-gcc gmp-devel perl python3 moreutils wget nodejs make hostname procps-ng pass libsecret pinentry crudini cmake

gcc_source="/opt/rh/gcc-toolset-9/enable"
if test -f $gcc_source; then
   source gcc_source
fi

set -eu

echo "--------------installing golang---------------------------"
curl https://dl.google.com/go/go1.16.4.linux-amd64.tar.gz --output $HOME/go.tar.gz
tar -C $HOME -xzf $HOME/go.tar.gz
rm $HOME/go.tar.gz
export PATH=$PATH:$HOME/go/bin
export GOPATH=$HOME/go
echo "export GOPATH=$HOME/go" >> ~/.bashrc
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
