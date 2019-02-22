#!/bin/sh
REGISTRY_URI=$1

rm config/gen/control.yaml

cat config/control.yaml | \
  sed 's~REGISTRY_URI_PLACEHOLDER~'"$REGISTRY_URI"'~' >> config/gen/control.yaml 
  
kubectl apply -f config/gen/control.yaml