#!/bin/bash

source test/utils.sh

echo "Testing onboarding for project $PROJECT ..."

# Test keptn create-project and onboard
rm -rf examples
git clone --branch master https://github.com/keptn/examples --single-branch
cd examples/onboarding-carts

echo "Creating a new project without git upstream"
keptn create project $PROJECT --shipyard=./shipyard.yaml
verify_test_step $? "keptn create project command failed."
sleep 10

###########################################
# onboard carts                           #
###########################################
keptn onboard service carts --project=$PROJECT --chart=./carts
verify_test_step $? "keptn onboard carts failed."
sleep 10

# add functional tests
keptn add-resource --project=sockshop --service=carts --stage=dev --resource=jmeter/basiccheck.jmx --resourceUri=jmeter/basiccheck.jmx
# add performance tests
keptn add-resource --project=sockshop --service=carts --stage=staging --resource=jmeter/basiccheck.jmx --resourceUri=jmeter/load.jmx

###########################################
# onboard carts-db                        #
###########################################
keptn onboard service carts-db --project=$PROJECT --chart=./carts-db --deployment-strategy=direct
verify_test_step $? "keptn onboard carts-db failed."
sleep 10

# check which namespaces exist
echo "Verifying that the following namespaces are available:"

verify_namespace_exists "$PROJECT-dev"
verify_namespace_exists "$PROJECT-staging"
verify_namespace_exists "$PROJECT-prod"

echo "Onboarding done!"

exit 0
