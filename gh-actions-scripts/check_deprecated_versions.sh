#!/bin/bash

set -e
ls -la $GITHUB_WORKSPACE/bin
./helm template installer/manifests/keptn | ./pluto detect-files -owide
./pluto detect-files -owide
set +e
