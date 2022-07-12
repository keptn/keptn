import { TestUtils } from '../../../_utils/test.utils';

const evaluationResultsData = [
  {
    data: {
      evaluation: {
        comparedEvents: ['01100b8e-6a03-45ba-b140-78142187f24f'],
        indicatorResults: [
          {
            displayName: 'Response time P95',
            keySli: false,
            passTargets: [
              {
                criteria: '<=+10%',
                targetValue: 400.4072038456992,
                violated: true,
              },
              {
                criteria: '<600',
                targetValue: 600,
                violated: false,
              },
            ],
            score: 0.5,
            status: 'warning',
            value: {
              metric: 'response_time_p95',
              success: true,
              value: 414.2608775629195,
            },
            warningTargets: [
              {
                criteria: '<=800',
                targetValue: 800,
                violated: false,
              },
            ],
          },
        ],
        result: 'fail',
        score: 50,
        sloFileContent:
          'LS0tDQpzcGVjX3ZlcnNpb246ICIxLjAiDQpjb21wYXJpc29uOg0KICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciDQogIGNvbXBhcmVfd2l0aDogInNpbmdsZV9yZXN1bHQiDQogIGluY2x1ZGVfcmVzdWx0X3dpdGhfc2NvcmU6ICJwYXNzIg0KICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxDQpmaWx0ZXI6DQpvYmplY3RpdmVzOg0KICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1Ig0KICAgIGRpc3BsYXlOYW1lOiAiUmVzcG9uc2UgdGltZSBQOTUiDQogICAga2V5X3NsaTogZmFsc2UNCiAgICBwYXNzOiAgICAgICAgICAgICAjIHBhc3MgaWYgKHJlbGF0aXZlIGNoYW5nZSA8PSAxMCUgQU5EIGFic29sdXRlIHZhbHVlIGlzIDwgNjAwbXMpDQogICAgICAtIGNyaXRlcmlhOg0KICAgICAgICAgIC0gIjw9KzEwJSIgICMgcmVsYXRpdmUgdmFsdWVzIHJlcXVpcmUgYSBwcmVmaXhlZCBzaWduIChwbHVzIG9yIG1pbnVzKQ0KICAgICAgICAgIC0gIjw2MDAiICAgICMgYWJzb2x1dGUgdmFsdWVzIG9ubHkgcmVxdWlyZSBhIGxvZ2ljYWwgb3BlcmF0b3INCiAgICB3YXJuaW5nOiAgICAgICAgICAjIGlmIHRoZSByZXNwb25zZSB0aW1lIGlzIGJlbG93IDgwMG1zLCB0aGUgcmVzdWx0IHNob3VsZCBiZSBhIHdhcm5pbmcNCiAgICAgIC0gY3JpdGVyaWE6DQogICAgICAgICAgLSAiPD04MDAiDQogICAgd2VpZ2h0OiAxDQp0b3RhbF9zY29yZToNCiAgcGFzczogIjkwJSINCiAgd2FybmluZzogIjc1JSI=',
        timeEnd: '2021-11-11T15:30:39Z',
        timeStart: '2021-11-11T15:28:22Z',
      },
      labels: {
        DtCreds: 'dynatrace',
      },
      message: 'Evaluation failed since the calculated score of 50 is below the warning value of 75',
      project: 'sockshop',
      result: 'fail',
      service: 'carts',
      stage: 'staging',
      status: 'succeeded',
    },
    id: '37aba4aa-7dbd-4946-b53d-f25c3b3703c5',
    source: 'lighthouse-service',
    specversion: '1.0',
    time: '2021-11-11T15:31:48.791Z',
    type: 'sh.keptn.event.evaluation.finished',
    shkeptncontext: '08547346-c845-4f49-acab-9f0b9301067e',
    shkeptnspecversion: '0.2.3',
    triggeredid: '8e22cc44-9456-4514-83ce-22eb91fa0121',
  },
  {
    data: {
      evaluation: {
        comparedEvents: ['80da2703-0df4-4995-80f6-39cae6fec4e5'],
        indicatorResults: [
          {
            displayName: 'Response time P95',
            keySli: false,
            passTargets: [
              {
                criteria: '<=+10%',
                targetValue: 417.60840032549925,
                violated: false,
              },
              {
                criteria: '<600',
                targetValue: 600,
                violated: false,
              },
            ],
            score: 1,
            status: 'pass',
            value: {
              metric: 'response_time_p95',
              success: true,
              value: 364.00654895063565,
            },
            warningTargets: [
              {
                criteria: '<=800',
                targetValue: 800,
                violated: false,
              },
            ],
          },
        ],
        result: 'pass',
        score: 100,
        sloFileContent:
          'LS0tDQpzcGVjX3ZlcnNpb246ICIxLjAiDQpjb21wYXJpc29uOg0KICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciDQogIGNvbXBhcmVfd2l0aDogInNpbmdsZV9yZXN1bHQiDQogIGluY2x1ZGVfcmVzdWx0X3dpdGhfc2NvcmU6ICJwYXNzIg0KICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxDQpmaWx0ZXI6DQpvYmplY3RpdmVzOg0KICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1Ig0KICAgIGRpc3BsYXlOYW1lOiAiUmVzcG9uc2UgdGltZSBQOTUiDQogICAga2V5X3NsaTogZmFsc2UNCiAgICBwYXNzOiAgICAgICAgICAgICAjIHBhc3MgaWYgKHJlbGF0aXZlIGNoYW5nZSA8PSAxMCUgQU5EIGFic29sdXRlIHZhbHVlIGlzIDwgNjAwbXMpDQogICAgICAtIGNyaXRlcmlhOg0KICAgICAgICAgIC0gIjw9KzEwJSIgICMgcmVsYXRpdmUgdmFsdWVzIHJlcXVpcmUgYSBwcmVmaXhlZCBzaWduIChwbHVzIG9yIG1pbnVzKQ0KICAgICAgICAgIC0gIjw2MDAiICAgICMgYWJzb2x1dGUgdmFsdWVzIG9ubHkgcmVxdWlyZSBhIGxvZ2ljYWwgb3BlcmF0b3INCiAgICB3YXJuaW5nOiAgICAgICAgICAjIGlmIHRoZSByZXNwb25zZSB0aW1lIGlzIGJlbG93IDgwMG1zLCB0aGUgcmVzdWx0IHNob3VsZCBiZSBhIHdhcm5pbmcNCiAgICAgIC0gY3JpdGVyaWE6DQogICAgICAgICAgLSAiPD04MDAiDQogICAgd2VpZ2h0OiAxDQp0b3RhbF9zY29yZToNCiAgcGFzczogIjkwJSINCiAgd2FybmluZzogIjc1JSI=',
        timeEnd: '2021-11-11T13:13:13Z',
        timeStart: '2021-11-11T13:10:58Z',
      },
      labels: {
        DtCreds: 'dynatrace',
      },
      project: 'sockshop',
      result: 'pass',
      service: 'carts',
      stage: 'staging',
      status: 'succeeded',
    },
    id: '01100b8e-6a03-45ba-b140-78142187f24f',
    source: 'lighthouse-service',
    specversion: '1.0',
    time: '2021-11-11T13:14:22.107Z',
    type: 'sh.keptn.event.evaluation.finished',
    shkeptncontext: '76f0b0af-0290-458e-82da-56bec6ec5868',
    shkeptnspecversion: '0.2.3',
    triggeredid: 'c6d579ba-e34a-44c1-826c-db8425dc8489',
  },
  {
    data: {
      evaluation: {
        indicatorResults: [
          {
            displayName: 'Response time P95',
            keySli: false,
            passTargets: [
              {
                criteria: '<=+10%',
                targetValue: 0,
                violated: false,
              },
              {
                criteria: '<600',
                targetValue: 600,
                violated: false,
              },
            ],
            score: 1,
            status: 'pass',
            value: {
              metric: 'response_time_p95',
              success: true,
              value: 379.64400029590837,
            },
            warningTargets: [
              {
                criteria: '<=800',
                targetValue: 800,
                violated: false,
              },
            ],
          },
        ],
        result: 'pass',
        score: 100,
        sloFileContent:
          'LS0tDQpzcGVjX3ZlcnNpb246ICIxLjAiDQpjb21wYXJpc29uOg0KICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciDQogIGNvbXBhcmVfd2l0aDogInNpbmdsZV9yZXN1bHQiDQogIGluY2x1ZGVfcmVzdWx0X3dpdGhfc2NvcmU6ICJwYXNzIg0KICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxDQpmaWx0ZXI6DQpvYmplY3RpdmVzOg0KICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1Ig0KICAgIGRpc3BsYXlOYW1lOiAiUmVzcG9uc2UgdGltZSBQOTUiDQogICAga2V5X3NsaTogZmFsc2UNCiAgICBwYXNzOiAgICAgICAgICAgICAjIHBhc3MgaWYgKHJlbGF0aXZlIGNoYW5nZSA8PSAxMCUgQU5EIGFic29sdXRlIHZhbHVlIGlzIDwgNjAwbXMpDQogICAgICAtIGNyaXRlcmlhOg0KICAgICAgICAgIC0gIjw9KzEwJSIgICMgcmVsYXRpdmUgdmFsdWVzIHJlcXVpcmUgYSBwcmVmaXhlZCBzaWduIChwbHVzIG9yIG1pbnVzKQ0KICAgICAgICAgIC0gIjw2MDAiICAgICMgYWJzb2x1dGUgdmFsdWVzIG9ubHkgcmVxdWlyZSBhIGxvZ2ljYWwgb3BlcmF0b3INCiAgICB3YXJuaW5nOiAgICAgICAgICAjIGlmIHRoZSByZXNwb25zZSB0aW1lIGlzIGJlbG93IDgwMG1zLCB0aGUgcmVzdWx0IHNob3VsZCBiZSBhIHdhcm5pbmcNCiAgICAgIC0gY3JpdGVyaWE6DQogICAgICAgICAgLSAiPD04MDAiDQogICAgd2VpZ2h0OiAxDQp0b3RhbF9zY29yZToNCiAgcGFzczogIjkwJSINCiAgd2FybmluZzogIjc1JSI=',
        timeEnd: '2021-11-09T15:20:53Z',
        timeStart: '2021-11-09T15:18:54Z',
      },
      labels: {
        DtCreds: 'dynatrace',
      },
      project: 'sockshop',
      result: 'pass',
      service: 'carts',
      stage: 'staging',
      status: 'succeeded',
    },
    id: '80da2703-0df4-4995-80f6-39cae6fec4e5',
    source: 'lighthouse-service',
    specversion: '1.0',
    time: '2021-11-09T15:22:56.501Z',
    type: 'sh.keptn.event.evaluation.finished',
    shkeptncontext: '77baf26f-f64d-4a68-9ab5-efde9276ee73',
    shkeptnspecversion: '0.2.3',
    triggeredid: 'b22b4575-248c-44fa-bc27-8167ea765db4',
  },
];

const results = TestUtils.mapTraces(evaluationResultsData);
export { results as EvaluationResultsResponseDataMock };
