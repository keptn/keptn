#!/bin/bash

## requires go 1.12+

if [ ! -z "$debugBuild" ]; then export BUILDFLAGS='-gcflags "all=-N -l"'; fi
go build -ldflags '-linkmode=external' -o jmeter-extended-service
