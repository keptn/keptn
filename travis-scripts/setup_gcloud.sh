#!/usr/bin/env bash

source ./travis-scripts/setup_functions.sh

# prints the full command before output of the command.
set -x

setup_gcloud

gcloud version
gsutil version