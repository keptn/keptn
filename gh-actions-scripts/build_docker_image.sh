#!/bin/bash

# ${IMAGE}=$1 ${FOLDER}=$2 ${GIT_SHA}=$3 ${VERSION}=$4 ${DATETIME}=$5
IMAGE=$1
FOLDER=$2
GIT_SHA=$3
VERSION=$4
DATETIME=$5

echo "Building Docker Image ${IMAGE}:${VERSION}.${DATETIME}"
cp MANIFEST ./${FOLDER}MANIFEST #$FOLDER contains / at the end
cp travis-scripts/entrypoint.sh ./${FOLDER}entrypoint.sh #$FOLDER contains / at the end

cd ./${FOLDER}

# uncomment certain lines from Dockerfile that are for Travis builds only
sed -i '/#travis-uncomment/s/^#travis-uncomment //g' Dockerfile
cat MANIFEST
docker build . -t "${IMAGE}:${GIT_SHA}" -t "${IMAGE}:${VERSION}.${DATETIME}" -t "${IMAGE}:${VERSION}" --build-arg version="${VERSION}"

if [[ $? -ne 0 ]]; then
  echo "Failed to build Docker Image ${IMAGE}:${VERSION}.${DATETIME}, exiting"
  echo "::error file=${FOLDER}/Dockerfile::Failed to build Docker Image"
  exit 1
fi

docker push "${IMAGE}"

# change back to previous directory
cd -
