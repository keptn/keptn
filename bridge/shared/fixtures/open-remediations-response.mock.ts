const openRemediationsResponseMock = {
  states: [
    {
      name: 'remediation',
      problemTitle: 'Failure rate increase',
      service: 'carts',
      project: 'sockshop',
      time: '2021-11-04T04:51:21.557Z',
      shkeptncontext: '35383737-3630-4639-b037-353138323631',
      state: 'started',
      stages: [
        {
          name: 'production',
          state: 'started',
          latestEvent: {
            type: 'sh.keptn.event.get-action.started',
            id: '7448420f-5b15-4777-9d39-cc8308e2b0c3',
            time: '2021-11-04T04:51:23.266Z',
          },
        },
      ],
    },
  ],
};

const remediationTraceResponse = {
  events: [
    {
      data: {
        labels: {
          'Problem URL': 'https://myURL.com',
        },
        problem: {
          ImpactedEntity: 'Failure rate increase on Web service ItemsController',
          PID: '5877606907518261221_1636001100000V2',
          ProblemDetails: {
            displayName: 'P-21118',
            endTime: -1,
            hasRootCause: true,
            id: '5877606907518261221_1636001100000V2',
            impactLevel: 'SERVICE',
            severityLevel: 'ERROR',
            startTime: 1636001400000,
            status: 'OPEN',
          },
          ProblemID: 'P-21118',
          ProblemTitle: 'Failure rate increase',
          ProblemURL: 'https://myURL.com',
          State: 'OPEN',
          Tags: 'keptn_service:carts, keptn_deployment:primary, keptn_project:sockshop, keptn_stage:production',
        },
        project: 'sockshop',
        service: 'carts',
        stage: 'production',
        temporaryData: {
          distributor: {
            subscriptionID: '',
          },
        },
      },
      id: 'f17573df-4efb-4143-9aa9-a41169920cf7',
      source: 'dynatrace-service',
      specversion: '1.0',
      time: '2021-11-04T04:51:20.329Z',
      type: 'sh.keptn.event.production.remediation.triggered',
      shkeptncontext: '35383737-3630-4639-b037-353138323631',
      shkeptnspecversion: '0.2.1',
    },
  ],
  pageSize: 50,
  totalCount: 1,
};

