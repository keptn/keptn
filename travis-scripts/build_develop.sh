#!/usr/bin/env bash

FOLDER=$1
IMAGE=$2
GIT_SHA=$3
DATE=$4

cd "${FOLDER}"

docker build . -t "${IMAGE}:${GIT_SHA}"
docker tag "${IMAGE}:${GIT_SHA}" "${IMAGE}:${DATE}"
docker tag "${IMAGE}:${GIT_SHA}" "${IMAGE}:latest"
docker push "${IMAGE}:${GIT_SHA}"
docker push "${IMAGE}:${DATE}"
docker push "${IMAGE}:latest"

cd ../..
