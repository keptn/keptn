#! /bin/bash

# Install readme generator 
#
# git clone https://github.com/bitnami-labs/readme-generator-for-helm (zip file included in this PR)
# npm install ./readme-generator-for-helm
#
# dependencies:
# npm 8.11.0
# nodejs 16.16.0

readme-generator --values=installer/manifests/keptn/values.yaml --readme=./bitnami.md

# Please be aware, the readme file needs to exists and needs to have a Parameters section, as only this section will be re-generated 