#!/bin/bash

## should be started from a Mac OS X 

# MAC
go build -o keptn
zip keptn-macOS.zip keptn
tar -zcvf keptn-macOS.tar.gz keptn
rm keptn

# Linux
env GOOS=linux GOARCH=amd64 go build -o keptn
zip keptn-linux.zip keptn
tar -zcvf keptn-linux.tar.gz keptn
rm keptn

# Windows
env GOOS=windows GOARCH=amd64 go build -o keptn.exe
zip keptn-windows.zip keptn.exe
tar -zcvf keptn-windows.tar.gz keptn.exe
rm keptn.exe