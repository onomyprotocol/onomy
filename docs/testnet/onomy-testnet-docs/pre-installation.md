
In theory, Onomy chain can be run on Windows and Mac. Binaries will be provided on the releases page and currently, scripts files are provided to make binaries.
I also suggest an open notepad or other document to keep track of the keys you will be generating.

## Bootstrapping steps and commands

Start by logging into your Linux server using ssh. The following commands are intended to be run on that machine

### Download/install Onomy chain binaries
To download and install binaries follow these steps

1. create a new directory to download binaries in and open that directory `mkdir binaries && cd binaries` 
2. Download binaries using wget and add executable permission
```
wget https://github.com/onomyprotocol/onomy/releases/download/v0.0.1/onomyd
wget https://github.com/onomyprotocol/onomy/releases/download/v0.0.1/gbt
wget https://github.com/onomyprotocol/onomy/releases/download/v0.0.1/geth
chmod +x *

```
3. You can now use these binaries, but in order to use them from anywhere in your terminal, you will need add them to $PATH variable
`
export PATH=$PATH:$HOME/binaries/
`

or If you have Fedora (Fedora 34) or Redhat (Red Hat Enterprise Linux 8.4 (Ootpa))
and you want to make binaries yourself, then follow these steps
```
sudo yum install -y git
git clone -b dev https://github.com/onomyprotocol/onomy.git
cd onomy/deploy/testnet
bash bin.sh
```
The second way may be unsafe because it used the latest version of the artifacts.

At specific points during the testnet you may be told to `update your orchestrator` or `update your onomyd binary`. In order to do that you can simply repeat the above instructions and then restart the affected software.

to check what version of the tools you have run `gbt --version`
