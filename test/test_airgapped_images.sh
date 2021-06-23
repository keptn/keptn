#!/bin/bash

# This integration test checks that all images that are pulled (can be retrieved by kubectl get events -n keptn)
# are coming from the local registry from K3d (passed as a parameter to this script), rather than anywhere else.
# Exception: k3s/rancher based images, as they are needed internally by K3d.

# shellcheck disable=SC1091
source test/utils.sh

EXPECTED_REGISTRY=${1}

KEPTN_NAMESPACE=${KEPTN_NAMESPACE:-keptn}

PULLED_IMAGES=$(kubectl get events -n "${KEPTN_NAMESPACE}" | awk -F'Pulling image ' '{ print $2 }' | xargs)


for IMAGE in ${PULLED_IMAGES}; do

  if [[ "$IMAGE" == "rancher/"* ]]; then
    # ignore rancher / k3s base images
    continue
  fi

  if [[ "$IMAGE" == "$EXPECTED_REGISTRY"* ]]; then
    # okay, as expected
    continue
  fi

  print_error "ERROR: Expected $EXPECTED_REGISTRY in image name, but found '$IMAGE'"
  print_error "Airgapped scenario failed! Please check which images Keptn helm charts are consuming."
  exit 1
done

exit 0
