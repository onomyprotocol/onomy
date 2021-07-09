FROM golang:1.16-alpine AS build-env

# Set up dependencies
ENV PACKAGES curl make git libc-dev bash gcc linux-headers eudev-dev python3 git

RUN apk add --no-cache $PACKAGES

ARG GRAVITY_DIR=/go/src/github.com/onomyprotocol/cosmos-gravity-bridge

# clone the repo
RUN git clone https://github.com/onomyprotocol/cosmos-gravity-bridge.git $GRAVITY_DIR

# set working dir
WORKDIR $GRAVITY_DIR/module

# Build gravity artifact
RUN make install

RUN ls build
# Final image
FROM alpine:edge

# Install ca-certificates
RUN apk add --update ca-certificates git jq
WORKDIR /root

ARG GRAVITY_DIR=/go/src/github.com/onomyprotocol/cosmos-gravity-bridge

# Copy over binaries from the build-env
COPY --from=build-env $GRAVITY_DIR/module/build/gravity /usr/bin/gravity
EXPOSE 26656 26657 1317 9090
# FROM leeway321/gravity-bridge-binary:latest AS builder

ARG CHAIN_ID="gravity-testnet"
ARG KEYRING="--keyring-backend test"
ARG GRAVITY_HOST="0.0.0.0"
ARG GRAVITY_NODE_CONFIG=".gravity/config/config.toml"
ARG GRAVITY_APP_CONFIG=".gravity/config/app.toml"
ARG PUBLIC_IP=${PUBLIC_IP}
ARG GIT_HUB_USER=${GIT_HUB_USER}
ARG GIT_HUB_PASS=${GIT_HUB_PASS}
ARG GIT_HUB_EMAIL=${GIT_HUB_EMAIL}
ARG GIT_HUB_BRANCH=${GIT_HUB_BRANCH}
# init one fullnode
RUN gravity init --chain-id=$CHAIN_ID fullnode

#RUN gravity keys add fullnode $KEYRING
RUN gravity keys add --output=json fullnode $KEYRING > fullnode.json
RUN gravity add-genesis-account $(gravity keys show fullnode -a $KEYRING) 10000000stake
RUN gravity eth_keys add --output=text --dry-run=true | grep address: | sed 's/address://g' > eth_key

# Change configuration in the toml file
RUN sed -i "s#\"tcp://127.0.0.1:26656\"#\"tcp://$GRAVITY_HOST:26656\"#g" $GRAVITY_NODE_CONFIG
RUN sed -i "s#\"tcp://127.0.0.1:26657\"#\"tcp://$GRAVITY_HOST:26657\"#g" $GRAVITY_NODE_CONFIG
RUN sed -i 's#addr_book_strict = true#addr_book_strict = false#g' $GRAVITY_NODE_CONFIG
RUN sed -i 's#external_address = ""#external_address = "tcp://'$GRAVITY_HOST:26656'"#g' $GRAVITY_NODE_CONFIG
RUN sed -i 's#enable = false#enable = true#g' $GRAVITY_APP_CONFIG
RUN sed -i 's#swagger = false#swagger = true#g' $GRAVITY_APP_CONFIG


##----------------------------- commit master genesis,json file-----
RUN git clone -b $GIT_HUB_BRANCH https://github.com/sunnyk56/onomy.git
WORKDIR /root/onomy
CMD sh cosmosNode/fullNode/start.sh $GIT_HUB_USER $GIT_HUB_PASS $GIT_HUB_EMAIL $GIT_HUB_BRANCH
