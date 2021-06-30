#!/bin/bash

##########################################################################################
# Analyzes each module of keptn/keptn mono repo and spits out a csv formatted file called
# "all-dependencies.txt" containing a list of dependencies + their version and licences
##########################################################################################

command -v go-licence-detector >/dev/null 2>&1 || { echo >&2 "This script requires 'go-licence-detector' to be installed. Please install it via 'go get go.elastic.co/go-licence-detector'"; exit 1; }

SCRIPT_FULL_PATH=$(dirname "$0")
cd "$SCRIPT_FULL_PATH" || return

TMP_DIR=$(mktemp -d 2>/dev/null || mktemp -d -t 'keptndeps')
echo "storing dependencies files in ${TMP_DIR}"

MODULES="api cli configuration-service distributor helm-service jmeter-service lighthouse-service mongodb-datastore remediation-service secret-service shipyard-controller statistics-service"
for MODULE in $MODULES; do
   echo "🔍 Analyzing dependencies in module $MODULE"
   ( cd ../"$MODULE" || return ; go list -m -json all | go-licence-detector -depsTemplate=../.licences/templates/dependencies.csv.tmpl -depsOut="${TMP_DIR}"/"${MODULE}"-dependencies.txt  -overrides=../.licences/overrides/overrides.json)
done

cat "$TMP_DIR"/*.txt | sort | uniq > all-dependencies.txt

echo
echo "👍 done. written results to ./all-dependencies.txt"