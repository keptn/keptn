#!/bin/bash

shopt -s globstar

# target branch/tag/...
TARGET=${1:-master}
ARTIFACT=${2:-go-utils}

echo "$ARTIFACT Target Commit/Branch/Tag: ${TARGET}"

# update go modules in all directories that contain a go.mod which contains utils
for file in ./**/*; do
  if [[ -f "$file/go.mod" ]]; then
    echo "Checking if $file/go.mod contains $ARTIFACT"
    if grep "github.com/keptn/$ARTIFACT" "$file/go.mod"; then
      echo "Yes, updating $ARTIFACT now..."
      cd "$file" || exit
      # fetch the desired version (this will update go.mod and go.sum)
      go get "github.com/keptn/$ARTIFACT@$TARGET" && \
      go get ./... && \
      go mod tidy
      cd - || exit
    fi
  fi
done

echo "Changed files:"
git status
shopt -u globstar
