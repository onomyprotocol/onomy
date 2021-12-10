# Steps to set up and start the monitoring

* Install docker and docker-compose

* Init chain node monitoring settings

```
./start-monitoring.sh
```

The script will ask you to provide the node ip, use the IP of your node.

If the monitoring is set up correctly you will see the output:

```
echo "Open the http://localhost:3000/ to access the Grafana, default credentials are: admin:admin"
```

You can use "stop-monitoring.sh" and "start-monitoring.sh" to restart the monitoring if needed. By default, the
docker-compose will start the container on the machine start.
