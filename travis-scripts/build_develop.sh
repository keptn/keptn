#!/bin/bash

# ${IMAGE}=$1 ${FOLDER}=$2 ${GIT_SHA}=$3 ${DATE}=$4
IMAGE=$1
FOLDER=$2
GIT_SHA=$3
DATE=$4

echo "Build ${IMAGE}"
cp MANIFEST ./${FOLDER}MANIFEST #$FOLDER contains / at the end
cp travis-scripts/entrypoint.sh ./${FOLDER}entrypoint.sh #$FOLDER contains / at the end
cd ./${FOLDER}
# uncomment certain lines from Dockerfile that are for travis builds only
sed -i '/#travis-uncomment/s/^#travis-uncomment //g' Dockerfile

cat MANIFEST
docker build . -t "${IMAGE}:${GIT_SHA}" -t "${IMAGE}:${DATE}" -t "${IMAGE}:latest" --build-arg version=$VERSION || travis_terminate 1
docker push "${IMAGE}"
