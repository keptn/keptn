#!/bin/bash

source test/utils.sh

echo "---------------------------------------------"
echo "- Sending new artifact for mongo            -"
echo "---------------------------------------------"
echo ""

# send new artifcat for database
keptn send event new-artifact --project=$PROJECT --service=carts-db --image=mongo
verify_test_step $? "Send event new-artifact for carts-db failed"

# wait until mongodb has been deployed
wait_for_deployment_in_namespace "carts-db" "$PROJECT-dev"
verify_test_step $? "Deployment carts-db not available, exiting..."

# okay, now we can start with carts
test/utils/send_new_artifact_sockshop.sh $PROJECT docker.io/keptnexamples/carts 0.10.1

# wait before sending the next artifact
echo "Waiting a little bit before we continue..."
sleep 30
echo "Continuing now!"

test/utils/send_new_artifact_sockshop.sh $PROJECT docker.io/keptnexamples/carts 0.10.3

exit 0
