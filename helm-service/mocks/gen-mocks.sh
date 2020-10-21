#!/bin/bash

set -euo pipefail

GO111MODULE=off go get -u github.com/golang/mock/mockgen

mockgen -package mocks -destination=./mock_chart_storer.go github.com/keptn/kubernetes-utils/pkg ChartStorer
mockgen -package mocks -destination=./mock_chart_packager.go github.com/keptn/kubernetes-utils/pkg ChartPackager