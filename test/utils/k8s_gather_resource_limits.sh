#!/bin/bash

NAMESPACE=${1:-"keptn"}

# Gathers resource limits by deployment for keptn namespace

echo -e "| Pod | Container | lim.cpu | lim.mem | req.cpu | req.mem | Images |"
echo -e "|-----|-----------|---------|---------|---------|---------|--------|"
kubectl get deployments -n $NAMESPACE | sed '1d' | awk '{print $1}' | sort | while read DEPLOYMENT; do
  kubectl get deployment $DEPLOYMENT -n $NAMESPACE -o jsonpath='{range .spec.template.spec.containers[*]}{"'$DEPLOYMENT'"}{" | "}{.name}{" | "}--{.resources.limits.cpu}{" | "}--{.resources.limits.memory}{" | "}--{.resources.requests.cpu}{" | "}--{.resources.requests.memory}{" | "}{.image}{" | "}{"\n"}{end}' | sed -E -e 's/--([0-9])/\1/g' -e 's/--/-/g'
done
