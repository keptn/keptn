#!/usr/bin/env bash

IMAGE=$1;
VERSION=$2;

docker build . -t "${IMAGE}:${VERSION}"
docker push "${IMAGE}:${VERSION}"
