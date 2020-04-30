#!/bin/bash

PROJECT=$1
ARTIFACT_IMAGE=$2
ARTIFACT_IMAGE_TAG=$3

source test/utils.sh

echo "---------------------------------------------------------------------"
echo "- Sending new artifact for ${ARTIFACT_IMAGE}:${ARTIFACT_IMAGE_TAG}"
echo "---------------------------------------------------------------------"
echo ""

# send new artifact for carts service
keptn send event new-artifact --project=$PROJECT --service=carts --image=${ARTIFACT_IMAGE} --tag=$ARTIFACT_IMAGE_TAG
verify_test_step $? "Send event new-artifact for carts failed"

# a new artifact for the carts service might take a while, so lets wait
sleep 30

# the following stages / namespaces should have some pods in it
kubectl get pods -n "$PROJECT-dev"
kubectl get pods -n "$PROJECT-staging"
kubectl get pods -n "$PROJECT-production"

echo "Verifying that services have been deployed to all stages..."

####################################
# Verify dev deployment            #
####################################
echo "---------------------------------------------"
echo "Checking dev deployment"
echo ""

wait_for_deployment_in_namespace "carts-db" "$PROJECT-dev"
wait_for_deployment_in_namespace "carts" "$PROJECT-dev"
verify_pod_in_namespace "carts" "$PROJECT-dev"
verify_test_step $? "Pod carts not found, exiting..."
verify_pod_in_namespace "carts-db" "$PROJECT-dev"
verify_test_step $? "Pod carts-db not found, exiting..."

# get URL for that deployment
DEV_URL=$(echo https://carts.${PROJECT}-dev.$(kubectl get cm keptn-domain -n keptn -o=jsonpath='{.data.app_domain}'))
# try to access that URL
wait_for_url $DEV_URL/health
verify_test_step $? "Trying to access $DEV_URL/health failed"

# verify image name of carts deployment
verify_image_of_deployment "carts" "$PROJECT-dev" ${ARTIFACT_IMAGE}:$ARTIFACT_IMAGE_TAG
verify_test_step $? "Wrong image for deployment carts in $PROJECT-dev"

echo "It might take a while for the service to be available on staging - waiting a bit"
sleep 30
echo "Still waiting..."
sleep 30

####################################
# Verify staging deployment        #
####################################
echo "---------------------------------------------"
echo "Checking staging deployment"
echo ""

wait_for_deployment_in_namespace "carts-db" "$PROJECT-staging"
wait_for_deployment_in_namespace "carts" "$PROJECT-staging"
verify_pod_in_namespace "carts" "$PROJECT-staging"
verify_test_step $? "Pod carts not found, exiting..."
verify_pod_in_namespace "carts-primary" "$PROJECT-staging"
verify_test_step $? "Pod carts-primary not found, exiting..."
verify_pod_in_namespace "carts-db" "$PROJECT-staging"
verify_test_step $? "Pod carts-db not found, exiting..."

# get URL for that deployment
STAGING_URL=$(echo https://carts.${PROJECT}-staging.$(kubectl get cm keptn-domain -n keptn -o=jsonpath='{.data.app_domain}'))
# try to access that URL
wait_for_url $STAGING_URL/health
verify_test_step $? "Trying to access $STAGING_URL/health failed"

# verify image name of carts deployment
verify_image_of_deployment "carts" "$PROJECT-staging" ${ARTIFACT_IMAGE}:$ARTIFACT_IMAGE_TAG
verify_test_step $? "Wrong image for deployment carts in $PROJECT-staging"

echo "It might take a while for the service to be available on production - waiting a bit"
sleep 30
echo "Still waiting..."
sleep 30

####################################
# Verify produciton deployment     #
####################################
echo "---------------------------------------------"
echo "Checking production deployment"
echo ""

wait_for_deployment_in_namespace "carts-db" "$PROJECT-production"
wait_for_deployment_in_namespace "carts" "$PROJECT-production"
verify_pod_in_namespace "carts" "$PROJECT-production"
verify_test_step $? "Pod carts not found, exiting..."
verify_pod_in_namespace "carts-primary" "$PROJECT-production"
verify_test_step $? "Pod carts-primary not found, exiting..."
verify_pod_in_namespace "carts-db" "$PROJECT-production"
verify_test_step $? "Pod carts-db not found, exiting..."

# get URL for that deployment
PRODUCTION_URL=$(echo https://carts.${PROJECT}-production.$(kubectl get cm keptn-domain -n keptn -o=jsonpath='{.data.app_domain}'))
# try to access that URL
wait_for_url $PRODUCTION_URL/health
verify_test_step $? "Trying to access $PRODUCTION_URL/health failed"

# verify image name of carts deployment
verify_image_of_deployment "carts" "$PROJECT-production" ${ARTIFACT_IMAGE}:$ARTIFACT_IMAGE_TAG
verify_test_step $? "Wrong image for deployment carts in $PROJECT-production"

echo ""
echo "-----------------------------------------"
echo "- Deployment of ${ARTIFACT_IMAGE_TAG}         - "
echo "- looks good!                           -"
echo "-----------------------------------------"
echo ""
