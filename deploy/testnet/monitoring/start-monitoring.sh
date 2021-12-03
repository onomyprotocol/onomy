#!/bin/bash
set -eu

if [[ -z "${ONOMY-NODE-IP}" ]]; then
  echo "The variable ONOMY-NODE-IP is required: "
  exit
fi

docker-compose up -d

echo "Open the http://localhost:3000/ to access the Grafana, default credentials are: admin:admin:"
