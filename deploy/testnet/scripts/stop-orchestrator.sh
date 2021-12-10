#!/bin/bash
set -eu

echo "Stopping gbt orchestrator"

kill $(pidof gbt) && echo "gbt is stopped"
