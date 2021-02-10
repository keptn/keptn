#!/bin/bash

# run cli tests
cd ./cli || return
go test -race -v  -coverprofile=coverage.txt -covermode=atomic ./...
