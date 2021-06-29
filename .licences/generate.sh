#!/bin/bash

command -v go-licence-detector >/dev/null 2>&1 || { echo >&2 "This script requires 'go-licence-detector' to be installed. Please install it via 'go get go.elastic.co/go-licence-detector'"; exit 1; }

SCRIPT_FULL_PATH=$(dirname "$0")
cd "$SCRIPT_FULL_PATH" || return

MODULES="api cli configuration-service distributor helm-service jmeter-service lighthouse-service mongodb-datastore remediation-service secret-service shipyard-controller statistics-service"

for MODULE in $MODULES; do
   echo "🔍 Analyzing dependencies in module $MODULE"
   ( cd ../"$MODULE" || return ; go list -m -json all | go-licence-detector -depsTemplate=../.licences/templates/dependencies.md.tmpl -depsOut=../.licences/"${MODULE}"-dependencies.md  -overrides=../.licences/overrides/overrides.json)
done

echo
echo "👍 done"