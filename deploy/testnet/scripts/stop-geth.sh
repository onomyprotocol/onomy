#!/bin/bash
set -eu

echo "Stopping full node"

kill $(pidof geth) && echo "geth  is stopped"