#!/usr/bin/env bash

set -o errexit

if [ "$#" -lt 4 ]; then
  echo "Builds a go project at specified path, then docker builds it in supervisord or standalone mode, and optionally pushes to docker hub."
  echo "usage: $0 {docker-proj-name} {docker-tag} {path-from-root} {supervisord | standalone} [push]"
  exit 1
fi

name=$1
tag=$2
path=$3
mode=$4
docker_name=go_watchman_$name:$tag

set -x

go test ../$path

go build -o bin/$name github.com/Sotera/go_watchman/$path
docker build -f Dockerfile-$mode -t sotera/$docker_name --build-arg BIN_NAME=$name .
if [ "$5" == "push" ]; then
  docker push sotera/$docker_name
fi