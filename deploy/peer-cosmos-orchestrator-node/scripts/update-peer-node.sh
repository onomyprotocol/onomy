GRAVITY_GENESIS_FILE="/root/testchain/gravity/config/genesis.json"
GRAVITY_ASSETS="/root/assets"
BUCKET_MASTER_GENESIS_FILE="master/genesis.json"
BUCKET_MASTER="/root/onomy/master"

GIT_HUB_BRANCH=$1

echo "Get pull updates"
git pull origin $GIT_HUB_BRANCH

echo "updating EthGenesis.json in the root assets directory"
cp $BUCKET_MASTER/. $GRAVITY_ASSETS


echo "Copying genesis file"
rm -f $GRAVITY_GENESIS_FILE
touch $GRAVITY_GENESIS_FILE
cp $BUCKET_MASTER_GENESIS_FILE $GRAVITY_GENESIS_FILE

echo "Run the cosmos-run scripts"
sh /root/scripts/cosmos-run.sh