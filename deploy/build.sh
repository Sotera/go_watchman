#!/usr/bin/env bash

set -o errexit

if [ "$#" -lt 3 ]; then
  echo "usage: $0 {project} {tag} {supervisord | standalone} [push]"
  exit 1
fi

proj=$1
tag=$2
mode=$3
docker_proj=go_watchman_$proj:$tag

set -x

go test ../$proj
# executables expected in 'cmd' dir
go build -o bin/$proj github.com/Sotera/go_watchman/$proj/cmd
docker build -f Dockerfile-$mode -t sotera/$docker_proj --build-arg PROJ=$proj .
if [ "$4" == "push" ]; then
  docker push sotera/$docker_proj
fi