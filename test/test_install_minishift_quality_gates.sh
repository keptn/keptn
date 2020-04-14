#!/bin/bash

source test/utils.sh

echo "{
  \"openshiftUrl\": \"https://127.0.0.1:8443\",
  \"openshiftUser\": \"admin\",
  \"openshiftPassword\": \"admin\"
}" > creds.json

echo "Installing keptn on minishift cluster"

# Install keptn (using the develop version, which should point the :latest docker images)
keptn install --platform=openshift --keptn-installer-image=keptn/installer:latest --use-case=quality-gates --creds=creds.json --verbose
