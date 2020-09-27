#!/bin/bash

VERSION=${1:-develop}
KUBE_CONSTRAINTS=$2

cd ./cli/

if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    ############################################
    # Linux
    ############################################
    echo "Building Keptn CLI for Linux"
    env GOOS=linux GOARCH=amd64 go mod download
    env GOOS=linux GOARCH=amd64 go build -v -x -ldflags="-X 'main.Version=$VERSION' -X 'main.KubeServerVersionConstraints=$KUBE_CONSTRAINTS'" -o keptn

    if [ $? -ne 0 ]; then
    echo "Error compiling Keptn CLI, exiting ..."
    exit 1
    fi

elif [[ "$OSTYPE" == "darwin"* ]]; then
    ############################################
    # MAC OS
    ############################################
    echo "Building Keptn CLI for OSX"
    env GOOS=darwin GOARCH=amd64 go mod download
    env GOOS=darwin GOARCH=amd64 go build -v -x  -ldflags="-X 'main.Version=$VERSION' -X 'main.KubeServerVersionConstraints=$KUBE_CONSTRAINTS'" -o keptn

    if [ $? -ne 0 ]; then
    echo "Error compiling Keptn CLI, exiting ..."
    exit 1
    fi
fi