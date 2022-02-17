FROM onomy/tests-base:latest

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

WORKDIR /root/home

RUN bash $ONOMY_DIR/tests/tm-load-test/init-master.sh
WORKDIR $ONOMY_DIR/tests/tm-load-test
 
