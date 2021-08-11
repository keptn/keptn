#!/bin/bash

# target branch/tag/...
GO_UTILS_TARGET=${1:-master}

echo "Go-Utils Target Commit/Branch/Tag: ${GO_UTILS_TARGET}"

# update go modules in all directories that contain a go.mod which contains go-utils
for file in ./*; do
  if [[ -f "$file/go.mod" ]]; then
    echo "Checking if $file/go.mod contains go-utils"
    if grep "github.com/keptn/go-utils" "$file/go.mod"; then
      echo "Yes, updating go-utils now..."
      cd "$file" || exit
      # fetch the desired version (this will update go.mod and go.sum)
      go get "github.com/keptn/go-utils@$GO_UTILS_TARGET" && \
      go get ./... && \
      go mod tidy
      cd - || exit
    fi
  fi
done

echo "Changed files:"
git status
