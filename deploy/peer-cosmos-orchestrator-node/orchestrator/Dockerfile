FROM leeway321/gravity-bridge-binary:latest as builder

EXPOSE 61278
ARG GRAVITY_HOST=${GRAVITY_HOST}
ARG ETH_HOST=${ETH_HOST}
ARG GIT_HUB_BRANCH=${GIT_HUB_BRANCH}

# //TODO this clone may be change in future
RUN git clone -b $GIT_HUB_BRANCH https://github.com/sunnyk56/onomy.git
RUN cp -r /root/onomy/deploy/peer-cosmos-orchestrator-node/scripts scripts
RUN cp -r /root/onomy/master/assets assets
RUN cp -r /root/onomy/peerInfo/testchain testchain

RUN chmod -R +x /root/scripts
RUN chmod -R +x /root/assets
RUN chmod -R +x /root/testchain

CMD ["sh", "scripts/orchestrator-run.sh $GRAVITY_HOST  $ETH_HOST"]