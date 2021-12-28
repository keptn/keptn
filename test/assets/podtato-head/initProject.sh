#!/usr/bin/env bash

PROJECT="pod-tato-head"
IMAGE="ghcr.io/podtato-head/podtatoserver"
VERSION="$2"

case "$1" in
  "create-project")
    echo "Creating keptn project $PROJECT"
    echo keptn create project "${PROJECT}" --shipyard=./shipyard.yaml   
    keptn create project "${PROJECT}" --shipyard=./shipyard.yaml
    ;;
  "create-service")
    echo "Creating keptn service helloserver in project ${PROJECT}"
    keptn create service helloserver --project="${PROJECT}"
    ;;
  "add-resource")
    echo "Adding keptn resource to service helloserver in project ${PROJECT}"
    keptn add-resource --service=helloserver --project="${PROJECT}" --all-stages --resource=helm-charts/helloserver.tgz --resourceUri=helm-charts/helloserver.tgz
    ;;
  "first-deploy-service")
    echo "Deploying keptn service helloserver in project ${PROJECT}"
    keptn trigger delivery --project="${PROJECT}" --service=helloserver --image="${IMAGE}" --tag=v0.1.1
    ;;
  "deploy-service")
    echo "Deploying keptn service helloserver in project ${PROJECT}"
    echo keptn trigger delivery --project="${PROJECT}" --service=helloserver --image="${IMAGE}" --tag=v"${VERSION}"
    keptn trigger delivery --project="${PROJECT}" --service=helloserver --image="${IMAGE}" --tag=v"${VERSION}"
    ;;    
  "upgrade-service")
    echo "Upgrading keptn service helloserver in project ${PROJECT}"
    keptn trigger delivery --project="${PROJECT}" --service=helloserver --image="${IMAGE}" --tag=v0.1.0
    ;;
  "slow-build")
    echo "Deploying slow build version of helloserver in project ${PROJECT}"
    keptn trigger delivery --project="${PROJECT}" --service=helloserver --image="${IMAGE}" --tag=v0.1.2
    ;;
  "add-quality-gates")
    echo "Adding keptn quality-gates to project ${PROJECT}"
    keptn add-resource --project=${PROJECT} --stage=hardening --service=helloserver --resource=prometheus/sli.yaml --resourceUri=prometheus/sli.yaml
    keptn add-resource --project=${PROJECT} --stage=hardening --service=helloserver --resource=slo.yaml --resourceUri=slo.yaml
    ;;
  "add-jmeter-tests")
    echo "Adding jmeter load tests to project ${PROJECT}"
    keptn add-resource --project=${PROJECT} --stage=hardening --service=helloserver --resource=jmeter/load.jmx --resourceUri=jmeter/load.jmx
    keptn add-resource --project=${PROJECT} --stage=hardening --service=helloserver --resource=jmeter/jmeter.conf.yaml --resourceUri=jmeter/jmeter.conf.yaml
    ;;
esac