GRAVITY_CONFIG_FILE="/root/testchain/gravity/config"
GRAVITY_GENESIS_FILE="/root/testchain/gravity/config/genesis.json"
BUCKET_MASTER_GENESIS_FILE="master/genesis.json"

PEER_INFO="/root/gravity-bridge/peerInfo"
PEER_INFO_VALIDATOR_KEY="/root/gravity-bridge/peerInfo/validator_key.json"
PEER_INFO_ORCHESTRATOR_KEY="/root/gravity-bridge/peerInfo/orchestrator_key.json"

GIT_HUB_USER=$1
GIT_HUB_PASS=$2
GIT_HUB_EMAIL=$3
GIT_HUB_BRANCH=$4
GRAVITY_HOME_FLAG="--home /root/testchain/gravity"


STAKE_DENOM="stake"
NORMAL_DENOM="samoleans"
GRAVITY_GENESIS_COINS="100000000000$STAKE_DENOM,100000000000$NORMAL_DENOM"

echo "Get pull updates"
git pull origin $GIT_HUB_BRANCH

echo "extracting validator address"
validatorKey="$(jq .address $PEER_INFO_VALIDATOR_KEY)"
echo validatorKey

echo "Adding validator addresses to genesis files"
gravity $GRAVITY_HOME_FLAG add-genesis-account $validatorKey $GRAVITY_GENESIS_COINS

echo "Collecting gentxs files in config gentx"
cp $PEER_INFO/gentx/*.json $GRAVITY_CONFIG_FILE/gentx/


echo "Adding orchestrator keys to genesis"
GRAVITY_ORCHESTRATOR_KEY="$(jq .address $PEER_INFO_ORCHESTRATOR_KEY)"

jq ".app_state.auth.accounts += [{\"@type\": \"/cosmos.auth.v1beta1.BaseAccount\",\"address\": $GRAVITY_ORCHESTRATOR_KEY,\"pub_key\": null,\"account_number\": \"0\",\"sequence\": \"0\"}]" $GRAVITY_CONFIG_FILE/genesis.json | sponge $GRAVITY_CONFIG_FILE/genesis.json
jq ".app_state.bank.balances += [{\"address\": $GRAVITY_ORCHESTRATOR_KEY,\"coins\": [{\"denom\": \"$NORMAL_DENOM\",\"amount\": \"100000000000\"},{\"denom\": \"$STAKE_DENOM\",\"amount\": \"100000000000\"}]}]" $GRAVITY_CONFIG_FILE/genesis.json | sponge $GRAVITY_CONFIG_FILE/genesis.json


echo "Collecting gentxs"
gravity $GRAVITY_HOME_FLAG collect-gentxs

# update genesis file ------
#rm -r peerInfo
rm -f $BUCKET_MASTER_GENESIS_FILE
touch $BUCKET_MASTER_GENESIS_FILE
echo "Copying genesis file"
cp $GRAVITY_GENESIS_FILE $BUCKET_MASTER_GENESIS_FILE
echo "git add command"
git add .
echo "git add git config command"
git config --global user.email $GIT_HUB_EMAIL
git config --global user.name $GIT_HUB_USER
git remote set-url origin https://$GIT_HUB_USER:$GIT_HUB_PASS@github.com/sunnyk56/gravity-bridge.git
echo "git commit command"
git commit -m "update genesis file by master node"
echo "git push command"
git push origin $GIT_HUB_BRANCH

# Resets the blockchain database, removes address book files and start the node
gravity $GRAVITY_HOME_FLAG unsafe-reset-all
gravity $GRAVITY_HOME_FLAG --address tcp://0.0.0.0:26655 --rpc.laddr tcp://0.0.0.0:26657 --grpc.address 0.0.0.0:9090 --log_level error --p2p.laddr tcp://0.0.0.0:26656 --rpc.pprof_laddr 0.0.0.0:6060 start
