# Steps to run the sentry node

* Install dependencies from source code

```
./bin.sh
```

* Init sentry node

```
./init-sentry-node.sh
```

* Optionally expose monitoring

```
./expose-metrics.sh
```

* Start the node

Before running the script please set up "ulimit > 65535" ([Red Hat Enterprise Linux](set-ulimit-rhel8.md))

```
./start-node.sh
```

* Optionally run node exporter

```
./start-node-exporter.sh
```

* Start the node

Before running the script please set up "ulimit > 65535" ([Red Hat Enterprise Linux](set-ulimit-rhel8.md))

```
./start-node.sh
```

* Setup auto-start

Add to your crontab or /etc/init.d scripts:

* `start-node.sh`
* `start-node-exporter.sh`

***If you used the bin.sh installation and want to use the scripts for the auto-start, additionally you need to
add ```export PATH=$PATH:$ONOMY_HOME/bin``` to your scripts after the ```ONOMY_HOME=$HOME/.onomy```***