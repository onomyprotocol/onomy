FROM ubuntu:21.04

RUN apt-get update -y -q && apt-get upgrade -yq
# common (DEBIAN_FRONTEND is a fix for tzdata)
RUN DEBIAN_FRONTEND="noninteractive" apt-get install --no-install-recommends -yq software-properties-common \
                                                curl \
                                                build-essential \
                                                ca-certificates \
                                                tar \
                                                git
# scripts utils
RUN apt-get install --no-install-recommends -yq jq \
                                                moreutils

# install golang
RUN curl https://dl.google.com/go/go1.16.4.linux-amd64.tar.gz --output go.tar.gz
RUN tar -C /usr/local -xzf go.tar.gz
ENV PATH="/usr/local/go/bin:$PATH"
ENV GOPATH=/go
ENV PATH=$PATH:$GOPATH/bin

RUN GO111MODULE=on go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.41.1
RUN GO111MODULE=on go get github.com/vasi-stripe/gogroup/cmd/gogroup@v0.0.0-20200806161525-b5d7f67a97b5
RUN GO111MODULE=on go get mvdan.cc/gofumpt@v0.0.0-20200927160801-5bfeb2e70dd6
RUN GO111MODULE=on go get github.com/bufbuild/buf/cmd/buf@v0.56.0

RUN curl https://get.starport.network/starport@v0.19.3! | bash

RUN rm -rf /root/.cache/go-build/ /go/pkg/*
COPY entrypoint.sh /entrypoint.sh

RUN mkdir -p /root/.cache/ && \
    ln -s /cache/golangci-lint/ /root/.cache/golangci-lint && \
    ln -s /cache/go-build/ /root/.cache/go-build

WORKDIR /go/src/github.com/onomyprotocol/onomy

ENTRYPOINT ["bash", "/entrypoint.sh"]