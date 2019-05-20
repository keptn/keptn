#!/bin/bash

GIT_SHA=$1
TYPE=$2
NUMBER=$3
DATE=$4

# MAC
env GOOS=darwin GOARCH=amd64 go get ./...
env GOOS=darwin GOARCH=amd64 go build -o keptn
zip keptn-macOS-${GIT_SHA}-${TYPE}-${NUMBER}-${DATE}.zip keptn
tar -zcvf keptn-macOS-${GIT_SHA}-${TYPE}-${NUMBER}-${DATE}.tar.gz keptn
rm keptn

# Linux
env GOOS=linux GOARCH=amd64 go get ./...
env GOOS=linux GOARCH=amd64 go build -o keptn
zip keptn-linux-${GIT_SHA}-${TYPE}-${NUMBER}-${DATE}.zip keptn
tar -zcvf keptn-linux-${GIT_SHA}-${TYPE}-${NUMBER}-${DATE}.tar.gz keptn
rm keptn

gsutil cp keptn-linux-${GIT_SHA}-${TYPE}-${NUMBER}-${DATE}.zip gs://keptn-cli

# Windows
env GOOS=windows GOARCH=amd64 go get ./...
env GOOS=windows GOARCH=amd64 go build -o keptn.exe
zip keptn-windows-${GIT_SHA}-${TYPE}-${NUMBER}-${DATE}.zip keptn.exe
tar -zcvf keptn-windows-${GIT_SHA}-${TYPE}-${NUMBER}-${DATE}.tar.gz keptn.exe
rm keptn.exe

gsutil cp keptn-windows-${GIT_SHA}-${TYPE}-${NUMBER}-${DATE}.zip gs://keptn-cli

ls -lsa
