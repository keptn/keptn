#!/bin/bash

# shellcheck disable=SC1091
source test/utils.sh

echo "get services before delete services command:"
services=$(keptn get services --project="$PROJECT" -ojson)
echo "$services" | jq .

number_of_services=$(keptn get services --project="$PROJECT" -ojson | jq -r -s '. | length')

verify_value "number of services" "$number_of_services" 8

echo "Testing the keptn delete service command"
keptn delete service "$SERVICE" --project="$PROJECT"
verify_test_step $? "keptn delete service ${SERVICE} --project=${PROJECT} - failed"

echo "get services after delete services command:"
services=$(keptn get services --project="$PROJECT" -ojson)
echo "$services" | jq .

number_of_services=$(keptn get services --project="$PROJECT" -ojson | jq -r -s '. | length')

verify_value "number of services" "$number_of_services" 4

echo "Testing the keptn delete project command"
keptn delete project "$PROJECT"
verify_test_step $? "keptn delete project ${PROJECT} - failed"

# Note: delete project only deletes the project from configuration-service, but it will not "un-deploy" the project

exit 0
