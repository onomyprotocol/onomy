#Setting up constants
ONOMY_HOME=$HOME/.onomy
ONOMY_SRC=$ONOMY_HOME/src/onomy

ONOMY_VERSION="v0.0.3"

mkdir $HOME/.onomy
mkdir $ONOMY_HOME/src
mkdir $HOME/.onomy/bin

echo "-----------installing dependencies---------------"
sudo dnf -y update
sudo dnf -y copr enable ngompa/musl-libc
sudo dnf -y install https://dl.fedoraproject.org/pub/epel/epel-release-latest-8.noarch.rpm
sudo dnf -y install subscription-manager
sudo subscription-manager repos --enable codeready-builder-for-rhel-8-x86_64-rpms
sudo dnf -y --skip-broken install curl nano ca-certificates tar git jq gcc-c++ gcc-toolset-9 openssl-devel musl-devel musl-gcc gmp-devel perl python3 moreutils wget nodejs make hostname procps-ng
source "/opt/rh/gcc-toolset-9/enable"

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
cp onomyd $ONOMY_HOME/bin/onomyd

echo "-------------------adding binaries to path-----------------------"

export PATH=$PATH:$ONOMY_HOME/bin

echo "export PATH=$PATH" >> ~/.bashrc

echo "Onomy binaries are installed successfully."