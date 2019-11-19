#!/bin/bash

source test/utils.sh

ARTIFACT_IMAGE_TAG=0.10.1

echo "---------------------------------------------"
echo "Sending new artifact for mongo"

# send new artifcat for database
keptn send event new-artifact --project=$PROJECT --service=carts-db --image=mongo
verify_test_step $? "Send event new-artifact for carts-db failed"


echo "---------------------------------------------"
echo "Sending new artifact for docker.io/keptnexamples/carts:${ARTIFACT_IMAGE_TAG}"
echo ""

# send new artifact for carts service
keptn send event new-artifact --project=$PROJECT --service=carts --image=docker.io/keptnexamples/carts --tag=$ARTIFACT_IMAGE_TAG
verify_test_step $? "Send event new-artifact for carts failed"

# a new artifact for the carts service might take a while, so lets wait
sleep 10

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

wait_for_deployment_in_namespace "carts" "$PROJECT-dev"
verify_deployment_in_namespace "carts-db" "$PROJECT-dev"
verify_pod_in_namespace "carts" "$PROJECT-dev"
verify_pod_in_namespace "carts-db" "$PROJECT-dev"

# get URL for that deployment
DEV_URL=$(echo http://carts.sockshop-dev.$(kubectl get cm keptn-domain -n keptn -o=jsonpath='{.data.app_domain}'))
# try to access that URL
wait_for_url $DEV_URL/health
verify_test_step $? "Trying to access $DEV_URL/health failed"

# verify image name of carts deployment
verify_image_of_deployment "carts" "$PROJECT-dev" docker.io/keptnexamples/carts:$ARTIFACT_IMAGE_TAG
verify_test_step $? "Wrong image for deployment carts in $PROJECT-dev"

# It might take a while for the service to be available on staging - so lets wait a bit
sleep 10

####################################
# Verify staging deployment        #
####################################
echo "---------------------------------------------"
echo "Checking staging deployment"
echo ""

wait_for_deployment_in_namespace "carts" "$PROJECT-staging"
verify_deployment_in_namespace "carts-db" "$PROJECT-staging"
verify_pod_in_namespace "carts" "$PROJECT-staging"
verify_pod_in_namespace "carts-db" "$PROJECT-staging"

# get URL for that deployment
STAGING_URL=$(echo http://carts.sockshop-staging.$(kubectl get cm keptn-domain -n keptn -o=jsonpath='{.data.app_domain}'))
# try to access that URL
wait_for_url $STAGING_URL/health
verify_test_step $? "Trying to access $STAGING_URL/health failed"

# verify image name of carts deployment
verify_image_of_deployment "carts" "$PROJECT-staging" docker.io/keptnexamples/carts:$ARTIFACT_IMAGE_TAG
verify_test_step $? "Wrong image for deployment carts in $PROJECT-staging"

# It might take a while for the service to be available on production - so lets wait a bit
sleep 60

####################################
# Verify produciton deployment     #
####################################
echo "---------------------------------------------"
echo "Checking production deployment"
echo ""

wait_for_deployment_in_namespace "carts" "$PROJECT-production"
verify_deployment_in_namespace "carts-db" "$PROJECT-production"
verify_pod_in_namespace "carts" "$PROJECT-production"
verify_pod_in_namespace "carts-db" "$PROJECT-production"

# get URL for that deployment
PRODUCTION_URL=$(echo http://carts.sockshop-production.$(kubectl get cm keptn-domain -n keptn -o=jsonpath='{.data.app_domain}'))
# try to access that URL
wait_for_url $PRODUCTION_URL/health
verify_test_step $? "Trying to access $PRODUCTION_URL/health failed"

# verify image name of carts deployment
verify_image_of_deployment "carts" "$PROJECT-production" docker.io/keptnexamples/carts:$ARTIFACT_IMAGE_TAG
verify_test_step $? "Wrong image for deployment carts in $PROJECT-production"

echo ""
echo "-----------------------------------------"
echo "- Looks good!                           -"
echo "-----------------------------------------"
echo ""

exit 0
