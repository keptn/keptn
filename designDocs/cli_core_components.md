# Design decisions for core components and CLI

This document describes design decision regarding the keptn CLI.

## Install keptn

1. Execute `defineCredentials.sh` that defines:
    * DT_API_TOKEN
    * DT_PAAS_TOKEN
    * (?) GITHUB_TOKEN
    * (?) GITHUB_ORG

1. Execute `setupInfrastructure.sh` that returns:
    * Endpoint
    * KEPTN-API-TOKEN

## Core components

The following components are available after the keptn install:

### Event Broker

This component accepts incoming events from sources such as GitHub webhooks, and pushes them into an internal knative eventing queue. This queue will fan out incoming messages to subscribers.

### Control

The control component provides following endpoints:
* /auth
* /onboard
* ...

The payload of each request needs to be a CloudEvent. Each request needs to contain a request header attribute *x-keptn-signature* that holds a signature of the payload and the *KEPTN-API-TOKEN*: `sha1(payload || keptn-api-token)`

### Auth

The auth component manages the keptn-api-token and is used by the control component to authorize an incoming request.

## Use keptn via CLI

In the following, planned commands for the CLI are listed and explained:

* **keptn auth --endpoint --api-token**: Authenticates the keptn CLI against a keptn installation.

    *Example:* 
    ```console
    $ keptn auth --endpoint= --api-token=
    ```

* **keptn configure --org **: Configures the **GitHub organization** in the keptn installation.

    *Example:*
    ```console
    $ keptn configure --org= 
    ```

* **keptn create project**: Creates a new repository in the GitHub organization and initializes the repository with helm charts. For this now, the shipyard.yml file (see below) contains the number of stages and name of each stage. Based on that information, this command creates a branch for each stage.

    *Example:*
    ```console
    $ keptn create project sockshop shipyard.yml
    ```

* **keptn onboard service**: Takes an service.yaml and sends a CloudEvent (containing the yaml) to the control component.
   
     *Example: Takes file*
    ```console
    $ keptn onboard service --project=sockshop --file=carts.yaml
    ```

# Structure GitHub Organization

Github Organization "Dynatrace"<br/>
|<br/>
|- Repository: Pipline CI360<br/>
|<br/>
|- Repository: Pipline Licensing <br/>
|<br/>
|- Repository: Pipline Sockshop {master, dev, staging, production} <br/>
&nbsp;&nbsp;&nbsp; |- charts/<br/>
&nbsp;&nbsp;&nbsp; |- config/<br/>
&nbsp;&nbsp;&nbsp; |- values.yaml<br/>

# Shipyard

*Template:*
```yaml
stages: 
- name: 
    development_strategy: [direct, service_blue/green, application_blue/green]
    deployment_operator:
    test_strategy: [functional, health, continous_performance, production]
    test_operator: 
    validation_operator: 
    remediation_handler: [rollback]
    monitoring_handler: [dynatrace, prometheus]
    next: staging
```

*Example:*
```yaml
stages: 
- name: dev
    deployment_strategy: direct
    deployment_operator: jenkins-operator, slack
    test_strategy: functional
    test_operator: neotys_operator
    validation_operator: keptn.monspec-evaluator
    remediation_handler: // TBD    
    next: staging
- name: staging
    deployment_strategy: service_blue/green
    deployment_operator: jenkins-operator, slack
    test_strategy: continous_performance
    test_operator: neotys_operator
    validation_operator: keptn.monspec-evaluator
    remediation_handler: rollback
    next: production
- name: production
    deployment_strategy: application blue/green
    deployment_operator: jenkins-operator, slack
    test_strategy: production
    test_operator: neotys_operator
    validation_strategy: production
    remediation_handler: rollback
```
