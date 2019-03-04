#!/bin/sh

REGISTRY_URI=$(kubectl describe svc docker-registry -n cicd | grep IP: | sed 's~IP:[ \t]*~~')

rm -f config/gen/event-broker.yaml

cat config/event-broker.yaml | \
  sed 's~REGISTRY_URI_PLACEHOLDER~'"$REGISTRY_URI"'~' >> config/gen/event-broker.yaml 
  
kubectl delete -f config/gen/event-broker.yaml
kubectl apply -f config/gen/event-broker.yaml