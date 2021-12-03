#!/bin/bash
set -eu

if [[ -z "${ONOMY_NODE_IP}" ]]; then
  echo "The variable ONOMY-NODE-IP is required: "
  exit
fi

echo "Starting monitoring for: $ONOMY_NODE_IP"

docker-compose up -d

echo "Open the http://localhost:3000/ to access the Grafana, default credentials are: admin:admin"
