# How to run a Onomy testnet full node

A Onomy chain full node is just like any other Cosmos chain full node and unlike the validator flow requires no external software

## What do I need?

A Linux server with any modern Linux distribution, 4cores, 8gb of ram and at least 20gb of SSD storage.

In theory, Onomy chain can be run on Windows and Mac. Binaries will be provided on the releases page and currently, scripts files are provided to make binaries.
I also suggest an open notepad or other document to keep track of the keys you will be generating.

## Bootstrapping steps and commands

Start by logging into your Linux server using ssh. The following commands are intended to be run on that machine

### Download/install Onomy chain binaries
```
To download binary follow these commands
mkdir binaries
cd binaries
wget https://github.com/onomyprotocol/onomy/releases/download/v0.0.1/onomyd
wget https://github.com/onomyprotocol/onomy/releases/download/v0.0.1/gbt
wget https://github.com/onomyprotocol/onomy/releases/download/v0.0.1/geth
cd ..
chmod -R +x binaries
export PATH=$PATH:$HOME/binaries/


or If you have Fedora (Fedora 34) or Redhat (Red Hat Enterprise Linux 8.4 (Ootpa))
 and you want to make binaries yourself, then follow these steps

sudo yum install -y git
git clone -b dev https://github.com/onomyprotocol/onomy.git
cd onomy/deploy/testnet
bash bin.sh
```

### Initiate chain

```
git clone -b dev https://github.com/onomyprotocol/onomy.git
cd onomy/deploy/testnet
```

### Run the first time bootstrapping playbook and script

This script will run commands to generate keys and also store in files. You will need them later.

```
bash peer-validator/init.sh

Note: 1. Script will ask for enter validator name(Type any name for example validator1)
      2. Script will ask for enter node-id of any validator that is running in chain to add seed(Please enter 1302d0ed290d74d6f061fb8506e0e34f3f67f7ff)
      3. Script will ask for enter ip of validator for which you have added node-id(Please enter 147.182.128.38)
```

Now it's finally up and start the sycing block

### Check the status of the Onomy chain

You should be good to go! You can check the status of the three
Onomy chain by running.
```
curl http://localhost:26657/status
```
if catching_up is false means your node is fully synced

