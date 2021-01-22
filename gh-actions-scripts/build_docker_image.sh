#!/bin/bash

if [[ "$#" -ne 5 ]]; then
  echo "Usage: $0 IMAGE FOLDER GIT_SHA VERSION DATETIME"
  echo "      Example: $0 keptn/api api/ 1234abcd 1.2.3 20210101101210"
  exit 1
fi

# ${IMAGE}=$1 ${FOLDER}=$2 ${GIT_SHA}=$3 ${VERSION}=$4 ${DATETIME}=$5
IMAGE=$1
FOLDER=$2
GIT_SHA=$3
VERSION=$4
DATETIME=$5

# ensure trailing slash in folder
if [[ "${FOLDER}" != */ ]]; then
  echo "Please ensure that FOLDER has a trailing slash, e.g., api/"
  exit 1
fi

echo "Building Docker Image ${IMAGE}:${VERSION}.${DATETIME} from ${FOLDER}DOCKERFILE"
cp MANIFEST ./${FOLDER}MANIFEST #$FOLDER contains / at the end

if [[ $? -ne 0 ]]; then
  echo "::error file=${FOLDER}/Dockerfile::Could not find MANIFEST. Please create a file called MANIFEST"
  exit 1
fi

cp travis-scripts/entrypoint.sh ./${FOLDER}entrypoint.sh #$FOLDER contains / at the end

cd ./${FOLDER}

# uncomment certain lines from Dockerfile that are for Travis builds only
sed -i '/#travis-uncomment/s/^#travis-uncomment //g' Dockerfile
cat MANIFEST
docker build . -t "${IMAGE}:${VERSION}.${DATETIME}" -t "${IMAGE}:${VERSION}" --build-arg version="${VERSION}"

if [[ $? -ne 0 ]]; then
  echo "Failed to build Docker Image ${IMAGE}:${VERSION}.${DATETIME}, exiting"
  echo "::error file=${FOLDER}/Dockerfile::Failed to build Docker Image"
  exit 1
fi

# push all tags that we just built
docker push "${IMAGE}:${VERSION}.${DATETIME}" 
docker push "${IMAGE}:${VERSION}"

if [[ $? -ne 0 ]]; then
  echo "::warning file=${FOLDER}/Dockerfile::Failed to push ${IMAGE}:${VERSION}.${DATETIME} to DockerHub, continuing anyway"
  echo "* Failed to push ${IMAGE}:${VERSION}.${DATETIME} and ${IMAGE}:${VERSION} (Source: ${FOLDER})" >> ../docker_build_report/report.txt
else
  echo "* Pushed ${IMAGE}:${VERSION}.${DATETIME} and ${IMAGE}:${VERSION} (Source: ${FOLDER})" >> ../docker_build_report/report.txt
fi

# change back to previous directory
cd -
