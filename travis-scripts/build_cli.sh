#!/bin/bash

VERSION=$1
echo "$VERSION" > version

# MAC is not supported yet
env GOOS=darwin GOARCH=amd64 go get ./...
env GOOS=darwin GOARCH=amd64 go build -ldflags="-X 'main.Version=$VERSION'" -o keptn
zip keptn-macOS.zip keptn
tar -zcvf keptn-macOS.tar.gz keptn
rm keptn

gsutil cp keptn-macOS.zip gs://keptn-cli/${TAG}/keptn-macOS.zip
gsutil cp keptn-macOS.tar.gz gs://keptn-cli/${TAG}/keptn-macOS.tar.gz

# Linux
env GOOS=linux GOARCH=amd64 go get ./...
env GOOS=linux GOARCH=amd64 go build -ldflags="-X 'main.Version=$VERSION'" -o keptn
zip keptn-linux.zip keptn
tar -zcvf keptn-linux.tar.gz keptn
rm keptn

gsutil cp keptn-linux.zip gs://keptn-cli/${TAG}/keptn-linux.zip
gsutil cp keptn-linux.tar.gz gs://keptn-cli/${TAG}/keptn-linux.tar.gz

# Windows
env GOOS=windows GOARCH=amd64 go get ./...
env GOOS=windows GOARCH=amd64 go build -ldflags="-X 'main.Version=$VERSION'" -o keptn.exe
zip keptn-windows.zip keptn.exe
tar -zcvf keptn-windows.tar.gz keptn.exe
rm keptn.exe

gsutil cp keptn-windows.zip gs://keptn-cli/${TAG}/keptn-windows.zip
gsutil cp keptn-windows.tar.gz gs://keptn-cli/${TAG}/keptn-windows.tar.gz
