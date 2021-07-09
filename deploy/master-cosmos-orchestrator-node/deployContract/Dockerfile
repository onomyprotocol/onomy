FROM leeway321/gravity-bridge-binary:latest as builder

ARG GRAVITY_HOST=${GRAVITY_HOST}
ARG ETH_HOST=${ETH_HOST}
ARG GIT_HUB_USER=${GIT_HUB_USER}
ARG GIT_HUB_PASS=${GIT_HUB_PASS}
ARG GIT_HUB_EMAIL=${GIT_HUB_EMAIL}
ARG GIT_HUB_BRANCH=${GIT_HUB_BRANCH}

# purpose of clone the repo is to get updated Ethgenesis.json file and store contract address
# //TODO this clone may be change in future
RUN git clone -b $GIT_HUB_BRANCH https://github.com/sunnyk56/onomy.git
RUN cp -r /root/onomy/deploy/master-cosmos-orchestrator-node/scripts scripts
RUN cp -r /root/onomy/master/assets assets
RUN cp -r /root/onomy/master/testchain testchain


RUN chmod -R +x /root/scripts
RUN chmod -R +x /root/assets
RUN chmod -R +x /root/testchain

WORKDIR /root/onomy
RUN sh /root/onomy/deploy/master-cosmos-orchestrator-node/scripts/deploy-contract.sh $GIT_HUB_USER $GIT_HUB_PASS $GIT_HUB_EMAIL $GIT_HUB_BRANCH $GRAVITY_HOST $ETH_HOST

##----------------------------- commit save Contract address-----
#RUN sh deploy/master-cosmos-orchestrator-node/scripts/store-ethereum-contract-info.sh $GIT_HUB_USER $GIT_HUB_PASS $GIT_HUB_EMAIL $GIT_HUB_BRANCH