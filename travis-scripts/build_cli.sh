#!/bin/bash

VERSION=${1:-develop}
KUBE_CONSTRAINTS=$2

############################################
# MAC OS
############################################
echo "Building keptn cli for OSX"
env GOOS=darwin GOARCH=amd64 go mod download
env GOOS=darwin GOARCH=amd64 go build -v -x  -ldflags="-X 'main.Version=$VERSION' -X 'main.KubeServerVersionConstraints=$KUBE_CONSTRAINTS'" -o keptn

if [ $? -ne 0 ]; then
  echo 'Error compiling keptn cli, exiting...'
  exit 1
fi

# create archives
zip keptn-macOS.zip keptn
tar -zcvf keptn-macOS.tar.gz keptn
rm keptn

# upload to gcloud
if [ -n "$TAG" ]; then
  gsutil cp keptn-macOS.zip gs://keptn-cli/${TAG}/keptn-macOS.zip
  gsutil cp keptn-macOS.tar.gz gs://keptn-cli/${TAG}/keptn-macOS.tar.gz
fi

rm keptn-macOS.zip
rm keptn-macOS.tar.gz

############################################
# Linux
############################################
echo "Building keptn cli for linux"
env GOOS=linux GOARCH=amd64 go mod download
env GOOS=linux GOARCH=amd64 go build -v -x -ldflags="-X 'main.Version=$VERSION' -X 'main.KubeServerVersionConstraints=$KUBE_CONSTRAINTS'" -o keptn

if [ $? -ne 0 ]; then
  echo 'Error compiling keptn cli, exiting...'
  exit 1
fi

# create archives
zip keptn-linux.zip keptn
tar -zcvf keptn-linux.tar.gz keptn
rm keptn

# upload to gcloud
if [ -n "$TAG" ]; then
  gsutil cp keptn-linux.zip gs://keptn-cli/${TAG}/keptn-linux.zip
  gsutil cp keptn-linux.tar.gz gs://keptn-cli/${TAG}/keptn-linux.tar.gz
fi

rm keptn-linux.zip
rm keptn-linux.tar.gz

############################################
# Windows
############################################
echo "Building keptn cli for windows"
env GOOS=windows GOARCH=amd64 go mod download
env GOOS=windows GOARCH=amd64 go build -ldflags="-X 'main.Version=$VERSION' -X 'main.KubeServerVersionConstraints=$KUBE_CONSTRAINTS'" -o keptn.exe

if [ $? -ne 0 ]; then
  echo 'Error compiling keptn cli, exiting...'
  exit 1
fi

# create archives
zip keptn-windows.zip keptn.exe
tar -zcvf keptn-windows.tar.gz keptn.exe
rm keptn.exe

# upload to gcloud
if [ -n "$TAG" ]; then
  gsutil cp keptn-windows.zip gs://keptn-cli/${TAG}/keptn-windows.zip
  gsutil cp keptn-windows.tar.gz gs://keptn-cli/${TAG}/keptn-windows.tar.gz
fi

rm keptn-windows.zip
rm keptn-windows.tar.gz
