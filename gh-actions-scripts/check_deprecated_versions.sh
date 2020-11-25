#!/bin/bash

set -e
helm template installer/manifests/keptn | pluto detect-files -owide
pluto detect-files -owide
set +e
