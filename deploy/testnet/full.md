# Steps to run the full node

* Install dependencies from source code

```
./bin.sh
```

* Init full node

```
./init-full-node.sh
```

* Expose monitoring

```
./expose-metrics.sh
```

* Optionally allow cors requests

```
./allow-cors.sh
```

* Run node exporter

```
./start-node-exporter.sh
```

* Start the node

Before run the script please set up "ulimit > 65535" ([example-hrel8](set-ulimit-hrel8.md))

```
./start-node.sh
```

* Setup auto-start

Add to your crontab or /etc/init.d scripts:

* `start-node.sh`
* `start-node-exporter.sh`

***If you used the bin.sh installation and what to use the scripts for the auto-start, additionally you need to
add ```export PATH=$PATH:$ONOMY_HOME/bin``` to your scripts after the ```ONOMY_HOME=$HOME/.onomy```***
