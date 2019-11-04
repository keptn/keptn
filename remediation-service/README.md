# Remediation Service

The remediation-service is a keptn core component. It is responsible for remediation of onboarded services in response to issues detected during their runtime, e.g., by Prometheus. 

The service receives as an input a problem event. Upon this event the services tries to find a matching remediation action from the `remediation.yaml` file that has been onboarded for the affected service.
A corresponding configuration change will be created by the remediation service which will be applied by the keptn workflow to remediate the issue.
