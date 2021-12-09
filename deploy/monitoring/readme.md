# Steps set up and start the monitoring

## Expose the node metrics

If you didn't run the expose-metrics.sh before run it:

```
./expose-metrics.sh
```

This script will enable the prometheus metrics in your node config.

## Expose the OS metrics

In order to expose the OS metrics you need to install and run the "Node Exporter"

Installation instruction:

```
curl -LO https://github.com/prometheus/node_exporter/releases/download/v0.18.1/node_exporter-0.18.1.linux-amd64.tar.gz
tar -xvf node_exporter-0.18.1.linux-amd64.tar.gz
mv node_exporter-0.18.1.linux-amd64/node_exporter /bin/
node_exporter --version
```

Run script that runs node-exporter

```
./start-node-exporter.sh
```

Open the http://<server-IP>:9100/metrics to check that service is available.

## Setup auto-start

Add start-node-exporter.sh to your crontab or /etc/init.d in order to start automatically after the OS restart.

## If the node is running now, restart it by start-node.sh and then stop-node.sh

## Open the url: "NodeIP:26660"

You should get the metrics output.

## Set up the prometheus and grafana

### Install docker and docker-compose

### Init chain node monitoring settings

```
./start-monitoring.sh
```

The script will as you to provide the node ip, use the same IP of your node.

If the monitoring is set up correctly you will see the output:

```
echo "Open the http://localhost:3000/ to access the Grafana, default credentials are: admin:admin"
```

You can use "stop-monitoring.sh" and "start-monitoring.sh" to restart the monitoring if needed. By default, the
docker-compose will start the container on the machine start.


