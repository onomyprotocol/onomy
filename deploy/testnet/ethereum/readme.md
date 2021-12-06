# Steps to run the ethereum node

## Install dependencies from source code

```
./bin.sh
```

## Start the node

Before run the script please set up "ulimit > 65535":

```
./start-node.sh
```

## Setup auto-start

Add start-node.sh to your crontab or /etc/init.d in order to start automatically after the OS restart.

If you used the bin.sh installation then additionally you need to add 
```
export PATH=$PATH:$ONOMY_HOME/bin
```

In your start scripts (after the ONOMY_HOME initialization)