const sequenceDeliveryResponseMock = {
  states: [
    {
      name: 'delivery',
      service: 'carts',
      project: 'sockshop',
      time: '2021-10-13T10:45:06.621Z',
      shkeptncontext: '2c0e568b-8bd3-4726-a188-e528423813ed',
      state: 'finished',
      stages: [
        {
          name: 'dev',
          image: 'docker.io/keptnexamples/carts:0.12.1',
          state: 'finished',
          latestEvaluation: {
            result: 'pass',
            score: 0,
          },
          latestEvent: {
            type: 'sh.keptn.event.dev.delivery.finished',
            id: '861c0880-0555-44e5-b7c3-478e4a7edae4',
            time: '2021-10-13T10:49:30.001Z',
          },
        },
        {
          name: 'staging',
          image: 'docker.io/keptnexamples/carts:0.12.1',
          state: 'finished',
          latestEvaluation: {
            result: 'pass',
            score: 100,
          },
          latestEvent: {
            type: 'sh.keptn.event.staging.delivery.finished',
            id: 'f622d9ac-e673-4e72-b115-7926cc4a137f',
            time: '2021-10-13T10:59:45.002Z',
          },
        },
        {
          name: 'production',
          image: 'docker.io/keptnexamples/carts:0.12.1',
          state: 'finished',
          latestEvent: {
            type: 'sh.keptn.event.production.delivery.finished',
            id: '46ec4f2d-21c4-49bd-94c5-3976c1808abc',
            time: '2021-10-13T11:03:12.212Z',
          },
        },
      ],
    },
  ],
};

export { sequenceDeliveryResponseMock as SequenceDeliveryResponseMock };
