#!/bin/bash

# shellcheck disable=SC1091
source test/utils.sh

echo "---------------------------------------------"
echo "- Trigger delivery for mongo            -"
echo "---------------------------------------------"
echo ""

# trigger delivery for database (include tag in image parameter to test if combining image/tag works)
keptn trigger delivery --project="$PROJECT" --service=carts-db --image=mongo:latest --sequence=delivery-direct
verify_test_step $? "keptn trigger delivery --project=${PROJECT} --service=carts-db --image=mongo - failed"

# wait until mongodb has been deployed
wait_for_deployment_in_namespace "carts-db" "$PROJECT-dev"
verify_test_step $? "Deployment carts-db not available, exiting ..."

# trigger delivery for carts
test/utils/trigger_delivery_sockshop.sh "$PROJECT" docker.io/keptnexamples/carts 0.10.1 delivery

# wait before sending the next artifact
echo "Waiting 30sec before continue ..."
sleep 30

echo "Trigger delivery now"
test/utils/trigger_delivery_sockshop.sh "$PROJECT" docker.io/keptnexamples/carts 0.10.3 delivery

exit 0
