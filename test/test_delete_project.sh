#!/bin/bash

source test/utils.sh

echo "Testing the keptn delete project command"

keptn delete project $PROJECT
verify_test_step $? "keptn delete project ${PROJECT} - failed"

# Note: delete project only deletes the project from configuration-service, but it will not "un-deploy" the project

exit 0
