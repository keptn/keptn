#!/bin/bash

GIT_SHA=$1
DATE=$2

# MAC
env GOOS=darwin GOARCH=amd64 go get ./...
env GOOS=darwin GOARCH=amd64 go build -o keptn
zip keptn-macOS-${DATE}-latest.zip keptn
tar -zcvf keptn-macOS-${DATE}-latest.tar.gz keptn
rm keptn

# Linux
env GOOS=linux GOARCH=amd64 go get ./...
env GOOS=linux GOARCH=amd64 go build -o keptn
zip keptn-linux-${DATE}-latest.zip keptn
tar -zcvf keptn-linux-${DATE}-latest.tar.gz keptn
rm keptn

gsutil cp keptn-linux-${DATE}-latest.zip gs://keptn-cli

# Windows
env GOOS=windows GOARCH=amd64 go get ./...
env GOOS=windows GOARCH=amd64 go build -o keptn.exe
zip keptn-windows-${DATE}-latest.zip keptn.exe
tar -zcvf keptn-windows-${DATE}-latest.tar.gz keptn.exe
rm keptn.exe

gsutil cp keptn-windows-${DATE}-latest.zip gs://keptn-cli

ls -lsa
