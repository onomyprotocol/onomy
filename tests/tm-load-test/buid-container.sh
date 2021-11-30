#!/bin/bash
set -eux
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

REPOFOLDER=$DIR/../..
pushd $REPOFOLDER

docker build -f tests/tm-load-test/Dockerfile  -t cosmoschain .
