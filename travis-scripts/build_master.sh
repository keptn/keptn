#!/usr/bin/env bash

FOLDER=$1
IMAGE=$2
VERSION=$3

cd "${FOLDER}"

docker build . -t "${IMAGE}:${VERSION}"
docker push "${IMAGE}:${VERSION}"
