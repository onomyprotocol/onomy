FROM ubuntu:latest

RUN apt-get update -y -q && apt-get upgrade -yq
# common (DEBIAN_FRONTEND is a fix for tzdata)
RUN DEBIAN_FRONTEND="noninteractive" apt-get install --no-install-recommends -yq software-properties-common \
                                                curl \
                                                build-essential \
                                                ca-certificates \
                                                tar \
                                                git

ARG GRAVITY_DIR=/go/src/github.com/onomyprotocol/gravity-bridge
# clone the repo //TODO this clone may be change in future
# I have tried using repo https://github.com/althea-net/cosmos-gravity-bridge.git but there is no
# command make build-linux found in makesfiles, so we forked anddd thhis command
#RUN git clone https://github.com/althea-net/cosmos-gravity-bridge.git $GRAVITY_DIR
RUN git clone https://github.com/ChaudharyShiva18/cosmos-gravity-bridge.git $GRAVITY_DIR


# ------------------ gravity module ------------------
# install golang
RUN curl https://dl.google.com/go/go1.16.4.linux-amd64.tar.gz --output go.tar.gz
RUN tar -C /usr/local -xzf go.tar.gz
ENV PATH="/usr/local/go/bin:$PATH"
ENV GOPATH=/go
ENV PATH=$PATH:$GOPATH/bin

# build gravity artifact
RUN pwd
RUN ls
WORKDIR $GRAVITY_DIR/module
RUN pwd
RUN ls
# Build gravity artifact
RUN make build-linux
# copy binary to bin
RUN cp $GRAVITY_DIR/module/build/gravity /usr/bin/gravity

# ------------------ orchestrator ------------------
# install rust
RUN curl https://sh.rustup.rs -sSf | bash -s -- -y
ENV PATH="/root/.cargo/bin:${PATH}"
# check cargo works
RUN cargo version

# build orchestrator artifact
RUN pwd
RUN ls
WORKDIR $GRAVITY_DIR/orchestrator
RUN pwd
RUN ls
RUN rustup target add x86_64-unknown-linux-musl
RUN apt-get install --no-install-recommends -yq musl-tools
RUN cargo build --target=x86_64-unknown-linux-musl --release  --all
RUN cp $GRAVITY_DIR/orchestrator/target/x86_64-unknown-linux-musl/release/orchestrator /usr/bin/orchestrator

# ------------------ Solidity ------------------
RUN apt-get install --no-install-recommends -yq python3 \
                                                make \
                                                g++

RUN curl -sL https://deb.nodesource.com/setup_15.x  | bash -
RUN apt-get -y install nodejs
# compile contract
RUN pwd
RUN ls
WORKDIR $GRAVITY_DIR/solidity
RUN pwd
RUN ls
RUN npm ci
RUN chmod -R +x scripts
RUN npm run typechain

# ------------------ ethereum go (geth) ------------------
RUN add-apt-repository -y ppa:ethereum/ethereum
RUN apt-get install --no-install-recommends -yq ethereum
# ------------------ soft to run the scripts ------------------
RUN apt-get install --no-install-recommends -yq jq \
                                                moreutils

# go to root home
RUN pwd
RUN ls
#WORKDIR /root/home
WORKDIR /root
RUN pwd
RUN ls
CMD ["bash"]
