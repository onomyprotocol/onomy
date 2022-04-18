FROM onomy/dev-base:latest

# ------------------ gravity bridge/deployer ------------------

ARG ETH_BRIDGE_VERSION="v0.0.6"
ADD https://github.com/onomyprotocol/cosmos-gravity-bridge/releases/download/$ETH_BRIDGE_VERSION/Gravity.json /root/home/gravity/solidity/Gravity.json
ADD https://github.com/onomyprotocol/cosmos-gravity-bridge/releases/download/$ETH_BRIDGE_VERSION/TestERC20A.json /root/home/gravity/solidity/TestERC20A.json
ADD https://github.com/onomyprotocol/cosmos-gravity-bridge/releases/download/$ETH_BRIDGE_VERSION/contract-deployer /root/home/gravity/solidity/contract-deployer
RUN chmod +x /root/home/gravity/solidity/contract-deployer

ADD https://github.com/onomyprotocol/cosmos-gravity-bridge/releases/download/$ETH_BRIDGE_VERSION/gbt /usr/bin/gbt
RUN chmod +x /usr/bin/gbt

# ------------------ test dependencies ------------------

WORKDIR /root/home/wnom
# the build is called form the ../ folder so the line copies ../ folder
COPY env env
WORKDIR /root/home/wnom/env
RUN npm install
RUN mkdir logs
RUN touch logs/eth.log
RUN touch logs/orchestrator.log
CMD tail -f logs/eth.log & tail -f logs/orchestrator.log