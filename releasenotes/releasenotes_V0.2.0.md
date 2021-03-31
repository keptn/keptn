# Release Notes 0.2.0

## Release Goal

The goal of this release is to provide a keptn installation that comes with automated quality gates evaluation, different deployment strategies, a dedicated keptn CLI and integration with 3rd party vendors to allow you to build your cloud-native delivery.

### Automated setup
This release comes with scripts that automatically install keptn on your GKE cluster and set up all needed components including Istio, Knative, Kibana and all keptn services. Additionally, keptn creates namespaces for the different stages your service will go through (e.g., dev, staging, production) to ensure resource isolation.

### Easy onboarding of custom services
With the keptn CLI custom services can be onboarded so keptn takes care of delivery. For onboarding, a so-called `shipyard` file is needed that defines which stages have to be created and also defines the quality gates that have to be passed for artifacts to be promoted from one stage to the other. This release comes with a demo service that can be onboarded and defines the following stages in its shipyard file: dev, staging, production, as well as a deployment and testing strategy for each stage.

### Deployments with automated quality gate evaluation

- *Quality Gates*: Once services have been onboarded, keptn takes care that only built artifacts that pass the quality gates will be promoted to next stage. The quality gates are automatically validated by our Pitometer component, which already provides integrations for Prometheus and Dynatrace monitoring but also allows developers to integrate their own data sources. The actual validation in Pitometer is done by leveraging a so-called `perfspec` file that defines checks for given metrics, e.g., the response time of a service has to be lower than the defined threshold. 

- *Deployment Strategies*: In addition to quality gates, keptn provides different deployment strategies for onboarded services. For example, a direct deployment directly replaces the old version of the artifact with its new version, while a blue/green deployment strategy deploys the artifact in a separate version (e.g., blue) while retaining the previous version of the artifact (e.g., green). This allows performance tests and the quality gate evaluation on the blue version while the green version might still handle user traffic. Only if the blue version passes the evaluation it will replace the green version by rerouting all traffic to the blue version. This also allows an instant fallback to the previous version if some (production) issues are detected.

### Runbook automation and self-healing

While keptn provides means to ensure high quality of all onboarded services as they get promoted through the different stages, an issue free production environment cannot be guaranteed. For example, changes at runtime might cause problems that have to be handled quickly and reliably.

In this release, keptn provides an integration with ServiceNow and Dynatrace to trigger workflows in ServiceNow if Dynatrace detects problems in your environment. In the demo that is shipped with this release, the service is prepared to allow for configuration changes at runtime, which will introduce an increase of the failure rate. This increase is detected by Dynatrace and a problem notification is sent to keptn. In keptn, the ServiceNow service will create an incident in the provided ServiceNow instance which will trigger the workflow to revert the configuration change in the production environment. 

### Unbreakable delivery pipelines

keptn provides an unbreakable delivery pipeline, which prevents that bad code changes are impacting your end users by ensuring that only high quality code can be promoted to subsequent stages and eventually into production.

Thereby it relies on three concepts:

- Shift-Left is the ability to pull data for specific entities (processes, services, or applications) through an automation API and feed it into the tools that are used to decide on whether to stop the pipeline or keep it running,

- Shift-Right is the ability to push deployment information and metadata to your monitoring solution (e.g., to differentiate BLUE vs GREEN deployments), to push build or revision numbers of a deployment, or to notify about configuration changes,

- Self-Healing is the ability for smart auto-remediation that addresses the root cause of a problem and not the symptom.

## Version capabilities:

keptn is installed by using these images from the [keptn Docker Hub registry](https://hub.docker.com/u/keptn):

- keptn/keptn-authenticator:0.2.0
- keptn/keptn-control:0.2.0
- keptn/keptn-event-broker:0.2.0
- keptn/keptn-event-broker-ext:0.2.0
- keptn/pitometer-service:0.1.0
- keptn/servicenow-service:0.1.0
- keptn/github-service:0.1.0
- keptn/jenkins-service:0.1.0
  - keptn/jenkins-0.4.0

## Known limitations:

- installation currently only on GKE (more platforms to come)
- no multi-tenant functionality yet (only one GitHub organization can be configured with the keptn server)
- for use cases that require Dynatrace: support for Dynatrace SaaS tenants only (managed support to come)
- keptn CLI output not reliably reflecting success/error of keptn services: the CLI only reflects the successful acknowledgment of the CLI command but not its successful execution
