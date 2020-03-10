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

verify_test_step $? "ERROR: Command keptn get event evaluation-done --keptn-context=${keptn_context_id}) returned a non-zero exit code"

echo $response | grep "no SLI-provider configured for project $PROJECT"

if [[ $? -eq 0 ]]; then
  echo "Result was as expected (no SLI-provider found)"
else
  echo "Got result"
  echo $response
  echo "ERROR: Expected string 'no SLI-provider configured for project $PROJECT' in output"
  exit 1
fi

# okay so far, now we install a SLI provider
kubectl apply -f https://raw.githubusercontent.com/keptn-contrib/dynatrace-sli-service/0.3.0/deploy/service.yaml
sleep 10

wait_for_deployment_in_namespace "dynatrace-sli-service" "keptn"

echo "Sending 'configure monitoring'..."

# now configure monitoring for Dynatrace
keptn configure monitoring dynatrace --project=$PROJECT
sleep 5
# this should set the configmap

kubectl -n keptn get configmap lighthouse-config-$PROJECT -oyaml
verify_test_step $? "ERROR: Could not find configmap lighthouse-config-$PROJECT (this is expected to be created by keptn configure monitoring dynatrace --project=$PROJECT)"

# now that this is set, let's send the start evaluation command again
response=$(keptn send event start-evaluation --project=$PROJECT --stage=hardening --service=catalogue --timeframe=5m)

echo $response

keptn_context_id=$(echo $response | awk -F'Keptn context:' '{ print $2 }' | xargs)

sleep 10
# parse output of above command and extract keptn context)
response=$(keptn get event evaluation-done --keptn-context=${keptn_context_id})

echo $response

verify_test_step $? "ERROR: Command keptn get event evaluation-done --keptn-context=${keptn_context_id}) returned a non-zero exit code"



exit 0
