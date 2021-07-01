import {Trace} from "../../_models/trace";

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
    "evaluationHistory": [
      {
      "traces": [],
      "data": {
        "evaluation": {
          "gitCommit": "",
          "indicatorResults": [
            {
              "displayName": "Response time P95 A",
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
                "message": "Error when fetching SLI config for response_time_p95a Unsupported SLI metric response_time_p95a.",
                "metric": "response_time_p95a",
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
            },
            {
              "displayName": "Response time P95 B",
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
                "message": "Error when fetching SLI config for response_time_p95b Unsupported SLI metric response_time_p95b.",
                "metric": "response_time_p95b",
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
          "sloFileContent": "LS0tDQpzcGVjX3ZlcnNpb246ICIxLjAiDQpjb21wYXJpc29uOg0KICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciDQogIGNvbXBhcmVfd2l0aDogInNpbmdsZV9yZXN1bHQiDQogIGluY2x1ZGVfcmVzdWx0X3dpdGhfc2NvcmU6ICJwYXNzIg0KICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxDQpmaWx0ZXI6DQpvYmplY3RpdmVzOg0KICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1YSINCiAgICBkaXNwbGF5TmFtZTogIlJlc3BvbnNlIHRpbWUgUDk1Ig0KICAgIGtleV9zbGk6IGZhbHNlDQogICAgcGFzczogICAgICAgICAgICAgIyBwYXNzIGlmIChyZWxhdGl2ZSBjaGFuZ2UgPD0gMTAlIEFORCBhYnNvbHV0ZSB2YWx1ZSBpcyA8IDYwMG1zKQ0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykNCiAgICAgICAgICAtICI8NjAwIiAgICAjIGFic29sdXRlIHZhbHVlcyBvbmx5IHJlcXVpcmUgYSBsb2dpY2FsIG9wZXJhdG9yDQogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nDQogICAgICAtIGNyaXRlcmlhOg0KICAgICAgICAgIC0gIjw9ODAwIg0KICAgIHdlaWdodDogMQ0KICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1YiINCiAgICBkaXNwbGF5TmFtZTogIlJlc3BvbnNlIHRpbWUgUDk1Ig0KICAgIGtleV9zbGk6IGZhbHNlDQogICAgcGFzczogICAgICAgICAgICAgIyBwYXNzIGlmIChyZWxhdGl2ZSBjaGFuZ2UgPD0gMTAlIEFORCBhYnNvbHV0ZSB2YWx1ZSBpcyA8IDYwMG1zKQ0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykNCiAgICAgICAgICAtICI8NjAwIiAgICAjIGFic29sdXRlIHZhbHVlcyBvbmx5IHJlcXVpcmUgYSBsb2dpY2FsIG9wZXJhdG9yDQogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nDQogICAgICAtIGNyaXRlcmlhOg0KICAgICAgICAgIC0gIjw9ODAwIg0KICAgIHdlaWdodDogMQ0KICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1YyINCiAgICBkaXNwbGF5TmFtZTogIlJlc3BvbnNlIHRpbWUgUDk1Ig0KICAgIGtleV9zbGk6IGZhbHNlDQogICAgcGFzczogICAgICAgICAgICAgIyBwYXNzIGlmIChyZWxhdGl2ZSBjaGFuZ2UgPD0gMTAlIEFORCBhYnNvbHV0ZSB2YWx1ZSBpcyA8IDYwMG1zKQ0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykNCiAgICAgICAgICAtICI8NjAwIiAgICAjIGFic29sdXRlIHZhbHVlcyBvbmx5IHJlcXVpcmUgYSBsb2dpY2FsIG9wZXJhdG9yDQogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nDQogICAgICAtIGNyaXRlcmlhOg0KICAgICAgICAgIC0gIjw9ODAwIg0KICAgIHdlaWdodDogMQ0KICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1ZCINCiAgICBkaXNwbGF5TmFtZTogIlJlc3BvbnNlIHRpbWUgUDk1Ig0KICAgIGtleV9zbGk6IGZhbHNlDQogICAgcGFzczogICAgICAgICAgICAgIyBwYXNzIGlmIChyZWxhdGl2ZSBjaGFuZ2UgPD0gMTAlIEFORCBhYnNvbHV0ZSB2YWx1ZSBpcyA8IDYwMG1zKQ0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykNCiAgICAgICAgICAtICI8NjAwIiAgICAjIGFic29sdXRlIHZhbHVlcyBvbmx5IHJlcXVpcmUgYSBsb2dpY2FsIG9wZXJhdG9yDQogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nDQogICAgICAtIGNyaXRlcmlhOg0KICAgICAgICAgIC0gIjw9ODAwIg0KICAgIHdlaWdodDogMQ0KICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1ZSINCiAgICBkaXNwbGF5TmFtZTogIlJlc3BvbnNlIHRpbWUgUDk1Ig0KICAgIGtleV9zbGk6IGZhbHNlDQogICAgcGFzczogICAgICAgICAgICAgIyBwYXNzIGlmIChyZWxhdGl2ZSBjaGFuZ2UgPD0gMTAlIEFORCBhYnNvbHV0ZSB2YWx1ZSBpcyA8IDYwMG1zKQ0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykNCiAgICAgICAgICAtICI8NjAwIiAgICAjIGFic29sdXRlIHZhbHVlcyBvbmx5IHJlcXVpcmUgYSBsb2dpY2FsIG9wZXJhdG9yDQogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nDQogICAgICAtIGNyaXRlcmlhOg0KICAgICAgICAgIC0gIjw9ODAwIg0KICAgIHdlaWdodDogMQ0KICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1ZiINCiAgICBkaXNwbGF5TmFtZTogIlJlc3BvbnNlIHRpbWUgUDk1Ig0KICAgIGtleV9zbGk6IGZhbHNlDQogICAgcGFzczogICAgICAgICAgICAgIyBwYXNzIGlmIChyZWxhdGl2ZSBjaGFuZ2UgPD0gMTAlIEFORCBhYnNvbHV0ZSB2YWx1ZSBpcyA8IDYwMG1zKQ0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykNCiAgICAgICAgICAtICI8NjAwIiAgICAjIGFic29sdXRlIHZhbHVlcyBvbmx5IHJlcXVpcmUgYSBsb2dpY2FsIG9wZXJhdG9yDQogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nDQogICAgICAtIGNyaXRlcmlhOg0KICAgICAgICAgIC0gIjw9ODAwIg0KICAgIHdlaWdodDogMQ0KICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1ZyINCiAgICBkaXNwbGF5TmFtZTogIlJlc3BvbnNlIHRpbWUgUDk1Ig0KICAgIGtleV9zbGk6IGZhbHNlDQogICAgcGFzczogICAgICAgICAgICAgIyBwYXNzIGlmIChyZWxhdGl2ZSBjaGFuZ2UgPD0gMTAlIEFORCBhYnNvbHV0ZSB2YWx1ZSBpcyA8IDYwMG1zKQ0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykNCiAgICAgICAgICAtICI8NjAwIiAgICAjIGFic29sdXRlIHZhbHVlcyBvbmx5IHJlcXVpcmUgYSBsb2dpY2FsIG9wZXJhdG9yDQogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nDQogICAgICAtIGNyaXRlcmlhOg0KICAgICAgICAgIC0gIjw9ODAwIg0KICAgIHdlaWdodDogMQ0KICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1aCINCiAgICBkaXNwbGF5TmFtZTogIlJlc3BvbnNlIHRpbWUgUDk1Ig0KICAgIGtleV9zbGk6IGZhbHNlDQogICAgcGFzczogICAgICAgICAgICAgIyBwYXNzIGlmIChyZWxhdGl2ZSBjaGFuZ2UgPD0gMTAlIEFORCBhYnNvbHV0ZSB2YWx1ZSBpcyA8IDYwMG1zKQ0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykNCiAgICAgICAgICAtICI8NjAwIiAgICAjIGFic29sdXRlIHZhbHVlcyBvbmx5IHJlcXVpcmUgYSBsb2dpY2FsIG9wZXJhdG9yDQogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nDQogICAgICAtIGNyaXRlcmlhOg0KICAgICAgICAgIC0gIjw9ODAwIg0KICAgIHdlaWdodDogMQ0KICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1aSINCiAgICBkaXNwbGF5TmFtZTogIlJlc3BvbnNlIHRpbWUgUDk1Ig0KICAgIGtleV9zbGk6IGZhbHNlDQogICAgcGFzczogICAgICAgICAgICAgIyBwYXNzIGlmIChyZWxhdGl2ZSBjaGFuZ2UgPD0gMTAlIEFORCBhYnNvbHV0ZSB2YWx1ZSBpcyA8IDYwMG1zKQ0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykNCiAgICAgICAgICAtICI8NjAwIiAgICAjIGFic29sdXRlIHZhbHVlcyBvbmx5IHJlcXVpcmUgYSBsb2dpY2FsIG9wZXJhdG9yDQogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nDQogICAgICAtIGNyaXRlcmlhOg0KICAgICAgICAgIC0gIjw9ODAwIg0KICAgIHdlaWdodDogMQ0KICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1aiINCiAgICBkaXNwbGF5TmFtZTogIlJlc3BvbnNlIHRpbWUgUDk1Ig0KICAgIGtleV9zbGk6IGZhbHNlDQogICAgcGFzczogICAgICAgICAgICAgIyBwYXNzIGlmIChyZWxhdGl2ZSBjaGFuZ2UgPD0gMTAlIEFORCBhYnNvbHV0ZSB2YWx1ZSBpcyA8IDYwMG1zKQ0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykNCiAgICAgICAgICAtICI8NjAwIiAgICAjIGFic29sdXRlIHZhbHVlcyBvbmx5IHJlcXVpcmUgYSBsb2dpY2FsIG9wZXJhdG9yDQogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nDQogICAgICAtIGNyaXRlcmlhOg0KICAgICAgICAgIC0gIjw9ODAwIg0KICAgIHdlaWdodDogMQ0KICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1ayINCiAgICBkaXNwbGF5TmFtZTogIlJlc3BvbnNlIHRpbWUgUDk1Ig0KICAgIGtleV9zbGk6IGZhbHNlDQogICAgcGFzczogICAgICAgICAgICAgIyBwYXNzIGlmIChyZWxhdGl2ZSBjaGFuZ2UgPD0gMTAlIEFORCBhYnNvbHV0ZSB2YWx1ZSBpcyA8IDYwMG1zKQ0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykNCiAgICAgICAgICAtICI8NjAwIiAgICAjIGFic29sdXRlIHZhbHVlcyBvbmx5IHJlcXVpcmUgYSBsb2dpY2FsIG9wZXJhdG9yDQogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nDQogICAgICAtIGNyaXRlcmlhOg0KICAgICAgICAgIC0gIjw9ODAwIg0KICAgIHdlaWdodDogMQ0KICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1bCINCiAgICBkaXNwbGF5TmFtZTogIlJlc3BvbnNlIHRpbWUgUDk1Ig0KICAgIGtleV9zbGk6IGZhbHNlDQogICAgcGFzczogICAgICAgICAgICAgIyBwYXNzIGlmIChyZWxhdGl2ZSBjaGFuZ2UgPD0gMTAlIEFORCBhYnNvbHV0ZSB2YWx1ZSBpcyA8IDYwMG1zKQ0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykNCiAgICAgICAgICAtICI8NjAwIiAgICAjIGFic29sdXRlIHZhbHVlcyBvbmx5IHJlcXVpcmUgYSBsb2dpY2FsIG9wZXJhdG9yDQogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nDQogICAgICAtIGNyaXRlcmlhOg0KICAgICAgICAgIC0gIjw9ODAwIg0KICAgIHdlaWdodDogMQ0KICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1bSINCiAgICBkaXNwbGF5TmFtZTogIlJlc3BvbnNlIHRpbWUgUDk1Ig0KICAgIGtleV9zbGk6IGZhbHNlDQogICAgcGFzczogICAgICAgICAgICAgIyBwYXNzIGlmIChyZWxhdGl2ZSBjaGFuZ2UgPD0gMTAlIEFORCBhYnNvbHV0ZSB2YWx1ZSBpcyA8IDYwMG1zKQ0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykNCiAgICAgICAgICAtICI8NjAwIiAgICAjIGFic29sdXRlIHZhbHVlcyBvbmx5IHJlcXVpcmUgYSBsb2dpY2FsIG9wZXJhdG9yDQogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nDQogICAgICAtIGNyaXRlcmlhOg0KICAgICAgICAgIC0gIjw9ODAwIg0KICAgIHdlaWdodDogMQ0KICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1biINCiAgICBkaXNwbGF5TmFtZTogIlJlc3BvbnNlIHRpbWUgUDk1Ig0KICAgIGtleV9zbGk6IGZhbHNlDQogICAgcGFzczogICAgICAgICAgICAgIyBwYXNzIGlmIChyZWxhdGl2ZSBjaGFuZ2UgPD0gMTAlIEFORCBhYnNvbHV0ZSB2YWx1ZSBpcyA8IDYwMG1zKQ0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykNCiAgICAgICAgICAtICI8NjAwIiAgICAjIGFic29sdXRlIHZhbHVlcyBvbmx5IHJlcXVpcmUgYSBsb2dpY2FsIG9wZXJhdG9yDQogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nDQogICAgICAtIGNyaXRlcmlhOg0KICAgICAgICAgIC0gIjw9ODAwIg0KICAgIHdlaWdodDogMQ0KICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1byINCiAgICBkaXNwbGF5TmFtZTogIlJlc3BvbnNlIHRpbWUgUDk1Ig0KICAgIGtleV9zbGk6IGZhbHNlDQogICAgcGFzczogICAgICAgICAgICAgIyBwYXNzIGlmIChyZWxhdGl2ZSBjaGFuZ2UgPD0gMTAlIEFORCBhYnNvbHV0ZSB2YWx1ZSBpcyA8IDYwMG1zKQ0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykNCiAgICAgICAgICAtICI8NjAwIiAgICAjIGFic29sdXRlIHZhbHVlcyBvbmx5IHJlcXVpcmUgYSBsb2dpY2FsIG9wZXJhdG9yDQogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nDQogICAgICAtIGNyaXRlcmlhOg0KICAgICAgICAgIC0gIjw9ODAwIg0KICAgIHdlaWdodDogMQ0KICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1cCINCiAgICBkaXNwbGF5TmFtZTogIlJlc3BvbnNlIHRpbWUgUDk1Ig0KICAgIGtleV9zbGk6IGZhbHNlDQogICAgcGFzczogICAgICAgICAgICAgIyBwYXNzIGlmIChyZWxhdGl2ZSBjaGFuZ2UgPD0gMTAlIEFORCBhYnNvbHV0ZSB2YWx1ZSBpcyA8IDYwMG1zKQ0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykNCiAgICAgICAgICAtICI8NjAwIiAgICAjIGFic29sdXRlIHZhbHVlcyBvbmx5IHJlcXVpcmUgYSBsb2dpY2FsIG9wZXJhdG9yDQogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nDQogICAgICAtIGNyaXRlcmlhOg0KICAgICAgICAgIC0gIjw9ODAwIg0KICAgIHdlaWdodDogMQ0KDQp0b3RhbF9zY29yZToNCiAgcGFzczogIjkwJSINCiAgd2FybmluZzogIjc1JSINCg==",
          "timeEnd": "2021-06-22T14:50:47.000Z",
          "timeStart": "2021-06-22T14:45:47.000Z"
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
      "id": "625413dd-a5b7-4cc8-8ce4-bcd5ead2a484",
      "source": "lighthouse-service",
      "specversion": "1.0",
      "time": "2021-06-22T14:50:48.080Z",
      "type": "sh.keptn.event.evaluation.finished",
      "shkeptncontext": "7036154d-d366-445e-b9a4-2b62649b0bfd",
      "shkeptnspecversion": "0.2.1",
      "triggeredid": "a1392234-6671-4681-8e54-65614557a667",
      "plainEvent": {
        "data": {
          "evaluation": {
            "gitCommit": "",
            "indicatorResults": [
              {
                "displayName": "Response time P95",
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
                  "message": "Error when fetching SLI config for response_time_p95a Unsupported SLI metric response_time_p95a.",
                  "metric": "response_time_p95a",
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
              },
              {
                "displayName": "Response time P95",
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
                  "message": "Error when fetching SLI config for response_time_p95b Unsupported SLI metric response_time_p95b.",
                  "metric": "response_time_p95b",
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
              },
              {
                "displayName": "Response time P95",
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
                  "message": "Error when fetching SLI config for response_time_p95c Unsupported SLI metric response_time_p95c.",
                  "metric": "response_time_p95c",
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
              },
              {
                "displayName": "Response time P95",
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
                  "message": "Error when fetching SLI config for response_time_p95d Unsupported SLI metric response_time_p95d.",
                  "metric": "response_time_p95d",
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
              },
              {
                "displayName": "Response time P95",
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
                  "message": "Error when fetching SLI config for response_time_p95e Unsupported SLI metric response_time_p95e.",
                  "metric": "response_time_p95e",
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
              },
              {
                "displayName": "Response time P95",
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
                  "message": "Error when fetching SLI config for response_time_p95f Unsupported SLI metric response_time_p95f.",
                  "metric": "response_time_p95f",
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
              },
              {
                "displayName": "Response time P95",
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
                  "message": "Error when fetching SLI config for response_time_p95g Unsupported SLI metric response_time_p95g.",
                  "metric": "response_time_p95g",
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
              },
              {
                "displayName": "Response time P95",
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
                  "message": "Error when fetching SLI config for response_time_p95h Unsupported SLI metric response_time_p95h.",
                  "metric": "response_time_p95h",
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
              },
              {
                "displayName": "Response time P95",
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
                  "message": "Error when fetching SLI config for response_time_p95i Unsupported SLI metric response_time_p95i.",
                  "metric": "response_time_p95i",
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
              },
              {
                "displayName": "Response time P95",
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
                  "message": "Error when fetching SLI config for response_time_p95j Unsupported SLI metric response_time_p95j.",
                  "metric": "response_time_p95j",
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
              },
              {
                "displayName": "Response time P95",
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
                  "message": "Error when fetching SLI config for response_time_p95k Unsupported SLI metric response_time_p95k.",
                  "metric": "response_time_p95k",
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
              },
              {
                "displayName": "Response time P95",
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
                  "message": "Error when fetching SLI config for response_time_p95l Unsupported SLI metric response_time_p95l.",
                  "metric": "response_time_p95l",
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
              },
              {
                "displayName": "Response time P95",
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
                  "message": "Error when fetching SLI config for response_time_p95m Unsupported SLI metric response_time_p95m.",
                  "metric": "response_time_p95m",
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
              },
              {
                "displayName": "Response time P95",
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
                  "message": "Error when fetching SLI config for response_time_p95n Unsupported SLI metric response_time_p95n.",
                  "metric": "response_time_p95n",
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
              },
              {
                "displayName": "Response time P95",
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
                  "message": "Error when fetching SLI config for response_time_p95o Unsupported SLI metric response_time_p95o.",
                  "metric": "response_time_p95o",
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
              },
              {
                "displayName": "Response time P95",
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
                  "message": "Error when fetching SLI config for response_time_p95p Unsupported SLI metric response_time_p95p.",
                  "metric": "response_time_p95p",
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
            "sloFileContent": "LS0tDQpzcGVjX3ZlcnNpb246ICIxLjAiDQpjb21wYXJpc29uOg0KICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciDQogIGNvbXBhcmVfd2l0aDogInNpbmdsZV9yZXN1bHQiDQogIGluY2x1ZGVfcmVzdWx0X3dpdGhfc2NvcmU6ICJwYXNzIg0KICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxDQpmaWx0ZXI6DQpvYmplY3RpdmVzOg0KICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1YSINCiAgICBkaXNwbGF5TmFtZTogIlJlc3BvbnNlIHRpbWUgUDk1Ig0KICAgIGtleV9zbGk6IGZhbHNlDQogICAgcGFzczogICAgICAgICAgICAgIyBwYXNzIGlmIChyZWxhdGl2ZSBjaGFuZ2UgPD0gMTAlIEFORCBhYnNvbHV0ZSB2YWx1ZSBpcyA8IDYwMG1zKQ0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykNCiAgICAgICAgICAtICI8NjAwIiAgICAjIGFic29sdXRlIHZhbHVlcyBvbmx5IHJlcXVpcmUgYSBsb2dpY2FsIG9wZXJhdG9yDQogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nDQogICAgICAtIGNyaXRlcmlhOg0KICAgICAgICAgIC0gIjw9ODAwIg0KICAgIHdlaWdodDogMQ0KICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1YiINCiAgICBkaXNwbGF5TmFtZTogIlJlc3BvbnNlIHRpbWUgUDk1Ig0KICAgIGtleV9zbGk6IGZhbHNlDQogICAgcGFzczogICAgICAgICAgICAgIyBwYXNzIGlmIChyZWxhdGl2ZSBjaGFuZ2UgPD0gMTAlIEFORCBhYnNvbHV0ZSB2YWx1ZSBpcyA8IDYwMG1zKQ0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykNCiAgICAgICAgICAtICI8NjAwIiAgICAjIGFic29sdXRlIHZhbHVlcyBvbmx5IHJlcXVpcmUgYSBsb2dpY2FsIG9wZXJhdG9yDQogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nDQogICAgICAtIGNyaXRlcmlhOg0KICAgICAgICAgIC0gIjw9ODAwIg0KICAgIHdlaWdodDogMQ0KICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1YyINCiAgICBkaXNwbGF5TmFtZTogIlJlc3BvbnNlIHRpbWUgUDk1Ig0KICAgIGtleV9zbGk6IGZhbHNlDQogICAgcGFzczogICAgICAgICAgICAgIyBwYXNzIGlmIChyZWxhdGl2ZSBjaGFuZ2UgPD0gMTAlIEFORCBhYnNvbHV0ZSB2YWx1ZSBpcyA8IDYwMG1zKQ0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykNCiAgICAgICAgICAtICI8NjAwIiAgICAjIGFic29sdXRlIHZhbHVlcyBvbmx5IHJlcXVpcmUgYSBsb2dpY2FsIG9wZXJhdG9yDQogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nDQogICAgICAtIGNyaXRlcmlhOg0KICAgICAgICAgIC0gIjw9ODAwIg0KICAgIHdlaWdodDogMQ0KICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1ZCINCiAgICBkaXNwbGF5TmFtZTogIlJlc3BvbnNlIHRpbWUgUDk1Ig0KICAgIGtleV9zbGk6IGZhbHNlDQogICAgcGFzczogICAgICAgICAgICAgIyBwYXNzIGlmIChyZWxhdGl2ZSBjaGFuZ2UgPD0gMTAlIEFORCBhYnNvbHV0ZSB2YWx1ZSBpcyA8IDYwMG1zKQ0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykNCiAgICAgICAgICAtICI8NjAwIiAgICAjIGFic29sdXRlIHZhbHVlcyBvbmx5IHJlcXVpcmUgYSBsb2dpY2FsIG9wZXJhdG9yDQogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nDQogICAgICAtIGNyaXRlcmlhOg0KICAgICAgICAgIC0gIjw9ODAwIg0KICAgIHdlaWdodDogMQ0KICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1ZSINCiAgICBkaXNwbGF5TmFtZTogIlJlc3BvbnNlIHRpbWUgUDk1Ig0KICAgIGtleV9zbGk6IGZhbHNlDQogICAgcGFzczogICAgICAgICAgICAgIyBwYXNzIGlmIChyZWxhdGl2ZSBjaGFuZ2UgPD0gMTAlIEFORCBhYnNvbHV0ZSB2YWx1ZSBpcyA8IDYwMG1zKQ0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykNCiAgICAgICAgICAtICI8NjAwIiAgICAjIGFic29sdXRlIHZhbHVlcyBvbmx5IHJlcXVpcmUgYSBsb2dpY2FsIG9wZXJhdG9yDQogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nDQogICAgICAtIGNyaXRlcmlhOg0KICAgICAgICAgIC0gIjw9ODAwIg0KICAgIHdlaWdodDogMQ0KICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1ZiINCiAgICBkaXNwbGF5TmFtZTogIlJlc3BvbnNlIHRpbWUgUDk1Ig0KICAgIGtleV9zbGk6IGZhbHNlDQogICAgcGFzczogICAgICAgICAgICAgIyBwYXNzIGlmIChyZWxhdGl2ZSBjaGFuZ2UgPD0gMTAlIEFORCBhYnNvbHV0ZSB2YWx1ZSBpcyA8IDYwMG1zKQ0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykNCiAgICAgICAgICAtICI8NjAwIiAgICAjIGFic29sdXRlIHZhbHVlcyBvbmx5IHJlcXVpcmUgYSBsb2dpY2FsIG9wZXJhdG9yDQogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nDQogICAgICAtIGNyaXRlcmlhOg0KICAgICAgICAgIC0gIjw9ODAwIg0KICAgIHdlaWdodDogMQ0KICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1ZyINCiAgICBkaXNwbGF5TmFtZTogIlJlc3BvbnNlIHRpbWUgUDk1Ig0KICAgIGtleV9zbGk6IGZhbHNlDQogICAgcGFzczogICAgICAgICAgICAgIyBwYXNzIGlmIChyZWxhdGl2ZSBjaGFuZ2UgPD0gMTAlIEFORCBhYnNvbHV0ZSB2YWx1ZSBpcyA8IDYwMG1zKQ0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykNCiAgICAgICAgICAtICI8NjAwIiAgICAjIGFic29sdXRlIHZhbHVlcyBvbmx5IHJlcXVpcmUgYSBsb2dpY2FsIG9wZXJhdG9yDQogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nDQogICAgICAtIGNyaXRlcmlhOg0KICAgICAgICAgIC0gIjw9ODAwIg0KICAgIHdlaWdodDogMQ0KICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1aCINCiAgICBkaXNwbGF5TmFtZTogIlJlc3BvbnNlIHRpbWUgUDk1Ig0KICAgIGtleV9zbGk6IGZhbHNlDQogICAgcGFzczogICAgICAgICAgICAgIyBwYXNzIGlmIChyZWxhdGl2ZSBjaGFuZ2UgPD0gMTAlIEFORCBhYnNvbHV0ZSB2YWx1ZSBpcyA8IDYwMG1zKQ0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykNCiAgICAgICAgICAtICI8NjAwIiAgICAjIGFic29sdXRlIHZhbHVlcyBvbmx5IHJlcXVpcmUgYSBsb2dpY2FsIG9wZXJhdG9yDQogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nDQogICAgICAtIGNyaXRlcmlhOg0KICAgICAgICAgIC0gIjw9ODAwIg0KICAgIHdlaWdodDogMQ0KICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1aSINCiAgICBkaXNwbGF5TmFtZTogIlJlc3BvbnNlIHRpbWUgUDk1Ig0KICAgIGtleV9zbGk6IGZhbHNlDQogICAgcGFzczogICAgICAgICAgICAgIyBwYXNzIGlmIChyZWxhdGl2ZSBjaGFuZ2UgPD0gMTAlIEFORCBhYnNvbHV0ZSB2YWx1ZSBpcyA8IDYwMG1zKQ0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykNCiAgICAgICAgICAtICI8NjAwIiAgICAjIGFic29sdXRlIHZhbHVlcyBvbmx5IHJlcXVpcmUgYSBsb2dpY2FsIG9wZXJhdG9yDQogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nDQogICAgICAtIGNyaXRlcmlhOg0KICAgICAgICAgIC0gIjw9ODAwIg0KICAgIHdlaWdodDogMQ0KICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1aiINCiAgICBkaXNwbGF5TmFtZTogIlJlc3BvbnNlIHRpbWUgUDk1Ig0KICAgIGtleV9zbGk6IGZhbHNlDQogICAgcGFzczogICAgICAgICAgICAgIyBwYXNzIGlmIChyZWxhdGl2ZSBjaGFuZ2UgPD0gMTAlIEFORCBhYnNvbHV0ZSB2YWx1ZSBpcyA8IDYwMG1zKQ0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykNCiAgICAgICAgICAtICI8NjAwIiAgICAjIGFic29sdXRlIHZhbHVlcyBvbmx5IHJlcXVpcmUgYSBsb2dpY2FsIG9wZXJhdG9yDQogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nDQogICAgICAtIGNyaXRlcmlhOg0KICAgICAgICAgIC0gIjw9ODAwIg0KICAgIHdlaWdodDogMQ0KICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1ayINCiAgICBkaXNwbGF5TmFtZTogIlJlc3BvbnNlIHRpbWUgUDk1Ig0KICAgIGtleV9zbGk6IGZhbHNlDQogICAgcGFzczogICAgICAgICAgICAgIyBwYXNzIGlmIChyZWxhdGl2ZSBjaGFuZ2UgPD0gMTAlIEFORCBhYnNvbHV0ZSB2YWx1ZSBpcyA8IDYwMG1zKQ0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykNCiAgICAgICAgICAtICI8NjAwIiAgICAjIGFic29sdXRlIHZhbHVlcyBvbmx5IHJlcXVpcmUgYSBsb2dpY2FsIG9wZXJhdG9yDQogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nDQogICAgICAtIGNyaXRlcmlhOg0KICAgICAgICAgIC0gIjw9ODAwIg0KICAgIHdlaWdodDogMQ0KICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1bCINCiAgICBkaXNwbGF5TmFtZTogIlJlc3BvbnNlIHRpbWUgUDk1Ig0KICAgIGtleV9zbGk6IGZhbHNlDQogICAgcGFzczogICAgICAgICAgICAgIyBwYXNzIGlmIChyZWxhdGl2ZSBjaGFuZ2UgPD0gMTAlIEFORCBhYnNvbHV0ZSB2YWx1ZSBpcyA8IDYwMG1zKQ0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykNCiAgICAgICAgICAtICI8NjAwIiAgICAjIGFic29sdXRlIHZhbHVlcyBvbmx5IHJlcXVpcmUgYSBsb2dpY2FsIG9wZXJhdG9yDQogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nDQogICAgICAtIGNyaXRlcmlhOg0KICAgICAgICAgIC0gIjw9ODAwIg0KICAgIHdlaWdodDogMQ0KICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1bSINCiAgICBkaXNwbGF5TmFtZTogIlJlc3BvbnNlIHRpbWUgUDk1Ig0KICAgIGtleV9zbGk6IGZhbHNlDQogICAgcGFzczogICAgICAgICAgICAgIyBwYXNzIGlmIChyZWxhdGl2ZSBjaGFuZ2UgPD0gMTAlIEFORCBhYnNvbHV0ZSB2YWx1ZSBpcyA8IDYwMG1zKQ0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykNCiAgICAgICAgICAtICI8NjAwIiAgICAjIGFic29sdXRlIHZhbHVlcyBvbmx5IHJlcXVpcmUgYSBsb2dpY2FsIG9wZXJhdG9yDQogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nDQogICAgICAtIGNyaXRlcmlhOg0KICAgICAgICAgIC0gIjw9ODAwIg0KICAgIHdlaWdodDogMQ0KICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1biINCiAgICBkaXNwbGF5TmFtZTogIlJlc3BvbnNlIHRpbWUgUDk1Ig0KICAgIGtleV9zbGk6IGZhbHNlDQogICAgcGFzczogICAgICAgICAgICAgIyBwYXNzIGlmIChyZWxhdGl2ZSBjaGFuZ2UgPD0gMTAlIEFORCBhYnNvbHV0ZSB2YWx1ZSBpcyA8IDYwMG1zKQ0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykNCiAgICAgICAgICAtICI8NjAwIiAgICAjIGFic29sdXRlIHZhbHVlcyBvbmx5IHJlcXVpcmUgYSBsb2dpY2FsIG9wZXJhdG9yDQogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nDQogICAgICAtIGNyaXRlcmlhOg0KICAgICAgICAgIC0gIjw9ODAwIg0KICAgIHdlaWdodDogMQ0KICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1byINCiAgICBkaXNwbGF5TmFtZTogIlJlc3BvbnNlIHRpbWUgUDk1Ig0KICAgIGtleV9zbGk6IGZhbHNlDQogICAgcGFzczogICAgICAgICAgICAgIyBwYXNzIGlmIChyZWxhdGl2ZSBjaGFuZ2UgPD0gMTAlIEFORCBhYnNvbHV0ZSB2YWx1ZSBpcyA8IDYwMG1zKQ0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykNCiAgICAgICAgICAtICI8NjAwIiAgICAjIGFic29sdXRlIHZhbHVlcyBvbmx5IHJlcXVpcmUgYSBsb2dpY2FsIG9wZXJhdG9yDQogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nDQogICAgICAtIGNyaXRlcmlhOg0KICAgICAgICAgIC0gIjw9ODAwIg0KICAgIHdlaWdodDogMQ0KICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1cCINCiAgICBkaXNwbGF5TmFtZTogIlJlc3BvbnNlIHRpbWUgUDk1Ig0KICAgIGtleV9zbGk6IGZhbHNlDQogICAgcGFzczogICAgICAgICAgICAgIyBwYXNzIGlmIChyZWxhdGl2ZSBjaGFuZ2UgPD0gMTAlIEFORCBhYnNvbHV0ZSB2YWx1ZSBpcyA8IDYwMG1zKQ0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykNCiAgICAgICAgICAtICI8NjAwIiAgICAjIGFic29sdXRlIHZhbHVlcyBvbmx5IHJlcXVpcmUgYSBsb2dpY2FsIG9wZXJhdG9yDQogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyA4MDBtcywgdGhlIHJlc3VsdCBzaG91bGQgYmUgYSB3YXJuaW5nDQogICAgICAtIGNyaXRlcmlhOg0KICAgICAgICAgIC0gIjw9ODAwIg0KICAgIHdlaWdodDogMQ0KDQp0b3RhbF9zY29yZToNCiAgcGFzczogIjkwJSINCiAgd2FybmluZzogIjc1JSINCg==",
            "timeEnd": "2021-06-22T14:50:47.000Z",
            "timeStart": "2021-06-22T14:45:47.000Z"
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
        "id": "625413dd-a5b7-4cc8-8ce4-bcd5ead2a484",
        "source": "lighthouse-service",
        "specversion": "1.0",
        "time": "2021-06-22T14:50:48.080Z",
        "type": "sh.keptn.event.evaluation.finished",
        "shkeptncontext": "7036154d-d366-445e-b9a4-2b62649b0bfd",
        "shkeptnspecversion": "0.2.1",
        "triggeredid": "a1392234-6671-4681-8e54-65614557a667"
      },
      "heatmapLabel": "2021-06-22 16:50"
    }
    ],
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
  }
}

export const Evaluations: Trace = Trace.fromJSON(evaluationData);
