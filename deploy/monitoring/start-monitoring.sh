#!/bin/bash
set -eu

# Switch sed command in the case of linux
fsed() {
  if [ `uname` = 'Linux' ]; then
    sed -i "$@"
  else
    sed -i '' "$@"
  fi
}

if [[ ! -f "prometheus/prometheus.yml" ]]
then
  ONOMY_NODE_IP=
  while [[ $ONOMY_NODE_IP = "" ]]; do
     read -r -p "Enter the ip of the node you want to include for the monitoring:" ONOMY_NODE_IP
  done

  cp prometheus/template.prometheus.yml prometheus/prometheus.yml
  fsed "s#ONOMY_NODE_IP#$ONOMY_NODE_IP#g" prometheus/prometheus.yml
fi

echo "Starting monitoring"

docker-compose up -d

echo "Open the http://localhost:3000/ to access the Grafana, default credentials are: admin:admin"
