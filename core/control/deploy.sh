#!/bin/sh
CHANNEL_URI=$(kubectl describe channel keptn-channel -n keptn | grep "Hostname:" | sed 's~[ \t]*Hostname:[ \t]*~~')

rm -f config/gen/control.yaml

cat config/control.yaml | \
  sed 's~CHANNEL_URI_PLACEHOLDER~'"$CHANNEL_URI"'~' >> config/gen/control.yaml 
  
kubectl delete -f config/gen/control.yaml --ignore-not-found
kubectl apply -f config/gen/control.yaml