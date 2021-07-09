GRAVITY_GENESIS_FILE="/root/.gravity/config"
GRAVITY_CONFIG_FILE="/root/.gravity/config/config.toml"
BUCKET_MASTER_GENESIS_FILE="/root/mainNode/testing/genesis.json"
BUCKET_MASTER_SEED_FILE="/root/mainNode/testing/seed"
MAIN_NODE="/root/mainNode"

GIT_HUB_USER=$1
GIT_HUB_PASS=$2
GIT_HUB_EMAIL=$3
GIT_HUB_BRANCH=$4

echo "Get pull updates"
git clone -b $GIT_HUB_BRANCH https://github.com/sunnyk56/onomy.git $MAIN_NODE

echo "Copying genesis file"
rm $GRAVITY_GENESIS_FILE/genesis.json
cp $BUCKET_MASTER_GENESIS_FILE  $GRAVITY_GENESIS_FILE
seed=$(cat $BUCKET_MASTER_SEED_FILE)

sed -i 's#seeds = ""#seeds = "'$seed'"#g' $GRAVITY_CONFIG_FILE

rm -r $MAIN_NODE

# Resets the blockchain database, removes address book files and start the node
gravity unsafe-reset-all
gravity --home /root/.gravity/ --address tcp://0.0.0.0:26655 --rpc.laddr tcp://0.0.0.0:26657 --grpc.address 0.0.0.0:9090 --log_level error --p2p.laddr tcp://0.0.0.0:26656 --rpc.pprof_laddr 0.0.0.0:6060 start