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

const sequenceDeliveryTillStagingResponseMock = {
  states: [
    {
      name: 'delivery',
      service: 'carts',
      project: 'sockshop',
      time: '2021-10-13T10:45:06.621Z',
      shkeptncontext: '2c0e568b-8bd3-4726-a188-e528423813ed',
      state: 'started',
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
          state: 'started',
          latestEvaluation: {
            result: 'pass',
            score: 100,
          },
          latestEvent: {
            type: 'sh.keptn.event.approval.started',
            id: '7d93986c-37bb-41b0-8e97-2e8ee7eee6c0',
            time: '2021-10-13T10:59:45.002Z',
          },
        },
      ],
    },
  ],
};

const sequenceResponseURLFallback = {
  states: [
    {
      name: 'delivery',
      service: 'carts',
      project: 'sockshop',
      time: '2021-11-05T12:20:06.463Z',
      shkeptncontext: 'keptnContext',
      state: 'finished',
      stages: [
        {
          name: 'dev',
          state: 'finished',
          latestEvent: {
            type: 'sh.keptn.event.deployment.finished',
            id: 'eventId',
            time: '2021-11-10T11:59:58.655Z',
          },
          latestFailedEvent: {
            type: 'sh.keptn.event.deployment.finished',
            id: 'eventId',
            time: '2021-11-10T11:59:58.655Z',
          },
        },
      ],
    },
  ],
  totalCount: 1,
};

const sequenceResponseEvaluationFallback = {
  states: [
    {
      name: 'delivery',
      service: 'carts',
      project: 'sockshop',
      time: '2021-11-05T12:20:06.463Z',
      shkeptncontext: 'keptnContext',
      state: 'finished',
      stages: [
        {
          name: 'dev',
          state: 'finished',
          latestEvent: {
            type: 'sh.keptn.event.evaluation.finished',
            id: 'eventId',
            time: '2021-11-10T11:59:58.655Z',
          },
        },
      ],
    },
  ],
  totalCount: 1,
};

