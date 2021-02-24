#!/bin/bash
# shellcheck disable=SC2181

# shellcheck disable=SC1091
source test/utils.sh

echo "y" | keptn uninstall
verify_test_step $? "keptn uninstall - failed"

# verify namespace keptn has been removed
kubectl -n keptn get namespace keptn 2> /dev/null

if [[ $? -eq 0 ]]; then
  echo "Found namespace keptn"
  echo "keptn uninstall - failed"
  exit 1
fi

# delete the namespaces for projects that we onboarded (if they exist)
echo "Deleting namespaces $PROJECT-dev $PROJECT-staging $PROJECT-production"
kubectl delete namespace $PROJECT-dev $PROJECT-staging $PROJECT-production || true

# wait for the actual deletion
sleep 60

echo "Keptn has been uninstalled. This is what is left on the cluster:"
kubectl get all --all-namespaces
