#!/bin/bash

NAMESPACE=${1:-"keptn"}

# Gathers resource limits by deployment for keptn namespace
# | Deployment                          	| Memory (requested) 	| CPU (requested) 	| Memory (limit) 	| CPU (limit)
KUBE_GET_DEPS="kubectl get deployments -n $NAMESPACE"
echo -e ""
echo -e "| Pod | Container | Memory (requested) | CPU (requested) | Memory (limit) | CPU (limit) | Images |"
echo -e "|-----|-----------|--------------------|-----------------|----------------|-------------|--------|"
$KUBE_GET_DEPS | sed '1d' | awk '{print $1}' | sort | while read -r DEPLOYMENT; do
  # shellcheck disable=SC2086
  $KUBE_GET_DEPS "$DEPLOYMENT" -o jsonpath='{range .spec.template.spec.containers[*]}{"| "}{"'$DEPLOYMENT'"}{" | "}{.name}{" | "}--{.resources.requests.memory}{" | "}--{.resources.requests.cpu}{" | "}--{.resources.limits.memory}{" | "}--{.resources.limits.cpu}{" | "}{.image}{" | "}{"\n"}{end}' | sed -E -e 's/--([0-9])/\1/g' -e 's/--/-/g'
done

echo -e ""
echo -e "Summary (whole cluster):"
echo -e "\`\`\`"
echo -e '$ kubectl describe node | grep -A5 "Allocated"'
kubectl describe node | grep -A5 "Allocated"
echo -e "\`\`\`"
echo -e "Please note: Depending on the setup, the above includes usage for Istio as well as the Kubernetes control-plane"
echo -e ""

# print PVC data
echo -e "| Name | Size |"
echo -e "|------|------|"
kubectl get pvc -n "$NAMESPACE" | sed '1d' | awk '{print $1}' | sort | while read -r PVC; do
  # shellcheck disable=SC2086
  kubectl get pvc "$PVC" -n "$NAMESPACE" -o jsonpath='{range .spec}{"| "}{"'$PVC'"}{" | "}--{.resources.requests.storage}{" | "}{"\n"}{end}' | sed -E -e 's/--([0-9])/\1/g' -e 's/--/-/g'
done

echo -e ""
