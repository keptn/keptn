#!/bin/bash

# Print complete yaml of the keptn installer
helm template ../dist/keptn-installer/keptn-*.tgz --dry-run --name-template=keptn > "$HELM_TEMPLATE"

# yq: Extract container name and resource limits/requests for all resources
# jq: Bring yq result into nice format that can be read by the tablemark-cli
yq eval-all '[
{
  "name": .metadata.name,
  "container": .spec.template.spec.containers.[]
} |
{
  "name": .name,
  "container": .container.name,
  "resources": .container.resources
}
]' "$HELM_TEMPLATE" -o=json | \
jq '[
.[] | {
    name: .name,
    container_name: .container,
    cpu_request: .resources.requests.cpu,
    cpu_limit: .resources.limits.cpu,
    mem_request: .resources.requests.memory,
    mem_limit: .resources.limits.memory
  }
]' > "$RESOURCE_JSON"

# Generate markdown table from JSON
npx tablemark-cli@v2.0.0 "$RESOURCE_JSON" -c "Name" -c "Container Name" -c "Image" -c "CPU Request" -c "CPU Limit" -c "RAM Request" -c "RAM Limit" > "$RESOURCE_MARKDOWN"

# Attach resource stats to release notes
{
  echo ""
  echo "<details>"
  echo "<summary>Kubernetes Resource Data</summary>"
  echo "### Resource Stats"
  cat "$RESOURCE_MARKDOWN"
  echo "</details>"
} >> "$RELEASE_NOTES_FILE"
