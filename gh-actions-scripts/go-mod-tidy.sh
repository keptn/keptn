#!/bin/bash

# update go modules in all directories that contain a go.mod
for file in ./{**/,}*; do
  if [[ -f "$file/go.mod" ]]; then
    cd "$file" || exit
    go mod tidy
    cd - || exit
  fi
done
