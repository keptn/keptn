#!/bin/bash

set -e

DOCKER_ORG_NAME="keptn"
DOCKER_VERSION=${DOCKER_VERSION:-$(git describe --abbrev=1 --tags || echo "dev")}
DOCKER_ROOT_PATH="./"

cd "$DOCKER_ROOT_PATH"

# List of docker images
mapfile -t DOCKERFILE_LIST < <(find . -name 'Dockerfile')

# List of excluded docker images
DOCKERFILE_EXCLUSIONS=(
  "upgrader"
)

# List of docker name overrides
declare -A OVERRIDES
OVERRIDES["bridge"]="bridge2"

# String tokenizer using the '/' delimiter
tokenizeString() {
  IFS='/'
  read -ra DOCKERFILE_PATH <<<"$1"
  IFS=' '
}

for dockerfile in "${DOCKERFILE_LIST[@]}"; do
  tokenizeString "$dockerfile"

  TOKENS_SIZE=${#DOCKERFILE_PATH[@]}
  DOCKER_PATH=${dockerfile%/*}
  DOCKER_DIR=${DOCKERFILE_PATH[$TOKENS_SIZE-2]}
  DOCKER_PARENT_DIR=${DOCKERFILE_PATH[$TOKENS_SIZE-3]}

  # Apply Dockerfile exclusions
  if printf '%s\n' "${DOCKERFILE_EXCLUSIONS[@]}" | grep -q -P "^$DOCKER_PARENT_DIR$"; then
    echo "Excluded $dockerfile"
    continue
  fi

  # Apply Dockerfile name overrides
  if [[ -v OVERRIDES[${DOCKER_DIR}] ]]; then
    DOCKER_NAME="${OVERRIDES[${DOCKER_DIR}]}"
    echo "Overridden docker tag $DOCKER_DIR -> $DOCKER_NAME"
  else
    DOCKER_NAME=$DOCKER_DIR
  fi

  # Build docker image
  DOCKER_TAG="$DOCKER_ORG_NAME/$DOCKER_NAME:$DOCKER_VERSION"

  echo "Building docker image $DOCKER_TAG using $dockerfile"
  docker build -t "$DOCKER_TAG" "$DOCKER_PATH"
done
