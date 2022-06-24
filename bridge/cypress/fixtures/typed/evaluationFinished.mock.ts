const evaluationFinishedMock = (slo?: string): unknown => ({
  events: [
    {
      data: {
        evaluation: {
          comparedEvents: [],
          gitCommit: '',
          indicatorResults: [
            {
              displayName: '',
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
                comparedValue: 0,
                metric: 'http_response_time_seconds_main_page_sum',
                success: true,
                value: 2.00052919040708,
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
              displayName: '',
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
                comparedValue: 0,
                metric: 'request_throughput',
                success: true,
                value: 18.42,
              },
              warningTargets: null,
            },
            {
              displayName: '',
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
                comparedValue: 8,
                metric: 'go_routines',
                success: true,
                value: 88,
              },
              warningTargets: null,
            },
          ],
          result: 'fail',
          score: 33.99999999999999,
          sloFileContent: slo,
          timeEnd: '2022-02-09T09:29:47.625Z',
          timeStart: '2022-02-09T09:28:55.866Z',
        },
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
  ],
});

export { evaluationFinishedMock as EvaluationFinishedMock };
