#!/bin/bash

VERSION=${VERSION:-$(git describe --abbrev=1 --tags || echo "dev")}
KUBE_CONSTRAINTS=${KUBE_CONSTRAINT:-""}
OUTPUT_EXECUTEABLE_NAME=${OUTPUT_EXECUTEABLE_NAME:-"keptn"}

cd ./cli/ || return


echo "Building Keptn CLI"
env go mod download
env go build -v -x -ldflags="-X 'main.Version=$VERSION' -X 'main.KubeServerVersionConstraints=$KUBE_CONSTRAINTS'" -o "${OUTPUT_EXECUTEABLE_NAME}"

if [ $? -ne 0 ]; then
    echo "Error compiling Keptn CLI, exiting ..."
    exit 1
fi
