#!/bin/bash

# shellcheck disable=SC1091
source test/utils.sh

function debuglogs() {
  echo "::group::Namespaces"
  kubectl get namespaces
  echo "::endgroup::"
  echo "::group::Pods"
  echo "Show pods in all sockshop-* namespaces"
  kubectl get pods -n "$PROJECT-dev"
  kubectl get pods -n "$PROJECT-staging"
  kubectl get pods -n "$PROJECT-prod-a"
  kubectl get pods -n "$PROJECT-prod-b"
  echo "::endgroup::"
  echo "::group::Deployments"
  echo "Show deployments in all sockshop-* namespaces"
  kubectl get deployments -n "$PROJECT-dev" -owide
  kubectl get deployments -n "$PROJECT-staging" -owide
  kubectl get deployments -n "$PROJECT-prod-a" -owide
  kubectl get deployments -n "$PROJECT-prod-b" -owide
  echo "::endgroup::"
  echo "::group::CloudEvents Carts DB"
  echo "Show CloudEvents for service carts-db"
  keptn get event --project="$PROJECT" --service=carts-db
  echo "::endgroup::"
  echo "::group::CloudEvents Carts"
  echo "Show CloudEvents for service carts"
  keptn get event --project="$PROJECT" --service=carts
  echo "::endgroup::"
}
trap debuglogs EXIT SIGINT

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
verify_test_step $? "Delivery of carts 0.10.1 failed"

# wait before sending the next artifact
echo "Waiting 30sec before continue ..."
sleep 30

echo "Trigger delivery now"
test/utils/trigger_delivery_sockshop.sh "$PROJECT" docker.io/keptnexamples/carts 0.10.3 delivery
verify_test_step $? "Delivery of carts 0.10.3 failed"

exit 0
