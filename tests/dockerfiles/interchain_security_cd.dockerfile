FROM fedora:38
RUN dnf install -y git make cmake gcc gcc-c++ which iproute iputils procps-ng vim-minimal tmux net-tools htop tar jq npm openssl-devel perl rust cargo golang
RUN go install cosmossdk.io/tools/cosmovisor/cmd/cosmovisor@latest
ENV PATH=$PATH:/root/go/bin

# without the leading 'v'
ENV ICS_VERSION=2.0.0-rc1
ENV DAEMON_NAME="interchain-security-cd"
ENV DAEMON_HOME="/root/.interchain-security-c"
ENV CONSUMER_CURRENT_VERSION=v07-Theta

ADD https://github.com/cosmos/interchain-security/archive/refs/tags/v$ICS_VERSION.tar.gz /root/v$ICS_VERSION.tar.gz
RUN cd /root && tar -xvf ./v$ICS_VERSION.tar.gz
RUN cd /root/interchain-security-$ICS_VERSION && go build ./cmd/interchain-security-cd

RUN mkdir -p $DAEMON_HOME/cosmovisor/genesis/$CONSUMER_CURRENT_VERSION/bin/
RUN mv /root/interchain-security-$ICS_VERSION/interchain-security-cd $DAEMON_HOME/cosmovisor/genesis/$CONSUMER_CURRENT_VERSION/bin/$DAEMON_NAME

# for manual testing
RUN chmod +x $DAEMON_HOME/cosmovisor/genesis/$CONSUMER_CURRENT_VERSION/bin/$DAEMON_NAME

# set up symbolic links
RUN cosmovisor init $DAEMON_HOME/cosmovisor/genesis/$CONSUMER_CURRENT_VERSION/bin/$DAEMON_NAME

# some commands don't like if the data directory does not exist
RUN mkdir $DAEMON_HOME/data
