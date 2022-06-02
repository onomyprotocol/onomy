FROM fedora:33
RUN dnf install -y git make gcc gcc-c++ which iproute iputils procps-ng vim-minimal tmux net-tools htop tar jq npm openssl-devel perl rust cargo

RUN curl https://dl.google.com/go/go1.16.4.linux-amd64.tar.gz --output /root/go.tar.gz
RUN tar -C /root -xzf /root/go.tar.gz
RUN rm /root/go.tar.gz
ENV PATH=$PATH:/root/go/bin
ENV GOPATH=/root/go

ADD https://gethstore.blob.core.windows.net/builds/geth-linux-amd64-1.10.10-bb74230f.tar.gz /geth/
RUN cd /geth && tar -xvf * && mv /geth/**/geth /usr/bin/geth
# this is the Gravity bridge release, update it along with your go.mod entry
ARG ETH_BRIDGE_VERSION="v0.0.7"
ADD https://github.com/onomyprotocol/cosmos-gravity-bridge/releases/download/$ETH_BRIDGE_VERSION/test-runner /usr/bin/test-runner
RUN chmod +x /usr/bin/test-runner
ADD https://github.com/onomyprotocol/cosmos-gravity-bridge/releases/download/$ETH_BRIDGE_VERSION/Gravity.json /solidity/Gravity.json
ADD https://github.com/onomyprotocol/cosmos-gravity-bridge/releases/download/$ETH_BRIDGE_VERSION/TestERC20A.json /solidity/TestERC20A.json
ADD https://github.com/onomyprotocol/cosmos-gravity-bridge/releases/download/$ETH_BRIDGE_VERSION/TestERC20B.json /solidity/TestERC20B.json
ADD https://github.com/onomyprotocol/cosmos-gravity-bridge/releases/download/$ETH_BRIDGE_VERSION/TestERC20C.json /solidity/TestERC20C.json
ADD https://github.com/onomyprotocol/cosmos-gravity-bridge/releases/download/$ETH_BRIDGE_VERSION/contract-deployer /solidity/contract-deployer
RUN chmod +x /solidity/contract-deployer
# the actual source code for this repo, this **only** includes checked in files!
# this is a bit of a pain but it does speed things up a lot
ADD onomy.tar.gz /
# build steps for all codebases in this repo, must be below the add statement
RUN pushd /onomy/ && make install