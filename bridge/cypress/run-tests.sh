#!/bin/bash

yarn test:ui
test_outcome=$?
mv ./dist/cypress/screenshots /shared/screenshots

if [[ $test_outcome -ne 0 ]]; then
  exit 1
fi
