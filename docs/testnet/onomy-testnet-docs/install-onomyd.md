# Installing onomy
The operating system you use for your node is entirely your personal preference. You will be able to compile the `onomyd` daemon on most modern linux distributions and recent versions of macOS.

For the tutorial, it is assumed that you are using an Ubuntu LTS release.

If you have chosen a different operating system, you will need to modify your commands to suit your operating system.

## Requirements

A Linux server with any modern Linux distribution, 2gb of ram and at least 20gb storage. Requirements are very minimal.

### Install pre-requisites
```bash:
# update the local package list and install any available upgrades

sudo apt-get update && sudo apt upgrade -y

# install toolchain

sudo apt-get install make build-essential gcc git jq -y
```

### Install Go
Go is required to build `onomy` from source.
```bash:
rm -rf /usr/local/go && tar -C /usr/local -xzf go1.17.1.linux-amd64.tar.gz
```

Please install Go v1.17 or later.

If you are in any way unsure about how to configure Go, then set these in the `.profile` in the user's home (i.e. `~/`) folder.

```bash:
# Updates environmental variables to include go
cat <<EOF>> ~/.profile
export GOROOT=/usr/local/go
export GOPATH=$HOME/go
export GO111MODULE=on
export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin
EOF
source ~/.profile
```

### Build Onomy from source
```
git clone https://github.com/onomyprotocol/onomy
cd onomy
git fetch
git checkout <version-tag>
```
Where `<version-tag>` is the current version (currently: `v0.0.1`).

Then install `onomy`:
```bash:
make install
```

Verify installation:
```bash:
onomyd version
```

Your output should be `<version-tag>`. 

## Next Steps

If you were successful, you can now [install a full node](onomy-testnet-docs/setting-up-a-fullnode-manual.md).
