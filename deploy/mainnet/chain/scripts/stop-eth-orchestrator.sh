#!/bin/bash
set -eu

echo "Stopping eth orchestrator"

pkill gbt && echo "eth orchestrator is stopped"
