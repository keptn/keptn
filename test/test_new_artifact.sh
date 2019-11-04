#!/bin/bash

source test/utils.sh

# send new artifcat for database
keptn send event new-artifact --project=$PROJECT --service=carts-db --image=mongo
verify_test_step $? "Send event new-artifact for carts-db failed"


# send new artifact for carts service
keptn send event new-artifact --project=$PROJECT --service=carts --image=docker.io/keptnexamples/carts --tag=0.9.1
verify_test_step $? "Send event new-artifact for carts failed"


sleep 10

# the following stages / namespaces should have some pods in it
kubectl get pods -n "$PROJECT-dev"
kubectl get pods -n "$PROJECT-staging"
kubectl get pods -n "$PROJECT-production"

echo "Verifying that pods have been deployed to all the stages..."

# the following will not work, as we only onboarded the service, but we didnt create a new artifact
verify_deployment_in_namespace "carts" "$PROJECT-dev"
verify_deployment_in_namespace "carts-db" "$PROJECT-dev"

verify_deployment_in_namespace "carts" "$PROJECT-staging"
verify_deployment_in_namespace "carts-db" "$PROJECT-staging"

verify_deployment_in_namespace "carts" "$PROJECT-production"
verify_deployment_in_namespace "carts-db" "$PROJECT-production"

exit 0
