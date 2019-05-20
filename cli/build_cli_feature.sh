#!/bin/bash

GIT_SHA=$1
TYPE=$2
NUMBER=$3
DATE=$4

# MAC
env GOOS=darwin GOARCH=amd64 go get ./...
env GOOS=darwin GOARCH=amd64 go build -o keptn
zip keptn-macOS-${TYPE}-${NUMBER}-${DATE}.zip keptn
tar -zcvf keptn-macOS-${TYPE}-${NUMBER}-${DATE}.tar.gz keptn
rm keptn

# Linux
env GOOS=linux GOARCH=amd64 go get ./...
env GOOS=linux GOARCH=amd64 go build -o keptn
zip keptn-linux-${TYPE}-${NUMBER}-${DATE}.zip keptn
tar -zcvf keptn-linux-${TYPE}-${NUMBER}-${DATE}.tar.gz keptn
rm keptn

gsutil cp keptn-linux-${TYPE}-${NUMBER}-${DATE}.zip gs://keptn-cli

# Windows
env GOOS=windows GOARCH=amd64 go get ./...
env GOOS=windows GOARCH=amd64 go build -o keptn.exe
zip keptn-windows-${TYPE}-${NUMBER}-${DATE}.zip keptn.exe
tar -zcvf keptn-windows-${TYPE}-${NUMBER}-${DATE}.tar.gz keptn.exe
rm keptn.exe

gsutil cp keptn-windows-${TYPE}-${NUMBER}-${DATE}.zip gs://keptn-cli

ls -lsa
