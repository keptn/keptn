# Keptn Cloud Events

* [Create Project](#create-project)
* [Create Service](#create-service)
* [Configuration Change](#configuration-change)
* [Deployment Finished](#deployment-finished)
* [Tests Finished](#tests-finished)
* [Evaluation Done](#evaluation-done)
* [Problem](#problem)
* [Configure Monitoring](#configure-monitoring)

---

## Create Project

The *project create* event is sent when a new project should be created.

```json
{
  "type": "sh.keptn.internal.event.project.create",
  "specversion": "0.2",
  "source": "https://github.com/keptn/keptn/cli",
  "id": "49ac0dec-a83b-4bc1-9dc0-1f050c7e781b",
  "time": "2019-06-07T07:02:15.64489Z",
  "contenttype": "application/json",
  "shkeptncontext": "49ac0dec-a83b-4bc1-9dc0",
  "data": {
    "project": "sockshop",
    "gitUser": "scott",
    "gitToken": "token",
    "gitRemoteURL": "https://remote.url/project/repository",
    "shipyard": "stages:
      - name: \"staging\"
        deployment_strategy: \"blue_green_service\"
        test_strategy: \"performance\"
      - name: \"production\"
        deployment_strategy: \"blue_green_service\"
        remediation_strategy: \"automated\""
  }
}
```
([&uarr; up to index](#keptn-cloud-events))

## Create Service

The *service create* event is sent when a new service should be created.

```json
{
  "type": "sh.keptn.internal.event.service.create",
  "specversion": "0.2",
  "source": "https://github.com/keptn/keptn/cli",
  "id": "49ac0dec-a83b-4bc1-9dc0-1f050c7e781b",
  "time": "2019-06-07T07:02:15.64489Z",
  "contenttype": "application/json",
  "shkeptncontext": "49ac0dec-a83b-4bc1-9dc0",
  "data": {
    "project": "sockshop",
    "service": "carts",
    "helmChart": "string",
    "deploymentStrategies": "string"
  }
}
```
([&uarr; up to index](#keptn-cloud-events))

## Configuration Change

The *configuration change* event is sent when a desired state for a service is available and has to be updated in its configuration.

```json
{
  "type": "sh.keptn.event.configuration.change",
  "specversion": "0.2",
  "source": "https://github.com/keptn/keptn/cli",
  "id": "49ac0dec-a83b-4bc1-9dc0-1f050c7e781b",
  "time": "2019-06-07T07:02:15.64489Z",
  "contenttype": "application/json",
  "shkeptncontext":"49ac0dec-a83b-4bc1-9dc0-1f050c7e789b",
  "data": {
    "project": "sockshop",
    "stage": "staging",
    "service": "carts",
    "valuesCanary": {
      "image": "docker.io/keptnexamples/carts:0.9.1"
    }
  }
}
```
([&uarr; up to index](#keptn-cloud-events))

## Deployment Finished

The *deployment-finished* event is sent when a desired state of a service is deployed in a stage.

```json
{
  "type": "sh.keptn.events.deployment-finished",
  "specversion": "0.2",
  "source": "https://github.com/keptn/keptn/helm-service",
  "id": "49ac0dec-a83b-4bc1-9dc0-1f050c7e781b",
  "time": "2019-06-07T07:02:15.64489Z",
  "contenttype": "application/json",
  "shkeptncontext":"49ac0dec-a83b-4bc1-9dc0-1f050c7e789b",
  "data": {
    "project": "sockshop",
    "stage": "staging",
    "service": "carts",
    "testStrategy": "performance",
    "deploymentStrategy": "direct"
  }
}
```
([&uarr; up to index](#keptn-cloud-events))

## Tests Finished

The *tests-finished* event is sent when the tests for a service in a stage are finished.

```json
{
  "type": "sh.keptn.events.tests-finished",
  "specversion": "0.2",
  "source": "https://github.com/keptn/keptn/jmeter-service",
  "id": "49ac0dec-a83b-4bc1-9dc0-1f050c7e781b",
  "time": "2019-06-07T07:02:15.64489Z",
  "contenttype": "application/json",
  "shkeptncontext":"49ac0dec-a83b-4bc1-9dc0-1f050c7e789b",
  "data": {
    "project": "sockshop",
    "stage": "staging",
    "service": "carts",
    "testStrategy": "performance",
    "deploymentStrategy": "direct",
    "startedat": "2019-09-01 12:03"
  }
}
```
([&uarr; up to index](#keptn-cloud-events))

## Evaluation Done

The *evaluation-done* event is sent when the evaluation of the test execution is completed.

```json
{
  "type": "sh.keptn.events.evaluation-done",
  "specversion": "0.2",
  "source": "https://github.com/keptn/keptn/pitometer-service",
  "id": "49ac0dec-a83b-4bc1-9dc0-1f050c7e781b",
  "time": "2019-06-07T07:02:15.64489Z",
  "contenttype": "application/json",
  "shkeptncontext":"49ac0dec-a83b-4bc1-9dc0-1f050c7e789b",
  "data":{  
    "project":"sockshop",
    "service":"carts",
    "image":"keptnexamples/carts",
    "tag":"0.7.1",
    "stage":"dev",
    "teststrategy":"functional",
    "deploymentstrategy":"direct",
    "evaluationpassed":"false",
    "evaluationdetails": {
      "options": {
        "timeStart": 1558428643,
        "timeEnd": 1558429200
      },
      "totalScore": 50,
      "objectives": {
        "pass": 90,
        "warning": 75
      },
      "indicatorResults": [
        {
          "id": "ResponseTime_Service",
          "violations": [],
          "score": 50
        },
        {
          "id": "FailureRate_Service",
          "violations": [
            {
              "value": 12.55862428348098,
              "key": "SERVICE-34E5F883CB9DF269",
              "breach": "upperSevere",
              "threshold": 10
            }
          ],
          "score": 0
        }
      ],
      "result": "fail"
      }
   }
}
```
([&uarr; up to index](#keptn-cloud-events))

## Configure Monitoring

The *monitoring configure* event is sent when a monitoring solution need to be configured for a new service.

```json
{
  "type": "sh.keptn.event.monitoring.configure",
  "specversion": "0.2",
  "source": "https://github.com/keptn/keptn/jmeter-service",
  "id": "49ac0dec-a83b-4bc1-9dc0-1f050c7e781b",
  "time": "2019-06-07T07:02:15.64489Z",
  "contenttype": "application/json",
  "shkeptncontext":"49ac0dec-a83b-4bc1-9dc0-1f050c7e789b",
  "data": {
    "type": "prometheus",
    "project": "sockshop",
    "service": "carts",
    "serviceIndicators": "",
    "serviceObjectives": "",
    "remediation": ""
  }
}
```
([&uarr; up to index](#keptn-cloud-events))

## Problem

The *problem* event is sent when a monitored service causes any problem.

```json
{
  "type": "sh.keptn.events.problem",
  "specversion": "0.2",
  "source": "https://github.com/keptn/keptn/jmeter-service",
  "id": "49ac0dec-a83b-4bc1-9dc0-1f050c7e781b",
  "time": "2019-06-07T07:02:15.64489Z",
  "contenttype": "application/json",
  "shkeptncontext":"49ac0dec-a83b-4bc1-9dc0-1f050c7e789b",
  "data": {
    "state": "open",
    "problemID": "ad7a-139fyf-915da",
    "problemtitle": "problem title",
    "problemdetails": "problem details",
    "impactedEntity": "impacted entity",
  }
}
```
([&uarr; up to index](#keptn-cloud-events))