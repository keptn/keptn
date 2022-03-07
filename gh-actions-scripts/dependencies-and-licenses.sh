#!/bin/bash

TMP_DIR=$(mktemp -d 2>/dev/null || mktemp -d -t 'keptndeps')
MODULES=(
  api
  approval-service
  cli
  configuration-service
  distributor
  go-sdk
  helm-service
  jmeter-service
  lighthouse-service
  mongodb-datastore
  remediation-service
  secret-service
  shipyard-controller
  statistics-service
  webhook-service
  resource-service
)

for MODULE in "${MODULES[@]}"; do
   echo "🔍 Analyzing dependencies in module $MODULE"
   ( cd ./"$MODULE" || return ; go mod tidy > /dev/null 2>&1; go list -m -json all | go-licence-detector -depsTemplate=../.dependencies/templates/dependencies.csv.tmpl -depsOut="${TMP_DIR}"/"${MODULE}"-dependencies.txt  -overrides=../.dependencies/overrides/overrides.json )
done

echo "🔍 Analyzing dependencies in go-utils"
( cd ../go-utils || return ; go mod tidy > /dev/null 2>&1; go list -m -json all | go-licence-detector -depsTemplate=../keptn/.dependencies/templates/dependencies.csv.tmpl -depsOut="${TMP_DIR}"/go-utils-dependencies.txt )

echo "🔍 Analyzing dependencies in kubernetes-utils"
( cd ../kubernetes-utils || return ; go mod tidy > /dev/null 2>&1; go list -m -json all | go-licence-detector -depsTemplate=../keptn/.dependencies/templates/dependencies.csv.tmpl -depsOut="${TMP_DIR}"/kubernetes-utils-dependencies.txt )

cat "$TMP_DIR"/*.txt | sort | uniq > dependencies-and-licenses-go.txt

echo
echo "👍 done. written results to ./dependencies-and-licenses-go.txt"

cat dependencies-and-licenses-go.txt
