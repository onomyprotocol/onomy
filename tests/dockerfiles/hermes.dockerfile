FROM fedora:38
RUN dnf install -y git make cmake gcc gcc-c++ which iproute iputils procps-ng vim-minimal tmux net-tools htop tar jq npm openssl-devel perl rust cargo golang
RUN go install cosmossdk.io/tools/cosmovisor/cmd/cosmovisor@latest
ENV PATH=$PATH:/root/go/bin
ADD https://github.com/informalsystems/hermes/releases/download/v1.5.1/hermes-v1.5.1-x86_64-unknown-linux-gnu.tar.gz /root/.hermes/bin/
RUN cd /root/.hermes/bin/ && tar -vxf *
ENV PATH=$PATH:/root/.hermes/bin

ADD ./dockerfile_resources/hermes_config.toml /root/.hermes/config.toml
