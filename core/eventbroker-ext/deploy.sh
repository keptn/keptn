#!/bin/sh

REGISTRY_URI=$(kubectl describe svc docker-registry -n keptn | grep IP: | sed 's~IP:[ \t]*~~')

rm -f config/gen/event-broker-ext.yaml

cat config/event-broker-ext.yaml | \
  sed 's~REGISTRY_URI_PLACEHOLDER~'"$REGISTRY_URI"'~' >> config/gen/event-broker-ext.yaml 
  
kubectl delete -f config/gen/event-broker-ext.yaml --ignore-not-found
kubectl apply -f config/gen/event-broker-ext.yaml