#!/bin/bash

# ${IMAGE}=$1 ${FOLDER}=$2 ${GIT_SHA}=$3 ${TYPE}=$4 ${NUMBER}=$5 ${DATE}=$6
IMAGE=$1
FOLDER=$2
GIT_SHA=$3
TYPE=$4
NUMBER=$5
DATE=$6

echo "Build ${IMAGE}"
cp MANIFEST ./${FOLDER}MANIFEST #$FOLDER contains / at the end
cd ./${FOLDER}
cat MANIFEST
docker build . -t "${IMAGE}:${GIT_SHA}"
docker tag "${IMAGE}:${GIT_SHA}" "${IMAGE}:${TYPE}.${NUMBER}.${DATE}"
docker push "${IMAGE}:${GIT_SHA}"
docker push "${IMAGE}:${TYPE}.${NUMBER}.${DATE}"