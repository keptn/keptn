#!/bin/bash

source test/utils.sh

echo "---------------------------------------------"
echo "- Sending new artifact for mongo            -"
echo "---------------------------------------------"
echo ""

# send new artifcat for database
keptn send event new-artifact --project=$PROJECT --service=carts-db --image=mongo
verify_test_step $? "keptn send event new-artifact --project=${PROJECT} --service=carts-db --image=mongo - failed"

# wait until mongodb has been deployed
wait_for_deployment_in_namespace "carts-db" "$PROJECT-dev"
verify_test_step $? "Deployment carts-db not available, exiting ..."

# send new artifact for carts
test/utils/send_new_artifact_sockshop.sh $PROJECT docker.io/keptnexamples/carts 0.10.1

# wait before sending the next artifact
echo "Waiting 30sec before continue ..."
sleep 30

echo "Send new artifact now"
test/utils/send_new_artifact_sockshop.sh $PROJECT docker.io/keptnexamples/carts 0.10.3

exit 0
