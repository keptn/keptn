#!/bin/bash

echo "Exposing Keptn Bridge..."
# expose Keptn's Bridge (for easier troubleshooting/debugging afterwards)
kubectl -n keptn create secret generic bridge-credentials --from-literal="BASIC_AUTH_USERNAME=${NIGHTLY_BRIDGE_USERNAME}" --from-literal="BASIC_AUTH_PASSWORD=${NIGHTLY_BRIDGE_PASSWORD}"
keptn configure bridge --action=expose
# restart bridge pod to make use of the secret
kubectl -n keptn delete pods --selector=run=bridge
sleep 5
# verify that bridge is available
export BRIDGE_URL=$(echo https://bridge.keptn.$(kubectl get cm -n keptn keptn-domain -ojsonpath={.data.app_domain}))
curl "${BRIDGE_URL}" -k

if [[ $? != '0' ]]; then
    echo "Accessing Keptn Bridge failed"
    exit 1
fi
