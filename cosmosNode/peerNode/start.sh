GENTEX_FILE="/root/.gravity/config/gentx/."
VALIDATOR_FILE="/root/validator_key"
BUCKET_MASTER_GENTX_FILE="peerValidatorInfo/gentx"
BUCKET_MASTER_VALIDATOR_FILE="peerValidatorInfo/validator_key"

GIT_HUB_USER=$1
GIT_HUB_PASS=$2
GIT_HUB_EMAIL=$3
GIT_HUB_BRANCH=$4

echo "Get pull updates"
git pull origin $GIT_HUB_BRANCH

echo "Gentx file moving"
rm -r peerValidatorInfo
mkdir peerValidatorInfo
mkdir $BUCKET_MASTER_GENTX_FILE
touch $BUCKET_MASTER_VALIDATOR_FILE

echo "Copying validator file"
cp $VALIDATOR_FILE $BUCKET_MASTER_VALIDATOR_FILE

echo "Copying gentx file"
cp -r $GENTEX_FILE $BUCKET_MASTER_GENTX_FILE

echo "git add command"
git add .
echo "git add git config command"

git config --global user.email $GIT_HUB_EMAIL
git config --global user.name $GIT_HUB_USER
git remote set-url origin https://$GIT_HUB_USER:$GIT_HUB_PASS@github.com/sunnyk56/onomy.git

echo "git commit command"
git commit -m "added peerValidatorInfo folder"

echo "git fetch command"
git fetch

echo "git push command"
git push origin $GIT_HUB_BRANCH