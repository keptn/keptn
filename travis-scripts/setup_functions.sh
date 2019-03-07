
function setup_gcloud {
    if [ ! -d "$HOME/google-cloud-sdk/bin" ]; then rm -rf $HOME/google-cloud-sdk; export CLOUDSDK_CORE_DISABLE_PROMPTS=1; curl https://sdk.cloud.google.com | bash; fi
    source /home/travis/google-cloud-sdk/path.bash.inc
    gcloud --quiet version
    gcloud --quiet components update
    gcloud --quiet components update kubectl
    echo $GCLOUD_SERVICE_KEY | base64 --decode -i > ${HOME}/gcloud-service-key.json
    gcloud auth activate-service-account --key-file ${HOME}/gcloud-service-key.json
    gcloud container clusters get-credentials $CLUSTER_PR_STATUSCHECK_NAME --zone $CLUSTER_PR_STATUSCHECK_ZONE --project $PROJECT_NAME
    export GCLOUD_USER=$(gcloud config get-value account)

    kubectl create clusterrolebinding travis-cluster-admin-binding --clusterrole=cluster-admin --user=$GCLOUD_USER || true
    export REGISTRY_URL=$(kubectl describe svc docker-registry -n keptn | grep "IP:" | sed 's~IP:[ \t]*~~')
}

# Install yq
function install_yq {
    sudo add-apt-repository ppa:rmescandon/yq -y
    sudo apt update
    sudo apt install yq -y
}

function setup_knative {
    # First delete potential knative remainings
    kubectl delete ns knative-build knative-eventing knative-monitoring knative-serving knative-sources || true
    kubectl delete svc knative-ingressgateway -n istio-system || true
    kubectl delete deploy knative-ingressgateway -n istio-system || true
    pwd
    cd ./install/scripts/
    ./setupKnative.sh $JENKINS_USER $JENKINS_PASSWORD $REGISTRY_URL
    export EVENT_BROKER_NAME=$(kubectl describe ksvc event-broker -n keptn | grep -m 1 "Name:" | sed 's~Name:[ \t]*~~')
    ./../../test/assertEquals.sh $EVENT_BROKER_NAME event-broker
    
    export AUTHENTICATOR_NAME=$(kubectl describe ksvc authenticator -n keptn | grep -m 1 "Name:" | sed 's~Name:[ \t]*~~')
    ./../../test/assertEquals.sh $AUTHENTICATOR_NAME authenticator

    export CONTROL_NAME=$(kubectl describe ksvc control -n keptn | grep -m 1 "Name:" | sed 's~Name:[ \t]*~~')
    ./../../test/assertEquals.sh $CONTROL_NAME control
    cd ../..
}

function execute_core_component_tests {
    # execute unit tests for core components
    
    # Control
    cd ./core/control
    npm install
    npm run test
    
    # Auth
    cd ../auth
    npm install
    npm run test
    
    # Event Broker
    cd ../eventbroker
    npm install
    npm run test

    # Event Broker (ext)
    cd ../eventbroker-ext
    npm install
    npm run test
    
    cd ../..
}

function execute_cli_tests {

    cd cli
    ENDPOINT="$(kubectl get ksvc control -n keptn -o=yaml | yq r - status.domain)"
    while [ "$ENDPOINT" = "null" ]; do sleep 30; ENDPOINT="$(kubectl get ksvc control -n keptn -o=yaml | yq r - status.domain)"; echo "waiting for control service"; done
    printf "https://" > ~/.keptnmock
    kubectl get ksvc control -n keptn -o=yaml  | yq r - status.domain >> ~/.keptnmock

    SEC="$(kubectl get secret keptn-api-token  -n keptn -o=yaml | yq r - data.keptn-api-token | base64 --decode)"
    echo "${SEC}" >> ~/.keptnmock
        
        # execute GO tests
    go test ${gobuild_args} -timeout 240s ./...
    cd ..
}