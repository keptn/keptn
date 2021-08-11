#!/bin/bash

# target branch/tag/...
KUBERNETES_UTILS_TARGET=${1:-master}

echo "Kubernetes-Utils Target Commit/Branch/Tag: ${KUBERNETES_UTILS_TARGET}"

# update go modules in all directories that contain a go.mod which contains go-utils
for file in ./*; do
  if [[ -f "$file/go.mod" ]]; then
    echo "Checking if $file/go.mod contains go-utils"
    if grep "github.com/keptn/kubernetes-utils" "$file/go.mod"; then
      echo "Yes, updating kubernetes-utils now..."
      cd "$file" || exit
      # fetch the desired version (this will update go.mod and go.sum)
      go get "github.com/keptn/kubernetes-utils@$KUBERNETES_UTILS_TARGET" && \
      go get ./... && \
      go mod tidy
      cd - || exit
    fi
  fi
done

echo "Changed files:"
git status
