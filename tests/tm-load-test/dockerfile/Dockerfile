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

#----------------------building binary of tm-load-test custom client-------------------
RUN git clone -b customClient https://github.com/puneetsingh166/tm-load-test.git
RUN ls
WORKDIR tm-load-test/customclient/my-cosmos-tester
RUN go build
RUN cp my-cosmos-tester /usr/bin/my-cosmos-tester
RUN my-cosmos-tester --help


# ------------------ soft to run the scripts ------------------
RUN apt-get install --no-install-recommends -yq jq \
                                                moreutils
# ------------------ run node ------------------

WORKDIR /root/home
RUN cp -r $ONOMY_DIR/deploy/testnet/master-validator master-validator
WORKDIR /root/home
RUN  cp -r $ONOMY_DIR/tests/tm-load-test/run-load-test run-load-test  


 

 
