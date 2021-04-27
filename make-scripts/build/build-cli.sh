#!/bin/bash

VERSION=${VERSION:-$(git describe --abbrev=1 --tags || echo "dev")}
KUBE_CONSTRAINTS=${KUBE_CONSTRAINT:-">= 1.14, <= 1.20"}
OUTPUT_EXECUTABLE_NAME=${OUTPUT_EXECUTABLE_NAME:-"keptn"}

cd ./cli/ || return


echo "Building Keptn CLI"
env go mod download
env go build -v -x -ldflags="-X 'main.Version=$VERSION' -X 'main.KubeServerVersionConstraints=$KUBE_CONSTRAINTS'" -o "${OUTPUT_EXECUTABLE_NAME}"

# shellcheck disable=SC2181
if [ $? -ne 0 ]; then
    echo "Error compiling Keptn CLI, exiting ..."
    exit 1
fi
