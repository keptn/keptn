#!/bin/bash

# target branch/tag/...
LIBRARY_NAME=${1:-"github.com/keptn/go-utils"}
LIBRARY_TARGET=${2:-"master"}

echo "Target Branch/Tag: ${LIBRARY_NAME}@${LIBRARY_TARGET}"

# update go modules in all directories that contain a go.mod which contains go-utils
for file in *; do
  if [[ -f "$file/go.mod" ]]; then
    echo "Checking if $file/go.mod contains ${LIBRARY_NAME}"
    grep ${LIBRARY_NAME} "$file/go.mod"
    if [[ $? -eq 0 ]]; then
      echo "Yes, updating go-utils now..."
      cd $file || exit
      # fetch the desired version (this will update go.mod and go.sum)
      go get "${LIBRARY_NAME}@${LIBRARY_TARGET}"
      cd - || exit
    fi
  fi
done
# debug:
git status
