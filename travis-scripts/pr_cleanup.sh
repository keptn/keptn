#!/usr/bin/env bash

# prints the full command before output of the command.
set -x

cd ./core/
kubectl delete -f auth/config/authenticator.yaml
kubectl delete -f control/config/gen
kubectl delete -f eventbroker/config/event-broker.yaml
kubectl delete -f eventbroker-ext/config/event-broker-ext.yaml
kubectl delete svc,deployments --all -n keptn
kubectl delete channel --all -n keptn
kubectl delete clusterchannelprovisioner --all -n keptn
kubectl delete clusterrolebinding travis-cluster-admin-binding
kubectl delete svc knative-ingressgateway -n istio-system
kubectl delete deploy knative-ingressgateway -n istio-system
