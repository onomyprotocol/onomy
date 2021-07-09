GRAVITY_CHAIN_DATA="/root/testchain"
GRAVITY_ASSETS="/root/assets"
GRAVITY_GENESIS_FILE="/root/testchain/gravity/config/genesis.json"
GRAVITY_GENTX="/root/testchain/gravity/config/gentx"
BUCKET_PEER_CHAIN_DATA="peerInfo/testchain"
BUCKET_PEER_GENTX_DATA="peerInfo/gentx"
BUCKET_MASTER_ASSETS="master/assets"

GIT_HUB_USER=$1
GIT_HUB_PASS=$2
GIT_HUB_EMAIL=$3
GIT_HUB_BRANCH=$4

echo "Get pull updates"
git pull origin $GIT_HUB_BRANCH

echo "add master genesis,json file"
rm -r peerInfo
mkdir peerInfo
mkdir $BUCKET_PEER_CHAIN_DATA
mkdir $BUCKET_PEER_GENTX_DATA
mkdir $BUCKET_MASTER_ASSETS
echo "Copying gentx and testchain data in peerInfo and EthGenesis.json in master/assets"
cp -r $GRAVITY_GENTX/. $BUCKET_PEER_GENTX_DATA
cp -r $GRAVITY_CHAIN_DATA/. $BUCKET_PEER_CHAIN_DATA
cp -r $GRAVITY_ASSETS/. $BUCKET_MASTER_ASSETS
echo "git add command"
git add master
git add peerInfo
echo "git add git config command"
git config --global user.email $GIT_HUB_EMAIL
git config --global user.name $GIT_HUB_USER
# //TODO this repo name should be pass as parameter
git remote set-url origin https://$GIT_HUB_USER:$GIT_HUB_PASS@github.com/sunnyk56/onomy.git

echo "git commit command"
git commit -m "add gentx and testchain data in peerInfo and EthGenesis.json in master/assets by peer node"

echo "git fetch command"
git fetch
echo "git push command"
git push origin $GIT_HUB_BRANCH