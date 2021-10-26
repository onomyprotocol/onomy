#!/bin/bash
set -eux
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

REPOFOLDER=$DIR/../..
pushd $REPOFOLDER

docker build --pull -f tests/tm-load-test/dockerfile/Dockerfile  -t cosmoschain .
