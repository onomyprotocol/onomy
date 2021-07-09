FROM leeway321/gravity-bridge-binary:latest as builder

EXPOSE 26656 26657 61278 1317 9090 8545
ARG CHAIN_ID="gravity-testnet"
ARG KEYRING="--keyring-backend test"
ARG GRAVITY_HOST="0.0.0.0"
ARG GRAVITY_NODE_CONFIG=".gravity/config/config.toml"
ARG GRAVITY_APP_CONFIG=".gravity/config/app.toml"
ARG PUBLIC_IP=${PUBLIC_IP}
ARG GIT_HUB_USER=${GIT_HUB_USER}
ARG GIT_HUB_PASS=${GIT_HUB_PASS}
ARG GIT_HUB_EMAIL=${GIT_HUB_EMAIL}
ARG GIT_HUB_BRANCH=${GIT_HUB_BRANCH}

ARG GRAVITY_GENESIS_COINS="100000000000stake,100000000000samoleans"

# purpose of clone the repo is to get updated Ethgenesis.json file
# //TODO this clone may be change in future
RUN git clone -b $GIT_HUB_BRANCH https://github.com/sunnyk56/onomy.git
RUN cp -r /root/onomy/deploy/master-cosmos-orchestrator-node/scripts scripts
RUN cp -r /root/onomy/deploy/master-cosmos-orchestrator-node/assets assets

RUN chmod -R +x /root/scripts
RUN chmod -R +x /root/assets

RUN sh scripts/init.sh

##----------------------------- commit master genesis,json file and seed-----
RUN gravity --home /root/testchain/gravity tendermint show-node-id > seed
RUN sed -i 's#'$(gravity --home /root/testchain/gravity tendermint show-node-id)'#'$(gravity --home /root/testchain/gravity tendermint show-node-id)'@'$PUBLIC_IP':26656#g' seed

WORKDIR /root/onomy
RUN sh deploy/master-cosmos-orchestrator-node/scripts/store-master-info.sh $GIT_HUB_USER $GIT_HUB_PASS $GIT_HUB_EMAIL $GIT_HUB_BRANCH

WORKDIR /root

CMD ["sh", "scripts/run.sh"]