const remediationTracesResponse = {
  events: [
    {
      data: {
        labels: {
          DtCreds: 'dynatrace',
          'Problem URL': 'https://myUrl.com',
        },
        message: 'No more actions defined for problem type  in remediation.yaml file.',
        project: 'sockshop',
        result: 'fail',
        service: 'carts',
        stage: 'production',
        status: 'succeeded',
      },
      id: 'ab37343e-4307-45cf-bceb-b47a333f60c8',
      source: 'shipyard-controller',
      specversion: '1.0',
      time: '2021-11-18T13:08:49.797Z',
      type: 'sh.keptn.event.production.remediation.finished',
      shkeptncontext: '35383737-3630-4639-b037-353138323631',
      shkeptnspecversion: '0.2.3',
      triggeredid: 'bd783f59-3325-4d97-ae61-253e93ea259a',
    },
    {
      data: {
        labels: {
          DtCreds: 'dynatrace',
          'Problem URL': 'https://myUrl.com',
        },
        message: 'No more actions defined for problem type  in remediation.yaml file.',
        project: 'sockshop',
        result: 'fail',
        service: 'carts',
        stage: 'production',
        status: 'succeeded',
      },
      id: 'a160833d-12c7-43b0-a67b-81912e255192',
      source: 'remediation-service',
      specversion: '1.0',
      time: '2021-11-18T13:08:49.692Z',
      type: 'sh.keptn.event.get-action.finished',
      shkeptncontext: '35383737-3630-4639-b037-353138323631',
      triggeredid: '3009ef9b-61a2-4b3c-87a1-81336ecd4574',
    },
    {
      data: {
        labels: {
          DtCreds: 'dynatrace',
          'Problem URL': 'https://myUrl.com',
        },
        project: 'sockshop',
        service: 'carts',
        stage: 'production',
      },
      id: 'd08dea0a-034b-4bea-974b-2a24a5ce5620',
      source: 'remediation-service',
      specversion: '1.0',
      time: '2021-11-18T13:08:46.595Z',
      type: 'sh.keptn.event.get-action.started',
      shkeptncontext: '35383737-3630-4639-b037-353138323631',
      triggeredid: '3009ef9b-61a2-4b3c-87a1-81336ecd4574',
    },
    {
      data: {
        ActionIndex: 1,
        action: {
          action: 'toggle-feature',
          description: 'Toogle feature flag EnableItemCache to ON',
          name: 'Toogle feature flag',
          value: {
            EnableItemCache: 'on',
          },
        },
        evaluation: {
          comparedEvents: ['e2800526-c0f9-4817-9aec-9a3037d42837'],
          gitCommit: '',
          indicatorResults: [
            {
              displayName: 'Response time P90',
              keySli: false,
              passTargets: [
                {
                  criteria: '<=+10%',
                  targetValue: 0,
                  violated: true,
                },
                {
                  criteria: '<1000',
                  targetValue: 0,
                  violated: true,
                },
              ],
              score: 0,
              status: 'fail',
              value: {
                message:
                  'Dynatrace Metrics API returned 0 result values, expected 1 for query: .\nPlease ensure the response contains exactly one value (e.g., by using :merge(0):avg for the metric). Here is the output for troubleshooting: {"metricId":"builtin:service.response.time:merge(\\"dt.entity.service\\"):percentile(90)","data":[]}',
                metric: 'response_time_p90',
                success: false,
                value: 0,
              },
              warningTargets: [
                {
                  criteria: '<=1200',
                  targetValue: 0,
                  violated: true,
                },
              ],
            },
            {
              displayName: 'Problem open',
              keySli: true,
              passTargets: [
                {
                  criteria: '=0',
                  targetValue: 0,
                  violated: true,
                },
              ],
              score: 0,
              status: 'fail',
              value: {
                message:
                  'Problems API request http://myUrl.com was not successful: Dynatrace API returned error 403: Token is missing required scope. Use one of: problems.read (Read problems)',
                metric: 'problem_open',
                success: false,
                value: 0,
              },
              warningTargets: null,
            },
          ],
          result: 'fail',
          score: 0,
          sloFileContent:
            'LS0tCnNwZWNfdmVyc2lvbjogIjEuMCIKY29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciCiAgY29tcGFyZV93aXRoOiAic2luZ2xlX3Jlc3VsdCIKICBpbmNsdWRlX3Jlc3VsdF93aXRoX3Njb3JlOiAicGFzcyIKICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxCmZpbHRlcjoKb2JqZWN0aXZlczoKICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDkwIgogICAgZGlzcGxheU5hbWU6ICJSZXNwb25zZSB0aW1lIFA5MCIKICAgIGtleV9zbGk6IGZhbHNlCiAgICBwYXNzOiAgICAgICAgICAgICAjIHBhc3MgaWYgKHJlbGF0aXZlIGNoYW5nZSA8PSAxMCUgQU5EIGFic29sdXRlIHZhbHVlIGlzIDwgMTAwMCkKICAgICAgLSBjcml0ZXJpYToKICAgICAgICAgIC0gIjw9KzEwJSIgICMgcmVsYXRpdmUgdmFsdWVzIHJlcXVpcmUgYSBwcmVmaXhlZCBzaWduIChwbHVzIG9yIG1pbnVzKQogICAgICAgICAgLSAiPDEwMDAiICAgIyBhYnNvbHV0ZSB2YWx1ZXMgb25seSByZXF1aXJlIGEgbG9naWNhbCBvcGVyYXRvcgogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyAxMjAwbXMsIHRoZSByZXN1bHQgc2hvdWxkIGJlIGEgd2FybmluZwogICAgICAtIGNyaXRlcmlhOgogICAgICAgICAgLSAiPD0xMjAwIgogICAgd2VpZ2h0OiAxCiAgLSBzbGk6ICJwcm9ibGVtX29wZW4iCiAgICBkaXNwbGF5TmFtZTogIlByb2JsZW0gb3BlbiIKICAgIGtleV9zbGk6IHRydWUKICAgIHBhc3M6CiAgICAgIC0gY3JpdGVyaWE6CiAgICAgICAgICAtICI9MCIKICAgIHdlaWdodDogMQp0b3RhbF9zY29yZToKICBwYXNzOiAiOTAlIgogIHdhcm5pbmc6ICI0MCUiCg==',
          timeEnd: '2021-11-18T13:08:20.106Z',
          timeStart: '2021-11-18T12:53:20.106Z',
        },
        'get-action': null,
        labels: {
          DtCreds: 'dynatrace',
          'Problem URL': 'https://myUrl.com',
        },
        message: '',
        problem: {
          ImpactedEntity:
            'Connectivity problem on Process SpringBoot carts works.weave.socks.cart.CartApplication carts-* (carts-655bbd7474-xjsht)',
          PID: '-4365226942138395796_1635983460000V2',
          ProblemDetails: {
            displayName: 'P-21115',
            endTime: -1,
            hasRootCause: false,
            id: '-4365226942138395796_1635983460000V2',
            impactLevel: 'INFRASTRUCTURE',
            severityLevel: 'ERROR',
            startTime: 1635983760000,
            status: 'OPEN',
          },
          ProblemID: 'P-21115',
          ProblemTitle: 'Response time degradation',
          ProblemURL: 'https://myUrl.com',
          State: 'OPEN',
          Tags: 'keptn_service:carts, keptn_project:sockshop, keptn_stage:production, keptn_deployment:direct',
        },
        project: 'sockshop',
        result: '',
        service: 'carts',
        stage: 'production',
        status: '',
        temporaryData: {
          distributor: {
            subscriptionID: '',
          },
        },
      },
      id: '3009ef9b-61a2-4b3c-87a1-81336ecd4574',
      source: 'shipyard-controller',
      specversion: '1.0',
      time: '2021-11-18T13:08:46.590Z',
      type: 'sh.keptn.event.get-action.triggered',
      shkeptncontext: '35383737-3630-4639-b037-353138323631',
      shkeptnspecversion: '0.2.3',
    },
    {
      data: {
        ActionIndex: 1,
        action: {
          action: 'toggle-feature',
          description: 'Toogle feature flag EnableItemCache to ON',
          name: 'Toogle feature flag',
          value: {
            EnableItemCache: 'on',
          },
        },
        evaluation: {
          comparedEvents: ['e2800526-c0f9-4817-9aec-9a3037d42837'],
          gitCommit: '',
          indicatorResults: [
            {
              displayName: 'Response time P90',
              keySli: false,
              passTargets: [
                {
                  criteria: '<=+10%',
                  targetValue: 0,
                  violated: true,
                },
                {
                  criteria: '<1000',
                  targetValue: 0,
                  violated: true,
                },
              ],
              score: 0,
              status: 'fail',
              value: {
                message:
                  'Dynatrace Metrics API returned 0 result values, expected 1 for query: https://myUrl.com.\nPlease ensure the response contains exactly one value (e.g., by using :merge(0):avg for the metric). Here is the output for troubleshooting: {"metricId":"builtin:service.response.time:merge(\\"dt.entity.service\\"):percentile(90)","data":[]}',
                metric: 'response_time_p90',
                success: false,
                value: 0,
              },
              warningTargets: [
                {
                  criteria: '<=1200',
                  targetValue: 0,
                  violated: true,
                },
              ],
            },
            {
              displayName: 'Problem open',
              keySli: true,
              passTargets: [
                {
                  criteria: '=0',
                  targetValue: 0,
                  violated: true,
                },
              ],
              score: 0,
              status: 'fail',
              value: {
                message:
                  'Problems API request https://myUrl.com was not successful: Dynatrace API returned error 403: Token is missing required scope. Use one of: problems.read (Read problems)',
                metric: 'problem_open',
                success: false,
                value: 0,
              },
              warningTargets: null,
            },
          ],
          result: 'fail',
          score: 0,
          sloFileContent:
            'LS0tCnNwZWNfdmVyc2lvbjogIjEuMCIKY29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciCiAgY29tcGFyZV93aXRoOiAic2luZ2xlX3Jlc3VsdCIKICBpbmNsdWRlX3Jlc3VsdF93aXRoX3Njb3JlOiAicGFzcyIKICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxCmZpbHRlcjoKb2JqZWN0aXZlczoKICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDkwIgogICAgZGlzcGxheU5hbWU6ICJSZXNwb25zZSB0aW1lIFA5MCIKICAgIGtleV9zbGk6IGZhbHNlCiAgICBwYXNzOiAgICAgICAgICAgICAjIHBhc3MgaWYgKHJlbGF0aXZlIGNoYW5nZSA8PSAxMCUgQU5EIGFic29sdXRlIHZhbHVlIGlzIDwgMTAwMCkKICAgICAgLSBjcml0ZXJpYToKICAgICAgICAgIC0gIjw9KzEwJSIgICMgcmVsYXRpdmUgdmFsdWVzIHJlcXVpcmUgYSBwcmVmaXhlZCBzaWduIChwbHVzIG9yIG1pbnVzKQogICAgICAgICAgLSAiPDEwMDAiICAgIyBhYnNvbHV0ZSB2YWx1ZXMgb25seSByZXF1aXJlIGEgbG9naWNhbCBvcGVyYXRvcgogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyAxMjAwbXMsIHRoZSByZXN1bHQgc2hvdWxkIGJlIGEgd2FybmluZwogICAgICAtIGNyaXRlcmlhOgogICAgICAgICAgLSAiPD0xMjAwIgogICAgd2VpZ2h0OiAxCiAgLSBzbGk6ICJwcm9ibGVtX29wZW4iCiAgICBkaXNwbGF5TmFtZTogIlByb2JsZW0gb3BlbiIKICAgIGtleV9zbGk6IHRydWUKICAgIHBhc3M6CiAgICAgIC0gY3JpdGVyaWE6CiAgICAgICAgICAtICI9MCIKICAgIHdlaWdodDogMQp0b3RhbF9zY29yZToKICBwYXNzOiAiOTAlIgogIHdhcm5pbmc6ICI0MCUiCg==',
          timeEnd: '2021-11-18T13:08:20.106Z',
          timeStart: '2021-11-18T12:53:20.106Z',
        },
        labels: {
          DtCreds: 'dynatrace',
          'Problem URL': 'https://myUrl.com',
        },
        message: 'Evaluation failed since the calculated score of 0 is below the target value of 90',
        problem: {
          ImpactedEntity:
            'Connectivity problem on Process SpringBoot carts works.weave.socks.cart.CartApplication carts-* (carts-655bbd7474-xjsht)',
          PID: '-4365226942138395796_1635983460000V2',
          ProblemDetails: {
            displayName: 'P-21115',
            endTime: -1,
            hasRootCause: false,
            id: '-4365226942138395796_1635983460000V2',
            impactLevel: 'INFRASTRUCTURE',
            severityLevel: 'ERROR',
            startTime: 1635983760000,
            status: 'OPEN',
          },
          ProblemID: 'P-21115',
          ProblemTitle: 'Response time degradation',
          ProblemURL: 'https://myUrl.com',
          State: 'OPEN',
          Tags: 'keptn_service:carts, keptn_project:sockshop, keptn_stage:production, keptn_deployment:direct',
        },
        project: 'sockshop',
        result: '',
        service: 'carts',
        stage: 'production',
        status: '',
        temporaryData: {
          distributor: {
            subscriptionID: '',
          },
        },
      },
      id: 'bd783f59-3325-4d97-ae61-253e93ea259a',
      source: 'shipyard-controller',
      specversion: '1.0',
      time: '2021-11-18T13:08:43.790Z',
      type: 'sh.keptn.event.production.remediation.triggered',
      shkeptncontext: '35383737-3630-4639-b037-353138323631',
      shkeptnspecversion: '0.2.3',
    },
    {
      data: {
        labels: {
          DtCreds: 'dynatrace',
          'Problem URL': 'https://myUrl.com',
        },
        message: 'Evaluation failed since the calculated score of 0 is below the target value of 90',
        project: 'sockshop',
        result: 'fail',
        service: 'carts',
        stage: 'production',
        status: 'succeeded',
      },
      id: '8b382fa8-1cea-4e38-aa5f-9f39484c30c2',
      source: 'shipyard-controller',
      specversion: '1.0',
      time: '2021-11-18T13:08:43.491Z',
      type: 'sh.keptn.event.production.remediation.finished',
      shkeptncontext: '35383737-3630-4639-b037-353138323631',
      shkeptnspecversion: '0.2.3',
      triggeredid: 'd085209d-c2cf-4e43-b797-468b2b79a3ef',
    },
    {
      data: {
        evaluation: {
          comparedEvents: ['e2800526-c0f9-4817-9aec-9a3037d42837'],
          gitCommit: '',
          indicatorResults: [
            {
              displayName: 'Response time P90',
              keySli: false,
              passTargets: [
                {
                  criteria: '<=+10%',
                  targetValue: 0,
                  violated: true,
                },
                {
                  criteria: '<1000',
                  targetValue: 0,
                  violated: true,
                },
              ],
              score: 0,
              status: 'fail',
              value: {
                message:
                  'Dynatrace Metrics API returned 0 result values, expected 1 for query: https://myUrl.com.\nPlease ensure the response contains exactly one value (e.g., by using :merge(0):avg for the metric). Here is the output for troubleshooting: {"metricId":"builtin:service.response.time:merge(\\"dt.entity.service\\"):percentile(90)","data":[]}',
                metric: 'response_time_p90',
                success: false,
                value: 0,
              },
              warningTargets: [
                {
                  criteria: '<=1200',
                  targetValue: 0,
                  violated: true,
                },
              ],
            },
            {
              displayName: 'Problem open',
              keySli: true,
              passTargets: [
                {
                  criteria: '=0',
                  targetValue: 0,
                  violated: true,
                },
              ],
              score: 0,
              status: 'fail',
              value: {
                message:
                  'Problems API request https://myUrl.com was not successful: Dynatrace API returned error 403: Token is missing required scope. Use one of: problems.read (Read problems)',
                metric: 'problem_open',
                success: false,
                value: 0,
              },
              warningTargets: null,
            },
          ],
          result: 'fail',
          score: 0,
          sloFileContent:
            'LS0tCnNwZWNfdmVyc2lvbjogIjEuMCIKY29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciCiAgY29tcGFyZV93aXRoOiAic2luZ2xlX3Jlc3VsdCIKICBpbmNsdWRlX3Jlc3VsdF93aXRoX3Njb3JlOiAicGFzcyIKICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxCmZpbHRlcjoKb2JqZWN0aXZlczoKICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDkwIgogICAgZGlzcGxheU5hbWU6ICJSZXNwb25zZSB0aW1lIFA5MCIKICAgIGtleV9zbGk6IGZhbHNlCiAgICBwYXNzOiAgICAgICAgICAgICAjIHBhc3MgaWYgKHJlbGF0aXZlIGNoYW5nZSA8PSAxMCUgQU5EIGFic29sdXRlIHZhbHVlIGlzIDwgMTAwMCkKICAgICAgLSBjcml0ZXJpYToKICAgICAgICAgIC0gIjw9KzEwJSIgICMgcmVsYXRpdmUgdmFsdWVzIHJlcXVpcmUgYSBwcmVmaXhlZCBzaWduIChwbHVzIG9yIG1pbnVzKQogICAgICAgICAgLSAiPDEwMDAiICAgIyBhYnNvbHV0ZSB2YWx1ZXMgb25seSByZXF1aXJlIGEgbG9naWNhbCBvcGVyYXRvcgogICAgd2FybmluZzogICAgICAgICAgIyBpZiB0aGUgcmVzcG9uc2UgdGltZSBpcyBiZWxvdyAxMjAwbXMsIHRoZSByZXN1bHQgc2hvdWxkIGJlIGEgd2FybmluZwogICAgICAtIGNyaXRlcmlhOgogICAgICAgICAgLSAiPD0xMjAwIgogICAgd2VpZ2h0OiAxCiAgLSBzbGk6ICJwcm9ibGVtX29wZW4iCiAgICBkaXNwbGF5TmFtZTogIlByb2JsZW0gb3BlbiIKICAgIGtleV9zbGk6IHRydWUKICAgIHBhc3M6CiAgICAgIC0gY3JpdGVyaWE6CiAgICAgICAgICAtICI9MCIKICAgIHdlaWdodDogMQp0b3RhbF9zY29yZToKICBwYXNzOiAiOTAlIgogIHdhcm5pbmc6ICI0MCUiCg==',
          timeEnd: '2021-11-18T13:08:20.106Z',
          timeStart: '2021-11-18T12:53:20.106Z',
        },
        labels: {
          DtCreds: 'dynatrace',
          'Problem URL': 'https://myUrl.com',
        },
        message: 'Evaluation failed since the calculated score of 0 is below the target value of 90',
        project: 'sockshop',
        result: 'fail',
        service: 'carts',
        stage: 'production',
        status: 'succeeded',
      },
      id: 'cab7320e-f4e7-4f9a-b26d-4dc79b960297',
      source: 'lighthouse-service',
      specversion: '1.0',
      time: '2021-11-18T13:08:43.306Z',
      type: 'sh.keptn.event.evaluation.finished',
      shkeptncontext: '35383737-3630-4639-b037-353138323631',
      shkeptnspecversion: '0.2.3',
      triggeredid: '1fce9408-7c86-42aa-bbe4-4f8a79362d8f',
    },
    {
      data: {
        'get-sli': {
          end: '2021-11-18T13:08:20.106Z',
          indicatorValues: [
            {
              message:
                'Dynatrace Metrics API returned 0 result values, expected 1 for query: https://myUrl.com.\nPlease ensure the response contains exactly one value (e.g., by using :merge(0):avg for the metric). Here is the output for troubleshooting: {"metricId":"builtin:service.response.time:merge(\\"dt.entity.service\\"):percentile(90)","data":[]}',
              metric: 'response_time_p90',
              success: false,
              value: 0,
            },
            {
              message:
                'Problems API request https://myUrl.com was not successful: Dynatrace API returned error 403: Token is missing required scope. Use one of: problems.read (Read problems)',
              metric: 'problem_open',
              success: false,
              value: 0,
            },
          ],
          start: '2021-11-18T12:53:20.106Z',
        },
        labels: {
          DtCreds: 'dynatrace',
          'Problem URL': 'https://myUrl.com',
        },
        project: 'sockshop',
        result: 'pass',
        service: 'carts',
        stage: 'production',
        status: 'succeeded',
      },
      id: '05fe950b-d88c-49d8-8833-86f5bc00ad23',
      source: 'dynatrace-service',
      specversion: '1.0',
      time: '2021-11-18T13:08:39.796Z',
      type: 'sh.keptn.event.get-sli.finished',
      shkeptncontext: '35383737-3630-4639-b037-353138323631',
      shkeptnspecversion: '0.2.1',
      triggeredid: 'ab1fc561-5fc1-4e5f-94ff-b74ad0e81ae9',
    },
    {
      data: {
        labels: {
          'Problem URL': 'https://myUrl.com',
        },
        project: 'sockshop',
        result: 'pass',
        service: 'carts',
        stage: 'production',
        status: 'succeeded',
      },
      id: 'dfb667f4-ac6a-4c6a-b4ec-f9039db07923',
      source: 'dynatrace-service',
      specversion: '1.0',
      time: '2021-11-18T13:08:21.804Z',
      type: 'sh.keptn.event.get-sli.started',
      shkeptncontext: '35383737-3630-4639-b037-353138323631',
      shkeptnspecversion: '0.2.1',
      triggeredid: 'ab1fc561-5fc1-4e5f-94ff-b74ad0e81ae9',
    },
    {
      data: {
        deployment: '',
        'get-sli': {
          end: '2021-11-18T13:08:20.106Z',
          indicators: ['response_time_p90', 'problem_open'],
          sliProvider: 'dynatrace',
          start: '2021-11-18T12:53:20.106Z',
        },
        labels: {
          'Problem URL': 'https://myUrl.com',
        },
        project: 'sockshop',
        service: 'carts',
        stage: 'production',
      },
      id: 'ab1fc561-5fc1-4e5f-94ff-b74ad0e81ae9',
      source: 'lighthouse-service',
      specversion: '1.0',
      time: '2021-11-18T13:08:21.800Z',
      type: 'sh.keptn.event.get-sli.triggered',
      shkeptncontext: '35383737-3630-4639-b037-353138323631',
      shkeptnspecversion: '0.2.3',
    },
    {
      data: {
        labels: {
          'Problem URL': 'https://myUrl.com',
        },
        project: 'sockshop',
        result: 'pass',
        service: 'carts',
        stage: 'production',
        status: 'succeeded',
      },
      id: '36d4afd3-a9eb-4b31-9709-bbd9a784f601',
      source: 'lighthouse-service',
      specversion: '1.0',
      time: '2021-11-18T13:08:20.106Z',
      type: 'sh.keptn.event.evaluation.started',
      shkeptncontext: '35383737-3630-4639-b037-353138323631',
      shkeptnspecversion: '0.2.3',
      triggeredid: '1fce9408-7c86-42aa-bbe4-4f8a79362d8f',
    },
    {
      data: {
        ActionIndex: 1,
        action: {
          action: 'toggle-feature',
          description: 'Toogle feature flag EnableItemCache to ON',
          name: 'Toogle feature flag',
          value: {
            EnableItemCache: 'on',
          },
        },
        evaluation: {
          timeframe: '15m',
        },
        labels: {
          'Problem URL': 'https://myUrl.com',
        },
        message: '',
        problem: {
          ImpactedEntity:
            'Connectivity problem on Process SpringBoot carts works.weave.socks.cart.CartApplication carts-* (carts-655bbd7474-xjsht)',
          PID: '-4365226942138395796_1635983460000V2',
          ProblemDetails: {
            displayName: 'P-21115',
            endTime: -1,
            hasRootCause: false,
            id: '-4365226942138395796_1635983460000V2',
            impactLevel: 'INFRASTRUCTURE',
            severityLevel: 'ERROR',
            startTime: 1635983760000,
            status: 'OPEN',
          },
          ProblemID: 'P-21115',
          ProblemTitle: 'Response time degradation',
          ProblemURL: 'https://myUrl.com',
          State: 'OPEN',
          Tags: 'keptn_service:carts, keptn_project:sockshop, keptn_stage:production, keptn_deployment:direct',
        },
        project: 'sockshop',
        result: 'pass',
        service: 'carts',
        stage: 'production',
        status: 'succeeded',
        temporaryData: {
          distributor: {
            subscriptionID: '',
          },
        },
      },
      id: '1fce9408-7c86-42aa-bbe4-4f8a79362d8f',
      source: 'shipyard-controller',
      specversion: '1.0',
      time: '2021-11-18T13:08:20.098Z',
      type: 'sh.keptn.event.evaluation.triggered',
      shkeptncontext: '35383737-3630-4639-b037-353138323631',
    },
    {
      data: {
        action: {},
        labels: {
          'Problem URL': 'https://myUrl.com',
        },
        project: 'sockshop',
        result: 'pass',
        service: 'carts',
        stage: 'production',
        status: 'succeeded',
      },
      id: 'd46e6bca-4f8e-4627-8633-ad256b259c1d',
      source: 'unleash-service',
      specversion: '1.0',
      time: '2021-11-18T12:53:14.333Z',
      type: 'sh.keptn.event.action.finished',
      shkeptncontext: '35383737-3630-4639-b037-353138323631',
      triggeredid: '7f2c8734-dd49-4bdf-8ac2-4eff6c84bc71',
    },
    {
      data: {
        labels: {
          'Problem URL': 'https://myUrl.com',
        },
        project: 'sockshop',
        service: 'carts',
        stage: 'production',
      },
      id: '85a633d5-c574-4b6a-814c-218d3e585795',
      source: 'unleash-service',
      specversion: '1.0',
      time: '2021-11-18T12:53:14.311Z',
      type: 'sh.keptn.event.action.started',
      shkeptncontext: '35383737-3630-4639-b037-353138323631',
      triggeredid: '7f2c8734-dd49-4bdf-8ac2-4eff6c84bc71',
    },
    {
      data: {
        ActionIndex: 1,
        action: {
          action: 'toggle-feature',
          description: 'Toogle feature flag EnableItemCache to ON',
          name: 'Toogle feature flag',
          value: {
            EnableItemCache: 'on',
          },
        },
        labels: {
          'Problem URL': 'https://myUrl.com',
        },
        message: '',
        problem: {
          ImpactedEntity:
            'Connectivity problem on Process SpringBoot carts works.weave.socks.cart.CartApplication carts-* (carts-655bbd7474-xjsht)',
          PID: '-4365226942138395796_1635983460000V2',
          ProblemDetails: {
            displayName: 'P-21115',
            endTime: -1,
            hasRootCause: false,
            id: '-4365226942138395796_1635983460000V2',
            impactLevel: 'INFRASTRUCTURE',
            severityLevel: 'ERROR',
            startTime: 1635983760000,
            status: 'OPEN',
          },
          ProblemID: 'P-21115',
          ProblemTitle: 'Response time degradation',
          ProblemURL: 'https://myUrl.com',
          State: 'OPEN',
          Tags: 'keptn_service:carts, keptn_project:sockshop, keptn_stage:production, keptn_deployment:direct',
        },
        project: 'sockshop',
        result: 'pass',
        service: 'carts',
        stage: 'production',
        status: 'succeeded',
        temporaryData: {
          distributor: {
            subscriptionID: '',
          },
        },
      },
      id: '7f2c8734-dd49-4bdf-8ac2-4eff6c84bc71',
      source: 'shipyard-controller',
      specversion: '1.0',
      time: '2021-11-18T12:53:14.294Z',
      type: 'sh.keptn.event.action.triggered',
      shkeptncontext: '35383737-3630-4639-b037-353138323631',
      shkeptnspecversion: '0.2.3',
    },
    {
      data: {
        ActionIndex: 1,
        action: {
          action: 'toggle-feature',
          description: 'Toogle feature flag EnableItemCache to ON',
          name: 'Toogle feature flag',
          value: {
            EnableItemCache: 'on',
          },
        },
        labels: {
          'Problem URL': 'https://myUrl.com',
        },
        project: 'sockshop',
        result: 'pass',
        service: 'carts',
        stage: 'production',
        status: 'succeeded',
      },
      id: '42eb1bfc-3483-492b-a2e0-272a2a9ea6d4',
      source: 'remediation-service',
      specversion: '1.0',
      time: '2021-11-18T12:53:14.186Z',
      type: 'sh.keptn.event.get-action.finished',
      shkeptncontext: '35383737-3630-4639-b037-353138323631',
      triggeredid: '6d1bcc3c-2b21-4f33-8a8c-d5bf34ac3791',
    },
    {
      data: {
        labels: {
          'Problem URL': 'https://myUrl.com',
        },
        project: 'sockshop',
        service: 'carts',
        stage: 'production',
      },
      id: 'aaf9eb3f-9ead-4655-b997-d6d9b2351f46',
      source: 'remediation-service',
      specversion: '1.0',
      time: '2021-11-18T12:53:12.292Z',
      type: 'sh.keptn.event.get-action.started',
      shkeptncontext: '35383737-3630-4639-b037-353138323631',
      triggeredid: '6d1bcc3c-2b21-4f33-8a8c-d5bf34ac3791',
    },
    {
      data: {
        'get-action': null,
        labels: {
          'Problem URL': 'https://myUrl.com',
        },
        message: '',
        problem: {
          ImpactedEntity:
            'Connectivity problem on Process SpringBoot carts works.weave.socks.cart.CartApplication carts-* (carts-655bbd7474-xjsht)',
          PID: '-4365226942138395796_1635983460000V2',
          ProblemDetails: {
            displayName: 'P-21115',
            endTime: -1,
            hasRootCause: false,
            id: '-4365226942138395796_1635983460000V2',
            impactLevel: 'INFRASTRUCTURE',
            severityLevel: 'ERROR',
            startTime: 1635983760000,
            status: 'OPEN',
          },
          ProblemID: 'P-21115',
          ProblemTitle: 'Response time degradation',
          ProblemURL: 'https://myUrl.com',
          State: 'OPEN',
          Tags: 'keptn_service:carts, keptn_project:sockshop, keptn_stage:production, keptn_deployment:direct',
        },
        project: 'sockshop',
        result: '',
        service: 'carts',
        stage: 'production',
        status: '',
        temporaryData: {
          distributor: {
            subscriptionID: '',
          },
        },
      },
      id: '6d1bcc3c-2b21-4f33-8a8c-d5bf34ac3791',
      source: 'shipyard-controller',
      specversion: '1.0',
      time: '2021-11-18T12:53:12.285Z',
      type: 'sh.keptn.event.get-action.triggered',
      shkeptncontext: '35383737-3630-4639-b037-353138323631',
      shkeptnspecversion: '0.2.3',
    },
    {
      data: {
        labels: {
          'Problem URL': 'https://myUrl.com',
        },
        problem: {
          ImpactedEntity:
            'Connectivity problem on Process SpringBoot carts works.weave.socks.cart.CartApplication carts-* (carts-655bbd7474-xjsht)',
          PID: '-4365226942138395796_1635983460000V2',
          ProblemDetails: {
            displayName: 'P-21115',
            endTime: -1,
            hasRootCause: false,
            id: '-4365226942138395796_1635983460000V2',
            impactLevel: 'INFRASTRUCTURE',
            severityLevel: 'ERROR',
            startTime: 1635983760000,
            status: 'OPEN',
          },
          ProblemID: 'P-21115',
          ProblemTitle: 'Response time degradation',
          ProblemURL: 'https://myUrl.com',
          State: 'OPEN',
          Tags: 'keptn_service:carts, keptn_project:sockshop, keptn_stage:production, keptn_deployment:direct',
        },
        project: 'sockshop',
        service: 'carts',
        stage: 'production',
        temporaryData: {
          distributor: {
            subscriptionID: '',
          },
        },
      },
      id: 'd085209d-c2cf-4e43-b797-468b2b79a3ef',
      source: 'swagger-ui',
      specversion: '1.0',
      time: '2021-11-18T12:53:10.015Z',
      type: 'sh.keptn.event.production.remediation.triggered',
      shkeptncontext: '35383737-3630-4639-b037-353138323631',
      shkeptnspecversion: '0.2.3',
    },
  ],
  pageSize: 100,
  totalCount: 19,
};

const serviceRemediationResponse = {
  stages: [
    {
      name: 'production',
      remediations: [
        {
          name: 'remediation',
          problemTitle: 'Failure rate increase',
          service: 'carts',
          project: 'sockshop',
          time: '2021-11-04T04:51:21.557Z',
          shkeptncontext: '35383737-3630-4639-b037-353138323631',
          state: 'started',
          stages: [
            {
              name: 'production',
              actions: [],
              state: 'started',
              latestEvent: {
                type: 'sh.keptn.event.get-action.started',
                id: '7448420f-5b15-4777-9d39-cc8308e2b0c3',
                time: '2021-11-04T04:51:23.266Z',
              },
            },
          ],
        },
      ],
    },
  ],
};

export { openRemediationsResponseMock as OpenRemediationsResponse };
export { remediationTraceResponse as RemediationTraceResponse };
export { remediationTracesResponse as RemediationTracesResponse };
export { serviceRemediationResponse as ServiceRemediationResponse };
