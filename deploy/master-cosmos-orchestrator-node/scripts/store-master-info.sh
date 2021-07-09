GRAVITY_CHAIN_DATA="/root/testchain"
GRAVITY_ASSETS="/root/assets"
GRAVITY_GENESIS_FILE="/root/testchain/gravity/config/genesis.json"
GRAVITY_SEED_FILE="/root/seed"
BUCKET_MASTER_CHAIN_DATA="master/testchain"
BUCKET_MASTER_ASSETS="master/assets"
BUCKET_MASTER_GENESIS_FILE="master/genesis.json"
BUCKET_MASTER_SEED="master/seed"

GIT_HUB_USER=$1
GIT_HUB_PASS=$2
GIT_HUB_EMAIL=$3
GIT_HUB_BRANCH=$4

echo "Get pull updates"
git pull origin $GIT_HUB_BRANCH

echo "add master genesis,json file"
rm -r master
mkdir master
mkdir $BUCKET_MASTER_CHAIN_DATA
mkdir $BUCKET_MASTER_ASSETS
touch $BUCKET_MASTER_GENESIS_FILE
touch $BUCKET_MASTER_SEED
echo "Copying genesis file"
cp $GRAVITY_GENESIS_FILE $BUCKET_MASTER_GENESIS_FILE
cp $GRAVITY_SEED_FILE $BUCKET_MASTER_SEED
cp -r $GRAVITY_CHAIN_DATA/. $BUCKET_MASTER_CHAIN_DATA
cp -r $GRAVITY_ASSETS/. $BUCKET_MASTER_ASSETS
echo "git add command"
git add .
echo "git add git config command"
git config --global user.email $GIT_HUB_EMAIL
git config --global user.name $GIT_HUB_USER
# //TODO this repo name should be pass as parameter
git remote set-url origin https://$GIT_HUB_USER:$GIT_HUB_PASS@github.com/sunnyk56/onomy.git

echo "git commit command"
git commit -m "add genesis file"

echo "git fetch command"
git fetch
echo "git push command"
git push origin $GIT_HUB_BRANCH