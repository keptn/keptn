#!/bin/bash

go get go.elastic.co/go-licence-detector

SCRIPT_FULL_PATH=$(dirname "$0")
cd $SCRIPT_FULL_PATH

MODULES="api cli configuration-service distributor helm-service jmeter-service lighthouse-service mongodb-datastore remediation-service secret-service shipyard-controller statistics-service"

for MODULE in $MODULES; do
   echo "🔍 Analyzing dependencies in module $MODULE"
   ( cd ../$MODULE; go list -m -json all | go-licence-detector -depsTemplate=../.licences/templates/dependencies.md.tmpl -depsOut=../.licences/${MODULE}-dependencies.md  -overrides=../.licences/overrides/overrides.json)
done

echo
echo "👍 done"