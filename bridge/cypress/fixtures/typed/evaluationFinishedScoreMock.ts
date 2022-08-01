const evaluationFinishedScoredMock = (score1: number, score2: number): unknown => ({
  events: [
    {
      data: {
        evaluation: {
          indicatorResults: null,
          result: 'pass',
          score: score1,
          sloFileContent: '',
          timeEnd: '2022-02-08T11:44:10.718Z',
          timeStart: '2022-02-08T11:43:40.456Z',
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
      id: '25ab0f26-e6d8-48d5-a08f-08c8a136a687',
      source: 'lighthouse-service',
      specversion: '1.0',
      time: '2022-02-08T11:46:30.853Z',
      type: 'sh.keptn.event.evaluation.finished',
      shkeptncontext: '8f884b2a-2197-4e2f-8284-170ea0a66578',
      shkeptnspecversion: '0.2.3',
      triggeredid: '991eba72-c520-4da5-ba95-d5876101393a',
    },
    {
      data: {
        evaluation: {
          indicatorResults: null,
          result: 'pass',
          score: score2,
          sloFileContent: '',
          timeEnd: '2022-02-08T11:44:19.718Z',
          timeStart: '2022-02-08T11:43:51.456Z',
        },
        labels: {
          DtCreds: 'dynatrace',
          buildId: 'myBuildId',
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
  ],
});
export { evaluationFinishedScoredMock as EvaluationFinishedScoredMock };
