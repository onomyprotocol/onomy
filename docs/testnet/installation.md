# Installation

In theory, Onomy chain can be run on Windows and Mac. Binaries will be provided on the releases page and currently
scripts files are provided to make binaries. We also suggest an open notepad or other document to keep track of the keys
you will be generating.

## Bootstrapping steps and commands

Start by logging into your Linux server. The following commands are intended to be run on that machine. There are three
options to install the binaries:

1. [Compile binaries yourself using source code](#compileInstall)
2. [Download them directly from github and install](#downloadInstall)
3. [Install an RPM Package (Requires Fedore or CentOS machine)](#rpmInstall)

### <a name="compileInstall"></a> 1. Compile Onomy chain binaries

if you want to compile binaries yourself, then follow these steps

1. Clone Onomy repo. (You might need to install git using `dnf install git`).

```
git clone https://github.com/onomyprotocol/onomy.git
```

2. go to deploy/testnet/scripts and run bin.sh script file

```
cd onomy/deploy/testnet/scripts
sh bin.sh
```

### <a name="downloadInstall"></a> 2. Download/install Onomy chain binaries

To download and install binaries follow these steps

1. Create a new directory in your home directory which will save all the onomy
   packages. `mkdir -p $HOME/.onomy/bin && cd $HOME/.onomy/bin`
2. Download binaries using wget and add executable permission

```
cd $HOME/.onomy/bin
wget https://github.com/onomyprotocol/onomy/releases/download/v0.0.4/onomyd
chmod +x *
```

3. You can now use these binaries, but in order to use them from anywhere in your terminal, you will need add them to
   $PATH variable
   `
   export PATH=$PATH:$HOME/.onomy/bin
   `

### <a name="rpmInstall"></a> 3. Install RPM package

If you are using a Red Hat or Fedora based OS, you can use following RPM to install the binaries and download required
script files.

In order to install using RPM, use the following command:

```
sudo yum install "https://github.com/onomyprotocol/onomy/releases/download/v0.0.4/onomy.x86_64.rpm"
```

Running this command will take care of all the dependancies and it will install the required binaries
in `$HOME/.onomy/bin` directory. It will also add this path to $PATH variable.
