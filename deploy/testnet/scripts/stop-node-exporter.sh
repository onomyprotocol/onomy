#!/bin/bash
set -eu

echo "Stopping node_exporter"

kill $(pidof node_exporter) && echo "node_exporter is stopped"