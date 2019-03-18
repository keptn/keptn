#!/bin/sh
REGISTRY_URI=$(kubectl describe svc docker-registry -n keptn | grep IP: | sed 's~IP:[ \t]*~~')
CHANNEL_URI=$2

rm -f config/gen/control.yaml

cat config/control.yaml | \
  sed 's~CHANNEL_URI_PLACEHOLDER~'"$CHANNEL_URI"'~' | \
  sed 's~REGISTRY_URI_PLACEHOLDER~'"$REGISTRY_URI"'~' >> config/gen/control.yaml 
  
kubectl delete -f config/gen/control.yaml --ignore-not-found
kubectl apply -f config/gen/control.yaml