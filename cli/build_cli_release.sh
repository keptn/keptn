#!/bin/bash

GIT_SHA=$1
DATE=$2
VERSION=$3

# MAC
env GOOS=darwin GOARCH=amd64 go get ./...
env GOOS=darwin GOARCH=amd64 go build -o keptn
zip keptn-macOS-${VERSIN}-${DATE}.zip keptn
tar -zcvf keptn-macOS-${VERSIN}-${DATE}.tar.gz keptn
rm keptn

# Linux
env GOOS=linux GOARCH=amd64 go get ./...
env GOOS=linux GOARCH=amd64 go build -o keptn
zip keptn-linux-${VERSIN}-${DATE}.zip keptn
tar -zcvf keptn-linux-${VERSIN}-${DATE}.tar.gz keptn
rm keptn

gsutil cp keptn-linux-${VERSIN}-${DATE}.zip gs://keptn-cli

# Windows
env GOOS=windows GOARCH=amd64 go get ./...
env GOOS=windows GOARCH=amd64 go build -o keptn.exe
zip keptn-windows-${VERSIN}-${DATE}.zip keptn.exe
tar -zcvf keptn-windows-${VERSIN}-${DATE}.tar.gz keptn.exe
rm keptn.exe

gsutil cp keptn-windows-${VERSIN}-${DATE}.zip gs://keptn-cli

ls -lsa
