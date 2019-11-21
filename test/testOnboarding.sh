#!/bin/bash
source ./utils.sh

echo "Testing onboarding..."

PROJECT=sockshop

# Delete old project
git ls-remote https://github.com/$GITHUB_ORG_NIGHTLY/$PROJECT > /dev/null 2>&1
if [ $? = 0 ]; then 
    echo "Delete project $PROJECT" 
    GITHUB_USER=GITHUB_USER_NAME_NIGHTLY
    GITHUB_PASSWORD=GITHUB_TOKEN_NIGHTLY
    hub delete -y $GITHUB_ORG_NIGHTLY/$PROJECT
    echo "Finished deleting project $PROJECT"
else 
    echo "No project to delete"
fi

# Test keptn create-project and onboard
rm -rf examples
git clone https://github.com/keptn/examples
cd examples
cd onboarding-carts

echo "Creating a new project without git upstream"
keptn create project $PROJECT --shipyard=shipyard.yaml
verify_test_step $? "keptn create project command failed."
sleep 10

keptn onboard service carts --project=$PROJECT --chart=./carts
sleep 10

keptn onboard service carts-db --project=$PROJECT --chart=./carts-db --deployment-strategy=direct
sleep 10

# check which namespaces exist
kubectl get namespaces
kubectl get namespaces -L istio-injection

# the following stages should have been created
# kubectl get pods -n "$PROJECT-dev"
# kubectl get pods -n "$PROJECT-production"
# kubectl get pods -n "$PROJECT-staging"

# the following will not work, as we only onboarded the service, but we didnt create a new artifact
#wait_for_deployment_in_namespace "carts" "$PROJECT-dev"
#wait_for_deployment_in_namespace "carts-db" "$PROJECT-dev"

#wait_for_deployment_in_namespace "carts" "$PROJECT-dev"
#wait_for_deployment_in_namespace "carts-db" "$PROJECT-dev"


# newman only checks several github related things - not applicable any more as of version 0.5.0
#cd ../..
#npm install newman
#go get gopkg.in/mikefarah/yq.v2

#$GOPATH/bin/yq.v2 w keptn.postman_environment.json values[0].value $GITHUB_ORG_NIGHTLY | $GOPATH/bin/yq.v2  - w values[1].value $PROJECT | $GOPATH/bin/yq.v2  - w values[2].value $GITHUB_CLIENT_ID_NIGHTLY |  $GOPATH/bin/yq.v2  - w values[3].value $GITHUB_CLIENT_SECRET_NIGHTLY | $GOPATH/bin/yq.v2 - r -j > keptn.postman_environment_tmp.json
#rm keptn.postman_environment.json
#mv keptn.postman_environment_tmp.json keptn.postman_environment.json
#node_modules/.bin/newman run keptn.postman_collection.json -e keptn.postman_environment.json
