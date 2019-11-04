#!/bin/bash

# ${IMAGE}=$1 ${FOLDER}=$2 ${VERSION}=$3
IMAGE=$1
FOLDER=$2
VERSION=$3

echo "Build ${IMAGE}"
cp MANIFEST ./${FOLDER}MANIFEST #$FOLDER contains / at the end
cd ./${FOLDER}
cat MANIFEST
docker build . -t "${IMAGE}:${VERSION}" --build-arg version=$VERSION
docker push "${IMAGE}:${VERSION}"