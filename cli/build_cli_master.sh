#!/bin/bash

VERSION=$1

# MAC
env GOOS=darwin GOARCH=amd64 go get ./...
env GOOS=darwin GOARCH=amd64 go build -o keptn
zip keptn-macOS-${VERSION}.zip keptn
tar -zcvf keptn-macOS-${VERSION}.tar.gz keptn
rm keptn

# Linux
env GOOS=linux GOARCH=amd64 go get ./...
env GOOS=linux GOARCH=amd64 go build -o keptn
zip keptn-linux-${VERSION}.zip keptn
tar -zcvf keptn-linux-${VERSION}.tar.gz keptn
rm keptn

gsutil cp keptn-linux-${VERSION}.zip gs://keptn-cli

# Windows
env GOOS=windows GOARCH=amd64 go get ./...
env GOOS=windows GOARCH=amd64 go build -o keptn.exe
zip keptn-windows-${VERSION}.zip keptn.exe
tar -zcvf keptn-windows-${VERSION}.tar.gz keptn.exe
rm keptn.exe

gsutil cp keptn-windows-${VERSION}.zip gs://keptn-cli

ls -lsa
