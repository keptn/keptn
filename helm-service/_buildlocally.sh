#!/bin/bash

## requires go 1.12+

if [ ! -z "$debugBuild" ]; then export BUILDFLAGS='-gcflags "all=-N -l"'; fi

# CGO_ENABLED=0 GOOS=linux go test $BUILDFLAGS -v -o helm-service
CGO_ENABLED=0 GOOS=linux go build $BUILDFLAGS -v -o helm-service