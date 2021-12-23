FROM golang:1.16-alpine


RUN apk update  && apk upgrade 
# common (DEBIAN_FRONTEND is a fix for tzdata)
#install bash
RUN apk add --no-cache bash
RUN apk add --no-cache  curl gcc g++ make ca-certificates tar git

# ------------------ onomy chain ------------------

ARG ONOMY_DIR=/go/src/github.com/onomyprotocol/onomy/
# clone the repo
COPY . $ONOMY_DIR

# build onomyd artifact
WORKDIR $ONOMY_DIR
RUN make install

#----------------------building binary of tm-load-test custom client-------------------
WORKDIR $ONOMY_DIR
RUN make build-load-test
WORKDIR $ONOMY_DIR/build
RUN cp onomy-load-test /usr/bin/onomy-load-test
RUN onomy-load-test --help


# ------------------ soft to run the scripts ------------------
RUN apk add --no-cache  jq moreutils
# ------------------ run node ------------------

WORKDIR /root/home

RUN bash $ONOMY_DIR/tests/tm-load-test/init-master.sh
WORKDIR $ONOMY_DIR/tests/tm-load-test


 

 
