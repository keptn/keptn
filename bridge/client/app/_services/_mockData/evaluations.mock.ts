import { Trace } from '../../_models/trace';

const evaluationData = {
  "traces": [],
  "data": {
  "evaluation": {
    "comparedEvents": [
      "c38712d2-0e2f-4412-bd39-fed91f333510"
    ],
      "gitCommit": "",
      "indicatorResults": [
      {
        "displayName": "",
        "keySli": false,
        "passTargets": [
          {
            "criteria": "<=+10%",
            "targetValue": 392.5944454106811,
            "violated": false
          },
          {
            "criteria": "<600",
            "targetValue": 600,
            "violated": false
          }
        ],
        "score": 1,
        "status": "pass",
        "value": {
          "metric": "response_time_p95",
          "success": true,
          "value": 315.8696785366614
        },
        "warningTargets": [
          {
            "criteria": "<=800",
            "targetValue": 800,
            "violated": false
          }
        ]
      }
    ],
      "result": "pass",
      "score": 100,
      "sloFileContent": "LS0tCnNwZWNfdmVyc2lvbjogIjEuMCIKY29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciCiAgY29tcGFyZV93aXRoOiAic2luZ2xlX3Jlc3VsdCIKICBpbmNsdWRlX3Jlc3VsdF93aXRoX3Njb3JlOiAicGFzcyIKICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxCmZpbHRlcjoKb2JqZWN0aXZlczoKICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1IgogICAga2V5X3NsaTogZmFsc2UKICAgIHBhc3M6ICAgICAgICAgICAgICMgcGFzcyBpZiAocmVsYXRpdmUgY2hhbmdlIDw9IDEwJSBBTkQgYWJzb2x1dGUgdmFsdWUgaXMgPCA2MDBtcykKICAgICAgLSBjcml0ZXJpYToKICAgICAgICAgIC0gIjw9KzEwJSIgICMgcmVsYXRpdmUgdmFsdWVzIHJlcXVpcmUgYSBwcmVmaXhlZCBzaWduIChwbHVzIG9yIG1pbnVzKQogICAgICAgICAgLSAiPDYwMCIgICAgIyBhYnNvbHV0ZSB2YWx1ZXMgb25seSByZXF1aXJlIGEgbG9naWNhbCBvcGVyYXRvcgogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nCiAgICAgIC0gY3JpdGVyaWE6CiAgICAgICAgICAtICI8PTgwMCIKICAgIHdlaWdodDogMQp0b3RhbF9zY29yZToKICBwYXNzOiAiOTAlIgogIHdhcm5pbmc6ICI3NSUi",
      "timeEnd": "2021-04-08T15:45:02Z",
      "timeStart": "2021-04-08T15:43:09Z"
  },
  "labels": {
    "DtCreds": "dynatrace"
  },
  "project": "sockshop",
    "result": "pass",
    "service": "carts",
    "stage": "staging",
    "status": "succeeded",
    "evaluationHistory": [
    {
      "traces": [],
      "contenttype": "application/json",
      "data": {
        "project": "sockshop",
        "stage": "staging",
        "service": "carts",
        "result": "pass",
        "evaluation": {
          "timeStart": "2020-11-10T11:08:23Z",
          "timeEnd": "2020-11-10T11:10:07Z",
          "result": "pass",
          "score": 100,
          "sloFileContent": "LS0tCnNwZWNfdmVyc2lvbjogIjEuMCIKY29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciCiAgY29tcGFyZV93aXRoOiAic2luZ2xlX3Jlc3VsdCIKICBpbmNsdWRlX3Jlc3VsdF93aXRoX3Njb3JlOiAicGFzcyIKICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxCmZpbHRlcjoKb2JqZWN0aXZlczoKICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1IgogICAga2V5X3NsaTogZmFsc2UKICAgIHBhc3M6ICAgICAgICAgICAgICMgcGFzcyBpZiAocmVsYXRpdmUgY2hhbmdlIDw9IDEwJSBBTkQgYWJzb2x1dGUgdmFsdWUgaXMgPCA2MDBtcykKICAgICAgLSBjcml0ZXJpYToKICAgICAgICAgIC0gIjw9KzEwJSIgICMgcmVsYXRpdmUgdmFsdWVzIHJlcXVpcmUgYSBwcmVmaXhlZCBzaWduIChwbHVzIG9yIG1pbnVzKQogICAgICAgICAgLSAiPDYwMCIgICAgIyBhYnNvbHV0ZSB2YWx1ZXMgb25seSByZXF1aXJlIGEgbG9naWNhbCBvcGVyYXRvcgogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nCiAgICAgIC0gY3JpdGVyaWE6CiAgICAgICAgICAtICI8PTgwMCIKICAgIHdlaWdodDogMQp0b3RhbF9zY29yZToKICBwYXNzOiAiOTAlIgogIHdhcm5pbmc6ICI3NSUi",
          "indicatorResults": [
            {
              "score": 1,
              "value": {
                "metric": "response_time_p95",
                "value": 299.18637492576534,
                "success": true
              },
              "displayName": "",
              "passTargets": null,
              "warningTargets": null,
              "keySli": false,
              "status": "pass"
            },
            {
              "score": 4,
              "value": {
                "metric": "response_time_p90",
                "value": 250.18637492576534,
                "success": true
              },
              "displayName": "",
              "passTargets": null,
              "warningTargets": null,
              "keySli": false,
              "status": "pass"
            },
            {
              "score": 2,
              "value": {
                "metric": "response_time_p50",
                "value": 200.18637492576534,
                "success": true
              },
              "displayName": "",
              "passTargets": null,
              "warningTargets": null,
              "keySli": false,
              "status": "pass"
            }
          ],
          "comparedEvents": [
            "cfa408ce-f552-43c4-aff2-673b1e0548d2"
          ],
          "gitCommit": ""
        }
      },
      "id": "04266cc2-eeea-485b-85b3-f1dea50890ce",
      "source": "lighthouse-service",
      "specversion": "1.0",
      "time": "2020-11-10T11:12:12.364Z",
      "type": "sh.keptn.event.evaluation.finished",
      "shkeptncontext": "24534d07-3a68-4d2b-9a6e-56a1a1a55d2f",
      "plainEvent": {
        "contenttype": "application/json",
        "data": {
          "project": "sockshop",
          "stage": "staging",
          "service": "carts",
          "result": "pass",
          "evaluation": {
            "timeStart": "2020-11-10T11:08:23Z",
            "timeEnd": "2020-11-10T11:10:07Z",
            "result": "pass",
            "score": 100,
            "sloFileContent": "LS0tCnNwZWNfdmVyc2lvbjogIjEuMCIKY29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciCiAgY29tcGFyZV93aXRoOiAic2luZ2xlX3Jlc3VsdCIKICBpbmNsdWRlX3Jlc3VsdF93aXRoX3Njb3JlOiAicGFzcyIKICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxCmZpbHRlcjoKb2JqZWN0aXZlczoKICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1IgogICAga2V5X3NsaTogZmFsc2UKICAgIHBhc3M6ICAgICAgICAgICAgICMgcGFzcyBpZiAocmVsYXRpdmUgY2hhbmdlIDw9IDEwJSBBTkQgYWJzb2x1dGUgdmFsdWUgaXMgPCA2MDBtcykKICAgICAgLSBjcml0ZXJpYToKICAgICAgICAgIC0gIjw9KzEwJSIgICMgcmVsYXRpdmUgdmFsdWVzIHJlcXVpcmUgYSBwcmVmaXhlZCBzaWduIChwbHVzIG9yIG1pbnVzKQogICAgICAgICAgLSAiPDYwMCIgICAgIyBhYnNvbHV0ZSB2YWx1ZXMgb25seSByZXF1aXJlIGEgbG9naWNhbCBvcGVyYXRvcgogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nCiAgICAgIC0gY3JpdGVyaWE6CiAgICAgICAgICAtICI8PTgwMCIKICAgIHdlaWdodDogMQp0b3RhbF9zY29yZToKICBwYXNzOiAiOTAlIgogIHdhcm5pbmc6ICI3NSUi",
            "indicatorResults": [
              {
                "score": 1,
                "value": {
                  "metric": "response_time_p95",
                  "value": 299.18637492576534,
                  "success": true
                },
                "displayName": "",
                "passTargets": null,
                "warningTargets": null,
                "keySli": false,
                "status": "pass"
              }
            ],
            "comparedEvents": [
              "cfa408ce-f552-43c4-aff2-673b1e0548d2"
            ],
            "gitCommit": ""
          }
        },
        "id": "04266cc2-eeea-485b-85b3-f1dea50890ce",
        "source": "lighthouse-service",
        "specversion": "1.0",
        "time": "2020-11-10T11:12:12.364Z",
        "type": "sh.keptn.event.evaluation.finished",
        "shkeptncontext": "24534d07-3a68-4d2b-9a6e-56a1a1a55d2f"
      },
      "heatmapLabel": "2020-11-10 12:12"
    },
    {
      "traces": [],
      "contenttype": "application/json",
      "data": {
        "project": "sockshop",
        "stage": "staging",
        "service": "carts",
        "result": "fail",
        "evaluation": {
          "timeStart": "2020-11-10T11:11:36Z",
          "timeEnd": "2020-11-10T11:13:30Z",
          "result": "fail",
          "score": 50,
          "sloFileContent": "LS0tCnNwZWNfdmVyc2lvbjogIjEuMCIKY29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciCiAgY29tcGFyZV93aXRoOiAic2luZ2xlX3Jlc3VsdCIKICBpbmNsdWRlX3Jlc3VsdF93aXRoX3Njb3JlOiAicGFzcyIKICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxCmZpbHRlcjoKb2JqZWN0aXZlczoKICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1IgogICAga2V5X3NsaTogZmFsc2UKICAgIHBhc3M6ICAgICAgICAgICAgICMgcGFzcyBpZiAocmVsYXRpdmUgY2hhbmdlIDw9IDEwJSBBTkQgYWJzb2x1dGUgdmFsdWUgaXMgPCA2MDBtcykKICAgICAgLSBjcml0ZXJpYToKICAgICAgICAgIC0gIjw9KzEwJSIgICMgcmVsYXRpdmUgdmFsdWVzIHJlcXVpcmUgYSBwcmVmaXhlZCBzaWduIChwbHVzIG9yIG1pbnVzKQogICAgICAgICAgLSAiPDYwMCIgICAgIyBhYnNvbHV0ZSB2YWx1ZXMgb25seSByZXF1aXJlIGEgbG9naWNhbCBvcGVyYXRvcgogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nCiAgICAgIC0gY3JpdGVyaWE6CiAgICAgICAgICAtICI8PTgwMCIKICAgIHdlaWdodDogMQp0b3RhbF9zY29yZToKICBwYXNzOiAiOTAlIgogIHdhcm5pbmc6ICI3NSUi",
          "indicatorResults": [
            {
              "score": 0.5,
              "value": {
                "metric": "response_time_p95",
                "value": 353.3103470363961,
                "success": true
              },
              "displayName": "",
              "passTargets": null,
              "warningTargets": null,
              "keySli": false,
              "status": "warning"
            }
          ],
          "comparedEvents": [
            "04266cc2-eeea-485b-85b3-f1dea50890ce"
          ],
          "gitCommit": ""
        }
      },
      "id": "beefaec1-ecc8-43a0-8d8f-c1cb9b17c18a",
      "source": "lighthouse-service",
      "specversion": "1.0",
      "time": "2020-11-10T11:15:34.488Z",
      "type": "sh.keptn.event.evaluation.finished",
      "shkeptncontext": "448a9f61-c67f-433a-9e4f-fef71bb2a375",
      "plainEvent": {
        "contenttype": "application/json",
        "data": {
          "project": "sockshop",
          "stage": "staging",
          "service": "carts",
          "result": "fail",
          "evaluation": {
            "timeStart": "2020-11-10T11:11:36Z",
            "timeEnd": "2020-11-10T11:13:30Z",
            "result": "fail",
            "score": 50,
            "sloFileContent": "LS0tCnNwZWNfdmVyc2lvbjogIjEuMCIKY29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciCiAgY29tcGFyZV93aXRoOiAic2luZ2xlX3Jlc3VsdCIKICBpbmNsdWRlX3Jlc3VsdF93aXRoX3Njb3JlOiAicGFzcyIKICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxCmZpbHRlcjoKb2JqZWN0aXZlczoKICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1IgogICAga2V5X3NsaTogZmFsc2UKICAgIHBhc3M6ICAgICAgICAgICAgICMgcGFzcyBpZiAocmVsYXRpdmUgY2hhbmdlIDw9IDEwJSBBTkQgYWJzb2x1dGUgdmFsdWUgaXMgPCA2MDBtcykKICAgICAgLSBjcml0ZXJpYToKICAgICAgICAgIC0gIjw9KzEwJSIgICMgcmVsYXRpdmUgdmFsdWVzIHJlcXVpcmUgYSBwcmVmaXhlZCBzaWduIChwbHVzIG9yIG1pbnVzKQogICAgICAgICAgLSAiPDYwMCIgICAgIyBhYnNvbHV0ZSB2YWx1ZXMgb25seSByZXF1aXJlIGEgbG9naWNhbCBvcGVyYXRvcgogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nCiAgICAgIC0gY3JpdGVyaWE6CiAgICAgICAgICAtICI8PTgwMCIKICAgIHdlaWdodDogMQp0b3RhbF9zY29yZToKICBwYXNzOiAiOTAlIgogIHdhcm5pbmc6ICI3NSUi",
            "indicatorResults": [
              {
                "score": 0.5,
                "value": {
                  "metric": "response_time_p95",
                  "value": 353.3103470363961,
                  "success": true
                },
                "displayName": "",
                "passTargets": null,
                "warningTargets": null,
                "keySli": false,
                "status": "warning"
              }
            ],
            "comparedEvents": [
              "04266cc2-eeea-485b-85b3-f1dea50890ce"
            ],
            "gitCommit": ""
          }
        },
        "id": "beefaec1-ecc8-43a0-8d8f-c1cb9b17c18a",
        "source": "lighthouse-service",
        "specversion": "1.0",
        "time": "2020-11-10T11:15:34.488Z",
        "type": "sh.keptn.event.evaluation.finished",
        "shkeptncontext": "448a9f61-c67f-433a-9e4f-fef71bb2a375"
      },
      "heatmapLabel": "2020-11-10 12:15"
    },
    {
      "traces": [],
      "contenttype": "application/json",
      "data": {
        "project": "sockshop",
        "stage": "staging",
        "service": "carts",
        "result": "pass",
        "evaluation": {
          "timeStart": "2020-12-21T11:56:16Z",
          "timeEnd": "2020-12-21T11:58:09Z",
          "result": "pass",
          "score": 100,
          "sloFileContent": "LS0tCnNwZWNfdmVyc2lvbjogIjEuMCIKY29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciCiAgY29tcGFyZV93aXRoOiAic2luZ2xlX3Jlc3VsdCIKICBpbmNsdWRlX3Jlc3VsdF93aXRoX3Njb3JlOiAicGFzcyIKICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxCmZpbHRlcjoKb2JqZWN0aXZlczoKICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1IgogICAga2V5X3NsaTogZmFsc2UKICAgIHBhc3M6ICAgICAgICAgICAgICMgcGFzcyBpZiAocmVsYXRpdmUgY2hhbmdlIDw9IDEwJSBBTkQgYWJzb2x1dGUgdmFsdWUgaXMgPCA2MDBtcykKICAgICAgLSBjcml0ZXJpYToKICAgICAgICAgIC0gIjw9KzEwJSIgICMgcmVsYXRpdmUgdmFsdWVzIHJlcXVpcmUgYSBwcmVmaXhlZCBzaWduIChwbHVzIG9yIG1pbnVzKQogICAgICAgICAgLSAiPDYwMCIgICAgIyBhYnNvbHV0ZSB2YWx1ZXMgb25seSByZXF1aXJlIGEgbG9naWNhbCBvcGVyYXRvcgogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nCiAgICAgIC0gY3JpdGVyaWE6CiAgICAgICAgICAtICI8PTgwMCIKICAgIHdlaWdodDogMQp0b3RhbF9zY29yZToKICBwYXNzOiAiOTAlIgogIHdhcm5pbmc6ICI3NSUi",
          "indicatorResults": [
            {
              "score": 1,
              "value": {
                "metric": "response_time_p95",
                "value": 327.1371879223135,
                "success": true
              },
              "displayName": "",
              "passTargets": null,
              "warningTargets": null,
              "keySli": false,
              "status": "pass"
            }
          ],
          "comparedEvents": [
            "04266cc2-eeea-485b-85b3-f1dea50890ce"
          ],
          "gitCommit": ""
        }
      },
      "id": "99ef68f9-a2a8-426d-81a5-f4ed359506f9",
      "source": "lighthouse-service",
      "specversion": "1.0",
      "time": "2020-12-21T12:00:14.126Z",
      "type": "sh.keptn.event.evaluation.finished",
      "shkeptncontext": "4565c54d-4b7d-4f55-b51e-4797f2e6be05",
      "plainEvent": {
        "contenttype": "application/json",
        "data": {
          "project": "sockshop",
          "stage": "staging",
          "service": "carts",
          "result": "pass",
          "evaluation": {
            "timeStart": "2020-12-21T11:56:16Z",
            "timeEnd": "2020-12-21T11:58:09Z",
            "result": "pass",
            "score": 100,
            "sloFileContent": "LS0tCnNwZWNfdmVyc2lvbjogIjEuMCIKY29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciCiAgY29tcGFyZV93aXRoOiAic2luZ2xlX3Jlc3VsdCIKICBpbmNsdWRlX3Jlc3VsdF93aXRoX3Njb3JlOiAicGFzcyIKICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxCmZpbHRlcjoKb2JqZWN0aXZlczoKICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1IgogICAga2V5X3NsaTogZmFsc2UKICAgIHBhc3M6ICAgICAgICAgICAgICMgcGFzcyBpZiAocmVsYXRpdmUgY2hhbmdlIDw9IDEwJSBBTkQgYWJzb2x1dGUgdmFsdWUgaXMgPCA2MDBtcykKICAgICAgLSBjcml0ZXJpYToKICAgICAgICAgIC0gIjw9KzEwJSIgICMgcmVsYXRpdmUgdmFsdWVzIHJlcXVpcmUgYSBwcmVmaXhlZCBzaWduIChwbHVzIG9yIG1pbnVzKQogICAgICAgICAgLSAiPDYwMCIgICAgIyBhYnNvbHV0ZSB2YWx1ZXMgb25seSByZXF1aXJlIGEgbG9naWNhbCBvcGVyYXRvcgogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nCiAgICAgIC0gY3JpdGVyaWE6CiAgICAgICAgICAtICI8PTgwMCIKICAgIHdlaWdodDogMQp0b3RhbF9zY29yZToKICBwYXNzOiAiOTAlIgogIHdhcm5pbmc6ICI3NSUi",
            "indicatorResults": [
              {
                "score": 1,
                "value": {
                  "metric": "response_time_p95",
                  "value": 327.1371879223135,
                  "success": true
                },
                "displayName": "",
                "passTargets": null,
                "warningTargets": null,
                "keySli": false,
                "status": "pass"
              }
            ],
            "comparedEvents": [
              "04266cc2-eeea-485b-85b3-f1dea50890ce"
            ],
            "gitCommit": ""
          }
        },
        "id": "99ef68f9-a2a8-426d-81a5-f4ed359506f9",
        "source": "lighthouse-service",
        "specversion": "1.0",
        "time": "2020-12-21T12:00:14.126Z",
        "type": "sh.keptn.event.evaluation.finished",
        "shkeptncontext": "4565c54d-4b7d-4f55-b51e-4797f2e6be05"
      },
      "heatmapLabel": "2020-12-21 13:00"
    },
    {
      "traces": [],
      "contenttype": "application/json",
      "data": {
        "project": "sockshop",
        "stage": "staging",
        "service": "carts",
        "result": "pass",
        "evaluation": {
          "timeStart": "2020-12-21T13:09:38Z",
          "timeEnd": "2020-12-21T13:11:37Z",
          "result": "pass",
          "score": 100,
          "sloFileContent": "LS0tCnNwZWNfdmVyc2lvbjogIjEuMCIKY29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciCiAgY29tcGFyZV93aXRoOiAic2luZ2xlX3Jlc3VsdCIKICBpbmNsdWRlX3Jlc3VsdF93aXRoX3Njb3JlOiAicGFzcyIKICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxCmZpbHRlcjoKb2JqZWN0aXZlczoKICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1IgogICAga2V5X3NsaTogZmFsc2UKICAgIHBhc3M6ICAgICAgICAgICAgICMgcGFzcyBpZiAocmVsYXRpdmUgY2hhbmdlIDw9IDEwJSBBTkQgYWJzb2x1dGUgdmFsdWUgaXMgPCA2MDBtcykKICAgICAgLSBjcml0ZXJpYToKICAgICAgICAgIC0gIjw9KzEwJSIgICMgcmVsYXRpdmUgdmFsdWVzIHJlcXVpcmUgYSBwcmVmaXhlZCBzaWduIChwbHVzIG9yIG1pbnVzKQogICAgICAgICAgLSAiPDYwMCIgICAgIyBhYnNvbHV0ZSB2YWx1ZXMgb25seSByZXF1aXJlIGEgbG9naWNhbCBvcGVyYXRvcgogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nCiAgICAgIC0gY3JpdGVyaWE6CiAgICAgICAgICAtICI8PTgwMCIKICAgIHdlaWdodDogMQp0b3RhbF9zY29yZToKICBwYXNzOiAiOTAlIgogIHdhcm5pbmc6ICI3NSUi",
          "indicatorResults": [
            {
              "score": 1,
              "value": {
                "metric": "response_time_p95",
                "value": 348.6222080454321,
                "success": true
              },
              "displayName": "",
              "passTargets": null,
              "warningTargets": null,
              "keySli": false,
              "status": "pass"
            }
          ],
          "comparedEvents": [
            "99ef68f9-a2a8-426d-81a5-f4ed359506f9"
          ],
          "gitCommit": ""
        }
      },
      "id": "a9dcaf6a-f234-4971-80de-c8ad3e045585",
      "source": "lighthouse-service",
      "specversion": "1.0",
      "time": "2020-12-21T13:13:41.949Z",
      "type": "sh.keptn.event.evaluation.finished",
      "shkeptncontext": "f1499961-ef6b-4602-9d89-43adeca5d445",
      "plainEvent": {
        "contenttype": "application/json",
        "data": {
          "project": "sockshop",
          "stage": "staging",
          "service": "carts",
          "result": "pass",
          "evaluation": {
            "timeStart": "2020-12-21T13:09:38Z",
            "timeEnd": "2020-12-21T13:11:37Z",
            "result": "pass",
            "score": 100,
            "sloFileContent": "LS0tCnNwZWNfdmVyc2lvbjogIjEuMCIKY29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciCiAgY29tcGFyZV93aXRoOiAic2luZ2xlX3Jlc3VsdCIKICBpbmNsdWRlX3Jlc3VsdF93aXRoX3Njb3JlOiAicGFzcyIKICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxCmZpbHRlcjoKb2JqZWN0aXZlczoKICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1IgogICAga2V5X3NsaTogZmFsc2UKICAgIHBhc3M6ICAgICAgICAgICAgICMgcGFzcyBpZiAocmVsYXRpdmUgY2hhbmdlIDw9IDEwJSBBTkQgYWJzb2x1dGUgdmFsdWUgaXMgPCA2MDBtcykKICAgICAgLSBjcml0ZXJpYToKICAgICAgICAgIC0gIjw9KzEwJSIgICMgcmVsYXRpdmUgdmFsdWVzIHJlcXVpcmUgYSBwcmVmaXhlZCBzaWduIChwbHVzIG9yIG1pbnVzKQogICAgICAgICAgLSAiPDYwMCIgICAgIyBhYnNvbHV0ZSB2YWx1ZXMgb25seSByZXF1aXJlIGEgbG9naWNhbCBvcGVyYXRvcgogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nCiAgICAgIC0gY3JpdGVyaWE6CiAgICAgICAgICAtICI8PTgwMCIKICAgIHdlaWdodDogMQp0b3RhbF9zY29yZToKICBwYXNzOiAiOTAlIgogIHdhcm5pbmc6ICI3NSUi",
            "indicatorResults": [
              {
                "score": 1,
                "value": {
                  "metric": "response_time_p95",
                  "value": 348.6222080454321,
                  "success": true
                },
                "displayName": "",
                "passTargets": null,
                "warningTargets": null,
                "keySli": false,
                "status": "pass"
              }
            ],
            "comparedEvents": [
              "99ef68f9-a2a8-426d-81a5-f4ed359506f9"
            ],
            "gitCommit": ""
          }
        },
        "id": "a9dcaf6a-f234-4971-80de-c8ad3e045585",
        "source": "lighthouse-service",
        "specversion": "1.0",
        "time": "2020-12-21T13:13:41.949Z",
        "type": "sh.keptn.event.evaluation.finished",
        "shkeptncontext": "f1499961-ef6b-4602-9d89-43adeca5d445"
      },
      "heatmapLabel": "2020-12-21 14:13"
    },
    {
      "traces": [],
      "contenttype": "application/json",
      "data": {
        "project": "sockshop",
        "stage": "staging",
        "service": "carts",
        "result": "pass",
        "evaluation": {
          "timeStart": "2021-02-03T14:23:46Z",
          "timeEnd": "2021-02-03T14:25:34Z",
          "result": "pass",
          "score": 100,
          "sloFileContent": "LS0tCnNwZWNfdmVyc2lvbjogIjEuMCIKY29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciCiAgY29tcGFyZV93aXRoOiAic2luZ2xlX3Jlc3VsdCIKICBpbmNsdWRlX3Jlc3VsdF93aXRoX3Njb3JlOiAicGFzcyIKICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxCmZpbHRlcjoKb2JqZWN0aXZlczoKICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1IgogICAga2V5X3NsaTogZmFsc2UKICAgIHBhc3M6ICAgICAgICAgICAgICMgcGFzcyBpZiAocmVsYXRpdmUgY2hhbmdlIDw9IDEwJSBBTkQgYWJzb2x1dGUgdmFsdWUgaXMgPCA2MDBtcykKICAgICAgLSBjcml0ZXJpYToKICAgICAgICAgIC0gIjw9KzEwJSIgICMgcmVsYXRpdmUgdmFsdWVzIHJlcXVpcmUgYSBwcmVmaXhlZCBzaWduIChwbHVzIG9yIG1pbnVzKQogICAgICAgICAgLSAiPDYwMCIgICAgIyBhYnNvbHV0ZSB2YWx1ZXMgb25seSByZXF1aXJlIGEgbG9naWNhbCBvcGVyYXRvcgogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nCiAgICAgIC0gY3JpdGVyaWE6CiAgICAgICAgICAtICI8PTgwMCIKICAgIHdlaWdodDogMQp0b3RhbF9zY29yZToKICBwYXNzOiAiOTAlIgogIHdhcm5pbmc6ICI3NSUi",
          "indicatorResults": [
            {
              "score": 1,
              "value": {
                "metric": "response_time_p95",
                "value": 339.86815099185304,
                "success": true
              },
              "displayName": "",
              "passTargets": null,
              "warningTargets": null,
              "keySli": false,
              "status": "pass"
            }
          ],
          "comparedEvents": [
            "a9dcaf6a-f234-4971-80de-c8ad3e045585"
          ],
          "gitCommit": ""
        }
      },
      "id": "d5550753-d690-431e-96e0-8b0c3f7f7ca9",
      "source": "lighthouse-service",
      "specversion": "1.0",
      "time": "2021-02-03T14:27:38.964Z",
      "type": "sh.keptn.event.evaluation.finished",
      "shkeptncontext": "13bd94b0-d76b-4d14-8a56-3c10904d2ce8",
      "plainEvent": {
        "contenttype": "application/json",
        "data": {
          "project": "sockshop",
          "stage": "staging",
          "service": "carts",
          "result": "pass",
          "evaluation": {
            "timeStart": "2021-02-03T14:23:46Z",
            "timeEnd": "2021-02-03T14:25:34Z",
            "result": "pass",
            "score": 100,
            "sloFileContent": "LS0tCnNwZWNfdmVyc2lvbjogIjEuMCIKY29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciCiAgY29tcGFyZV93aXRoOiAic2luZ2xlX3Jlc3VsdCIKICBpbmNsdWRlX3Jlc3VsdF93aXRoX3Njb3JlOiAicGFzcyIKICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxCmZpbHRlcjoKb2JqZWN0aXZlczoKICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1IgogICAga2V5X3NsaTogZmFsc2UKICAgIHBhc3M6ICAgICAgICAgICAgICMgcGFzcyBpZiAocmVsYXRpdmUgY2hhbmdlIDw9IDEwJSBBTkQgYWJzb2x1dGUgdmFsdWUgaXMgPCA2MDBtcykKICAgICAgLSBjcml0ZXJpYToKICAgICAgICAgIC0gIjw9KzEwJSIgICMgcmVsYXRpdmUgdmFsdWVzIHJlcXVpcmUgYSBwcmVmaXhlZCBzaWduIChwbHVzIG9yIG1pbnVzKQogICAgICAgICAgLSAiPDYwMCIgICAgIyBhYnNvbHV0ZSB2YWx1ZXMgb25seSByZXF1aXJlIGEgbG9naWNhbCBvcGVyYXRvcgogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nCiAgICAgIC0gY3JpdGVyaWE6CiAgICAgICAgICAtICI8PTgwMCIKICAgIHdlaWdodDogMQp0b3RhbF9zY29yZToKICBwYXNzOiAiOTAlIgogIHdhcm5pbmc6ICI3NSUi",
            "indicatorResults": [
              {
                "score": 1,
                "value": {
                  "metric": "response_time_p95",
                  "value": 339.86815099185304,
                  "success": true
                },
                "displayName": "",
                "passTargets": null,
                "warningTargets": null,
                "keySli": false,
                "status": "pass"
              }
            ],
            "comparedEvents": [
              "a9dcaf6a-f234-4971-80de-c8ad3e045585"
            ],
            "gitCommit": ""
          }
        },
        "id": "d5550753-d690-431e-96e0-8b0c3f7f7ca9",
        "source": "lighthouse-service",
        "specversion": "1.0",
        "time": "2021-02-03T14:27:38.964Z",
        "type": "sh.keptn.event.evaluation.finished",
        "shkeptncontext": "13bd94b0-d76b-4d14-8a56-3c10904d2ce8"
      },
      "heatmapLabel": "2021-02-03 15:27"
    },
    {
      "traces": [],
      "data": {
        "evaluation": {
          "comparedEvents": [
            "d5550753-d690-431e-96e0-8b0c3f7f7ca9"
          ],
          "gitCommit": "",
          "indicatorResults": [
            {
              "keySli": false,
              "passTargets": [
                {
                  "criteria": "<=+10%",
                  "targetValue": 373.85496609103836,
                  "violated": false
                },
                {
                  "criteria": "<600",
                  "targetValue": 600,
                  "violated": false
                }
              ],
              "score": 1,
              "status": "pass",
              "value": {
                "metric": "response_time_p95",
                "success": true,
                "value": 370.20610238347683
              },
              "warningTargets": [
                {
                  "criteria": "<=800",
                  "targetValue": 800,
                  "violated": false
                }
              ]
            }
          ],
          "result": "pass",
          "score": 100,
          "sloFileContent": "LS0tCnNwZWNfdmVyc2lvbjogIjEuMCIKY29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciCiAgY29tcGFyZV93aXRoOiAic2luZ2xlX3Jlc3VsdCIKICBpbmNsdWRlX3Jlc3VsdF93aXRoX3Njb3JlOiAicGFzcyIKICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxCmZpbHRlcjoKb2JqZWN0aXZlczoKICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1IgogICAga2V5X3NsaTogZmFsc2UKICAgIHBhc3M6ICAgICAgICAgICAgICMgcGFzcyBpZiAocmVsYXRpdmUgY2hhbmdlIDw9IDEwJSBBTkQgYWJzb2x1dGUgdmFsdWUgaXMgPCA2MDBtcykKICAgICAgLSBjcml0ZXJpYToKICAgICAgICAgIC0gIjw9KzEwJSIgICMgcmVsYXRpdmUgdmFsdWVzIHJlcXVpcmUgYSBwcmVmaXhlZCBzaWduIChwbHVzIG9yIG1pbnVzKQogICAgICAgICAgLSAiPDYwMCIgICAgIyBhYnNvbHV0ZSB2YWx1ZXMgb25seSByZXF1aXJlIGEgbG9naWNhbCBvcGVyYXRvcgogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nCiAgICAgIC0gY3JpdGVyaWE6CiAgICAgICAgICAtICI8PTgwMCIKICAgIHdlaWdodDogMQp0b3RhbF9zY29yZToKICBwYXNzOiAiOTAlIgogIHdhcm5pbmc6ICI3NSUi",
          "timeEnd": "2021-03-15T13:07:50Z",
          "timeStart": "2021-03-15T13:05:47Z"
        },
        "labels": {
          "DtCreds": "dynatrace"
        },
        "project": "sockshop",
        "result": "pass",
        "service": "carts",
        "stage": "staging",
        "status": "succeeded"
      },
      "id": "e0b793e2-0e14-467c-b31d-3d08156651f2",
      "source": "lighthouse-service",
      "specversion": "1.0",
      "time": "2021-03-15T13:08:55.844Z",
      "type": "sh.keptn.event.evaluation.finished",
      "shkeptncontext": "030bf0a4-9975-404e-9b50-7f07d7eefb95",
      "triggeredid": "f841c2a3-d543-4845-9379-94c8d61cf22c",
      "plainEvent": {
        "data": {
          "evaluation": {
            "comparedEvents": [
              "d5550753-d690-431e-96e0-8b0c3f7f7ca9"
            ],
            "gitCommit": "",
            "indicatorResults": [
              {
                "keySli": false,
                "passTargets": [
                  {
                    "criteria": "<=+10%",
                    "targetValue": 373.85496609103836,
                    "violated": false
                  },
                  {
                    "criteria": "<600",
                    "targetValue": 600,
                    "violated": false
                  }
                ],
                "score": 1,
                "status": "pass",
                "value": {
                  "metric": "response_time_p95",
                  "success": true,
                  "value": 370.20610238347683
                },
                "warningTargets": [
                  {
                    "criteria": "<=800",
                    "targetValue": 800,
                    "violated": false
                  }
                ]
              }
            ],
            "result": "pass",
            "score": 100,
            "sloFileContent": "LS0tCnNwZWNfdmVyc2lvbjogIjEuMCIKY29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciCiAgY29tcGFyZV93aXRoOiAic2luZ2xlX3Jlc3VsdCIKICBpbmNsdWRlX3Jlc3VsdF93aXRoX3Njb3JlOiAicGFzcyIKICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxCmZpbHRlcjoKb2JqZWN0aXZlczoKICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1IgogICAga2V5X3NsaTogZmFsc2UKICAgIHBhc3M6ICAgICAgICAgICAgICMgcGFzcyBpZiAocmVsYXRpdmUgY2hhbmdlIDw9IDEwJSBBTkQgYWJzb2x1dGUgdmFsdWUgaXMgPCA2MDBtcykKICAgICAgLSBjcml0ZXJpYToKICAgICAgICAgIC0gIjw9KzEwJSIgICMgcmVsYXRpdmUgdmFsdWVzIHJlcXVpcmUgYSBwcmVmaXhlZCBzaWduIChwbHVzIG9yIG1pbnVzKQogICAgICAgICAgLSAiPDYwMCIgICAgIyBhYnNvbHV0ZSB2YWx1ZXMgb25seSByZXF1aXJlIGEgbG9naWNhbCBvcGVyYXRvcgogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nCiAgICAgIC0gY3JpdGVyaWE6CiAgICAgICAgICAtICI8PTgwMCIKICAgIHdlaWdodDogMQp0b3RhbF9zY29yZToKICBwYXNzOiAiOTAlIgogIHdhcm5pbmc6ICI3NSUi",
            "timeEnd": "2021-03-15T13:07:50Z",
            "timeStart": "2021-03-15T13:05:47Z"
          },
          "labels": {
            "DtCreds": "dynatrace"
          },
          "project": "sockshop",
          "result": "pass",
          "service": "carts",
          "stage": "staging",
          "status": "succeeded"
        },
        "id": "e0b793e2-0e14-467c-b31d-3d08156651f2",
        "source": "lighthouse-service",
        "specversion": "1.0",
        "time": "2021-03-15T13:08:55.844Z",
        "type": "sh.keptn.event.evaluation.finished",
        "shkeptncontext": "030bf0a4-9975-404e-9b50-7f07d7eefb95",
        "triggeredid": "f841c2a3-d543-4845-9379-94c8d61cf22c"
      },
      "heatmapLabel": "2021-03-15 14:08"
    },
    {
      "traces": [],
      "data": {
        "evaluation": {
          "comparedEvents": [
            "e0b793e2-0e14-467c-b31d-3d08156651f2"
          ],
          "gitCommit": "",
          "indicatorResults": [
            {
              "keySli": false,
              "passTargets": [
                {
                  "criteria": "<=+10%",
                  "targetValue": 407.2267126218245,
                  "violated": true
                },
                {
                  "criteria": "<600",
                  "targetValue": 600,
                  "violated": true
                }
              ],
              "score": 0,
              "status": "fail",
              "value": {
                "metric": "response_time_p95",
                "success": true,
                "value": 1082.6252505468217
              },
              "warningTargets": [
                {
                  "criteria": "<=800",
                  "targetValue": 800,
                  "violated": true
                }
              ]
            }
          ],
          "result": "fail",
          "score": 0,
          "sloFileContent": "LS0tCnNwZWNfdmVyc2lvbjogIjEuMCIKY29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciCiAgY29tcGFyZV93aXRoOiAic2luZ2xlX3Jlc3VsdCIKICBpbmNsdWRlX3Jlc3VsdF93aXRoX3Njb3JlOiAicGFzcyIKICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxCmZpbHRlcjoKb2JqZWN0aXZlczoKICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1IgogICAga2V5X3NsaTogZmFsc2UKICAgIHBhc3M6ICAgICAgICAgICAgICMgcGFzcyBpZiAocmVsYXRpdmUgY2hhbmdlIDw9IDEwJSBBTkQgYWJzb2x1dGUgdmFsdWUgaXMgPCA2MDBtcykKICAgICAgLSBjcml0ZXJpYToKICAgICAgICAgIC0gIjw9KzEwJSIgICMgcmVsYXRpdmUgdmFsdWVzIHJlcXVpcmUgYSBwcmVmaXhlZCBzaWduIChwbHVzIG9yIG1pbnVzKQogICAgICAgICAgLSAiPDYwMCIgICAgIyBhYnNvbHV0ZSB2YWx1ZXMgb25seSByZXF1aXJlIGEgbG9naWNhbCBvcGVyYXRvcgogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nCiAgICAgIC0gY3JpdGVyaWE6CiAgICAgICAgICAtICI8PTgwMCIKICAgIHdlaWdodDogMQp0b3RhbF9zY29yZToKICBwYXNzOiAiOTAlIgogIHdhcm5pbmc6ICI3NSUi",
          "timeEnd": "2021-03-16T06:37:25Z",
          "timeStart": "2021-03-16T06:28:15Z"
        },
        "labels": {
          "DtCreds": "dynatrace"
        },
        "project": "sockshop",
        "result": "fail",
        "service": "carts",
        "stage": "staging",
        "status": "succeeded"
      },
      "id": "f5896d5d-3d98-431b-b8a0-202cd6e0657d",
      "source": "lighthouse-service",
      "specversion": "1.0",
      "time": "2021-03-16T06:37:32.967Z",
      "type": "sh.keptn.event.evaluation.finished",
      "shkeptncontext": "bb648dd6-f8cf-481f-85e4-fc0f153d5919",
      "triggeredid": "6364f4ef-4a1b-46d9-8d05-dd5cf1d6f164",
      "plainEvent": {
        "data": {
          "evaluation": {
            "comparedEvents": [
              "e0b793e2-0e14-467c-b31d-3d08156651f2"
            ],
            "gitCommit": "",
            "indicatorResults": [
              {
                "keySli": false,
                "passTargets": [
                  {
                    "criteria": "<=+10%",
                    "targetValue": 407.2267126218245,
                    "violated": true
                  },
                  {
                    "criteria": "<600",
                    "targetValue": 600,
                    "violated": true
                  }
                ],
                "score": 0,
                "status": "fail",
                "value": {
                  "metric": "response_time_p95",
                  "success": true,
                  "value": 1082.6252505468217
                },
                "warningTargets": [
                  {
                    "criteria": "<=800",
                    "targetValue": 800,
                    "violated": true
                  }
                ]
              }
            ],
            "result": "fail",
            "score": 0,
            "sloFileContent": "LS0tCnNwZWNfdmVyc2lvbjogIjEuMCIKY29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciCiAgY29tcGFyZV93aXRoOiAic2luZ2xlX3Jlc3VsdCIKICBpbmNsdWRlX3Jlc3VsdF93aXRoX3Njb3JlOiAicGFzcyIKICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxCmZpbHRlcjoKb2JqZWN0aXZlczoKICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1IgogICAga2V5X3NsaTogZmFsc2UKICAgIHBhc3M6ICAgICAgICAgICAgICMgcGFzcyBpZiAocmVsYXRpdmUgY2hhbmdlIDw9IDEwJSBBTkQgYWJzb2x1dGUgdmFsdWUgaXMgPCA2MDBtcykKICAgICAgLSBjcml0ZXJpYToKICAgICAgICAgIC0gIjw9KzEwJSIgICMgcmVsYXRpdmUgdmFsdWVzIHJlcXVpcmUgYSBwcmVmaXhlZCBzaWduIChwbHVzIG9yIG1pbnVzKQogICAgICAgICAgLSAiPDYwMCIgICAgIyBhYnNvbHV0ZSB2YWx1ZXMgb25seSByZXF1aXJlIGEgbG9naWNhbCBvcGVyYXRvcgogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nCiAgICAgIC0gY3JpdGVyaWE6CiAgICAgICAgICAtICI8PTgwMCIKICAgIHdlaWdodDogMQp0b3RhbF9zY29yZToKICBwYXNzOiAiOTAlIgogIHdhcm5pbmc6ICI3NSUi",
            "timeEnd": "2021-03-16T06:37:25Z",
            "timeStart": "2021-03-16T06:28:15Z"
          },
          "labels": {
            "DtCreds": "dynatrace"
          },
          "project": "sockshop",
          "result": "fail",
          "service": "carts",
          "stage": "staging",
          "status": "succeeded"
        },
        "id": "f5896d5d-3d98-431b-b8a0-202cd6e0657d",
        "source": "lighthouse-service",
        "specversion": "1.0",
        "time": "2021-03-16T06:37:32.967Z",
        "type": "sh.keptn.event.evaluation.finished",
        "shkeptncontext": "bb648dd6-f8cf-481f-85e4-fc0f153d5919",
        "triggeredid": "6364f4ef-4a1b-46d9-8d05-dd5cf1d6f164"
      },
      "heatmapLabel": "2021-03-16 07:37"
    },
    {
      "traces": [],
      "data": {
        "evaluation": {
          "comparedEvents": [
            "e0b793e2-0e14-467c-b31d-3d08156651f2"
          ],
          "gitCommit": "",
          "indicatorResults": [
            {
              "keySli": false,
              "passTargets": [
                {
                  "criteria": "<=+10%",
                  "targetValue": 407.2267126218245,
                  "violated": false
                },
                {
                  "criteria": "<600",
                  "targetValue": 600,
                  "violated": false
                }
              ],
              "score": 1,
              "status": "pass",
              "value": {
                "metric": "response_time_p95",
                "success": true,
                "value": 334.5522038804587
              },
              "warningTargets": [
                {
                  "criteria": "<=800",
                  "targetValue": 800,
                  "violated": false
                }
              ]
            }
          ],
          "result": "pass",
          "score": 100,
          "sloFileContent": "LS0tCnNwZWNfdmVyc2lvbjogIjEuMCIKY29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciCiAgY29tcGFyZV93aXRoOiAic2luZ2xlX3Jlc3VsdCIKICBpbmNsdWRlX3Jlc3VsdF93aXRoX3Njb3JlOiAicGFzcyIKICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxCmZpbHRlcjoKb2JqZWN0aXZlczoKICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1IgogICAga2V5X3NsaTogZmFsc2UKICAgIHBhc3M6ICAgICAgICAgICAgICMgcGFzcyBpZiAocmVsYXRpdmUgY2hhbmdlIDw9IDEwJSBBTkQgYWJzb2x1dGUgdmFsdWUgaXMgPCA2MDBtcykKICAgICAgLSBjcml0ZXJpYToKICAgICAgICAgIC0gIjw9KzEwJSIgICMgcmVsYXRpdmUgdmFsdWVzIHJlcXVpcmUgYSBwcmVmaXhlZCBzaWduIChwbHVzIG9yIG1pbnVzKQogICAgICAgICAgLSAiPDYwMCIgICAgIyBhYnNvbHV0ZSB2YWx1ZXMgb25seSByZXF1aXJlIGEgbG9naWNhbCBvcGVyYXRvcgogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nCiAgICAgIC0gY3JpdGVyaWE6CiAgICAgICAgICAtICI8PTgwMCIKICAgIHdlaWdodDogMQp0b3RhbF9zY29yZToKICBwYXNzOiAiOTAlIgogIHdhcm5pbmc6ICI3NSUi",
          "timeEnd": "2021-03-29T10:15:58Z",
          "timeStart": "2021-03-29T10:13:52Z"
        },
        "labels": {
          "DtCreds": "dynatrace"
        },
        "project": "sockshop",
        "result": "pass",
        "service": "carts",
        "stage": "staging",
        "status": "succeeded"
      },
      "id": "acb7af04-6dea-4aa5-a6d8-e502d627eaaf",
      "source": "lighthouse-service",
      "specversion": "1.0",
      "time": "2021-03-29T10:17:06.322Z",
      "type": "sh.keptn.event.evaluation.finished",
      "shkeptncontext": "c65d28ef-2e65-4921-842e-07217dde2675",
      "triggeredid": "460a018a-a694-41e2-9d71-b43e891175b6",
      "plainEvent": {
        "data": {
          "evaluation": {
            "comparedEvents": [
              "e0b793e2-0e14-467c-b31d-3d08156651f2"
            ],
            "gitCommit": "",
            "indicatorResults": [
              {
                "keySli": false,
                "passTargets": [
                  {
                    "criteria": "<=+10%",
                    "targetValue": 407.2267126218245,
                    "violated": false
                  },
                  {
                    "criteria": "<600",
                    "targetValue": 600,
                    "violated": false
                  }
                ],
                "score": 1,
                "status": "pass",
                "value": {
                  "metric": "response_time_p95",
                  "success": true,
                  "value": 334.5522038804587
                },
                "warningTargets": [
                  {
                    "criteria": "<=800",
                    "targetValue": 800,
                    "violated": false
                  }
                ]
              }
            ],
            "result": "pass",
            "score": 100,
            "sloFileContent": "LS0tCnNwZWNfdmVyc2lvbjogIjEuMCIKY29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciCiAgY29tcGFyZV93aXRoOiAic2luZ2xlX3Jlc3VsdCIKICBpbmNsdWRlX3Jlc3VsdF93aXRoX3Njb3JlOiAicGFzcyIKICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxCmZpbHRlcjoKb2JqZWN0aXZlczoKICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1IgogICAga2V5X3NsaTogZmFsc2UKICAgIHBhc3M6ICAgICAgICAgICAgICMgcGFzcyBpZiAocmVsYXRpdmUgY2hhbmdlIDw9IDEwJSBBTkQgYWJzb2x1dGUgdmFsdWUgaXMgPCA2MDBtcykKICAgICAgLSBjcml0ZXJpYToKICAgICAgICAgIC0gIjw9KzEwJSIgICMgcmVsYXRpdmUgdmFsdWVzIHJlcXVpcmUgYSBwcmVmaXhlZCBzaWduIChwbHVzIG9yIG1pbnVzKQogICAgICAgICAgLSAiPDYwMCIgICAgIyBhYnNvbHV0ZSB2YWx1ZXMgb25seSByZXF1aXJlIGEgbG9naWNhbCBvcGVyYXRvcgogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nCiAgICAgIC0gY3JpdGVyaWE6CiAgICAgICAgICAtICI8PTgwMCIKICAgIHdlaWdodDogMQp0b3RhbF9zY29yZToKICBwYXNzOiAiOTAlIgogIHdhcm5pbmc6ICI3NSUi",
            "timeEnd": "2021-03-29T10:15:58Z",
            "timeStart": "2021-03-29T10:13:52Z"
          },
          "labels": {
            "DtCreds": "dynatrace"
          },
          "project": "sockshop",
          "result": "pass",
          "service": "carts",
          "stage": "staging",
          "status": "succeeded"
        },
        "id": "acb7af04-6dea-4aa5-a6d8-e502d627eaaf",
        "source": "lighthouse-service",
        "specversion": "1.0",
        "time": "2021-03-29T10:17:06.322Z",
        "type": "sh.keptn.event.evaluation.finished",
        "shkeptncontext": "c65d28ef-2e65-4921-842e-07217dde2675",
        "triggeredid": "460a018a-a694-41e2-9d71-b43e891175b6"
      },
      "heatmapLabel": "2021-03-29 12:17"
    },
    {
      "traces": [],
      "data": {
        "evaluation": {
          "comparedEvents": [
            "acb7af04-6dea-4aa5-a6d8-e502d627eaaf"
          ],
          "gitCommit": "",
          "indicatorResults": [
            {
              "displayName": "",
              "keySli": false,
              "passTargets": [
                {
                  "criteria": "<=+10%",
                  "targetValue": 368.00742426850456,
                  "violated": true
                },
                {
                  "criteria": "<600",
                  "targetValue": 600,
                  "violated": true
                }
              ],
              "score": 0,
              "status": "fail",
              "value": {
                "metric": "response_time_p95",
                "success": true,
                "value": 1077.5072951443556
              },
              "warningTargets": [
                {
                  "criteria": "<=800",
                  "targetValue": 800,
                  "violated": true
                }
              ]
            }
          ],
          "result": "fail",
          "score": 0,
          "sloFileContent": "LS0tCnNwZWNfdmVyc2lvbjogIjEuMCIKY29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciCiAgY29tcGFyZV93aXRoOiAic2luZ2xlX3Jlc3VsdCIKICBpbmNsdWRlX3Jlc3VsdF93aXRoX3Njb3JlOiAicGFzcyIKICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxCmZpbHRlcjoKb2JqZWN0aXZlczoKICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1IgogICAga2V5X3NsaTogZmFsc2UKICAgIHBhc3M6ICAgICAgICAgICAgICMgcGFzcyBpZiAocmVsYXRpdmUgY2hhbmdlIDw9IDEwJSBBTkQgYWJzb2x1dGUgdmFsdWUgaXMgPCA2MDBtcykKICAgICAgLSBjcml0ZXJpYToKICAgICAgICAgIC0gIjw9KzEwJSIgICMgcmVsYXRpdmUgdmFsdWVzIHJlcXVpcmUgYSBwcmVmaXhlZCBzaWduIChwbHVzIG9yIG1pbnVzKQogICAgICAgICAgLSAiPDYwMCIgICAgIyBhYnNvbHV0ZSB2YWx1ZXMgb25seSByZXF1aXJlIGEgbG9naWNhbCBvcGVyYXRvcgogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nCiAgICAgIC0gY3JpdGVyaWE6CiAgICAgICAgICAtICI8PTgwMCIKICAgIHdlaWdodDogMQp0b3RhbF9zY29yZToKICBwYXNzOiAiOTAlIgogIHdhcm5pbmc6ICI3NSUi",
          "timeEnd": "2021-03-29T14:08:31Z",
          "timeStart": "2021-03-29T13:59:41Z"
        },
        "labels": {
          "DtCreds": "dynatrace"
        },
        "project": "sockshop",
        "result": "fail",
        "service": "carts",
        "stage": "staging",
        "status": "succeeded"
      },
      "id": "52fb89c5-027b-4deb-bdf7-6b01e222907f",
      "source": "lighthouse-service",
      "specversion": "1.0",
      "time": "2021-03-29T14:08:39.265Z",
      "type": "sh.keptn.event.evaluation.finished",
      "shkeptncontext": "e5d78b1e-c393-4d03-8975-d330c801c63b",
      "triggeredid": "a33ac1ee-baf1-4745-bae4-2d31267a3ae4",
      "plainEvent": {
        "data": {
          "evaluation": {
            "comparedEvents": [
              "acb7af04-6dea-4aa5-a6d8-e502d627eaaf"
            ],
            "gitCommit": "",
            "indicatorResults": [
              {
                "displayName": "",
                "keySli": false,
                "passTargets": [
                  {
                    "criteria": "<=+10%",
                    "targetValue": 368.00742426850456,
                    "violated": true
                  },
                  {
                    "criteria": "<600",
                    "targetValue": 600,
                    "violated": true
                  }
                ],
                "score": 0,
                "status": "fail",
                "value": {
                  "metric": "response_time_p95",
                  "success": true,
                  "value": 1077.5072951443556
                },
                "warningTargets": [
                  {
                    "criteria": "<=800",
                    "targetValue": 800,
                    "violated": true
                  }
                ]
              }
            ],
            "result": "fail",
            "score": 0,
            "sloFileContent": "LS0tCnNwZWNfdmVyc2lvbjogIjEuMCIKY29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciCiAgY29tcGFyZV93aXRoOiAic2luZ2xlX3Jlc3VsdCIKICBpbmNsdWRlX3Jlc3VsdF93aXRoX3Njb3JlOiAicGFzcyIKICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxCmZpbHRlcjoKb2JqZWN0aXZlczoKICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1IgogICAga2V5X3NsaTogZmFsc2UKICAgIHBhc3M6ICAgICAgICAgICAgICMgcGFzcyBpZiAocmVsYXRpdmUgY2hhbmdlIDw9IDEwJSBBTkQgYWJzb2x1dGUgdmFsdWUgaXMgPCA2MDBtcykKICAgICAgLSBjcml0ZXJpYToKICAgICAgICAgIC0gIjw9KzEwJSIgICMgcmVsYXRpdmUgdmFsdWVzIHJlcXVpcmUgYSBwcmVmaXhlZCBzaWduIChwbHVzIG9yIG1pbnVzKQogICAgICAgICAgLSAiPDYwMCIgICAgIyBhYnNvbHV0ZSB2YWx1ZXMgb25seSByZXF1aXJlIGEgbG9naWNhbCBvcGVyYXRvcgogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nCiAgICAgIC0gY3JpdGVyaWE6CiAgICAgICAgICAtICI8PTgwMCIKICAgIHdlaWdodDogMQp0b3RhbF9zY29yZToKICBwYXNzOiAiOTAlIgogIHdhcm5pbmc6ICI3NSUi",
            "timeEnd": "2021-03-29T14:08:31Z",
            "timeStart": "2021-03-29T13:59:41Z"
          },
          "labels": {
            "DtCreds": "dynatrace"
          },
          "project": "sockshop",
          "result": "fail",
          "service": "carts",
          "stage": "staging",
          "status": "succeeded"
        },
        "id": "52fb89c5-027b-4deb-bdf7-6b01e222907f",
        "source": "lighthouse-service",
        "specversion": "1.0",
        "time": "2021-03-29T14:08:39.265Z",
        "type": "sh.keptn.event.evaluation.finished",
        "shkeptncontext": "e5d78b1e-c393-4d03-8975-d330c801c63b",
        "triggeredid": "a33ac1ee-baf1-4745-bae4-2d31267a3ae4"
      },
      "heatmapLabel": "2021-03-29 16:08"
    },
    {
      "traces": [],
      "data": {
        "evaluation": {
          "comparedEvents": [
            "acb7af04-6dea-4aa5-a6d8-e502d627eaaf"
          ],
          "gitCommit": "",
          "indicatorResults": [
            {
              "displayName": "",
              "keySli": false,
              "passTargets": [
                {
                  "criteria": "<=+10%",
                  "targetValue": 368.00742426850456,
                  "violated": false
                },
                {
                  "criteria": "<600",
                  "targetValue": 600,
                  "violated": false
                }
              ],
              "score": 1,
              "status": "pass",
              "value": {
                "metric": "response_time_p95",
                "success": true,
                "value": 365.2828487064428
              },
              "warningTargets": [
                {
                  "criteria": "<=800",
                  "targetValue": 800,
                  "violated": false
                }
              ]
            }
          ],
          "result": "pass",
          "score": 100,
          "sloFileContent": "LS0tCnNwZWNfdmVyc2lvbjogIjEuMCIKY29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciCiAgY29tcGFyZV93aXRoOiAic2luZ2xlX3Jlc3VsdCIKICBpbmNsdWRlX3Jlc3VsdF93aXRoX3Njb3JlOiAicGFzcyIKICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxCmZpbHRlcjoKb2JqZWN0aXZlczoKICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1IgogICAga2V5X3NsaTogZmFsc2UKICAgIHBhc3M6ICAgICAgICAgICAgICMgcGFzcyBpZiAocmVsYXRpdmUgY2hhbmdlIDw9IDEwJSBBTkQgYWJzb2x1dGUgdmFsdWUgaXMgPCA2MDBtcykKICAgICAgLSBjcml0ZXJpYToKICAgICAgICAgIC0gIjw9KzEwJSIgICMgcmVsYXRpdmUgdmFsdWVzIHJlcXVpcmUgYSBwcmVmaXhlZCBzaWduIChwbHVzIG9yIG1pbnVzKQogICAgICAgICAgLSAiPDYwMCIgICAgIyBhYnNvbHV0ZSB2YWx1ZXMgb25seSByZXF1aXJlIGEgbG9naWNhbCBvcGVyYXRvcgogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nCiAgICAgIC0gY3JpdGVyaWE6CiAgICAgICAgICAtICI8PTgwMCIKICAgIHdlaWdodDogMQp0b3RhbF9zY29yZToKICBwYXNzOiAiOTAlIgogIHdhcm5pbmc6ICI3NSUi",
          "timeEnd": "2021-03-29T15:28:13Z",
          "timeStart": "2021-03-29T15:26:11Z"
        },
        "labels": {
          "DtCreds": "dynatrace"
        },
        "project": "sockshop",
        "result": "pass",
        "service": "carts",
        "stage": "staging",
        "status": "succeeded"
      },
      "id": "8ca4250c-5442-45cb-b745-4b00f2284ea0",
      "source": "lighthouse-service",
      "specversion": "1.0",
      "time": "2021-03-29T15:29:20.851Z",
      "type": "sh.keptn.event.evaluation.finished",
      "shkeptncontext": "03efe46f-f640-4547-8136-96d6c066343c",
      "triggeredid": "91c3a353-c06b-4a55-bf78-16e6184d1c12",
      "plainEvent": {
        "data": {
          "evaluation": {
            "comparedEvents": [
              "acb7af04-6dea-4aa5-a6d8-e502d627eaaf"
            ],
            "gitCommit": "",
            "indicatorResults": [
              {
                "displayName": "",
                "keySli": false,
                "passTargets": [
                  {
                    "criteria": "<=+10%",
                    "targetValue": 368.00742426850456,
                    "violated": false
                  },
                  {
                    "criteria": "<600",
                    "targetValue": 600,
                    "violated": false
                  }
                ],
                "score": 1,
                "status": "pass",
                "value": {
                  "metric": "response_time_p95",
                  "success": true,
                  "value": 365.2828487064428
                },
                "warningTargets": [
                  {
                    "criteria": "<=800",
                    "targetValue": 800,
                    "violated": false
                  }
                ]
              }
            ],
            "result": "pass",
            "score": 100,
            "sloFileContent": "LS0tCnNwZWNfdmVyc2lvbjogIjEuMCIKY29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciCiAgY29tcGFyZV93aXRoOiAic2luZ2xlX3Jlc3VsdCIKICBpbmNsdWRlX3Jlc3VsdF93aXRoX3Njb3JlOiAicGFzcyIKICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxCmZpbHRlcjoKb2JqZWN0aXZlczoKICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1IgogICAga2V5X3NsaTogZmFsc2UKICAgIHBhc3M6ICAgICAgICAgICAgICMgcGFzcyBpZiAocmVsYXRpdmUgY2hhbmdlIDw9IDEwJSBBTkQgYWJzb2x1dGUgdmFsdWUgaXMgPCA2MDBtcykKICAgICAgLSBjcml0ZXJpYToKICAgICAgICAgIC0gIjw9KzEwJSIgICMgcmVsYXRpdmUgdmFsdWVzIHJlcXVpcmUgYSBwcmVmaXhlZCBzaWduIChwbHVzIG9yIG1pbnVzKQogICAgICAgICAgLSAiPDYwMCIgICAgIyBhYnNvbHV0ZSB2YWx1ZXMgb25seSByZXF1aXJlIGEgbG9naWNhbCBvcGVyYXRvcgogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nCiAgICAgIC0gY3JpdGVyaWE6CiAgICAgICAgICAtICI8PTgwMCIKICAgIHdlaWdodDogMQp0b3RhbF9zY29yZToKICBwYXNzOiAiOTAlIgogIHdhcm5pbmc6ICI3NSUi",
            "timeEnd": "2021-03-29T15:28:13Z",
            "timeStart": "2021-03-29T15:26:11Z"
          },
          "labels": {
            "DtCreds": "dynatrace"
          },
          "project": "sockshop",
          "result": "pass",
          "service": "carts",
          "stage": "staging",
          "status": "succeeded"
        },
        "id": "8ca4250c-5442-45cb-b745-4b00f2284ea0",
        "source": "lighthouse-service",
        "specversion": "1.0",
        "time": "2021-03-29T15:29:20.851Z",
        "type": "sh.keptn.event.evaluation.finished",
        "shkeptncontext": "03efe46f-f640-4547-8136-96d6c066343c",
        "triggeredid": "91c3a353-c06b-4a55-bf78-16e6184d1c12"
      },
      "heatmapLabel": "2021-03-29 17:29"
    },
    {
      "traces": [],
      "data": {
        "evaluation": {
          "comparedEvents": [
            "8ca4250c-5442-45cb-b745-4b00f2284ea0"
          ],
          "gitCommit": "",
          "indicatorResults": [
            {
              "displayName": "",
              "keySli": false,
              "passTargets": [
                {
                  "criteria": "<=+10%",
                  "targetValue": 401.8111335770871,
                  "violated": false
                },
                {
                  "criteria": "<600",
                  "targetValue": 600,
                  "violated": false
                }
              ],
              "score": 1,
              "status": "pass",
              "value": {
                "metric": "response_time_p95",
                "success": true,
                "value": 356.90404128243733
              },
              "warningTargets": [
                {
                  "criteria": "<=800",
                  "targetValue": 800,
                  "violated": false
                }
              ]
            }
          ],
          "result": "pass",
          "score": 100,
          "sloFileContent": "LS0tCnNwZWNfdmVyc2lvbjogIjEuMCIKY29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciCiAgY29tcGFyZV93aXRoOiAic2luZ2xlX3Jlc3VsdCIKICBpbmNsdWRlX3Jlc3VsdF93aXRoX3Njb3JlOiAicGFzcyIKICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxCmZpbHRlcjoKb2JqZWN0aXZlczoKICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1IgogICAga2V5X3NsaTogZmFsc2UKICAgIHBhc3M6ICAgICAgICAgICAgICMgcGFzcyBpZiAocmVsYXRpdmUgY2hhbmdlIDw9IDEwJSBBTkQgYWJzb2x1dGUgdmFsdWUgaXMgPCA2MDBtcykKICAgICAgLSBjcml0ZXJpYToKICAgICAgICAgIC0gIjw9KzEwJSIgICMgcmVsYXRpdmUgdmFsdWVzIHJlcXVpcmUgYSBwcmVmaXhlZCBzaWduIChwbHVzIG9yIG1pbnVzKQogICAgICAgICAgLSAiPDYwMCIgICAgIyBhYnNvbHV0ZSB2YWx1ZXMgb25seSByZXF1aXJlIGEgbG9naWNhbCBvcGVyYXRvcgogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nCiAgICAgIC0gY3JpdGVyaWE6CiAgICAgICAgICAtICI8PTgwMCIKICAgIHdlaWdodDogMQp0b3RhbF9zY29yZToKICBwYXNzOiAiOTAlIgogIHdhcm5pbmc6ICI3NSUi",
          "timeEnd": "2021-04-08T15:23:59Z",
          "timeStart": "2021-04-08T15:22:00Z"
        },
        "labels": {
          "DtCreds": "dynatrace"
        },
        "project": "sockshop",
        "result": "pass",
        "service": "carts",
        "stage": "staging",
        "status": "succeeded"
      },
      "id": "c38712d2-0e2f-4412-bd39-fed91f333510",
      "source": "lighthouse-service",
      "specversion": "1.0",
      "time": "2021-04-08T15:26:07.223Z",
      "type": "sh.keptn.event.evaluation.finished",
      "shkeptncontext": "777b25e4-4f47-4b43-a4d4-88c87037df4c",
      "triggeredid": "32f1be11-5936-447f-8f62-56429da702bc",
      "plainEvent": {
        "data": {
          "evaluation": {
            "comparedEvents": [
              "8ca4250c-5442-45cb-b745-4b00f2284ea0"
            ],
            "gitCommit": "",
            "indicatorResults": [
              {
                "displayName": "",
                "keySli": false,
                "passTargets": [
                  {
                    "criteria": "<=+10%",
                    "targetValue": 401.8111335770871,
                    "violated": false
                  },
                  {
                    "criteria": "<600",
                    "targetValue": 600,
                    "violated": false
                  }
                ],
                "score": 1,
                "status": "pass",
                "value": {
                  "metric": "response_time_p95",
                  "success": true,
                  "value": 356.90404128243733
                },
                "warningTargets": [
                  {
                    "criteria": "<=800",
                    "targetValue": 800,
                    "violated": false
                  }
                ]
              }
            ],
            "result": "pass",
            "score": 100,
            "sloFileContent": "LS0tCnNwZWNfdmVyc2lvbjogIjEuMCIKY29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciCiAgY29tcGFyZV93aXRoOiAic2luZ2xlX3Jlc3VsdCIKICBpbmNsdWRlX3Jlc3VsdF93aXRoX3Njb3JlOiAicGFzcyIKICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxCmZpbHRlcjoKb2JqZWN0aXZlczoKICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1IgogICAga2V5X3NsaTogZmFsc2UKICAgIHBhc3M6ICAgICAgICAgICAgICMgcGFzcyBpZiAocmVsYXRpdmUgY2hhbmdlIDw9IDEwJSBBTkQgYWJzb2x1dGUgdmFsdWUgaXMgPCA2MDBtcykKICAgICAgLSBjcml0ZXJpYToKICAgICAgICAgIC0gIjw9KzEwJSIgICMgcmVsYXRpdmUgdmFsdWVzIHJlcXVpcmUgYSBwcmVmaXhlZCBzaWduIChwbHVzIG9yIG1pbnVzKQogICAgICAgICAgLSAiPDYwMCIgICAgIyBhYnNvbHV0ZSB2YWx1ZXMgb25seSByZXF1aXJlIGEgbG9naWNhbCBvcGVyYXRvcgogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nCiAgICAgIC0gY3JpdGVyaWE6CiAgICAgICAgICAtICI8PTgwMCIKICAgIHdlaWdodDogMQp0b3RhbF9zY29yZToKICBwYXNzOiAiOTAlIgogIHdhcm5pbmc6ICI3NSUi",
            "timeEnd": "2021-04-08T15:23:59Z",
            "timeStart": "2021-04-08T15:22:00Z"
          },
          "labels": {
            "DtCreds": "dynatrace"
          },
          "project": "sockshop",
          "result": "pass",
          "service": "carts",
          "stage": "staging",
          "status": "succeeded"
        },
        "id": "c38712d2-0e2f-4412-bd39-fed91f333510",
        "source": "lighthouse-service",
        "specversion": "1.0",
        "time": "2021-04-08T15:26:07.223Z",
        "type": "sh.keptn.event.evaluation.finished",
        "shkeptncontext": "777b25e4-4f47-4b43-a4d4-88c87037df4c",
        "triggeredid": "32f1be11-5936-447f-8f62-56429da702bc"
      },
      "heatmapLabel": "2021-04-08 17:26"
    },
    {
      "traces": [],
      "data": {
        "evaluation": {
          "comparedEvents": [
            "c38712d2-0e2f-4412-bd39-fed91f333510"
          ],
          "gitCommit": "",
          "indicatorResults": [
            {
              "displayName": "",
              "keySli": false,
              "passTargets": [
                {
                  "criteria": "<=+10%",
                  "targetValue": 392.5944454106811,
                  "violated": false
                },
                {
                  "criteria": "<600",
                  "targetValue": 600,
                  "violated": false
                }
              ],
              "score": 1,
              "status": "pass",
              "value": {
                "metric": "response_time_p95",
                "success": true,
                "value": 315.8696785366614
              },
              "warningTargets": [
                {
                  "criteria": "<=800",
                  "targetValue": 800,
                  "violated": false
                }
              ]
            }
          ],
          "result": "pass",
          "score": 100,
          "sloFileContent": "LS0tCnNwZWNfdmVyc2lvbjogIjEuMCIKY29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciCiAgY29tcGFyZV93aXRoOiAic2luZ2xlX3Jlc3VsdCIKICBpbmNsdWRlX3Jlc3VsdF93aXRoX3Njb3JlOiAicGFzcyIKICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxCmZpbHRlcjoKb2JqZWN0aXZlczoKICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1IgogICAga2V5X3NsaTogZmFsc2UKICAgIHBhc3M6ICAgICAgICAgICAgICMgcGFzcyBpZiAocmVsYXRpdmUgY2hhbmdlIDw9IDEwJSBBTkQgYWJzb2x1dGUgdmFsdWUgaXMgPCA2MDBtcykKICAgICAgLSBjcml0ZXJpYToKICAgICAgICAgIC0gIjw9KzEwJSIgICMgcmVsYXRpdmUgdmFsdWVzIHJlcXVpcmUgYSBwcmVmaXhlZCBzaWduIChwbHVzIG9yIG1pbnVzKQogICAgICAgICAgLSAiPDYwMCIgICAgIyBhYnNvbHV0ZSB2YWx1ZXMgb25seSByZXF1aXJlIGEgbG9naWNhbCBvcGVyYXRvcgogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nCiAgICAgIC0gY3JpdGVyaWE6CiAgICAgICAgICAtICI8PTgwMCIKICAgIHdlaWdodDogMQp0b3RhbF9zY29yZToKICBwYXNzOiAiOTAlIgogIHdhcm5pbmc6ICI3NSUi",
          "timeEnd": "2021-04-08T15:45:02Z",
          "timeStart": "2021-04-08T15:43:09Z",
          "sloFileContentParsed": "---\nspec_version: \"1.0\"\ncomparison:\n  aggregate_function: \"avg\"\n  compare_with: \"single_result\"\n  include_result_with_score: \"pass\"\n  number_of_comparison_results: 1\nfilter:\nobjectives:\n  - sli: \"response_time_p95\"\n    key_sli: false\n    pass:             # pass if (relative change <= 10% AND absolute value is < 600ms)\n      - criteria:\n          - \"<=+10%\"  # relative values require a prefixed sign (plus or minus)\n          - \"<600\"    # absolute values only require a logical operator\n    warning:          # if the response time is below 800ms, the result should be a warning\n      - criteria:\n          - \"<=800\"\n    weight: 1\ntotal_score:\n  pass: \"90%\"\n  warning: \"75%\"",
          "score_pass": "90",
          "score_warning": "75",
          "compare_with": "single_result\n",
          "include_result_with_score": "pass\n",
          "number_of_comparison_results": 1,
          "number_of_missing_comparison_results": 0
        },
        "labels": {
          "DtCreds": "dynatrace"
        },
        "project": "sockshop",
        "result": "pass",
        "service": "carts",
        "stage": "staging",
        "status": "succeeded"
      },
      "id": "01b1eff1-5bd9-4955-b2ef-30fac990b761",
      "source": "lighthouse-service",
      "specversion": "1.0",
      "time": "2021-04-08T15:47:10.200Z",
      "type": "sh.keptn.event.evaluation.finished",
      "shkeptncontext": "468286b0-9ea8-450e-800d-51897947c668",
      "triggeredid": "7a718389-0c3a-4997-a917-9112df3f8c2a",
      "plainEvent": {
        "data": {
          "evaluation": {
            "comparedEvents": [
              "c38712d2-0e2f-4412-bd39-fed91f333510"
            ],
            "gitCommit": "",
            "indicatorResults": [
              {
                "displayName": "",
                "keySli": false,
                "passTargets": [
                  {
                    "criteria": "<=+10%",
                    "targetValue": 392.5944454106811,
                    "violated": false
                  },
                  {
                    "criteria": "<600",
                    "targetValue": 600,
                    "violated": false
                  }
                ],
                "score": 1,
                "status": "pass",
                "value": {
                  "metric": "response_time_p95",
                  "success": true,
                  "value": 315.8696785366614
                },
                "warningTargets": [
                  {
                    "criteria": "<=800",
                    "targetValue": 800,
                    "violated": false
                  }
                ]
              }
            ],
            "result": "pass",
            "score": 100,
            "sloFileContent": "LS0tCnNwZWNfdmVyc2lvbjogIjEuMCIKY29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciCiAgY29tcGFyZV93aXRoOiAic2luZ2xlX3Jlc3VsdCIKICBpbmNsdWRlX3Jlc3VsdF93aXRoX3Njb3JlOiAicGFzcyIKICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxCmZpbHRlcjoKb2JqZWN0aXZlczoKICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1IgogICAga2V5X3NsaTogZmFsc2UKICAgIHBhc3M6ICAgICAgICAgICAgICMgcGFzcyBpZiAocmVsYXRpdmUgY2hhbmdlIDw9IDEwJSBBTkQgYWJzb2x1dGUgdmFsdWUgaXMgPCA2MDBtcykKICAgICAgLSBjcml0ZXJpYToKICAgICAgICAgIC0gIjw9KzEwJSIgICMgcmVsYXRpdmUgdmFsdWVzIHJlcXVpcmUgYSBwcmVmaXhlZCBzaWduIChwbHVzIG9yIG1pbnVzKQogICAgICAgICAgLSAiPDYwMCIgICAgIyBhYnNvbHV0ZSB2YWx1ZXMgb25seSByZXF1aXJlIGEgbG9naWNhbCBvcGVyYXRvcgogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nCiAgICAgIC0gY3JpdGVyaWE6CiAgICAgICAgICAtICI8PTgwMCIKICAgIHdlaWdodDogMQp0b3RhbF9zY29yZToKICBwYXNzOiAiOTAlIgogIHdhcm5pbmc6ICI3NSUi",
            "timeEnd": "2021-04-08T15:45:02Z",
            "timeStart": "2021-04-08T15:43:09Z"
          },
          "labels": {
            "DtCreds": "dynatrace"
          },
          "project": "sockshop",
          "result": "pass",
          "service": "carts",
          "stage": "staging",
          "status": "succeeded"
        },
        "id": "01b1eff1-5bd9-4955-b2ef-30fac990b761",
        "source": "lighthouse-service",
        "specversion": "1.0",
        "time": "2021-04-08T15:47:10.200Z",
        "type": "sh.keptn.event.evaluation.finished",
        "shkeptncontext": "468286b0-9ea8-450e-800d-51897947c668",
        "triggeredid": "7a718389-0c3a-4997-a917-9112df3f8c2a"
      },
      "heatmapLabel": "2021-04-08 17:47"
    }
  ]
},
  "id": "01b1eff1-5bd9-4955-b2ef-30fac990b761",
  "source": "lighthouse-service",
  "specversion": "1.0",
  "time": "2021-04-08T15:47:10.200Z",
  "type": "sh.keptn.event.evaluation.finished",
  "shkeptncontext": "468286b0-9ea8-450e-800d-51897947c668",
  "triggeredid": "7a718389-0c3a-4997-a917-9112df3f8c2a",
  "plainEvent": {
  "data": {
    "evaluation": {
      "comparedEvents": [
        "c38712d2-0e2f-4412-bd39-fed91f333510"
      ],
        "gitCommit": "",
        "indicatorResults": [
        {
          "displayName": "",
          "keySli": false,
          "passTargets": [
            {
              "criteria": "<=+10%",
              "targetValue": 392.5944454106811,
              "violated": false
            },
            {
              "criteria": "<600",
              "targetValue": 600,
              "violated": false
            }
          ],
          "score": 1,
          "status": "pass",
          "value": {
            "metric": "response_time_p95",
            "success": true,
            "value": 315.8696785366614
          },
          "warningTargets": [
            {
              "criteria": "<=800",
              "targetValue": 800,
              "violated": false
            }
          ]
        }
      ],
        "result": "pass",
        "score": 100,
        "sloFileContent": "LS0tCnNwZWNfdmVyc2lvbjogIjEuMCIKY29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciCiAgY29tcGFyZV93aXRoOiAic2luZ2xlX3Jlc3VsdCIKICBpbmNsdWRlX3Jlc3VsdF93aXRoX3Njb3JlOiAicGFzcyIKICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxCmZpbHRlcjoKb2JqZWN0aXZlczoKICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1IgogICAga2V5X3NsaTogZmFsc2UKICAgIHBhc3M6ICAgICAgICAgICAgICMgcGFzcyBpZiAocmVsYXRpdmUgY2hhbmdlIDw9IDEwJSBBTkQgYWJzb2x1dGUgdmFsdWUgaXMgPCA2MDBtcykKICAgICAgLSBjcml0ZXJpYToKICAgICAgICAgIC0gIjw9KzEwJSIgICMgcmVsYXRpdmUgdmFsdWVzIHJlcXVpcmUgYSBwcmVmaXhlZCBzaWduIChwbHVzIG9yIG1pbnVzKQogICAgICAgICAgLSAiPDYwMCIgICAgIyBhYnNvbHV0ZSB2YWx1ZXMgb25seSByZXF1aXJlIGEgbG9naWNhbCBvcGVyYXRvcgogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nCiAgICAgIC0gY3JpdGVyaWE6CiAgICAgICAgICAtICI8PTgwMCIKICAgIHdlaWdodDogMQp0b3RhbF9zY29yZToKICBwYXNzOiAiOTAlIgogIHdhcm5pbmc6ICI3NSUi",
        "timeEnd": "2021-04-08T15:45:02Z",
        "timeStart": "2021-04-08T15:43:09Z"
    },
    "labels": {
      "DtCreds": "dynatrace"
    },
    "project": "sockshop",
      "result": "pass",
      "service": "carts",
      "stage": "staging",
      "status": "succeeded"
  },
  "id": "01b1eff1-5bd9-4955-b2ef-30fac990b761",
    "source": "lighthouse-service",
    "specversion": "1.0",
    "time": "2021-04-08T15:47:10.200Z",
    "type": "sh.keptn.event.evaluation.finished",
    "shkeptncontext": "468286b0-9ea8-450e-800d-51897947c668",
    "triggeredid": "7a718389-0c3a-4997-a917-9112df3f8c2a"
},
  "finished": true,
  "icon": "traffic-light",
  "label": "evaluation"
};

export const Evaluations: Trace = Trace.fromJSON(evaluationData);
