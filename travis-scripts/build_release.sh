#!/bin/bash

# ${IMAGE}=$1 ${FOLDER}=$2 ${GIT_SHA}=$3 ${DATE}=$4 ${VERSION}=$5
IMAGE=$1
FOLDER=$2
GIT_SHA=$3
DATE=$4
VERSION=$5

echo "Build ${IMAGE}"
cp MANIFEST ./${FOLDER}MANIFEST #$FOLDER contains / at the end
cp travis-scripts/entrypoint.sh ./${FOLDER}entrypoint.sh #$FOLDER contains / at the end
cd ./${FOLDER}
# uncomment certain lines from Dockerfile that are for travis builds only
sed -i '/#travis-uncomment/s/^#travis-uncomment //g' Dockerfile

cat MANIFEST
docker build . -t "${IMAGE}:${GIT_SHA}" --build-arg version=$VERSION
docker tag "${IMAGE}:${GIT_SHA}" "${IMAGE}:${VERSION}.${DATE}"
docker tag "${IMAGE}:${GIT_SHA}" "${IMAGE}:${VERSION}.latest"
docker push "${IMAGE}:${GIT_SHA}"
docker push "${IMAGE}:${VERSION}.${DATE}"
docker push "${IMAGE}:${VERSION}.latest"