#!/usr/bin/env bash

set -o errexit

if [ "$#" -lt 3 ]; then
  echo "usage: $0 project tag [push]"
  exit 1
fi

proj=$1
tag=$2
docker_proj=go_watchman_$proj:$tag

set -x

go install github.com/Sotera/go_watchman/$proj
cp $GOPATH/bin/$proj ./bin
docker build -t sotera/$docker_proj --build-arg PROJ=$proj .
if [ "$3" == "push" ]; then
  docker push sotera/$docker_proj
fi