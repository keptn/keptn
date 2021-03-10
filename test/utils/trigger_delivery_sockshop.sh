#!/bin/bash

PROJECT=$1
ARTIFACT_IMAGE=$2
ARTIFACT_IMAGE_TAG=$3
SEQUENCE=$4

# shellcheck disable=SC1091
source test/utils.sh

echo "---------------------------------------------------------------------"
echo "- Trigger delivery for ${ARTIFACT_IMAGE}:${ARTIFACT_IMAGE_TAG}"
echo "---------------------------------------------------------------------"
echo ""

# trigger delivery for carts service
keptn trigger delivery --project="$PROJECT" --service=carts --image="${ARTIFACT_IMAGE}" --tag="${ARTIFACT_IMAGE_TAG}" --sequence="$SEQUENCE"
verify_test_step $? "keptn trigger delivery --project=${PROJECT} --service=carts - failed"

sleep 30

# the following stages / namespaces should have some pods in it
kubectl get pods -n "$PROJECT-dev"
kubectl get pods -n "$PROJECT-staging"
kubectl get pods -n "$PROJECT-prod-a"
kubectl get pods -n "$PROJECT-prod-b"

echo "Verifying that services have been deployed to all stages ..."

####################################
# Verify dev deployment            #
####################################

verify_sockshop_deployment "${PROJECT}" "dev" "${ARTIFACT_IMAGE}" "${ARTIFACT_IMAGE_TAG}" "${KEPTN_NAMESPACE}" "false"

echo "It might take a while for the service to be available on staging - waiting a bit"
sleep 30
echo "Still waiting ..."
sleep 30

####################################
# Verify staging deployment        #
####################################
verify_sockshop_deployment "${PROJECT}" "staging" "${ARTIFACT_IMAGE}" "${ARTIFACT_IMAGE_TAG}" "${KEPTN_NAMESPACE}" "true"

echo "It might take a while for the service to be available on production - waiting a bit"
sleep 30
echo "Still waiting ..."
sleep 30

####################################
# Verify prod-a deployment         #
####################################
verify_sockshop_deployment "${PROJECT}" "prod-a" "${ARTIFACT_IMAGE}" "${ARTIFACT_IMAGE_TAG}" "${KEPTN_NAMESPACE}" "true"

####################################
# Verify prod-b deployment         #
####################################
verify_sockshop_deployment "${PROJECT}" "prod-b" "${ARTIFACT_IMAGE}" "${ARTIFACT_IMAGE_TAG}" "${KEPTN_NAMESPACE}" "true"

echo ""
echo "-----------------------------------------"
echo "- Deployment of ${ARTIFACT_IMAGE_TAG}   -"
echo "- looks good!                           -"
echo "-----------------------------------------"
echo ""
