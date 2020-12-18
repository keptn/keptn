#!/bin/bash

echo "Download the CLI"

# Download latest KEPTN cli for linux
wget "https://storage.googleapis.com/keptn-cli/${KEPTN_CLI_VERSION}/keptn-linux.zip"
unzip keptn-linux.zip

sudo mv keptn /usr/local/bin/keptn
