#!/bin/sh
REGISTRY_URI=$1

rm -f config/gen/authenticator.yaml

cat config/authenticator.yaml | \
  sed 's~REGISTRY_URI_PLACEHOLDER~'"$REGISTRY_URI"'~' >> config/gen/authenticator.yaml 
  
kubectl apply -f config/gen/authenticator.yaml