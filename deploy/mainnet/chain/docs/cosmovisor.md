# Cosmovisor

cosmovisor is a small process manager for Cosmos SDK application binaries that monitors the governance module for
incoming chain upgrade proposals. If it sees a proposal that gets approved, cosmovisor can automatically download the
new binary, stop the current binary, switch from the old binary to the new one, and restart the node with the new
binary.

Cosmovisor is designed to be used as a wrapper for onomyd:

* it will pass arguments to the onomy app (configured by DAEMON_NAME env variable). Running `cosmovisor run arg1 arg2`
  is same as running `onomyd arg1 arg2`
* it will manage the app by restarting and upgrading if needed

* Installation If you set up your node using bin.sh script, it has already installed cosmovisor for you. If you want to
  install it manually you can clone the onomy-sdk repo and compile cosmovisor.

```
git clone https://github.com/onomyprotocol/onomy-sdk
cd cosmovisor
make cosmovisor
```

cosmovisor binRY will be saved in cosmovisor directory, you can copy it to anywhere you like, ideally
to `$ONOMY_HOME/bin`

```
cp cosmovisor/cosmovisor $ONOMY_HOME/bin/
```

## Setting up cosmovisor

Once installed, you need to setup cosmovisor directory structure and some environment variable for it to be ready to
run. Following environment variables are used:

* `DAEMON_HOME` is the location where the cosmovisor/ directory is kept that contains the genesis binary, the upgrade
  binaries, and any additional auxiliary files associated with each binary (e.g. $ONOMY_HOME).
* `DAEMON_NAME` is the name of the binary itself (e.g. onomyd).

Once environment variables are set, you need to create a directory tree with following structure inside `$DAEMON_HOME`
directory:

```
.
├── current -> genesis or upgrades/<name> 	# This link is created by cosmovisor
├── genesis
│   └── bin
│       └── $DAEMON_NAME
└── upgrades
```

copy onomyd binary from `$ONOMY_HOME/bin` to `DAEMON_HOME/cosmovisor/genesis/bin/`.

## Running full node using cosmovisor

Run cosmovisor using `cosmovisor run` command. To start the full node run

```
cosmovisor run start
```

which is equivalent to

```
onomyd start
```

When cosmovisor is running, it will detect any queued upgrades and will switch binary from `genesis/bin/onomyd`
to `upgrades/<upgrade-name>/bin/onomyd`.

For more information on cosmovisor, please
visit [Cosmovisor](https://github.com/onomyprotocol/onomy-sdk/tree/master/cosmovisor)
