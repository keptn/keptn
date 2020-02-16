#!/bin/bash

VERSION=$(cat version | tr -d '[:space:]')
KUBE_CONSTRAINTS=${KUBE_CONSTRAINTS:->=1.13,<=1.15}
go build -ldflags="-X 'main.Version=$VERSION' -X 'main.KubeServerVersionConstraints=$KUBE_CONSTRAINTS'"  -o keptn

# MAC
env GOOS=darwin GOARCH=amd64 go build -ldflags="-X 'main.Version=$VERSION' -X 'main.KubeServerVersionConstraints=$KUBE_CONSTRAINTS'" -o keptn
zip keptn-macOS.zip keptn
tar -zcvf keptn-macOS.tar.gz keptn
rm keptn

# Linux build covered by Travis CI

#env GOOS=linux GOARCH=amd64 go build -ldflags="-X 'main.Version=$VERSION' -X 'main.KubeServerVersionConstraints=$KUBE_CONSTRAINTS'" -o keptn
#zip keptn-linux.zip keptn
#tar -zcvf keptn-linux.tar.gz keptn
#rm keptn

# Windows build covered by Travis CI

#env GOOS=windows GOARCH=amd64 go build -ldflags="-X 'main.Version=$VERSION' -X 'main.KubeServerVersionConstraints=$KUBE_CONSTRAINTS'" -o keptn.exe
#zip keptn-windows.zip keptn.exe
#tar -zcvf keptn-windows.tar.gz keptn.exe
#rm keptn.exe
