import { Injectable } from '@angular/core';
import {
  HttpEvent,
  HttpHandler,
  HttpInterceptor,
  HttpRequest,
  HttpResponse
} from "@angular/common/http";
import {Observable, of} from "rxjs";
import {map} from "rxjs/operators";

const problemEventMockData = [
  {
    "contenttype": "application/json",
    "data": {
      "file": "remediation.yaml",
      "fileChangesGeneratedChart": {
        "templates/simplenode-primary-deployment.yaml": "---\napiVersion: apps/v1\nkind: Deployment\nmetadata:\n  annotations:\n    deployment.kubernetes.io/revision: \"5\"\n  creationTimestamp: \"2020-01-27T08:55:39Z\"\n  generation: 5\n  name: simplenode-primary\n  resourceVersion: \"2168265\"\n  selfLink: /apis/apps/v1/namespaces/simpleproject-production/deployments/simplenode\n  uid: cabdc754-40e2-11ea-b7cf-42010a80002d\nspec:\n  progressDeadlineSeconds: 600\n  replicas: 3\n  revisionHistoryLimit: 10\n  selector:\n    matchLabels:\n      app: simplenode-primary\n  strategy:\n    rollingUpdate:\n      maxSurge: 25%\n      maxUnavailable: 0\n    type: RollingUpdate\n  template:\n    metadata:\n      creationTimestamp: null\n      labels:\n        app: simplenode-primary\n    spec:\n      containers:\n      - env:\n        - name: DT_CUSTOM_PROP\n          value: keptn_project=simpleproject keptn_service=simplenode keptn_stage=production\n            keptn_deployment=primary\n        - name: POD_NAME\n          valueFrom:\n            fieldRef:\n              apiVersion: v1\n              fieldPath: metadata.name\n        - name: DEPLOYMENT_NAME\n          valueFrom:\n            fieldRef:\n              apiVersion: v1\n              fieldPath: metadata.labels['deployment']\n        - name: CONTAINER_IMAGE\n          value: docker.io/bacherfl/simplenodeservice:4.0.0\n        - name: KEPTN_PROJECT\n          value: simpleproject\n        - name: KEPTN_STAGE\n          value: production\n        - name: KEPTN_SERVICE\n          value: simplenode\n        - name: KEPTN_DEPLOYMENT\n          value: canary\n        image: docker.io/bacherfl/simplenodeservice:4.0.0\n        imagePullPolicy: Always\n        livenessProbe:\n          failureThreshold: 3\n          httpGet:\n            path: /\n            port: 8080\n            scheme: HTTP\n          initialDelaySeconds: 30\n          periodSeconds: 10\n          successThreshold: 1\n          timeoutSeconds: 15\n        name: simplenode\n        ports:\n        - containerPort: 8080\n          name: http\n          protocol: TCP\n        readinessProbe:\n          failureThreshold: 3\n          httpGet:\n            path: /\n            port: 8080\n            scheme: HTTP\n          initialDelaySeconds: 30\n          periodSeconds: 10\n          successThreshold: 1\n          timeoutSeconds: 15\n        resources:\n          limits:\n            cpu: 100m\n            memory: 128Mi\n          requests:\n            cpu: 100m\n            memory: 128Mi\n        terminationMessagePath: /dev/termination-log\n        terminationMessagePolicy: File\n      dnsPolicy: ClusterFirst\n      restartPolicy: Always\n      schedulerName: default-scheduler\n      securityContext: {}\n      terminationGracePeriodSeconds: 30\nstatus: {}\n"
      },
      "labels": null,
      "project": "simpleproject",
      "service": "simplenode",
      "stage": "production"
    },
    "id": "d69ffa83-9c83-435b-b005-cc53438cad17",
    "source": "https://github.com/keptn/keptn/remediation-service",
    "specversion": "0.2",
    "time": "2020-01-28T11:14:49.894Z",
    "type": "Remediation file found",
    "shkeptncontext": "77fcb1db-a100-44b4-8d37-234703182fa6"
  },
  {
    "contenttype": "application/json",
    "data": {
      "file": "remediation.yaml",
      "action": "Promote flag set to: disabled",
      "fileChangesGeneratedChart": {
        "templates/simplenode-primary-deployment.yaml": "---\napiVersion: apps/v1\nkind: Deployment\nmetadata:\n  annotations:\n    deployment.kubernetes.io/revision: \"5\"\n  creationTimestamp: \"2020-01-27T08:55:39Z\"\n  generation: 5\n  name: simplenode-primary\n  resourceVersion: \"2168265\"\n  selfLink: /apis/apps/v1/namespaces/simpleproject-production/deployments/simplenode\n  uid: cabdc754-40e2-11ea-b7cf-42010a80002d\nspec:\n  progressDeadlineSeconds: 600\n  replicas: 3\n  revisionHistoryLimit: 10\n  selector:\n    matchLabels:\n      app: simplenode-primary\n  strategy:\n    rollingUpdate:\n      maxSurge: 25%\n      maxUnavailable: 0\n    type: RollingUpdate\n  template:\n    metadata:\n      creationTimestamp: null\n      labels:\n        app: simplenode-primary\n    spec:\n      containers:\n      - env:\n        - name: DT_CUSTOM_PROP\n          value: keptn_project=simpleproject keptn_service=simplenode keptn_stage=production\n            keptn_deployment=primary\n        - name: POD_NAME\n          valueFrom:\n            fieldRef:\n              apiVersion: v1\n              fieldPath: metadata.name\n        - name: DEPLOYMENT_NAME\n          valueFrom:\n            fieldRef:\n              apiVersion: v1\n              fieldPath: metadata.labels['deployment']\n        - name: CONTAINER_IMAGE\n          value: docker.io/bacherfl/simplenodeservice:4.0.0\n        - name: KEPTN_PROJECT\n          value: simpleproject\n        - name: KEPTN_STAGE\n          value: production\n        - name: KEPTN_SERVICE\n          value: simplenode\n        - name: KEPTN_DEPLOYMENT\n          value: canary\n        image: docker.io/bacherfl/simplenodeservice:4.0.0\n        imagePullPolicy: Always\n        livenessProbe:\n          failureThreshold: 3\n          httpGet:\n            path: /\n            port: 8080\n            scheme: HTTP\n          initialDelaySeconds: 30\n          periodSeconds: 10\n          successThreshold: 1\n          timeoutSeconds: 15\n        name: simplenode\n        ports:\n        - containerPort: 8080\n          name: http\n          protocol: TCP\n        readinessProbe:\n          failureThreshold: 3\n          httpGet:\n            path: /\n            port: 8080\n            scheme: HTTP\n          initialDelaySeconds: 30\n          periodSeconds: 10\n          successThreshold: 1\n          timeoutSeconds: 15\n        resources:\n          limits:\n            cpu: 100m\n            memory: 128Mi\n          requests:\n            cpu: 100m\n            memory: 128Mi\n        terminationMessagePath: /dev/termination-log\n        terminationMessagePolicy: File\n      dnsPolicy: ClusterFirst\n      restartPolicy: Always\n      schedulerName: default-scheduler\n      securityContext: {}\n      terminationGracePeriodSeconds: 30\nstatus: {}\n"
      },
      "labels": null,
      "project": "simpleproject",
      "service": "simplenode",
      "stage": "production"
    },
    "id": "d69ffa83-9c83-435b-b005-cc53438cad17",
    "source": "https://github.com/keptn/keptn/remediation-service",
    "specversion": "0.2",
    "time": "2020-02-03T11:14:49.894Z",
    "type": "Remediation file executed",
    "shkeptncontext": "77fcb1db-a100-44b4-8d37-234703182fa6"
  },
  {
    "contenttype": "application/json",
    "data": {
      "deploymentstrategy": "blue_green_service",
      "evaluationdetails": {
        "indicatorResults": [
          {
            "score": 0.5,
            "status": "warning",
            "targets": [
              {
                "criteria": "<=800",
                "targetValue": 800,
                "violated": false
              },
              {
                "criteria": "<=+10%",
                "targetValue": 0,
                "violated": false
              },
              {
                "criteria": "<600",
                "targetValue": 600,
                "violated": true
              }
            ],
            "value": {
              "metric": "response_time_p90",
              "success": true,
              "value": 638.4818203572551
            }
          },
          {
            "score": 0,
            "status": "info",
            "targets": null,
            "value": {
              "metric": "throughput",
              "success": true,
              "value": 1
            }
          },
          {
            "score": 2,
            "status": "pass",
            "targets": [
              {
                "criteria": "<=1%",
                "targetValue": 1,
                "violated": false
              }
            ],
            "value": {
              "metric": "error_rate",
              "success": true,
              "value": 0
            }
          },
          {
            "score": 0,
            "status": "info",
            "targets": null,
            "value": {
              "metric": "response_time_p50",
              "success": true,
              "value": 632.6766088362081
            }
          },
          {
            "score": 0,
            "status": "info",
            "targets": null,
            "value": {
              "metric": "response_time_p90",
              "success": true,
              "value": 638.4818203572551
            }
          }
        ],
        "result": "warning",
        "score": 83.33333333333334,
        "sloFileContent": "Y29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246IGF2ZwogIGNvbXBhcmVfd2l0aDogc2luZ2xlX3Jlc3VsdAogIGluY2x1ZGVfcmVzdWx0X3dpdGhfc2NvcmU6IHBhc3MKICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAzCmZpbHRlcjogbnVsbApvYmplY3RpdmVzOgotIGtleV9zbGk6IGZhbHNlCiAgcGFzczoKICAtIGNyaXRlcmlhOgogICAgLSA8PSsxMCUKICAgIC0gPDYwMAogIHNsaTogcmVzcG9uc2VfdGltZV9wOTAKICB3YXJuaW5nOgogIC0gY3JpdGVyaWE6CiAgICAtIDw9ODAwCiAgd2VpZ2h0OiAxCi0ga2V5X3NsaTogZmFsc2UKICBwYXNzOiBudWxsCiAgc2xpOiB0aHJvdWdocHV0CiAgd2FybmluZzogbnVsbAogIHdlaWdodDogMQotIGtleV9zbGk6IGZhbHNlCiAgcGFzczoKICAtIGNyaXRlcmlhOgogICAgLSA8PTElCiAgc2xpOiBlcnJvcl9yYXRlCiAgd2FybmluZzoKICAtIGNyaXRlcmlhOgogICAgLSA8PTIlCiAgd2VpZ2h0OiAyCi0ga2V5X3NsaTogZmFsc2UKICBwYXNzOiBudWxsCiAgc2xpOiByZXNwb25zZV90aW1lX3A1MAogIHdhcm5pbmc6IG51bGwKICB3ZWlnaHQ6IDEKLSBrZXlfc2xpOiBmYWxzZQogIHBhc3M6IG51bGwKICBzbGk6IHJlc3BvbnNlX3RpbWVfcDkwCiAgd2FybmluZzogbnVsbAogIHdlaWdodDogMQpzcGVjX3ZlcnNpb246IDAuMS4wCnRvdGFsX3Njb3JlOgogIHBhc3M6IDkwJQogIHdhcm5pbmc6IDc1JQo=",
        "timeEnd": "2020-02-03T12:57:34Z",
        "timeStart": "2020-02-03T12:55:34Z"
      },
      "labels": null,
      "project": "simpleproject",
      "result": "warning",
      "service": "simplenode",
      "stage": "production",
      "teststrategy": "real-user"
    },
    "id": "fdbbdccd-1ed8-42e0-a718-ce5eaa5c915d",
    "source": "lighthouse-service",
    "specversion": "0.2",
    "time": "2020-01-28T11:27:36.391Z",
    "type": "sh.keptn.events.evaluation-done",
    "shkeptncontext": "77fcb1db-a100-44b4-8d37-234703182fa6"
  },
  {
    "contenttype": "application/json",
    "data": {
      "ImpactedEntity": "Remediation successful",
      "PID": "6140028070610368294",
      "ProblemDetails": {
        "displayName": "294",
        "endTime": -1,
        "hasRootCause": true,
        "id": "6140028070610368294_1580209560000V2",
        "impactLevel": "SERVICE",
        "severityLevel": "PERFORMANCE",
        "startTime": 1580209560000,
        "status": "OPEN"
      },
      "ProblemID": "294",
      "ProblemTitle": "Response time degradation",
      "State": "OPEN",
      "labels": null,
      "project": "simpleproject",
      "service": "simplenode",
      "stage": "production"
    },
    "id": "9511b28b-1b08-4ee3-936c-3b75c0bfe476",
    "source": "dynatrace",
    "specversion": "0.2",
    "time": "2020-01-28T11:14:49.689Z",
    "type": "sh.keptn.event.problem.close",
    "shkeptncontext": "77fcb1db-a100-44b4-8d37-234703182fa6"
  }
];
const projectsMockData = [
  {
    "projectName": "Dynatrace SaaS"
  },
  {
    "projectName": "AppMon SaaS"
  },
  {
    "projectName": "Keptn Bridge"
  },
  {
    "projectName": "Dynatrace Merch Shop"
  }
];
const stagesMockData = {
  "Dynatrace SaaS": [{"stageName":"dev"},{"stageName":"staging"},{"stageName":"production"}],
  "AppMon SaaS": [{"stageName":"dev"},{"stageName":"staging"},{"stageName":"production"}],
  "Dynatrace Merch Shop": [{"stageName":"dev"},{"stageName":"production"}],
  "Keptn Bridge": [{"stageName":"dev"},{"stageName":"staging"},{"stageName":"production"}],
};
const servicesMockData = {
  "Dynatrace SaaS": [{"serviceName":"configuration-service"},{"serviceName":"monitoring-service"}, {"serviceName":"analysis-service"}, {"serviceName":"davis-service"},{"serviceName":"casandra-db"}, {"serviceName":"configuration-service"},{"serviceName":"monitoring-service"}, {"serviceName":"analysis-service"}, {"serviceName":"davis-service"},{"serviceName":"casandra-db"}],
  "AppMon SaaS": [{"serviceName":"configuration-service"},{"serviceName":"monitoring-service"}, {"serviceName":"analysis-service"}, {"serviceName":"davis-service"},{"serviceName":"casandra-db"}, {"serviceName":"configuration-service"},{"serviceName":"monitoring-service"}, {"serviceName":"analysis-service"}],
  "Dynatrace Merch Shop": [{"serviceName":"carts-service"},{"serviceName":"payment-service"}, {"serviceName":"products-service"}, {"serviceName":"frontend"},{"serviceName":"mongo-db"}],
  "Keptn Bridge": [{"serviceName":"bridge"},{"serviceName":"api-service"}, {"serviceName":"mongo-db"}],
};
const rootEventsMockData = {
  "Keptn Bridge": [
    {
      "contenttype": "application/json",
      "data": {
        "canary": {
          "action": "set",
          "value": 100
        },
        "eventContext": null,
        "labels": null,
        "project": "ck-sockshop",
        "service": "carts",
        "stage": "",
        "valuesCanary": {
          "image": "docker.io/keptn/bridge:0.10.4"
        }
      },
      "id": "a52a64bf-6c45-4c4f-bde5-4911cf9d46e2",
      "source": "https://github.com/keptn/keptn/api",
      "specversion": "0.2",
      "time": "2020-02-03T08:34:15.194Z",
      "type": "sh.keptn.event.configuration.change",
      "shkeptncontext": "b826e543-1853-4624-90ef-2eaf841ff528"
    },
    {
      "contenttype": "application/json",
      "data": {
        "valuesCanary": {
          "image": "remediation.yaml executed"
        },
        "ImpactedEntities": [
          {
            "entity": "SERVICE-65FBE40CCBCEA563",
            "name": "carts.production.primary",
            "type": "SERVICE"
          }
        ],
        "ImpactedEntity": "2 service problems on Web service carts.production.primary",
        "PID": "-1113367559322804549",
        "ProblemDetails": {
          "affectedCounts": {
            "APPLICATION": 0,
            "ENVIRONMENT": 0,
            "INFRASTRUCTURE": 0,
            "SERVICE": 1
          },
          "commentCount": 2,
          "displayName": "549",
          "endTime": -1,
          "hasRootCause": true,
          "id": "-1113367559322804549_1580714280000V2",
          "impactLevel": "SERVICE",
          "rankedEvents": [
            {
              "affectedRequestsPerMinute": 121.2,
              "endTime": 1580715000000,
              "entityId": "SERVICE-65FBE40CCBCEA563",
              "entityName": "carts.production.primary",
              "eventType": "SERVICE_RESPONSE_TIME_DEGRADED",
              "impactLevel": "SERVICE",
              "isRootCause": true,
              "percentile": "50th",
              "referenceResponseTime50thPercentile": 0,
              "referenceResponseTime90thPercentile": 0,
              "service": "carts.production.primary",
              "serviceMethodGroup": "Default requests",
              "severities": [
                {
                  "context": "RESPONSE_TIME_50TH_PERCENTILE",
                  "unit": "MicroSecond (µs)",
                  "value": 453924
                }
              ],
              "severityLevel": "PERFORMANCE",
              "startTime": 1580714280000,
              "status": "CLOSED",
              "userDefined50thPercentileThreshold": 100000,
              "userDefined90thPercentileThreshold": 1000000
            },
            {
              "affectedRequestsPerMinute": 141,
              "endTime": -1,
              "entityId": "SERVICE-65FBE40CCBCEA563",
              "entityName": "carts.production.primary",
              "eventType": "SERVICE_RESPONSE_TIME_DEGRADED",
              "impactLevel": "SERVICE",
              "isRootCause": true,
              "percentile": "50th",
              "referenceResponseTime50thPercentile": 0,
              "referenceResponseTime90thPercentile": 0,
              "service": "carts.production.primary",
              "serviceMethod": "addToCart",
              "severities": [
                {
                  "context": "RESPONSE_TIME_50TH_PERCENTILE",
                  "unit": "MicroSecond (µs)",
                  "value": 446051
                }
              ],
              "severityLevel": "PERFORMANCE",
              "startTime": 1580718300000,
              "status": "OPEN",
              "userDefined50thPercentileThreshold": 100000,
              "userDefined90thPercentileThreshold": 1000000
            },
            {
              "affectedRequestsPerMinute": 353,
              "endTime": -1,
              "entityId": "SERVICE-65FBE40CCBCEA563",
              "entityName": "carts.production.primary",
              "eventType": "FAILURE_RATE_INCREASED",
              "impactLevel": "SERVICE",
              "isRootCause": true,
              "service": "carts.production.primary",
              "serviceMethodGroup": "Default requests",
              "severities": [
                {
                  "context": "FAILURE_RATE",
                  "unit": "Ratio",
                  "value": 1
                }
              ],
              "severityLevel": "ERROR",
              "startTime": 1580718240000,
              "status": "OPEN"
            },
            {
              "affectedRequestsPerMinute": 206.2,
              "endTime": 1580717340000,
              "entityId": "SERVICE-65FBE40CCBCEA563",
              "entityName": "carts.production.primary",
              "eventType": "SERVICE_RESPONSE_TIME_DEGRADED",
              "impactLevel": "SERVICE",
              "isRootCause": true,
              "percentile": "50th",
              "referenceResponseTime50thPercentile": 0,
              "referenceResponseTime90thPercentile": 0,
              "service": "carts.production.primary",
              "serviceMethodGroup": "Default requests",
              "severities": [
                {
                  "context": "RESPONSE_TIME_50TH_PERCENTILE",
                  "unit": "MicroSecond (µs)",
                  "value": 450204.5
                }
              ],
              "severityLevel": "PERFORMANCE",
              "startTime": 1580715540000,
              "status": "CLOSED",
              "userDefined50thPercentileThreshold": 100000,
              "userDefined90thPercentileThreshold": 1000000
            },
            {
              "affectedRequestsPerMinute": 353,
              "endTime": -1,
              "entityId": "SERVICE-65FBE40CCBCEA563",
              "entityName": "carts.production.primary",
              "eventType": "FAILURE_RATE_INCREASED",
              "impactLevel": "SERVICE",
              "isRootCause": true,
              "service": "carts.production.primary",
              "serviceMethod": "addToCart",
              "severities": [
                {
                  "context": "FAILURE_RATE",
                  "unit": "Ratio",
                  "value": 1
                }
              ],
              "severityLevel": "ERROR",
              "startTime": 1580718240000,
              "status": "OPEN"
            },
            {
              "affectedRequestsPerMinute": 141,
              "endTime": -1,
              "entityId": "SERVICE-65FBE40CCBCEA563",
              "entityName": "carts.production.primary",
              "eventType": "SERVICE_RESPONSE_TIME_DEGRADED",
              "impactLevel": "SERVICE",
              "isRootCause": true,
              "percentile": "50th",
              "referenceResponseTime50thPercentile": 0,
              "referenceResponseTime90thPercentile": 0,
              "service": "carts.production.primary",
              "serviceMethodGroup": "Default requests",
              "severities": [
                {
                  "context": "RESPONSE_TIME_50TH_PERCENTILE",
                  "unit": "MicroSecond (µs)",
                  "value": 446051
                }
              ],
              "severityLevel": "PERFORMANCE",
              "startTime": 1580718300000,
              "status": "OPEN",
              "userDefined50thPercentileThreshold": 100000,
              "userDefined90thPercentileThreshold": 1000000
            }
          ],
          "rankedImpacts": [
            {
              "entityId": "SERVICE-65FBE40CCBCEA563",
              "entityName": "carts.production.primary",
              "eventType": "SERVICE_RESPONSE_TIME_DEGRADED",
              "impactLevel": "SERVICE",
              "severityLevel": "PERFORMANCE"
            },
            {
              "entityId": "SERVICE-65FBE40CCBCEA563",
              "entityName": "carts.production.primary",
              "eventType": "SERVICE_RESPONSE_TIME_DEGRADED",
              "impactLevel": "SERVICE",
              "severityLevel": "PERFORMANCE"
            },
            {
              "entityId": "SERVICE-65FBE40CCBCEA563",
              "entityName": "carts.production.primary",
              "eventType": "FAILURE_RATE_INCREASED",
              "impactLevel": "SERVICE",
              "severityLevel": "ERROR"
            },
            {
              "entityId": "SERVICE-65FBE40CCBCEA563",
              "entityName": "carts.production.primary",
              "eventType": "SERVICE_RESPONSE_TIME_DEGRADED",
              "impactLevel": "SERVICE",
              "severityLevel": "PERFORMANCE"
            },
            {
              "entityId": "SERVICE-65FBE40CCBCEA563",
              "entityName": "carts.production.primary",
              "eventType": "FAILURE_RATE_INCREASED",
              "impactLevel": "SERVICE",
              "severityLevel": "ERROR"
            },
            {
              "entityId": "SERVICE-65FBE40CCBCEA563",
              "entityName": "carts.production.primary",
              "eventType": "SERVICE_RESPONSE_TIME_DEGRADED",
              "impactLevel": "SERVICE",
              "severityLevel": "PERFORMANCE"
            }
          ],
          "recoveredCounts": {
            "APPLICATION": 0,
            "ENVIRONMENT": 0,
            "INFRASTRUCTURE": 0,
            "SERVICE": 0
          },
          "severityLevel": "ERROR",
          "startTime": 1580714280000,
          "status": "OPEN",
          "tagsOfAffectedEntities": [
            {
              "context": "CONTEXTLESS",
              "key": "keptn_service",
              "value": "carts"
            },
            {
              "context": "CONTEXTLESS",
              "key": "keptn_project",
              "value": "ck-sockshop"
            },
            {
              "context": "CONTEXTLESS",
              "key": "keptn_stage",
              "value": "production"
            },
            {
              "context": "CONTEXTLESS",
              "key": "keptn_deployment",
              "value": "primary"
            }
          ]
        },
        "ProblemID": "549",
        "ProblemTitle": "2 service problems",
        "State": "OPEN",
        "Tags": "keptn_service:carts, keptn_deployment:primary, keptn_stage:production, keptn_project:ck-sockshop",
        "eventContext": null
      },
      "id": "07f22861-be2f-4a45-810f-2ddcccc91278",
      "source": "https://github.com/keptn/keptn/api",
      "specversion": "0.2",
      "time": "2020-01-31T07:59:04.230Z",
      "type": "sh.keptn.events.problem",
      "shkeptncontext": "5f599164-49da-4032-8161-8f1b530fbae9"
    },
    {
      "contenttype": "application/json",
      "data": {
        "canary": {
          "action": "set",
          "value": 100
        },
        "eventContext": null,
        "labels": null,
        "project": "ck-sockshop",
        "service": "carts",
        "stage": "",
        "valuesCanary": {
          "image": "docker.io/keptn/bridge:0.10.3"
        }
      },
      "id": "07a01805-7a37-4b02-96d6-88b90aab6085",
      "source": "https://github.com/keptn/keptn/api",
      "specversion": "0.2",
      "time": "2020-01-30T09:14:32.551Z",
      "type": "sh.keptn.event.configuration.change",
      "shkeptncontext": "9fe0df6e-1ce8-4c47-a1cf-80d1f187bd80"
    },
    {
      "contenttype": "application/json",
      "data": {
        "canary": {
          "action": "set",
          "value": 100
        },
        "eventContext": null,
        "labels": null,
        "project": "ck-sockshop",
        "service": "carts",
        "stage": "",
        "valuesCanary": {
          "image": "docker.io/keptn/bridge:0.10.2"
        }
      },
      "id": "d45045d2-9b66-45db-8e49-8a233edbef95",
      "source": "https://github.com/keptn/keptn/api",
      "specversion": "0.2",
      "time": "2020-01-24T09:49:38.095Z",
      "type": "sh.keptn.event.configuration.change",
      "shkeptncontext": "a8c014cd-b08d-4731-955b-8a7f50496083"
    },
    {
      "contenttype": "application/json",
      "data": {
        "canary": {
          "action": "set",
          "value": 100
        },
        "eventContext": null,
        "labels": null,
        "project": "ck-sockshop",
        "service": "carts",
        "stage": "",
        "valuesCanary": {
          "image": "docker.io/keptn/bridge:0.10.1"
        }
      },
      "id": "6ce1c3f8-9c6f-4304-a3c2-205502245090",
      "source": "https://github.com/keptn/keptn/api",
      "specversion": "0.2",
      "time": "2020-01-24T09:26:00.536Z",
      "type": "sh.keptn.event.configuration.change",
      "shkeptncontext": "5f334de2-924d-419a-8111-d07bd41ba483"
    },
    {
      "contenttype": "application/json",
      "data": {
        "deploymentStrategies": {},
        "eventContext": null,
        "helmChart": "H4sIFAAAAAAA/ykAK2FIUjBjSE02THk5NWIzVjBkUzVpWlM5Nk9WVjZNV2xqYW5keVRRbz1IZWxtAOxXW2/qRhDOs3/FKC/nKWZNIESW+hAReg5tAItLpKqq0MaewDZr73Z3jQ464r9XvmJs0lRqLo3E92I8Mzs7s8yMv/WpMrrVX1Nl7C0N+dkbgBBCet1u+iSE1J/EaV+dOZdXTrvd6ZG2c0ac9iXpngF5i2DqiLWh6oz8573qyX0SUMnuUWkmIhc2jhWg9hWTJn2/gW/IQ/CT4oBHoUCj2jAfIS0aK6IhuvnvTeGD2I5NrI/O6oR/i6z/N5THqN9qALzU/6R3We//bqdz6v/3AAvpCl0IhP+Eymai9YTSRPidhpKjbqXl4RLbIbZjKZSc+bQv4si44Hx06Ce8ArL+NxhKTg3qVoCSi22I0SvSgRf6/7J7Ve//bueyfer/98DFxQVYVRJApdStjWM9sShw4bYsBytEQwNqqGsBVD/9WqKfyPLpoF348QPs++yLUh0ZsNtZANooanC1TZYAKME5i1YLGVCDmQggpN8XEd1QxukDRxdIKjdbiS5MqwsSd8jRN0Jla0Nq/PUdfUCuC2dUyiJSgKLOc+siI8ht+cHK2lqAItMEvogMZRGq0vzi4FQKF/l8Pa+cSSqC3e780MiLOfcEZ/7WheHjWBhPoU4OvrCSIpnF5Wux39oYWQoBpBJG+IK7MO97FXkZsCeUceGaXJNSi9Gm6fd2vuwvZvPJaOlNJ1VPKVlw4Tz9UiylEn+ib36qJJgq7FwBux1kljl7bFoWtHJvaejqmJ3JDi632o+qpuled3DQRW7e5HY5vhkN6ln9rEToVoQAjwx5MMXHQ2ku96hZu3Be1JGdOG9udjvw7ia/jQbj+VvsmdXs71/2GX/5oxlDfzKe3wzHg+lyOLr52ojhpQIt3Pw68ObjpB5+GfTnzzjJ7pJjGv6Ti9n8SBSvcPpaUv/IX5BvOpjeD/vHkk97trlsMb4b3My+pQsH0+ViendkbdJ/bqsVRxypXtv58yLATYtKtnfK2QYj1NpT4gGrCSUOvqI5zFGmybXWSLlZH2qa/QvAImYY5bfI6XaGvogC7cJV1UKiYiIodU5VZ1iIIjZ7ZbdUKqQB+4xhaxErH3U1PM5CZnS9lnwZJ34JCWuKEEOhti60Sed6xCpKhX/FqJ/x1H3ekUPanRH7f12O6/wvn8WvehV86f7X7nVr/K/TI70T/3sPNPhfSf1mWSVYByypQfwSYcbM+jzWBtUwoQoZVUl0DZZSzIHspUFWDFUrNAccpaR4GU2rkLKPPrwTTjjhhE+MvwMAAP//UakOzQAaAAA=",
        "project": "ck-sockshop",
        "service": "carts"
      },
      "id": "f95686c3-bce8-4477-9391-6602fdb982f9",
      "source": "https://github.com/keptn/keptn/api",
      "specversion": "0.2",
      "time": "2020-01-23T07:41:05.391Z",
      "type": "sh.keptn.internal.event.service.create",
      "shkeptncontext": "7f3f031e-05c4-4c62-9904-f13e67ffcbc0"
    }
  ]
};
const tracesMockData = {
  "a8c014cd-b08d-4731-955b-8a7f50496083": [{"contenttype":"application/json","data":{"canary":{"action":"set","value":100},"eventContext":null,"labels":null,"project":"ck-sockshop","service":"carts","stage":"","valuesCanary":{"image":"docker.io/keptnexamples/carts:0.10.2"}},"id":"d45045d2-9b66-45db-8e49-8a233edbef95","source":"https://github.com/keptn/keptn/api","specversion":"0.2","time":"2020-01-24T09:49:38.095Z","type":"sh.keptn.event.configuration.change","shkeptncontext":"a8c014cd-b08d-4731-955b-8a7f50496083"},{"contenttype":"application/json","data":{"deploymentstrategy":"direct","image":"docker.io/keptnexamples/carts","labels":null,"project":"ck-sockshop","service":"carts","stage":"dev","tag":"0.10.2","teststrategy":"functional"},"id":"5871b2a3-18d1-4278-8504-50410bbf6884","source":"helm-service","specversion":"0.2","time":"2020-01-24T09:50:56.469Z","type":"sh.keptn.events.deployment-finished","shkeptncontext":"a8c014cd-b08d-4731-955b-8a7f50496083"},{"contenttype":"application/json","data":{"deploymentstrategy":"direct","end":"2020-01-24T09:51:03Z","labels":null,"project":"ck-sockshop","result":"pass","service":"carts","stage":"dev","start":"2020-01-24T09:50:56Z","teststrategy":"functional"},"id":"4874ff08-77c4-4f23-9ba5-49c9c66fb401","source":"jmeter-service","specversion":"0.2","time":"2020-01-24T09:51:03.744Z","type":"sh.keptn.events.tests-finished","shkeptncontext":"a8c014cd-b08d-4731-955b-8a7f50496083"},{"contenttype":"application/json","data":{"deploymentstrategy":"direct","evaluationdetails":{"indicatorResults":null,"result":"no evaluation performed by lighthouse service (TestStrategy=functional)","score":0,"sloFileContent":"","timeEnd":"2020-01-24T09:51:03Z","timeStart":"2020-01-24T09:50:56Z"},"labels":null,"project":"ck-sockshop","result":"pass","service":"carts","stage":"dev","teststrategy":"functional"},"id":"aaac1a0b-ef2e-4273-94c4-0b870ef13be7","source":"lighthouse-service","specversion":"0.2","time":"2020-01-24T09:51:03.767Z","type":"sh.keptn.events.evaluation-done","shkeptncontext":"a8c014cd-b08d-4731-955b-8a7f50496083"},{"contenttype":"application/json","data":{"canary":{"action":"promote"},"labels":null,"project":"ck-sockshop","service":"carts","stage":"dev"},"id":"f7235a63-7cb7-482d-bb4b-f32aa2199f71","source":"gatekeeper-service","specversion":"0.2","time":"2020-01-24T09:51:03.788Z","type":"sh.keptn.event.configuration.change","shkeptncontext":"a8c014cd-b08d-4731-955b-8a7f50496083"},{"contenttype":"application/json","data":{"canary":{"action":"set","value":100},"labels":null,"project":"ck-sockshop","service":"carts","stage":"staging","valuesCanary":{"image":"docker.io/keptnexamples/carts:0.10.2"}},"id":"adc8f951-03ac-4308-acf5-cfdda22253bf","source":"gatekeeper-service","specversion":"0.2","time":"2020-01-24T09:51:05.846Z","type":"sh.keptn.event.configuration.change","shkeptncontext":"a8c014cd-b08d-4731-955b-8a7f50496083"},{"contenttype":"application/json","data":{"deploymentstrategy":"blue_green_service","image":"docker.io/keptnexamples/carts","labels":null,"project":"ck-sockshop","service":"carts","stage":"staging","tag":"0.10.2","teststrategy":"performance"},"id":"e6c84ed4-53af-4a0e-bd03-0026d7af6e74","source":"helm-service","specversion":"0.2","time":"2020-01-24T09:52:35.668Z","type":"sh.keptn.events.deployment-finished","shkeptncontext":"a8c014cd-b08d-4731-955b-8a7f50496083"},{"contenttype":"application/json","data":{"deploymentstrategy":"blue_green_service","end":"2020-01-24T10:09:41Z","labels":null,"project":"ck-sockshop","result":"pass","service":"carts","stage":"staging","start":"2020-01-24T09:52:35Z","teststrategy":"performance"},"id":"1c95bdf3-4b76-44d4-94b6-0ee0346a9f78","source":"jmeter-service","specversion":"0.2","time":"2020-01-24T10:09:41.510Z","type":"sh.keptn.events.tests-finished","shkeptncontext":"a8c014cd-b08d-4731-955b-8a7f50496083"},{"contenttype":"application/json","data":{"customFilters":[],"deployment":"canary","deploymentstrategy":"blue_green_service","end":"2020-01-24T10:09:41Z","indicators":["response_time_p95","response_time_p50","error_rate","throughput"],"labels":null,"project":"ck-sockshop","service":"carts","sliProvider":"dynatrace","stage":"staging","start":"2020-01-24T09:52:35Z","teststrategy":"performance"},"id":"e69eda07-8776-4a72-b8b0-7e385cac4981","source":"lighthouse-service","specversion":"0.2","time":"2020-01-24T10:09:42.353Z","type":"sh.keptn.internal.event.get-sli","shkeptncontext":"a8c014cd-b08d-4731-955b-8a7f50496083"},{"contenttype":"application/json","data":{"deployment":"canary","deploymentstrategy":"blue_green_service","end":"2020-01-24T10:09:41Z","indicatorValues":[{"metric":"response_time_p95","success":true,"value":2022.0465677022669},{"metric":"response_time_p50","success":true,"value":2015.5157333082732},{"metric":"error_rate","success":true,"value":0},{"metric":"throughput","success":true,"value":4}],"labels":null,"project":"ck-sockshop","service":"carts","stage":"staging","start":"2020-01-24T09:52:35Z","teststrategy":"performance"},"id":"f5da24c1-2ba0-4773-bc1b-9d9708cbbf51","source":"dynatrace-sli-service","specversion":"0.2","time":"2020-01-24T10:11:45.175Z","type":"sh.keptn.internal.event.get-sli.done","shkeptncontext":"a8c014cd-b08d-4731-955b-8a7f50496083"},{"contenttype":"application/json","data":{"deploymentstrategy":"blue_green_service","evaluationdetails":{"indicatorResults":[{"score":0,"status":"fail","targets":[{"criteria":"<=900","targetValue":900,"violated":true},{"criteria":"<600","targetValue":600,"violated":true}],"value":{"metric":"response_time_p95","success":true,"value":2022.0465677022669}},{"score":0,"status":"fail","targets":[{"criteria":"<=800","targetValue":800,"violated":true},{"criteria":"<300","targetValue":300,"violated":true}],"value":{"metric":"response_time_p50","success":true,"value":2015.5157333082732}},{"score":1,"status":"pass","targets":[{"criteria":"=0","targetValue":0,"violated":false}],"value":{"metric":"error_rate","success":true,"value":0}},{"score":0,"status":"info","targets":null,"value":{"metric":"throughput","success":true,"value":4}}],"result":"fail","score":33.33333333333333,"sloFileContent":"Y29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246IGF2ZwogIGNvbXBhcmVfd2l0aDogc2V2ZXJhbF9yZXN1bHRzCiAgaW5jbHVkZV9yZXN1bHRfd2l0aF9zY29yZTogcGFzcwogIG51bWJlcl9vZl9jb21wYXJpc29uX3Jlc3VsdHM6IDMKZmlsdGVyOiBudWxsCm9iamVjdGl2ZXM6Ci0ga2V5X3NsaTogZmFsc2UKICBwYXNzOgogIC0gY3JpdGVyaWE6CiAgICAtIDw2MDAKICBzbGk6IHJlc3BvbnNlX3RpbWVfcDk1CiAgd2FybmluZzoKICAtIGNyaXRlcmlhOgogICAgLSA8PTkwMAogIHdlaWdodDogMQotIGtleV9zbGk6IGZhbHNlCiAgcGFzczoKICAtIGNyaXRlcmlhOgogICAgLSA8MzAwCiAgc2xpOiByZXNwb25zZV90aW1lX3A1MAogIHdhcm5pbmc6CiAgLSBjcml0ZXJpYToKICAgIC0gPD04MDAKICB3ZWlnaHQ6IDEKLSBrZXlfc2xpOiBmYWxzZQogIHBhc3M6CiAgLSBjcml0ZXJpYToKICAgIC0gPTAKICBzbGk6IGVycm9yX3JhdGUKICB3YXJuaW5nOiBudWxsCiAgd2VpZ2h0OiAxCi0ga2V5X3NsaTogZmFsc2UKICBwYXNzOiBudWxsCiAgc2xpOiB0aHJvdWdocHV0CiAgd2FybmluZzogbnVsbAogIHdlaWdodDogMQpzcGVjX3ZlcnNpb246IDAuMS4xCnRvdGFsX3Njb3JlOgogIHBhc3M6IDkwJQogIHdhcm5pbmc6IDc1JQo=","timeEnd":"2020-01-24T10:09:41Z","timeStart":"2020-01-24T09:52:35Z"},"labels":null,"project":"ck-sockshop","result":"fail","service":"carts","stage":"staging","teststrategy":"performance"},"id":"efabc546-48d1-48d1-a29e-b8dfd25cf3a0","source":"lighthouse-service","specversion":"0.2","time":"2020-01-24T10:11:46.085Z","type":"sh.keptn.events.evaluation-done","shkeptncontext":"a8c014cd-b08d-4731-955b-8a7f50496083"},{"contenttype":"application/json","data":{"canary":{"action":"discard"},"labels":null,"project":"ck-sockshop","service":"carts","stage":"staging"},"id":"0afb886f-5b3a-492a-bd0e-3b308a2c44b7","source":"gatekeeper-service","specversion":"0.2","time":"2020-01-24T10:11:46.105Z","type":"sh.keptn.event.configuration.change","shkeptncontext":"a8c014cd-b08d-4731-955b-8a7f50496083"}]
}

