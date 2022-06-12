FROM ubuntu:latest

RUN apt-get update -y -q && apt-get upgrade -yq
# common (DEBIAN_FRONTEND is a fix for tzdata)
RUN DEBIAN_FRONTEND="noninteractive" apt-get install --no-install-recommends -yq software-properties-common \
                                                curl \
                                                build-essential \
                                                ca-certificates \
                                                tar \
                                                git \
                                                crudini
# scripts utils
RUN apt-get install --no-install-recommends -yq jq \
                                                iputils-ping \
                                                moreutils

# ------------------go ------------------

# install golang
RUN curl https://dl.google.com/go/go1.16.4.linux-amd64.tar.gz --output go.tar.gz
RUN tar -C /usr/local -xzf go.tar.gz
ENV PATH="/usr/local/go/bin:$PATH"
ENV GOPATH=/go
ENV PATH=$PATH:$GOPATH/bin

# ------------------ nodejs ------------------

RUN curl -sL https://deb.nodesource.com/setup_15.x  | bash -
RUN apt-get -y install nodejs

# ------------------ ethereum go (geth) ------------------

RUN add-apt-repository -y ppa:ethereum/ethereum
RUN apt-get install --no-install-recommends -yq ethereum