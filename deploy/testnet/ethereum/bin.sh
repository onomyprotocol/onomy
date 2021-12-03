#Setting up constants
ONOMY_HOME=$HOME/.onomy
GETH_SRC=$ONOMY_HOME/src/go-ethereum

#Creating Directories
mkdir $HOME/.onomy
mkdir $ONOMY_HOME/src
mkdir $HOME/.onomy/bin

echo "-----------installing dependencies---------------"
sudo dnf -y update
sudo dnf -y copr enable ngompa/musl-libc
sudo dnf -y install https://dl.fedoraproject.org/pub/epel/epel-release-latest-8.noarch.rpm
sudo dnf -y install subscription-manager
sudo subscription-manager repos --enable codeready-builder-for-rhel-8-x86_64-rpms
sudo dnf -y install curl nano ca-certificates tar git jq gcc-c++ openssl-devel musl-devel musl-gcc gmp-devel perl python3 moreutils wget nodejs make hostname procps-ng

echo "--------------installing golang---------------------------"
curl https://dl.google.com/go/go1.16.4.linux-amd64.tar.gz --output $HOME/go.tar.gz
tar -C $HOME -xzf $HOME/go.tar.gz
export PATH=$PATH:$HOME/go/bin
export GOPATH=$HOME/go
echo "export GOPATH=$HOME/go" >> ~/.bashrc
go version

echo "-------------------installing geth-----------------------"
git clone https://github.com/ethereum/go-ethereum $GETH_SRC
cd $GETH_SRC
make geth
cp build/bin/geth $ONOMY_HOME/bin/geth

echo "-------------------adding binaries to path-----------------------"

export PATH=$PATH:$ONOMY_HOME/bin

echo "export PATH=$PATH" >> ~/.bashrc

echo "Ethereum binaries are installed successfully."