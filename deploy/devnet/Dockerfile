FROM onomy/tests-base:latest

# ------------------ onomy chain ------------------

ARG ONOMY_DIR=/go/src/github.com/onomyprotocol/onomy/
# clone the repo
COPY . $ONOMY_DIR

# build onomyd artifact
WORKDIR $ONOMY_DIR
RUN make install

WORKDIR /root/home
RUN cp -r $ONOMY_DIR/deploy/devnet/scripts scripts
RUN mkdir onomy
RUN cp -r $ONOMY_DIR/deploy/devnet/.onomy $HOME/.onomy

CMD /bin/bash scripts/run.sh