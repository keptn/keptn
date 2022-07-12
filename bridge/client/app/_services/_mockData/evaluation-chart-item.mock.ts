import { TestUtils } from '../../_utils/test.utils';

const evaluationChartItemMock = [
  {
    metricName: 'Score',
    name: 'Score',
    type: 'column',
    data: [
      {
        y: 100,
        evaluationData: {
          traces: [],
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
                    value: 298.9114676593789,
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
                'LS0tCnNwZWNfdmVyc2lvbjogIjEuMCIKY29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciCiAgY29tcGFyZV93aXRoOiAic2luZ2xlX3Jlc3VsdCIKICBpbmNsdWRlX3Jlc3VsdF93aXRoX3Njb3JlOiAicGFzcyIKICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxCmZpbHRlcjoKb2JqZWN0aXZlczoKICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1IgogICAgZGlzcGxheU5hbWU6ICJSZXNwb25zZSB0aW1lIFA5NSIKICAgIGtleV9zbGk6IGZhbHNlCiAgICBwYXNzOiAgICAgICAgICAgICAjIHBhc3MgaWYgKHJlbGF0aXZlIGNoYW5nZSA8PSAxMCUgQU5EIGFic29sdXRlIHZhbHVlIGlzIDwgNjAwbXMpCiAgICAgIC0gY3JpdGVyaWE6CiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykKICAgICAgICAgIC0gIjw2MDAiICAgICMgYWJzb2x1dGUgdmFsdWVzIG9ubHkgcmVxdWlyZSBhIGxvZ2ljYWwgb3BlcmF0b3IKICAgIHdhcm5pbmc6ICAgICAgICAgICMgaWYgdGhlIHJlc3BvbnNlIHRpbWUgaXMgYmVsb3cgODAwbXMsIHRoZSByZXN1bHQgc2hvdWxkIGJlIGEgd2FybmluZwogICAgICAtIGNyaXRlcmlhOgogICAgICAgICAgLSAiPD04MDAiCiAgICB3ZWlnaHQ6IDEKdG90YWxfc2NvcmU6CiAgcGFzczogIjkwJSIKICB3YXJuaW5nOiAiNzUlIg==',
              timeEnd: '2021-09-30T11:41:15Z',
              timeStart: '2021-09-30T11:39:08Z',
              sloFileContentParsed:
                '---\nspec_version: "1.0"\ncomparison:\n  aggregate_function: "avg"\n  compare_with: "single_result"\n  include_result_with_score: "pass"\n  number_of_comparison_results: 1\nfilter:\nobjectives:\n  - sli: "response_time_p95"\n    displayName: "Response time P95"\n    key_sli: false\n    pass:             # pass if (relative change <= 10% AND absolute value is < 600ms)\n      - criteria:\n          - "<=+10%"  # relative values require a prefixed sign (plus or minus)\n          - "<600"    # absolute values only require a logical operator\n    warning:          # if the response time is below 800ms, the result should be a warning\n      - criteria:\n          - "<=800"\n    weight: 1\ntotal_score:\n  pass: "90%"\n  warning: "75%"',
              score_pass: '90',
              score_warning: '75',
              compare_with: 'single_result\n',
              include_result_with_score: 'pass\n',
              number_of_comparison_results: 0,
            },
            labels: {
              DtCreds: 'dynatrace',
            },
            project: 'sockshop',
            result: 'pass',
            service: 'carts',
            stage: 'staging',
            status: 'succeeded',
            temporaryData: {
              distributor: {
                subscriptionID: '',
              },
            },
          },
          id: '40cfcfc2-9cc6-4fce-ba62-5f5fec3e3e3b',
          source: 'lighthouse-service',
          specversion: '1.0',
          time: '2021-09-30T11:42:29.616Z',
          type: 'sh.keptn.event.evaluation.finished',
          shkeptncontext: '369b15c8-5410-49de-87f7-08a2f9f16243',
          shkeptnspecversion: '0.2.3',
          triggeredid: '2d2e799f-6c25-4b26-83ab-12953335f6d4',
          plainEvent: {
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
                      value: 298.9114676593789,
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
                  'LS0tCnNwZWNfdmVyc2lvbjogIjEuMCIKY29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciCiAgY29tcGFyZV93aXRoOiAic2luZ2xlX3Jlc3VsdCIKICBpbmNsdWRlX3Jlc3VsdF93aXRoX3Njb3JlOiAicGFzcyIKICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxCmZpbHRlcjoKb2JqZWN0aXZlczoKICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1IgogICAgZGlzcGxheU5hbWU6ICJSZXNwb25zZSB0aW1lIFA5NSIKICAgIGtleV9zbGk6IGZhbHNlCiAgICBwYXNzOiAgICAgICAgICAgICAjIHBhc3MgaWYgKHJlbGF0aXZlIGNoYW5nZSA8PSAxMCUgQU5EIGFic29sdXRlIHZhbHVlIGlzIDwgNjAwbXMpCiAgICAgIC0gY3JpdGVyaWE6CiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykKICAgICAgICAgIC0gIjw2MDAiICAgICMgYWJzb2x1dGUgdmFsdWVzIG9ubHkgcmVxdWlyZSBhIGxvZ2ljYWwgb3BlcmF0b3IKICAgIHdhcm5pbmc6ICAgICAgICAgICMgaWYgdGhlIHJlc3BvbnNlIHRpbWUgaXMgYmVsb3cgODAwbXMsIHRoZSByZXN1bHQgc2hvdWxkIGJlIGEgd2FybmluZwogICAgICAtIGNyaXRlcmlhOgogICAgICAgICAgLSAiPD04MDAiCiAgICB3ZWlnaHQ6IDEKdG90YWxfc2NvcmU6CiAgcGFzczogIjkwJSIKICB3YXJuaW5nOiAiNzUlIg==',
                timeEnd: '2021-09-30T11:41:15Z',
                timeStart: '2021-09-30T11:39:08Z',
              },
              labels: {
                DtCreds: 'dynatrace',
              },
              project: 'sockshop',
              result: 'pass',
              service: 'carts',
              stage: 'staging',
              status: 'succeeded',
              temporaryData: {
                distributor: {
                  subscriptionID: '',
                },
              },
            },
            id: '40cfcfc2-9cc6-4fce-ba62-5f5fec3e3e3b',
            source: 'lighthouse-service',
            specversion: '1.0',
            time: '2021-09-30T11:42:29.616Z',
            type: 'sh.keptn.event.evaluation.finished',
            shkeptncontext: '369b15c8-5410-49de-87f7-08a2f9f16243',
            shkeptnspecversion: '0.2.3',
            triggeredid: '2d2e799f-6c25-4b26-83ab-12953335f6d4',
          },
          heatmapLabel: '2021-09-30 13:42',
        },
        color: '#7dc540',
        name: '2021-09-30 13:42',
        label: '2021-09-30 13:42',
        x: 0,
      },
      {
        y: 0,
        evaluationData: {
          traces: [],
          data: {
            evaluation: {
              comparedEvents: ['40cfcfc2-9cc6-4fce-ba62-5f5fec3e3e3b'],
              indicatorResults: [
                {
                  displayName: 'Response time P95',
                  keySli: false,
                  passTargets: [
                    {
                      criteria: '<=+10%',
                      targetValue: 328.8026144253168,
                      violated: true,
                    },
                    {
                      criteria: '<600',
                      targetValue: 600,
                      violated: true,
                    },
                  ],
                  score: 0,
                  status: 'fail',
                  value: {
                    metric: 'response_time_p95',
                    success: true,
                    value: 1091.8809833281025,
                  },
                  warningTargets: [
                    {
                      criteria: '<=800',
                      targetValue: 800,
                      violated: true,
                    },
                  ],
                },
              ],
              result: 'fail',
              score: 0,
              sloFileContent:
                'LS0tCnNwZWNfdmVyc2lvbjogIjEuMCIKY29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciCiAgY29tcGFyZV93aXRoOiAic2luZ2xlX3Jlc3VsdCIKICBpbmNsdWRlX3Jlc3VsdF93aXRoX3Njb3JlOiAicGFzcyIKICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxCmZpbHRlcjoKb2JqZWN0aXZlczoKICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1IgogICAgZGlzcGxheU5hbWU6ICJSZXNwb25zZSB0aW1lIFA5NSIKICAgIGtleV9zbGk6IGZhbHNlCiAgICBwYXNzOiAgICAgICAgICAgICAjIHBhc3MgaWYgKHJlbGF0aXZlIGNoYW5nZSA8PSAxMCUgQU5EIGFic29sdXRlIHZhbHVlIGlzIDwgNjAwbXMpCiAgICAgIC0gY3JpdGVyaWE6CiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykKICAgICAgICAgIC0gIjw2MDAiICAgICMgYWJzb2x1dGUgdmFsdWVzIG9ubHkgcmVxdWlyZSBhIGxvZ2ljYWwgb3BlcmF0b3IKICAgIHdhcm5pbmc6ICAgICAgICAgICMgaWYgdGhlIHJlc3BvbnNlIHRpbWUgaXMgYmVsb3cgODAwbXMsIHRoZSByZXN1bHQgc2hvdWxkIGJlIGEgd2FybmluZwogICAgICAtIGNyaXRlcmlhOgogICAgICAgICAgLSAiPD04MDAiCiAgICB3ZWlnaHQ6IDEKdG90YWxfc2NvcmU6CiAgcGFzczogIjkwJSIKICB3YXJuaW5nOiAiNzUlIg==',
              timeEnd: '2021-10-01T09:25:48Z',
              timeStart: '2021-10-01T09:16:31Z',
              sloFileContentParsed:
                '---\nspec_version: "1.0"\ncomparison:\n  aggregate_function: "avg"\n  compare_with: "single_result"\n  include_result_with_score: "pass"\n  number_of_comparison_results: 1\nfilter:\nobjectives:\n  - sli: "response_time_p95"\n    displayName: "Response time P95"\n    key_sli: false\n    pass:             # pass if (relative change <= 10% AND absolute value is < 600ms)\n      - criteria:\n          - "<=+10%"  # relative values require a prefixed sign (plus or minus)\n          - "<600"    # absolute values only require a logical operator\n    warning:          # if the response time is below 800ms, the result should be a warning\n      - criteria:\n          - "<=800"\n    weight: 1\ntotal_score:\n  pass: "90%"\n  warning: "75%"',
              score_pass: '90',
              score_warning: '75',
              compare_with: 'single_result\n',
              include_result_with_score: 'pass\n',
              number_of_comparison_results: 1,
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
            temporaryData: {
              distributor: {
                subscriptionID: '',
              },
            },
          },
          id: 'f623de99-a716-4b9b-90a2-8df1054e92c7',
          source: 'lighthouse-service',
          specversion: '1.0',
          time: '2021-10-01T09:26:15.608Z',
          type: 'sh.keptn.event.evaluation.finished',
          shkeptncontext: '081ce5a3-e0f9-4c4d-bfcf-5893f85bfaff',
          shkeptnspecversion: '0.2.3',
          triggeredid: 'f564369b-2ee6-4c32-8488-03508ed05b5a',
          plainEvent: {
            data: {
              evaluation: {
                comparedEvents: ['40cfcfc2-9cc6-4fce-ba62-5f5fec3e3e3b'],
                indicatorResults: [
                  {
                    displayName: 'Response time P95',
                    keySli: false,
                    passTargets: [
                      {
                        criteria: '<=+10%',
                        targetValue: 328.8026144253168,
                        violated: true,
                      },
                      {
                        criteria: '<600',
                        targetValue: 600,
                        violated: true,
                      },
                    ],
                    score: 0,
                    status: 'fail',
                    value: {
                      metric: 'response_time_p95',
                      success: true,
                      value: 1091.8809833281025,
                    },
                    warningTargets: [
                      {
                        criteria: '<=800',
                        targetValue: 800,
                        violated: true,
                      },
                    ],
                  },
                ],
                result: 'fail',
                score: 0,
                sloFileContent:
                  'LS0tCnNwZWNfdmVyc2lvbjogIjEuMCIKY29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciCiAgY29tcGFyZV93aXRoOiAic2luZ2xlX3Jlc3VsdCIKICBpbmNsdWRlX3Jlc3VsdF93aXRoX3Njb3JlOiAicGFzcyIKICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxCmZpbHRlcjoKb2JqZWN0aXZlczoKICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1IgogICAgZGlzcGxheU5hbWU6ICJSZXNwb25zZSB0aW1lIFA5NSIKICAgIGtleV9zbGk6IGZhbHNlCiAgICBwYXNzOiAgICAgICAgICAgICAjIHBhc3MgaWYgKHJlbGF0aXZlIGNoYW5nZSA8PSAxMCUgQU5EIGFic29sdXRlIHZhbHVlIGlzIDwgNjAwbXMpCiAgICAgIC0gY3JpdGVyaWE6CiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykKICAgICAgICAgIC0gIjw2MDAiICAgICMgYWJzb2x1dGUgdmFsdWVzIG9ubHkgcmVxdWlyZSBhIGxvZ2ljYWwgb3BlcmF0b3IKICAgIHdhcm5pbmc6ICAgICAgICAgICMgaWYgdGhlIHJlc3BvbnNlIHRpbWUgaXMgYmVsb3cgODAwbXMsIHRoZSByZXN1bHQgc2hvdWxkIGJlIGEgd2FybmluZwogICAgICAtIGNyaXRlcmlhOgogICAgICAgICAgLSAiPD04MDAiCiAgICB3ZWlnaHQ6IDEKdG90YWxfc2NvcmU6CiAgcGFzczogIjkwJSIKICB3YXJuaW5nOiAiNzUlIg==',
                timeEnd: '2021-10-01T09:25:48Z',
                timeStart: '2021-10-01T09:16:31Z',
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
              temporaryData: {
                distributor: {
                  subscriptionID: '',
                },
              },
            },
            id: 'f623de99-a716-4b9b-90a2-8df1054e92c7',
            source: 'lighthouse-service',
            specversion: '1.0',
            time: '2021-10-01T09:26:15.608Z',
            type: 'sh.keptn.event.evaluation.finished',
            shkeptncontext: '081ce5a3-e0f9-4c4d-bfcf-5893f85bfaff',
            shkeptnspecversion: '0.2.3',
            triggeredid: 'f564369b-2ee6-4c32-8488-03508ed05b5a',
          },
          heatmapLabel: '2021-10-01 11:26',
        },
        color: '#dc172a',
        name: '2021-10-01 11:26',
        label: '2021-10-01 11:26',
        x: 1,
      },
      {
        y: 100,
        evaluationData: {
          traces: [],
          data: {
            evaluation: {
              comparedEvents: ['40cfcfc2-9cc6-4fce-ba62-5f5fec3e3e3b'],
              indicatorResults: [
                {
                  displayName: 'Response time P95',
                  keySli: false,
                  passTargets: [
                    {
                      criteria: '<=+10%',
                      targetValue: 328.8026144253168,
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
                    value: 306.3587409268154,
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
                'LS0tCnNwZWNfdmVyc2lvbjogIjEuMCIKY29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciCiAgY29tcGFyZV93aXRoOiAic2luZ2xlX3Jlc3VsdCIKICBpbmNsdWRlX3Jlc3VsdF93aXRoX3Njb3JlOiAicGFzcyIKICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxCmZpbHRlcjoKb2JqZWN0aXZlczoKICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1IgogICAgZGlzcGxheU5hbWU6ICJSZXNwb25zZSB0aW1lIFA5NSIKICAgIGtleV9zbGk6IGZhbHNlCiAgICBwYXNzOiAgICAgICAgICAgICAjIHBhc3MgaWYgKHJlbGF0aXZlIGNoYW5nZSA8PSAxMCUgQU5EIGFic29sdXRlIHZhbHVlIGlzIDwgNjAwbXMpCiAgICAgIC0gY3JpdGVyaWE6CiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykKICAgICAgICAgIC0gIjw2MDAiICAgICMgYWJzb2x1dGUgdmFsdWVzIG9ubHkgcmVxdWlyZSBhIGxvZ2ljYWwgb3BlcmF0b3IKICAgIHdhcm5pbmc6ICAgICAgICAgICMgaWYgdGhlIHJlc3BvbnNlIHRpbWUgaXMgYmVsb3cgODAwbXMsIHRoZSByZXN1bHQgc2hvdWxkIGJlIGEgd2FybmluZwogICAgICAtIGNyaXRlcmlhOgogICAgICAgICAgLSAiPD04MDAiCiAgICB3ZWlnaHQ6IDEKdG90YWxfc2NvcmU6CiAgcGFzczogIjkwJSIKICB3YXJuaW5nOiAiNzUlIg==',
              timeEnd: '2021-10-01T10:26:34Z',
              timeStart: '2021-10-01T10:24:15Z',
              sloFileContentParsed:
                '---\nspec_version: "1.0"\ncomparison:\n  aggregate_function: "avg"\n  compare_with: "single_result"\n  include_result_with_score: "pass"\n  number_of_comparison_results: 1\nfilter:\nobjectives:\n  - sli: "response_time_p95"\n    displayName: "Response time P95"\n    key_sli: false\n    pass:             # pass if (relative change <= 10% AND absolute value is < 600ms)\n      - criteria:\n          - "<=+10%"  # relative values require a prefixed sign (plus or minus)\n          - "<600"    # absolute values only require a logical operator\n    warning:          # if the response time is below 800ms, the result should be a warning\n      - criteria:\n          - "<=800"\n    weight: 1\ntotal_score:\n  pass: "90%"\n  warning: "75%"',
              score_pass: '90',
              score_warning: '75',
              compare_with: 'single_result\n',
              include_result_with_score: 'pass\n',
              number_of_comparison_results: 1,
            },
            labels: {
              DtCreds: 'dynatrace',
            },
            project: 'sockshop',
            result: 'pass',
            service: 'carts',
            stage: 'staging',
            status: 'succeeded',
            temporaryData: {
              distributor: {
                subscriptionID: '',
              },
            },
          },
          id: '3344487d-e384-4cd9-a0e0-fcf157a33ad6',
          source: 'lighthouse-service',
          specversion: '1.0',
          time: '2021-10-01T10:27:49.318Z',
          type: 'sh.keptn.event.evaluation.finished',
          shkeptncontext: 'e9dd1b2f-8383-48a9-a2d4-36cb8ff370c7',
          shkeptnspecversion: '0.2.3',
          triggeredid: '61de0ee6-b535-42bb-ac30-d0f359298bb0',
          plainEvent: {
            data: {
              evaluation: {
                comparedEvents: ['40cfcfc2-9cc6-4fce-ba62-5f5fec3e3e3b'],
                indicatorResults: [
                  {
                    displayName: 'Response time P95',
                    keySli: false,
                    passTargets: [
                      {
                        criteria: '<=+10%',
                        targetValue: 328.8026144253168,
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
                      value: 306.3587409268154,
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
                  'LS0tCnNwZWNfdmVyc2lvbjogIjEuMCIKY29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciCiAgY29tcGFyZV93aXRoOiAic2luZ2xlX3Jlc3VsdCIKICBpbmNsdWRlX3Jlc3VsdF93aXRoX3Njb3JlOiAicGFzcyIKICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxCmZpbHRlcjoKb2JqZWN0aXZlczoKICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1IgogICAgZGlzcGxheU5hbWU6ICJSZXNwb25zZSB0aW1lIFA5NSIKICAgIGtleV9zbGk6IGZhbHNlCiAgICBwYXNzOiAgICAgICAgICAgICAjIHBhc3MgaWYgKHJlbGF0aXZlIGNoYW5nZSA8PSAxMCUgQU5EIGFic29sdXRlIHZhbHVlIGlzIDwgNjAwbXMpCiAgICAgIC0gY3JpdGVyaWE6CiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykKICAgICAgICAgIC0gIjw2MDAiICAgICMgYWJzb2x1dGUgdmFsdWVzIG9ubHkgcmVxdWlyZSBhIGxvZ2ljYWwgb3BlcmF0b3IKICAgIHdhcm5pbmc6ICAgICAgICAgICMgaWYgdGhlIHJlc3BvbnNlIHRpbWUgaXMgYmVsb3cgODAwbXMsIHRoZSByZXN1bHQgc2hvdWxkIGJlIGEgd2FybmluZwogICAgICAtIGNyaXRlcmlhOgogICAgICAgICAgLSAiPD04MDAiCiAgICB3ZWlnaHQ6IDEKdG90YWxfc2NvcmU6CiAgcGFzczogIjkwJSIKICB3YXJuaW5nOiAiNzUlIg==',
                timeEnd: '2021-10-01T10:26:34Z',
                timeStart: '2021-10-01T10:24:15Z',
              },
              labels: {
                DtCreds: 'dynatrace',
              },
              project: 'sockshop',
              result: 'pass',
              service: 'carts',
              stage: 'staging',
              status: 'succeeded',
              temporaryData: {
                distributor: {
                  subscriptionID: '',
                },
              },
            },
            id: '3344487d-e384-4cd9-a0e0-fcf157a33ad6',
            source: 'lighthouse-service',
            specversion: '1.0',
            time: '2021-10-01T10:27:49.318Z',
            type: 'sh.keptn.event.evaluation.finished',
            shkeptncontext: 'e9dd1b2f-8383-48a9-a2d4-36cb8ff370c7',
            shkeptnspecversion: '0.2.3',
            triggeredid: '61de0ee6-b535-42bb-ac30-d0f359298bb0',
          },
          heatmapLabel: '2021-10-01 12:27',
        },
        color: '#7dc540',
        name: '2021-10-01 12:27',
        label: '2021-10-01 12:27',
        x: 2,
      },
      {
        y: 100,
        evaluationData: {
          traces: [],
          data: {
            evaluation: {
              comparedEvents: ['3344487d-e384-4cd9-a0e0-fcf157a33ad6'],
              indicatorResults: [
                {
                  displayName: 'Response time P95',
                  keySli: false,
                  passTargets: [
                    {
                      criteria: '<=+10%',
                      targetValue: 336.9946150194969,
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
                    value: 23.12433526069513,
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
                'LS0tCnNwZWNfdmVyc2lvbjogIjEuMCIKY29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciCiAgY29tcGFyZV93aXRoOiAic2luZ2xlX3Jlc3VsdCIKICBpbmNsdWRlX3Jlc3VsdF93aXRoX3Njb3JlOiAicGFzcyIKICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxCmZpbHRlcjoKb2JqZWN0aXZlczoKICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1IgogICAgZGlzcGxheU5hbWU6ICJSZXNwb25zZSB0aW1lIFA5NSIKICAgIGtleV9zbGk6IGZhbHNlCiAgICBwYXNzOiAgICAgICAgICAgICAjIHBhc3MgaWYgKHJlbGF0aXZlIGNoYW5nZSA8PSAxMCUgQU5EIGFic29sdXRlIHZhbHVlIGlzIDwgNjAwbXMpCiAgICAgIC0gY3JpdGVyaWE6CiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykKICAgICAgICAgIC0gIjw2MDAiICAgICMgYWJzb2x1dGUgdmFsdWVzIG9ubHkgcmVxdWlyZSBhIGxvZ2ljYWwgb3BlcmF0b3IKICAgIHdhcm5pbmc6ICAgICAgICAgICMgaWYgdGhlIHJlc3BvbnNlIHRpbWUgaXMgYmVsb3cgODAwbXMsIHRoZSByZXN1bHQgc2hvdWxkIGJlIGEgd2FybmluZwogICAgICAtIGNyaXRlcmlhOgogICAgICAgICAgLSAiPD04MDAiCiAgICB3ZWlnaHQ6IDEKdG90YWxfc2NvcmU6CiAgcGFzczogIjkwJSIKICB3YXJuaW5nOiAiNzUlIg==',
              timeEnd: '2021-10-01T10:32:47.372Z',
              timeStart: '2021-10-01T10:27:47.372Z',
              sloFileContentParsed:
                '---\nspec_version: "1.0"\ncomparison:\n  aggregate_function: "avg"\n  compare_with: "single_result"\n  include_result_with_score: "pass"\n  number_of_comparison_results: 1\nfilter:\nobjectives:\n  - sli: "response_time_p95"\n    displayName: "Response time P95"\n    key_sli: false\n    pass:             # pass if (relative change <= 10% AND absolute value is < 600ms)\n      - criteria:\n          - "<=+10%"  # relative values require a prefixed sign (plus or minus)\n          - "<600"    # absolute values only require a logical operator\n    warning:          # if the response time is below 800ms, the result should be a warning\n      - criteria:\n          - "<=800"\n    weight: 1\ntotal_score:\n  pass: "90%"\n  warning: "75%"',
              score_pass: '90',
              score_warning: '75',
              compare_with: 'single_result\n',
              include_result_with_score: 'pass\n',
              number_of_comparison_results: 1,
            },
            labels: {
              DtCreds: 'dynatrace',
            },
            project: 'sockshop',
            result: 'pass',
            service: 'carts',
            stage: 'staging',
            status: 'succeeded',
            temporaryData: {
              distributor: {
                subscriptionID: '',
              },
            },
          },
          id: '27c733aa-9bb0-44b4-9574-ef04e38eb4c4',
          source: 'lighthouse-service',
          specversion: '1.0',
          time: '2021-10-01T10:33:19.708Z',
          type: 'sh.keptn.event.evaluation.finished',
          shkeptncontext: '199287a3-b723-479b-8083-8b3ff3482531',
          shkeptnspecversion: '0.2.3',
          triggeredid: '53ac99d5-2828-470a-b067-442562cf3bca',
          plainEvent: {
            data: {
              evaluation: {
                comparedEvents: ['3344487d-e384-4cd9-a0e0-fcf157a33ad6'],
                indicatorResults: [
                  {
                    displayName: 'Response time P95',
                    keySli: false,
                    passTargets: [
                      {
                        criteria: '<=+10%',
                        targetValue: 336.9946150194969,
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
                      value: 23.12433526069513,
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
                  'LS0tCnNwZWNfdmVyc2lvbjogIjEuMCIKY29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciCiAgY29tcGFyZV93aXRoOiAic2luZ2xlX3Jlc3VsdCIKICBpbmNsdWRlX3Jlc3VsdF93aXRoX3Njb3JlOiAicGFzcyIKICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxCmZpbHRlcjoKb2JqZWN0aXZlczoKICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1IgogICAgZGlzcGxheU5hbWU6ICJSZXNwb25zZSB0aW1lIFA5NSIKICAgIGtleV9zbGk6IGZhbHNlCiAgICBwYXNzOiAgICAgICAgICAgICAjIHBhc3MgaWYgKHJlbGF0aXZlIGNoYW5nZSA8PSAxMCUgQU5EIGFic29sdXRlIHZhbHVlIGlzIDwgNjAwbXMpCiAgICAgIC0gY3JpdGVyaWE6CiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykKICAgICAgICAgIC0gIjw2MDAiICAgICMgYWJzb2x1dGUgdmFsdWVzIG9ubHkgcmVxdWlyZSBhIGxvZ2ljYWwgb3BlcmF0b3IKICAgIHdhcm5pbmc6ICAgICAgICAgICMgaWYgdGhlIHJlc3BvbnNlIHRpbWUgaXMgYmVsb3cgODAwbXMsIHRoZSByZXN1bHQgc2hvdWxkIGJlIGEgd2FybmluZwogICAgICAtIGNyaXRlcmlhOgogICAgICAgICAgLSAiPD04MDAiCiAgICB3ZWlnaHQ6IDEKdG90YWxfc2NvcmU6CiAgcGFzczogIjkwJSIKICB3YXJuaW5nOiAiNzUlIg==',
                timeEnd: '2021-10-01T10:32:47.372Z',
                timeStart: '2021-10-01T10:27:47.372Z',
              },
              labels: {
                DtCreds: 'dynatrace',
              },
              project: 'sockshop',
              result: 'pass',
              service: 'carts',
              stage: 'staging',
              status: 'succeeded',
              temporaryData: {
                distributor: {
                  subscriptionID: '',
                },
              },
            },
            id: '27c733aa-9bb0-44b4-9574-ef04e38eb4c4',
            source: 'lighthouse-service',
            specversion: '1.0',
            time: '2021-10-01T10:33:19.708Z',
            type: 'sh.keptn.event.evaluation.finished',
            shkeptncontext: '199287a3-b723-479b-8083-8b3ff3482531',
            shkeptnspecversion: '0.2.3',
            triggeredid: '53ac99d5-2828-470a-b067-442562cf3bca',
          },
          heatmapLabel: '2021-10-01 12:33',
        },
        color: '#7dc540',
        name: '2021-10-01 12:33',
        label: '2021-10-01 12:33',
        x: 3,
      },
      {
        y: 50,
        evaluationData: {
          traces: [],
          data: {
            evaluation: {
              comparedEvents: ['27c733aa-9bb0-44b4-9574-ef04e38eb4c4'],
              indicatorResults: [
                {
                  displayName: 'Response time P95',
                  keySli: false,
                  passTargets: [
                    {
                      criteria: '<=+10%',
                      targetValue: 25.43676878676464,
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
                    value: 301.7440018752324,
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
                'LS0tCnNwZWNfdmVyc2lvbjogIjEuMCIKY29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciCiAgY29tcGFyZV93aXRoOiAic2luZ2xlX3Jlc3VsdCIKICBpbmNsdWRlX3Jlc3VsdF93aXRoX3Njb3JlOiAicGFzcyIKICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxCmZpbHRlcjoKb2JqZWN0aXZlczoKICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1IgogICAgZGlzcGxheU5hbWU6ICJSZXNwb25zZSB0aW1lIFA5NSIKICAgIGtleV9zbGk6IGZhbHNlCiAgICBwYXNzOiAgICAgICAgICAgICAjIHBhc3MgaWYgKHJlbGF0aXZlIGNoYW5nZSA8PSAxMCUgQU5EIGFic29sdXRlIHZhbHVlIGlzIDwgNjAwbXMpCiAgICAgIC0gY3JpdGVyaWE6CiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykKICAgICAgICAgIC0gIjw2MDAiICAgICMgYWJzb2x1dGUgdmFsdWVzIG9ubHkgcmVxdWlyZSBhIGxvZ2ljYWwgb3BlcmF0b3IKICAgIHdhcm5pbmc6ICAgICAgICAgICMgaWYgdGhlIHJlc3BvbnNlIHRpbWUgaXMgYmVsb3cgODAwbXMsIHRoZSByZXN1bHQgc2hvdWxkIGJlIGEgd2FybmluZwogICAgICAtIGNyaXRlcmlhOgogICAgICAgICAgLSAiPD04MDAiCiAgICB3ZWlnaHQ6IDEKdG90YWxfc2NvcmU6CiAgcGFzczogIjkwJSIKICB3YXJuaW5nOiAiNzUlIg==',
              timeEnd: '2021-10-04T11:04:33Z',
              timeStart: '2021-10-04T11:02:17Z',
              sloFileContentParsed:
                '---\nspec_version: "1.0"\ncomparison:\n  aggregate_function: "avg"\n  compare_with: "single_result"\n  include_result_with_score: "pass"\n  number_of_comparison_results: 1\nfilter:\nobjectives:\n  - sli: "response_time_p95"\n    displayName: "Response time P95"\n    key_sli: false\n    pass:             # pass if (relative change <= 10% AND absolute value is < 600ms)\n      - criteria:\n          - "<=+10%"  # relative values require a prefixed sign (plus or minus)\n          - "<600"    # absolute values only require a logical operator\n    warning:          # if the response time is below 800ms, the result should be a warning\n      - criteria:\n          - "<=800"\n    weight: 1\ntotal_score:\n  pass: "90%"\n  warning: "75%"',
              score_pass: '90',
              score_warning: '75',
              compare_with: 'single_result\n',
              include_result_with_score: 'pass\n',
              number_of_comparison_results: 1,
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
            temporaryData: {
              distributor: {
                subscriptionID: '',
              },
            },
          },
          id: '0c7d757d-0b7b-447c-a801-71d1b7f51784',
          source: 'lighthouse-service',
          specversion: '1.0',
          time: '2021-10-04T11:05:46.915Z',
          type: 'sh.keptn.event.evaluation.finished',
          shkeptncontext: 'eb7d7e4f-f7cf-4a15-80ed-0167235b290e',
          shkeptnspecversion: '0.2.3',
          triggeredid: '4d6bea6d-4432-4573-959b-0b2141565e41',
          plainEvent: {
            data: {
              evaluation: {
                comparedEvents: ['27c733aa-9bb0-44b4-9574-ef04e38eb4c4'],
                indicatorResults: [
                  {
                    displayName: 'Response time P95',
                    keySli: false,
                    passTargets: [
                      {
                        criteria: '<=+10%',
                        targetValue: 25.43676878676464,
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
                      value: 301.7440018752324,
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
                  'LS0tCnNwZWNfdmVyc2lvbjogIjEuMCIKY29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciCiAgY29tcGFyZV93aXRoOiAic2luZ2xlX3Jlc3VsdCIKICBpbmNsdWRlX3Jlc3VsdF93aXRoX3Njb3JlOiAicGFzcyIKICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxCmZpbHRlcjoKb2JqZWN0aXZlczoKICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1IgogICAgZGlzcGxheU5hbWU6ICJSZXNwb25zZSB0aW1lIFA5NSIKICAgIGtleV9zbGk6IGZhbHNlCiAgICBwYXNzOiAgICAgICAgICAgICAjIHBhc3MgaWYgKHJlbGF0aXZlIGNoYW5nZSA8PSAxMCUgQU5EIGFic29sdXRlIHZhbHVlIGlzIDwgNjAwbXMpCiAgICAgIC0gY3JpdGVyaWE6CiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykKICAgICAgICAgIC0gIjw2MDAiICAgICMgYWJzb2x1dGUgdmFsdWVzIG9ubHkgcmVxdWlyZSBhIGxvZ2ljYWwgb3BlcmF0b3IKICAgIHdhcm5pbmc6ICAgICAgICAgICMgaWYgdGhlIHJlc3BvbnNlIHRpbWUgaXMgYmVsb3cgODAwbXMsIHRoZSByZXN1bHQgc2hvdWxkIGJlIGEgd2FybmluZwogICAgICAtIGNyaXRlcmlhOgogICAgICAgICAgLSAiPD04MDAiCiAgICB3ZWlnaHQ6IDEKdG90YWxfc2NvcmU6CiAgcGFzczogIjkwJSIKICB3YXJuaW5nOiAiNzUlIg==',
                timeEnd: '2021-10-04T11:04:33Z',
                timeStart: '2021-10-04T11:02:17Z',
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
              temporaryData: {
                distributor: {
                  subscriptionID: '',
                },
              },
            },
            id: '0c7d757d-0b7b-447c-a801-71d1b7f51784',
            source: 'lighthouse-service',
            specversion: '1.0',
            time: '2021-10-04T11:05:46.915Z',
            type: 'sh.keptn.event.evaluation.finished',
            shkeptncontext: 'eb7d7e4f-f7cf-4a15-80ed-0167235b290e',
            shkeptnspecversion: '0.2.3',
            triggeredid: '4d6bea6d-4432-4573-959b-0b2141565e41',
          },
          heatmapLabel: '2021-10-04 13:05',
        },
        color: '#dc172a',
        name: '2021-10-04 13:05',
        label: '2021-10-04 13:05',
        x: 4,
      },
      {
        y: 100,
        evaluationData: {
          traces: [],
          data: {
            evaluation: {
              comparedEvents: ['27c733aa-9bb0-44b4-9574-ef04e38eb4c4'],
              indicatorResults: [
                {
                  displayName: 'Response time P95',
                  keySli: false,
                  passTargets: [
                    {
                      criteria: '<=+10%',
                      targetValue: 25.43676878676464,
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
                    value: 9.280349173333333,
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
                'LS0tCnNwZWNfdmVyc2lvbjogIjEuMCIKY29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciCiAgY29tcGFyZV93aXRoOiAic2luZ2xlX3Jlc3VsdCIKICBpbmNsdWRlX3Jlc3VsdF93aXRoX3Njb3JlOiAicGFzcyIKICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxCmZpbHRlcjoKb2JqZWN0aXZlczoKICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1IgogICAgZGlzcGxheU5hbWU6ICJSZXNwb25zZSB0aW1lIFA5NSIKICAgIGtleV9zbGk6IGZhbHNlCiAgICBwYXNzOiAgICAgICAgICAgICAjIHBhc3MgaWYgKHJlbGF0aXZlIGNoYW5nZSA8PSAxMCUgQU5EIGFic29sdXRlIHZhbHVlIGlzIDwgNjAwbXMpCiAgICAgIC0gY3JpdGVyaWE6CiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykKICAgICAgICAgIC0gIjw2MDAiICAgICMgYWJzb2x1dGUgdmFsdWVzIG9ubHkgcmVxdWlyZSBhIGxvZ2ljYWwgb3BlcmF0b3IKICAgIHdhcm5pbmc6ICAgICAgICAgICMgaWYgdGhlIHJlc3BvbnNlIHRpbWUgaXMgYmVsb3cgODAwbXMsIHRoZSByZXN1bHQgc2hvdWxkIGJlIGEgd2FybmluZwogICAgICAtIGNyaXRlcmlhOgogICAgICAgICAgLSAiPD04MDAiCiAgICB3ZWlnaHQ6IDEKdG90YWxfc2NvcmU6CiAgcGFzczogIjkwJSIKICB3YXJuaW5nOiAiNzUlIg==',
              timeEnd: '2021-10-07T15:01:00.582Z',
              timeStart: '2021-10-07T14:56:00.582Z',
              sloFileContentParsed:
                '---\nspec_version: "1.0"\ncomparison:\n  aggregate_function: "avg"\n  compare_with: "single_result"\n  include_result_with_score: "pass"\n  number_of_comparison_results: 1\nfilter:\nobjectives:\n  - sli: "response_time_p95"\n    displayName: "Response time P95"\n    key_sli: false\n    pass:             # pass if (relative change <= 10% AND absolute value is < 600ms)\n      - criteria:\n          - "<=+10%"  # relative values require a prefixed sign (plus or minus)\n          - "<600"    # absolute values only require a logical operator\n    warning:          # if the response time is below 800ms, the result should be a warning\n      - criteria:\n          - "<=800"\n    weight: 1\ntotal_score:\n  pass: "90%"\n  warning: "75%"',
              score_pass: '90',
              score_warning: '75',
              compare_with: 'single_result\n',
              include_result_with_score: 'pass\n',
              number_of_comparison_results: 1,
            },
            labels: {
              DtCreds: 'dynatrace',
            },
            project: 'sockshop',
            result: 'pass',
            service: 'carts',
            stage: 'staging',
            status: 'succeeded',
            temporaryData: {
              distributor: {
                subscriptionID: '',
              },
            },
          },
          id: '53359c6f-fa31-4e8b-9fde-e003b3ea57ec',
          source: 'lighthouse-service',
          specversion: '1.0',
          time: '2021-10-07T15:01:24.717Z',
          type: 'sh.keptn.event.evaluation.finished',
          shkeptncontext: '68239e0e-e338-483a-be72-9f47cca3eaeb',
          shkeptnspecversion: '0.2.3',
          triggeredid: 'f0b1c0e1-3e33-48fe-b626-b52e508e9150',
          plainEvent: {
            data: {
              evaluation: {
                comparedEvents: ['27c733aa-9bb0-44b4-9574-ef04e38eb4c4'],
                indicatorResults: [
                  {
                    displayName: 'Response time P95',
                    keySli: false,
                    passTargets: [
                      {
                        criteria: '<=+10%',
                        targetValue: 25.43676878676464,
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
                      value: 9.280349173333333,
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
                  'LS0tCnNwZWNfdmVyc2lvbjogIjEuMCIKY29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciCiAgY29tcGFyZV93aXRoOiAic2luZ2xlX3Jlc3VsdCIKICBpbmNsdWRlX3Jlc3VsdF93aXRoX3Njb3JlOiAicGFzcyIKICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxCmZpbHRlcjoKb2JqZWN0aXZlczoKICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1IgogICAgZGlzcGxheU5hbWU6ICJSZXNwb25zZSB0aW1lIFA5NSIKICAgIGtleV9zbGk6IGZhbHNlCiAgICBwYXNzOiAgICAgICAgICAgICAjIHBhc3MgaWYgKHJlbGF0aXZlIGNoYW5nZSA8PSAxMCUgQU5EIGFic29sdXRlIHZhbHVlIGlzIDwgNjAwbXMpCiAgICAgIC0gY3JpdGVyaWE6CiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykKICAgICAgICAgIC0gIjw2MDAiICAgICMgYWJzb2x1dGUgdmFsdWVzIG9ubHkgcmVxdWlyZSBhIGxvZ2ljYWwgb3BlcmF0b3IKICAgIHdhcm5pbmc6ICAgICAgICAgICMgaWYgdGhlIHJlc3BvbnNlIHRpbWUgaXMgYmVsb3cgODAwbXMsIHRoZSByZXN1bHQgc2hvdWxkIGJlIGEgd2FybmluZwogICAgICAtIGNyaXRlcmlhOgogICAgICAgICAgLSAiPD04MDAiCiAgICB3ZWlnaHQ6IDEKdG90YWxfc2NvcmU6CiAgcGFzczogIjkwJSIKICB3YXJuaW5nOiAiNzUlIg==',
                timeEnd: '2021-10-07T15:01:00.582Z',
                timeStart: '2021-10-07T14:56:00.582Z',
              },
              labels: {
                DtCreds: 'dynatrace',
              },
              project: 'sockshop',
              result: 'pass',
              service: 'carts',
              stage: 'staging',
              status: 'succeeded',
              temporaryData: {
                distributor: {
                  subscriptionID: '',
                },
              },
            },
            id: '53359c6f-fa31-4e8b-9fde-e003b3ea57ec',
            source: 'lighthouse-service',
            specversion: '1.0',
            time: '2021-10-07T15:01:24.717Z',
            type: 'sh.keptn.event.evaluation.finished',
            shkeptncontext: '68239e0e-e338-483a-be72-9f47cca3eaeb',
            shkeptnspecversion: '0.2.3',
            triggeredid: 'f0b1c0e1-3e33-48fe-b626-b52e508e9150',
          },
          heatmapLabel: '2021-10-07 17:01',
        },
        color: '#7dc540',
        name: '2021-10-07 17:01',
        label: '2021-10-07 17:01',
        x: 5,
      },
      {
        y: 100,
        evaluationData: {
          traces: [],
          data: {
            evaluation: {
              comparedEvents: ['53359c6f-fa31-4e8b-9fde-e003b3ea57ec'],
              indicatorResults: [
                {
                  displayName: 'Response time P95',
                  keySli: false,
                  passTargets: [
                    {
                      criteria: '<=+10%',
                      targetValue: 10.208384090666666,
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
                    value: 9.329947741347496,
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
                'LS0tCnNwZWNfdmVyc2lvbjogIjEuMCIKY29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciCiAgY29tcGFyZV93aXRoOiAic2luZ2xlX3Jlc3VsdCIKICBpbmNsdWRlX3Jlc3VsdF93aXRoX3Njb3JlOiAicGFzcyIKICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxCmZpbHRlcjoKb2JqZWN0aXZlczoKICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1IgogICAgZGlzcGxheU5hbWU6ICJSZXNwb25zZSB0aW1lIFA5NSIKICAgIGtleV9zbGk6IGZhbHNlCiAgICBwYXNzOiAgICAgICAgICAgICAjIHBhc3MgaWYgKHJlbGF0aXZlIGNoYW5nZSA8PSAxMCUgQU5EIGFic29sdXRlIHZhbHVlIGlzIDwgNjAwbXMpCiAgICAgIC0gY3JpdGVyaWE6CiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykKICAgICAgICAgIC0gIjw2MDAiICAgICMgYWJzb2x1dGUgdmFsdWVzIG9ubHkgcmVxdWlyZSBhIGxvZ2ljYWwgb3BlcmF0b3IKICAgIHdhcm5pbmc6ICAgICAgICAgICMgaWYgdGhlIHJlc3BvbnNlIHRpbWUgaXMgYmVsb3cgODAwbXMsIHRoZSByZXN1bHQgc2hvdWxkIGJlIGEgd2FybmluZwogICAgICAtIGNyaXRlcmlhOgogICAgICAgICAgLSAiPD04MDAiCiAgICB3ZWlnaHQ6IDEKdG90YWxfc2NvcmU6CiAgcGFzczogIjkwJSIKICB3YXJuaW5nOiAiNzUlIg==',
              timeEnd: '2021-10-07T15:02:02.371Z',
              timeStart: '2021-10-07T14:57:02.371Z',
              sloFileContentParsed:
                '---\nspec_version: "1.0"\ncomparison:\n  aggregate_function: "avg"\n  compare_with: "single_result"\n  include_result_with_score: "pass"\n  number_of_comparison_results: 1\nfilter:\nobjectives:\n  - sli: "response_time_p95"\n    displayName: "Response time P95"\n    key_sli: false\n    pass:             # pass if (relative change <= 10% AND absolute value is < 600ms)\n      - criteria:\n          - "<=+10%"  # relative values require a prefixed sign (plus or minus)\n          - "<600"    # absolute values only require a logical operator\n    warning:          # if the response time is below 800ms, the result should be a warning\n      - criteria:\n          - "<=800"\n    weight: 1\ntotal_score:\n  pass: "90%"\n  warning: "75%"',
              score_pass: '90',
              score_warning: '75',
              compare_with: 'single_result\n',
              include_result_with_score: 'pass\n',
              number_of_comparison_results: 1,
            },
            labels: {
              DtCreds: 'dynatrace',
            },
            project: 'sockshop',
            result: 'pass',
            service: 'carts',
            stage: 'staging',
            status: 'succeeded',
            temporaryData: {
              distributor: {
                subscriptionID: '',
              },
            },
          },
          id: 'd7fba175-96b1-4a22-943a-78d871971925',
          source: 'lighthouse-service',
          specversion: '1.0',
          time: '2021-10-07T15:02:27.215Z',
          type: 'sh.keptn.event.evaluation.finished',
          shkeptncontext: '6f0636f7-8f01-4c2a-99b1-1cbf9ff8bfec',
          shkeptnspecversion: '0.2.3',
          triggeredid: '1ee69f43-4a51-46c1-a454-55c925ec4519',
          plainEvent: {
            data: {
              evaluation: {
                comparedEvents: ['53359c6f-fa31-4e8b-9fde-e003b3ea57ec'],
                indicatorResults: [
                  {
                    displayName: 'Response time P95',
                    keySli: false,
                    passTargets: [
                      {
                        criteria: '<=+10%',
                        targetValue: 10.208384090666666,
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
                      value: 9.329947741347496,
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
                  'LS0tCnNwZWNfdmVyc2lvbjogIjEuMCIKY29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciCiAgY29tcGFyZV93aXRoOiAic2luZ2xlX3Jlc3VsdCIKICBpbmNsdWRlX3Jlc3VsdF93aXRoX3Njb3JlOiAicGFzcyIKICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxCmZpbHRlcjoKb2JqZWN0aXZlczoKICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1IgogICAgZGlzcGxheU5hbWU6ICJSZXNwb25zZSB0aW1lIFA5NSIKICAgIGtleV9zbGk6IGZhbHNlCiAgICBwYXNzOiAgICAgICAgICAgICAjIHBhc3MgaWYgKHJlbGF0aXZlIGNoYW5nZSA8PSAxMCUgQU5EIGFic29sdXRlIHZhbHVlIGlzIDwgNjAwbXMpCiAgICAgIC0gY3JpdGVyaWE6CiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykKICAgICAgICAgIC0gIjw2MDAiICAgICMgYWJzb2x1dGUgdmFsdWVzIG9ubHkgcmVxdWlyZSBhIGxvZ2ljYWwgb3BlcmF0b3IKICAgIHdhcm5pbmc6ICAgICAgICAgICMgaWYgdGhlIHJlc3BvbnNlIHRpbWUgaXMgYmVsb3cgODAwbXMsIHRoZSByZXN1bHQgc2hvdWxkIGJlIGEgd2FybmluZwogICAgICAtIGNyaXRlcmlhOgogICAgICAgICAgLSAiPD04MDAiCiAgICB3ZWlnaHQ6IDEKdG90YWxfc2NvcmU6CiAgcGFzczogIjkwJSIKICB3YXJuaW5nOiAiNzUlIg==',
                timeEnd: '2021-10-07T15:02:02.371Z',
                timeStart: '2021-10-07T14:57:02.371Z',
              },
              labels: {
                DtCreds: 'dynatrace',
              },
              project: 'sockshop',
              result: 'pass',
              service: 'carts',
              stage: 'staging',
              status: 'succeeded',
              temporaryData: {
                distributor: {
                  subscriptionID: '',
                },
              },
            },
            id: 'd7fba175-96b1-4a22-943a-78d871971925',
            source: 'lighthouse-service',
            specversion: '1.0',
            time: '2021-10-07T15:02:27.215Z',
            type: 'sh.keptn.event.evaluation.finished',
            shkeptncontext: '6f0636f7-8f01-4c2a-99b1-1cbf9ff8bfec',
            shkeptnspecversion: '0.2.3',
            triggeredid: '1ee69f43-4a51-46c1-a454-55c925ec4519',
          },
          heatmapLabel: '2021-10-07 17:02',
        },
        color: '#7dc540',
        name: '2021-10-07 17:02',
        label: '2021-10-07 17:02',
        x: 6,
      },
    ],
    cursor: 'pointer',
    turboThreshold: 0,
  },
  {
    name: 'Score',
    metricName: 'Score',
    type: 'line',
    data: [
      {
        y: 100,
        evaluationData: {
          traces: [],
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
                    value: 298.9114676593789,
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
                'LS0tCnNwZWNfdmVyc2lvbjogIjEuMCIKY29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciCiAgY29tcGFyZV93aXRoOiAic2luZ2xlX3Jlc3VsdCIKICBpbmNsdWRlX3Jlc3VsdF93aXRoX3Njb3JlOiAicGFzcyIKICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxCmZpbHRlcjoKb2JqZWN0aXZlczoKICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1IgogICAgZGlzcGxheU5hbWU6ICJSZXNwb25zZSB0aW1lIFA5NSIKICAgIGtleV9zbGk6IGZhbHNlCiAgICBwYXNzOiAgICAgICAgICAgICAjIHBhc3MgaWYgKHJlbGF0aXZlIGNoYW5nZSA8PSAxMCUgQU5EIGFic29sdXRlIHZhbHVlIGlzIDwgNjAwbXMpCiAgICAgIC0gY3JpdGVyaWE6CiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykKICAgICAgICAgIC0gIjw2MDAiICAgICMgYWJzb2x1dGUgdmFsdWVzIG9ubHkgcmVxdWlyZSBhIGxvZ2ljYWwgb3BlcmF0b3IKICAgIHdhcm5pbmc6ICAgICAgICAgICMgaWYgdGhlIHJlc3BvbnNlIHRpbWUgaXMgYmVsb3cgODAwbXMsIHRoZSByZXN1bHQgc2hvdWxkIGJlIGEgd2FybmluZwogICAgICAtIGNyaXRlcmlhOgogICAgICAgICAgLSAiPD04MDAiCiAgICB3ZWlnaHQ6IDEKdG90YWxfc2NvcmU6CiAgcGFzczogIjkwJSIKICB3YXJuaW5nOiAiNzUlIg==',
              timeEnd: '2021-09-30T11:41:15Z',
              timeStart: '2021-09-30T11:39:08Z',
              sloFileContentParsed:
                '---\nspec_version: "1.0"\ncomparison:\n  aggregate_function: "avg"\n  compare_with: "single_result"\n  include_result_with_score: "pass"\n  number_of_comparison_results: 1\nfilter:\nobjectives:\n  - sli: "response_time_p95"\n    displayName: "Response time P95"\n    key_sli: false\n    pass:             # pass if (relative change <= 10% AND absolute value is < 600ms)\n      - criteria:\n          - "<=+10%"  # relative values require a prefixed sign (plus or minus)\n          - "<600"    # absolute values only require a logical operator\n    warning:          # if the response time is below 800ms, the result should be a warning\n      - criteria:\n          - "<=800"\n    weight: 1\ntotal_score:\n  pass: "90%"\n  warning: "75%"',
              score_pass: '90',
              score_warning: '75',
              compare_with: 'single_result\n',
              include_result_with_score: 'pass\n',
              number_of_comparison_results: 0,
            },
            labels: {
              DtCreds: 'dynatrace',
            },
            project: 'sockshop',
            result: 'pass',
            service: 'carts',
            stage: 'staging',
            status: 'succeeded',
            temporaryData: {
              distributor: {
                subscriptionID: '',
              },
            },
          },
          id: '40cfcfc2-9cc6-4fce-ba62-5f5fec3e3e3b',
          source: 'lighthouse-service',
          specversion: '1.0',
          time: '2021-09-30T11:42:29.616Z',
          type: 'sh.keptn.event.evaluation.finished',
          shkeptncontext: '369b15c8-5410-49de-87f7-08a2f9f16243',
          shkeptnspecversion: '0.2.3',
          triggeredid: '2d2e799f-6c25-4b26-83ab-12953335f6d4',
          plainEvent: {
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
                      value: 298.9114676593789,
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
                  'LS0tCnNwZWNfdmVyc2lvbjogIjEuMCIKY29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciCiAgY29tcGFyZV93aXRoOiAic2luZ2xlX3Jlc3VsdCIKICBpbmNsdWRlX3Jlc3VsdF93aXRoX3Njb3JlOiAicGFzcyIKICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxCmZpbHRlcjoKb2JqZWN0aXZlczoKICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1IgogICAgZGlzcGxheU5hbWU6ICJSZXNwb25zZSB0aW1lIFA5NSIKICAgIGtleV9zbGk6IGZhbHNlCiAgICBwYXNzOiAgICAgICAgICAgICAjIHBhc3MgaWYgKHJlbGF0aXZlIGNoYW5nZSA8PSAxMCUgQU5EIGFic29sdXRlIHZhbHVlIGlzIDwgNjAwbXMpCiAgICAgIC0gY3JpdGVyaWE6CiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykKICAgICAgICAgIC0gIjw2MDAiICAgICMgYWJzb2x1dGUgdmFsdWVzIG9ubHkgcmVxdWlyZSBhIGxvZ2ljYWwgb3BlcmF0b3IKICAgIHdhcm5pbmc6ICAgICAgICAgICMgaWYgdGhlIHJlc3BvbnNlIHRpbWUgaXMgYmVsb3cgODAwbXMsIHRoZSByZXN1bHQgc2hvdWxkIGJlIGEgd2FybmluZwogICAgICAtIGNyaXRlcmlhOgogICAgICAgICAgLSAiPD04MDAiCiAgICB3ZWlnaHQ6IDEKdG90YWxfc2NvcmU6CiAgcGFzczogIjkwJSIKICB3YXJuaW5nOiAiNzUlIg==',
                timeEnd: '2021-09-30T11:41:15Z',
                timeStart: '2021-09-30T11:39:08Z',
              },
              labels: {
                DtCreds: 'dynatrace',
              },
              project: 'sockshop',
              result: 'pass',
              service: 'carts',
              stage: 'staging',
              status: 'succeeded',
              temporaryData: {
                distributor: {
                  subscriptionID: '',
                },
              },
            },
            id: '40cfcfc2-9cc6-4fce-ba62-5f5fec3e3e3b',
            source: 'lighthouse-service',
            specversion: '1.0',
            time: '2021-09-30T11:42:29.616Z',
            type: 'sh.keptn.event.evaluation.finished',
            shkeptncontext: '369b15c8-5410-49de-87f7-08a2f9f16243',
            shkeptnspecversion: '0.2.3',
            triggeredid: '2d2e799f-6c25-4b26-83ab-12953335f6d4',
          },
          heatmapLabel: '2021-09-30 13:42',
        },
        color: '#7dc540',
        name: '2021-09-30 13:42',
        label: '2021-09-30 13:42',
        x: 0,
      },
      {
        y: 0,
        evaluationData: {
          traces: [],
          data: {
            evaluation: {
              comparedEvents: ['40cfcfc2-9cc6-4fce-ba62-5f5fec3e3e3b'],
              indicatorResults: [
                {
                  displayName: 'Response time P95',
                  keySli: false,
                  passTargets: [
                    {
                      criteria: '<=+10%',
                      targetValue: 328.8026144253168,
                      violated: true,
                    },
                    {
                      criteria: '<600',
                      targetValue: 600,
                      violated: true,
                    },
                  ],
                  score: 0,
                  status: 'fail',
                  value: {
                    metric: 'response_time_p95',
                    success: true,
                    value: 1091.8809833281025,
                  },
                  warningTargets: [
                    {
                      criteria: '<=800',
                      targetValue: 800,
                      violated: true,
                    },
                  ],
                },
              ],
              result: 'fail',
              score: 0,
              sloFileContent:
                'LS0tCnNwZWNfdmVyc2lvbjogIjEuMCIKY29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciCiAgY29tcGFyZV93aXRoOiAic2luZ2xlX3Jlc3VsdCIKICBpbmNsdWRlX3Jlc3VsdF93aXRoX3Njb3JlOiAicGFzcyIKICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxCmZpbHRlcjoKb2JqZWN0aXZlczoKICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1IgogICAgZGlzcGxheU5hbWU6ICJSZXNwb25zZSB0aW1lIFA5NSIKICAgIGtleV9zbGk6IGZhbHNlCiAgICBwYXNzOiAgICAgICAgICAgICAjIHBhc3MgaWYgKHJlbGF0aXZlIGNoYW5nZSA8PSAxMCUgQU5EIGFic29sdXRlIHZhbHVlIGlzIDwgNjAwbXMpCiAgICAgIC0gY3JpdGVyaWE6CiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykKICAgICAgICAgIC0gIjw2MDAiICAgICMgYWJzb2x1dGUgdmFsdWVzIG9ubHkgcmVxdWlyZSBhIGxvZ2ljYWwgb3BlcmF0b3IKICAgIHdhcm5pbmc6ICAgICAgICAgICMgaWYgdGhlIHJlc3BvbnNlIHRpbWUgaXMgYmVsb3cgODAwbXMsIHRoZSByZXN1bHQgc2hvdWxkIGJlIGEgd2FybmluZwogICAgICAtIGNyaXRlcmlhOgogICAgICAgICAgLSAiPD04MDAiCiAgICB3ZWlnaHQ6IDEKdG90YWxfc2NvcmU6CiAgcGFzczogIjkwJSIKICB3YXJuaW5nOiAiNzUlIg==',
              timeEnd: '2021-10-01T09:25:48Z',
              timeStart: '2021-10-01T09:16:31Z',
              sloFileContentParsed:
                '---\nspec_version: "1.0"\ncomparison:\n  aggregate_function: "avg"\n  compare_with: "single_result"\n  include_result_with_score: "pass"\n  number_of_comparison_results: 1\nfilter:\nobjectives:\n  - sli: "response_time_p95"\n    displayName: "Response time P95"\n    key_sli: false\n    pass:             # pass if (relative change <= 10% AND absolute value is < 600ms)\n      - criteria:\n          - "<=+10%"  # relative values require a prefixed sign (plus or minus)\n          - "<600"    # absolute values only require a logical operator\n    warning:          # if the response time is below 800ms, the result should be a warning\n      - criteria:\n          - "<=800"\n    weight: 1\ntotal_score:\n  pass: "90%"\n  warning: "75%"',
              score_pass: '90',
              score_warning: '75',
              compare_with: 'single_result\n',
              include_result_with_score: 'pass\n',
              number_of_comparison_results: 1,
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
            temporaryData: {
              distributor: {
                subscriptionID: '',
              },
            },
          },
          id: 'f623de99-a716-4b9b-90a2-8df1054e92c7',
          source: 'lighthouse-service',
          specversion: '1.0',
          time: '2021-10-01T09:26:15.608Z',
          type: 'sh.keptn.event.evaluation.finished',
          shkeptncontext: '081ce5a3-e0f9-4c4d-bfcf-5893f85bfaff',
          shkeptnspecversion: '0.2.3',
          triggeredid: 'f564369b-2ee6-4c32-8488-03508ed05b5a',
          plainEvent: {
            data: {
              evaluation: {
                comparedEvents: ['40cfcfc2-9cc6-4fce-ba62-5f5fec3e3e3b'],
                indicatorResults: [
                  {
                    displayName: 'Response time P95',
                    keySli: false,
                    passTargets: [
                      {
                        criteria: '<=+10%',
                        targetValue: 328.8026144253168,
                        violated: true,
                      },
                      {
                        criteria: '<600',
                        targetValue: 600,
                        violated: true,
                      },
                    ],
                    score: 0,
                    status: 'fail',
                    value: {
                      metric: 'response_time_p95',
                      success: true,
                      value: 1091.8809833281025,
                    },
                    warningTargets: [
                      {
                        criteria: '<=800',
                        targetValue: 800,
                        violated: true,
                      },
                    ],
                  },
                ],
                result: 'fail',
                score: 0,
                sloFileContent:
                  'LS0tCnNwZWNfdmVyc2lvbjogIjEuMCIKY29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciCiAgY29tcGFyZV93aXRoOiAic2luZ2xlX3Jlc3VsdCIKICBpbmNsdWRlX3Jlc3VsdF93aXRoX3Njb3JlOiAicGFzcyIKICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxCmZpbHRlcjoKb2JqZWN0aXZlczoKICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1IgogICAgZGlzcGxheU5hbWU6ICJSZXNwb25zZSB0aW1lIFA5NSIKICAgIGtleV9zbGk6IGZhbHNlCiAgICBwYXNzOiAgICAgICAgICAgICAjIHBhc3MgaWYgKHJlbGF0aXZlIGNoYW5nZSA8PSAxMCUgQU5EIGFic29sdXRlIHZhbHVlIGlzIDwgNjAwbXMpCiAgICAgIC0gY3JpdGVyaWE6CiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykKICAgICAgICAgIC0gIjw2MDAiICAgICMgYWJzb2x1dGUgdmFsdWVzIG9ubHkgcmVxdWlyZSBhIGxvZ2ljYWwgb3BlcmF0b3IKICAgIHdhcm5pbmc6ICAgICAgICAgICMgaWYgdGhlIHJlc3BvbnNlIHRpbWUgaXMgYmVsb3cgODAwbXMsIHRoZSByZXN1bHQgc2hvdWxkIGJlIGEgd2FybmluZwogICAgICAtIGNyaXRlcmlhOgogICAgICAgICAgLSAiPD04MDAiCiAgICB3ZWlnaHQ6IDEKdG90YWxfc2NvcmU6CiAgcGFzczogIjkwJSIKICB3YXJuaW5nOiAiNzUlIg==',
                timeEnd: '2021-10-01T09:25:48Z',
                timeStart: '2021-10-01T09:16:31Z',
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
              temporaryData: {
                distributor: {
                  subscriptionID: '',
                },
              },
            },
            id: 'f623de99-a716-4b9b-90a2-8df1054e92c7',
            source: 'lighthouse-service',
            specversion: '1.0',
            time: '2021-10-01T09:26:15.608Z',
            type: 'sh.keptn.event.evaluation.finished',
            shkeptncontext: '081ce5a3-e0f9-4c4d-bfcf-5893f85bfaff',
            shkeptnspecversion: '0.2.3',
            triggeredid: 'f564369b-2ee6-4c32-8488-03508ed05b5a',
          },
          heatmapLabel: '2021-10-01 11:26',
        },
        color: '#dc172a',
        name: '2021-10-01 11:26',
        label: '2021-10-01 11:26',
        x: 1,
      },
      {
        y: 100,
        evaluationData: {
          traces: [],
          data: {
            evaluation: {
              comparedEvents: ['40cfcfc2-9cc6-4fce-ba62-5f5fec3e3e3b'],
              indicatorResults: [
                {
                  displayName: 'Response time P95',
                  keySli: false,
                  passTargets: [
                    {
                      criteria: '<=+10%',
                      targetValue: 328.8026144253168,
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
                    value: 306.3587409268154,
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
                'LS0tCnNwZWNfdmVyc2lvbjogIjEuMCIKY29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciCiAgY29tcGFyZV93aXRoOiAic2luZ2xlX3Jlc3VsdCIKICBpbmNsdWRlX3Jlc3VsdF93aXRoX3Njb3JlOiAicGFzcyIKICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxCmZpbHRlcjoKb2JqZWN0aXZlczoKICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1IgogICAgZGlzcGxheU5hbWU6ICJSZXNwb25zZSB0aW1lIFA5NSIKICAgIGtleV9zbGk6IGZhbHNlCiAgICBwYXNzOiAgICAgICAgICAgICAjIHBhc3MgaWYgKHJlbGF0aXZlIGNoYW5nZSA8PSAxMCUgQU5EIGFic29sdXRlIHZhbHVlIGlzIDwgNjAwbXMpCiAgICAgIC0gY3JpdGVyaWE6CiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykKICAgICAgICAgIC0gIjw2MDAiICAgICMgYWJzb2x1dGUgdmFsdWVzIG9ubHkgcmVxdWlyZSBhIGxvZ2ljYWwgb3BlcmF0b3IKICAgIHdhcm5pbmc6ICAgICAgICAgICMgaWYgdGhlIHJlc3BvbnNlIHRpbWUgaXMgYmVsb3cgODAwbXMsIHRoZSByZXN1bHQgc2hvdWxkIGJlIGEgd2FybmluZwogICAgICAtIGNyaXRlcmlhOgogICAgICAgICAgLSAiPD04MDAiCiAgICB3ZWlnaHQ6IDEKdG90YWxfc2NvcmU6CiAgcGFzczogIjkwJSIKICB3YXJuaW5nOiAiNzUlIg==',
              timeEnd: '2021-10-01T10:26:34Z',
              timeStart: '2021-10-01T10:24:15Z',
              sloFileContentParsed:
                '---\nspec_version: "1.0"\ncomparison:\n  aggregate_function: "avg"\n  compare_with: "single_result"\n  include_result_with_score: "pass"\n  number_of_comparison_results: 1\nfilter:\nobjectives:\n  - sli: "response_time_p95"\n    displayName: "Response time P95"\n    key_sli: false\n    pass:             # pass if (relative change <= 10% AND absolute value is < 600ms)\n      - criteria:\n          - "<=+10%"  # relative values require a prefixed sign (plus or minus)\n          - "<600"    # absolute values only require a logical operator\n    warning:          # if the response time is below 800ms, the result should be a warning\n      - criteria:\n          - "<=800"\n    weight: 1\ntotal_score:\n  pass: "90%"\n  warning: "75%"',
              score_pass: '90',
              score_warning: '75',
              compare_with: 'single_result\n',
              include_result_with_score: 'pass\n',
              number_of_comparison_results: 1,
            },
            labels: {
              DtCreds: 'dynatrace',
            },
            project: 'sockshop',
            result: 'pass',
            service: 'carts',
            stage: 'staging',
            status: 'succeeded',
            temporaryData: {
              distributor: {
                subscriptionID: '',
              },
            },
          },
          id: '3344487d-e384-4cd9-a0e0-fcf157a33ad6',
          source: 'lighthouse-service',
          specversion: '1.0',
          time: '2021-10-01T10:27:49.318Z',
          type: 'sh.keptn.event.evaluation.finished',
          shkeptncontext: 'e9dd1b2f-8383-48a9-a2d4-36cb8ff370c7',
          shkeptnspecversion: '0.2.3',
          triggeredid: '61de0ee6-b535-42bb-ac30-d0f359298bb0',
          plainEvent: {
            data: {
              evaluation: {
                comparedEvents: ['40cfcfc2-9cc6-4fce-ba62-5f5fec3e3e3b'],
                indicatorResults: [
                  {
                    displayName: 'Response time P95',
                    keySli: false,
                    passTargets: [
                      {
                        criteria: '<=+10%',
                        targetValue: 328.8026144253168,
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
                      value: 306.3587409268154,
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
                  'LS0tCnNwZWNfdmVyc2lvbjogIjEuMCIKY29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciCiAgY29tcGFyZV93aXRoOiAic2luZ2xlX3Jlc3VsdCIKICBpbmNsdWRlX3Jlc3VsdF93aXRoX3Njb3JlOiAicGFzcyIKICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxCmZpbHRlcjoKb2JqZWN0aXZlczoKICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1IgogICAgZGlzcGxheU5hbWU6ICJSZXNwb25zZSB0aW1lIFA5NSIKICAgIGtleV9zbGk6IGZhbHNlCiAgICBwYXNzOiAgICAgICAgICAgICAjIHBhc3MgaWYgKHJlbGF0aXZlIGNoYW5nZSA8PSAxMCUgQU5EIGFic29sdXRlIHZhbHVlIGlzIDwgNjAwbXMpCiAgICAgIC0gY3JpdGVyaWE6CiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykKICAgICAgICAgIC0gIjw2MDAiICAgICMgYWJzb2x1dGUgdmFsdWVzIG9ubHkgcmVxdWlyZSBhIGxvZ2ljYWwgb3BlcmF0b3IKICAgIHdhcm5pbmc6ICAgICAgICAgICMgaWYgdGhlIHJlc3BvbnNlIHRpbWUgaXMgYmVsb3cgODAwbXMsIHRoZSByZXN1bHQgc2hvdWxkIGJlIGEgd2FybmluZwogICAgICAtIGNyaXRlcmlhOgogICAgICAgICAgLSAiPD04MDAiCiAgICB3ZWlnaHQ6IDEKdG90YWxfc2NvcmU6CiAgcGFzczogIjkwJSIKICB3YXJuaW5nOiAiNzUlIg==',
                timeEnd: '2021-10-01T10:26:34Z',
                timeStart: '2021-10-01T10:24:15Z',
              },
              labels: {
                DtCreds: 'dynatrace',
              },
              project: 'sockshop',
              result: 'pass',
              service: 'carts',
              stage: 'staging',
              status: 'succeeded',
              temporaryData: {
                distributor: {
                  subscriptionID: '',
                },
              },
            },
            id: '3344487d-e384-4cd9-a0e0-fcf157a33ad6',
            source: 'lighthouse-service',
            specversion: '1.0',
            time: '2021-10-01T10:27:49.318Z',
            type: 'sh.keptn.event.evaluation.finished',
            shkeptncontext: 'e9dd1b2f-8383-48a9-a2d4-36cb8ff370c7',
            shkeptnspecversion: '0.2.3',
            triggeredid: '61de0ee6-b535-42bb-ac30-d0f359298bb0',
          },
          heatmapLabel: '2021-10-01 12:27',
        },
        color: '#7dc540',
        name: '2021-10-01 12:27',
        label: '2021-10-01 12:27',
        x: 2,
      },
      {
        y: 100,
        evaluationData: {
          traces: [],
          data: {
            evaluation: {
              comparedEvents: ['3344487d-e384-4cd9-a0e0-fcf157a33ad6'],
              indicatorResults: [
                {
                  displayName: 'Response time P95',
                  keySli: false,
                  passTargets: [
                    {
                      criteria: '<=+10%',
                      targetValue: 336.9946150194969,
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
                    value: 23.12433526069513,
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
                'LS0tCnNwZWNfdmVyc2lvbjogIjEuMCIKY29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciCiAgY29tcGFyZV93aXRoOiAic2luZ2xlX3Jlc3VsdCIKICBpbmNsdWRlX3Jlc3VsdF93aXRoX3Njb3JlOiAicGFzcyIKICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxCmZpbHRlcjoKb2JqZWN0aXZlczoKICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1IgogICAgZGlzcGxheU5hbWU6ICJSZXNwb25zZSB0aW1lIFA5NSIKICAgIGtleV9zbGk6IGZhbHNlCiAgICBwYXNzOiAgICAgICAgICAgICAjIHBhc3MgaWYgKHJlbGF0aXZlIGNoYW5nZSA8PSAxMCUgQU5EIGFic29sdXRlIHZhbHVlIGlzIDwgNjAwbXMpCiAgICAgIC0gY3JpdGVyaWE6CiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykKICAgICAgICAgIC0gIjw2MDAiICAgICMgYWJzb2x1dGUgdmFsdWVzIG9ubHkgcmVxdWlyZSBhIGxvZ2ljYWwgb3BlcmF0b3IKICAgIHdhcm5pbmc6ICAgICAgICAgICMgaWYgdGhlIHJlc3BvbnNlIHRpbWUgaXMgYmVsb3cgODAwbXMsIHRoZSByZXN1bHQgc2hvdWxkIGJlIGEgd2FybmluZwogICAgICAtIGNyaXRlcmlhOgogICAgICAgICAgLSAiPD04MDAiCiAgICB3ZWlnaHQ6IDEKdG90YWxfc2NvcmU6CiAgcGFzczogIjkwJSIKICB3YXJuaW5nOiAiNzUlIg==',
              timeEnd: '2021-10-01T10:32:47.372Z',
              timeStart: '2021-10-01T10:27:47.372Z',
              sloFileContentParsed:
                '---\nspec_version: "1.0"\ncomparison:\n  aggregate_function: "avg"\n  compare_with: "single_result"\n  include_result_with_score: "pass"\n  number_of_comparison_results: 1\nfilter:\nobjectives:\n  - sli: "response_time_p95"\n    displayName: "Response time P95"\n    key_sli: false\n    pass:             # pass if (relative change <= 10% AND absolute value is < 600ms)\n      - criteria:\n          - "<=+10%"  # relative values require a prefixed sign (plus or minus)\n          - "<600"    # absolute values only require a logical operator\n    warning:          # if the response time is below 800ms, the result should be a warning\n      - criteria:\n          - "<=800"\n    weight: 1\ntotal_score:\n  pass: "90%"\n  warning: "75%"',
              score_pass: '90',
              score_warning: '75',
              compare_with: 'single_result\n',
              include_result_with_score: 'pass\n',
              number_of_comparison_results: 1,
            },
            labels: {
              DtCreds: 'dynatrace',
            },
            project: 'sockshop',
            result: 'pass',
            service: 'carts',
            stage: 'staging',
            status: 'succeeded',
            temporaryData: {
              distributor: {
                subscriptionID: '',
              },
            },
          },
          id: '27c733aa-9bb0-44b4-9574-ef04e38eb4c4',
          source: 'lighthouse-service',
          specversion: '1.0',
          time: '2021-10-01T10:33:19.708Z',
          type: 'sh.keptn.event.evaluation.finished',
          shkeptncontext: '199287a3-b723-479b-8083-8b3ff3482531',
          shkeptnspecversion: '0.2.3',
          triggeredid: '53ac99d5-2828-470a-b067-442562cf3bca',
          plainEvent: {
            data: {
              evaluation: {
                comparedEvents: ['3344487d-e384-4cd9-a0e0-fcf157a33ad6'],
                indicatorResults: [
                  {
                    displayName: 'Response time P95',
                    keySli: false,
                    passTargets: [
                      {
                        criteria: '<=+10%',
                        targetValue: 336.9946150194969,
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
                      value: 23.12433526069513,
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
                  'LS0tCnNwZWNfdmVyc2lvbjogIjEuMCIKY29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciCiAgY29tcGFyZV93aXRoOiAic2luZ2xlX3Jlc3VsdCIKICBpbmNsdWRlX3Jlc3VsdF93aXRoX3Njb3JlOiAicGFzcyIKICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxCmZpbHRlcjoKb2JqZWN0aXZlczoKICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1IgogICAgZGlzcGxheU5hbWU6ICJSZXNwb25zZSB0aW1lIFA5NSIKICAgIGtleV9zbGk6IGZhbHNlCiAgICBwYXNzOiAgICAgICAgICAgICAjIHBhc3MgaWYgKHJlbGF0aXZlIGNoYW5nZSA8PSAxMCUgQU5EIGFic29sdXRlIHZhbHVlIGlzIDwgNjAwbXMpCiAgICAgIC0gY3JpdGVyaWE6CiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykKICAgICAgICAgIC0gIjw2MDAiICAgICMgYWJzb2x1dGUgdmFsdWVzIG9ubHkgcmVxdWlyZSBhIGxvZ2ljYWwgb3BlcmF0b3IKICAgIHdhcm5pbmc6ICAgICAgICAgICMgaWYgdGhlIHJlc3BvbnNlIHRpbWUgaXMgYmVsb3cgODAwbXMsIHRoZSByZXN1bHQgc2hvdWxkIGJlIGEgd2FybmluZwogICAgICAtIGNyaXRlcmlhOgogICAgICAgICAgLSAiPD04MDAiCiAgICB3ZWlnaHQ6IDEKdG90YWxfc2NvcmU6CiAgcGFzczogIjkwJSIKICB3YXJuaW5nOiAiNzUlIg==',
                timeEnd: '2021-10-01T10:32:47.372Z',
                timeStart: '2021-10-01T10:27:47.372Z',
              },
              labels: {
                DtCreds: 'dynatrace',
              },
              project: 'sockshop',
              result: 'pass',
              service: 'carts',
              stage: 'staging',
              status: 'succeeded',
              temporaryData: {
                distributor: {
                  subscriptionID: '',
                },
              },
            },
            id: '27c733aa-9bb0-44b4-9574-ef04e38eb4c4',
            source: 'lighthouse-service',
            specversion: '1.0',
            time: '2021-10-01T10:33:19.708Z',
            type: 'sh.keptn.event.evaluation.finished',
            shkeptncontext: '199287a3-b723-479b-8083-8b3ff3482531',
            shkeptnspecversion: '0.2.3',
            triggeredid: '53ac99d5-2828-470a-b067-442562cf3bca',
          },
          heatmapLabel: '2021-10-01 12:33',
        },
        color: '#7dc540',
        name: '2021-10-01 12:33',
        label: '2021-10-01 12:33',
        x: 3,
      },
      {
        y: 50,
        evaluationData: {
          traces: [],
          data: {
            evaluation: {
              comparedEvents: ['27c733aa-9bb0-44b4-9574-ef04e38eb4c4'],
              indicatorResults: [
                {
                  displayName: 'Response time P95',
                  keySli: false,
                  passTargets: [
                    {
                      criteria: '<=+10%',
                      targetValue: 25.43676878676464,
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
                    value: 301.7440018752324,
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
                'LS0tCnNwZWNfdmVyc2lvbjogIjEuMCIKY29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciCiAgY29tcGFyZV93aXRoOiAic2luZ2xlX3Jlc3VsdCIKICBpbmNsdWRlX3Jlc3VsdF93aXRoX3Njb3JlOiAicGFzcyIKICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxCmZpbHRlcjoKb2JqZWN0aXZlczoKICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1IgogICAgZGlzcGxheU5hbWU6ICJSZXNwb25zZSB0aW1lIFA5NSIKICAgIGtleV9zbGk6IGZhbHNlCiAgICBwYXNzOiAgICAgICAgICAgICAjIHBhc3MgaWYgKHJlbGF0aXZlIGNoYW5nZSA8PSAxMCUgQU5EIGFic29sdXRlIHZhbHVlIGlzIDwgNjAwbXMpCiAgICAgIC0gY3JpdGVyaWE6CiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykKICAgICAgICAgIC0gIjw2MDAiICAgICMgYWJzb2x1dGUgdmFsdWVzIG9ubHkgcmVxdWlyZSBhIGxvZ2ljYWwgb3BlcmF0b3IKICAgIHdhcm5pbmc6ICAgICAgICAgICMgaWYgdGhlIHJlc3BvbnNlIHRpbWUgaXMgYmVsb3cgODAwbXMsIHRoZSByZXN1bHQgc2hvdWxkIGJlIGEgd2FybmluZwogICAgICAtIGNyaXRlcmlhOgogICAgICAgICAgLSAiPD04MDAiCiAgICB3ZWlnaHQ6IDEKdG90YWxfc2NvcmU6CiAgcGFzczogIjkwJSIKICB3YXJuaW5nOiAiNzUlIg==',
              timeEnd: '2021-10-04T11:04:33Z',
              timeStart: '2021-10-04T11:02:17Z',
              sloFileContentParsed:
                '---\nspec_version: "1.0"\ncomparison:\n  aggregate_function: "avg"\n  compare_with: "single_result"\n  include_result_with_score: "pass"\n  number_of_comparison_results: 1\nfilter:\nobjectives:\n  - sli: "response_time_p95"\n    displayName: "Response time P95"\n    key_sli: false\n    pass:             # pass if (relative change <= 10% AND absolute value is < 600ms)\n      - criteria:\n          - "<=+10%"  # relative values require a prefixed sign (plus or minus)\n          - "<600"    # absolute values only require a logical operator\n    warning:          # if the response time is below 800ms, the result should be a warning\n      - criteria:\n          - "<=800"\n    weight: 1\ntotal_score:\n  pass: "90%"\n  warning: "75%"',
              score_pass: '90',
              score_warning: '75',
              compare_with: 'single_result\n',
              include_result_with_score: 'pass\n',
              number_of_comparison_results: 1,
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
            temporaryData: {
              distributor: {
                subscriptionID: '',
              },
            },
          },
          id: '0c7d757d-0b7b-447c-a801-71d1b7f51784',
          source: 'lighthouse-service',
          specversion: '1.0',
          time: '2021-10-04T11:05:46.915Z',
          type: 'sh.keptn.event.evaluation.finished',
          shkeptncontext: 'eb7d7e4f-f7cf-4a15-80ed-0167235b290e',
          shkeptnspecversion: '0.2.3',
          triggeredid: '4d6bea6d-4432-4573-959b-0b2141565e41',
          plainEvent: {
            data: {
              evaluation: {
                comparedEvents: ['27c733aa-9bb0-44b4-9574-ef04e38eb4c4'],
                indicatorResults: [
                  {
                    displayName: 'Response time P95',
                    keySli: false,
                    passTargets: [
                      {
                        criteria: '<=+10%',
                        targetValue: 25.43676878676464,
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
                      value: 301.7440018752324,
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
                  'LS0tCnNwZWNfdmVyc2lvbjogIjEuMCIKY29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciCiAgY29tcGFyZV93aXRoOiAic2luZ2xlX3Jlc3VsdCIKICBpbmNsdWRlX3Jlc3VsdF93aXRoX3Njb3JlOiAicGFzcyIKICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxCmZpbHRlcjoKb2JqZWN0aXZlczoKICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1IgogICAgZGlzcGxheU5hbWU6ICJSZXNwb25zZSB0aW1lIFA5NSIKICAgIGtleV9zbGk6IGZhbHNlCiAgICBwYXNzOiAgICAgICAgICAgICAjIHBhc3MgaWYgKHJlbGF0aXZlIGNoYW5nZSA8PSAxMCUgQU5EIGFic29sdXRlIHZhbHVlIGlzIDwgNjAwbXMpCiAgICAgIC0gY3JpdGVyaWE6CiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykKICAgICAgICAgIC0gIjw2MDAiICAgICMgYWJzb2x1dGUgdmFsdWVzIG9ubHkgcmVxdWlyZSBhIGxvZ2ljYWwgb3BlcmF0b3IKICAgIHdhcm5pbmc6ICAgICAgICAgICMgaWYgdGhlIHJlc3BvbnNlIHRpbWUgaXMgYmVsb3cgODAwbXMsIHRoZSByZXN1bHQgc2hvdWxkIGJlIGEgd2FybmluZwogICAgICAtIGNyaXRlcmlhOgogICAgICAgICAgLSAiPD04MDAiCiAgICB3ZWlnaHQ6IDEKdG90YWxfc2NvcmU6CiAgcGFzczogIjkwJSIKICB3YXJuaW5nOiAiNzUlIg==',
                timeEnd: '2021-10-04T11:04:33Z',
                timeStart: '2021-10-04T11:02:17Z',
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
              temporaryData: {
                distributor: {
                  subscriptionID: '',
                },
              },
            },
            id: '0c7d757d-0b7b-447c-a801-71d1b7f51784',
            source: 'lighthouse-service',
            specversion: '1.0',
            time: '2021-10-04T11:05:46.915Z',
            type: 'sh.keptn.event.evaluation.finished',
            shkeptncontext: 'eb7d7e4f-f7cf-4a15-80ed-0167235b290e',
            shkeptnspecversion: '0.2.3',
            triggeredid: '4d6bea6d-4432-4573-959b-0b2141565e41',
          },
          heatmapLabel: '2021-10-04 13:05',
        },
        color: '#dc172a',
        name: '2021-10-04 13:05',
        label: '2021-10-04 13:05',
        x: 4,
      },
      {
        y: 100,
        evaluationData: {
          traces: [],
          data: {
            evaluation: {
              comparedEvents: ['27c733aa-9bb0-44b4-9574-ef04e38eb4c4'],
              indicatorResults: [
                {
                  displayName: 'Response time P95',
                  keySli: false,
                  passTargets: [
                    {
                      criteria: '<=+10%',
                      targetValue: 25.43676878676464,
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
                    value: 9.280349173333333,
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
                'LS0tCnNwZWNfdmVyc2lvbjogIjEuMCIKY29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciCiAgY29tcGFyZV93aXRoOiAic2luZ2xlX3Jlc3VsdCIKICBpbmNsdWRlX3Jlc3VsdF93aXRoX3Njb3JlOiAicGFzcyIKICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxCmZpbHRlcjoKb2JqZWN0aXZlczoKICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1IgogICAgZGlzcGxheU5hbWU6ICJSZXNwb25zZSB0aW1lIFA5NSIKICAgIGtleV9zbGk6IGZhbHNlCiAgICBwYXNzOiAgICAgICAgICAgICAjIHBhc3MgaWYgKHJlbGF0aXZlIGNoYW5nZSA8PSAxMCUgQU5EIGFic29sdXRlIHZhbHVlIGlzIDwgNjAwbXMpCiAgICAgIC0gY3JpdGVyaWE6CiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykKICAgICAgICAgIC0gIjw2MDAiICAgICMgYWJzb2x1dGUgdmFsdWVzIG9ubHkgcmVxdWlyZSBhIGxvZ2ljYWwgb3BlcmF0b3IKICAgIHdhcm5pbmc6ICAgICAgICAgICMgaWYgdGhlIHJlc3BvbnNlIHRpbWUgaXMgYmVsb3cgODAwbXMsIHRoZSByZXN1bHQgc2hvdWxkIGJlIGEgd2FybmluZwogICAgICAtIGNyaXRlcmlhOgogICAgICAgICAgLSAiPD04MDAiCiAgICB3ZWlnaHQ6IDEKdG90YWxfc2NvcmU6CiAgcGFzczogIjkwJSIKICB3YXJuaW5nOiAiNzUlIg==',
              timeEnd: '2021-10-07T15:01:00.582Z',
              timeStart: '2021-10-07T14:56:00.582Z',
              sloFileContentParsed:
                '---\nspec_version: "1.0"\ncomparison:\n  aggregate_function: "avg"\n  compare_with: "single_result"\n  include_result_with_score: "pass"\n  number_of_comparison_results: 1\nfilter:\nobjectives:\n  - sli: "response_time_p95"\n    displayName: "Response time P95"\n    key_sli: false\n    pass:             # pass if (relative change <= 10% AND absolute value is < 600ms)\n      - criteria:\n          - "<=+10%"  # relative values require a prefixed sign (plus or minus)\n          - "<600"    # absolute values only require a logical operator\n    warning:          # if the response time is below 800ms, the result should be a warning\n      - criteria:\n          - "<=800"\n    weight: 1\ntotal_score:\n  pass: "90%"\n  warning: "75%"',
              score_pass: '90',
              score_warning: '75',
              compare_with: 'single_result\n',
              include_result_with_score: 'pass\n',
              number_of_comparison_results: 1,
            },
            labels: {
              DtCreds: 'dynatrace',
            },
            project: 'sockshop',
            result: 'pass',
            service: 'carts',
            stage: 'staging',
            status: 'succeeded',
            temporaryData: {
              distributor: {
                subscriptionID: '',
              },
            },
          },
          id: '53359c6f-fa31-4e8b-9fde-e003b3ea57ec',
          source: 'lighthouse-service',
          specversion: '1.0',
          time: '2021-10-07T15:01:24.717Z',
          type: 'sh.keptn.event.evaluation.finished',
          shkeptncontext: '68239e0e-e338-483a-be72-9f47cca3eaeb',
          shkeptnspecversion: '0.2.3',
          triggeredid: 'f0b1c0e1-3e33-48fe-b626-b52e508e9150',
          plainEvent: {
            data: {
              evaluation: {
                comparedEvents: ['27c733aa-9bb0-44b4-9574-ef04e38eb4c4'],
                indicatorResults: [
                  {
                    displayName: 'Response time P95',
                    keySli: false,
                    passTargets: [
                      {
                        criteria: '<=+10%',
                        targetValue: 25.43676878676464,
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
                      value: 9.280349173333333,
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
                  'LS0tCnNwZWNfdmVyc2lvbjogIjEuMCIKY29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciCiAgY29tcGFyZV93aXRoOiAic2luZ2xlX3Jlc3VsdCIKICBpbmNsdWRlX3Jlc3VsdF93aXRoX3Njb3JlOiAicGFzcyIKICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxCmZpbHRlcjoKb2JqZWN0aXZlczoKICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1IgogICAgZGlzcGxheU5hbWU6ICJSZXNwb25zZSB0aW1lIFA5NSIKICAgIGtleV9zbGk6IGZhbHNlCiAgICBwYXNzOiAgICAgICAgICAgICAjIHBhc3MgaWYgKHJlbGF0aXZlIGNoYW5nZSA8PSAxMCUgQU5EIGFic29sdXRlIHZhbHVlIGlzIDwgNjAwbXMpCiAgICAgIC0gY3JpdGVyaWE6CiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykKICAgICAgICAgIC0gIjw2MDAiICAgICMgYWJzb2x1dGUgdmFsdWVzIG9ubHkgcmVxdWlyZSBhIGxvZ2ljYWwgb3BlcmF0b3IKICAgIHdhcm5pbmc6ICAgICAgICAgICMgaWYgdGhlIHJlc3BvbnNlIHRpbWUgaXMgYmVsb3cgODAwbXMsIHRoZSByZXN1bHQgc2hvdWxkIGJlIGEgd2FybmluZwogICAgICAtIGNyaXRlcmlhOgogICAgICAgICAgLSAiPD04MDAiCiAgICB3ZWlnaHQ6IDEKdG90YWxfc2NvcmU6CiAgcGFzczogIjkwJSIKICB3YXJuaW5nOiAiNzUlIg==',
                timeEnd: '2021-10-07T15:01:00.582Z',
                timeStart: '2021-10-07T14:56:00.582Z',
              },
              labels: {
                DtCreds: 'dynatrace',
              },
              project: 'sockshop',
              result: 'pass',
              service: 'carts',
              stage: 'staging',
              status: 'succeeded',
              temporaryData: {
                distributor: {
                  subscriptionID: '',
                },
              },
            },
            id: '53359c6f-fa31-4e8b-9fde-e003b3ea57ec',
            source: 'lighthouse-service',
            specversion: '1.0',
            time: '2021-10-07T15:01:24.717Z',
            type: 'sh.keptn.event.evaluation.finished',
            shkeptncontext: '68239e0e-e338-483a-be72-9f47cca3eaeb',
            shkeptnspecversion: '0.2.3',
            triggeredid: 'f0b1c0e1-3e33-48fe-b626-b52e508e9150',
          },
          heatmapLabel: '2021-10-07 17:01',
        },
        color: '#7dc540',
        name: '2021-10-07 17:01',
        label: '2021-10-07 17:01',
        x: 5,
      },
      {
        y: 100,
        evaluationData: {
          traces: [],
          data: {
            evaluation: {
              comparedEvents: ['53359c6f-fa31-4e8b-9fde-e003b3ea57ec'],
              indicatorResults: [
                {
                  displayName: 'Response time P95',
                  keySli: false,
                  passTargets: [
                    {
                      criteria: '<=+10%',
                      targetValue: 10.208384090666666,
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
                    value: 9.329947741347496,
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
                'LS0tCnNwZWNfdmVyc2lvbjogIjEuMCIKY29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciCiAgY29tcGFyZV93aXRoOiAic2luZ2xlX3Jlc3VsdCIKICBpbmNsdWRlX3Jlc3VsdF93aXRoX3Njb3JlOiAicGFzcyIKICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxCmZpbHRlcjoKb2JqZWN0aXZlczoKICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1IgogICAgZGlzcGxheU5hbWU6ICJSZXNwb25zZSB0aW1lIFA5NSIKICAgIGtleV9zbGk6IGZhbHNlCiAgICBwYXNzOiAgICAgICAgICAgICAjIHBhc3MgaWYgKHJlbGF0aXZlIGNoYW5nZSA8PSAxMCUgQU5EIGFic29sdXRlIHZhbHVlIGlzIDwgNjAwbXMpCiAgICAgIC0gY3JpdGVyaWE6CiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykKICAgICAgICAgIC0gIjw2MDAiICAgICMgYWJzb2x1dGUgdmFsdWVzIG9ubHkgcmVxdWlyZSBhIGxvZ2ljYWwgb3BlcmF0b3IKICAgIHdhcm5pbmc6ICAgICAgICAgICMgaWYgdGhlIHJlc3BvbnNlIHRpbWUgaXMgYmVsb3cgODAwbXMsIHRoZSByZXN1bHQgc2hvdWxkIGJlIGEgd2FybmluZwogICAgICAtIGNyaXRlcmlhOgogICAgICAgICAgLSAiPD04MDAiCiAgICB3ZWlnaHQ6IDEKdG90YWxfc2NvcmU6CiAgcGFzczogIjkwJSIKICB3YXJuaW5nOiAiNzUlIg==',
              timeEnd: '2021-10-07T15:02:02.371Z',
              timeStart: '2021-10-07T14:57:02.371Z',
              sloFileContentParsed:
                '---\nspec_version: "1.0"\ncomparison:\n  aggregate_function: "avg"\n  compare_with: "single_result"\n  include_result_with_score: "pass"\n  number_of_comparison_results: 1\nfilter:\nobjectives:\n  - sli: "response_time_p95"\n    displayName: "Response time P95"\n    key_sli: false\n    pass:             # pass if (relative change <= 10% AND absolute value is < 600ms)\n      - criteria:\n          - "<=+10%"  # relative values require a prefixed sign (plus or minus)\n          - "<600"    # absolute values only require a logical operator\n    warning:          # if the response time is below 800ms, the result should be a warning\n      - criteria:\n          - "<=800"\n    weight: 1\ntotal_score:\n  pass: "90%"\n  warning: "75%"',
              score_pass: '90',
              score_warning: '75',
              compare_with: 'single_result\n',
              include_result_with_score: 'pass\n',
              number_of_comparison_results: 1,
            },
            labels: {
              DtCreds: 'dynatrace',
            },
            project: 'sockshop',
            result: 'pass',
            service: 'carts',
            stage: 'staging',
            status: 'succeeded',
            temporaryData: {
              distributor: {
                subscriptionID: '',
              },
            },
          },
          id: 'd7fba175-96b1-4a22-943a-78d871971925',
          source: 'lighthouse-service',
          specversion: '1.0',
          time: '2021-10-07T15:02:27.215Z',
          type: 'sh.keptn.event.evaluation.finished',
          shkeptncontext: '6f0636f7-8f01-4c2a-99b1-1cbf9ff8bfec',
          shkeptnspecversion: '0.2.3',
          triggeredid: '1ee69f43-4a51-46c1-a454-55c925ec4519',
          plainEvent: {
            data: {
              evaluation: {
                comparedEvents: ['53359c6f-fa31-4e8b-9fde-e003b3ea57ec'],
                indicatorResults: [
                  {
                    displayName: 'Response time P95',
                    keySli: false,
                    passTargets: [
                      {
                        criteria: '<=+10%',
                        targetValue: 10.208384090666666,
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
                      value: 9.329947741347496,
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
                  'LS0tCnNwZWNfdmVyc2lvbjogIjEuMCIKY29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciCiAgY29tcGFyZV93aXRoOiAic2luZ2xlX3Jlc3VsdCIKICBpbmNsdWRlX3Jlc3VsdF93aXRoX3Njb3JlOiAicGFzcyIKICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxCmZpbHRlcjoKb2JqZWN0aXZlczoKICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1IgogICAgZGlzcGxheU5hbWU6ICJSZXNwb25zZSB0aW1lIFA5NSIKICAgIGtleV9zbGk6IGZhbHNlCiAgICBwYXNzOiAgICAgICAgICAgICAjIHBhc3MgaWYgKHJlbGF0aXZlIGNoYW5nZSA8PSAxMCUgQU5EIGFic29sdXRlIHZhbHVlIGlzIDwgNjAwbXMpCiAgICAgIC0gY3JpdGVyaWE6CiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykKICAgICAgICAgIC0gIjw2MDAiICAgICMgYWJzb2x1dGUgdmFsdWVzIG9ubHkgcmVxdWlyZSBhIGxvZ2ljYWwgb3BlcmF0b3IKICAgIHdhcm5pbmc6ICAgICAgICAgICMgaWYgdGhlIHJlc3BvbnNlIHRpbWUgaXMgYmVsb3cgODAwbXMsIHRoZSByZXN1bHQgc2hvdWxkIGJlIGEgd2FybmluZwogICAgICAtIGNyaXRlcmlhOgogICAgICAgICAgLSAiPD04MDAiCiAgICB3ZWlnaHQ6IDEKdG90YWxfc2NvcmU6CiAgcGFzczogIjkwJSIKICB3YXJuaW5nOiAiNzUlIg==',
                timeEnd: '2021-10-07T15:02:02.371Z',
                timeStart: '2021-10-07T14:57:02.371Z',
              },
              labels: {
                DtCreds: 'dynatrace',
              },
              project: 'sockshop',
              result: 'pass',
              service: 'carts',
              stage: 'staging',
              status: 'succeeded',
              temporaryData: {
                distributor: {
                  subscriptionID: '',
                },
              },
            },
            id: 'd7fba175-96b1-4a22-943a-78d871971925',
            source: 'lighthouse-service',
            specversion: '1.0',
            time: '2021-10-07T15:02:27.215Z',
            type: 'sh.keptn.event.evaluation.finished',
            shkeptncontext: '6f0636f7-8f01-4c2a-99b1-1cbf9ff8bfec',
            shkeptnspecversion: '0.2.3',
            triggeredid: '1ee69f43-4a51-46c1-a454-55c925ec4519',
          },
          heatmapLabel: '2021-10-07 17:02',
        },
        color: '#7dc540',
        name: '2021-10-07 17:02',
        label: '2021-10-07 17:02',
        x: 6,
      },
    ],
    cursor: 'pointer',
    visible: false,
    turboThreshold: 0,
  },
  {
    metricName: 'response_time_p95',
    name: 'Response time P95',
    type: 'line',
    yAxis: 1,
    data: [
      {
        y: 298.9114676593789,
        indicatorResult: {
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
            value: 298.9114676593789,
          },
          warningTargets: [
            {
              criteria: '<=800',
              targetValue: 800,
              violated: false,
            },
          ],
        },
        evaluationData: {
          traces: [],
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
                    value: 298.9114676593789,
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
                'LS0tCnNwZWNfdmVyc2lvbjogIjEuMCIKY29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciCiAgY29tcGFyZV93aXRoOiAic2luZ2xlX3Jlc3VsdCIKICBpbmNsdWRlX3Jlc3VsdF93aXRoX3Njb3JlOiAicGFzcyIKICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxCmZpbHRlcjoKb2JqZWN0aXZlczoKICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1IgogICAgZGlzcGxheU5hbWU6ICJSZXNwb25zZSB0aW1lIFA5NSIKICAgIGtleV9zbGk6IGZhbHNlCiAgICBwYXNzOiAgICAgICAgICAgICAjIHBhc3MgaWYgKHJlbGF0aXZlIGNoYW5nZSA8PSAxMCUgQU5EIGFic29sdXRlIHZhbHVlIGlzIDwgNjAwbXMpCiAgICAgIC0gY3JpdGVyaWE6CiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykKICAgICAgICAgIC0gIjw2MDAiICAgICMgYWJzb2x1dGUgdmFsdWVzIG9ubHkgcmVxdWlyZSBhIGxvZ2ljYWwgb3BlcmF0b3IKICAgIHdhcm5pbmc6ICAgICAgICAgICMgaWYgdGhlIHJlc3BvbnNlIHRpbWUgaXMgYmVsb3cgODAwbXMsIHRoZSByZXN1bHQgc2hvdWxkIGJlIGEgd2FybmluZwogICAgICAtIGNyaXRlcmlhOgogICAgICAgICAgLSAiPD04MDAiCiAgICB3ZWlnaHQ6IDEKdG90YWxfc2NvcmU6CiAgcGFzczogIjkwJSIKICB3YXJuaW5nOiAiNzUlIg==',
              timeEnd: '2021-09-30T11:41:15Z',
              timeStart: '2021-09-30T11:39:08Z',
              sloFileContentParsed:
                '---\nspec_version: "1.0"\ncomparison:\n  aggregate_function: "avg"\n  compare_with: "single_result"\n  include_result_with_score: "pass"\n  number_of_comparison_results: 1\nfilter:\nobjectives:\n  - sli: "response_time_p95"\n    displayName: "Response time P95"\n    key_sli: false\n    pass:             # pass if (relative change <= 10% AND absolute value is < 600ms)\n      - criteria:\n          - "<=+10%"  # relative values require a prefixed sign (plus or minus)\n          - "<600"    # absolute values only require a logical operator\n    warning:          # if the response time is below 800ms, the result should be a warning\n      - criteria:\n          - "<=800"\n    weight: 1\ntotal_score:\n  pass: "90%"\n  warning: "75%"',
              score_pass: '90',
              score_warning: '75',
              compare_with: 'single_result\n',
              include_result_with_score: 'pass\n',
              number_of_comparison_results: 0,
            },
            labels: {
              DtCreds: 'dynatrace',
            },
            project: 'sockshop',
            result: 'pass',
            service: 'carts',
            stage: 'staging',
            status: 'succeeded',
            temporaryData: {
              distributor: {
                subscriptionID: '',
              },
            },
          },
          id: '40cfcfc2-9cc6-4fce-ba62-5f5fec3e3e3b',
          source: 'lighthouse-service',
          specversion: '1.0',
          time: '2021-09-30T11:42:29.616Z',
          type: 'sh.keptn.event.evaluation.finished',
          shkeptncontext: '369b15c8-5410-49de-87f7-08a2f9f16243',
          shkeptnspecversion: '0.2.3',
          triggeredid: '2d2e799f-6c25-4b26-83ab-12953335f6d4',
          plainEvent: {
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
                      value: 298.9114676593789,
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
                  'LS0tCnNwZWNfdmVyc2lvbjogIjEuMCIKY29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciCiAgY29tcGFyZV93aXRoOiAic2luZ2xlX3Jlc3VsdCIKICBpbmNsdWRlX3Jlc3VsdF93aXRoX3Njb3JlOiAicGFzcyIKICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxCmZpbHRlcjoKb2JqZWN0aXZlczoKICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1IgogICAgZGlzcGxheU5hbWU6ICJSZXNwb25zZSB0aW1lIFA5NSIKICAgIGtleV9zbGk6IGZhbHNlCiAgICBwYXNzOiAgICAgICAgICAgICAjIHBhc3MgaWYgKHJlbGF0aXZlIGNoYW5nZSA8PSAxMCUgQU5EIGFic29sdXRlIHZhbHVlIGlzIDwgNjAwbXMpCiAgICAgIC0gY3JpdGVyaWE6CiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykKICAgICAgICAgIC0gIjw2MDAiICAgICMgYWJzb2x1dGUgdmFsdWVzIG9ubHkgcmVxdWlyZSBhIGxvZ2ljYWwgb3BlcmF0b3IKICAgIHdhcm5pbmc6ICAgICAgICAgICMgaWYgdGhlIHJlc3BvbnNlIHRpbWUgaXMgYmVsb3cgODAwbXMsIHRoZSByZXN1bHQgc2hvdWxkIGJlIGEgd2FybmluZwogICAgICAtIGNyaXRlcmlhOgogICAgICAgICAgLSAiPD04MDAiCiAgICB3ZWlnaHQ6IDEKdG90YWxfc2NvcmU6CiAgcGFzczogIjkwJSIKICB3YXJuaW5nOiAiNzUlIg==',
                timeEnd: '2021-09-30T11:41:15Z',
                timeStart: '2021-09-30T11:39:08Z',
              },
              labels: {
                DtCreds: 'dynatrace',
              },
              project: 'sockshop',
              result: 'pass',
              service: 'carts',
              stage: 'staging',
              status: 'succeeded',
              temporaryData: {
                distributor: {
                  subscriptionID: '',
                },
              },
            },
            id: '40cfcfc2-9cc6-4fce-ba62-5f5fec3e3e3b',
            source: 'lighthouse-service',
            specversion: '1.0',
            time: '2021-09-30T11:42:29.616Z',
            type: 'sh.keptn.event.evaluation.finished',
            shkeptncontext: '369b15c8-5410-49de-87f7-08a2f9f16243',
            shkeptnspecversion: '0.2.3',
            triggeredid: '2d2e799f-6c25-4b26-83ab-12953335f6d4',
          },
          heatmapLabel: '2021-09-30 13:42',
        },
        name: '2021-09-30 13:42',
        x: 0,
      },
      {
        y: 1091.8809833281025,
        indicatorResult: {
          displayName: 'Response time P95',
          keySli: false,
          passTargets: [
            {
              criteria: '<=+10%',
              targetValue: 328.8026144253168,
              violated: true,
            },
            {
              criteria: '<600',
              targetValue: 600,
              violated: true,
            },
          ],
          score: 0,
          status: 'fail',
          value: {
            metric: 'response_time_p95',
            success: true,
            value: 1091.8809833281025,
          },
          warningTargets: [
            {
              criteria: '<=800',
              targetValue: 800,
              violated: true,
            },
          ],
        },
        evaluationData: {
          traces: [],
          data: {
            evaluation: {
              comparedEvents: ['40cfcfc2-9cc6-4fce-ba62-5f5fec3e3e3b'],
              indicatorResults: [
                {
                  displayName: 'Response time P95',
                  keySli: false,
                  passTargets: [
                    {
                      criteria: '<=+10%',
                      targetValue: 328.8026144253168,
                      violated: true,
                    },
                    {
                      criteria: '<600',
                      targetValue: 600,
                      violated: true,
                    },
                  ],
                  score: 0,
                  status: 'fail',
                  value: {
                    metric: 'response_time_p95',
                    success: true,
                    value: 1091.8809833281025,
                  },
                  warningTargets: [
                    {
                      criteria: '<=800',
                      targetValue: 800,
                      violated: true,
                    },
                  ],
                },
              ],
              result: 'fail',
              score: 0,
              sloFileContent:
                'LS0tCnNwZWNfdmVyc2lvbjogIjEuMCIKY29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciCiAgY29tcGFyZV93aXRoOiAic2luZ2xlX3Jlc3VsdCIKICBpbmNsdWRlX3Jlc3VsdF93aXRoX3Njb3JlOiAicGFzcyIKICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxCmZpbHRlcjoKb2JqZWN0aXZlczoKICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1IgogICAgZGlzcGxheU5hbWU6ICJSZXNwb25zZSB0aW1lIFA5NSIKICAgIGtleV9zbGk6IGZhbHNlCiAgICBwYXNzOiAgICAgICAgICAgICAjIHBhc3MgaWYgKHJlbGF0aXZlIGNoYW5nZSA8PSAxMCUgQU5EIGFic29sdXRlIHZhbHVlIGlzIDwgNjAwbXMpCiAgICAgIC0gY3JpdGVyaWE6CiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykKICAgICAgICAgIC0gIjw2MDAiICAgICMgYWJzb2x1dGUgdmFsdWVzIG9ubHkgcmVxdWlyZSBhIGxvZ2ljYWwgb3BlcmF0b3IKICAgIHdhcm5pbmc6ICAgICAgICAgICMgaWYgdGhlIHJlc3BvbnNlIHRpbWUgaXMgYmVsb3cgODAwbXMsIHRoZSByZXN1bHQgc2hvdWxkIGJlIGEgd2FybmluZwogICAgICAtIGNyaXRlcmlhOgogICAgICAgICAgLSAiPD04MDAiCiAgICB3ZWlnaHQ6IDEKdG90YWxfc2NvcmU6CiAgcGFzczogIjkwJSIKICB3YXJuaW5nOiAiNzUlIg==',
              timeEnd: '2021-10-01T09:25:48Z',
              timeStart: '2021-10-01T09:16:31Z',
              sloFileContentParsed:
                '---\nspec_version: "1.0"\ncomparison:\n  aggregate_function: "avg"\n  compare_with: "single_result"\n  include_result_with_score: "pass"\n  number_of_comparison_results: 1\nfilter:\nobjectives:\n  - sli: "response_time_p95"\n    displayName: "Response time P95"\n    key_sli: false\n    pass:             # pass if (relative change <= 10% AND absolute value is < 600ms)\n      - criteria:\n          - "<=+10%"  # relative values require a prefixed sign (plus or minus)\n          - "<600"    # absolute values only require a logical operator\n    warning:          # if the response time is below 800ms, the result should be a warning\n      - criteria:\n          - "<=800"\n    weight: 1\ntotal_score:\n  pass: "90%"\n  warning: "75%"',
              score_pass: '90',
              score_warning: '75',
              compare_with: 'single_result\n',
              include_result_with_score: 'pass\n',
              number_of_comparison_results: 1,
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
            temporaryData: {
              distributor: {
                subscriptionID: '',
              },
            },
          },
          id: 'f623de99-a716-4b9b-90a2-8df1054e92c7',
          source: 'lighthouse-service',
          specversion: '1.0',
          time: '2021-10-01T09:26:15.608Z',
          type: 'sh.keptn.event.evaluation.finished',
          shkeptncontext: '081ce5a3-e0f9-4c4d-bfcf-5893f85bfaff',
          shkeptnspecversion: '0.2.3',
          triggeredid: 'f564369b-2ee6-4c32-8488-03508ed05b5a',
          plainEvent: {
            data: {
              evaluation: {
                comparedEvents: ['40cfcfc2-9cc6-4fce-ba62-5f5fec3e3e3b'],
                indicatorResults: [
                  {
                    displayName: 'Response time P95',
                    keySli: false,
                    passTargets: [
                      {
                        criteria: '<=+10%',
                        targetValue: 328.8026144253168,
                        violated: true,
                      },
                      {
                        criteria: '<600',
                        targetValue: 600,
                        violated: true,
                      },
                    ],
                    score: 0,
                    status: 'fail',
                    value: {
                      metric: 'response_time_p95',
                      success: true,
                      value: 1091.8809833281025,
                    },
                    warningTargets: [
                      {
                        criteria: '<=800',
                        targetValue: 800,
                        violated: true,
                      },
                    ],
                  },
                ],
                result: 'fail',
                score: 0,
                sloFileContent:
                  'LS0tCnNwZWNfdmVyc2lvbjogIjEuMCIKY29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciCiAgY29tcGFyZV93aXRoOiAic2luZ2xlX3Jlc3VsdCIKICBpbmNsdWRlX3Jlc3VsdF93aXRoX3Njb3JlOiAicGFzcyIKICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxCmZpbHRlcjoKb2JqZWN0aXZlczoKICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1IgogICAgZGlzcGxheU5hbWU6ICJSZXNwb25zZSB0aW1lIFA5NSIKICAgIGtleV9zbGk6IGZhbHNlCiAgICBwYXNzOiAgICAgICAgICAgICAjIHBhc3MgaWYgKHJlbGF0aXZlIGNoYW5nZSA8PSAxMCUgQU5EIGFic29sdXRlIHZhbHVlIGlzIDwgNjAwbXMpCiAgICAgIC0gY3JpdGVyaWE6CiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykKICAgICAgICAgIC0gIjw2MDAiICAgICMgYWJzb2x1dGUgdmFsdWVzIG9ubHkgcmVxdWlyZSBhIGxvZ2ljYWwgb3BlcmF0b3IKICAgIHdhcm5pbmc6ICAgICAgICAgICMgaWYgdGhlIHJlc3BvbnNlIHRpbWUgaXMgYmVsb3cgODAwbXMsIHRoZSByZXN1bHQgc2hvdWxkIGJlIGEgd2FybmluZwogICAgICAtIGNyaXRlcmlhOgogICAgICAgICAgLSAiPD04MDAiCiAgICB3ZWlnaHQ6IDEKdG90YWxfc2NvcmU6CiAgcGFzczogIjkwJSIKICB3YXJuaW5nOiAiNzUlIg==',
                timeEnd: '2021-10-01T09:25:48Z',
                timeStart: '2021-10-01T09:16:31Z',
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
              temporaryData: {
                distributor: {
                  subscriptionID: '',
                },
              },
            },
            id: 'f623de99-a716-4b9b-90a2-8df1054e92c7',
            source: 'lighthouse-service',
            specversion: '1.0',
            time: '2021-10-01T09:26:15.608Z',
            type: 'sh.keptn.event.evaluation.finished',
            shkeptncontext: '081ce5a3-e0f9-4c4d-bfcf-5893f85bfaff',
            shkeptnspecversion: '0.2.3',
            triggeredid: 'f564369b-2ee6-4c32-8488-03508ed05b5a',
          },
          heatmapLabel: '2021-10-01 11:26',
        },
        name: '2021-10-01 11:26',
        x: 1,
      },
      {
        y: 306.3587409268154,
        indicatorResult: {
          displayName: 'Response time P95',
          keySli: false,
          passTargets: [
            {
              criteria: '<=+10%',
              targetValue: 328.8026144253168,
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
            value: 306.3587409268154,
          },
          warningTargets: [
            {
              criteria: '<=800',
              targetValue: 800,
              violated: false,
            },
          ],
        },
        evaluationData: {
          traces: [],
          data: {
            evaluation: {
              comparedEvents: ['40cfcfc2-9cc6-4fce-ba62-5f5fec3e3e3b'],
              indicatorResults: [
                {
                  displayName: 'Response time P95',
                  keySli: false,
                  passTargets: [
                    {
                      criteria: '<=+10%',
                      targetValue: 328.8026144253168,
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
                    value: 306.3587409268154,
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
                'LS0tCnNwZWNfdmVyc2lvbjogIjEuMCIKY29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciCiAgY29tcGFyZV93aXRoOiAic2luZ2xlX3Jlc3VsdCIKICBpbmNsdWRlX3Jlc3VsdF93aXRoX3Njb3JlOiAicGFzcyIKICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxCmZpbHRlcjoKb2JqZWN0aXZlczoKICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1IgogICAgZGlzcGxheU5hbWU6ICJSZXNwb25zZSB0aW1lIFA5NSIKICAgIGtleV9zbGk6IGZhbHNlCiAgICBwYXNzOiAgICAgICAgICAgICAjIHBhc3MgaWYgKHJlbGF0aXZlIGNoYW5nZSA8PSAxMCUgQU5EIGFic29sdXRlIHZhbHVlIGlzIDwgNjAwbXMpCiAgICAgIC0gY3JpdGVyaWE6CiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykKICAgICAgICAgIC0gIjw2MDAiICAgICMgYWJzb2x1dGUgdmFsdWVzIG9ubHkgcmVxdWlyZSBhIGxvZ2ljYWwgb3BlcmF0b3IKICAgIHdhcm5pbmc6ICAgICAgICAgICMgaWYgdGhlIHJlc3BvbnNlIHRpbWUgaXMgYmVsb3cgODAwbXMsIHRoZSByZXN1bHQgc2hvdWxkIGJlIGEgd2FybmluZwogICAgICAtIGNyaXRlcmlhOgogICAgICAgICAgLSAiPD04MDAiCiAgICB3ZWlnaHQ6IDEKdG90YWxfc2NvcmU6CiAgcGFzczogIjkwJSIKICB3YXJuaW5nOiAiNzUlIg==',
              timeEnd: '2021-10-01T10:26:34Z',
              timeStart: '2021-10-01T10:24:15Z',
              sloFileContentParsed:
                '---\nspec_version: "1.0"\ncomparison:\n  aggregate_function: "avg"\n  compare_with: "single_result"\n  include_result_with_score: "pass"\n  number_of_comparison_results: 1\nfilter:\nobjectives:\n  - sli: "response_time_p95"\n    displayName: "Response time P95"\n    key_sli: false\n    pass:             # pass if (relative change <= 10% AND absolute value is < 600ms)\n      - criteria:\n          - "<=+10%"  # relative values require a prefixed sign (plus or minus)\n          - "<600"    # absolute values only require a logical operator\n    warning:          # if the response time is below 800ms, the result should be a warning\n      - criteria:\n          - "<=800"\n    weight: 1\ntotal_score:\n  pass: "90%"\n  warning: "75%"',
              score_pass: '90',
              score_warning: '75',
              compare_with: 'single_result\n',
              include_result_with_score: 'pass\n',
              number_of_comparison_results: 1,
            },
            labels: {
              DtCreds: 'dynatrace',
            },
            project: 'sockshop',
            result: 'pass',
            service: 'carts',
            stage: 'staging',
            status: 'succeeded',
            temporaryData: {
              distributor: {
                subscriptionID: '',
              },
            },
          },
          id: '3344487d-e384-4cd9-a0e0-fcf157a33ad6',
          source: 'lighthouse-service',
          specversion: '1.0',
          time: '2021-10-01T10:27:49.318Z',
          type: 'sh.keptn.event.evaluation.finished',
          shkeptncontext: 'e9dd1b2f-8383-48a9-a2d4-36cb8ff370c7',
          shkeptnspecversion: '0.2.3',
          triggeredid: '61de0ee6-b535-42bb-ac30-d0f359298bb0',
          plainEvent: {
            data: {
              evaluation: {
                comparedEvents: ['40cfcfc2-9cc6-4fce-ba62-5f5fec3e3e3b'],
                indicatorResults: [
                  {
                    displayName: 'Response time P95',
                    keySli: false,
                    passTargets: [
                      {
                        criteria: '<=+10%',
                        targetValue: 328.8026144253168,
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
                      value: 306.3587409268154,
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
                  'LS0tCnNwZWNfdmVyc2lvbjogIjEuMCIKY29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciCiAgY29tcGFyZV93aXRoOiAic2luZ2xlX3Jlc3VsdCIKICBpbmNsdWRlX3Jlc3VsdF93aXRoX3Njb3JlOiAicGFzcyIKICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxCmZpbHRlcjoKb2JqZWN0aXZlczoKICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1IgogICAgZGlzcGxheU5hbWU6ICJSZXNwb25zZSB0aW1lIFA5NSIKICAgIGtleV9zbGk6IGZhbHNlCiAgICBwYXNzOiAgICAgICAgICAgICAjIHBhc3MgaWYgKHJlbGF0aXZlIGNoYW5nZSA8PSAxMCUgQU5EIGFic29sdXRlIHZhbHVlIGlzIDwgNjAwbXMpCiAgICAgIC0gY3JpdGVyaWE6CiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykKICAgICAgICAgIC0gIjw2MDAiICAgICMgYWJzb2x1dGUgdmFsdWVzIG9ubHkgcmVxdWlyZSBhIGxvZ2ljYWwgb3BlcmF0b3IKICAgIHdhcm5pbmc6ICAgICAgICAgICMgaWYgdGhlIHJlc3BvbnNlIHRpbWUgaXMgYmVsb3cgODAwbXMsIHRoZSByZXN1bHQgc2hvdWxkIGJlIGEgd2FybmluZwogICAgICAtIGNyaXRlcmlhOgogICAgICAgICAgLSAiPD04MDAiCiAgICB3ZWlnaHQ6IDEKdG90YWxfc2NvcmU6CiAgcGFzczogIjkwJSIKICB3YXJuaW5nOiAiNzUlIg==',
                timeEnd: '2021-10-01T10:26:34Z',
                timeStart: '2021-10-01T10:24:15Z',
              },
              labels: {
                DtCreds: 'dynatrace',
              },
              project: 'sockshop',
              result: 'pass',
              service: 'carts',
              stage: 'staging',
              status: 'succeeded',
              temporaryData: {
                distributor: {
                  subscriptionID: '',
                },
              },
            },
            id: '3344487d-e384-4cd9-a0e0-fcf157a33ad6',
            source: 'lighthouse-service',
            specversion: '1.0',
            time: '2021-10-01T10:27:49.318Z',
            type: 'sh.keptn.event.evaluation.finished',
            shkeptncontext: 'e9dd1b2f-8383-48a9-a2d4-36cb8ff370c7',
            shkeptnspecversion: '0.2.3',
            triggeredid: '61de0ee6-b535-42bb-ac30-d0f359298bb0',
          },
          heatmapLabel: '2021-10-01 12:27',
        },
        name: '2021-10-01 12:27',
        x: 2,
      },
      {
        y: 23.12433526069513,
        indicatorResult: {
          displayName: 'Response time P95',
          keySli: false,
          passTargets: [
            {
              criteria: '<=+10%',
              targetValue: 336.9946150194969,
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
            value: 23.12433526069513,
          },
          warningTargets: [
            {
              criteria: '<=800',
              targetValue: 800,
              violated: false,
            },
          ],
        },
        evaluationData: {
          traces: [],
          data: {
            evaluation: {
              comparedEvents: ['3344487d-e384-4cd9-a0e0-fcf157a33ad6'],
              indicatorResults: [
                {
                  displayName: 'Response time P95',
                  keySli: false,
                  passTargets: [
                    {
                      criteria: '<=+10%',
                      targetValue: 336.9946150194969,
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
                    value: 23.12433526069513,
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
                'LS0tCnNwZWNfdmVyc2lvbjogIjEuMCIKY29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciCiAgY29tcGFyZV93aXRoOiAic2luZ2xlX3Jlc3VsdCIKICBpbmNsdWRlX3Jlc3VsdF93aXRoX3Njb3JlOiAicGFzcyIKICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxCmZpbHRlcjoKb2JqZWN0aXZlczoKICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1IgogICAgZGlzcGxheU5hbWU6ICJSZXNwb25zZSB0aW1lIFA5NSIKICAgIGtleV9zbGk6IGZhbHNlCiAgICBwYXNzOiAgICAgICAgICAgICAjIHBhc3MgaWYgKHJlbGF0aXZlIGNoYW5nZSA8PSAxMCUgQU5EIGFic29sdXRlIHZhbHVlIGlzIDwgNjAwbXMpCiAgICAgIC0gY3JpdGVyaWE6CiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykKICAgICAgICAgIC0gIjw2MDAiICAgICMgYWJzb2x1dGUgdmFsdWVzIG9ubHkgcmVxdWlyZSBhIGxvZ2ljYWwgb3BlcmF0b3IKICAgIHdhcm5pbmc6ICAgICAgICAgICMgaWYgdGhlIHJlc3BvbnNlIHRpbWUgaXMgYmVsb3cgODAwbXMsIHRoZSByZXN1bHQgc2hvdWxkIGJlIGEgd2FybmluZwogICAgICAtIGNyaXRlcmlhOgogICAgICAgICAgLSAiPD04MDAiCiAgICB3ZWlnaHQ6IDEKdG90YWxfc2NvcmU6CiAgcGFzczogIjkwJSIKICB3YXJuaW5nOiAiNzUlIg==',
              timeEnd: '2021-10-01T10:32:47.372Z',
              timeStart: '2021-10-01T10:27:47.372Z',
              sloFileContentParsed:
                '---\nspec_version: "1.0"\ncomparison:\n  aggregate_function: "avg"\n  compare_with: "single_result"\n  include_result_with_score: "pass"\n  number_of_comparison_results: 1\nfilter:\nobjectives:\n  - sli: "response_time_p95"\n    displayName: "Response time P95"\n    key_sli: false\n    pass:             # pass if (relative change <= 10% AND absolute value is < 600ms)\n      - criteria:\n          - "<=+10%"  # relative values require a prefixed sign (plus or minus)\n          - "<600"    # absolute values only require a logical operator\n    warning:          # if the response time is below 800ms, the result should be a warning\n      - criteria:\n          - "<=800"\n    weight: 1\ntotal_score:\n  pass: "90%"\n  warning: "75%"',
              score_pass: '90',
              score_warning: '75',
              compare_with: 'single_result\n',
              include_result_with_score: 'pass\n',
              number_of_comparison_results: 1,
            },
            labels: {
              DtCreds: 'dynatrace',
            },
            project: 'sockshop',
            result: 'pass',
            service: 'carts',
            stage: 'staging',
            status: 'succeeded',
            temporaryData: {
              distributor: {
                subscriptionID: '',
              },
            },
          },
          id: '27c733aa-9bb0-44b4-9574-ef04e38eb4c4',
          source: 'lighthouse-service',
          specversion: '1.0',
          time: '2021-10-01T10:33:19.708Z',
          type: 'sh.keptn.event.evaluation.finished',
          shkeptncontext: '199287a3-b723-479b-8083-8b3ff3482531',
          shkeptnspecversion: '0.2.3',
          triggeredid: '53ac99d5-2828-470a-b067-442562cf3bca',
          plainEvent: {
            data: {
              evaluation: {
                comparedEvents: ['3344487d-e384-4cd9-a0e0-fcf157a33ad6'],
                indicatorResults: [
                  {
                    displayName: 'Response time P95',
                    keySli: false,
                    passTargets: [
                      {
                        criteria: '<=+10%',
                        targetValue: 336.9946150194969,
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
                      value: 23.12433526069513,
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
                  'LS0tCnNwZWNfdmVyc2lvbjogIjEuMCIKY29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciCiAgY29tcGFyZV93aXRoOiAic2luZ2xlX3Jlc3VsdCIKICBpbmNsdWRlX3Jlc3VsdF93aXRoX3Njb3JlOiAicGFzcyIKICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxCmZpbHRlcjoKb2JqZWN0aXZlczoKICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1IgogICAgZGlzcGxheU5hbWU6ICJSZXNwb25zZSB0aW1lIFA5NSIKICAgIGtleV9zbGk6IGZhbHNlCiAgICBwYXNzOiAgICAgICAgICAgICAjIHBhc3MgaWYgKHJlbGF0aXZlIGNoYW5nZSA8PSAxMCUgQU5EIGFic29sdXRlIHZhbHVlIGlzIDwgNjAwbXMpCiAgICAgIC0gY3JpdGVyaWE6CiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykKICAgICAgICAgIC0gIjw2MDAiICAgICMgYWJzb2x1dGUgdmFsdWVzIG9ubHkgcmVxdWlyZSBhIGxvZ2ljYWwgb3BlcmF0b3IKICAgIHdhcm5pbmc6ICAgICAgICAgICMgaWYgdGhlIHJlc3BvbnNlIHRpbWUgaXMgYmVsb3cgODAwbXMsIHRoZSByZXN1bHQgc2hvdWxkIGJlIGEgd2FybmluZwogICAgICAtIGNyaXRlcmlhOgogICAgICAgICAgLSAiPD04MDAiCiAgICB3ZWlnaHQ6IDEKdG90YWxfc2NvcmU6CiAgcGFzczogIjkwJSIKICB3YXJuaW5nOiAiNzUlIg==',
                timeEnd: '2021-10-01T10:32:47.372Z',
                timeStart: '2021-10-01T10:27:47.372Z',
              },
              labels: {
                DtCreds: 'dynatrace',
              },
              project: 'sockshop',
              result: 'pass',
              service: 'carts',
              stage: 'staging',
              status: 'succeeded',
              temporaryData: {
                distributor: {
                  subscriptionID: '',
                },
              },
            },
            id: '27c733aa-9bb0-44b4-9574-ef04e38eb4c4',
            source: 'lighthouse-service',
            specversion: '1.0',
            time: '2021-10-01T10:33:19.708Z',
            type: 'sh.keptn.event.evaluation.finished',
            shkeptncontext: '199287a3-b723-479b-8083-8b3ff3482531',
            shkeptnspecversion: '0.2.3',
            triggeredid: '53ac99d5-2828-470a-b067-442562cf3bca',
          },
          heatmapLabel: '2021-10-01 12:33',
        },
        name: '2021-10-01 12:33',
        x: 3,
      },
      {
        y: 301.7440018752324,
        indicatorResult: {
          displayName: 'Response time P95',
          keySli: false,
          passTargets: [
            {
              criteria: '<=+10%',
              targetValue: 25.43676878676464,
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
            value: 301.7440018752324,
          },
          warningTargets: [
            {
              criteria: '<=800',
              targetValue: 800,
              violated: false,
            },
          ],
        },
        evaluationData: {
          traces: [],
          data: {
            evaluation: {
              comparedEvents: ['27c733aa-9bb0-44b4-9574-ef04e38eb4c4'],
              indicatorResults: [
                {
                  displayName: 'Response time P95',
                  keySli: false,
                  passTargets: [
                    {
                      criteria: '<=+10%',
                      targetValue: 25.43676878676464,
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
                    value: 301.7440018752324,
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
                'LS0tCnNwZWNfdmVyc2lvbjogIjEuMCIKY29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciCiAgY29tcGFyZV93aXRoOiAic2luZ2xlX3Jlc3VsdCIKICBpbmNsdWRlX3Jlc3VsdF93aXRoX3Njb3JlOiAicGFzcyIKICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxCmZpbHRlcjoKb2JqZWN0aXZlczoKICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1IgogICAgZGlzcGxheU5hbWU6ICJSZXNwb25zZSB0aW1lIFA5NSIKICAgIGtleV9zbGk6IGZhbHNlCiAgICBwYXNzOiAgICAgICAgICAgICAjIHBhc3MgaWYgKHJlbGF0aXZlIGNoYW5nZSA8PSAxMCUgQU5EIGFic29sdXRlIHZhbHVlIGlzIDwgNjAwbXMpCiAgICAgIC0gY3JpdGVyaWE6CiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykKICAgICAgICAgIC0gIjw2MDAiICAgICMgYWJzb2x1dGUgdmFsdWVzIG9ubHkgcmVxdWlyZSBhIGxvZ2ljYWwgb3BlcmF0b3IKICAgIHdhcm5pbmc6ICAgICAgICAgICMgaWYgdGhlIHJlc3BvbnNlIHRpbWUgaXMgYmVsb3cgODAwbXMsIHRoZSByZXN1bHQgc2hvdWxkIGJlIGEgd2FybmluZwogICAgICAtIGNyaXRlcmlhOgogICAgICAgICAgLSAiPD04MDAiCiAgICB3ZWlnaHQ6IDEKdG90YWxfc2NvcmU6CiAgcGFzczogIjkwJSIKICB3YXJuaW5nOiAiNzUlIg==',
              timeEnd: '2021-10-04T11:04:33Z',
              timeStart: '2021-10-04T11:02:17Z',
              sloFileContentParsed:
                '---\nspec_version: "1.0"\ncomparison:\n  aggregate_function: "avg"\n  compare_with: "single_result"\n  include_result_with_score: "pass"\n  number_of_comparison_results: 1\nfilter:\nobjectives:\n  - sli: "response_time_p95"\n    displayName: "Response time P95"\n    key_sli: false\n    pass:             # pass if (relative change <= 10% AND absolute value is < 600ms)\n      - criteria:\n          - "<=+10%"  # relative values require a prefixed sign (plus or minus)\n          - "<600"    # absolute values only require a logical operator\n    warning:          # if the response time is below 800ms, the result should be a warning\n      - criteria:\n          - "<=800"\n    weight: 1\ntotal_score:\n  pass: "90%"\n  warning: "75%"',
              score_pass: '90',
              score_warning: '75',
              compare_with: 'single_result\n',
              include_result_with_score: 'pass\n',
              number_of_comparison_results: 1,
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
            temporaryData: {
              distributor: {
                subscriptionID: '',
              },
            },
          },
          id: '0c7d757d-0b7b-447c-a801-71d1b7f51784',
          source: 'lighthouse-service',
          specversion: '1.0',
          time: '2021-10-04T11:05:46.915Z',
          type: 'sh.keptn.event.evaluation.finished',
          shkeptncontext: 'eb7d7e4f-f7cf-4a15-80ed-0167235b290e',
          shkeptnspecversion: '0.2.3',
          triggeredid: '4d6bea6d-4432-4573-959b-0b2141565e41',
          plainEvent: {
            data: {
              evaluation: {
                comparedEvents: ['27c733aa-9bb0-44b4-9574-ef04e38eb4c4'],
                indicatorResults: [
                  {
                    displayName: 'Response time P95',
                    keySli: false,
                    passTargets: [
                      {
                        criteria: '<=+10%',
                        targetValue: 25.43676878676464,
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
                      value: 301.7440018752324,
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
                  'LS0tCnNwZWNfdmVyc2lvbjogIjEuMCIKY29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciCiAgY29tcGFyZV93aXRoOiAic2luZ2xlX3Jlc3VsdCIKICBpbmNsdWRlX3Jlc3VsdF93aXRoX3Njb3JlOiAicGFzcyIKICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxCmZpbHRlcjoKb2JqZWN0aXZlczoKICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1IgogICAgZGlzcGxheU5hbWU6ICJSZXNwb25zZSB0aW1lIFA5NSIKICAgIGtleV9zbGk6IGZhbHNlCiAgICBwYXNzOiAgICAgICAgICAgICAjIHBhc3MgaWYgKHJlbGF0aXZlIGNoYW5nZSA8PSAxMCUgQU5EIGFic29sdXRlIHZhbHVlIGlzIDwgNjAwbXMpCiAgICAgIC0gY3JpdGVyaWE6CiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykKICAgICAgICAgIC0gIjw2MDAiICAgICMgYWJzb2x1dGUgdmFsdWVzIG9ubHkgcmVxdWlyZSBhIGxvZ2ljYWwgb3BlcmF0b3IKICAgIHdhcm5pbmc6ICAgICAgICAgICMgaWYgdGhlIHJlc3BvbnNlIHRpbWUgaXMgYmVsb3cgODAwbXMsIHRoZSByZXN1bHQgc2hvdWxkIGJlIGEgd2FybmluZwogICAgICAtIGNyaXRlcmlhOgogICAgICAgICAgLSAiPD04MDAiCiAgICB3ZWlnaHQ6IDEKdG90YWxfc2NvcmU6CiAgcGFzczogIjkwJSIKICB3YXJuaW5nOiAiNzUlIg==',
                timeEnd: '2021-10-04T11:04:33Z',
                timeStart: '2021-10-04T11:02:17Z',
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
              temporaryData: {
                distributor: {
                  subscriptionID: '',
                },
              },
            },
            id: '0c7d757d-0b7b-447c-a801-71d1b7f51784',
            source: 'lighthouse-service',
            specversion: '1.0',
            time: '2021-10-04T11:05:46.915Z',
            type: 'sh.keptn.event.evaluation.finished',
            shkeptncontext: 'eb7d7e4f-f7cf-4a15-80ed-0167235b290e',
            shkeptnspecversion: '0.2.3',
            triggeredid: '4d6bea6d-4432-4573-959b-0b2141565e41',
          },
          heatmapLabel: '2021-10-04 13:05',
        },
        name: '2021-10-04 13:05',
        x: 4,
      },
      {
        y: 9.280349173333333,
        indicatorResult: {
          displayName: 'Response time P95',
          keySli: false,
          passTargets: [
            {
              criteria: '<=+10%',
              targetValue: 25.43676878676464,
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
            value: 9.280349173333333,
          },
          warningTargets: [
            {
              criteria: '<=800',
              targetValue: 800,
              violated: false,
            },
          ],
        },
        evaluationData: {
          traces: [],
          data: {
            evaluation: {
              comparedEvents: ['27c733aa-9bb0-44b4-9574-ef04e38eb4c4'],
              indicatorResults: [
                {
                  displayName: 'Response time P95',
                  keySli: false,
                  passTargets: [
                    {
                      criteria: '<=+10%',
                      targetValue: 25.43676878676464,
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
                    value: 9.280349173333333,
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
                'LS0tCnNwZWNfdmVyc2lvbjogIjEuMCIKY29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciCiAgY29tcGFyZV93aXRoOiAic2luZ2xlX3Jlc3VsdCIKICBpbmNsdWRlX3Jlc3VsdF93aXRoX3Njb3JlOiAicGFzcyIKICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxCmZpbHRlcjoKb2JqZWN0aXZlczoKICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1IgogICAgZGlzcGxheU5hbWU6ICJSZXNwb25zZSB0aW1lIFA5NSIKICAgIGtleV9zbGk6IGZhbHNlCiAgICBwYXNzOiAgICAgICAgICAgICAjIHBhc3MgaWYgKHJlbGF0aXZlIGNoYW5nZSA8PSAxMCUgQU5EIGFic29sdXRlIHZhbHVlIGlzIDwgNjAwbXMpCiAgICAgIC0gY3JpdGVyaWE6CiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykKICAgICAgICAgIC0gIjw2MDAiICAgICMgYWJzb2x1dGUgdmFsdWVzIG9ubHkgcmVxdWlyZSBhIGxvZ2ljYWwgb3BlcmF0b3IKICAgIHdhcm5pbmc6ICAgICAgICAgICMgaWYgdGhlIHJlc3BvbnNlIHRpbWUgaXMgYmVsb3cgODAwbXMsIHRoZSByZXN1bHQgc2hvdWxkIGJlIGEgd2FybmluZwogICAgICAtIGNyaXRlcmlhOgogICAgICAgICAgLSAiPD04MDAiCiAgICB3ZWlnaHQ6IDEKdG90YWxfc2NvcmU6CiAgcGFzczogIjkwJSIKICB3YXJuaW5nOiAiNzUlIg==',
              timeEnd: '2021-10-07T15:01:00.582Z',
              timeStart: '2021-10-07T14:56:00.582Z',
              sloFileContentParsed:
                '---\nspec_version: "1.0"\ncomparison:\n  aggregate_function: "avg"\n  compare_with: "single_result"\n  include_result_with_score: "pass"\n  number_of_comparison_results: 1\nfilter:\nobjectives:\n  - sli: "response_time_p95"\n    displayName: "Response time P95"\n    key_sli: false\n    pass:             # pass if (relative change <= 10% AND absolute value is < 600ms)\n      - criteria:\n          - "<=+10%"  # relative values require a prefixed sign (plus or minus)\n          - "<600"    # absolute values only require a logical operator\n    warning:          # if the response time is below 800ms, the result should be a warning\n      - criteria:\n          - "<=800"\n    weight: 1\ntotal_score:\n  pass: "90%"\n  warning: "75%"',
              score_pass: '90',
              score_warning: '75',
              compare_with: 'single_result\n',
              include_result_with_score: 'pass\n',
              number_of_comparison_results: 1,
            },
            labels: {
              DtCreds: 'dynatrace',
            },
            project: 'sockshop',
            result: 'pass',
            service: 'carts',
            stage: 'staging',
            status: 'succeeded',
            temporaryData: {
              distributor: {
                subscriptionID: '',
              },
            },
          },
          id: '53359c6f-fa31-4e8b-9fde-e003b3ea57ec',
          source: 'lighthouse-service',
          specversion: '1.0',
          time: '2021-10-07T15:01:24.717Z',
          type: 'sh.keptn.event.evaluation.finished',
          shkeptncontext: '68239e0e-e338-483a-be72-9f47cca3eaeb',
          shkeptnspecversion: '0.2.3',
          triggeredid: 'f0b1c0e1-3e33-48fe-b626-b52e508e9150',
          plainEvent: {
            data: {
              evaluation: {
                comparedEvents: ['27c733aa-9bb0-44b4-9574-ef04e38eb4c4'],
                indicatorResults: [
                  {
                    displayName: 'Response time P95',
                    keySli: false,
                    passTargets: [
                      {
                        criteria: '<=+10%',
                        targetValue: 25.43676878676464,
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
                      value: 9.280349173333333,
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
                  'LS0tCnNwZWNfdmVyc2lvbjogIjEuMCIKY29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciCiAgY29tcGFyZV93aXRoOiAic2luZ2xlX3Jlc3VsdCIKICBpbmNsdWRlX3Jlc3VsdF93aXRoX3Njb3JlOiAicGFzcyIKICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxCmZpbHRlcjoKb2JqZWN0aXZlczoKICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1IgogICAgZGlzcGxheU5hbWU6ICJSZXNwb25zZSB0aW1lIFA5NSIKICAgIGtleV9zbGk6IGZhbHNlCiAgICBwYXNzOiAgICAgICAgICAgICAjIHBhc3MgaWYgKHJlbGF0aXZlIGNoYW5nZSA8PSAxMCUgQU5EIGFic29sdXRlIHZhbHVlIGlzIDwgNjAwbXMpCiAgICAgIC0gY3JpdGVyaWE6CiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykKICAgICAgICAgIC0gIjw2MDAiICAgICMgYWJzb2x1dGUgdmFsdWVzIG9ubHkgcmVxdWlyZSBhIGxvZ2ljYWwgb3BlcmF0b3IKICAgIHdhcm5pbmc6ICAgICAgICAgICMgaWYgdGhlIHJlc3BvbnNlIHRpbWUgaXMgYmVsb3cgODAwbXMsIHRoZSByZXN1bHQgc2hvdWxkIGJlIGEgd2FybmluZwogICAgICAtIGNyaXRlcmlhOgogICAgICAgICAgLSAiPD04MDAiCiAgICB3ZWlnaHQ6IDEKdG90YWxfc2NvcmU6CiAgcGFzczogIjkwJSIKICB3YXJuaW5nOiAiNzUlIg==',
                timeEnd: '2021-10-07T15:01:00.582Z',
                timeStart: '2021-10-07T14:56:00.582Z',
              },
              labels: {
                DtCreds: 'dynatrace',
              },
              project: 'sockshop',
              result: 'pass',
              service: 'carts',
              stage: 'staging',
              status: 'succeeded',
              temporaryData: {
                distributor: {
                  subscriptionID: '',
                },
              },
            },
            id: '53359c6f-fa31-4e8b-9fde-e003b3ea57ec',
            source: 'lighthouse-service',
            specversion: '1.0',
            time: '2021-10-07T15:01:24.717Z',
            type: 'sh.keptn.event.evaluation.finished',
            shkeptncontext: '68239e0e-e338-483a-be72-9f47cca3eaeb',
            shkeptnspecversion: '0.2.3',
            triggeredid: 'f0b1c0e1-3e33-48fe-b626-b52e508e9150',
          },
          heatmapLabel: '2021-10-07 17:01',
        },
        name: '2021-10-07 17:01',
        x: 5,
      },
      {
        y: 9.329947741347496,
        indicatorResult: {
          displayName: 'Response time P95',
          keySli: false,
          passTargets: [
            {
              criteria: '<=+10%',
              targetValue: 10.208384090666666,
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
            value: 9.329947741347496,
          },
          warningTargets: [
            {
              criteria: '<=800',
              targetValue: 800,
              violated: false,
            },
          ],
        },
        evaluationData: {
          traces: [],
          data: {
            evaluation: {
              comparedEvents: ['53359c6f-fa31-4e8b-9fde-e003b3ea57ec'],
              indicatorResults: [
                {
                  displayName: 'Response time P95',
                  keySli: false,
                  passTargets: [
                    {
                      criteria: '<=+10%',
                      targetValue: 10.208384090666666,
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
                    value: 9.329947741347496,
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
                'LS0tCnNwZWNfdmVyc2lvbjogIjEuMCIKY29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciCiAgY29tcGFyZV93aXRoOiAic2luZ2xlX3Jlc3VsdCIKICBpbmNsdWRlX3Jlc3VsdF93aXRoX3Njb3JlOiAicGFzcyIKICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxCmZpbHRlcjoKb2JqZWN0aXZlczoKICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1IgogICAgZGlzcGxheU5hbWU6ICJSZXNwb25zZSB0aW1lIFA5NSIKICAgIGtleV9zbGk6IGZhbHNlCiAgICBwYXNzOiAgICAgICAgICAgICAjIHBhc3MgaWYgKHJlbGF0aXZlIGNoYW5nZSA8PSAxMCUgQU5EIGFic29sdXRlIHZhbHVlIGlzIDwgNjAwbXMpCiAgICAgIC0gY3JpdGVyaWE6CiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykKICAgICAgICAgIC0gIjw2MDAiICAgICMgYWJzb2x1dGUgdmFsdWVzIG9ubHkgcmVxdWlyZSBhIGxvZ2ljYWwgb3BlcmF0b3IKICAgIHdhcm5pbmc6ICAgICAgICAgICMgaWYgdGhlIHJlc3BvbnNlIHRpbWUgaXMgYmVsb3cgODAwbXMsIHRoZSByZXN1bHQgc2hvdWxkIGJlIGEgd2FybmluZwogICAgICAtIGNyaXRlcmlhOgogICAgICAgICAgLSAiPD04MDAiCiAgICB3ZWlnaHQ6IDEKdG90YWxfc2NvcmU6CiAgcGFzczogIjkwJSIKICB3YXJuaW5nOiAiNzUlIg==',
              timeEnd: '2021-10-07T15:02:02.371Z',
              timeStart: '2021-10-07T14:57:02.371Z',
              sloFileContentParsed:
                '---\nspec_version: "1.0"\ncomparison:\n  aggregate_function: "avg"\n  compare_with: "single_result"\n  include_result_with_score: "pass"\n  number_of_comparison_results: 1\nfilter:\nobjectives:\n  - sli: "response_time_p95"\n    displayName: "Response time P95"\n    key_sli: false\n    pass:             # pass if (relative change <= 10% AND absolute value is < 600ms)\n      - criteria:\n          - "<=+10%"  # relative values require a prefixed sign (plus or minus)\n          - "<600"    # absolute values only require a logical operator\n    warning:          # if the response time is below 800ms, the result should be a warning\n      - criteria:\n          - "<=800"\n    weight: 1\ntotal_score:\n  pass: "90%"\n  warning: "75%"',
              score_pass: '90',
              score_warning: '75',
              compare_with: 'single_result\n',
              include_result_with_score: 'pass\n',
              number_of_comparison_results: 1,
            },
            labels: {
              DtCreds: 'dynatrace',
            },
            project: 'sockshop',
            result: 'pass',
            service: 'carts',
            stage: 'staging',
            status: 'succeeded',
            temporaryData: {
              distributor: {
                subscriptionID: '',
              },
            },
          },
          id: 'd7fba175-96b1-4a22-943a-78d871971925',
          source: 'lighthouse-service',
          specversion: '1.0',
          time: '2021-10-07T15:02:27.215Z',
          type: 'sh.keptn.event.evaluation.finished',
          shkeptncontext: '6f0636f7-8f01-4c2a-99b1-1cbf9ff8bfec',
          shkeptnspecversion: '0.2.3',
          triggeredid: '1ee69f43-4a51-46c1-a454-55c925ec4519',
          plainEvent: {
            data: {
              evaluation: {
                comparedEvents: ['53359c6f-fa31-4e8b-9fde-e003b3ea57ec'],
                indicatorResults: [
                  {
                    displayName: 'Response time P95',
                    keySli: false,
                    passTargets: [
                      {
                        criteria: '<=+10%',
                        targetValue: 10.208384090666666,
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
                      value: 9.329947741347496,
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
                  'LS0tCnNwZWNfdmVyc2lvbjogIjEuMCIKY29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciCiAgY29tcGFyZV93aXRoOiAic2luZ2xlX3Jlc3VsdCIKICBpbmNsdWRlX3Jlc3VsdF93aXRoX3Njb3JlOiAicGFzcyIKICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxCmZpbHRlcjoKb2JqZWN0aXZlczoKICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1IgogICAgZGlzcGxheU5hbWU6ICJSZXNwb25zZSB0aW1lIFA5NSIKICAgIGtleV9zbGk6IGZhbHNlCiAgICBwYXNzOiAgICAgICAgICAgICAjIHBhc3MgaWYgKHJlbGF0aXZlIGNoYW5nZSA8PSAxMCUgQU5EIGFic29sdXRlIHZhbHVlIGlzIDwgNjAwbXMpCiAgICAgIC0gY3JpdGVyaWE6CiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykKICAgICAgICAgIC0gIjw2MDAiICAgICMgYWJzb2x1dGUgdmFsdWVzIG9ubHkgcmVxdWlyZSBhIGxvZ2ljYWwgb3BlcmF0b3IKICAgIHdhcm5pbmc6ICAgICAgICAgICMgaWYgdGhlIHJlc3BvbnNlIHRpbWUgaXMgYmVsb3cgODAwbXMsIHRoZSByZXN1bHQgc2hvdWxkIGJlIGEgd2FybmluZwogICAgICAtIGNyaXRlcmlhOgogICAgICAgICAgLSAiPD04MDAiCiAgICB3ZWlnaHQ6IDEKdG90YWxfc2NvcmU6CiAgcGFzczogIjkwJSIKICB3YXJuaW5nOiAiNzUlIg==',
                timeEnd: '2021-10-07T15:02:02.371Z',
                timeStart: '2021-10-07T14:57:02.371Z',
              },
              labels: {
                DtCreds: 'dynatrace',
              },
              project: 'sockshop',
              result: 'pass',
              service: 'carts',
              stage: 'staging',
              status: 'succeeded',
              temporaryData: {
                distributor: {
                  subscriptionID: '',
                },
              },
            },
            id: 'd7fba175-96b1-4a22-943a-78d871971925',
            source: 'lighthouse-service',
            specversion: '1.0',
            time: '2021-10-07T15:02:27.215Z',
            type: 'sh.keptn.event.evaluation.finished',
            shkeptncontext: '6f0636f7-8f01-4c2a-99b1-1cbf9ff8bfec',
            shkeptnspecversion: '0.2.3',
            triggeredid: '1ee69f43-4a51-46c1-a454-55c925ec4519',
          },
          heatmapLabel: '2021-10-07 17:02',
        },
        name: '2021-10-07 17:02',
        x: 6,
      },
    ],
    visible: false,
    turboThreshold: 0,
  },
];

const mocked = TestUtils.mapTraces(evaluationChartItemMock);

export { mocked as EvaluationChartItemMock };
