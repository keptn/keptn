import { Injectable } from '@angular/core';
import {HttpEvent, HttpHandler, HttpInterceptor, HttpRequest, HttpResponse} from "@angular/common/http";
import {Observable} from "rxjs";
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
  }, {
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

@Injectable({
  providedIn: 'root'
})
export class HttpMockInterceptor implements HttpInterceptor {

  intercept(request: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {
    return next.handle(request).pipe(
      map((event: HttpEvent<any>) => {
        if (event instanceof HttpResponse) {
          if(request.url.indexOf("/api/traces/") != -1) {
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
          }
        }
        return event;
      }));
  }
}
