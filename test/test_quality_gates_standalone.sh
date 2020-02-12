#!/bin/bash

source test/utils.sh

echo "Testing quality gates standalone for project $PROJECT ..."

# Test keptn create-project and create service
rm -rf examples
git clone --branch master https://github.com/keptn/examples --single-branch
cd examples/onboarding-carts

echo "Creating a new project without git upstream"
keptn create project $PROJECT --shipyard=./shipyard-quality-gates.yaml
verify_test_step $? "keptn create project command failed."
sleep 10

###########################################
# create service catalogue                #
###########################################
keptn create service catalogue --project=$PROJECT
verify_test_step $? "keptn create service catalogue failed."
sleep 10

# add SLO file for service
keptn add-resource --project=$PROJECT --stage=hardening --service=catalogue --resource=slo-quality-gates.yaml --resourceUri=slo.yaml
verify_test_step $? "keptn add-resource failed."

# send start evaluation command
response=$(keptn send event start-evaluation --project=$PROJECT --stage=hardening --service=catalogue --timeframe=5m)

echo $response

keptn_context_id=$(echo $response | awk -F'Keptn context:' '{ print $2 }' | xargs)

sleep 10
# parse output of above command and extract keptn context)
response=$(keptn get event evaluation-done --keptn-context=${keptn_context_id})

# ToDo: right now get event evaluation-done does not work as lighthouse does not send the event when there is no sli provider
#       see issue https://github.com/keptn/keptn/issues/1212

#if [[ $? -eq 0 ]]; then
#  echo "Got result"
#  echo $response
#else
#  exit 1
#fi

exit 0
