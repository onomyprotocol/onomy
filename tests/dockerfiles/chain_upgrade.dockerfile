FROM fedora:38
RUN dnf install -y git make cmake gcc gcc-c++ which iproute iputils procps-ng vim-minimal tmux net-tools htop tar jq npm openssl-devel perl rust cargo golang
RUN go install cosmossdk.io/tools/cosmovisor/cmd/cosmovisor@latest
ENV PATH=$PATH:/root/go/bin

ENV DAEMON_NAME="onomyd"
ENV DAEMON_HOME="/root/.onomy"
ENV CURRENT_VERSION=v1.1.4
ENV UPGRADE_VERSION=v1.1.5dev
ENV REPO="onomyprotocol/onomy"

ADD https://github.com/$REPO/releases/download/$CURRENT_VERSION/$DAEMON_NAME $DAEMON_HOME/cosmovisor/genesis/$CURRENT_VERSION/bin/$DAEMON_NAME
ADD ./dockerfile_resources/$DAEMON_NAME $DAEMON_HOME/cosmovisor/upgrades/$UPGRADE_VERSION/bin/$DAEMON_NAME
#ADD https://github.com/$REPO/releases/download/$UPGRADE_VERSION/$DAEMON_NAME $DAEMON_HOME/cosmovisor/upgrades/$UPGRADE_VERSION/bin/$DAEMON_NAME
#ADD https://github.com/$REPO/releases/download/$UPGRADE_VERSION/$DAEMON_NAME /tmp/$DAEMON_NAME

# for manual testing
RUN chmod +x $DAEMON_HOME/cosmovisor/genesis/$CURRENT_VERSION/bin/$DAEMON_NAME
RUN chmod +x $DAEMON_HOME/cosmovisor/upgrades/$UPGRADE_VERSION/bin/$DAEMON_NAME

# set up symbolic links
RUN cosmovisor init $DAEMON_HOME/cosmovisor/genesis/$CURRENT_VERSION/bin/$DAEMON_NAME

# some commands don't like if the data directory does not exist
RUN mkdir $DAEMON_HOME/data
