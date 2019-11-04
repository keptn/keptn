#!/bin/bash

# ${IMAGE}=$1 ${FOLDER}=$2
IMAGE=$1
FOLDER=$2

echo "Build ${IMAGE}"
cd ./${FOLDER}
go test ./...