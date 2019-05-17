#!/usr/bin/env bash

FOLDER=$1
IMAGE=$2
GIT_SHA=$3
DATE=$4
VERSION=$5

cd "${FOLDER}"

docker build . -t "${IMAGE}:${GIT_SHA}"
docker tag "${IMAGE}:${GIT_SHA}" "${IMAGE}:${VERSION}.${DATE}"
docker tag "${IMAGE}:${GIT_SHA}" "${IMAGE}:${VERSION}.latest"
docker push "${IMAGE}:${GIT_SHA}"
docker push "${IMAGE}:${VERSION}.${DATE}"
docker push "${IMAGE}:${VERSION}.latest"

cd ../..
