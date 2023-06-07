FROM fedora:38
RUN dnf install -y git make cmake gcc gcc-c++ which iproute iputils procps-ng vim-minimal tmux net-tools htop tar jq npm openssl-devel perl rust cargo golang
RUN go install cosmossdk.io/tools/cosmovisor/cmd/cosmovisor@latest
ENV PATH=$PATH:/root/go/bin

ENV DAEMON_NAME="onomyd"
ENV DAEMON_HOME="/root/.onomy"
# the previous version
ENV ONOMY_CURRENT_VERSION=v1.0.3.5
# the version that currently is implemented by this repository's state
ENV ONOMY_UPGRADE_VERSION=v1.1.0

# note that one has to go under `genesis/` and the other under `upgrades/`
ADD https://github.com/onomyprotocol/onomy/releases/download/$ONOMY_CURRENT_VERSION/$DAEMON_NAME $DAEMON_HOME/cosmovisor/genesis/$ONOMY_CURRENT_VERSION/bin/$DAEMON_NAME
ADD ./dockerfile_resources/$DAEMON_NAME $DAEMON_HOME/cosmovisor/upgrades/$ONOMY_UPGRADE_VERSION/bin/$DAEMON_NAME

# for manual testing
RUN chmod +x $DAEMON_HOME/cosmovisor/genesis/$ONOMY_CURRENT_VERSION/bin/$DAEMON_NAME
RUN chmod +x $DAEMON_HOME/cosmovisor/upgrades/$ONOMY_UPGRADE_VERSION/bin/$DAEMON_NAME

# set up symbolic links
RUN cosmovisor init $DAEMON_HOME/cosmovisor/genesis/$ONOMY_CURRENT_VERSION/bin/$DAEMON_NAME

# some commands don't like if the data directory does not exist
RUN mkdir $DAEMON_HOME/data
