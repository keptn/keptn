#!/usr/bin/env bash

IMAGE=$1;
GIT_SHA=$2;
TYPE=$3;
NUMBER=$4;
DATE=$5;

docker build . -t "${IMAGE}:${GIT_SHA}"
docker tag "${IMAGE}:${GIT_SHA}" "${IMAGE}:${TYPE}.${NUMBER}.${DATE}"
docker push "${IMAGE}:${GIT_SHA}"
docker push "${IMAGE}:${TYPE}.${NUMBER}.${DATE}"
