import {Trace} from "../../_models/trace";
import {traceData} from "./traces-mock";

const evaluationData = [
  {
    "traces": [],
    "data": {
      "evaluation": {
        "comparedEvents": [
          "8ad74029-5ceb-4680-b697-d1701078faff"
        ],
        "gitCommit": "",
        "indicatorResults": [
          {
            "displayName": "This is my very very very very very very very very very very very very long displayName",
            "keySli": false,
            "passTargets": [
              {
                "criteria": "<=+10%",
                "targetValue": 11.663352564204615,
                "violated": true
              },
              {
                "criteria": "<600",
                "targetValue": 600,
                "violated": false
              }
            ],
            "score": 0.5,
            "status": "warning",
            "value": {
              "metric": "response_time_p95",
              "success": true,
              "value": 331.4225933914111
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
        "result": "fail",
        "score": 50,
        "sloFileContent": "LS0tDQpzcGVjX3ZlcnNpb246ICIxLjAiDQpjb21wYXJpc29uOg0KICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciDQogIGNvbXBhcmVfd2l0aDogInNpbmdsZV9yZXN1bHQiDQogIGluY2x1ZGVfcmVzdWx0X3dpdGhfc2NvcmU6ICJwYXNzIg0KICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxDQpmaWx0ZXI6DQpvYmplY3RpdmVzOg0KICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1Ig0KICAgIGRpc3BsYXlOYW1lOiAiVGhpcyBpcyBteSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSBsb25nIGRpc3BsYXlOYW1lIg0KICAgIGtleV9zbGk6IGZhbHNlDQogICAgcGFzczogICAgICAgICAgICAgIyBwYXNzIGlmIChyZWxhdGl2ZSBjaGFuZ2UgPD0gMTAlIEFORCBhYnNvbHV0ZSB2YWx1ZSBpcyA8IDYwMG1zKQ0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykNCiAgICAgICAgICAtICI8NjAwIiAgICAjIGFic29sdXRlIHZhbHVlcyBvbmx5IHJlcXVpcmUgYSBsb2dpY2FsIG9wZXJhdG9yDQogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nDQogICAgICAtIGNyaXRlcmlhOg0KICAgICAgICAgIC0gIjw9ODAwIg0KICAgIHdlaWdodDogMQ0KdG90YWxfc2NvcmU6DQogIHBhc3M6ICI5MCUiDQogIHdhcm5pbmc6ICI3NSUi",
        "timeEnd": "2021-03-16T12:46:58Z",
        "timeStart": "2021-03-16T12:44:54Z",
        "sloFileContentParsed": "---\r\nspec_version: \"1.0\"\r\ncomparison:\r\n  aggregate_function: \"avg\"\r\n  compare_with: \"single_result\"\r\n  include_result_with_score: \"pass\"\r\n  number_of_comparison_results: 1\r\nfilter:\r\nobjectives:\r\n  - sli: \"response_time_p95\"\r\n    displayName: \"This is my very very very very very very very very very very very very long displayName\"\r\n    key_sli: false\r\n    pass:             # pass if (relative change <= 10% AND absolute value is < 600ms)\r\n      - criteria:\r\n          - \"<=+10%\"  # relative values require a prefixed sign (plus or minus)\r\n          - \"<600\"    # absolute values only require a logical operator\r\n    warning:          # if the response time is below 800ms, the result should be a warning\r\n      - criteria:\r\n          - \"<=800\"\r\n    weight: 1\r\ntotal_score:\r\n  pass: \"90%\"\r\n  warning: \"75%\"",
        "score_pass": "90",
        "score_warning": "75",
        "compare_with": "single_result\r\n",
        "include_result_with_score": "pass\r\n",
        "number_of_comparison_results": 1,
        "number_of_missing_comparison_results": 0
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
    "id": "843a4a6b-4dec-4328-8579-6af3d1a6c7cb",
    "source": "lighthouse-service",
    "specversion": "1.0",
    "time": "2021-03-16T12:48:05.732Z",
    "type": "sh.keptn.event.evaluation.finished",
    "shkeptncontext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
    "triggeredid": "7c105021-3a50-47c7-aaa9-2e6286b17d89",
    "plainEvent": {
      "data": {
        "evaluation": {
          "comparedEvents": [
            "8ad74029-5ceb-4680-b697-d1701078faff"
          ],
          "gitCommit": "",
          "indicatorResults": [
            {
              "displayName": "This is my very very very very very very very very very very very very long displayName",
              "keySli": false,
              "passTargets": [
                {
                  "criteria": "<=+10%",
                  "targetValue": 11.663352564204615,
                  "violated": true
                },
                {
                  "criteria": "<600",
                  "targetValue": 600,
                  "violated": false
                }
              ],
              "score": 0.5,
              "status": "warning",
              "value": {
                "metric": "response_time_p95",
                "success": true,
                "value": 331.4225933914111
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
          "result": "fail",
          "score": 50,
          "sloFileContent": "LS0tDQpzcGVjX3ZlcnNpb246ICIxLjAiDQpjb21wYXJpc29uOg0KICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciDQogIGNvbXBhcmVfd2l0aDogInNpbmdsZV9yZXN1bHQiDQogIGluY2x1ZGVfcmVzdWx0X3dpdGhfc2NvcmU6ICJwYXNzIg0KICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxDQpmaWx0ZXI6DQpvYmplY3RpdmVzOg0KICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1Ig0KICAgIGRpc3BsYXlOYW1lOiAiVGhpcyBpcyBteSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSBsb25nIGRpc3BsYXlOYW1lIg0KICAgIGtleV9zbGk6IGZhbHNlDQogICAgcGFzczogICAgICAgICAgICAgIyBwYXNzIGlmIChyZWxhdGl2ZSBjaGFuZ2UgPD0gMTAlIEFORCBhYnNvbHV0ZSB2YWx1ZSBpcyA8IDYwMG1zKQ0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykNCiAgICAgICAgICAtICI8NjAwIiAgICAjIGFic29sdXRlIHZhbHVlcyBvbmx5IHJlcXVpcmUgYSBsb2dpY2FsIG9wZXJhdG9yDQogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nDQogICAgICAtIGNyaXRlcmlhOg0KICAgICAgICAgIC0gIjw9ODAwIg0KICAgIHdlaWdodDogMQ0KdG90YWxfc2NvcmU6DQogIHBhc3M6ICI5MCUiDQogIHdhcm5pbmc6ICI3NSUi",
          "timeEnd": "2021-03-16T12:46:58Z",
          "timeStart": "2021-03-16T12:44:54Z"
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
      "id": "843a4a6b-4dec-4328-8579-6af3d1a6c7cb",
      "source": "lighthouse-service",
      "specversion": "1.0",
      "time": "2021-03-16T12:48:05.732Z",
      "type": "sh.keptn.event.evaluation.finished",
      "shkeptncontext": "6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432",
      "triggeredid": "7c105021-3a50-47c7-aaa9-2e6286b17d89"
    },
    "heatmapLabel": "2021-03-16 13:48"
  },
  {
    "traces": [],
    "data": {
      "evaluation": {
        "comparedEvents": [
          "8ad74029-5ceb-4680-b697-d1701078faff"
        ],
        "gitCommit": "",
        "indicatorResults": [
          {
            "keySli": false,
            "passTargets": [
              {
                "criteria": "<=+10%",
                "targetValue": 11.663352564204615,
                "violated": true
              },
              {
                "criteria": "<600",
                "targetValue": 600,
                "violated": false
              }
            ],
            "score": 0.5,
            "status": "warning",
            "value": {
              "metric": "response_time_p95",
              "success": true,
              "value": 274.8900806697535
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
        "result": "fail",
        "score": 50,
        "sloFileContent": "LS0tDQpzcGVjX3ZlcnNpb246ICIxLjAiDQpjb21wYXJpc29uOg0KICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciDQogIGNvbXBhcmVfd2l0aDogInNpbmdsZV9yZXN1bHQiDQogIGluY2x1ZGVfcmVzdWx0X3dpdGhfc2NvcmU6ICJwYXNzIg0KICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxDQpmaWx0ZXI6DQpvYmplY3RpdmVzOg0KICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1Ig0KICAgIGRpc3BsYXlOYW1lOiAiVGhpcyBpcyBteSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSBsb25nIGRpc3BsYXlOYW1lIg0KICAgIGtleV9zbGk6IGZhbHNlDQogICAgcGFzczogICAgICAgICAgICAgIyBwYXNzIGlmIChyZWxhdGl2ZSBjaGFuZ2UgPD0gMTAlIEFORCBhYnNvbHV0ZSB2YWx1ZSBpcyA8IDYwMG1zKQ0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykNCiAgICAgICAgICAtICI8NjAwIiAgICAjIGFic29sdXRlIHZhbHVlcyBvbmx5IHJlcXVpcmUgYSBsb2dpY2FsIG9wZXJhdG9yDQogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nDQogICAgICAtIGNyaXRlcmlhOg0KICAgICAgICAgIC0gIjw9ODAwIg0KICAgIHdlaWdodDogMQ0KdG90YWxfc2NvcmU6DQogIHBhc3M6ICI5MCUiDQogIHdhcm5pbmc6ICI3NSUi",
        "timeEnd": "2021-03-12T11:54:58Z",
        "timeStart": "2021-03-12T11:53:12Z"
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
    "id": "b28d36f6-b632-4170-b494-39152537d0ef",
    "source": "lighthouse-service",
    "specversion": "1.0",
    "time": "2021-03-12T11:57:05.133Z",
    "type": "sh.keptn.event.evaluation.finished",
    "shkeptncontext": "5f6c85f9-dd2b-4eb6-bd96-55ea15a09a3a",
    "triggeredid": "97e49b9f-8161-4c79-b68a-3894744cfe91",
    "plainEvent": {
      "data": {
        "evaluation": {
          "comparedEvents": [
            "8ad74029-5ceb-4680-b697-d1701078faff"
          ],
          "gitCommit": "",
          "indicatorResults": [
            {
              "keySli": false,
              "passTargets": [
                {
                  "criteria": "<=+10%",
                  "targetValue": 11.663352564204615,
                  "violated": true
                },
                {
                  "criteria": "<600",
                  "targetValue": 600,
                  "violated": false
                }
              ],
              "score": 0.5,
              "status": "warning",
              "value": {
                "metric": "response_time_p95",
                "success": true,
                "value": 274.8900806697535
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
          "result": "fail",
          "score": 50,
          "sloFileContent": "LS0tDQpzcGVjX3ZlcnNpb246ICIxLjAiDQpjb21wYXJpc29uOg0KICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciDQogIGNvbXBhcmVfd2l0aDogInNpbmdsZV9yZXN1bHQiDQogIGluY2x1ZGVfcmVzdWx0X3dpdGhfc2NvcmU6ICJwYXNzIg0KICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxDQpmaWx0ZXI6DQpvYmplY3RpdmVzOg0KICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1Ig0KICAgIGRpc3BsYXlOYW1lOiAiVGhpcyBpcyBteSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSBsb25nIGRpc3BsYXlOYW1lIg0KICAgIGtleV9zbGk6IGZhbHNlDQogICAgcGFzczogICAgICAgICAgICAgIyBwYXNzIGlmIChyZWxhdGl2ZSBjaGFuZ2UgPD0gMTAlIEFORCBhYnNvbHV0ZSB2YWx1ZSBpcyA8IDYwMG1zKQ0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykNCiAgICAgICAgICAtICI8NjAwIiAgICAjIGFic29sdXRlIHZhbHVlcyBvbmx5IHJlcXVpcmUgYSBsb2dpY2FsIG9wZXJhdG9yDQogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nDQogICAgICAtIGNyaXRlcmlhOg0KICAgICAgICAgIC0gIjw9ODAwIg0KICAgIHdlaWdodDogMQ0KdG90YWxfc2NvcmU6DQogIHBhc3M6ICI5MCUiDQogIHdhcm5pbmc6ICI3NSUi",
          "timeEnd": "2021-03-12T11:54:58Z",
          "timeStart": "2021-03-12T11:53:12Z"
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
      "id": "b28d36f6-b632-4170-b494-39152537d0ef",
      "source": "lighthouse-service",
      "specversion": "1.0",
      "time": "2021-03-12T11:57:05.133Z",
      "type": "sh.keptn.event.evaluation.finished",
      "shkeptncontext": "5f6c85f9-dd2b-4eb6-bd96-55ea15a09a3a",
      "triggeredid": "97e49b9f-8161-4c79-b68a-3894744cfe91"
    },
    "heatmapLabel": "2021-03-12 12:57"
  },
  {
    "traces": [],
    "data": {
      "evaluation": {
        "comparedEvents": [
          "8ad74029-5ceb-4680-b697-d1701078faff"
        ],
        "gitCommit": "",
        "indicatorResults": [
          {
            "keySli": false,
            "passTargets": [
              {
                "criteria": "<=+10%",
                "targetValue": 0,
                "violated": true
              },
              {
                "criteria": "<600",
                "targetValue": 0,
                "violated": true
              }
            ],
            "score": 0,
            "status": "fail",
            "value": {
              "message": "error validating time range: Supplied end-time 2021-03-12 11:09:00 +0000 UTC is too far (>120seconds) in the future (now: 2021-03-12 10:10:15.267498606 +0000 UTC m=+606470.155451374 - diff in sec: -3524.732501394)\n",
              "metric": "response_time_p95",
              "success": false,
              "value": 0
            },
            "warningTargets": [
              {
                "criteria": "<=800",
                "targetValue": 0,
                "violated": true
              }
            ]
          }
        ],
        "result": "fail",
        "score": 0,
        "sloFileContent": "LS0tDQpzcGVjX3ZlcnNpb246ICIxLjAiDQpjb21wYXJpc29uOg0KICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciDQogIGNvbXBhcmVfd2l0aDogInNpbmdsZV9yZXN1bHQiDQogIGluY2x1ZGVfcmVzdWx0X3dpdGhfc2NvcmU6ICJwYXNzIg0KICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxDQpmaWx0ZXI6DQpvYmplY3RpdmVzOg0KICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1Ig0KICAgIGRpc3BsYXlOYW1lOiAiVGhpcyBpcyBteSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSBsb25nIGRpc3BsYXlOYW1lIg0KICAgIGtleV9zbGk6IGZhbHNlDQogICAgcGFzczogICAgICAgICAgICAgIyBwYXNzIGlmIChyZWxhdGl2ZSBjaGFuZ2UgPD0gMTAlIEFORCBhYnNvbHV0ZSB2YWx1ZSBpcyA8IDYwMG1zKQ0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykNCiAgICAgICAgICAtICI8NjAwIiAgICAjIGFic29sdXRlIHZhbHVlcyBvbmx5IHJlcXVpcmUgYSBsb2dpY2FsIG9wZXJhdG9yDQogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nDQogICAgICAtIGNyaXRlcmlhOg0KICAgICAgICAgIC0gIjw9ODAwIg0KICAgIHdlaWdodDogMQ0KdG90YWxfc2NvcmU6DQogIHBhc3M6ICI5MCUiDQogIHdhcm5pbmc6ICI3NSUi",
        "timeEnd": "2021-03-12T11:09:00.000Z",
        "timeStart": "2021-03-12T11:04:00.000Z"
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
    "id": "9f797791-50eb-445d-a00a-ae0002ab98f8",
    "source": "lighthouse-service",
    "specversion": "1.0",
    "time": "2021-03-12T10:10:16.035Z",
    "type": "sh.keptn.event.evaluation.finished",
    "shkeptncontext": "8e452c9b-b8ac-4310-9f8e-bb62acfe4ca4",
    "triggeredid": "4a26765d-53c6-4259-9331-6591e973a3af",
    "plainEvent": {
      "data": {
        "evaluation": {
          "comparedEvents": [
            "8ad74029-5ceb-4680-b697-d1701078faff"
          ],
          "gitCommit": "",
          "indicatorResults": [
            {
              "keySli": false,
              "passTargets": [
                {
                  "criteria": "<=+10%",
                  "targetValue": 0,
                  "violated": true
                },
                {
                  "criteria": "<600",
                  "targetValue": 0,
                  "violated": true
                }
              ],
              "score": 0,
              "status": "fail",
              "value": {
                "message": "error validating time range: Supplied end-time 2021-03-12 11:09:00 +0000 UTC is too far (>120seconds) in the future (now: 2021-03-12 10:10:15.267498606 +0000 UTC m=+606470.155451374 - diff in sec: -3524.732501394)\n",
                "metric": "response_time_p95",
                "success": false,
                "value": 0
              },
              "warningTargets": [
                {
                  "criteria": "<=800",
                  "targetValue": 0,
                  "violated": true
                }
              ]
            }
          ],
          "result": "fail",
          "score": 0,
          "sloFileContent": "LS0tDQpzcGVjX3ZlcnNpb246ICIxLjAiDQpjb21wYXJpc29uOg0KICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciDQogIGNvbXBhcmVfd2l0aDogInNpbmdsZV9yZXN1bHQiDQogIGluY2x1ZGVfcmVzdWx0X3dpdGhfc2NvcmU6ICJwYXNzIg0KICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxDQpmaWx0ZXI6DQpvYmplY3RpdmVzOg0KICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1Ig0KICAgIGRpc3BsYXlOYW1lOiAiVGhpcyBpcyBteSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSBsb25nIGRpc3BsYXlOYW1lIg0KICAgIGtleV9zbGk6IGZhbHNlDQogICAgcGFzczogICAgICAgICAgICAgIyBwYXNzIGlmIChyZWxhdGl2ZSBjaGFuZ2UgPD0gMTAlIEFORCBhYnNvbHV0ZSB2YWx1ZSBpcyA8IDYwMG1zKQ0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykNCiAgICAgICAgICAtICI8NjAwIiAgICAjIGFic29sdXRlIHZhbHVlcyBvbmx5IHJlcXVpcmUgYSBsb2dpY2FsIG9wZXJhdG9yDQogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nDQogICAgICAtIGNyaXRlcmlhOg0KICAgICAgICAgIC0gIjw9ODAwIg0KICAgIHdlaWdodDogMQ0KdG90YWxfc2NvcmU6DQogIHBhc3M6ICI5MCUiDQogIHdhcm5pbmc6ICI3NSUi",
          "timeEnd": "2021-03-12T11:09:00.000Z",
          "timeStart": "2021-03-12T11:04:00.000Z"
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
      "id": "9f797791-50eb-445d-a00a-ae0002ab98f8",
      "source": "lighthouse-service",
      "specversion": "1.0",
      "time": "2021-03-12T10:10:16.035Z",
      "type": "sh.keptn.event.evaluation.finished",
      "shkeptncontext": "8e452c9b-b8ac-4310-9f8e-bb62acfe4ca4",
      "triggeredid": "4a26765d-53c6-4259-9331-6591e973a3af"
    },
    "heatmapLabel": "2021-03-12 11:10"
  },
  {
    "traces": [],
    "data": {
      "evaluation": {
        "comparedEvents": [
          "d9e78aee-f5e8-4d02-a802-2d268f262860"
        ],
        "gitCommit": "",
        "indicatorResults": [
          {
            "keySli": false,
            "passTargets": [
              {
                "criteria": "<=+10%",
                "targetValue": 11.461298340581093,
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
              "value": 10.603047785640559
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
        "sloFileContent": "LS0tDQpzcGVjX3ZlcnNpb246ICIxLjAiDQpjb21wYXJpc29uOg0KICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciDQogIGNvbXBhcmVfd2l0aDogInNpbmdsZV9yZXN1bHQiDQogIGluY2x1ZGVfcmVzdWx0X3dpdGhfc2NvcmU6ICJwYXNzIg0KICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxDQpmaWx0ZXI6DQpvYmplY3RpdmVzOg0KICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1Ig0KICAgIGRpc3BsYXlOYW1lOiAiVGhpcyBpcyBteSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSBsb25nIGRpc3BsYXlOYW1lIg0KICAgIGtleV9zbGk6IGZhbHNlDQogICAgcGFzczogICAgICAgICAgICAgIyBwYXNzIGlmIChyZWxhdGl2ZSBjaGFuZ2UgPD0gMTAlIEFORCBhYnNvbHV0ZSB2YWx1ZSBpcyA8IDYwMG1zKQ0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykNCiAgICAgICAgICAtICI8NjAwIiAgICAjIGFic29sdXRlIHZhbHVlcyBvbmx5IHJlcXVpcmUgYSBsb2dpY2FsIG9wZXJhdG9yDQogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nDQogICAgICAtIGNyaXRlcmlhOg0KICAgICAgICAgIC0gIjw9ODAwIg0KICAgIHdlaWdodDogMQ0KdG90YWxfc2NvcmU6DQogIHBhc3M6ICI5MCUiDQogIHdhcm5pbmc6ICI3NSUi",
        "timeEnd": "2021-03-12T10:09:00.000Z",
        "timeStart": "2021-03-12T10:04:00.000Z"
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
    "id": "8ad74029-5ceb-4680-b697-d1701078faff",
    "source": "lighthouse-service",
    "specversion": "1.0",
    "time": "2021-03-12T10:09:11.430Z",
    "type": "sh.keptn.event.evaluation.finished",
    "shkeptncontext": "a1aa2bce-fd76-4a98-8242-7593e9e931fb",
    "triggeredid": "53419ada-f9f2-4b93-a9e8-396af7638c46",
    "plainEvent": {
      "data": {
        "evaluation": {
          "comparedEvents": [
            "d9e78aee-f5e8-4d02-a802-2d268f262860"
          ],
          "gitCommit": "",
          "indicatorResults": [
            {
              "keySli": false,
              "passTargets": [
                {
                  "criteria": "<=+10%",
                  "targetValue": 11.461298340581093,
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
                "value": 10.603047785640559
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
          "sloFileContent": "LS0tDQpzcGVjX3ZlcnNpb246ICIxLjAiDQpjb21wYXJpc29uOg0KICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciDQogIGNvbXBhcmVfd2l0aDogInNpbmdsZV9yZXN1bHQiDQogIGluY2x1ZGVfcmVzdWx0X3dpdGhfc2NvcmU6ICJwYXNzIg0KICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxDQpmaWx0ZXI6DQpvYmplY3RpdmVzOg0KICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1Ig0KICAgIGRpc3BsYXlOYW1lOiAiVGhpcyBpcyBteSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSBsb25nIGRpc3BsYXlOYW1lIg0KICAgIGtleV9zbGk6IGZhbHNlDQogICAgcGFzczogICAgICAgICAgICAgIyBwYXNzIGlmIChyZWxhdGl2ZSBjaGFuZ2UgPD0gMTAlIEFORCBhYnNvbHV0ZSB2YWx1ZSBpcyA8IDYwMG1zKQ0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykNCiAgICAgICAgICAtICI8NjAwIiAgICAjIGFic29sdXRlIHZhbHVlcyBvbmx5IHJlcXVpcmUgYSBsb2dpY2FsIG9wZXJhdG9yDQogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nDQogICAgICAtIGNyaXRlcmlhOg0KICAgICAgICAgIC0gIjw9ODAwIg0KICAgIHdlaWdodDogMQ0KdG90YWxfc2NvcmU6DQogIHBhc3M6ICI5MCUiDQogIHdhcm5pbmc6ICI3NSUi",
          "timeEnd": "2021-03-12T10:09:00.000Z",
          "timeStart": "2021-03-12T10:04:00.000Z"
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
      "id": "8ad74029-5ceb-4680-b697-d1701078faff",
      "source": "lighthouse-service",
      "specversion": "1.0",
      "time": "2021-03-12T10:09:11.430Z",
      "type": "sh.keptn.event.evaluation.finished",
      "shkeptncontext": "a1aa2bce-fd76-4a98-8242-7593e9e931fb",
      "triggeredid": "53419ada-f9f2-4b93-a9e8-396af7638c46"
    },
    "heatmapLabel": "2021-03-12 11:09"
  },
  {
    "traces": [],
    "data": {
      "evaluation": {
        "comparedEvents": [
          "e7ca246d-b440-4d10-9005-04361caa5cda"
        ],
        "gitCommit": "",
        "indicatorResults": [
          {
            "keySli": false,
            "passTargets": [
              {
                "criteria": "<=+10%",
                "targetValue": 433.8437787504008,
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
              "value": 10.419362127800994
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
        "sloFileContent": "LS0tDQpzcGVjX3ZlcnNpb246ICIxLjAiDQpjb21wYXJpc29uOg0KICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciDQogIGNvbXBhcmVfd2l0aDogInNpbmdsZV9yZXN1bHQiDQogIGluY2x1ZGVfcmVzdWx0X3dpdGhfc2NvcmU6ICJwYXNzIg0KICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxDQpmaWx0ZXI6DQpvYmplY3RpdmVzOg0KICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1Ig0KICAgIGRpc3BsYXlOYW1lOiAiVGhpcyBpcyBteSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSBsb25nIGRpc3BsYXlOYW1lIg0KICAgIGtleV9zbGk6IGZhbHNlDQogICAgcGFzczogICAgICAgICAgICAgIyBwYXNzIGlmIChyZWxhdGl2ZSBjaGFuZ2UgPD0gMTAlIEFORCBhYnNvbHV0ZSB2YWx1ZSBpcyA8IDYwMG1zKQ0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykNCiAgICAgICAgICAtICI8NjAwIiAgICAjIGFic29sdXRlIHZhbHVlcyBvbmx5IHJlcXVpcmUgYSBsb2dpY2FsIG9wZXJhdG9yDQogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nDQogICAgICAtIGNyaXRlcmlhOg0KICAgICAgICAgIC0gIjw9ODAwIg0KICAgIHdlaWdodDogMQ0KdG90YWxfc2NvcmU6DQogIHBhc3M6ICI5MCUiDQogIHdhcm5pbmc6ICI3NSUi",
        "timeEnd": "2021-03-12T10:08:00.000Z",
        "timeStart": "2021-03-12T10:03:00.000Z"
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
    "id": "d9e78aee-f5e8-4d02-a802-2d268f262860",
    "source": "lighthouse-service",
    "specversion": "1.0",
    "time": "2021-03-12T10:08:05.786Z",
    "type": "sh.keptn.event.evaluation.finished",
    "shkeptncontext": "d316195e-9c1f-48e1-95c8-810969e4df0e",
    "triggeredid": "93a75e28-c1ac-4d0c-80de-582135e933e4",
    "plainEvent": {
      "data": {
        "evaluation": {
          "comparedEvents": [
            "e7ca246d-b440-4d10-9005-04361caa5cda"
          ],
          "gitCommit": "",
          "indicatorResults": [
            {
              "keySli": false,
              "passTargets": [
                {
                  "criteria": "<=+10%",
                  "targetValue": 433.8437787504008,
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
                "value": 10.419362127800994
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
          "sloFileContent": "LS0tDQpzcGVjX3ZlcnNpb246ICIxLjAiDQpjb21wYXJpc29uOg0KICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciDQogIGNvbXBhcmVfd2l0aDogInNpbmdsZV9yZXN1bHQiDQogIGluY2x1ZGVfcmVzdWx0X3dpdGhfc2NvcmU6ICJwYXNzIg0KICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxDQpmaWx0ZXI6DQpvYmplY3RpdmVzOg0KICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1Ig0KICAgIGRpc3BsYXlOYW1lOiAiVGhpcyBpcyBteSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSBsb25nIGRpc3BsYXlOYW1lIg0KICAgIGtleV9zbGk6IGZhbHNlDQogICAgcGFzczogICAgICAgICAgICAgIyBwYXNzIGlmIChyZWxhdGl2ZSBjaGFuZ2UgPD0gMTAlIEFORCBhYnNvbHV0ZSB2YWx1ZSBpcyA8IDYwMG1zKQ0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykNCiAgICAgICAgICAtICI8NjAwIiAgICAjIGFic29sdXRlIHZhbHVlcyBvbmx5IHJlcXVpcmUgYSBsb2dpY2FsIG9wZXJhdG9yDQogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nDQogICAgICAtIGNyaXRlcmlhOg0KICAgICAgICAgIC0gIjw9ODAwIg0KICAgIHdlaWdodDogMQ0KdG90YWxfc2NvcmU6DQogIHBhc3M6ICI5MCUiDQogIHdhcm5pbmc6ICI3NSUi",
          "timeEnd": "2021-03-12T10:08:00.000Z",
          "timeStart": "2021-03-12T10:03:00.000Z"
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
      "id": "d9e78aee-f5e8-4d02-a802-2d268f262860",
      "source": "lighthouse-service",
      "specversion": "1.0",
      "time": "2021-03-12T10:08:05.786Z",
      "type": "sh.keptn.event.evaluation.finished",
      "shkeptncontext": "d316195e-9c1f-48e1-95c8-810969e4df0e",
      "triggeredid": "93a75e28-c1ac-4d0c-80de-582135e933e4"
    },
    "heatmapLabel": "2021-03-12 11:08"
  },
  {
    "traces": [],
    "data": {
      "evaluation": {
        "comparedEvents": [
          "e7ca246d-b440-4d10-9005-04361caa5cda"
        ],
        "gitCommit": "",
        "indicatorResults": [
          {
            "keySli": false,
            "passTargets": [
              {
                "criteria": "<=+10%",
                "targetValue": 0,
                "violated": true
              },
              {
                "criteria": "<600",
                "targetValue": 0,
                "violated": true
              }
            ],
            "score": 0,
            "status": "fail",
            "value": {
              "message": "error validating time range: Supplied end-time 2021-03-12 10:11:00 +0000 UTC is too far (>120seconds) in the future (now: 2021-03-12 10:06:28.499073332 +0000 UTC m=+606243.387026108 - diff in sec: -271.500926668)\n",
              "metric": "response_time_p95",
              "success": false,
              "value": 0
            },
            "warningTargets": [
              {
                "criteria": "<=800",
                "targetValue": 0,
                "violated": true
              }
            ]
          }
        ],
        "result": "fail",
        "score": 0,
        "sloFileContent": "LS0tDQpzcGVjX3ZlcnNpb246ICIxLjAiDQpjb21wYXJpc29uOg0KICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciDQogIGNvbXBhcmVfd2l0aDogInNpbmdsZV9yZXN1bHQiDQogIGluY2x1ZGVfcmVzdWx0X3dpdGhfc2NvcmU6ICJwYXNzIg0KICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxDQpmaWx0ZXI6DQpvYmplY3RpdmVzOg0KICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1Ig0KICAgIGRpc3BsYXlOYW1lOiAiVGhpcyBpcyBteSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSBsb25nIGRpc3BsYXlOYW1lIg0KICAgIGtleV9zbGk6IGZhbHNlDQogICAgcGFzczogICAgICAgICAgICAgIyBwYXNzIGlmIChyZWxhdGl2ZSBjaGFuZ2UgPD0gMTAlIEFORCBhYnNvbHV0ZSB2YWx1ZSBpcyA8IDYwMG1zKQ0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykNCiAgICAgICAgICAtICI8NjAwIiAgICAjIGFic29sdXRlIHZhbHVlcyBvbmx5IHJlcXVpcmUgYSBsb2dpY2FsIG9wZXJhdG9yDQogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nDQogICAgICAtIGNyaXRlcmlhOg0KICAgICAgICAgIC0gIjw9ODAwIg0KICAgIHdlaWdodDogMQ0KdG90YWxfc2NvcmU6DQogIHBhc3M6ICI5MCUiDQogIHdhcm5pbmc6ICI3NSUi",
        "timeEnd": "2021-03-12T10:11:00.000Z",
        "timeStart": "2021-03-12T10:06:00.000Z"
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
    "id": "82d10352-2c02-4e80-9d94-4bdaafa3eb28",
    "source": "lighthouse-service",
    "specversion": "1.0",
    "time": "2021-03-12T10:06:29.300Z",
    "type": "sh.keptn.event.evaluation.finished",
    "shkeptncontext": "3299408b-62d5-4c81-8594-a107b7d059ff",
    "triggeredid": "2b827df9-7ef4-4004-a642-25e9f1efc08c",
    "plainEvent": {
      "data": {
        "evaluation": {
          "comparedEvents": [
            "e7ca246d-b440-4d10-9005-04361caa5cda"
          ],
          "gitCommit": "",
          "indicatorResults": [
            {
              "keySli": false,
              "passTargets": [
                {
                  "criteria": "<=+10%",
                  "targetValue": 0,
                  "violated": true
                },
                {
                  "criteria": "<600",
                  "targetValue": 0,
                  "violated": true
                }
              ],
              "score": 0,
              "status": "fail",
              "value": {
                "message": "error validating time range: Supplied end-time 2021-03-12 10:11:00 +0000 UTC is too far (>120seconds) in the future (now: 2021-03-12 10:06:28.499073332 +0000 UTC m=+606243.387026108 - diff in sec: -271.500926668)\n",
                "metric": "response_time_p95",
                "success": false,
                "value": 0
              },
              "warningTargets": [
                {
                  "criteria": "<=800",
                  "targetValue": 0,
                  "violated": true
                }
              ]
            }
          ],
          "result": "fail",
          "score": 0,
          "sloFileContent": "LS0tDQpzcGVjX3ZlcnNpb246ICIxLjAiDQpjb21wYXJpc29uOg0KICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciDQogIGNvbXBhcmVfd2l0aDogInNpbmdsZV9yZXN1bHQiDQogIGluY2x1ZGVfcmVzdWx0X3dpdGhfc2NvcmU6ICJwYXNzIg0KICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxDQpmaWx0ZXI6DQpvYmplY3RpdmVzOg0KICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1Ig0KICAgIGRpc3BsYXlOYW1lOiAiVGhpcyBpcyBteSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSBsb25nIGRpc3BsYXlOYW1lIg0KICAgIGtleV9zbGk6IGZhbHNlDQogICAgcGFzczogICAgICAgICAgICAgIyBwYXNzIGlmIChyZWxhdGl2ZSBjaGFuZ2UgPD0gMTAlIEFORCBhYnNvbHV0ZSB2YWx1ZSBpcyA8IDYwMG1zKQ0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykNCiAgICAgICAgICAtICI8NjAwIiAgICAjIGFic29sdXRlIHZhbHVlcyBvbmx5IHJlcXVpcmUgYSBsb2dpY2FsIG9wZXJhdG9yDQogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nDQogICAgICAtIGNyaXRlcmlhOg0KICAgICAgICAgIC0gIjw9ODAwIg0KICAgIHdlaWdodDogMQ0KdG90YWxfc2NvcmU6DQogIHBhc3M6ICI5MCUiDQogIHdhcm5pbmc6ICI3NSUi",
          "timeEnd": "2021-03-12T10:11:00.000Z",
          "timeStart": "2021-03-12T10:06:00.000Z"
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
      "id": "82d10352-2c02-4e80-9d94-4bdaafa3eb28",
      "source": "lighthouse-service",
      "specversion": "1.0",
      "time": "2021-03-12T10:06:29.300Z",
      "type": "sh.keptn.event.evaluation.finished",
      "shkeptncontext": "3299408b-62d5-4c81-8594-a107b7d059ff",
      "triggeredid": "2b827df9-7ef4-4004-a642-25e9f1efc08c"
    },
    "heatmapLabel": "2021-03-12 11:06"
  },
  {
    "traces": [],
    "data": {
      "evaluation": {
        "comparedEvents": [
          "e7ca246d-b440-4d10-9005-04361caa5cda"
        ],
        "gitCommit": "",
        "indicatorResults": [
          {
            "keySli": false,
            "passTargets": [
              {
                "criteria": "<=+10%",
                "targetValue": 0,
                "violated": true
              },
              {
                "criteria": "<600",
                "targetValue": 0,
                "violated": true
              }
            ],
            "score": 0,
            "status": "fail",
            "value": {
              "message": "error validating time range: Supplied end-time 2021-03-12 11:03:00 +0000 UTC is too far (>120seconds) in the future (now: 2021-03-12 10:05:15.771595629 +0000 UTC m=+606170.659548404 - diff in sec: -3464.228404371)\n",
              "metric": "response_time_p95",
              "success": false,
              "value": 0
            },
            "warningTargets": [
              {
                "criteria": "<=800",
                "targetValue": 0,
                "violated": true
              }
            ]
          }
        ],
        "result": "fail",
        "score": 0,
        "sloFileContent": "LS0tDQpzcGVjX3ZlcnNpb246ICIxLjAiDQpjb21wYXJpc29uOg0KICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciDQogIGNvbXBhcmVfd2l0aDogInNpbmdsZV9yZXN1bHQiDQogIGluY2x1ZGVfcmVzdWx0X3dpdGhfc2NvcmU6ICJwYXNzIg0KICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxDQpmaWx0ZXI6DQpvYmplY3RpdmVzOg0KICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1Ig0KICAgIGRpc3BsYXlOYW1lOiAiVGhpcyBpcyBteSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSBsb25nIGRpc3BsYXlOYW1lIg0KICAgIGtleV9zbGk6IGZhbHNlDQogICAgcGFzczogICAgICAgICAgICAgIyBwYXNzIGlmIChyZWxhdGl2ZSBjaGFuZ2UgPD0gMTAlIEFORCBhYnNvbHV0ZSB2YWx1ZSBpcyA8IDYwMG1zKQ0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykNCiAgICAgICAgICAtICI8NjAwIiAgICAjIGFic29sdXRlIHZhbHVlcyBvbmx5IHJlcXVpcmUgYSBsb2dpY2FsIG9wZXJhdG9yDQogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nDQogICAgICAtIGNyaXRlcmlhOg0KICAgICAgICAgIC0gIjw9ODAwIg0KICAgIHdlaWdodDogMQ0KdG90YWxfc2NvcmU6DQogIHBhc3M6ICI5MCUiDQogIHdhcm5pbmc6ICI3NSUi",
        "timeEnd": "2021-03-12T11:03:00.000Z",
        "timeStart": "2021-03-12T10:58:00.000Z"
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
    "id": "af121976-d640-4a41-9fdd-6a2887b7b796",
    "source": "lighthouse-service",
    "specversion": "1.0",
    "time": "2021-03-12T10:05:16.544Z",
    "type": "sh.keptn.event.evaluation.finished",
    "shkeptncontext": "76207fe5-5b72-40db-a74b-e220ae0892c2",
    "triggeredid": "6ada8320-8855-4e12-95c1-0b7fddc91e8d",
    "plainEvent": {
      "data": {
        "evaluation": {
          "comparedEvents": [
            "e7ca246d-b440-4d10-9005-04361caa5cda"
          ],
          "gitCommit": "",
          "indicatorResults": [
            {
              "keySli": false,
              "passTargets": [
                {
                  "criteria": "<=+10%",
                  "targetValue": 0,
                  "violated": true
                },
                {
                  "criteria": "<600",
                  "targetValue": 0,
                  "violated": true
                }
              ],
              "score": 0,
              "status": "fail",
              "value": {
                "message": "error validating time range: Supplied end-time 2021-03-12 11:03:00 +0000 UTC is too far (>120seconds) in the future (now: 2021-03-12 10:05:15.771595629 +0000 UTC m=+606170.659548404 - diff in sec: -3464.228404371)\n",
                "metric": "response_time_p95",
                "success": false,
                "value": 0
              },
              "warningTargets": [
                {
                  "criteria": "<=800",
                  "targetValue": 0,
                  "violated": true
                }
              ]
            }
          ],
          "result": "fail",
          "score": 0,
          "sloFileContent": "LS0tDQpzcGVjX3ZlcnNpb246ICIxLjAiDQpjb21wYXJpc29uOg0KICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciDQogIGNvbXBhcmVfd2l0aDogInNpbmdsZV9yZXN1bHQiDQogIGluY2x1ZGVfcmVzdWx0X3dpdGhfc2NvcmU6ICJwYXNzIg0KICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxDQpmaWx0ZXI6DQpvYmplY3RpdmVzOg0KICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1Ig0KICAgIGRpc3BsYXlOYW1lOiAiVGhpcyBpcyBteSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSBsb25nIGRpc3BsYXlOYW1lIg0KICAgIGtleV9zbGk6IGZhbHNlDQogICAgcGFzczogICAgICAgICAgICAgIyBwYXNzIGlmIChyZWxhdGl2ZSBjaGFuZ2UgPD0gMTAlIEFORCBhYnNvbHV0ZSB2YWx1ZSBpcyA8IDYwMG1zKQ0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykNCiAgICAgICAgICAtICI8NjAwIiAgICAjIGFic29sdXRlIHZhbHVlcyBvbmx5IHJlcXVpcmUgYSBsb2dpY2FsIG9wZXJhdG9yDQogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nDQogICAgICAtIGNyaXRlcmlhOg0KICAgICAgICAgIC0gIjw9ODAwIg0KICAgIHdlaWdodDogMQ0KdG90YWxfc2NvcmU6DQogIHBhc3M6ICI5MCUiDQogIHdhcm5pbmc6ICI3NSUi",
          "timeEnd": "2021-03-12T11:03:00.000Z",
          "timeStart": "2021-03-12T10:58:00.000Z"
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
      "id": "af121976-d640-4a41-9fdd-6a2887b7b796",
      "source": "lighthouse-service",
      "specversion": "1.0",
      "time": "2021-03-12T10:05:16.544Z",
      "type": "sh.keptn.event.evaluation.finished",
      "shkeptncontext": "76207fe5-5b72-40db-a74b-e220ae0892c2",
      "triggeredid": "6ada8320-8855-4e12-95c1-0b7fddc91e8d"
    },
    "heatmapLabel": "2021-03-12 11:05"
  },
  {
    "traces": [],
    "data": {
      "evaluation": {
        "comparedEvents": [
          "e7ca246d-b440-4d10-9005-04361caa5cda"
        ],
        "gitCommit": "",
        "indicatorResults": [
          {
            "keySli": false,
            "passTargets": [
              {
                "criteria": "<=+10%",
                "targetValue": 0,
                "violated": true
              },
              {
                "criteria": "<600",
                "targetValue": 0,
                "violated": true
              }
            ],
            "score": 0,
            "status": "fail",
            "value": {
              "message": "error validating time range: Supplied end-time 2021-03-12 11:05:00 +0000 UTC is too far (>120seconds) in the future (now: 2021-03-12 10:04:31.956048604 +0000 UTC m=+606126.844001381 - diff in sec: -3628.043951396)\n",
              "metric": "response_time_p95",
              "success": false,
              "value": 0
            },
            "warningTargets": [
              {
                "criteria": "<=800",
                "targetValue": 0,
                "violated": true
              }
            ]
          }
        ],
        "result": "fail",
        "score": 0,
        "sloFileContent": "LS0tDQpzcGVjX3ZlcnNpb246ICIxLjAiDQpjb21wYXJpc29uOg0KICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciDQogIGNvbXBhcmVfd2l0aDogInNpbmdsZV9yZXN1bHQiDQogIGluY2x1ZGVfcmVzdWx0X3dpdGhfc2NvcmU6ICJwYXNzIg0KICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxDQpmaWx0ZXI6DQpvYmplY3RpdmVzOg0KICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1Ig0KICAgIGRpc3BsYXlOYW1lOiAiVGhpcyBpcyBteSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSBsb25nIGRpc3BsYXlOYW1lIg0KICAgIGtleV9zbGk6IGZhbHNlDQogICAgcGFzczogICAgICAgICAgICAgIyBwYXNzIGlmIChyZWxhdGl2ZSBjaGFuZ2UgPD0gMTAlIEFORCBhYnNvbHV0ZSB2YWx1ZSBpcyA8IDYwMG1zKQ0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykNCiAgICAgICAgICAtICI8NjAwIiAgICAjIGFic29sdXRlIHZhbHVlcyBvbmx5IHJlcXVpcmUgYSBsb2dpY2FsIG9wZXJhdG9yDQogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nDQogICAgICAtIGNyaXRlcmlhOg0KICAgICAgICAgIC0gIjw9ODAwIg0KICAgIHdlaWdodDogMQ0KdG90YWxfc2NvcmU6DQogIHBhc3M6ICI5MCUiDQogIHdhcm5pbmc6ICI3NSUi",
        "timeEnd": "2021-03-12T11:05:00.000Z",
        "timeStart": "2021-03-12T11:00:00.000Z"
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
    "id": "4c1c2acd-dc6f-4008-82d7-b52e241631ac",
    "source": "lighthouse-service",
    "specversion": "1.0",
    "time": "2021-03-12T10:04:33.964Z",
    "type": "sh.keptn.event.evaluation.finished",
    "shkeptncontext": "f323e368-a40b-4d98-9098-56b1f56b4ee4",
    "triggeredid": "5f62799f-5d0d-4f5d-a3c7-b62b53b74375",
    "plainEvent": {
      "data": {
        "evaluation": {
          "comparedEvents": [
            "e7ca246d-b440-4d10-9005-04361caa5cda"
          ],
          "gitCommit": "",
          "indicatorResults": [
            {
              "keySli": false,
              "passTargets": [
                {
                  "criteria": "<=+10%",
                  "targetValue": 0,
                  "violated": true
                },
                {
                  "criteria": "<600",
                  "targetValue": 0,
                  "violated": true
                }
              ],
              "score": 0,
              "status": "fail",
              "value": {
                "message": "error validating time range: Supplied end-time 2021-03-12 11:05:00 +0000 UTC is too far (>120seconds) in the future (now: 2021-03-12 10:04:31.956048604 +0000 UTC m=+606126.844001381 - diff in sec: -3628.043951396)\n",
                "metric": "response_time_p95",
                "success": false,
                "value": 0
              },
              "warningTargets": [
                {
                  "criteria": "<=800",
                  "targetValue": 0,
                  "violated": true
                }
              ]
            }
          ],
          "result": "fail",
          "score": 0,
          "sloFileContent": "LS0tDQpzcGVjX3ZlcnNpb246ICIxLjAiDQpjb21wYXJpc29uOg0KICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciDQogIGNvbXBhcmVfd2l0aDogInNpbmdsZV9yZXN1bHQiDQogIGluY2x1ZGVfcmVzdWx0X3dpdGhfc2NvcmU6ICJwYXNzIg0KICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxDQpmaWx0ZXI6DQpvYmplY3RpdmVzOg0KICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1Ig0KICAgIGRpc3BsYXlOYW1lOiAiVGhpcyBpcyBteSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSBsb25nIGRpc3BsYXlOYW1lIg0KICAgIGtleV9zbGk6IGZhbHNlDQogICAgcGFzczogICAgICAgICAgICAgIyBwYXNzIGlmIChyZWxhdGl2ZSBjaGFuZ2UgPD0gMTAlIEFORCBhYnNvbHV0ZSB2YWx1ZSBpcyA8IDYwMG1zKQ0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykNCiAgICAgICAgICAtICI8NjAwIiAgICAjIGFic29sdXRlIHZhbHVlcyBvbmx5IHJlcXVpcmUgYSBsb2dpY2FsIG9wZXJhdG9yDQogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nDQogICAgICAtIGNyaXRlcmlhOg0KICAgICAgICAgIC0gIjw9ODAwIg0KICAgIHdlaWdodDogMQ0KdG90YWxfc2NvcmU6DQogIHBhc3M6ICI5MCUiDQogIHdhcm5pbmc6ICI3NSUi",
          "timeEnd": "2021-03-12T11:05:00.000Z",
          "timeStart": "2021-03-12T11:00:00.000Z"
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
      "id": "4c1c2acd-dc6f-4008-82d7-b52e241631ac",
      "source": "lighthouse-service",
      "specversion": "1.0",
      "time": "2021-03-12T10:04:33.964Z",
      "type": "sh.keptn.event.evaluation.finished",
      "shkeptncontext": "f323e368-a40b-4d98-9098-56b1f56b4ee4",
      "triggeredid": "5f62799f-5d0d-4f5d-a3c7-b62b53b74375"
    },
    "heatmapLabel": "2021-03-12 11:04 (2)"
  },
  {
    "traces": [],
    "data": {
      "evaluation": {
        "comparedEvents": [
          "e7ca246d-b440-4d10-9005-04361caa5cda"
        ],
        "gitCommit": "",
        "indicatorResults": [
          {
            "keySli": false,
            "passTargets": [
              {
                "criteria": "<=+10%",
                "targetValue": 0,
                "violated": true
              },
              {
                "criteria": "<600",
                "targetValue": 0,
                "violated": true
              }
            ],
            "score": 0,
            "status": "fail",
            "value": {
              "message": "error validating time range: Supplied end-time 2021-03-12 11:08:00 +0000 UTC is too far (>120seconds) in the future (now: 2021-03-12 10:04:01.163150505 +0000 UTC m=+606096.051103284 - diff in sec: -3838.836849495)\n",
              "metric": "response_time_p95",
              "success": false,
              "value": 0
            },
            "warningTargets": [
              {
                "criteria": "<=800",
                "targetValue": 0,
                "violated": true
              }
            ]
          }
        ],
        "result": "fail",
        "score": 0,
        "sloFileContent": "LS0tDQpzcGVjX3ZlcnNpb246ICIxLjAiDQpjb21wYXJpc29uOg0KICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciDQogIGNvbXBhcmVfd2l0aDogInNpbmdsZV9yZXN1bHQiDQogIGluY2x1ZGVfcmVzdWx0X3dpdGhfc2NvcmU6ICJwYXNzIg0KICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxDQpmaWx0ZXI6DQpvYmplY3RpdmVzOg0KICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1Ig0KICAgIGRpc3BsYXlOYW1lOiAiVGhpcyBpcyBteSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSBsb25nIGRpc3BsYXlOYW1lIg0KICAgIGtleV9zbGk6IGZhbHNlDQogICAgcGFzczogICAgICAgICAgICAgIyBwYXNzIGlmIChyZWxhdGl2ZSBjaGFuZ2UgPD0gMTAlIEFORCBhYnNvbHV0ZSB2YWx1ZSBpcyA8IDYwMG1zKQ0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykNCiAgICAgICAgICAtICI8NjAwIiAgICAjIGFic29sdXRlIHZhbHVlcyBvbmx5IHJlcXVpcmUgYSBsb2dpY2FsIG9wZXJhdG9yDQogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nDQogICAgICAtIGNyaXRlcmlhOg0KICAgICAgICAgIC0gIjw9ODAwIg0KICAgIHdlaWdodDogMQ0KdG90YWxfc2NvcmU6DQogIHBhc3M6ICI5MCUiDQogIHdhcm5pbmc6ICI3NSUi",
        "timeEnd": "2021-03-12T11:08:00.000Z",
        "timeStart": "2021-03-12T11:03:00.000Z"
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
    "id": "52d46191-76f3-4e06-bf37-afcd514a77e4",
    "source": "lighthouse-service",
    "specversion": "1.0",
    "time": "2021-03-12T10:04:03.449Z",
    "type": "sh.keptn.event.evaluation.finished",
    "shkeptncontext": "c1a905ba-712c-4351-ace0-86273004ff68",
    "triggeredid": "3e57f805-f24a-4505-81e8-adf66f7b3f88",
    "plainEvent": {
      "data": {
        "evaluation": {
          "comparedEvents": [
            "e7ca246d-b440-4d10-9005-04361caa5cda"
          ],
          "gitCommit": "",
          "indicatorResults": [
            {
              "keySli": false,
              "passTargets": [
                {
                  "criteria": "<=+10%",
                  "targetValue": 0,
                  "violated": true
                },
                {
                  "criteria": "<600",
                  "targetValue": 0,
                  "violated": true
                }
              ],
              "score": 0,
              "status": "fail",
              "value": {
                "message": "error validating time range: Supplied end-time 2021-03-12 11:08:00 +0000 UTC is too far (>120seconds) in the future (now: 2021-03-12 10:04:01.163150505 +0000 UTC m=+606096.051103284 - diff in sec: -3838.836849495)\n",
                "metric": "response_time_p95",
                "success": false,
                "value": 0
              },
              "warningTargets": [
                {
                  "criteria": "<=800",
                  "targetValue": 0,
                  "violated": true
                }
              ]
            }
          ],
          "result": "fail",
          "score": 0,
          "sloFileContent": "LS0tDQpzcGVjX3ZlcnNpb246ICIxLjAiDQpjb21wYXJpc29uOg0KICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciDQogIGNvbXBhcmVfd2l0aDogInNpbmdsZV9yZXN1bHQiDQogIGluY2x1ZGVfcmVzdWx0X3dpdGhfc2NvcmU6ICJwYXNzIg0KICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxDQpmaWx0ZXI6DQpvYmplY3RpdmVzOg0KICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1Ig0KICAgIGRpc3BsYXlOYW1lOiAiVGhpcyBpcyBteSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSBsb25nIGRpc3BsYXlOYW1lIg0KICAgIGtleV9zbGk6IGZhbHNlDQogICAgcGFzczogICAgICAgICAgICAgIyBwYXNzIGlmIChyZWxhdGl2ZSBjaGFuZ2UgPD0gMTAlIEFORCBhYnNvbHV0ZSB2YWx1ZSBpcyA8IDYwMG1zKQ0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykNCiAgICAgICAgICAtICI8NjAwIiAgICAjIGFic29sdXRlIHZhbHVlcyBvbmx5IHJlcXVpcmUgYSBsb2dpY2FsIG9wZXJhdG9yDQogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nDQogICAgICAtIGNyaXRlcmlhOg0KICAgICAgICAgIC0gIjw9ODAwIg0KICAgIHdlaWdodDogMQ0KdG90YWxfc2NvcmU6DQogIHBhc3M6ICI5MCUiDQogIHdhcm5pbmc6ICI3NSUi",
          "timeEnd": "2021-03-12T11:08:00.000Z",
          "timeStart": "2021-03-12T11:03:00.000Z"
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
      "id": "52d46191-76f3-4e06-bf37-afcd514a77e4",
      "source": "lighthouse-service",
      "specversion": "1.0",
      "time": "2021-03-12T10:04:03.449Z",
      "type": "sh.keptn.event.evaluation.finished",
      "shkeptncontext": "c1a905ba-712c-4351-ace0-86273004ff68",
      "triggeredid": "3e57f805-f24a-4505-81e8-adf66f7b3f88"
    },
    "heatmapLabel": "2021-03-12 11:04 (1)"
  },
  {
    "traces": [],
    "data": {
      "evaluation": {
        "comparedEvents": [
          "e7ca246d-b440-4d10-9005-04361caa5cda"
        ],
        "gitCommit": "",
        "indicatorResults": [
          {
            "keySli": false,
            "passTargets": [
              {
                "criteria": "<=+10%",
                "targetValue": 0,
                "violated": true
              },
              {
                "criteria": "<600",
                "targetValue": 0,
                "violated": true
              }
            ],
            "score": 0,
            "status": "fail",
            "value": {
              "message": "error validating time range: Supplied end-time 2021-03-12 11:08:00 +0000 UTC is too far (>120seconds) in the future (now: 2021-03-12 10:02:23.430074967 +0000 UTC m=+605998.318027744 - diff in sec: -3936.569925033)\n",
              "metric": "response_time_p95",
              "success": false,
              "value": 0
            },
            "warningTargets": [
              {
                "criteria": "<=800",
                "targetValue": 0,
                "violated": true
              }
            ]
          }
        ],
        "result": "fail",
        "score": 0,
        "sloFileContent": "LS0tDQpzcGVjX3ZlcnNpb246ICIxLjAiDQpjb21wYXJpc29uOg0KICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciDQogIGNvbXBhcmVfd2l0aDogInNpbmdsZV9yZXN1bHQiDQogIGluY2x1ZGVfcmVzdWx0X3dpdGhfc2NvcmU6ICJwYXNzIg0KICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxDQpmaWx0ZXI6DQpvYmplY3RpdmVzOg0KICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1Ig0KICAgIGRpc3BsYXlOYW1lOiAiVGhpcyBpcyBteSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSBsb25nIGRpc3BsYXlOYW1lIg0KICAgIGtleV9zbGk6IGZhbHNlDQogICAgcGFzczogICAgICAgICAgICAgIyBwYXNzIGlmIChyZWxhdGl2ZSBjaGFuZ2UgPD0gMTAlIEFORCBhYnNvbHV0ZSB2YWx1ZSBpcyA8IDYwMG1zKQ0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykNCiAgICAgICAgICAtICI8NjAwIiAgICAjIGFic29sdXRlIHZhbHVlcyBvbmx5IHJlcXVpcmUgYSBsb2dpY2FsIG9wZXJhdG9yDQogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nDQogICAgICAtIGNyaXRlcmlhOg0KICAgICAgICAgIC0gIjw9ODAwIg0KICAgIHdlaWdodDogMQ0KdG90YWxfc2NvcmU6DQogIHBhc3M6ICI5MCUiDQogIHdhcm5pbmc6ICI3NSUi",
        "timeEnd": "2021-03-12T11:08:00.000Z",
        "timeStart": "2021-03-12T11:03:00.000Z"
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
    "id": "075a4a74-825d-466e-b018-24df57cef52a",
    "source": "lighthouse-service",
    "specversion": "1.0",
    "time": "2021-03-12T10:02:24.197Z",
    "type": "sh.keptn.event.evaluation.finished",
    "shkeptncontext": "4db3f2c7-09cb-4c79-aff9-4fcd3377f8dc",
    "triggeredid": "d2e88968-5d91-4dab-a66b-289a301923c0",
    "plainEvent": {
      "data": {
        "evaluation": {
          "comparedEvents": [
            "e7ca246d-b440-4d10-9005-04361caa5cda"
          ],
          "gitCommit": "",
          "indicatorResults": [
            {
              "keySli": false,
              "passTargets": [
                {
                  "criteria": "<=+10%",
                  "targetValue": 0,
                  "violated": true
                },
                {
                  "criteria": "<600",
                  "targetValue": 0,
                  "violated": true
                }
              ],
              "score": 0,
              "status": "fail",
              "value": {
                "message": "error validating time range: Supplied end-time 2021-03-12 11:08:00 +0000 UTC is too far (>120seconds) in the future (now: 2021-03-12 10:02:23.430074967 +0000 UTC m=+605998.318027744 - diff in sec: -3936.569925033)\n",
                "metric": "response_time_p95",
                "success": false,
                "value": 0
              },
              "warningTargets": [
                {
                  "criteria": "<=800",
                  "targetValue": 0,
                  "violated": true
                }
              ]
            }
          ],
          "result": "fail",
          "score": 0,
          "sloFileContent": "LS0tDQpzcGVjX3ZlcnNpb246ICIxLjAiDQpjb21wYXJpc29uOg0KICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciDQogIGNvbXBhcmVfd2l0aDogInNpbmdsZV9yZXN1bHQiDQogIGluY2x1ZGVfcmVzdWx0X3dpdGhfc2NvcmU6ICJwYXNzIg0KICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxDQpmaWx0ZXI6DQpvYmplY3RpdmVzOg0KICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1Ig0KICAgIGRpc3BsYXlOYW1lOiAiVGhpcyBpcyBteSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSBsb25nIGRpc3BsYXlOYW1lIg0KICAgIGtleV9zbGk6IGZhbHNlDQogICAgcGFzczogICAgICAgICAgICAgIyBwYXNzIGlmIChyZWxhdGl2ZSBjaGFuZ2UgPD0gMTAlIEFORCBhYnNvbHV0ZSB2YWx1ZSBpcyA8IDYwMG1zKQ0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykNCiAgICAgICAgICAtICI8NjAwIiAgICAjIGFic29sdXRlIHZhbHVlcyBvbmx5IHJlcXVpcmUgYSBsb2dpY2FsIG9wZXJhdG9yDQogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nDQogICAgICAtIGNyaXRlcmlhOg0KICAgICAgICAgIC0gIjw9ODAwIg0KICAgIHdlaWdodDogMQ0KdG90YWxfc2NvcmU6DQogIHBhc3M6ICI5MCUiDQogIHdhcm5pbmc6ICI3NSUi",
          "timeEnd": "2021-03-12T11:08:00.000Z",
          "timeStart": "2021-03-12T11:03:00.000Z"
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
      "id": "075a4a74-825d-466e-b018-24df57cef52a",
      "source": "lighthouse-service",
      "specversion": "1.0",
      "time": "2021-03-12T10:02:24.197Z",
      "type": "sh.keptn.event.evaluation.finished",
      "shkeptncontext": "4db3f2c7-09cb-4c79-aff9-4fcd3377f8dc",
      "triggeredid": "d2e88968-5d91-4dab-a66b-289a301923c0"
    },
    "heatmapLabel": "2021-03-12 11:02"
  },
  {
    "traces": [],
    "data": {
      "evaluation": {
        "comparedEvents": [
          "e7ca246d-b440-4d10-9005-04361caa5cda"
        ],
        "gitCommit": "",
        "indicatorResults": [
          {
            "keySli": false,
            "passTargets": [
              {
                "criteria": "<=+10%",
                "targetValue": 0,
                "violated": true
              },
              {
                "criteria": "<600",
                "targetValue": 0,
                "violated": true
              }
            ],
            "score": 0,
            "status": "fail",
            "value": {
              "message": "error validating time range: Supplied end-time 2021-03-12 11:06:00 +0000 UTC is too far (>120seconds) in the future (now: 2021-03-12 10:01:05.701536734 +0000 UTC m=+605920.589489490 - diff in sec: -3894.298463266)\n",
              "metric": "response_time_p95",
              "success": false,
              "value": 0
            },
            "warningTargets": [
              {
                "criteria": "<=800",
                "targetValue": 0,
                "violated": true
              }
            ]
          }
        ],
        "result": "fail",
        "score": 0,
        "sloFileContent": "LS0tDQpzcGVjX3ZlcnNpb246ICIxLjAiDQpjb21wYXJpc29uOg0KICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciDQogIGNvbXBhcmVfd2l0aDogInNpbmdsZV9yZXN1bHQiDQogIGluY2x1ZGVfcmVzdWx0X3dpdGhfc2NvcmU6ICJwYXNzIg0KICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxDQpmaWx0ZXI6DQpvYmplY3RpdmVzOg0KICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1Ig0KICAgIGRpc3BsYXlOYW1lOiAiVGhpcyBpcyBteSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSBsb25nIGRpc3BsYXlOYW1lIg0KICAgIGtleV9zbGk6IGZhbHNlDQogICAgcGFzczogICAgICAgICAgICAgIyBwYXNzIGlmIChyZWxhdGl2ZSBjaGFuZ2UgPD0gMTAlIEFORCBhYnNvbHV0ZSB2YWx1ZSBpcyA8IDYwMG1zKQ0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykNCiAgICAgICAgICAtICI8NjAwIiAgICAjIGFic29sdXRlIHZhbHVlcyBvbmx5IHJlcXVpcmUgYSBsb2dpY2FsIG9wZXJhdG9yDQogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nDQogICAgICAtIGNyaXRlcmlhOg0KICAgICAgICAgIC0gIjw9ODAwIg0KICAgIHdlaWdodDogMQ0KdG90YWxfc2NvcmU6DQogIHBhc3M6ICI5MCUiDQogIHdhcm5pbmc6ICI3NSUi",
        "timeEnd": "2021-03-12T11:06:00.000Z",
        "timeStart": "2021-03-12T11:01:00.000Z"
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
    "id": "fe619205-2c2b-439e-a218-a7e3d50e7c1b",
    "source": "lighthouse-service",
    "specversion": "1.0",
    "time": "2021-03-12T10:01:06.449Z",
    "type": "sh.keptn.event.evaluation.finished",
    "shkeptncontext": "8eef8989-2955-4343-84b5-9e9fa1227307",
    "triggeredid": "01b43a3d-4672-4b80-baae-018b7014e419",
    "plainEvent": {
      "data": {
        "evaluation": {
          "comparedEvents": [
            "e7ca246d-b440-4d10-9005-04361caa5cda"
          ],
          "gitCommit": "",
          "indicatorResults": [
            {
              "keySli": false,
              "passTargets": [
                {
                  "criteria": "<=+10%",
                  "targetValue": 0,
                  "violated": true
                },
                {
                  "criteria": "<600",
                  "targetValue": 0,
                  "violated": true
                }
              ],
              "score": 0,
              "status": "fail",
              "value": {
                "message": "error validating time range: Supplied end-time 2021-03-12 11:06:00 +0000 UTC is too far (>120seconds) in the future (now: 2021-03-12 10:01:05.701536734 +0000 UTC m=+605920.589489490 - diff in sec: -3894.298463266)\n",
                "metric": "response_time_p95",
                "success": false,
                "value": 0
              },
              "warningTargets": [
                {
                  "criteria": "<=800",
                  "targetValue": 0,
                  "violated": true
                }
              ]
            }
          ],
          "result": "fail",
          "score": 0,
          "sloFileContent": "LS0tDQpzcGVjX3ZlcnNpb246ICIxLjAiDQpjb21wYXJpc29uOg0KICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciDQogIGNvbXBhcmVfd2l0aDogInNpbmdsZV9yZXN1bHQiDQogIGluY2x1ZGVfcmVzdWx0X3dpdGhfc2NvcmU6ICJwYXNzIg0KICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxDQpmaWx0ZXI6DQpvYmplY3RpdmVzOg0KICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1Ig0KICAgIGRpc3BsYXlOYW1lOiAiVGhpcyBpcyBteSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSBsb25nIGRpc3BsYXlOYW1lIg0KICAgIGtleV9zbGk6IGZhbHNlDQogICAgcGFzczogICAgICAgICAgICAgIyBwYXNzIGlmIChyZWxhdGl2ZSBjaGFuZ2UgPD0gMTAlIEFORCBhYnNvbHV0ZSB2YWx1ZSBpcyA8IDYwMG1zKQ0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykNCiAgICAgICAgICAtICI8NjAwIiAgICAjIGFic29sdXRlIHZhbHVlcyBvbmx5IHJlcXVpcmUgYSBsb2dpY2FsIG9wZXJhdG9yDQogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nDQogICAgICAtIGNyaXRlcmlhOg0KICAgICAgICAgIC0gIjw9ODAwIg0KICAgIHdlaWdodDogMQ0KdG90YWxfc2NvcmU6DQogIHBhc3M6ICI5MCUiDQogIHdhcm5pbmc6ICI3NSUi",
          "timeEnd": "2021-03-12T11:06:00.000Z",
          "timeStart": "2021-03-12T11:01:00.000Z"
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
      "id": "fe619205-2c2b-439e-a218-a7e3d50e7c1b",
      "source": "lighthouse-service",
      "specversion": "1.0",
      "time": "2021-03-12T10:01:06.449Z",
      "type": "sh.keptn.event.evaluation.finished",
      "shkeptncontext": "8eef8989-2955-4343-84b5-9e9fa1227307",
      "triggeredid": "01b43a3d-4672-4b80-baae-018b7014e419"
    },
    "heatmapLabel": "2021-03-12 11:01"
  },
  {
    "traces": [],
    "data": {
      "evaluation": {
        "comparedEvents": [
          "b7004bff-c52d-4fb2-9292-fc1cb9b19424"
        ],
        "gitCommit": "",
        "indicatorResults": [
          {
            "keySli": false,
            "passTargets": [
              {
                "criteria": "<=+10%",
                "targetValue": 412.4303120396647,
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
              "value": 394.4034352276371
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
        "sloFileContent": "LS0tDQpzcGVjX3ZlcnNpb246ICIxLjAiDQpjb21wYXJpc29uOg0KICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciDQogIGNvbXBhcmVfd2l0aDogInNpbmdsZV9yZXN1bHQiDQogIGluY2x1ZGVfcmVzdWx0X3dpdGhfc2NvcmU6ICJwYXNzIg0KICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxDQpmaWx0ZXI6DQpvYmplY3RpdmVzOg0KICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1Ig0KICAgIGRpc3BsYXlOYW1lOiAiVGhpcyBpcyBteSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSBsb25nIGRpc3BsYXlOYW1lIg0KICAgIGtleV9zbGk6IGZhbHNlDQogICAgcGFzczogICAgICAgICAgICAgIyBwYXNzIGlmIChyZWxhdGl2ZSBjaGFuZ2UgPD0gMTAlIEFORCBhYnNvbHV0ZSB2YWx1ZSBpcyA8IDYwMG1zKQ0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykNCiAgICAgICAgICAtICI8NjAwIiAgICAjIGFic29sdXRlIHZhbHVlcyBvbmx5IHJlcXVpcmUgYSBsb2dpY2FsIG9wZXJhdG9yDQogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nDQogICAgICAtIGNyaXRlcmlhOg0KICAgICAgICAgIC0gIjw9ODAwIg0KICAgIHdlaWdodDogMQ0KdG90YWxfc2NvcmU6DQogIHBhc3M6ICI5MCUiDQogIHdhcm5pbmc6ICI3NSUi",
        "timeEnd": "2021-03-11T14:54:34Z",
        "timeStart": "2021-03-11T14:52:31Z"
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
    "id": "e7ca246d-b440-4d10-9005-04361caa5cda",
    "source": "lighthouse-service",
    "specversion": "1.0",
    "time": "2021-03-11T14:55:40.689Z",
    "type": "sh.keptn.event.evaluation.finished",
    "shkeptncontext": "3061342b-03eb-485c-8586-de1b50c16451",
    "triggeredid": "125cae9c-b2ff-41e3-ab2a-1d694af0c252",
    "plainEvent": {
      "data": {
        "evaluation": {
          "comparedEvents": [
            "b7004bff-c52d-4fb2-9292-fc1cb9b19424"
          ],
          "gitCommit": "",
          "indicatorResults": [
            {
              "keySli": false,
              "passTargets": [
                {
                  "criteria": "<=+10%",
                  "targetValue": 412.4303120396647,
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
                "value": 394.4034352276371
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
          "sloFileContent": "LS0tDQpzcGVjX3ZlcnNpb246ICIxLjAiDQpjb21wYXJpc29uOg0KICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciDQogIGNvbXBhcmVfd2l0aDogInNpbmdsZV9yZXN1bHQiDQogIGluY2x1ZGVfcmVzdWx0X3dpdGhfc2NvcmU6ICJwYXNzIg0KICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxDQpmaWx0ZXI6DQpvYmplY3RpdmVzOg0KICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1Ig0KICAgIGRpc3BsYXlOYW1lOiAiVGhpcyBpcyBteSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSBsb25nIGRpc3BsYXlOYW1lIg0KICAgIGtleV9zbGk6IGZhbHNlDQogICAgcGFzczogICAgICAgICAgICAgIyBwYXNzIGlmIChyZWxhdGl2ZSBjaGFuZ2UgPD0gMTAlIEFORCBhYnNvbHV0ZSB2YWx1ZSBpcyA8IDYwMG1zKQ0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykNCiAgICAgICAgICAtICI8NjAwIiAgICAjIGFic29sdXRlIHZhbHVlcyBvbmx5IHJlcXVpcmUgYSBsb2dpY2FsIG9wZXJhdG9yDQogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nDQogICAgICAtIGNyaXRlcmlhOg0KICAgICAgICAgIC0gIjw9ODAwIg0KICAgIHdlaWdodDogMQ0KdG90YWxfc2NvcmU6DQogIHBhc3M6ICI5MCUiDQogIHdhcm5pbmc6ICI3NSUi",
          "timeEnd": "2021-03-11T14:54:34Z",
          "timeStart": "2021-03-11T14:52:31Z"
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
      "id": "e7ca246d-b440-4d10-9005-04361caa5cda",
      "source": "lighthouse-service",
      "specversion": "1.0",
      "time": "2021-03-11T14:55:40.689Z",
      "type": "sh.keptn.event.evaluation.finished",
      "shkeptncontext": "3061342b-03eb-485c-8586-de1b50c16451",
      "triggeredid": "125cae9c-b2ff-41e3-ab2a-1d694af0c252"
    },
    "heatmapLabel": "2021-03-11 15:55"
  },
  {
    "traces": [],
    "data": {
      "evaluation": {
        "comparedEvents": [
          "f43be3d7-a5c3-425c-9714-f15d9f507ab3"
        ],
        "gitCommit": "",
        "indicatorResults": [
          {
            "keySli": false,
            "passTargets": [
              {
                "criteria": "<=+10%",
                "targetValue": 393.9961174456488,
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
              "value": 374.9366473087861
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
        "sloFileContent": "LS0tDQpzcGVjX3ZlcnNpb246ICIxLjAiDQpjb21wYXJpc29uOg0KICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciDQogIGNvbXBhcmVfd2l0aDogInNpbmdsZV9yZXN1bHQiDQogIGluY2x1ZGVfcmVzdWx0X3dpdGhfc2NvcmU6ICJwYXNzIg0KICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxDQpmaWx0ZXI6DQpvYmplY3RpdmVzOg0KICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1Ig0KICAgIGRpc3BsYXlOYW1lOiAiVGhpcyBpcyBteSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSBsb25nIGRpc3BsYXlOYW1lIg0KICAgIGtleV9zbGk6IGZhbHNlDQogICAgcGFzczogICAgICAgICAgICAgIyBwYXNzIGlmIChyZWxhdGl2ZSBjaGFuZ2UgPD0gMTAlIEFORCBhYnNvbHV0ZSB2YWx1ZSBpcyA8IDYwMG1zKQ0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykNCiAgICAgICAgICAtICI8NjAwIiAgICAjIGFic29sdXRlIHZhbHVlcyBvbmx5IHJlcXVpcmUgYSBsb2dpY2FsIG9wZXJhdG9yDQogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nDQogICAgICAtIGNyaXRlcmlhOg0KICAgICAgICAgIC0gIjw9ODAwIg0KICAgIHdlaWdodDogMQ0KdG90YWxfc2NvcmU6DQogIHBhc3M6ICI5MCUiDQogIHdhcm5pbmc6ICI3NSUi",
        "timeEnd": "2021-03-11T08:39:33Z",
        "timeStart": "2021-03-11T08:37:31Z"
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
    "id": "b7004bff-c52d-4fb2-9292-fc1cb9b19424",
    "source": "lighthouse-service",
    "specversion": "1.0",
    "time": "2021-03-11T08:40:39.649Z",
    "type": "sh.keptn.event.evaluation.finished",
    "shkeptncontext": "69978125-6132-438e-ac1d-bf46a6269ccc",
    "triggeredid": "e78b2138-26f7-4fdb-af2b-f13c729b8470",
    "plainEvent": {
      "data": {
        "evaluation": {
          "comparedEvents": [
            "f43be3d7-a5c3-425c-9714-f15d9f507ab3"
          ],
          "gitCommit": "",
          "indicatorResults": [
            {
              "keySli": false,
              "passTargets": [
                {
                  "criteria": "<=+10%",
                  "targetValue": 393.9961174456488,
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
                "value": 374.9366473087861
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
          "sloFileContent": "LS0tDQpzcGVjX3ZlcnNpb246ICIxLjAiDQpjb21wYXJpc29uOg0KICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciDQogIGNvbXBhcmVfd2l0aDogInNpbmdsZV9yZXN1bHQiDQogIGluY2x1ZGVfcmVzdWx0X3dpdGhfc2NvcmU6ICJwYXNzIg0KICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxDQpmaWx0ZXI6DQpvYmplY3RpdmVzOg0KICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1Ig0KICAgIGRpc3BsYXlOYW1lOiAiVGhpcyBpcyBteSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSBsb25nIGRpc3BsYXlOYW1lIg0KICAgIGtleV9zbGk6IGZhbHNlDQogICAgcGFzczogICAgICAgICAgICAgIyBwYXNzIGlmIChyZWxhdGl2ZSBjaGFuZ2UgPD0gMTAlIEFORCBhYnNvbHV0ZSB2YWx1ZSBpcyA8IDYwMG1zKQ0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykNCiAgICAgICAgICAtICI8NjAwIiAgICAjIGFic29sdXRlIHZhbHVlcyBvbmx5IHJlcXVpcmUgYSBsb2dpY2FsIG9wZXJhdG9yDQogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nDQogICAgICAtIGNyaXRlcmlhOg0KICAgICAgICAgIC0gIjw9ODAwIg0KICAgIHdlaWdodDogMQ0KdG90YWxfc2NvcmU6DQogIHBhc3M6ICI5MCUiDQogIHdhcm5pbmc6ICI3NSUi",
          "timeEnd": "2021-03-11T08:39:33Z",
          "timeStart": "2021-03-11T08:37:31Z"
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
      "id": "b7004bff-c52d-4fb2-9292-fc1cb9b19424",
      "source": "lighthouse-service",
      "specversion": "1.0",
      "time": "2021-03-11T08:40:39.649Z",
      "type": "sh.keptn.event.evaluation.finished",
      "shkeptncontext": "69978125-6132-438e-ac1d-bf46a6269ccc",
      "triggeredid": "e78b2138-26f7-4fdb-af2b-f13c729b8470"
    },
    "heatmapLabel": "2021-03-11 09:40"
  },
  {
    "traces": [],
    "data": {
      "evaluation": {
        "comparedEvents": [
          "f338f850-808f-4992-87e8-6fd7b275e186"
        ],
        "gitCommit": "",
        "indicatorResults": [
          {
            "displayName": "This is my very very very very very very very very very very very very long displayName",
            "keySli": false,
            "passTargets": [
              {
                "criteria": "<=+10%",
                "targetValue": 419.7260685366458,
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
              "value": 358.17828858695344
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
        "sloFileContent": "LS0tDQpzcGVjX3ZlcnNpb246ICIxLjAiDQpjb21wYXJpc29uOg0KICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciDQogIGNvbXBhcmVfd2l0aDogInNpbmdsZV9yZXN1bHQiDQogIGluY2x1ZGVfcmVzdWx0X3dpdGhfc2NvcmU6ICJwYXNzIg0KICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxDQpmaWx0ZXI6DQpvYmplY3RpdmVzOg0KICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1Ig0KICAgIGRpc3BsYXlOYW1lOiAiVGhpcyBpcyBteSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSBsb25nIGRpc3BsYXlOYW1lIg0KICAgIGtleV9zbGk6IGZhbHNlDQogICAgcGFzczogICAgICAgICAgICAgIyBwYXNzIGlmIChyZWxhdGl2ZSBjaGFuZ2UgPD0gMTAlIEFORCBhYnNvbHV0ZSB2YWx1ZSBpcyA8IDYwMG1zKQ0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykNCiAgICAgICAgICAtICI8NjAwIiAgICAjIGFic29sdXRlIHZhbHVlcyBvbmx5IHJlcXVpcmUgYSBsb2dpY2FsIG9wZXJhdG9yDQogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nDQogICAgICAtIGNyaXRlcmlhOg0KICAgICAgICAgIC0gIjw9ODAwIg0KICAgIHdlaWdodDogMQ0KdG90YWxfc2NvcmU6DQogIHBhc3M6ICI5MCUiDQogIHdhcm5pbmc6ICI3NSUi",
        "timeEnd": "2021-03-05T09:55:14Z",
        "timeStart": "2021-03-05T09:53:17Z"
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
    "id": "f43be3d7-a5c3-425c-9714-f15d9f507ab3",
    "source": "lighthouse-service",
    "specversion": "1.0",
    "time": "2021-03-05T09:57:26.779Z",
    "type": "sh.keptn.event.evaluation.finished",
    "shkeptncontext": "7107b223-2102-406d-b2dc-c73b92471840",
    "triggeredid": "052a0dce-1841-460f-a604-2a04dfcdc18e",
    "plainEvent": {
      "data": {
        "evaluation": {
          "comparedEvents": [
            "f338f850-808f-4992-87e8-6fd7b275e186"
          ],
          "gitCommit": "",
          "indicatorResults": [
            {
              "displayName": "This is my very very very very very very very very very very very very long displayName",
              "keySli": false,
              "passTargets": [
                {
                  "criteria": "<=+10%",
                  "targetValue": 419.7260685366458,
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
                "value": 358.17828858695344
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
          "sloFileContent": "LS0tDQpzcGVjX3ZlcnNpb246ICIxLjAiDQpjb21wYXJpc29uOg0KICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciDQogIGNvbXBhcmVfd2l0aDogInNpbmdsZV9yZXN1bHQiDQogIGluY2x1ZGVfcmVzdWx0X3dpdGhfc2NvcmU6ICJwYXNzIg0KICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxDQpmaWx0ZXI6DQpvYmplY3RpdmVzOg0KICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1Ig0KICAgIGRpc3BsYXlOYW1lOiAiVGhpcyBpcyBteSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSBsb25nIGRpc3BsYXlOYW1lIg0KICAgIGtleV9zbGk6IGZhbHNlDQogICAgcGFzczogICAgICAgICAgICAgIyBwYXNzIGlmIChyZWxhdGl2ZSBjaGFuZ2UgPD0gMTAlIEFORCBhYnNvbHV0ZSB2YWx1ZSBpcyA8IDYwMG1zKQ0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykNCiAgICAgICAgICAtICI8NjAwIiAgICAjIGFic29sdXRlIHZhbHVlcyBvbmx5IHJlcXVpcmUgYSBsb2dpY2FsIG9wZXJhdG9yDQogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nDQogICAgICAtIGNyaXRlcmlhOg0KICAgICAgICAgIC0gIjw9ODAwIg0KICAgIHdlaWdodDogMQ0KdG90YWxfc2NvcmU6DQogIHBhc3M6ICI5MCUiDQogIHdhcm5pbmc6ICI3NSUi",
          "timeEnd": "2021-03-05T09:55:14Z",
          "timeStart": "2021-03-05T09:53:17Z"
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
      "id": "f43be3d7-a5c3-425c-9714-f15d9f507ab3",
      "source": "lighthouse-service",
      "specversion": "1.0",
      "time": "2021-03-05T09:57:26.779Z",
      "type": "sh.keptn.event.evaluation.finished",
      "shkeptncontext": "7107b223-2102-406d-b2dc-c73b92471840",
      "triggeredid": "052a0dce-1841-460f-a604-2a04dfcdc18e"
    },
    "heatmapLabel": "2021-03-05 10:57"
  },
  {
    "traces": [],
    "data": {
      "evaluation": {
        "comparedEvents": [
          "f338f850-808f-4992-87e8-6fd7b275e186"
        ],
        "gitCommit": "",
        "indicatorResults": [
          {
            "displayName": "This is my very very very very very very very very very very very very long displayName",
            "keySli": false,
            "passTargets": [
              {
                "criteria": "<=+10%",
                "targetValue": 0,
                "violated": true
              },
              {
                "criteria": "<600",
                "targetValue": 0,
                "violated": true
              }
            ],
            "score": 0,
            "status": "fail",
            "value": {
              "message": "error validating time range: Supplied end-time 2021-03-04 14:56:40 +0000 UTC is too far (>120seconds) in the future (now: 2021-03-04 14:51:50.977607608 +0000 UTC m=+264277.975640644 - diff in sec: -289.022392392)\n",
              "metric": "response_time_p95",
              "success": false,
              "value": 0
            },
            "warningTargets": [
              {
                "criteria": "<=800",
                "targetValue": 0,
                "violated": true
              }
            ]
          }
        ],
        "result": "fail",
        "score": 0,
        "sloFileContent": "LS0tDQpzcGVjX3ZlcnNpb246ICIxLjAiDQpjb21wYXJpc29uOg0KICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciDQogIGNvbXBhcmVfd2l0aDogInNpbmdsZV9yZXN1bHQiDQogIGluY2x1ZGVfcmVzdWx0X3dpdGhfc2NvcmU6ICJwYXNzIg0KICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxDQpmaWx0ZXI6DQpvYmplY3RpdmVzOg0KICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1Ig0KICAgIGRpc3BsYXlOYW1lOiAiVGhpcyBpcyBteSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSBsb25nIGRpc3BsYXlOYW1lIg0KICAgIGtleV9zbGk6IGZhbHNlDQogICAgcGFzczogICAgICAgICAgICAgIyBwYXNzIGlmIChyZWxhdGl2ZSBjaGFuZ2UgPD0gMTAlIEFORCBhYnNvbHV0ZSB2YWx1ZSBpcyA8IDYwMG1zKQ0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykNCiAgICAgICAgICAtICI8NjAwIiAgICAjIGFic29sdXRlIHZhbHVlcyBvbmx5IHJlcXVpcmUgYSBsb2dpY2FsIG9wZXJhdG9yDQogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nDQogICAgICAtIGNyaXRlcmlhOg0KICAgICAgICAgIC0gIjw9ODAwIg0KICAgIHdlaWdodDogMQ0KdG90YWxfc2NvcmU6DQogIHBhc3M6ICI5MCUiDQogIHdhcm5pbmc6ICI3NSUi",
        "timeEnd": "2021-03-04T14:56:40.000Z",
        "timeStart": "2021-03-04T14:51:40.000Z"
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
    "id": "f5380b65-2f42-4a17-83a6-bb235ebcf50c",
    "source": "lighthouse-service",
    "specversion": "1.0",
    "time": "2021-03-04T14:51:51.717Z",
    "type": "sh.keptn.event.evaluation.finished",
    "shkeptncontext": "281f8faf-f86e-4790-a334-900dc1600846",
    "triggeredid": "ad752e0c-8241-43b1-9432-70c8eb86d349",
    "plainEvent": {
      "data": {
        "evaluation": {
          "comparedEvents": [
            "f338f850-808f-4992-87e8-6fd7b275e186"
          ],
          "gitCommit": "",
          "indicatorResults": [
            {
              "displayName": "This is my very very very very very very very very very very very very long displayName",
              "keySli": false,
              "passTargets": [
                {
                  "criteria": "<=+10%",
                  "targetValue": 0,
                  "violated": true
                },
                {
                  "criteria": "<600",
                  "targetValue": 0,
                  "violated": true
                }
              ],
              "score": 0,
              "status": "fail",
              "value": {
                "message": "error validating time range: Supplied end-time 2021-03-04 14:56:40 +0000 UTC is too far (>120seconds) in the future (now: 2021-03-04 14:51:50.977607608 +0000 UTC m=+264277.975640644 - diff in sec: -289.022392392)\n",
                "metric": "response_time_p95",
                "success": false,
                "value": 0
              },
              "warningTargets": [
                {
                  "criteria": "<=800",
                  "targetValue": 0,
                  "violated": true
                }
              ]
            }
          ],
          "result": "fail",
          "score": 0,
          "sloFileContent": "LS0tDQpzcGVjX3ZlcnNpb246ICIxLjAiDQpjb21wYXJpc29uOg0KICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciDQogIGNvbXBhcmVfd2l0aDogInNpbmdsZV9yZXN1bHQiDQogIGluY2x1ZGVfcmVzdWx0X3dpdGhfc2NvcmU6ICJwYXNzIg0KICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxDQpmaWx0ZXI6DQpvYmplY3RpdmVzOg0KICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1Ig0KICAgIGRpc3BsYXlOYW1lOiAiVGhpcyBpcyBteSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSBsb25nIGRpc3BsYXlOYW1lIg0KICAgIGtleV9zbGk6IGZhbHNlDQogICAgcGFzczogICAgICAgICAgICAgIyBwYXNzIGlmIChyZWxhdGl2ZSBjaGFuZ2UgPD0gMTAlIEFORCBhYnNvbHV0ZSB2YWx1ZSBpcyA8IDYwMG1zKQ0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykNCiAgICAgICAgICAtICI8NjAwIiAgICAjIGFic29sdXRlIHZhbHVlcyBvbmx5IHJlcXVpcmUgYSBsb2dpY2FsIG9wZXJhdG9yDQogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nDQogICAgICAtIGNyaXRlcmlhOg0KICAgICAgICAgIC0gIjw9ODAwIg0KICAgIHdlaWdodDogMQ0KdG90YWxfc2NvcmU6DQogIHBhc3M6ICI5MCUiDQogIHdhcm5pbmc6ICI3NSUi",
          "timeEnd": "2021-03-04T14:56:40.000Z",
          "timeStart": "2021-03-04T14:51:40.000Z"
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
      "id": "f5380b65-2f42-4a17-83a6-bb235ebcf50c",
      "source": "lighthouse-service",
      "specversion": "1.0",
      "time": "2021-03-04T14:51:51.717Z",
      "type": "sh.keptn.event.evaluation.finished",
      "shkeptncontext": "281f8faf-f86e-4790-a334-900dc1600846",
      "triggeredid": "ad752e0c-8241-43b1-9432-70c8eb86d349"
    },
    "heatmapLabel": "2021-03-04 15:51"
  },
  {
    "traces": [],
    "data": {
      "evaluation": {
        "comparedEvents": [
          "f338f850-808f-4992-87e8-6fd7b275e186"
        ],
        "gitCommit": "",
        "indicatorResults": [
          {
            "displayName": "This is my very very very very very very very very very very very very long displayName",
            "keySli": false,
            "passTargets": [
              {
                "criteria": "<=+10%",
                "targetValue": 419.7260685366458,
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
              "value": 1076.4413638443743
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
        "sloFileContent": "LS0tDQpzcGVjX3ZlcnNpb246ICIxLjAiDQpjb21wYXJpc29uOg0KICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciDQogIGNvbXBhcmVfd2l0aDogInNpbmdsZV9yZXN1bHQiDQogIGluY2x1ZGVfcmVzdWx0X3dpdGhfc2NvcmU6ICJwYXNzIg0KICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxDQpmaWx0ZXI6DQpvYmplY3RpdmVzOg0KICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1Ig0KICAgIGRpc3BsYXlOYW1lOiAiVGhpcyBpcyBteSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSBsb25nIGRpc3BsYXlOYW1lIg0KICAgIGtleV9zbGk6IGZhbHNlDQogICAgcGFzczogICAgICAgICAgICAgIyBwYXNzIGlmIChyZWxhdGl2ZSBjaGFuZ2UgPD0gMTAlIEFORCBhYnNvbHV0ZSB2YWx1ZSBpcyA8IDYwMG1zKQ0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykNCiAgICAgICAgICAtICI8NjAwIiAgICAjIGFic29sdXRlIHZhbHVlcyBvbmx5IHJlcXVpcmUgYSBsb2dpY2FsIG9wZXJhdG9yDQogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nDQogICAgICAtIGNyaXRlcmlhOg0KICAgICAgICAgIC0gIjw9ODAwIg0KICAgIHdlaWdodDogMQ0KdG90YWxfc2NvcmU6DQogIHBhc3M6ICI5MCUiDQogIHdhcm5pbmc6ICI3NSUi",
        "timeEnd": "2021-03-04T14:29:19Z",
        "timeStart": "2021-03-04T14:20:32Z"
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
    "id": "8628a025-1812-48ed-95c1-decb879396ef",
    "source": "lighthouse-service",
    "specversion": "1.0",
    "time": "2021-03-04T14:29:26.223Z",
    "type": "sh.keptn.event.evaluation.finished",
    "shkeptncontext": "2246a327-09c9-4e6e-9e9b-1adce44cb45b",
    "triggeredid": "8f8f5f4f-4c21-412a-928a-70d98fde530e",
    "plainEvent": {
      "data": {
        "evaluation": {
          "comparedEvents": [
            "f338f850-808f-4992-87e8-6fd7b275e186"
          ],
          "gitCommit": "",
          "indicatorResults": [
            {
              "displayName": "This is my very very very very very very very very very very very very long displayName",
              "keySli": false,
              "passTargets": [
                {
                  "criteria": "<=+10%",
                  "targetValue": 419.7260685366458,
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
                "value": 1076.4413638443743
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
          "sloFileContent": "LS0tDQpzcGVjX3ZlcnNpb246ICIxLjAiDQpjb21wYXJpc29uOg0KICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciDQogIGNvbXBhcmVfd2l0aDogInNpbmdsZV9yZXN1bHQiDQogIGluY2x1ZGVfcmVzdWx0X3dpdGhfc2NvcmU6ICJwYXNzIg0KICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxDQpmaWx0ZXI6DQpvYmplY3RpdmVzOg0KICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1Ig0KICAgIGRpc3BsYXlOYW1lOiAiVGhpcyBpcyBteSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSBsb25nIGRpc3BsYXlOYW1lIg0KICAgIGtleV9zbGk6IGZhbHNlDQogICAgcGFzczogICAgICAgICAgICAgIyBwYXNzIGlmIChyZWxhdGl2ZSBjaGFuZ2UgPD0gMTAlIEFORCBhYnNvbHV0ZSB2YWx1ZSBpcyA8IDYwMG1zKQ0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykNCiAgICAgICAgICAtICI8NjAwIiAgICAjIGFic29sdXRlIHZhbHVlcyBvbmx5IHJlcXVpcmUgYSBsb2dpY2FsIG9wZXJhdG9yDQogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nDQogICAgICAtIGNyaXRlcmlhOg0KICAgICAgICAgIC0gIjw9ODAwIg0KICAgIHdlaWdodDogMQ0KdG90YWxfc2NvcmU6DQogIHBhc3M6ICI5MCUiDQogIHdhcm5pbmc6ICI3NSUi",
          "timeEnd": "2021-03-04T14:29:19Z",
          "timeStart": "2021-03-04T14:20:32Z"
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
      "id": "8628a025-1812-48ed-95c1-decb879396ef",
      "source": "lighthouse-service",
      "specversion": "1.0",
      "time": "2021-03-04T14:29:26.223Z",
      "type": "sh.keptn.event.evaluation.finished",
      "shkeptncontext": "2246a327-09c9-4e6e-9e9b-1adce44cb45b",
      "triggeredid": "8f8f5f4f-4c21-412a-928a-70d98fde530e"
    },
    "heatmapLabel": "2021-03-04 15:29"
  },
  {
    "traces": [],
    "data": {
      "evaluation": {
        "comparedEvents": [
          "1f82880e-4edc-4170-a53e-3bf1a83c8e89"
        ],
        "gitCommit": "",
        "indicatorResults": [
          {
            "displayName": "This is my very very very very very very very very very very very very long displayName",
            "keySli": false,
            "passTargets": [
              {
                "criteria": "<=+10%",
                "targetValue": 414.56522943851104,
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
              "value": 381.5691532151326
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
        "sloFileContent": "LS0tDQpzcGVjX3ZlcnNpb246ICIxLjAiDQpjb21wYXJpc29uOg0KICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciDQogIGNvbXBhcmVfd2l0aDogInNpbmdsZV9yZXN1bHQiDQogIGluY2x1ZGVfcmVzdWx0X3dpdGhfc2NvcmU6ICJwYXNzIg0KICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxDQpmaWx0ZXI6DQpvYmplY3RpdmVzOg0KICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1Ig0KICAgIGRpc3BsYXlOYW1lOiAiVGhpcyBpcyBteSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSBsb25nIGRpc3BsYXlOYW1lIg0KICAgIGtleV9zbGk6IGZhbHNlDQogICAgcGFzczogICAgICAgICAgICAgIyBwYXNzIGlmIChyZWxhdGl2ZSBjaGFuZ2UgPD0gMTAlIEFORCBhYnNvbHV0ZSB2YWx1ZSBpcyA8IDYwMG1zKQ0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykNCiAgICAgICAgICAtICI8NjAwIiAgICAjIGFic29sdXRlIHZhbHVlcyBvbmx5IHJlcXVpcmUgYSBsb2dpY2FsIG9wZXJhdG9yDQogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nDQogICAgICAtIGNyaXRlcmlhOg0KICAgICAgICAgIC0gIjw9ODAwIg0KICAgIHdlaWdodDogMQ0KdG90YWxfc2NvcmU6DQogIHBhc3M6ICI5MCUiDQogIHdhcm5pbmc6ICI3NSUi",
        "timeEnd": "2021-03-02T13:00:06Z",
        "timeStart": "2021-03-02T12:58:11Z"
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
    "id": "f338f850-808f-4992-87e8-6fd7b275e186",
    "source": "lighthouse-service",
    "specversion": "1.0",
    "time": "2021-03-02T13:02:13.230Z",
    "type": "sh.keptn.event.evaluation.finished",
    "shkeptncontext": "702c32a2-034f-4ac2-b7c3-903088fedb5d",
    "triggeredid": "b981b276-8cd4-44b2-a590-77faa2a60736",
    "plainEvent": {
      "data": {
        "evaluation": {
          "comparedEvents": [
            "1f82880e-4edc-4170-a53e-3bf1a83c8e89"
          ],
          "gitCommit": "",
          "indicatorResults": [
            {
              "displayName": "This is my very very very very very very very very very very very very long displayName",
              "keySli": false,
              "passTargets": [
                {
                  "criteria": "<=+10%",
                  "targetValue": 414.56522943851104,
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
                "value": 381.5691532151326
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
          "sloFileContent": "LS0tDQpzcGVjX3ZlcnNpb246ICIxLjAiDQpjb21wYXJpc29uOg0KICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciDQogIGNvbXBhcmVfd2l0aDogInNpbmdsZV9yZXN1bHQiDQogIGluY2x1ZGVfcmVzdWx0X3dpdGhfc2NvcmU6ICJwYXNzIg0KICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxDQpmaWx0ZXI6DQpvYmplY3RpdmVzOg0KICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1Ig0KICAgIGRpc3BsYXlOYW1lOiAiVGhpcyBpcyBteSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSBsb25nIGRpc3BsYXlOYW1lIg0KICAgIGtleV9zbGk6IGZhbHNlDQogICAgcGFzczogICAgICAgICAgICAgIyBwYXNzIGlmIChyZWxhdGl2ZSBjaGFuZ2UgPD0gMTAlIEFORCBhYnNvbHV0ZSB2YWx1ZSBpcyA8IDYwMG1zKQ0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykNCiAgICAgICAgICAtICI8NjAwIiAgICAjIGFic29sdXRlIHZhbHVlcyBvbmx5IHJlcXVpcmUgYSBsb2dpY2FsIG9wZXJhdG9yDQogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nDQogICAgICAtIGNyaXRlcmlhOg0KICAgICAgICAgIC0gIjw9ODAwIg0KICAgIHdlaWdodDogMQ0KdG90YWxfc2NvcmU6DQogIHBhc3M6ICI5MCUiDQogIHdhcm5pbmc6ICI3NSUi",
          "timeEnd": "2021-03-02T13:00:06Z",
          "timeStart": "2021-03-02T12:58:11Z"
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
      "id": "f338f850-808f-4992-87e8-6fd7b275e186",
      "source": "lighthouse-service",
      "specversion": "1.0",
      "time": "2021-03-02T13:02:13.230Z",
      "type": "sh.keptn.event.evaluation.finished",
      "shkeptncontext": "702c32a2-034f-4ac2-b7c3-903088fedb5d",
      "triggeredid": "b981b276-8cd4-44b2-a590-77faa2a60736"
    },
    "heatmapLabel": "2021-03-02 14:02"
  },
  {
    "traces": [],
    "data": {
      "evaluation": {
        "comparedEvents": [
          "90200421-420a-4087-b86b-12185d702237"
        ],
        "gitCommit": "",
        "indicatorResults": [
          {
            "displayName": "This is my very very very very very very very very very very very very long displayName",
            "keySli": false,
            "passTargets": [
              {
                "criteria": "<=+10%",
                "targetValue": 0,
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
              "value": 376.8774813077373
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
        "sloFileContent": "LS0tDQpzcGVjX3ZlcnNpb246ICIxLjAiDQpjb21wYXJpc29uOg0KICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciDQogIGNvbXBhcmVfd2l0aDogInNpbmdsZV9yZXN1bHQiDQogIGluY2x1ZGVfcmVzdWx0X3dpdGhfc2NvcmU6ICJwYXNzIg0KICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxDQpmaWx0ZXI6DQpvYmplY3RpdmVzOg0KICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1Ig0KICAgIGRpc3BsYXlOYW1lOiAiVGhpcyBpcyBteSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSBsb25nIGRpc3BsYXlOYW1lIg0KICAgIGtleV9zbGk6IGZhbHNlDQogICAgcGFzczogICAgICAgICAgICAgIyBwYXNzIGlmIChyZWxhdGl2ZSBjaGFuZ2UgPD0gMTAlIEFORCBhYnNvbHV0ZSB2YWx1ZSBpcyA8IDYwMG1zKQ0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykNCiAgICAgICAgICAtICI8NjAwIiAgICAjIGFic29sdXRlIHZhbHVlcyBvbmx5IHJlcXVpcmUgYSBsb2dpY2FsIG9wZXJhdG9yDQogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nDQogICAgICAtIGNyaXRlcmlhOg0KICAgICAgICAgIC0gIjw9ODAwIg0KICAgIHdlaWdodDogMQ0KdG90YWxfc2NvcmU6DQogIHBhc3M6ICI5MCUiDQogIHdhcm5pbmc6ICI3NSUi",
        "timeEnd": "2021-03-02T11:29:31Z",
        "timeStart": "2021-03-02T11:27:34Z"
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
    "id": "1f82880e-4edc-4170-a53e-3bf1a83c8e89",
    "source": "lighthouse-service",
    "specversion": "1.0",
    "time": "2021-03-02T11:31:37.843Z",
    "type": "sh.keptn.event.evaluation.finished",
    "shkeptncontext": "cdb13092-3281-4a7c-84c9-db11e7af6930",
    "triggeredid": "5b4730e8-34e7-46e7-9710-c63eef9ee847",
    "plainEvent": {
      "data": {
        "evaluation": {
          "comparedEvents": [
            "90200421-420a-4087-b86b-12185d702237"
          ],
          "gitCommit": "",
          "indicatorResults": [
            {
              "displayName": "This is my very very very very very very very very very very very very long displayName",
              "keySli": false,
              "passTargets": [
                {
                  "criteria": "<=+10%",
                  "targetValue": 0,
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
                "value": 376.8774813077373
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
          "sloFileContent": "LS0tDQpzcGVjX3ZlcnNpb246ICIxLjAiDQpjb21wYXJpc29uOg0KICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciDQogIGNvbXBhcmVfd2l0aDogInNpbmdsZV9yZXN1bHQiDQogIGluY2x1ZGVfcmVzdWx0X3dpdGhfc2NvcmU6ICJwYXNzIg0KICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxDQpmaWx0ZXI6DQpvYmplY3RpdmVzOg0KICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1Ig0KICAgIGRpc3BsYXlOYW1lOiAiVGhpcyBpcyBteSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSBsb25nIGRpc3BsYXlOYW1lIg0KICAgIGtleV9zbGk6IGZhbHNlDQogICAgcGFzczogICAgICAgICAgICAgIyBwYXNzIGlmIChyZWxhdGl2ZSBjaGFuZ2UgPD0gMTAlIEFORCBhYnNvbHV0ZSB2YWx1ZSBpcyA8IDYwMG1zKQ0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykNCiAgICAgICAgICAtICI8NjAwIiAgICAjIGFic29sdXRlIHZhbHVlcyBvbmx5IHJlcXVpcmUgYSBsb2dpY2FsIG9wZXJhdG9yDQogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nDQogICAgICAtIGNyaXRlcmlhOg0KICAgICAgICAgIC0gIjw9ODAwIg0KICAgIHdlaWdodDogMQ0KdG90YWxfc2NvcmU6DQogIHBhc3M6ICI5MCUiDQogIHdhcm5pbmc6ICI3NSUi",
          "timeEnd": "2021-03-02T11:29:31Z",
          "timeStart": "2021-03-02T11:27:34Z"
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
      "id": "1f82880e-4edc-4170-a53e-3bf1a83c8e89",
      "source": "lighthouse-service",
      "specversion": "1.0",
      "time": "2021-03-02T11:31:37.843Z",
      "type": "sh.keptn.event.evaluation.finished",
      "shkeptncontext": "cdb13092-3281-4a7c-84c9-db11e7af6930",
      "triggeredid": "5b4730e8-34e7-46e7-9710-c63eef9ee847"
    },
    "heatmapLabel": "2021-03-02 12:31"
  },
  {
    "traces": [],
    "data": {
      "evaluation": {
        "gitCommit": "",
        "indicatorResults": null,
        "result": "no evaluation performed by lighthouse because no SLO file configured for project sockshop",
        "score": 0,
        "sloFileContent": "",
        "timeEnd": "2021-03-02T11:16:13Z",
        "timeStart": "2021-03-02T11:14:07Z"
      },
      "labels": {
        "DtCreds": "dynatrace"
      },
      "project": "sockshop",
      "result": "pass",
      "service": "carts",
      "stage": "staging"
    },
    "id": "90200421-420a-4087-b86b-12185d702237",
    "source": "lighthouse-service",
    "specversion": "1.0",
    "time": "2021-03-02T11:17:19.362Z",
    "type": "sh.keptn.event.evaluation.finished",
    "shkeptncontext": "a51a93b2-e50d-4daa-87fd-f0ad756649bd",
    "triggeredid": "11596a76-d224-4550-bf78-be0aa4a85d51",
    "plainEvent": {
      "data": {
        "evaluation": {
          "gitCommit": "",
          "indicatorResults": null,
          "result": "no evaluation performed by lighthouse because no SLO file configured for project sockshop",
          "score": 0,
          "sloFileContent": "",
          "timeEnd": "2021-03-02T11:16:13Z",
          "timeStart": "2021-03-02T11:14:07Z"
        },
        "labels": {
          "DtCreds": "dynatrace"
        },
        "project": "sockshop",
        "result": "pass",
        "service": "carts",
        "stage": "staging"
      },
      "id": "90200421-420a-4087-b86b-12185d702237",
      "source": "lighthouse-service",
      "specversion": "1.0",
      "time": "2021-03-02T11:17:19.362Z",
      "type": "sh.keptn.event.evaluation.finished",
      "shkeptncontext": "a51a93b2-e50d-4daa-87fd-f0ad756649bd",
      "triggeredid": "11596a76-d224-4550-bf78-be0aa4a85d51"
    },
    "heatmapLabel": "2021-03-02 12:17"
  },
  {
    "traces": [],
    "data": {
      "evaluation": {
        "comparedEvents": [
          "1163b9bc-4c4e-4299-a3f7-23a52aaf6312"
        ],
        "gitCommit": "",
        "indicatorResults": [
          {
            "DisplayName": "This is my very very very very very very very very very very very very long displayName",
            "keySli": false,
            "passTargets": [
              {
                "criteria": "<=+10%",
                "targetValue": 416.86674986015555,
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
              "value": 381.0334802140529
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
        "sloFileContent": "LS0tDQpzcGVjX3ZlcnNpb246ICIxLjAiDQpjb21wYXJpc29uOg0KICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciDQogIGNvbXBhcmVfd2l0aDogInNpbmdsZV9yZXN1bHQiDQogIGluY2x1ZGVfcmVzdWx0X3dpdGhfc2NvcmU6ICJwYXNzIg0KICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxDQpmaWx0ZXI6DQpvYmplY3RpdmVzOg0KICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1Ig0KICAgIGRpc3BsYXlOYW1lOiAiVGhpcyBpcyBteSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSBsb25nIGRpc3BsYXlOYW1lIg0KICAgIGtleV9zbGk6IGZhbHNlDQogICAgcGFzczogICAgICAgICAgICAgIyBwYXNzIGlmIChyZWxhdGl2ZSBjaGFuZ2UgPD0gMTAlIEFORCBhYnNvbHV0ZSB2YWx1ZSBpcyA8IDYwMG1zKQ0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykNCiAgICAgICAgICAtICI8NjAwIiAgICAjIGFic29sdXRlIHZhbHVlcyBvbmx5IHJlcXVpcmUgYSBsb2dpY2FsIG9wZXJhdG9yDQogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nDQogICAgICAtIGNyaXRlcmlhOg0KICAgICAgICAgIC0gIjw9ODAwIg0KICAgIHdlaWdodDogMQ0KdG90YWxfc2NvcmU6DQogIHBhc3M6ICI5MCUiDQogIHdhcm5pbmc6ICI3NSUi",
        "timeEnd": "2021-03-02T10:33:34Z",
        "timeStart": "2021-03-02T10:31:39Z"
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
    "id": "30745d37-c942-49dc-907d-5215898fa5a7",
    "source": "lighthouse-service",
    "specversion": "1.0",
    "time": "2021-03-02T10:35:40.859Z",
    "type": "sh.keptn.event.evaluation.finished",
    "shkeptncontext": "7c6de281-5429-4f3a-b476-dea064a1c14c",
    "triggeredid": "c9ecfd7c-c85c-43db-bfd5-cad01015bb7c",
    "plainEvent": {
      "data": {
        "evaluation": {
          "comparedEvents": [
            "1163b9bc-4c4e-4299-a3f7-23a52aaf6312"
          ],
          "gitCommit": "",
          "indicatorResults": [
            {
              "DisplayName": "This is my very very very very very very very very very very very very long displayName",
              "keySli": false,
              "passTargets": [
                {
                  "criteria": "<=+10%",
                  "targetValue": 416.86674986015555,
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
                "value": 381.0334802140529
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
          "sloFileContent": "LS0tDQpzcGVjX3ZlcnNpb246ICIxLjAiDQpjb21wYXJpc29uOg0KICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciDQogIGNvbXBhcmVfd2l0aDogInNpbmdsZV9yZXN1bHQiDQogIGluY2x1ZGVfcmVzdWx0X3dpdGhfc2NvcmU6ICJwYXNzIg0KICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxDQpmaWx0ZXI6DQpvYmplY3RpdmVzOg0KICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1Ig0KICAgIGRpc3BsYXlOYW1lOiAiVGhpcyBpcyBteSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSB2ZXJ5IHZlcnkgdmVyeSBsb25nIGRpc3BsYXlOYW1lIg0KICAgIGtleV9zbGk6IGZhbHNlDQogICAgcGFzczogICAgICAgICAgICAgIyBwYXNzIGlmIChyZWxhdGl2ZSBjaGFuZ2UgPD0gMTAlIEFORCBhYnNvbHV0ZSB2YWx1ZSBpcyA8IDYwMG1zKQ0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykNCiAgICAgICAgICAtICI8NjAwIiAgICAjIGFic29sdXRlIHZhbHVlcyBvbmx5IHJlcXVpcmUgYSBsb2dpY2FsIG9wZXJhdG9yDQogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nDQogICAgICAtIGNyaXRlcmlhOg0KICAgICAgICAgIC0gIjw9ODAwIg0KICAgIHdlaWdodDogMQ0KdG90YWxfc2NvcmU6DQogIHBhc3M6ICI5MCUiDQogIHdhcm5pbmc6ICI3NSUi",
          "timeEnd": "2021-03-02T10:33:34Z",
          "timeStart": "2021-03-02T10:31:39Z"
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
      "id": "30745d37-c942-49dc-907d-5215898fa5a7",
      "source": "lighthouse-service",
      "specversion": "1.0",
      "time": "2021-03-02T10:35:40.859Z",
      "type": "sh.keptn.event.evaluation.finished",
      "shkeptncontext": "7c6de281-5429-4f3a-b476-dea064a1c14c",
      "triggeredid": "c9ecfd7c-c85c-43db-bfd5-cad01015bb7c"
    },
    "heatmapLabel": "2021-03-02 11:35"
  },
  {
    "traces": [],
    "data": {
      "evaluation": {
        "comparedEvents": [
          "3d920466-17e7-4680-9a00-85c9763d9ead"
        ],
        "gitCommit": "",
        "indicatorResults": [
          {
            "keySli": false,
            "passTargets": [
              {
                "criteria": "<=+10%",
                "targetValue": 415.1789849888457,
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
              "value": 378.96977260014137
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
        "sloFileContent": "LS0tDQpzcGVjX3ZlcnNpb246ICIxLjAiDQpjb21wYXJpc29uOg0KICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciDQogIGNvbXBhcmVfd2l0aDogInNpbmdsZV9yZXN1bHQiDQogIGluY2x1ZGVfcmVzdWx0X3dpdGhfc2NvcmU6ICJwYXNzIg0KICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxDQpmaWx0ZXI6DQpvYmplY3RpdmVzOg0KICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1Ig0KICAgIGtleV9zbGk6IGZhbHNlDQogICAgcGFzczogICAgICAgICAgICAgIyBwYXNzIGlmIChyZWxhdGl2ZSBjaGFuZ2UgPD0gMTAlIEFORCBhYnNvbHV0ZSB2YWx1ZSBpcyA8IDYwMG1zKQ0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykNCiAgICAgICAgICAtICI8NjAwIiAgICAjIGFic29sdXRlIHZhbHVlcyBvbmx5IHJlcXVpcmUgYSBsb2dpY2FsIG9wZXJhdG9yDQogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nDQogICAgICAtIGNyaXRlcmlhOg0KICAgICAgICAgIC0gIjw9ODAwIg0KICAgIHdlaWdodDogMQ0KdG90YWxfc2NvcmU6DQogIHBhc3M6ICI5MCUiDQogIHdhcm5pbmc6ICI3NSUi",
        "timeEnd": "2021-03-01T14:43:51Z",
        "timeStart": "2021-03-01T14:41:45Z"
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
    "id": "1163b9bc-4c4e-4299-a3f7-23a52aaf6312",
    "source": "lighthouse-service",
    "specversion": "1.0",
    "time": "2021-03-01T14:44:57.759Z",
    "type": "sh.keptn.event.evaluation.finished",
    "shkeptncontext": "469e19ad-1de6-4c9e-bab6-d2bbacbba2fe",
    "triggeredid": "19cc9fab-fbad-4b8f-862b-a5853fb34e3b",
    "plainEvent": {
      "data": {
        "evaluation": {
          "comparedEvents": [
            "3d920466-17e7-4680-9a00-85c9763d9ead"
          ],
          "gitCommit": "",
          "indicatorResults": [
            {
              "keySli": false,
              "passTargets": [
                {
                  "criteria": "<=+10%",
                  "targetValue": 415.1789849888457,
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
                "value": 378.96977260014137
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
          "sloFileContent": "LS0tDQpzcGVjX3ZlcnNpb246ICIxLjAiDQpjb21wYXJpc29uOg0KICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciDQogIGNvbXBhcmVfd2l0aDogInNpbmdsZV9yZXN1bHQiDQogIGluY2x1ZGVfcmVzdWx0X3dpdGhfc2NvcmU6ICJwYXNzIg0KICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxDQpmaWx0ZXI6DQpvYmplY3RpdmVzOg0KICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1Ig0KICAgIGtleV9zbGk6IGZhbHNlDQogICAgcGFzczogICAgICAgICAgICAgIyBwYXNzIGlmIChyZWxhdGl2ZSBjaGFuZ2UgPD0gMTAlIEFORCBhYnNvbHV0ZSB2YWx1ZSBpcyA8IDYwMG1zKQ0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykNCiAgICAgICAgICAtICI8NjAwIiAgICAjIGFic29sdXRlIHZhbHVlcyBvbmx5IHJlcXVpcmUgYSBsb2dpY2FsIG9wZXJhdG9yDQogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nDQogICAgICAtIGNyaXRlcmlhOg0KICAgICAgICAgIC0gIjw9ODAwIg0KICAgIHdlaWdodDogMQ0KdG90YWxfc2NvcmU6DQogIHBhc3M6ICI5MCUiDQogIHdhcm5pbmc6ICI3NSUi",
          "timeEnd": "2021-03-01T14:43:51Z",
          "timeStart": "2021-03-01T14:41:45Z"
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
      "id": "1163b9bc-4c4e-4299-a3f7-23a52aaf6312",
      "source": "lighthouse-service",
      "specversion": "1.0",
      "time": "2021-03-01T14:44:57.759Z",
      "type": "sh.keptn.event.evaluation.finished",
      "shkeptncontext": "469e19ad-1de6-4c9e-bab6-d2bbacbba2fe",
      "triggeredid": "19cc9fab-fbad-4b8f-862b-a5853fb34e3b"
    },
    "heatmapLabel": "2021-03-01 15:44"
  },
  {
    "traces": [],
    "data": {
      "evaluation": {
        "comparedEvents": [
          "6b63eb7b-601f-48ca-b5a9-2adc65c8468a"
        ],
        "gitCommit": "",
        "indicatorResults": [
          {
            "keySli": false,
            "passTargets": [
              {
                "criteria": "<=+10%",
                "targetValue": 0,
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
              "value": 377.4354408989506
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
        "sloFileContent": "LS0tDQpzcGVjX3ZlcnNpb246ICIxLjAiDQpjb21wYXJpc29uOg0KICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciDQogIGNvbXBhcmVfd2l0aDogInNpbmdsZV9yZXN1bHQiDQogIGluY2x1ZGVfcmVzdWx0X3dpdGhfc2NvcmU6ICJwYXNzIg0KICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxDQpmaWx0ZXI6DQpvYmplY3RpdmVzOg0KICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1Ig0KICAgIGtleV9zbGk6IGZhbHNlDQogICAgcGFzczogICAgICAgICAgICAgIyBwYXNzIGlmIChyZWxhdGl2ZSBjaGFuZ2UgPD0gMTAlIEFORCBhYnNvbHV0ZSB2YWx1ZSBpcyA8IDYwMG1zKQ0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykNCiAgICAgICAgICAtICI8NjAwIiAgICAjIGFic29sdXRlIHZhbHVlcyBvbmx5IHJlcXVpcmUgYSBsb2dpY2FsIG9wZXJhdG9yDQogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nDQogICAgICAtIGNyaXRlcmlhOg0KICAgICAgICAgIC0gIjw9ODAwIg0KICAgIHdlaWdodDogMQ0KdG90YWxfc2NvcmU6DQogIHBhc3M6ICI5MCUiDQogIHdhcm5pbmc6ICI3NSUi",
        "timeEnd": "2021-03-01T13:46:49Z",
        "timeStart": "2021-03-01T13:44:44Z"
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
    "id": "3d920466-17e7-4680-9a00-85c9763d9ead",
    "source": "lighthouse-service",
    "specversion": "1.0",
    "time": "2021-03-01T13:47:55.419Z",
    "type": "sh.keptn.event.evaluation.finished",
    "shkeptncontext": "df4fc786-7704-451c-b230-2d6b26ebfd82",
    "triggeredid": "9eca21b5-97f9-4cdd-a674-918b532c14d8",
    "plainEvent": {
      "data": {
        "evaluation": {
          "comparedEvents": [
            "6b63eb7b-601f-48ca-b5a9-2adc65c8468a"
          ],
          "gitCommit": "",
          "indicatorResults": [
            {
              "keySli": false,
              "passTargets": [
                {
                  "criteria": "<=+10%",
                  "targetValue": 0,
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
                "value": 377.4354408989506
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
          "sloFileContent": "LS0tDQpzcGVjX3ZlcnNpb246ICIxLjAiDQpjb21wYXJpc29uOg0KICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciDQogIGNvbXBhcmVfd2l0aDogInNpbmdsZV9yZXN1bHQiDQogIGluY2x1ZGVfcmVzdWx0X3dpdGhfc2NvcmU6ICJwYXNzIg0KICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxDQpmaWx0ZXI6DQpvYmplY3RpdmVzOg0KICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1Ig0KICAgIGtleV9zbGk6IGZhbHNlDQogICAgcGFzczogICAgICAgICAgICAgIyBwYXNzIGlmIChyZWxhdGl2ZSBjaGFuZ2UgPD0gMTAlIEFORCBhYnNvbHV0ZSB2YWx1ZSBpcyA8IDYwMG1zKQ0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykNCiAgICAgICAgICAtICI8NjAwIiAgICAjIGFic29sdXRlIHZhbHVlcyBvbmx5IHJlcXVpcmUgYSBsb2dpY2FsIG9wZXJhdG9yDQogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nDQogICAgICAtIGNyaXRlcmlhOg0KICAgICAgICAgIC0gIjw9ODAwIg0KICAgIHdlaWdodDogMQ0KdG90YWxfc2NvcmU6DQogIHBhc3M6ICI5MCUiDQogIHdhcm5pbmc6ICI3NSUi",
          "timeEnd": "2021-03-01T13:46:49Z",
          "timeStart": "2021-03-01T13:44:44Z"
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
      "id": "3d920466-17e7-4680-9a00-85c9763d9ead",
      "source": "lighthouse-service",
      "specversion": "1.0",
      "time": "2021-03-01T13:47:55.419Z",
      "type": "sh.keptn.event.evaluation.finished",
      "shkeptncontext": "df4fc786-7704-451c-b230-2d6b26ebfd82",
      "triggeredid": "9eca21b5-97f9-4cdd-a674-918b532c14d8"
    },
    "heatmapLabel": "2021-03-01 14:47"
  },
  {
    "traces": [],
    "data": {
      "evaluation": {
        "gitCommit": "",
        "indicatorResults": null,
        "result": "no evaluation performed by lighthouse because no SLI-provider configured for project sockshop",
        "score": 0,
        "sloFileContent": "",
        "timeEnd": "2021-03-01T13:19:11Z",
        "timeStart": "2021-03-01T13:17:19Z"
      },
      "message": "no evaluation performed by lighthouse because no SLI-provider configured for project sockshop",
      "project": "sockshop",
      "result": "pass",
      "service": "carts",
      "stage": "staging",
      "status": "succeeded"
    },
    "id": "6b63eb7b-601f-48ca-b5a9-2adc65c8468a",
    "source": "lighthouse-service",
    "specversion": "1.0",
    "time": "2021-03-01T13:19:12.922Z",
    "type": "sh.keptn.event.evaluation.finished",
    "shkeptncontext": "b0516983-f716-427b-a73f-73bf19fcefcd",
    "triggeredid": "35b4ada8-2ffd-4444-a111-26f980362030",
    "plainEvent": {
      "data": {
        "evaluation": {
          "gitCommit": "",
          "indicatorResults": null,
          "result": "no evaluation performed by lighthouse because no SLI-provider configured for project sockshop",
          "score": 0,
          "sloFileContent": "",
          "timeEnd": "2021-03-01T13:19:11Z",
          "timeStart": "2021-03-01T13:17:19Z"
        },
        "message": "no evaluation performed by lighthouse because no SLI-provider configured for project sockshop",
        "project": "sockshop",
        "result": "pass",
        "service": "carts",
        "stage": "staging",
        "status": "succeeded"
      },
      "id": "6b63eb7b-601f-48ca-b5a9-2adc65c8468a",
      "source": "lighthouse-service",
      "specversion": "1.0",
      "time": "2021-03-01T13:19:12.922Z",
      "type": "sh.keptn.event.evaluation.finished",
      "shkeptncontext": "b0516983-f716-427b-a73f-73bf19fcefcd",
      "triggeredid": "35b4ada8-2ffd-4444-a111-26f980362030"
    },
    "heatmapLabel": "2021-03-01 14:19"
  }
];

const Evaluations = traceData.map(evaluation => Trace.fromJSON(evaluation));

export {Evaluations};
