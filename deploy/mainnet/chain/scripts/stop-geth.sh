#!/bin/bash
set -eu

echo "Stopping geth"

pkill geth && echo "geth  is stopped"