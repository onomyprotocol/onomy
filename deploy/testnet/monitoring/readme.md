# Steps set up and start the monitoring

## Set up the node

If you didn't run the expose-metrics.sh before run it:

```
./expose-metrics.sh
```

This script will enable the prometheus metrics in your node config.

## If the node is running now, restart it by start-node.sh and then stop-node.sh

## Open the url: "NodeIP:26660"

You should get the metrics output.

## Set up the prometheus and grafana

### install docker and docker-compose

### Init monitoring settings

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


