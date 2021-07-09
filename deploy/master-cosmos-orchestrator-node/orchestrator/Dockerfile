FROM leeway321/gravity-bridge-binary:latest as builder

EXPOSE 61278
ARG PUBLIC_IP=${PUBLIC_IP}
ARG GIT_HUB_USER=${GIT_HUB_USER}
ARG GIT_HUB_PASS=${GIT_HUB_PASS}
ARG GIT_HUB_EMAIL=${GIT_HUB_EMAIL}
ARG GIT_HUB_BRANCH=${GIT_HUB_BRANCH}

# purpose of clone the repo is to get updated Ethgenesis.json file
# //TODO this clone may be change in future
RUN git clone -b $GIT_HUB_BRANCH https://github.com/sunnyk56/onomy.git
RUN cp -r /root/onomy/deploy/master-cosmos-orchestrator-node/scripts scripts
#RUN cp -r /root/onomy/master/assets assets
RUN cp -r /root/onomy/master/testchain testchain

RUN chmod -R +x /root/scripts
RUN chmod -R +x /root/testchain

CMD ["sh", "scripts/orchestrator-run.sh"]