#!/bin/sh
REGISTRY_URI=$(kubectl describe svc docker-registry -n keptn | grep IP: | sed 's~IP:[ \t]*~~')
CHANNEL_URI=$(kubectl describe channel keptn-channel -n keptn | grep "Hostname:" | sed 's~[ \t]*Hostname:[ \t]*~~')

rm -f config/gen/control-build.yaml

cat config/control-build.yaml | \
  sed 's~CHANNEL_URI_PLACEHOLDER~'"$CHANNEL_URI"'~' | \
  sed 's~REGISTRY_URI_PLACEHOLDER~'"$REGISTRY_URI"'~' >> config/gen/control-build.yaml 
  
kubectl delete -f config/gen/control-build.yaml --ignore-not-found
kubectl apply -f config/gen/control-build.yaml