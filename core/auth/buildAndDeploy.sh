#!/bin/sh
REGISTRY_URI=$(kubectl describe svc docker-registry -n keptn | grep IP: | sed 's~IP:[ \t]*~~')

rm -f config/gen/authenticator-build.yaml

cat config/authenticator-build.yaml | \
  sed 's~REGISTRY_URI_PLACEHOLDER~'"$REGISTRY_URI"'~' >> config/gen/authenticator-build.yaml 
  
kubectl delete -f config/gen/authenticator-build.yaml --ignore-not-found
kubectl apply -f config/gen/authenticator-build.yaml