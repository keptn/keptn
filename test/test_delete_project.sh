#!/bin/bash

source test/utils.sh

echo "Testing deletion of project"

keptn delete project $PROJECT
verify_test_step $? "Deleting project failed"

# Note: delete project only deletes the project from configuration-service, but it will not "un-deploy" the project

exit 0
