#!/bin/bash

## requires go 1.12+

# shellcheck disable=SC2154
if [ -n "$debugBuild" ]; then export BUILDFLAGS='-gcflags "all=-N -l"'; fi
go build -ldflags '-linkmode=external' -o jmeter-extended-service
