#Setting up constants
ONOMY_HOME=$HOME/.onomy
GRAVITY_SRC=$ONOMY_HOME/src/gravity
ONOMY_SRC=$ONOMY_HOME/src/onomy
GETH_SRC=$ONOMY_HOME/src/go-ethereum

GRAVITY_VERSION="v0.0.0-20210915184851-orch-nomarket"
ONOMY_VERSION="v0.0.3"

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
sudo dnf -y install curl nano ca-certificates tar git jq gcc-c++ gcc-toolset-9 openssl-devel musl-devel musl-gcc gmp-devel perl python3 moreutils wget nodejs make hostname procps-ng

echo "--------------installing rust---------------------------"
curl https://sh.rustup.rs -sSf | bash -s -- -y
export PATH=$HOME/.cargo/bin:$PATH
cargo version

echo "--------------installing golang---------------------------"
curl https://dl.google.com/go/go1.16.4.linux-amd64.tar.gz --output $HOME/go.tar.gz
tar -C $HOME -xzf $HOME/go.tar.gz
export PATH=$PATH:$HOME/go/bin
export GOPATH=$HOME/go
echo "export GOPATH=$HOME/go" >> ~/.bashrc
go version

echo "----------------------installing onomy---------------"
git clone -b $ONOMY_VERSION https://github.com/onomyprotocol/onomy.git $ONOMY_SRC
cd $ONOMY_SRC
make build
cp onomyd $ONOMY_HOME/bin/onomyd

echo "----------------installing gravity gbt-------------"
git clone -b $GRAVITY_VERSION https://github.com/onomyprotocol/cosmos-gravity-bridge.git $GRAVITY_SRC
cd $GRAVITY_SRC/orchestrator
rustup target add x86_64-unknown-linux-musl
cargo build --target=x86_64-unknown-linux-musl --release  --all
cp $GRAVITY_SRC/orchestrator/target/x86_64-unknown-linux-musl/release/gbt $ONOMY_HOME/bin/gbt

echo "---------------installing gravity solidity-------------------"
cd $GRAVITY_SRC/solidity
npm ci
chmod -R +x scripts
npm run typechain
npm run compile-deployer

echo "-------------------installing geth-----------------------"
git clone https://github.com/ethereum/go-ethereum $GETH_SRC
cd $GETH_SRC
make geth
cp build/bin/geth $ONOMY_HOME/bin/geth

echo "-------------------adding binaries to path-----------------------"

export PATH=$PATH:$ONOMY_HOME/bin

echo "export PATH=$PATH" >> ~/.bashrc

echo "Onomy binaries are installed successfully."