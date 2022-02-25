FROM onomy/dev-base:latest

# ------------------ onomy chain ------------------

ARG ONOMY_DIR=/go/src/github.com/onomyprotocol/onomy/
# clone the repo
COPY . $ONOMY_DIR

# build onomyd artifact
WORKDIR $ONOMY_DIR
RUN make install

WORKDIR /root/home
RUN cp -r $ONOMY_DIR/deploy/devnet/scripts scripts
RUN mkdir onomy
RUN cp -r $ONOMY_DIR/deploy/devnet/.onomy $HOME/.onomy

CMD /bin/bash scripts/run.sh

# ------------------ gravity bridge/deployer ------------------

ARG GRAVITY_BRIDGE_VERSION="v0.0.5"
ADD https://github.com/onomyprotocol/cosmos-gravity-bridge/releases/download/$GRAVITY_BRIDGE_VERSION/Gravity.json /root/home/gravity/solidity/Gravity.json
ADD https://github.com/onomyprotocol/cosmos-gravity-bridge/releases/download/$GRAVITY_BRIDGE_VERSION/TestERC20A.json /root/home/gravity/solidity/TestERC20A.json
ADD https://github.com/onomyprotocol/cosmos-gravity-bridge/releases/download/$GRAVITY_BRIDGE_VERSION/contract-deployer /root/home/gravity/solidity/contract-deployer
RUN chmod +x /root/home/gravity/solidity/contract-deployer

ADD https://github.com/onomyprotocol/cosmos-gravity-bridge/releases/download/$GRAVITY_BRIDGE_VERSION/gbt /usr/bin/gbt
RUN chmod +x /usr/bin/gbt

