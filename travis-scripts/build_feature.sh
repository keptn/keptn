#!/usr/bin/env bash

FOLDER=$1
IMAGE=$2
GIT_SHA=$3
TYPE=$4
NUMBER=$5
DATE=$6

cd "${FOLDER}"

docker build . -t "${IMAGE}:${GIT_SHA}"
docker tag "${IMAGE}:${GIT_SHA}" "${IMAGE}:${TYPE}.${NUMBER}.${DATE}"
docker push "${IMAGE}:${GIT_SHA}"
docker push "${IMAGE}:${TYPE}.${NUMBER}.${DATE}"

cd ../..
