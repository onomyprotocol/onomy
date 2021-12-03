#!/bin/bash
set -eu

echo "Starting monitoring"

docker-compose up -d

echo "Open the http://localhost:3000/ to access the Grafana, default credentials are: admin:admin"
