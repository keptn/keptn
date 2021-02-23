#!/bin/bash

source test/utils.sh

KEPTN_EXAMPLES_BRANCH=${KEPTN_EXAMPLES_BRANCH:-"master"}

echo "Testing onboarding for project $PROJECT"


# verify that the project does not exist yet via the Keptn API
response=$(curl -X GET "${KEPTN_ENDPOINT}/controlPlane/v1/project/${PROJECT}" -H  "accept: application/json" -H  "x-token: ${KEPTN_API_TOKEN}" -k 2>/dev/null | jq -r '.projectName')

if [[ "$response" == "${PROJECT}" ]]; then
  echo "Project ${PROJECT} already exists. Please delete it using:"
  echo "keptn delete project ${PROJECT}"
  exit 2
fi

# test keptn create-project and onboard
rm -rf examples
git clone --branch ${KEPTN_EXAMPLES_BRANCH} https://github.com/keptn/examples --single-branch
cd examples/onboarding-carts

echo "Creating a new project without Git upstream"
keptn create project $PROJECT --shipyard=../../test/assets/shipyard_onboard_service.yaml
verify_test_step $? "keptn create project ${PROJECT} failed."
sleep 10

###########################################
# onboard carts                           #
###########################################
keptn onboard service carts --project=$PROJECT --chart=./carts
verify_test_step $? "keptn onboard service carts failed"
sleep 10

# add functional tests
keptn add-resource --project=$PROJECT --service=carts --stage=dev --resource=jmeter/basiccheck.jmx --resourceUri=jmeter/basiccheck.jmx
# add performance tests
keptn add-resource --project=$PROJECT --service=carts --stage=staging --resource=jmeter/basiccheck.jmx --resourceUri=jmeter/load.jmx

###########################################
# onboard carts-db                        #
###########################################
keptn onboard service carts-db --project=$PROJECT --chart=./carts-db
verify_test_step $? "keptn onboard service carts-db failed"

echo "Onboarding done ✓"

exit 0
