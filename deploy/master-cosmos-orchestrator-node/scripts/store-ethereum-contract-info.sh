GRAVITY_CHAIN_DATA="/root/testchain"
BUCKET_MASTER_CHAIN_DATA="master/testchain"

GIT_HUB_USER=$1
GIT_HUB_PASS=$2
GIT_HUB_EMAIL=$3
GIT_HUB_BRANCH=$4

echo "Get pull updates"
git pull origin $GIT_HUB_BRANCH

echo "add master contract information file"
rm -r $BUCKET_MASTER_CHAIN_DATA
mkdir $BUCKET_MASTER_CHAIN_DATA
echo "Copying contract file"
cp -r $GRAVITY_CHAIN_DATA/. $BUCKET_MASTER_CHAIN_DATA
echo "git add command"
git add .
echo "git add git config command"
git config --global user.email $GIT_HUB_EMAIL
git config --global user.name $GIT_HUB_USER
# //TODO this repo name should be pass as parameter
git remote set-url origin https://$GIT_HUB_USER:$GIT_HUB_PASS@github.com/sunnyk56/onomy.git

echo "git commit command"
git commit -m "add contract file"

echo "git fetch command"
git fetch
echo "git push command"
git push origin $GIT_HUB_BRANCH