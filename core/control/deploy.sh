#!/bin/bash
CHANNEL_URI=$(kubectl describe channel keptn-channel -n keptn | grep "Hostname:" | sed 's~[ \t]*Hostname:[ \t]*~~')

if [[ -z "${CHANNEL_URI}" ]]; then
  echo "[keptn|ERROR] $(date +"[%m-%d %H:%M:%S]") CHANNEL_URI could not be derived from keptn-channel description."
  echo "[keptn|ERROR] $(date +"[%m-%d %H:%M:%S]") Stopping keptn installation. Already created resources are not deleted; execute the uninstallKeptn.sh script to clean-up.."
  exit 1
fi

rm -f config/gen/control.yaml

cat config/control.yaml | \
  sed 's~CHANNEL_URI_PLACEHOLDER~'"$CHANNEL_URI"'~' >> config/gen/control.yaml 
  
kubectl delete -f config/gen/control.yaml --ignore-not-found
kubectl apply -f config/gen/control.yaml
