# Installation

In theory, Onomy chain can be run on Windows and Mac. Binaries will be provided on the releases page and currently
scripts files are provided to make binaries. We also suggest an open notepad or other document to keep track of the keys
you will be generating.

## Bootstrapping steps and commands

Start by logging into your Linux server. The following commands are intended to be run on that machine. 

1. [Install binaries using scripts](#installWithScripts)
2. [Install compiled binaries](#installCompiled)

### <a name="installWithScripts"></a> 1. Install binaries using scripts.

1. Clone Onomy repo. (You might need to install git using `dnf install git`).

```
git clone https://github.com/onomyprotocol/onomy.git
```

2. Run the installation script

```
cd onomy/deploy/scripts
```

* Install dependencies from compiled binaries

  For mainnet
    ```
    ./onomy/deploy/scripts/bin-mainnet.sh
    ```
  For testnet
    ```
    ./onomy/deploy/scripts/bin-testnet.sh
    ```

* Or compile and install

  For mainnet
    ```
    ./onomy/deploy/scripts/bin-mainnet-from-sources.sh
    ```
  For testnet
    ```
    ./onomy/deploy/scripts/bin-testnet-from-sources.sh
    ```

### <a name="installCompiled"></a> 2. Install compiled binaries

To download and install binaries follow these steps

1. Create a new directory in your home directory which will save all the onomy
   packages. `mkdir -p $HOME/.onomy/bin && cd $HOME/.onomy/bin`
2. Download binaries using wget and add executable permission

* Create new bin dir

```
cd $HOME/.onomy/bin
```

* Download binaries

For mainnet

   ```
   wget https://github.com/onomyprotocol/onomy/releases/download/v1.0.0/onomyd
   ```

For testnet

   ```
   wget https://github.com/onomyprotocol/onomy/releases/download/v0.1.0/onomyd
   ```

* Make binaries executable

```
chmod +x *
```

3. You can now use these binaries, but in order to use them from anywhere in your terminal, you will need add them to
   $PATH variable
   `
   export PATH=$PATH:$HOME/.onomy/bin
   `


