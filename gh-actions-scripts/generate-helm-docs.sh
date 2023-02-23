#!/bin/bash
# Readme generator for Keptn Helm Chart
#
# This script will install the readme generator if it's not installed already
# and then it will generate the readme from the local helm values
#
# Dependencies:
# Node >=16

echo "Checking if readme generator is installed already..."
if [[ $(npm list -g | grep -c 'readme-generator-for-helm') -eq 0 ]]; then
  echo "Readme Generator not installed, installing now..."
  git clone https://github.com/bitnami-labs/readme-generator-for-helm.git
  cd ./readme-generator-for-helm || exit
  npm ci
  cd ..
  npm install -g ./readme-generator-for-helm
else
  echo "Readme Generator is already installed, continuing..."
fi

echo "Generating readme now..."
readme-generator --values=./installer/manifests/keptn/values.yaml --readme=./installer/manifests/keptn/README.md

# Please be aware, the readme file needs to exist and needs to have a Parameters section, as only this section will be re-generated
