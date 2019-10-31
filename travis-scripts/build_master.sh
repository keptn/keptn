#!/bin/bash

# ${IMAGE}=$1 ${FOLDER}=$2 ${VERSION}=$3
IMAGE=$1
FOLDER=$2
VERSION=$3

echo "Build ${IMAGE}"
cp MANIFEST ./${FOLDER}MANIFEST #$FOLDER contains / at the end
cp travis-scripts/entrypoint.sh ./${FOLDER}entrypoint.sh #$FOLDER contains / at the end
cd ./${FOLDER}
# uncomment certain lines from Dockerfile that are for travis builds only
sed -i '/#travis-uncomment/s/^#travis-uncomment //g' Dockerfile

cat MANIFEST
docker build . -t "${IMAGE}:${VERSION}" --build-arg version=$VERSION
docker push "${IMAGE}:${VERSION}"