#!/bin/bash
set -eu

echo "Stopping full node"

kill $(pidof onomyd) && echo "onomyd is stopped"
