#!/bin/bash

DATE="$(date +'%Y%m%d.%H%M')"
GIT_SHA="$(git rev-parse --short HEAD)"
REPO_URL=https://github.com/$GITHUB_REPOSITORY
JOB_URL=${REPO_URL}/actions/runs/${GITHUB_RUN_ID}

sed -i 's~MANIFEST_REPOSITORY~'"$REPO_URL"'~' MANIFEST
sed -i 's~MANIFEST_BRANCH~'"$BRANCH"'~' MANIFEST
sed -i 's~MANIFEST_COMMIT~'"$GIT_SHA"'~' MANIFEST
sed -i 's~MANIFEST_TRAVIS_JOB_URL~'"$JOB_URL"'~' MANIFEST
sed -i 's~MANIFEST_DATE~'"$DATE"'~' MANIFEST
