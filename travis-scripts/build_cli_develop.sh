#!/bin/bash

GIT_SHA=$1
DATE=$2

FOLDER="${DATE}-latest"

# MAC
env GOOS=darwin GOARCH=amd64 go get ./...
env GOOS=darwin GOARCH=amd64 go build -o keptn
zip keptn-macOS-latest.zip keptn
tar -zcvf keptn-macOS-latest.tar.gz keptn
rm keptn

# Linux
env GOOS=linux GOARCH=amd64 go get ./...
env GOOS=linux GOARCH=amd64 go build -o keptn
zip keptn-linux-latest.zip keptn
tar -zcvf keptn-linux-latest.tar.gz keptn
rm keptn

gsutil cp keptn-linux.zip gs://keptn-cli/${FOLDER}/

# Windows
env GOOS=windows GOARCH=amd64 go get ./...
env GOOS=windows GOARCH=amd64 go build -o keptn.exe
zip keptn-windows-latest.zip keptn.exe
tar -zcvf keptn-windows-latest.tar.gz keptn.exe
rm keptn.exe

gsutil cp keptn-windows.zip gs://keptn-cli/${FOLDER}/

ls -lsa
