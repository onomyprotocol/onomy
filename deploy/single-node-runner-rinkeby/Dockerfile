FROM ubuntu:latest

RUN apt-get update -y -q && apt-get upgrade -yq
# common (DEBIAN_FRONTEND is a fix for tzdata)
RUN DEBIAN_FRONTEND="noninteractive" apt-get install --no-install-recommends -yq software-properties-common \
                                                curl \
                                                build-essential \
                                                ca-certificates \
                                                tar \
                                                git

# ------------------ clone repository ------------------

ARG ONOMY_DIR=/go/src/github.com/onomyprotocol/onomy/
# clone the repo
COPY . $ONOMY_DIR

# ------------------ onomy chain ------------------

# install golang
RUN curl https://dl.google.com/go/go1.16.4.linux-amd64.tar.gz --output go.tar.gz
RUN tar -C /usr/local -xzf go.tar.gz
ENV PATH="/usr/local/go/bin:$PATH"
ENV GOPATH=/go
ENV PATH=$PATH:$GOPATH/bin

# build onomyd artifact
WORKDIR $ONOMY_DIR
RUN make build
# copy binary to bin
RUN cp onomyd /usr/bin/onomyd

# ------------------ gravity repo ------------------

ARG GRAVITY_DIR=/go/src/github.com/onomyprotocol/cosmos-gravity-bridge/
# clone the repo
RUN mkdir $GRAVITY_DIR
RUN git clone https://github.com/onomyprotocol/cosmos-gravity-bridge.git $GRAVITY_DIR

# ------------------ orchestrator ------------------

# install rust
RUN curl https://sh.rustup.rs -sSf | bash -s -- -y
ENV PATH="/root/.cargo/bin:${PATH}"
# check cargo works
RUN cargo version

# build orchestrator artifact
WORKDIR $GRAVITY_DIR/orchestrator
RUN rustup target add x86_64-unknown-linux-musl
RUN apt-get install --no-install-recommends -yq musl-tools
RUN cargo build --target=x86_64-unknown-linux-musl --release  --all
# copy binary to bin
RUN ls target/x86_64-unknown-linux-musl
RUN ls target/x86_64-unknown-linux-musl/release
RUN cp target/x86_64-unknown-linux-musl/release/gbt /usr/bin/gbt

# ------------------ contract-deployer ------------------

RUN apt-get install --no-install-recommends -yq python3 \
                                                make \
                                                g++

RUN curl -sL https://deb.nodesource.com/setup_15.x  | bash -
RUN apt-get -y install nodejs

# build contract-deployer
WORKDIR $GRAVITY_DIR/solidity
RUN npm ci
RUN chmod -R +x scripts
RUN npm run typechain

# go to root home
WORKDIR /root/home

# ------------------ ethereum go (geth) ------------------
RUN add-apt-repository -y ppa:ethereum/ethereum
RUN apt-get install --no-install-recommends -yq ethereum

# ------------------ soft to run the scripts ------------------
RUN apt-get install --no-install-recommends -yq jq \
                                                moreutils

# ------------------ install fauset ------------------

RUN curl https://get.starport.network/faucet! | bash

# ------------------ run node ------------------

RUN cp -r $ONOMY_DIR/deploy/single-node-runner-rinkeby/scripts scripts
RUN mkdir onomy
RUN cp -r $ONOMY_DIR/deploy/single-node-runner-rinkeby/onomy_home onomy/onomy

CMD /bin/bash scripts/run.sh