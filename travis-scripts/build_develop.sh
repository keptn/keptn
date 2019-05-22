#!/usr/bin/env bash

IMAGE=$1
GIT_SHA=$2
DATE=$3

docker build . -t "${IMAGE}:${GIT_SHA}"
docker tag "${IMAGE}:${GIT_SHA}" "${IMAGE}:${DATE}"
docker tag "${IMAGE}:${GIT_SHA}" "${IMAGE}:latest"
docker push "${IMAGE}:${GIT_SHA}"
docker push "${IMAGE}:${DATE}"
docker push "${IMAGE}:latest"
