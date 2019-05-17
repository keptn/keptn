#!/usr/bin/env bash

IMAGE=$1
GIT_SHA=$2
DATE=$3
VERSION=$4

docker build . -t "${IMAGE}:${GIT_SHA}"
docker tag "${IMAGE}:${GIT_SHA}" "${IMAGE}:${VERSION}.${DATE}"
docker tag "${IMAGE}:${GIT_SHA}" "${IMAGE}:${VERSION}.latest"
docker push "${IMAGE}:${GIT_SHA}"
docker push "${IMAGE}:${VERSION}.${DATE}"
docker push "${IMAGE}:${VERSION}.latest"
