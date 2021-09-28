echo "-----------Installing_dependencies---------------"

 dnf -y update
 dnf -y install curl
 dnf -y install nano
 dnf -y install ca-certificates
 dnf -y install tar
 dnf -y install git
 dnf -y install jq
 dnf -y install gcc-c++
 dnf -y copr enable ngompa/musl-libc
 dnf -y install musl-devel
 dnf -y install musl-gcc

 dnf -y install golang
 dnf -y install gmp-devel
 dnf -y install perl python3
 dnf -y install moreutils
 dnf -y install wget
 dnf -y install screen

HOME_DIRECTORY=$HOME

echo "--------------installing_rust---------------------------"

curl https://sh.rustup.rs -sSf | bash -s -- -y
export PATH=$HOME_DIRECTORY/.cargo/bin:$PATH
cargo version

echo "----------------cloning_repository-------------------"

GRAVITY_DIR=$HOME_DIRECTORY/gravity
ONOMY_DIR=$HOME_DIRECTORY/go/onomy
git clone https://github.com/onomyprotocol/cosmos-gravity-bridge.git $GRAVITY_DIR
git clone -b dev https://github.com/onomyprotocol/onomy.git $ONOMY_DIR

echo "--------------install_golang---------------------------"
curl https://dl.google.com/go/go1.16.4.linux-amd64.tar.gz --output $HOME_DIRECTORY/go.tar.gz
tar -C $HOME_DIRECTORY -xzf $HOME_DIRECTORY/go.tar.gz
export PATH=$PATH:$HOME_DIRECTORY/go/bin

echo "----------------------building_gravity_artifact---------------"
#cd $GRAVITY_DIR/module
#make install
cd $ONOMY_DIR
make build
cp onomyd $HOME_DIRECTORY/go/bin/onomyd

echo "----------------building_orchestrator_artifact-------------"
cd $GRAVITY_DIR/orchestrator
rustup target add x86_64-unknown-linux-musl
cargo build --target=x86_64-unknown-linux-musl --release  --all
cp $GRAVITY_DIR/orchestrator/target/x86_64-unknown-linux-musl/release/gbt $HOME_DIRECTORY/go/bin/gbt


echo "---------------Installing_solidity-------------------"
cd $GRAVITY_DIR/solidity
dnf -y install nodejs
npm ci
chmod -R +x scripts
npm run typechain


echo "-------------------making_geth-----------------------"
cd $HOME_DIRECTORY
git clone https://github.com/ethereum/go-ethereum
cd go-ethereum/
make geth
cp build/bin/geth $HOME_DIRECTORY/go/bin/geth

echo "------------------ install fauset ------------------"
curl https://get.starport.network/faucet! | bash

# dnf -y install https://dl.fedoraproject.org/pub/epel/epel-release-latest-8.noarch.rpm
# sudo subscription-manager repos --enable codeready-builder-for-rhel-8-x86_64-rpms
export PATH=$PATH:$HOME_DIRECTORY/go/bin
cd $HOME_DIRECTORY



