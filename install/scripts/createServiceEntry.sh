#!/bin/bash

entries=$(curl https://$1.live.dynatrace.com/api/v1/deployment/installer/agent/connectioninfo?Api-Token=$2 | jq -r '.communicationEndpoints[]')

rm -f ../manifests/gen/service_entries_oneagent.yml
rm -f ../manifests/gen/hosts
rm -f ../manifests/gen/service_entries

cat ../manifests/istio/service_entries_tpl/header_tmpl >> ../manifests/gen/service_entries_oneagent.yml

echo -e "  - $1.live.dynatrace.com" >> ../manifests/gen/hosts
cat ../manifests/istio/service_entries_tpl/service_entry_tmpl | sed 's~ENDPOINT_PLACEHOLDER~'"$1"'.live.dynatrace.com~' >> ../manifests/gen/service_entries

for row in $entries; do
    row=$(echo $row | sed 's~https://~~')
    row=$(echo $row | sed 's~/communication~~')
    echo -e "  - $row" >> ../manifests/gen/hosts
    cat ../manifests/istio/service_entries_tpl/service_entry_tmpl | sed 's~ENDPOINT_PLACEHOLDER~'"$row"'~' >> ../manifests/gen/service_entries
done

# Add virtual service section
cat ../manifests/gen/hosts >> ../manifests/gen/service_entries_oneagent.yml
cat ../manifests/istio/service_entries_tpl/virtual_service_tmpl >> ../manifests/gen/service_entries_oneagent.yml

# Add tls list
cat ../manifests/gen/hosts >> ../manifests/gen/service_entries_oneagent.yml
cat ../manifests/istio/service_entries_tpl/tls_tmpl >> ../manifests/gen/service_entries_oneagent.yml

# Attach list of service entries at the end
cat ../manifests/gen/service_entries >> ../manifests/gen/service_entries_oneagent.yml

# Apply service entries
kubectl apply -f ../manifests/gen/service_entries_oneagent.yml