const sequencesResponses = {
  states: [
    {
      name: 'delivery',
      service: 'carts',
      project: 'sockshop',
      time: '2021-11-05T12:20:06.463Z',
      shkeptncontext: '2e21574c-dcf7-4275-b677-6bc19214acd5',
      state: 'started',
      stages: [
        {
          name: 'dev',
          image: 'docker.io/keptnexamples/carts:0.12.3',
          state: 'finished',
          latestEvent: {
            type: 'sh.keptn.event.evaluation.started',
            id: 'b7d4225d-26d5-440a-b3c0-3265e2a1735b',
            time: '2021-11-10T11:59:58.655Z',
          },
        },
        {
          name: 'staging',
          image: 'docker.io/keptnexamples/carts:0.12.3',
          state: 'finished',
          latestEvaluation: {
            result: 'fail',
            score: 50,
          },
          latestEvent: {
            type: 'sh.keptn.event.staging.rollback.finished',
            id: '2471ecc3-c407-48be-ae57-7a0e2607fa6d',
            time: '2021-11-10T12:05:31.859Z',
          },
          latestFailedEvent: {
            type: 'sh.keptn.event.staging.delivery.finished',
            id: '4fc419f3-8a65-4c6f-81d2-7150bbef0548',
            time: '2021-11-10T12:05:07.664Z',
          },
        },
      ],
    },
    {
      name: 'delivery-direct',
      service: 'carts-db',
      project: 'sockshop',
      time: '2021-10-12T11:12:55.904Z',
      shkeptncontext: '0cc574e9-3d47-4a29-81b7-84faf33bdc9c',
      state: 'finished',
      stages: [
        {
          name: 'dev',
          image: 'docker.io/mongo:4.2.2',
          state: 'finished',
          latestEvent: {
            type: 'sh.keptn.event.dev.delivery-direct.finished',
            id: 'e170b6cc-27ab-4f9a-a8c9-86cb37c2babd',
            time: '2021-10-12T11:13:34.304Z',
          },
        },
        {
          name: 'staging',
          image: 'docker.io/mongo:4.2.2',
          state: 'finished',
          latestEvent: {
            type: 'sh.keptn.event.staging.delivery-direct.finished',
            id: '7e2ba68e-9780-43ca-9aa7-d181329370e6',
            time: '2021-10-12T11:14:33.202Z',
          },
        },
        {
          name: 'production',
          image: 'docker.io/mongo:4.2.2',
          state: 'finished',
          latestEvent: {
            type: 'sh.keptn.event.production.delivery-direct.finished',
            id: 'aca3bcab-4863-4b87-a1e9-6c328c077348',
            time: '2021-10-12T11:14:59.804Z',
          },
        },
      ],
    },
    {
      name: 'delivery',
      service: 'carts',
      project: 'sockshop',
      time: '2021-10-29T08:40:02.123Z',
      shkeptncontext: '29af69cc-ea85-4358-b169-ce29034d9c81',
      state: 'finished',
      stages: [
        {
          name: 'dev',
          image: 'docker.io/keptnexamples/carts:0.12.3',
          state: 'finished',
          latestEvaluation: {
            result: 'pass',
            score: 0,
          },
          latestEvent: {
            type: 'sh.keptn.event.dev.delivery.finished',
            id: '0f44f059-d2ef-4cbe-875f-7345f90711fb',
            time: '2021-11-05T10:42:19.568Z',
          },
        },
        {
          name: 'staging',
          image: 'docker.io/keptnexamples/carts:0.12.3',
          state: 'finished',
          latestEvaluation: {
            result: 'pass',
            score: 100,
          },
          latestEvent: {
            type: 'sh.keptn.event.staging.rollback.finished',
            id: '236a6da9-d76c-4d6c-bb0c-6e529c6912fb',
            time: '2021-11-05T10:53:00.603Z',
          },
          latestFailedEvent: {
            type: 'sh.keptn.event.staging.delivery.finished',
            id: 'f8d890e0-2951-42c4-8b97-d512af6f068a',
            time: '2021-11-05T10:52:31.260Z',
          },
        },
      ],
    },
    {
      name: 'delivery-direct',
      service: 'carts-db',
      project: 'sockshop',
      time: '2021-10-12T11:12:55.904Z',
      shkeptncontext: '0cc574e9-3d47-4a29-81b7-84faf33bdc9c',
      state: 'finished',
      stages: [
        {
          name: 'dev',
          image: 'docker.io/mongo:4.2.2',
          state: 'finished',
          latestEvent: {
            type: 'sh.keptn.event.dev.delivery-direct.finished',
            id: 'e170b6cc-27ab-4f9a-a8c9-86cb37c2babd',
            time: '2021-10-12T11:13:34.304Z',
          },
        },
        {
          name: 'staging',
          image: 'docker.io/mongo:4.2.2',
          state: 'finished',
          latestEvent: {
            type: 'sh.keptn.event.staging.delivery-direct.finished',
            id: '7e2ba68e-9780-43ca-9aa7-d181329370e6',
            time: '2021-10-12T11:14:33.202Z',
          },
        },
        {
          name: 'production',
          image: 'docker.io/mongo:4.2.2',
          state: 'finished',
          latestEvent: {
            type: 'sh.keptn.event.production.delivery-direct.finished',
            id: 'aca3bcab-4863-4b87-a1e9-6c328c077348',
            time: '2021-10-12T11:14:59.804Z',
          },
        },
      ],
    },
    {
      name: 'remediation',
      problemTitle: 'Failure rate increase',
      service: 'carts',
      project: 'sockshop',
      time: '2021-11-04T04:51:21.557Z',
      shkeptncontext: '35383737-3630-4639-b037-353138323631',
      state: 'finished',
      stages: [
        {
          name: 'production',
          state: 'finished',
          latestEvent: {
            type: 'sh.keptn.event.production.remediation.finished',
            id: '7448420f-5b15-4777-9d39-cc8308e2b0c3',
            time: '2021-11-04T04:51:23.266Z',
          },
          latestFailedEvent: {
            type: 'sh.keptn.event.production.remediation.finished',
            id: '7448420f-5b15-4777-9d39-cc8308e2b0c3',
            time: '2021-11-04T04:51:23.266Z',
          },
        },
      ],
    },
    {
      name: 'delivery-direct',
      service: 'carts-db',
      project: 'sockshop',
      time: '2021-10-12T11:12:55.904Z',
      shkeptncontext: '0cc574e9-3d47-4a29-81b7-84faf33bdc9c',
      state: 'finished',
      stages: [
        {
          name: 'dev',
          image: 'docker.io/mongo:4.2.2',
          state: 'finished',
          latestEvent: {
            type: 'sh.keptn.event.dev.delivery-direct.finished',
            id: 'e170b6cc-27ab-4f9a-a8c9-86cb37c2babd',
            time: '2021-10-12T11:13:34.304Z',
          },
        },
        {
          name: 'staging',
          image: 'docker.io/mongo:4.2.2',
          state: 'finished',
          latestEvent: {
            type: 'sh.keptn.event.staging.delivery-direct.finished',
            id: '7e2ba68e-9780-43ca-9aa7-d181329370e6',
            time: '2021-10-12T11:14:33.202Z',
          },
        },
        {
          name: 'production',
          image: 'docker.io/mongo:4.2.2',
          state: 'finished',
          latestEvent: {
            type: 'sh.keptn.event.production.delivery-direct.finished',
            id: 'aca3bcab-4863-4b87-a1e9-6c328c077348',
            time: '2021-10-12T11:14:59.804Z',
          },
        },
      ],
    },
  ],
  totalCount: 6,
};

export { sequenceDeliveryResponseMock as SequenceDeliveryResponseMock };
export { sequenceDeliveryTillStagingResponseMock as SequenceDeliveryTillStagingResponseMock };
export { sequencesResponses as SequencesResponses };
export { sequenceResponseURLFallback as SequenceResponseURLFallback };
export { sequenceResponseEvaluationFallback as SequenceResponseEvaluationFallback };
