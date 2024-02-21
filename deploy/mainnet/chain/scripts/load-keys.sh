#!/bin/bash
set -eu

echo "Loading orchestrator key"
pass keyring-onomy/eth-orchestrator-eth-private-key > /dev/null

echo "Required keys are loaded"