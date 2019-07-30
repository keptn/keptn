#!/bin/bash

# ${IMAGE}=$1 ${FOLDER}=$2 ${GIT_SHA}=$3 ${DATE}=$4
IMAGE=$1
FOLDER=$2
GIT_SHA=$3
DATE=$4

echo "Build ${IMAGE}"
cp MANIFEST ./${FOLDER}MANIFEST #$FOLDER contains / at the end
cd ./${FOLDER}
cat MANIFEST
docker build . -t "${IMAGE}:${GIT_SHA}"
docker tag "${IMAGE}:${GIT_SHA}" "${IMAGE}:${DATE}"
docker tag "${IMAGE}:${GIT_SHA}" "${IMAGE}:latest"
docker push "${IMAGE}:${GIT_SHA}"
docker push "${IMAGE}:${DATE}"
docker push "${IMAGE}:latest"