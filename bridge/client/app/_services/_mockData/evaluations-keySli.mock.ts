import { Trace } from '../../_models/trace';

const evaluationsData = [
  {
    data: {
      evaluation: {
        comparedEvents: ['182d10b8-b68d-49d4-86cd-5521352d7a42'],
        gitCommit: '',
        indicatorResults: [
          {
            displayName: 'HTTP Respone Time',
            keySli: false,
            passTargets: [
              {
                criteria: '<=0.4',
                targetValue: 0.4,
                violated: true,
              },
            ],
            score: 0,
            status: 'fail',
            value: {
              metric: 'http_response_time_seconds_main_page_sum',
              success: true,
              value: 0.40115315268094,
            },
            warningTargets: [
              {
                criteria: '<=0.1',
                targetValue: 0.1,
                violated: true,
              },
            ],
          },
          {
            displayName: 'Request Throughput',
            keySli: false,
            passTargets: [
              {
                criteria: '>=-80%',
                targetValue: 0,
                violated: false,
              },
            ],
            score: 1,
            status: 'pass',
            value: {
              metric: 'request_throughput',
              success: true,
              value: 169.23333333333332,
            },
            warningTargets: null,
          },
          {
            displayName: 'Go Routines',
            keySli: true,
            passTargets: [
              {
                criteria: '<=90',
                targetValue: 90,
                violated: false,
              },
            ],
            score: 1,
            status: 'pass',
            value: {
              metric: 'go_routines',
              success: true,
              value: 70,
            },
            warningTargets: null,
          },
          {
            displayName: 'Failing Requests',
            keySli: false,
            passTargets: [
              {
                criteria: '< 10',
                targetValue: 10,
                violated: true,
              },
            ],
            score: 0,
            status: 'fail',
            value: {
              metric: 'failing_request',
              success: true,
              value: 11,
            },
            warningTargets: null,
          },
        ],
        result: 'fail',
        score: 50,
        sloFileContent:
          'LS0tDQpzcGVjX3ZlcnNpb246ICcwLjEuMCcNCmNvbXBhcmlzb246DQogIGNvbXBhcmVfd2l0aDogInNpbmdsZV9yZXN1bHQiDQogIGluY2x1ZGVfcmVzdWx0X3dpdGhfc2NvcmU6ICJwYXNzIg0KICBhZ2dyZWdhdGVfZnVuY3Rpb246IGF2Zw0Kb2JqZWN0aXZlczoNCiAgLSBzbGk6IGh0dHBfcmVzcG9uc2VfdGltZV9zZWNvbmRzX21haW5fcGFnZV9zdW0NCiAgICBwYXNzOg0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PTEiDQogICAgd2FybmluZzoNCiAgICAgIC0gY3JpdGVyaWE6DQogICAgICAgICAgLSAiPD0wLjUiDQogIC0gc2xpOiByZXF1ZXN0X3Rocm91Z2hwdXQNCiAgICBwYXNzOg0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMDAlIg0KICAgICAgICAgIC0gIj49LTgwJSINCiAgLSBzbGk6IGdvX3JvdXRpbmVzDQogICAgcGFzczoNCiAgICAgIC0gY3JpdGVyaWE6DQogICAgICAgICAgLSAiPD0xMDAiDQp0b3RhbF9zY29yZToNCiAgcGFzczogIjkwJSINCiAgd2FybmluZzogIjc1JSI=',
        timeEnd: '2022-02-09T09:29:47.625Z',
        timeStart: '2022-02-09T09:28:55.866Z',
      },
      message: 'Evaluation failed since the calculated score of 75 is below the target value of 90',
      project: 'sockshop',
      result: 'fail',
      service: 'carts',
      stage: 'staging',
      status: 'succeeded',
    },
    id: '1500b971-bfc3-4e20-8dc1-a624e0faf963',
    source: 'lighthouse-service',
    specversion: '1.0',
    time: '2022-02-09T09:30:07.865Z',
    type: 'sh.keptn.event.evaluation.finished',
    shkeptncontext: 'da740469-9920-4e0c-b304-0fd4b18d17c2',
    shkeptnspecversion: '0.2.3',
    triggeredid: '09f53f15-c89d-4a4a-82cc-138985696172',
  },
  {
    data: {
      evaluation: {
        comparedEvents: ['182d10b8-b68d-49d4-86cd-5521352d7a42'],
        gitCommit: '',
        indicatorResults: [
          {
            displayName: 'HTTP Respone Time',
            keySli: false,
            passTargets: [
              {
                criteria: '<=0.4',
                targetValue: 0.4,
                violated: false,
              },
            ],
            score: 1,
            status: 'pass',
            value: {
              metric: 'http_response_time_seconds_main_page_sum',
              success: true,
              value: 0.20115315268094103,
            },
            warningTargets: [
              {
                criteria: '<=0.1',
                targetValue: 0.1,
                violated: true,
              },
            ],
          },
          {
            displayName: 'Request Throughput',
            keySli: false,
            passTargets: [
              {
                criteria: '>=-80%',
                targetValue: 0,
                violated: false,
              },
            ],
            score: 1,
            status: 'pass',
            value: {
              metric: 'request_throughput',
              success: true,
              value: 169.23333333333332,
            },
            warningTargets: null,
          },
          {
            displayName: 'Go Routines',
            keySli: true,
            passTargets: [
              {
                criteria: '<=90',
                targetValue: 90,
                violated: true,
              },
            ],
            score: 0,
            status: 'fail',
            value: {
              metric: 'go_routines',
              success: true,
              value: 253,
            },
            warningTargets: null,
          },
          {
            displayName: 'Failing Requests',
            keySli: false,
            passTargets: [
              {
                criteria: '< 10',
                targetValue: 10,
                violated: false,
              },
            ],
            score: 1,
            status: 'pass',
            value: {
              metric: 'failing_request',
              success: true,
              value: 0,
            },
            warningTargets: null,
          },
        ],
        result: 'fail',
        score: 75,
        sloFileContent:
          'LS0tDQpzcGVjX3ZlcnNpb246ICcwLjEuMCcNCmNvbXBhcmlzb246DQogIGNvbXBhcmVfd2l0aDogInNpbmdsZV9yZXN1bHQiDQogIGluY2x1ZGVfcmVzdWx0X3dpdGhfc2NvcmU6ICJwYXNzIg0KICBhZ2dyZWdhdGVfZnVuY3Rpb246IGF2Zw0Kb2JqZWN0aXZlczoNCiAgLSBzbGk6IGh0dHBfcmVzcG9uc2VfdGltZV9zZWNvbmRzX21haW5fcGFnZV9zdW0NCiAgICBwYXNzOg0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PTEiDQogICAgd2FybmluZzoNCiAgICAgIC0gY3JpdGVyaWE6DQogICAgICAgICAgLSAiPD0wLjUiDQogIC0gc2xpOiByZXF1ZXN0X3Rocm91Z2hwdXQNCiAgICBwYXNzOg0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMDAlIg0KICAgICAgICAgIC0gIj49LTgwJSINCiAgLSBzbGk6IGdvX3JvdXRpbmVzDQogICAgcGFzczoNCiAgICAgIC0gY3JpdGVyaWE6DQogICAgICAgICAgLSAiPD0xMDAiDQp0b3RhbF9zY29yZToNCiAgcGFzczogIjkwJSINCiAgd2FybmluZzogIjc1JSI=',
        timeEnd: '2022-02-08T13:01:08.748Z',
        timeStart: '2022-02-08T13:00:17.055Z',
      },
      message: 'Evaluation failed since the calculated score of 0 is below the warning value of 75',
      project: 'sockshop',
      result: 'fail',
      service: 'carts',
      stage: 'staging',
      status: 'succeeded',
    },
    id: '52b4b2c7-fa49-41f3-9b5c-b9aea2370bb4',
    source: 'lighthouse-service',
    specversion: '1.0',
    time: '2022-02-08T13:01:27.868Z',
    type: 'sh.keptn.event.evaluation.finished',
    shkeptncontext: '515f2595-39b6-4535-8b16-d05b837c7a61',
    shkeptnspecversion: '0.2.3',
    triggeredid: '0c796cd5-4619-4613-b7f8-d1717706acc1',
  },
  {
    data: {
      evaluation: {
        comparedEvents: ['c1b2761f-5b6d-4bdc-9bb7-4661a05ea3b2'],
        gitCommit: '',
        indicatorResults: [
          {
            displayName: 'HTTP Respone Time',
            keySli: false,
            passTargets: [
              {
                criteria: '<=0.4',
                targetValue: 0.4,
                violated: false,
              },
            ],
            score: 1,
            status: 'pass',
            value: {
              metric: 'http_response_time_seconds_main_page_sum',
              success: true,
              value: 0.20115315268094103,
            },
            warningTargets: [
              {
                criteria: '<=0.1',
                targetValue: 0.1,
                violated: true,
              },
            ],
          },
          {
            displayName: 'Request Throughput',
            keySli: false,
            passTargets: [
              {
                criteria: '>=-80%',
                targetValue: 0,
                violated: false,
              },
            ],
            score: 1,
            status: 'pass',
            value: {
              metric: 'request_throughput',
              success: true,
              value: 169.23333333333332,
            },
            warningTargets: null,
          },
          {
            displayName: 'Go Routines',
            keySli: true,
            passTargets: [
              {
                criteria: '<=90',
                targetValue: 90,
                violated: false,
              },
            ],
            score: 1,
            status: 'pass',
            value: {
              metric: 'go_routines',
              success: true,
              value: 70,
            },
            warningTargets: null,
          },
          {
            displayName: 'Failing Requests',
            keySli: false,
            passTargets: [
              {
                criteria: '< 10',
                targetValue: 10,
                violated: false,
              },
            ],
            score: 1,
            status: 'pass',
            value: {
              metric: 'failing_request',
              success: true,
              value: 0,
            },
            warningTargets: null,
          },
        ],
        result: 'pass',
        score: 100,
        sloFileContent:
          'LS0tDQpzcGVjX3ZlcnNpb246ICcwLjEuMCcNCmNvbXBhcmlzb246DQogIGNvbXBhcmVfd2l0aDogInNpbmdsZV9yZXN1bHQiDQogIGluY2x1ZGVfcmVzdWx0X3dpdGhfc2NvcmU6ICJwYXNzIg0KICBhZ2dyZWdhdGVfZnVuY3Rpb246IGF2Zw0Kb2JqZWN0aXZlczoNCiAgLSBzbGk6IGh0dHBfcmVzcG9uc2VfdGltZV9zZWNvbmRzX21haW5fcGFnZV9zdW0NCiAgICBwYXNzOg0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PTEiDQogICAgd2FybmluZzoNCiAgICAgIC0gY3JpdGVyaWE6DQogICAgICAgICAgLSAiPD0wLjUiDQogIC0gc2xpOiByZXF1ZXN0X3Rocm91Z2hwdXQNCiAgICBwYXNzOg0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMDAlIg0KICAgICAgICAgIC0gIj49LTgwJSINCiAgLSBzbGk6IGdvX3JvdXRpbmVzDQogICAgcGFzczoNCiAgICAgIC0gY3JpdGVyaWE6DQogICAgICAgICAgLSAiPD0xMDAiDQp0b3RhbF9zY29yZToNCiAgcGFzczogIjkwJSINCiAgd2FybmluZzogIjc1JSI=',
        timeEnd: '2022-02-08T12:57:01.881Z',
        timeStart: '2022-02-08T12:56:29.467Z',
      },
      project: 'sockshop',
      result: 'pass',
      service: 'carts',
      stage: 'staging',
      status: 'succeeded',
    },
    id: '182d10b8-b68d-49d4-86cd-5521352d7a42',
    source: 'lighthouse-service',
    specversion: '1.0',
    time: '2022-02-08T12:57:22.665Z',
    type: 'sh.keptn.event.evaluation.finished',
    shkeptncontext: 'fa66eea5-53a8-45b6-aefe-ef03c08b61e4',
    shkeptnspecversion: '0.2.3',
    triggeredid: '9525e089-44db-4f75-b5bb-3fd6e1141af4',
  },
  {
    data: {
      evaluation: {
        comparedEvents: ['91a77341-fe5e-43e1-a8a7-be9761b9cee5'],
        gitCommit: '',
        indicatorResults: [
          {
            displayName: 'HTTP Respone Time',
            keySli: false,
            passTargets: [
              {
                criteria: '<=1',
                targetValue: 1,
                violated: false,
              },
            ],
            score: 1,
            status: 'pass',
            value: {
              metric: 'http_response_time_seconds_main_page_sum',
              success: true,
              value: 0,
            },
            warningTargets: [
              {
                criteria: '<=0.5',
                targetValue: 0.5,
                violated: false,
              },
            ],
          },
          {
            displayName: 'Request Throughput',
            keySli: false,
            passTargets: [
              {
                criteria: '<=+100%',
                targetValue: 0,
                violated: false,
              },
              {
                criteria: '>=-80%',
                targetValue: 0,
                violated: false,
              },
            ],
            score: 1,
            status: 'pass',
            value: {
              metric: 'request_throughput',
              success: true,
              value: 0,
            },
            warningTargets: null,
          },
          {
            displayName: 'Go Routines',
            keySli: false,
            passTargets: [
              {
                criteria: '<=100',
                targetValue: 100,
                violated: false,
              },
            ],
            score: 1,
            status: 'pass',
            value: {
              metric: 'go_routines',
              success: true,
              value: 10,
            },
            warningTargets: null,
          },
        ],
        result: 'pass',
        score: 100,
        sloFileContent:
          'LS0tDQpzcGVjX3ZlcnNpb246ICcwLjEuMCcNCmNvbXBhcmlzb246DQogIGNvbXBhcmVfd2l0aDogInNpbmdsZV9yZXN1bHQiDQogIGluY2x1ZGVfcmVzdWx0X3dpdGhfc2NvcmU6ICJwYXNzIg0KICBhZ2dyZWdhdGVfZnVuY3Rpb246IGF2Zw0Kb2JqZWN0aXZlczoNCiAgLSBzbGk6IGh0dHBfcmVzcG9uc2VfdGltZV9zZWNvbmRzX21haW5fcGFnZV9zdW0NCiAgICBwYXNzOg0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PTEiDQogICAgd2FybmluZzoNCiAgICAgIC0gY3JpdGVyaWE6DQogICAgICAgICAgLSAiPD0wLjUiDQogIC0gc2xpOiByZXF1ZXN0X3Rocm91Z2hwdXQNCiAgICBwYXNzOg0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMDAlIg0KICAgICAgICAgIC0gIj49LTgwJSINCiAgLSBzbGk6IGdvX3JvdXRpbmVzDQogICAgcGFzczoNCiAgICAgIC0gY3JpdGVyaWE6DQogICAgICAgICAgLSAiPD0xMDAiDQp0b3RhbF9zY29yZToNCiAgcGFzczogIjkwJSINCiAgd2FybmluZzogIjc1JSI=',
        timeEnd: '2022-02-08T12:42:41.467Z',
        timeStart: '2022-02-08T12:42:12.657Z',
      },
      project: 'sockshop',
      result: 'pass',
      service: 'carts',
      stage: 'staging',
      status: 'succeeded',
    },
    id: 'c1b2761f-5b6d-4bdc-9bb7-4661a05ea3b2',
    source: 'lighthouse-service',
    specversion: '1.0',
    time: '2022-02-08T12:43:00.256Z',
    type: 'sh.keptn.event.evaluation.finished',
    shkeptncontext: '07091ebf-4cbf-4b56-ad33-c2ffc21f619e',
    shkeptnspecversion: '0.2.3',
    triggeredid: 'fd2fd1e6-19b9-44c2-8731-759ec628bf3d',
  },
  {
    data: {
      evaluation: {
        comparedEvents: ['91a77341-fe5e-43e1-a8a7-be9761b9cee5'],
        gitCommit: '',
        indicatorResults: [
          {
            displayName: 'HTTP Respone Time',
            keySli: false,
            passTargets: [
              {
                criteria: '<=1',
                targetValue: 1,
                violated: true,
              },
            ],
            score: 0,
            status: 'fail',
            value: {
              metric: 'http_response_time_seconds_main_page_sum',
              success: true,
              value: 2.0004689763871784,
            },
            warningTargets: [
              {
                criteria: '<=0.5',
                targetValue: 0.5,
                violated: true,
              },
            ],
          },
          {
            displayName: 'Request Throughput',
            keySli: false,
            passTargets: [
              {
                criteria: '<=+100%',
                targetValue: 0,
                violated: true,
              },
              {
                criteria: '>=-80%',
                targetValue: 0,
                violated: false,
              },
            ],
            score: 0,
            status: 'fail',
            value: {
              metric: 'request_throughput',
              success: true,
              value: 20.593531972789116,
            },
            warningTargets: null,
          },
          {
            displayName: 'Go Routines',
            keySli: false,
            passTargets: [
              {
                criteria: '<=100',
                targetValue: 100,
                violated: false,
              },
            ],
            score: 1,
            status: 'pass',
            value: {
              metric: 'go_routines',
              success: true,
              value: 77,
            },
            warningTargets: null,
          },
        ],
        result: 'fail',
        score: 33.99999999999999,
        sloFileContent:
          'LS0tDQpzcGVjX3ZlcnNpb246ICcwLjEuMCcNCmNvbXBhcmlzb246DQogIGNvbXBhcmVfd2l0aDogInNpbmdsZV9yZXN1bHQiDQogIGluY2x1ZGVfcmVzdWx0X3dpdGhfc2NvcmU6ICJwYXNzIg0KICBhZ2dyZWdhdGVfZnVuY3Rpb246IGF2Zw0Kb2JqZWN0aXZlczoNCiAgLSBzbGk6IGh0dHBfcmVzcG9uc2VfdGltZV9zZWNvbmRzX21haW5fcGFnZV9zdW0NCiAgICBwYXNzOg0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PTEiDQogICAgd2FybmluZzoNCiAgICAgIC0gY3JpdGVyaWE6DQogICAgICAgICAgLSAiPD0wLjUiDQogIC0gc2xpOiByZXF1ZXN0X3Rocm91Z2hwdXQNCiAgICBwYXNzOg0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMDAlIg0KICAgICAgICAgIC0gIj49LTgwJSINCiAgLSBzbGk6IGdvX3JvdXRpbmVzDQogICAgcGFzczoNCiAgICAgIC0gY3JpdGVyaWE6DQogICAgICAgICAgLSAiPD0xMDAiDQp0b3RhbF9zY29yZToNCiAgcGFzczogIjkwJSINCiAgd2FybmluZzogIjc1JSI=',
        timeEnd: '2022-02-08T12:33:27.282Z',
        timeStart: '2022-02-08T12:32:39.065Z',
      },
      message: 'Evaluation failed since the calculated score of 33.99999999999999 is below the warning value of 75',
      project: 'sockshop',
      result: 'fail',
      service: 'carts',
      stage: 'staging',
      status: 'succeeded',
    },
    id: '8a549059-8dcd-43ea-adff-b7c2ea9a0d99',
    source: 'lighthouse-service',
    specversion: '1.0',
    time: '2022-02-08T12:33:45.652Z',
    type: 'sh.keptn.event.evaluation.finished',
    shkeptncontext: '124fa9d0-cbfd-493b-b9eb-1a8d806a133f',
    shkeptnspecversion: '0.2.3',
    triggeredid: 'f368c6d6-56d2-44ff-8f5a-3178b73464ab',
  },
  {
    data: {
      evaluation: {
        comparedEvents: ['0c275de6-fbde-457e-b2f2-60f32a6acf54'],
        gitCommit: '',
        indicatorResults: [
          {
            displayName: 'HTTP Respone Time',
            keySli: false,
            passTargets: [
              {
                criteria: '<=1',
                targetValue: 1,
                violated: false,
              },
            ],
            score: 1,
            status: 'pass',
            value: {
              metric: 'http_response_time_seconds_main_page_sum',
              success: true,
              value: 0,
            },
            warningTargets: [
              {
                criteria: '<=0.5',
                targetValue: 0.5,
                violated: false,
              },
            ],
          },
          {
            displayName: 'Request Throughput',
            keySli: false,
            passTargets: [
              {
                criteria: '<=+100%',
                targetValue: 0,
                violated: false,
              },
              {
                criteria: '>=-80%',
                targetValue: 0,
                violated: false,
              },
            ],
            score: 1,
            status: 'pass',
            value: {
              metric: 'request_throughput',
              success: true,
              value: 0,
            },
            warningTargets: null,
          },
          {
            displayName: 'Go Routines',
            keySli: false,
            passTargets: [
              {
                criteria: '<=100',
                targetValue: 100,
                violated: false,
              },
            ],
            score: 1,
            status: 'pass',
            value: {
              metric: 'go_routines',
              success: true,
              value: 8,
            },
            warningTargets: null,
          },
        ],
        result: 'pass',
        score: 100,
        sloFileContent:
          'LS0tDQpzcGVjX3ZlcnNpb246ICcwLjEuMCcNCmNvbXBhcmlzb246DQogIGNvbXBhcmVfd2l0aDogInNpbmdsZV9yZXN1bHQiDQogIGluY2x1ZGVfcmVzdWx0X3dpdGhfc2NvcmU6ICJwYXNzIg0KICBhZ2dyZWdhdGVfZnVuY3Rpb246IGF2Zw0Kb2JqZWN0aXZlczoNCiAgLSBzbGk6IGh0dHBfcmVzcG9uc2VfdGltZV9zZWNvbmRzX21haW5fcGFnZV9zdW0NCiAgICBwYXNzOg0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PTEiDQogICAgd2FybmluZzoNCiAgICAgIC0gY3JpdGVyaWE6DQogICAgICAgICAgLSAiPD0wLjUiDQogIC0gc2xpOiByZXF1ZXN0X3Rocm91Z2hwdXQNCiAgICBwYXNzOg0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMDAlIg0KICAgICAgICAgIC0gIj49LTgwJSINCiAgLSBzbGk6IGdvX3JvdXRpbmVzDQogICAgcGFzczoNCiAgICAgIC0gY3JpdGVyaWE6DQogICAgICAgICAgLSAiPD0xMDAiDQp0b3RhbF9zY29yZToNCiAgcGFzczogIjkwJSINCiAgd2FybmluZzogIjc1JSI=',
        timeEnd: '2022-02-08T12:30:40.397Z',
        timeStart: '2022-02-08T12:30:12.455Z',
      },
      project: 'sockshop',
      result: 'pass',
      service: 'carts',
      stage: 'staging',
      status: 'succeeded',
    },
    id: '91a77341-fe5e-43e1-a8a7-be9761b9cee5',
    source: 'lighthouse-service',
    specversion: '1.0',
    time: '2022-02-08T12:30:57.052Z',
    type: 'sh.keptn.event.evaluation.finished',
    shkeptncontext: '420f0822-9e51-48ec-8856-345c0b347924',
    shkeptnspecversion: '0.2.3',
    triggeredid: 'd9492566-2c58-4e10-bdb5-2dfc99518fc2',
  },
  {
    data: {
      evaluation: {
        comparedEvents: ['0c275de6-fbde-457e-b2f2-60f32a6acf54'],
        gitCommit: '',
        indicatorResults: [
          {
            displayName: 'Response time P95',
            keySli: false,
            passTargets: [
              {
                criteria: '<=+10%',
                targetValue: 0,
                violated: true,
              },
              {
                criteria: '<600',
                targetValue: 0,
                violated: true,
              },
            ],
            score: 0,
            status: 'fail',
            value: {
              message: 'Dynatrace Metrics API returned zero data points',
              metric: 'response_time_p95',
              success: false,
              value: 0,
            },
            warningTargets: [
              {
                criteria: '<=800',
                targetValue: 0,
                violated: true,
              },
            ],
          },
        ],
        result: 'fail',
        score: 0,
        sloFileContent:
          'LS0tCnNwZWNfdmVyc2lvbjogIjEuMCIKY29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciCiAgY29tcGFyZV93aXRoOiAic2luZ2xlX3Jlc3VsdCIKICBpbmNsdWRlX3Jlc3VsdF93aXRoX3Njb3JlOiAicGFzcyIKICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxCmZpbHRlcjoKb2JqZWN0aXZlczoKICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1IgogICAgZGlzcGxheU5hbWU6ICJSZXNwb25zZSB0aW1lIFA5NSIKICAgIGtleV9zbGk6IGZhbHNlCiAgICBwYXNzOiAgICAgICAgICAgICAjIHBhc3MgaWYgKHJlbGF0aXZlIGNoYW5nZSA8PSAxMCUgQU5EIGFic29sdXRlIHZhbHVlIGlzIDwgNjAwbXMpCiAgICAgIC0gY3JpdGVyaWE6CiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykKICAgICAgICAgIC0gIjw2MDAiICAgICMgYWJzb2x1dGUgdmFsdWVzIG9ubHkgcmVxdWlyZSBhIGxvZ2ljYWwgb3BlcmF0b3IKICAgIHdhcm5pbmc6ICAgICAgICAgICMgaWYgdGhlIHJlc3BvbnNlIHRpbWUgaXMgYmVsb3cgODAwbXMsIHRoZSByZXN1bHQgc2hvdWxkIGJlIGEgd2FybmluZwogICAgICAtIGNyaXRlcmlhOgogICAgICAgICAgLSAiPD04MDAiCiAgICB3ZWlnaHQ6IDEKdG90YWxfc2NvcmU6CiAgcGFzczogIjkwJSIKICB3YXJuaW5nOiAiNzUlIg==',
        timeEnd: '2022-02-08T12:15:02.695Z',
        timeStart: '2022-02-08T12:14:34.757Z',
      },
      labels: {
        DtCreds: 'dynatrace',
      },
      message: 'Evaluation failed since the calculated score of 0 is below the warning value of 75',
      project: 'sockshop',
      result: 'fail',
      service: 'carts',
      stage: 'staging',
      status: 'succeeded',
    },
    id: 'e074893e-a7f9-4fa8-9e7e-898937a3d2b6',
    source: 'lighthouse-service',
    specversion: '1.0',
    time: '2022-02-08T12:17:18.553Z',
    type: 'sh.keptn.event.evaluation.finished',
    shkeptncontext: '9c5065bb-7ad5-4d02-b976-f75ec13ac0d1',
    shkeptnspecversion: '0.2.3',
    triggeredid: '5203cda7-7aad-43af-bba1-96844b8acb07',
  },
  {
    data: {
      evaluation: {
        gitCommit: '',
        indicatorResults: null,
        result: 'pass',
        score: 0,
        sloFileContent: '',
        timeEnd: '2022-02-08T12:02:03.779Z',
        timeStart: '2022-02-08T12:01:34.471Z',
      },
      labels: {
        DtCreds: 'dynatrace',
      },
      message: 'no evaluation performed by lighthouse because no SLO file configured for project sockshop',
      project: 'sockshop',
      result: 'pass',
      service: 'carts',
      stage: 'staging',
      status: 'succeeded',
    },
    id: '0c275de6-fbde-457e-b2f2-60f32a6acf54',
    source: 'lighthouse-service',
    specversion: '1.0',
    time: '2022-02-08T12:04:10.458Z',
    type: 'sh.keptn.event.evaluation.finished',
    shkeptncontext: '6305ad09-6a69-4b4f-8a6e-0b98f332cec3',
    shkeptnspecversion: '0.2.3',
    triggeredid: '841c01b6-5687-4f8a-9f16-378fbe3a8d17',
  },
  {
    data: {
      evaluation: {
        gitCommit: '',
        indicatorResults: null,
        result: 'pass',
        score: 0,
        sloFileContent: '',
        timeEnd: '2022-02-08T11:54:50.988Z',
        timeStart: '2022-02-08T11:54:23.963Z',
      },
      labels: {
        DtCreds: 'dynatrace',
      },
      message: 'no evaluation performed by lighthouse because no SLO file configured for project sockshop',
      project: 'sockshop',
      result: 'pass',
      service: 'carts',
      stage: 'staging',
      status: 'succeeded',
    },
    id: '880ac00e-3f65-476f-a1f9-4325919146ac',
    source: 'lighthouse-service',
    specversion: '1.0',
    time: '2022-02-08T11:56:58.953Z',
    type: 'sh.keptn.event.evaluation.finished',
    shkeptncontext: '93046ada-f07e-44e0-bd14-3f471c71467c',
    shkeptnspecversion: '0.2.3',
    triggeredid: '443055e4-db15-4858-ae72-9b6455a81b92',
  },
  {
    data: {
      evaluation: {
        gitCommit: '',
        indicatorResults: null,
        result: 'pass',
        score: 0,
        sloFileContent: '',
        timeEnd: '2022-02-08T11:44:19.718Z',
        timeStart: '2022-02-08T11:43:51.456Z',
      },
      labels: {
        DtCreds: 'dynatrace',
      },
      message: 'no evaluation performed by lighthouse because no SLO file configured for project sockshop',
      project: 'sockshop',
      result: 'pass',
      service: 'carts',
      stage: 'staging',
      status: 'succeeded',
    },
    id: '25ab0f26-e6d8-48d5-a08f-08c8a136a688',
    source: 'lighthouse-service',
    specversion: '1.0',
    time: '2022-02-08T11:46:35.853Z',
    type: 'sh.keptn.event.evaluation.finished',
    shkeptncontext: '8f884b2a-2197-4e2f-8284-170ea0a66579',
    shkeptnspecversion: '0.2.3',
    triggeredid: '991eba72-c520-4da5-ba95-d5876101393c',
  },
];

const evaluations: Trace[] = evaluationsData.map((evaluationData) => Trace.fromJSON(evaluationData));
export { evaluations as EvaluationsKeySliMock };
