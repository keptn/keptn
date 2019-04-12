#!/bin/sh

REGISTRY_URI=$(kubectl describe svc docker-registry -n keptn | grep IP: | sed 's~IP:[ \t]*~~')

rm -f config/gen/event-broker-build.yaml

cat config/event-broker-build.yaml | \
  sed 's~REGISTRY_URI_PLACEHOLDER~'"$REGISTRY_URI"'~' >> config/gen/event-broker-build.yaml 

kubectl delete -f config/gen/event-broker-build.yaml --ignore-not-found
kubectl apply -f config/gen/event-broker-build.yaml
