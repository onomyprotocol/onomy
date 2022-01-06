# Steps to run the full node

* Make the scripts executable

```
chmod +x *
```

* Install dependencies from source code

```
./bin.sh
```

* Init full node

```
./init-full-node.sh
```

* Optionally expose monitoring

```
./expose-metrics.sh
```

* Optionally allow cors requests

```
./allow-cors.sh
```

* Start the node

Before running the script please set up "ulimit > 65535" ([Red Hat Enterprise Linux](set-ulimit-rhel8.md))

```
./start-onomyd.sh &>> $HOME/.onomy/logs/onomyd.log &
```