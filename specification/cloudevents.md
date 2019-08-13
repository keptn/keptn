# CloudEvents used by keptn

Internal Events
* [Create Project](#create-project)
* [Create Project (from CLI)](#create-project-from-cli)
* [Create Service](#create-service)
* [Create Project (from CLI)](#create-service-from-cli)
* [Generic done event](#generic-done)

Keptn Events
* [Start Change Configuration](#start-change-configuration)
* [Start Apply Configuration](#start-apply-configuration)
* [Start Tests](#start-tests)
* [Start Evaluation](#start-evaluation)
* [Start Deployment](#start-deployment)
* [Start Release](#start-release)
* [Start Operations](#start-operations)
* [Done](#done)

---

# Internal Events

## Create Project

The *create project* event is sent when a new project should be created.
```json
{
  "type": "sh.keptn.internal.events.project.create",
  "specversion": "0.2",
  "source": "https://github.com/keptn/keptn/cli",
  "id": "49ac0dec-a83b-4bc1-9dc0-1f050c7e781b",
  "time": "2019-06-07T07:02:15.64489Z",
  "contenttype": "application/json",
  "shkeptncontext":"adf-afefea-feafa",
  "shkeptnphaseid":"k5dg-565h-9o87",
  "shkeptnphase":"INIT",
  "shkeptnstepid":"a3ce-f013-4h78",
  "shkeptnstep":"CREATE_PROJECT",
  "data": {
    "project": "sockshop",
    "gituser": "hermann",
    "gittoken": "adfadsf",
    "gitremoteurl": "https://remote.url/",
    "stages": [ 
          {
            "name": "dev",
            "deployment_strategy": "direct",
            "test_strategy": "functional"
          },
          {
            "name": "production",
            "deployment_strategy": "blue_green_service"
          }
      ]
  }
}
```
([&uarr; up to index](#cloudevents-used-by-keptn))

## Create Project from CLI

The *create project* sent from the CLI when a new project should be created.
```json
{
  "type": "sh.keptn.events.project.create",
  "specversion": "0.2",
  "source": "https://github.com/keptn/keptn/cli",
  "id": "49ac0dec-a83b-4bc1-9dc0-1f050c7e781b",
  "time": "2019-06-07T07:02:15.64489Z",
  "contenttype": "application/json",
  "data": {
    "project": "sockshop",
    "gituser": "hermann",
    "gittoken": "adfadsf",
    "gitremoteurl": "https://remote.url/",
    "stages": [ 
          {
            "name": "dev",
            "deployment_strategy": "direct",
            "test_strategy": "functional"
          },
          {
            "name": "production",
            "deployment_strategy": "blue_green_service"
          }
      ]
  }
}
```
([&uarr; up to index](#cloudevents-used-by-keptn))


## Create Service

The *create service* event is sent when a new service should be created.

```json
{
  "type": "sh.keptn.internal.events.service.create",
  "specversion": "0.2",
  "source": "https://github.com/keptn/keptn/cli",
  "id": "49ac0dec-a83b-4bc1-9dc0-1f050c7e781b",
  "time": "2019-06-07T07:02:15.64489Z",
  "contenttype": "application/json",
  "shkeptncontext":"adf-afefea-feafa",
  "shkeptnphaseid":"k5dg-565h-9o87",
  "shkeptnphase":"INIT",
  "shkeptnstepid":"a3ce-f013-4h78",
  "shkeptnstep":"CREATE_SERVICE",
  "data": {
    "project": "sockshop",
    "servicename": "carts",
    "values": "string",
    "manifest": "string",
    "templates": {
      "service": "string",
      "deployment": "string"
    }
  }
}
```
([&uarr; up to index](#cloudevents-used-by-keptn))


## Create Service from CLI

The *create project* event sent from the CLI when a new service should be created.

```json
{
  "type": "sh.keptn.events.service.create",
  "specversion": "0.2",
  "source": "https://github.com/keptn/keptn/cli",
  "id": "49ac0dec-a83b-4bc1-9dc0-1f050c7e781b",
  "time": "2019-06-07T07:02:15.64489Z",
  "contenttype": "application/json",
  "data": {
    "project": "sockshop",
    "servicename": "carts",
    "values": "string",
    "manifest": "string",
    "templates": {
      "service": "string",
      "deployment": "string"
    }
  }
}
```
([&uarr; up to index](#cloudevents-used-by-keptn))


## Generic done

```json
{
  "type": "sh.keptn.events.done",
  "specversion": "0.2",
  "source": "https://github.com/keptn/keptn/shipyard-service",
  "id": "49ac0dec-a83b-4bc1-9dc0-1f050c7e781b",
  "time": "2019-06-07T07:02:15.64489Z",
  "contenttype": "application/json",
  "shkeptncontext":"49ac0dec-a83b-4bc1-9dc0-1f050c7e789b",
  "shkeptnphaseid":"k5dg-565h-9o87",
  "shkeptnphase":"INIT",
  "shkeptnstepid":"a3ce-grdv-qwed",
  "shkeptnstep":"CREATE_PROJECT",
  "data": {
    "result": "success|error",
    "message": "successmessage|errormessage",
    "version": "git-commit-sha"
  }
}
```
([&uarr; up to index](#cloudevents-used-by-keptn))


# Keptn Events

## Start Change Configuration 

The *start change configuration* event is sent when a desired state for a service is available and has to be updated in its configuration.

```json
{
  "type": "sh.keptn.events.configuration.change",
  "specversion": "0.2",
  "source": "https://github.com/keptn/workflow-engine",
  "id": "49ac0dec-a83b-4bc1-9dc0-1f050c7e781b",
  "time": "2019-06-07T07:02:15.64489Z",
  "contenttype": "application/json",
  "shkeptncontext":"49ac0dec-a83b-4bc1-9dc0-1f050c7e789b",
  "shkeptnphaseid":"k5dg-565h-9o87",
  "shkeptnphase":"DEPLOYMENT",
  "shkeptnstepid":"a3ce-f013-4h78",
  "shkeptnstep":"CHANGE_CONFIGURATION",
  "data": {
    "project": "sockshop",
    "service": "carts",
    "image": "keptnexamples/carts",
    "tag": "0.7.1"
  }
}
```
([&uarr; up to index](#cloudevents-used-by-keptn))

## Start Apply Configuration

The *start apply configuration* event is sent when the service configuration is updated and has to be applied to a Kubernetes namespace.

```json
{
  "type": "sh.keptn.events.configuration.apply",
  "specversion": "0.2",
  "source": "https://github.com/keptn/workflow-engine",
  "id": "49ac0dec-a83b-4bc1-9dc0-1f050c7e781b",
  "time": "2019-06-07T07:02:15.64489Z",
  "contenttype": "application/json",
  "shkeptncontext":"49ac0dec-a83b-4bc1-9dc0-1f050c7e789b",
  "shkeptnphaseid":"k5dg-565h-9o87",
  "shkeptnphase":"DEPLOYMENT",
  "shkeptnstepid":"a3ce-grdv-3241",
  "shkeptnstep":"APPLY_CONFIGURATION",
  "data": {
    "project": "sockshop",
    "service": "carts",
    "image": "keptnexamples/carts",
    "tag": "0.7.1",
    "stage":"dev"
  }
}
```
([&uarr; up to index](#cloudevents-used-by-keptn))

## Start Tests

The *start test* event is sent when a serivce configuration is applied to a Kuberentes namespace and has to be tested.

```json
{
  "type": "sh.keptn.events.tests.start",
  "specversion": "0.2",
  "source": "https://github.com/keptn/workflow-engine",
  "id": "49ac0dec-a83b-4bc1-9dc0-1f050c7e781b",
  "time": "2019-06-07T07:02:15.64489Z",
  "contenttype": "application/json",
  "shkeptncontext":"49ac0dec-a83b-4bc1-9dc0-1f050c7e789b",
  "shkeptnphaseid":"k5dg-565h-9o87",
  "shkeptnphase":"DEPLOYMENT",
  "shkeptnstepid":"a3ce-grdv-3241",
  "shkeptnstep":"TESTS",
  "data": {
    "project": "sockshop",
    "service": "carts",
    "image": "keptnexamples/carts",
    "tag": "0.7.1",
    "stage":"dev",
    "teststrategy":"functional"
  }
}
```
([&uarr; up to index](#cloudevents-used-by-keptn))

## Start Evaluation

The *start evaluation* event is sent when a test for a service is completed and has to be evaluated.

```json
{
  "type": "sh.keptn.events.evaluation.start",
  "specversion": "0.2",
  "source": "https://github.com/keptn/workflow-engine",
  "id": "49ac0dec-a83b-4bc1-9dc0-1f050c7e781b",
  "time": "2019-06-07T07:02:15.64489Z",
  "contenttype": "application/json",
  "shkeptncontext":"49ac0dec-a83b-4bc1-9dc0-1f050c7e789b",
  "shkeptnphaseid":"k5dg-565h-9o87",
  "shkeptnphase":"DEPLOYMENT",
  "shkeptnstepid":"a3ce-grdv-qwed",
  "shkeptnstep":"TESTS",
  "data": {
    "project": "sockshop",
    "service": "carts",
    "image": "keptnexamples/carts",
    "tag": "0.7.1",
    "stage":"dev",
    "teststrategy":"functional"
  }
}
```
([&uarr; up to index](#cloudevents-used-by-keptn))

## Start Deployment

The *start deployment* event is sent when a new version of an artifact, i.e., a new version of a container image, is available and should be deployed. 

```json
{
  "type": "sh.keptn.events.deployment.start",
  "specversion": "0.2",
  "source": "https://github.com/keptn/workflow-engine",
  "id": "49ac0dec-a83b-4bc1-9dc0-1f050c7e781b",
  "time": "2019-06-07T07:02:15.64489Z",
  "contenttype": "application/json",
  "shkeptncontext":"49ac0dec-a83b-4bc1-9dc0-1f050c7e789b",
  "shkeptnphaseid":"k5dg-565h-9o87",
  "shkeptnphase":"DEPLOYMENT",
  "shkeptnstepid":"TODO:TBD",
  "shkeptnstep":"TODO:TBD",
  "data": {
    "project": "sockshop",
    "service": "carts",
    "image": "keptnexamples/carts",
    "tag": "0.7.1"
  }
}
```
([&uarr; up to index](#cloudevents-used-by-keptn))


## Start Release

The *start release* event is sent when a new deployment of an artifact is available and should be released to end-user.

```json
{
  "type": "sh.keptn.events.release.start",
  "specversion": "0.2",
  "source": "https://github.com/keptn/workflow-engine",
  "id": "49ac0dec-a83b-4bc1-9dc0-1f050c7e781b",
  "time": "2019-06-07T07:02:15.64489Z",
  "contenttype": "application/json",
  "shkeptncontext":"49ac0dec-a83b-4bc1-9dc0-1f050c7e789b",
  "shkeptnphaseid":"k5dg-565h-9o87",
  "shkeptnphase":"RELEASE",
  "shkeptnstepid":"TODO:TBD",
  "shkeptnstep":"TODO:TBD",
  "data": {
    "project": "sockshop",
    "service": "carts",
    "stage": "dev"
  }
}
```
([&uarr; up to index](#cloudevents-used-by-keptn))

## Start Operations (TBD!)

The *start operations* event is sent when a problem related to a keptn-managed service appears.

```json
{
  "type": "sh.keptn.events.start-operations",
  "specversion": "0.2",
  "source": "https://github.com/keptn/workflow-engine",
  "id": "49ac0dec-a83b-4bc1-9dc0-1f050c7e781b",
  "time": "2019-06-07T07:02:15.64489Z",
  "contenttype": "application/json",
  "shkeptncontext":"49ac0dec-a83b-4bc1-9dc0-1f050c7e789b",
  "shkeptnphaseid":"k5dg-565h-9o87",
  "shkeptnphase":"OPERATIONS",
  "shkeptnstepid":"TODO:TBD",
  "shkeptnstep":"TODO:TBD",
  "data": {
    "service": "carts",
    "stage": "dev"
  }
}
```
([&uarr; up to index](#cloudevents-used-by-keptn))


## Done

The *done* event is sent by each keptn service as response to the execution of a *.start or *.create event. 

```json
{
  "type": "sh.keptn.events.done",
  "specversion": "0.2",
  "source": "https://github.com/keptn/workflow-engine",
  "id": "49ac0dec-a83b-4bc1-9dc0-1f050c7e781b",
  "time": "2019-06-07T07:02:15.64489Z",
  "contenttype": "application/json",
  "shkeptncontext":"49ac0dec-a83b-4bc1-9dc0-1f050c7e789b",
  "shkeptnphaseid":"k5dg-565h-9o87",
  "shkeptnphase":"DEPLOYMENT",
  "shkeptnstepid":"a3ce-grdv-qwed",
  "shkeptnstep":"TESTS",
  "data": {
    "project": "sockshop",
    "service": "carts",
    "image": "keptnexamples/carts",
    "tag": "0.7.1",
    "stage":"dev",
    "teststrategy":"functional"
  }
}
```
([&uarr; up to index](#cloudevents-used-by-keptn))


# Postman collection of keptn CloudEvents
For testing purposes a postman collection is provided that contains all keptn CloudEvents. For now, the postman collection doesn't contain the *start-\** events.

https://github.com/keptn/keptn-specification/blob/master/keptn-events.postman_collection.json

---
---

## New Artifact (deprecated)

The *new artifact* event is triggered when a new version of an artifact i.e. a new version of a container image is available. 

```json
{
  "type": "sh.keptn.events.new-artifact",
  "specversion": "0.2",
  "source": "https://github.com/keptn/keptn/cli#new-artifact",
  "id": "49ac0dec-a83b-4bc1-9dc0-1f050c7e780b",
  "time": "2019-06-07T07:00:13.64489Z",
  "contenttype": "application/json",
  "data": {
    "project": "sockshop",
    "service": "carts",
    "image": "keptnexamples/carts",
    "tag": "0.7.1"
  }
}
```

## Start Configuration Change (deprecated)

Upon receiving a *new artifact* event a *start configuration change* event is dispatched that holds information about the updated artifact. A service takes care of updating the configuration in the version control system of choice.

```json
{
  "type": "sh.keptn.events.start-configuration-change",
  "specversion": "0.2",
  "source": "https://github.com/keptn/keptn-service#start-configuration-change",
  "id": "49ac0dec-a83b-4bc1-9dc0-1f050c7e781b",
  "time": "2019-06-07T07:02:15.64489Z",
  "contenttype": "application/json",
  "shkeptncontext":"49ac0dec-a83b-4bc1-9dc0-1f050c7e789b",
  "data": {
    "project": "sockshop",
    "service": "carts",
    "image": "keptnexamples/carts",
    "tag": "0.7.1"
  }
}
```

## Configuration Change Done (deprecated)

The *configuration change done* event is triggered when the GitOps configuration of a stage has been changed. 

```json
{  
   "type":"sh.keptn.events.configuration-change-done",
   "specversion":"0.2",
   "source": "https://github.com/keptn/github-service",
   "id": "49ac0dec-a83b-4bc1-9dc0-1f050c7e782b",
   "time":"2019-03-25T15:21:27+00:00",
   "contenttype":"application/json",
   "shkeptncontext":"db51be80-4fee-41af-bb53-1b093d2b694c",
   "data":{  
      "githuborg":"keptn-tiger",
      "project":"sockshop",
      "service":"carts",
      "image":"keptnexamples/carts",
      "tag":"0.7.1",
      "stage":"dev"
   },
}
```

## Start Deployment (deprecated)

Upon receiving a *configuration change done* event a *start deployment* event is dispatched that indicates that the changed configuration can now be applied.

```json
{  
   "type":"sh.keptn.events.start-deployment",
   "specversion":"0.2",
   "source": "https://github.com/keptn/keptn-service#start-deployment",
   "id": "49ac0dec-a83b-4bc1-9dc0-1f050c7e783b",
   "time":"2019-03-25T15:21:27+00:00",
   "contenttype":"application/json",
   "shkeptncontext":"db51be80-4fee-41af-bb53-1b093d2b694c",
   "data":{  
      "githuborg":"keptn-tiger",
      "project":"sockshop",
      "service":"carts",
      "image":"keptnexamples/carts",
      "tag":"0.7.1",
      "stage":"dev",
      "teststrategy":"functional",
      "deploymentstrategy":"direct"
   },
}
```

## Deployment Done (deprecated)

The *deployment done* event is sent by a deployment provider when it has finished applying a configuration change.

```json
{  
   "type":"sh.keptn.events.deployment-done",
   "specversion":"0.2",
   "source":"https://github.com/keptn/helm-service",
   "id":"49ac0dec-a83b-4bc1-9dc0-1f050c7e789b",
   "time":"20190325-15:22:50.560",
   "contenttype":"application/json",
   "shkeptncontext":"db51be80-4fee-41af-bb53-1b093d2b694c",
   "data":{  
      "githuborg":"keptn-tiger",
      "project":"sockshop",
      "service":"carts",
      "image":"keptnexamples/carts",
      "tag":"0.7.1",
      "stage":"dev",
      "teststrategy":"functional",
      "deploymentstrategy":"direct"
   }
}
```

## Start Tests (deprecated)

Upon receiving a *deployment done* event a *start tests* event is dispatched that indicates that the deployed artifact can now be tested.

```json
{  
   "type":"sh.keptn.events.start-tests",
   "specversion":"0.2",
   "source":"https://github.com/keptn/keptn-service#start-tests",
   "id":"49ac0dec-a83b-4bc1-9dc0-1f050c7e789b",
   "time":"20190325-15:22:50.560",
   "contenttype":"application/json",
   "shkeptncontext":"db51be80-4fee-41af-bb53-1b093d2b694c",
   "data":{  
      "githuborg":"keptn-tiger",
      "project":"sockshop",
      "service":"carts",
      "image":"keptnexamples/carts",
      "tag":"0.7.1",
      "stage":"dev",
      "teststrategy":"functional",
      "deploymentstrategy":"direct"
   }
}
```

## Tests Done (deprecated)

The *tests done* event is triggered when a test provider is finished executing tests. If a test provider can evaluate the result of the test on its own, e.g. in case of a functional test, the test provider can dispatch an *evaluation done* event with details about the test results.

```json
{  
   "type":"sh.keptn.events.tests-done",
   "specversion":"0.2",
   "source":"https://github.com/keptn/jmeter-service",
   "id":"49ac0dec-a83b-4bc1-9dc0-1f050c7e789b",
   "time":"20190325-15:22:50.560",
   "contenttype":"application/json",
   "shkeptncontext":"db51be80-4fee-41af-bb53-1b093d2b694c",
   "data":{  
      "githuborg":"keptn-tiger",
      "project":"sockshop",
      "service":"carts",
      "image":"keptnexamples/carts",
      "tag":"0.7.1",
      "stage":"dev",
      "teststrategy":"functional",
      "deploymentstrategy":"direct"
   }
}
```

## Start Evaluation (deprecated)

Upon receiving a *tests done* event a *start evaluation* event is dispatched that indicates that the test results can now be evaluated.

```json
{  
   "type":"sh.keptn.events.start-evaluation",
   "specversion":"0.2",
   "source":"https://github.com/keptn/keptn-service#start-evaluation",
   "id":"49ac0dec-a83b-4bc1-9dc0-1f050c7e789b",
   "time":"20190325-15:22:50.560",
   "contenttype":"application/json",
   "shkeptncontext":"db51be80-4fee-41af-bb53-1b093d2b694c",
   "data":{  
      "githuborg":"keptn-tiger",
      "project":"sockshop",
      "service":"carts",
      "image":"keptnexamples/carts",
      "tag":"0.7.1",
      "stage":"dev",
      "teststrategy":"functional",
      "deploymentstrategy":"direct"
   }
}
```

##  Evaluation Done (deprecated)

The *evaluation done* event is triggered when test results have been evaluated. 

```json
{  
   "type":"sh.keptn.events.evaluation-done",
   "specversion":"0.2",
   "source":"https://github.com/keptn/keptn-service#start-evaluation",
   "id":"49ac0dec-a83b-4bc1-9dc0-1f050c7e789b",
   "time":"20190325-15:22:50.560",
   "contenttype":"application/json",
   "shkeptncontext":"db51be80-4fee-41af-bb53-1b093d2b694c",
   "data":{  
      "githuborg":"keptn-tiger",
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

## Problem Opened (deprecated)

The *problem opened* is triggered when a monitoring provider detects a new problem.

```json
{
   "type":"sh.keptn.events.problem-opened",
   "specversion":"0.2",
   "source":"https://github.com/keptn/keptn-service#start-evaluation",
   "id":"49ac0dec-a83b-4bc1-9dc0-1f050c7e789b",
   "time":"20190325-15:22:50.560",
   "contenttype":"application/json",
   "shkeptncontext":"db51be80-4fee-41af-bb53-1b093d2b694c",
   "data": {
      "State":"{State}",
      "ProblemID":"{ProblemID}",
      "PID":"{PID}",
      "ProblemTitle":"{ProblemTitle}",
      "ProblemDetails":{ProblemDetailsJSON},
      "ImpactedEntities":{ImpactedEntities},
      "ImpactedEntity":"{ImpactedEntity}"
   }
}
```

## Problem Resolved (deprecated)

The *problem resolved* is triggered when a monitoring provider detects that a previously opened problem has been resolved.

```json
{
   "type":"sh.keptn.events.problem-resolved",
   "specversion":"0.2",
   "source":"https://github.com/keptn/keptn-service#start-evaluation",
   "id":"49ac0dec-a83b-4bc1-9dc0-1f050c7e789b",
   "time":"20190325-15:22:50.560",
   "contenttype":"application/json",
   "shkeptncontext":"db51be80-4fee-41af-bb53-1b093d2b694c",
   "data": {
      "State":"{State}",
      "ProblemID":"{ProblemID}",
      "PID":"{PID}",
      "ProblemTitle":"{ProblemTitle}",
      "ProblemDetails":{ProblemDetailsJSON},
      "ImpactedEntities":{ImpactedEntities},
      "ImpactedEntity":"{ImpactedEntity}"
   }
}
```
