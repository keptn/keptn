#!/bin/bash

# ${IMAGE}=$1 ${FOLDER}=$2 ${GIT_SHA}=$3 ${TYPE}=$4 ${NUMBER}=$5 ${DATE}=$6
IMAGE=$1
FOLDER=$2
GIT_SHA=$3
TYPE=$4
NUMBER=$5
DATE=$6

echo "Build keptn ${IMAGE}"
cp MANIFEST ./${FOLDER}
cd ./${FOLDER}
cat MANIFEST
docker build . -t "${IMAGE}:${GIT_SHA}"
docker tag "${IMAGE}:${GIT_SHA}" "${API_IMAGE}:${TYPE}.${NUMBER}.${DATE}"
docker push "${API_IMAGE}:${GIT_SHA}"
docker push "${API_IMAGE}:${TYPE}.${NUMBER}.${DATE}"