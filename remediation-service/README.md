# Remediation Service

The RemediationService is a Keptn service that handles `sh.keptn.event.get-action.triggered` events. First the "remediation.yaml"
resource uploaded for the relevant service is downloaded from the ConfigurationManagementService.
Based on the content of the received event, and the "remediation.yaml" file the next action will be determined and sent
out as payload of the `sh.keptn.event.get-action.finished` event.
