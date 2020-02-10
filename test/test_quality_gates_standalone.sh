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
keptn send event start-evaluation --project=$PROJECT --stage=hardening --service=catalogue --timeframe=5m

sleep 10
# Todo: parse output of above command and extract keptn context)
# keptn get event evaluation-done --keptn-context=...

exit 0
