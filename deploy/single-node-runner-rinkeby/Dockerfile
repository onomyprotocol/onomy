FROM ubuntu:latest

RUN apt-get update -y -q && apt-get upgrade -yq
# common (DEBIAN_FRONTEND is a fix for tzdata)
RUN DEBIAN_FRONTEND="noninteractive" apt-get install --no-install-recommends -yq software-properties-common \
                                                curl \
                                                build-essential \
                                                ca-certificates \
                                                tar \
                                                git

# ------------------ onomy chain ------------------

ARG ONOMY_DIR=/go/src/github.com/onomyprotocol/onomy/
# clone the repo
COPY . $ONOMY_DIR

# install golang
RUN curl https://dl.google.com/go/go1.16.4.linux-amd64.tar.gz --output go.tar.gz
RUN tar -C /usr/local -xzf go.tar.gz
ENV PATH="/usr/local/go/bin:$PATH"
ENV GOPATH=/go
ENV PATH=$PATH:$GOPATH/bin

# build onomyd artifact
WORKDIR $ONOMY_DIR
RUN make install


# ------------------ contract-deployer ------------------

ADD https://github.com/onomyprotocol/cosmos-gravity-bridge/releases/download/v0.0.0-20210915184851-orch-nomarket/Gravity.json /root/home/solidity/Gravity.json
ADD https://github.com/onomyprotocol/cosmos-gravity-bridge/releases/download/v0.0.0-20210915184851-orch-nomarket/TestERC20A.json /root/home/solidity/TestERC20A.json
ADD https://github.com/onomyprotocol/cosmos-gravity-bridge/releases/download/v0.0.0-20210915184851-orch-nomarket/contract-deployer /usr/bin/contract-deployer
RUN chmod +x /usr/bin/contract-deployer

# ------------------ orchestrator ------------------

ADD https://github.com/onomyprotocol/cosmos-gravity-bridge/releases/download/v0.0.0-20210915184851-orch-nomarket/gbt /usr/bin/gbt
RUN chmod +x /usr/bin/gbt

# ------------------ ethereum go (geth) ------------------
RUN add-apt-repository -y ppa:ethereum/ethereum
RUN apt-get install --no-install-recommends -yq ethereum

# ------------------ soft to run the scripts ------------------
RUN apt-get install --no-install-recommends -yq jq \
                                                moreutils

# ------------------ install fauset ------------------

RUN curl https://get.starport.network/faucet! | bash

# ------------------ run node ------------------

WORKDIR /root/home
RUN cp -r $ONOMY_DIR/deploy/single-node-runner-rinkeby/scripts scripts
RUN mkdir onomy
RUN cp -r $ONOMY_DIR/deploy/single-node-runner-rinkeby/onomy_home onomy/onomy

CMD /bin/bash scripts/run.sh