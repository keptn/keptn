
function setup_gcloud {
    if [ ! -d "$HOME/google-cloud-sdk/bin" ]; then rm -rf $HOME/google-cloud-sdk; export CLOUDSDK_CORE_DISABLE_PROMPTS=1; curl https://sdk.cloud.google.com | bash; fi
    source /home/travis/google-cloud-sdk/path.bash.inc
    gcloud --quiet version
    gcloud --quiet components update
    gcloud --quiet components update kubectl
    echo $GCLOUD_SERVICE_KEY | base64 --decode -i > ${HOME}/gcloud-service-key.json
    gcloud auth activate-service-account --key-file ${HOME}/gcloud-service-key.json
    verify_step $? "gcloud authentication failed."
}

function setup_glcoud_pr {
    gcloud --quiet config set project $PROJECT_NAME
    gcloud container clusters get-credentials $CLUSTER_PR_STATUSCHECK_NAME --zone $CLUSTER_PR_STATUSCHECK_ZONE --project $PROJECT_NAME
    export GCLOUD_USER=$(gcloud config get-value account)

    kubectl create clusterrolebinding travis-cluster-admin-binding --clusterrole=cluster-admin --user=$GCLOUD_USER || true
    # export REGISTRY_URL=$(kubectl describe svc docker-registry -n keptn | grep "IP:" | sed 's~IP:[ \t]*~~')
}

function install_helm {
    curl https://storage.googleapis.com/kubernetes-helm/helm-v2.12.3-linux-amd64.tar.gz --output helm-v2.12.3-linux-amd64.tar.gz
    tar -zxvf helm-v2.12.3-linux-amd64.tar.gz
    sudo mv linux-amd64/helm /usr/local/bin/helm
}

function install_yq {
    sudo add-apt-repository ppa:rmescandon/yq -y
    sudo apt update
    sudo apt install yq -y
}

function install_hub {
    sudo wget https://github.com/github/hub/releases/download/v2.6.0/hub-linux-amd64-2.6.0.tgz
    tar -xzf hub-linux-amd64-2.6.0.tgz
    sudo cp hub-linux-amd64-2.6.0/bin/hub /bin/
}

function install_sed {
    sudo apt install --reinstall sed
}

function setup_knative {    
    cd ./installer/scripts/
    ./setupKnative.sh $CLUSTER_NAME_NIGHTLY ${CLOUDSDK_COMPUTE_ZONE}
    cd ../..
}

function uninstall_keptn {
    cd ./installer/scripts
    ./uninstallKeptn.sh
    cd ../..
}

function setup_knative_pr {    
    cd ./installer/scripts/
    CLUSTER_IPV4_CIDR=$(gcloud container clusters describe ${CLUSTER_PR_STATUSCHECK_NAME} --zone=${CLUSTER_PR_STATUSCHECK_ZONE} | yq r - clusterIpv4Cidr)
    SERVICES_IPV4_CIDR=$(gcloud container clusters describe ${CLUSTER_PR_STATUSCHECK_NAME} --zone=${CLUSTER_PR_STATUSCHECK_ZONE} | yq r - servicesIpv4Cidr)
    source ./setupKnative.sh $CLUSTER_IPV4_CIDR $SERVICES_IPV4_CIDR
    cd ../..
}

function setup_keptn_pr {    
    cd ./installer/scripts/
    source ./setupKeptn.sh
    cd ../..
}

function export_names {
    export EVENT_BROKER_NAME=$(kubectl describe ksvc event-broker -n keptn | grep -m 1 "Name:" | sed 's~Name:[ \t]*~~')
    ./test/assertEquals.sh $EVENT_BROKER_NAME event-broker

    export EVENT_BROKER_EXT_NAME=$(kubectl describe ksvc event-broker-ext -n keptn | grep -m 1 "Name:" | sed 's~Name:[ \t]*~~')
    ./test/assertEquals.sh $EVENT_BROKER_EXT_NAME event-broker-ext
    
    export AUTHENTICATOR_NAME=$(kubectl describe ksvc authenticator -n keptn | grep -m 1 "Name:" | sed 's~Name:[ \t]*~~')
    ./test/assertEquals.sh $AUTHENTICATOR_NAME authenticator

    export CONTROL_NAME=$(kubectl describe ksvc control -n keptn | grep -m 1 "Name:" | sed 's~Name:[ \t]*~~')
    ./test/assertEquals.sh $CONTROL_NAME control
}

function verify_step() {
  if [[ $1 != '0' ]]; then
    print_error "$2"
    exit 1
  fi
}