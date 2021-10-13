
In theory, Onomy chain can be run on Windows and Mac. Binaries will be provided on the releases page and currently scripts files are provided to make binaries.
I also suggest an open notepad or other document to keep track of the keys you will be generating.

## Bootstrapping steps and commands

Start by logging into your Linux server using ssh. The following commands are intended to be run on that machine. There are three options to install the binaries
1. [Download them directly from github and install](#downloadInstall)
2. [Compile binaries yourself using source code](#compileInstall)
3. [Install an RPM Package (Requires Fedore or CentOS machine)][#rpmInstall]

### <a name="downloadInstall"></a> 1. Download/install Onomy chain binaries 
To download and install binaries follow these steps

1. create a new directory in your home directory which will save all the onomy packages. `mkdir -p $HOME/.onomy/bin && cd $HOME/.onomy/bin` 
2. Download binaries using wget and add executable permission
```
wget https://github.com/onomyprotocol/onomy/releases/download/v0.0.1/onomyd
wget https://github.com/onomyprotocol/onomy/releases/download/v0.0.1/gbt
wget https://github.com/onomyprotocol/onomy/releases/download/v0.0.1/geth
chmod +x *
```
3. You can now use these binaries, but in order to use them from anywhere in your terminal, you will need add them to $PATH variable
`
export PATH=$PATH:$HOME/.onomy/bin
`

### <a name="compileInstall"></a> 2. Compile Onomy chain binaries 
If you have Fedora 34 or Red Hat Enterprise Linux 8 and you want to make binaries yourself, then follow these steps
1. Clone Onomy repo. (You might need to install git using `dnf install git`).
```
git clone -b dev https://github.com/onomyprotocol/onomy.git
```
2. go to deploy/testnet and run bin.sh script file
```
cd onomy/deploy/testnet
bash bin.sh
```
The second way may be unsafe because it used the latest version of the artifacts.

At specific points during the testnet you may be told to `update your orchestrator` or `update your onomyd binary`. In order to do that you can simply repeat the above instructions and then restart the affected software.

to check what version of the tools you have run `gbt --version`

### <a name="rpmInstall"></a> 3. Install RPM package

In order to install using RPM, use the following command:
```
sudo yum install "https://github.com/onomyprotocol/onomy/releases/download/v0.0.1/onomy-0.1.0-1x86_64.rpm"
```

Running this command will take care of all the dependancies and it will install the required binaries in `$HOME/.onomy/bin` directory. It will also add this path to $PATH variable.
