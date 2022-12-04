#!/bin/bash
set -eu

echo "Stopping onomy node"

pkill onomyd && echo "onomyd is stopped"
