FROM leeway321/gravity-bridge-binary:latest as builder

EXPOSE 26656 26657 1317 9090
ARG PUBLIC_IP=${PUBLIC_IP}
ARG GIT_HUB_USER=${GIT_HUB_USER}
ARG GIT_HUB_PASS=${GIT_HUB_PASS}
ARG GIT_HUB_EMAIL=${GIT_HUB_EMAIL}
ARG GIT_HUB_BRANCH=${GIT_HUB_BRANCH}

ARG GRAVITY_GENESIS_COINS="100000000000stake,100000000000samoleans"

# purpose of clone the repo is to commit and push genesis file and seed in this repo to use by other node
# //TODO this clone may be change in future
RUN git clone -b $GIT_HUB_BRANCH https://github.com/sunnyk56/onomy.git
RUN cp -r /root/onomy/deploy/peer-cosmos-orchestrator-node/scripts scripts
RUN cp -r /root/onomy/master/assets assets

RUN chmod -R +x /root/scripts
RUN chmod -R +x /root/assets

RUN sh scripts/init.sh

##----------------------------- commit master genesis,json file and seed-----
RUN gravity --home /root/testchain/gravity tendermint show-node-id > seed
RUN sed -i 's#'$(gravity --home /root/testchain/gravity tendermint show-node-id)'#'$(gravity --home /root/testchain/gravity tendermint show-node-id)'@'$PUBLIC_IP':26656#g' seed

WORKDIR /root/onomy
RUN sh deploy/peer-cosmos-orchestrator-node/scripts/store-peer-info.sh $GIT_HUB_USER $GIT_HUB_PASS $GIT_HUB_EMAIL $GIT_HUB_BRANCH

#WORKDIR /root

#CMD ["sh", "/root/scripts/cosmos-run.sh"]