#!/bin/bash

helm template ../dist/keptn-installer/keptn-*.tgz --dry-run > "$HELM_TEMPLATE"
yq eval-all '[
  .spec.template.spec.containers.[] |
  {
    "name": .name,
    "resources": .resources
  }
]' "$HELM_TEMPLATE" -o=json | \
jq '[
  .[] | {
    name: .name,
    cpu_request: .resources.requests.cpu,
    cpu_limit: .resources.limits.cpu,
    mem_request: .resources.requests.memory,
    mem_limit: .resources.limits.memory
  }
] |
unique_by(.name)' > "$RESOURCE_JSON"
npx tablemark-cli@v2.0.0 "$RESOURCE_JSON" -c "Name" -c "CPU Request" -c "CPU Limit" -c "RAM Request" -c "RAM Limit" > "$RESOURCE_MARKDOWN"

{
  echo ""
  echo "### Resource Stats"
} >> "$RELEASE_NOTES_FILE"
cat "$RESOURCE_MARKDOWN" >> "$RELEASE_NOTES_FILE"
