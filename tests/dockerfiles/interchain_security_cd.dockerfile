FROM fedora:38
RUN dnf install -y git make cmake gcc gcc-c++ which iproute iputils procps-ng vim-minimal tmux net-tools htop tar jq npm openssl-devel perl rust cargo golang
RUN go install cosmossdk.io/tools/cosmovisor/cmd/cosmovisor@latest
ENV PATH=$PATH:/root/go/bin

ENV DAEMON_NAME="interchain-security-cd"
ENV DAEMON_HOME="/root/.interchain-security-c"
ENV MARKET_CURRENT_VERSION=v07-Theta

# TODO
ADD ./dockerfile_resources/$DAEMON_NAME $DAEMON_HOME/cosmovisor/genesis/$MARKET_CURRENT_VERSION/bin/$DAEMON_NAME

# for manual testing
RUN chmod +x $DAEMON_HOME/cosmovisor/genesis/$MARKET_CURRENT_VERSION/bin/$DAEMON_NAME

# set up symbolic links
RUN cosmovisor init $DAEMON_HOME/cosmovisor/genesis/$MARKET_CURRENT_VERSION/bin/$DAEMON_NAME

# some commands don't like if the data directory does not exist
RUN mkdir $DAEMON_HOME/data
