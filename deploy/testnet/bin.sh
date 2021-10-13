#Setting up constants
GRAVITY_DIR=$HOME/.onomy/gravity
ONOMY_HOME=$HOME/.onomy
#Creating Directories
mkdir $HOME/.onomy
mkdir $HOME/.onomy/bin
mkdir $HOME/.onomy/gravity

echo "-----------Installing_dependencies---------------"
dnf -y update
dnf -y copr enable ngompa/musl-libc
dnf -y install https://dl.fedoraproject.org/pub/epel/epel-release-latest-8.noarch.rpm
sudo subscription-manager repos --enable codeready-builder-for-rhel-8-x86_64-rpms
dnf -y install curl nano ca-certificates tar git jq gcc-c++ musl-devel musl-gcc golang gmp-devel perl python3 moreutils wget nodejs make

echo "--------------installing_rust---------------------------"
curl https://sh.rustup.rs -sSf | bash -s -- -y
export PATH=$HOME/.cargo/bin:$PATH
cargo version

echo "----------------cloning_repository-------------------"
git clone -b v0.0.0-20210915184851-orch-nomarket https://github.com/onomyprotocol/cosmos-gravity-bridge.git $GRAVITY_DIR
git clone -b v0.0.1 https://github.com/onomyprotocol/onomy.git $ONOMY_HOME/onomy

echo "--------------install_golang---------------------------"
curl https://dl.google.com/go/go1.16.4.linux-amd64.tar.gz --output $HOME/go.tar.gz
tar -C $HOME -xzf $HOME/go.tar.gz
export PATH=$PATH:$HOME/go/bin

echo "----------------------building_gravity_artifact---------------"
#cd $GRAVITY_DIR/module
#make install
cd $ONOMY_HOME/onomy
make build
cp onomyd $ONOMY_HOME/bin/onomyd

echo "----------------building_orchestrator_artifact-------------"
cd $GRAVITY_DIR/orchestrator
rustup target add x86_64-unknown-linux-musl
cargo build --target=x86_64-unknown-linux-musl --release  --all
cp $GRAVITY_DIR/orchestrator/target/x86_64-unknown-linux-musl/release/gbt $ONOMY_HOME/bin/gbt


echo "---------------Installing_solidity-------------------"
cd $GRAVITY_DIR/solidity
npm ci
chmod -R +x scripts
npm run typechain


echo "-------------------making_geth-----------------------"
cd $HOME
git clone https://github.com/ethereum/go-ethereum
cd go-ethereum/
make geth
cp build/bin/geth $ONOMY_HOME/bin/geth

echo "------------------ install fauset ------------------"
curl https://get.starport.network/faucet! | bash
cd $HOME

echo "export PATH=$PATH:$HOME/.cargo/bin:$ONOMY_HOME/bin" >> $HOME/.bashrc

