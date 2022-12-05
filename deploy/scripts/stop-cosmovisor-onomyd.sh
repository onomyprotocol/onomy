#!/bin/bash
set -eu

echo "Stopping onomy node"

pkill cosmovisor && echo "cosmovisor-onomyd is stopped"