@Injectable({
  providedIn: 'root'
})
export class HttpMockInterceptor implements HttpInterceptor {

  intercept(request: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {
    switch(request.url) {
      case '/api/project/AppMon SaaS/stage':
        return of(new HttpResponse({ status: 200, body: stagesMockData['AppMon SaaS'] }));
      case '/api/project/AppMon SaaS/stage/dev/service':
      case '/api/project/AppMon SaaS/stage/staging/service':
      case '/api/project/AppMon SaaS/stage/production/service':
        return of(new HttpResponse({ status: 200, body: servicesMockData['AppMon SaaS'] }));
      case '/api/project/Dynatrace Merch Shop/stage':
        return of(new HttpResponse({ status: 200, body: stagesMockData['Dynatrace Merch Shop'] }));
      case '/api/project/Dynatrace Merch Shop/stage/dev/service':
      case '/api/project/Dynatrace Merch Shop/stage/staging/service':
      case '/api/project/Dynatrace Merch Shop/stage/production/service':
        return of(new HttpResponse({ status: 200, body: servicesMockData['Dynatrace Merch Shop'] }));
      case '/api/project/Keptn Bridge/stage':
        return of(new HttpResponse({ status: 200, body: stagesMockData['Keptn Bridge'] }));
      case '/api/project/Keptn Bridge/stage/dev/service':
      case '/api/project/Keptn Bridge/stage/staging/service':
      case '/api/project/Keptn Bridge/stage/production/service':
        return of(new HttpResponse({ status: 200, body: servicesMockData['Keptn Bridge'] }));
      case '/api/project/Dynatrace SaaS/stage':
        return of(new HttpResponse({ status: 200, body: stagesMockData['Dynatrace SaaS'] }));
      case '/api/project/Dynatrace SaaS/stage/dev/service':
      case '/api/project/Dynatrace SaaS/stage/staging/service':
      case '/api/project/Dynatrace SaaS/stage/production/service':
        return of(new HttpResponse({ status: 200, body: servicesMockData['Dynatrace SaaS'] }));

      case '/api/roots/Keptn Bridge/bridge':
        return of(new HttpResponse({ status: 200, body: rootEventsMockData['Keptn Bridge'] }));
    }

    if(request.url.indexOf('/api/traces/a8c014cd-b08d-4731-955b-8a7f50496083') != -1) {
      console.log("traces for a8c014cd-b08d-4731-955b-8a7f50496083");
      return of(new HttpResponse({ status: 200, body: tracesMockData['a8c014cd-b08d-4731-955b-8a7f50496083'] }));
    }

    return next.handle(request).pipe(
      map((event: HttpEvent<any>) => {
        if (event instanceof HttpResponse) {
          if(request.url.indexOf('/api/traces/') != -1) {
            if(event.body && event.body[event.body.length-1].type == 'sh.keptn.event.problem.open') {
              let time = event.body[event.body.length-1].time;
              let stage = event.body[event.body.length-1].data.stage;

              let mockEvents = JSON.parse(JSON.stringify( problemEventMockData ));
              mockEvents
                .map((event) => {
                  event.time = time;
                  event.data.stage = stage;
                  return event;
                });
              event = event.clone({ body: event.body.concat(mockEvents) });
            }
          } else if(request.url == '/api/project') {
            event = event.clone({ body: projectsMockData });
          }
        }
        return event;
      }));
  }
}
