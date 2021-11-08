const openRemediationsResponseMock = {
  states: [
    {
      name: 'remediation',
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
          'Problem URL':
            'https://sjb57563.sprint.dynatracelabs.com/#problems/problemdetails;pid=5877606907518261221_1636001100000V2',
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
          ProblemURL:
            'https://sjb57563.sprint.dynatracelabs.com/#problems/problemdetails;pid=5877606907518261221_1636001100000V2',
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

export { openRemediationsResponseMock as OpenRemediationsResponse };
export { remediationTraceResponse as RemediationTraceResponse };
