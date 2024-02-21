#!/bin/bash
set -eu

echo "Stopping node exporter"

pkill node_exporter && echo "node exporter is stopped"