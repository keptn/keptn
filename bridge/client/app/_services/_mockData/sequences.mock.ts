import { Sequence } from '../../_models/sequence';

let sequencesData = [
  {
    name: 'delivery',
    service: 'carts',
    project: 'sockshop',
    time: '2021-08-02T08:00:21.813Z',
    shkeptncontext: '80f87047-d5c2-4b3e-8aac-f376ff309ed5',
    state: 'finished',
    stages: [
      {
        name: 'dev',
        image: 'docker.io/keptnexamples/carts:0.12.1',
        latestEvent: {
          type: 'sh.keptn.event.dev.delivery.finished',
          id: '852ab5f7-c9c1-4ef9-93b1-3394219b8009',
          time: '2021-08-02T08:03:29.173Z',
        },
        latestFailedEvent: {
          type: 'sh.keptn.event.dev.delivery.finished',
          id: '852ab5f7-c9c1-4ef9-93b1-3394219b8009',
          time: '2021-08-02T08:03:29.173Z',
        },
      },
    ],
  },
  {
    name: 'remediation',
    service: 'carts',
    project: 'sockshop',
    time: '2021-07-31T00:09:15.891Z',
    shkeptncontext: '2d343931-3339-4134-b234-333831333236',
    state: 'finished',
    stages: [
      {
        name: 'production',
        latestEvent: {
          type: 'sh.keptn.event.production.remediation.finished',
          id: '64e1e73c-974d-43e9-91d1-82d916a9be1a',
          time: '2021-07-31T00:09:17.950Z',
        },
        latestFailedEvent: {
          type: 'sh.keptn.event.production.remediation.finished',
          id: '64e1e73c-974d-43e9-91d1-82d916a9be1a',
          time: '2021-07-31T00:09:17.950Z',
        },
      },
    ],
  },
  {
    name: 'delivery',
    service: 'carts',
    project: 'sockshop',
    time: '2021-07-28T11:43:04.585Z',
    shkeptncontext: '722782ba-5ad9-4d8b-a79a-a38c27ea40f1',
    state: 'finished',
    stages: [
      {
        name: 'dev',
        image: 'docker.io/keptnexamples/carts:0.12.1',
        latestEvent: {
          type: 'sh.keptn.event.dev.delivery.finished',
          id: 'c742ee26-0a90-40ae-9f72-7707366e2922',
          time: '2021-07-28T11:46:09.450Z',
        },
        latestFailedEvent: {
          type: 'sh.keptn.event.dev.delivery.finished',
          id: 'c742ee26-0a90-40ae-9f72-7707366e2922',
          time: '2021-07-28T11:46:09.450Z',
        },
      },
    ],
  },
  {
    name: 'delivery',
    service: 'carts',
    project: 'sockshop',
    time: '2021-07-28T11:38:00.783Z',
    shkeptncontext: '75dfd139-2dd3-41df-8dc1-0cb3426f4506',
    state: 'paused',
    stages: [
      {
        name: 'dev',
        image: 'docker.io/keptnexamples/carts:0.12.1',
        latestEvent: {
          type: 'sh.keptn.event.dev.delivery.finished',
          id: 'b0679db1-9830-481e-bd99-90f09dfcba89',
          time: '2021-07-28T11:41:09.547Z',
        },
        latestFailedEvent: {
          type: 'sh.keptn.event.dev.delivery.finished',
          id: 'b0679db1-9830-481e-bd99-90f09dfcba89',
          time: '2021-07-28T11:41:09.547Z',
        },
      },
    ],
  },
  {
    name: 'evaluation',
    service: 'carts',
    project: 'sockshop',
    time: '2021-07-28T06:51:32.244Z',
    shkeptncontext: '9509bdf0-fcc9-452d-8287-c66f28c42858',
    state: 'finished',
    stages: [
      {
        name: 'dev',
        latestEvaluation: {
          result: 'pass',
          score: 0,
        },
        latestEvent: {
          type: 'sh.keptn.event.dev.evaluation.finished',
          id: '83516013-cf51-468d-a46d-51ecdbbbcadf',
          time: '2021-07-28T06:51:41.945Z',
        },
      },
    ],
  },
  {
    name: 'remediation',
    service: 'carts',
    project: 'sockshop',
    time: '2021-07-15T19:23:15.121Z',
    shkeptncontext: '2d313534-3438-4536-b731-393431363639',
    state: 'triggered',
    stages: [
      {
        name: 'dev',
        latestEvent: {
          type: 'sh.keptn.event.dev.remediation.triggered',
          id: '744763a9-697d-4ebe-afd4-2fb0c518b9d2',
          time: '2021-07-15T19:23:15.121Z',
        },
      },
    ],
  },
  {
    name: 'remediation',
    service: 'carts',
    project: 'sockshop',
    time: '2021-07-15T15:23:16.573Z',
    shkeptncontext: '2d383636-3132-4632-b133-323636343331',
    state: 'triggered',
    stages: [
      {
        name: 'dev',
        latestEvent: {
          type: 'sh.keptn.events.problem',
          id: '8f446a88-0012-4fad-a483-c4d4c736a8c6',
          time: '2021-07-15T19:26:15.640Z',
        },
      },
    ],
  },
  {
    name: 'delivery',
    service: 'carts',
    project: 'sockshop',
    time: '2021-07-15T15:17:08.530Z',
    shkeptncontext: '2ae520ec-31ce-4be0-a284-cfbb77a473a3',
    state: 'finished',
    stages: [
      {
        name: 'dev',
        image: 'docker.io/keptnexamples/carts:0.12.2',
        latestEvaluation: {
          result: 'pass',
          score: 0,
        },
        latestEvent: {
          type: 'sh.keptn.event.dev.delivery.finished',
          id: '4a6d31f2-a6fd-41c4-b673-430c834c6442',
          time: '2021-07-15T15:20:57.191Z',
        },
      },
      {
        name: 'staging',
        image: 'docker.io/keptnexamples/carts:0.12.2',
        latestEvent: {
          type: 'sh.keptn.event.staging.rollback.finished',
          id: 'b05b8f69-4854-46cd-82d7-69ce3ee73652',
          time: '2021-07-15T15:27:14.208Z',
        },
        latestFailedEvent: {
          type: 'sh.keptn.event.staging.rollback.finished',
          id: 'b05b8f69-4854-46cd-82d7-69ce3ee73652',
          time: '2021-07-15T15:27:14.208Z',
        },
      },
    ],
  },
  {
    name: 'delivery',
    service: 'carts',
    project: 'sockshop',
    time: '2021-07-15T15:11:20.113Z',
    shkeptncontext: '007649c6-a704-4c21-9ad1-add9bbf20941',
    state: 'finished',
    stages: [
      {
        name: 'dev',
        image: 'docker.io/keptnexamples/carts:0.12.1',
        latestEvaluation: {
          result: 'pass',
          score: 0,
        },
        latestEvent: {
          type: 'sh.keptn.event.dev.delivery.finished',
          id: 'b066e786-53a2-4eea-95af-dfba709f8a27',
          time: '2021-07-15T15:15:17.187Z',
        },
      },
      {
        name: 'staging',
        image: 'docker.io/keptnexamples/carts:0.12.1',
        latestEvent: {
          type: 'sh.keptn.event.staging.rollback.finished',
          id: '0e990b0a-97f1-4d31-9249-97d721603d33',
          time: '2021-07-15T15:21:36.128Z',
        },
        latestFailedEvent: {
          type: 'sh.keptn.event.staging.rollback.finished',
          id: '0e990b0a-97f1-4d31-9249-97d721603d33',
          time: '2021-07-15T15:21:36.128Z',
        },
      },
    ],
  },
  {
    name: 'delivery',
    service: 'carts',
    project: 'sockshop',
    time: '2021-07-15T12:20:24.826Z',
    shkeptncontext: 'd4165d39-e686-46ef-bf5b-d9edcb375e10',
    state: 'finished',
    stages: [
      {
        name: 'dev',
        image: 'docker.io/keptnexamples/carts:0.12.3',
        latestEvaluation: {
          result: 'pass',
          score: 0,
        },
        latestEvent: {
          type: 'sh.keptn.event.dev.delivery.finished',
          id: 'ce06258d-070e-46b8-99f2-d5206c7c4fff',
          time: '2021-07-15T12:24:10.497Z',
        },
      },
      {
        name: 'staging',
        image: 'docker.io/keptnexamples/carts:0.12.3',
        latestEvent: {
          type: 'sh.keptn.event.staging.rollback.finished',
          id: 'eb1e88c7-50f6-4e05-971e-4436e2faf82c',
          time: '2021-07-15T12:30:22.034Z',
        },
        latestFailedEvent: {
          type: 'sh.keptn.event.staging.rollback.finished',
          id: 'eb1e88c7-50f6-4e05-971e-4436e2faf82c',
          time: '2021-07-15T12:30:22.034Z',
        },
      },
    ],
  },
  {
    name: 'delivery',
    service: 'carts',
    project: 'sockshop',
    time: '2021-07-15T11:58:51.729Z',
    shkeptncontext: 'e28592c6-d857-44fe-aea6-e65de02929bf',
    state: 'finished',
    stages: [
      {
        name: 'dev',
        image: 'docker.io/keptnexamples/carts:0.12.3',
        latestEvaluation: {
          result: 'pass',
          score: 0,
        },
        latestEvent: {
          type: 'sh.keptn.event.dev.delivery.finished',
          id: '0838ff2c-f736-41f1-9d7e-5e420f07830d',
          time: '2021-07-15T12:02:45.561Z',
        },
      },
      {
        name: 'staging',
        image: 'docker.io/keptnexamples/carts:0.12.3',
        latestEvaluation: {
          result: 'pass',
          score: 0,
        },
        latestEvent: {
          type: 'sh.keptn.event.staging.rollback.finished',
          id: '5d16ec19-e0ef-453d-b312-03e5988a0905',
          time: '2021-07-15T12:13:40.942Z',
        },
        latestFailedEvent: {
          type: 'sh.keptn.event.staging.rollback.finished',
          id: '5d16ec19-e0ef-453d-b312-03e5988a0905',
          time: '2021-07-15T12:13:40.942Z',
        },
      },
    ],
  },
  {
    name: 'delivery',
    service: 'carts',
    project: 'sockshop',
    time: '2021-07-13T11:28:56.573Z',
    shkeptncontext: '8657fcab-7fd6-4e7a-98a6-e22fea3be97e',
    state: 'finished',
    stages: [
      {
        name: 'dev',
        image: 'docker.io/keptnexamples/carts:0.12.1',
        latestEvaluation: {
          result: 'pass',
          score: 0,
        },
        latestEvent: {
          type: 'sh.keptn.event.dev.delivery.finished',
          id: '1089b1cb-ca34-4203-8a72-ae2a52961689',
          time: '2021-07-13T11:32:47.062Z',
        },
      },
      {
        name: 'staging',
        image: 'docker.io/keptnexamples/carts:0.12.1',
        latestEvaluation: {
          result: 'pass',
          score: 0,
        },
        latestEvent: {
          type: 'sh.keptn.event.staging.delivery.finished',
          id: 'b2de8b58-c717-4c6f-be6c-6375a602713d',
          time: '2021-07-13T11:39:00.165Z',
        },
      },
      {
        name: 'production',
        image: 'docker.io/keptnexamples/carts:0.12.1',
        latestEvent: {
          type: 'sh.keptn.event.production.delivery.finished',
          id: 'd91c4870-b95c-4ef0-944e-19ccafd4b23d',
          time: '2021-07-13T11:41:58.250Z',
        },
      },
    ],
  },
  {
    name: 'delivery-direct',
    service: 'carts-db',
    project: 'sockshop',
    time: '2021-07-13T06:56:12.168Z',
    shkeptncontext: '57c4b0ac-1756-4da8-9f8b-7cf3515a8b13',
    state: 'finished',
    stages: [
      {
        name: 'dev',
        image: 'docker.io/mongo:4.2.2',
        latestEvent: {
          type: 'sh.keptn.event.dev.delivery-direct.finished',
          id: 'c28f8ed8-4b13-47bd-98d5-c0e43eb5df1a',
          time: '2021-07-13T06:56:20.174Z',
        },
      },
      {
        name: 'staging',
        image: 'docker.io/mongo:4.2.2',
        latestEvent: {
          type: 'sh.keptn.event.staging.delivery-direct.finished',
          id: '867ded48-14a8-4848-8a52-d9655973a94a',
          time: '2021-07-13T06:56:30.369Z',
        },
      },
      {
        name: 'production',
        image: 'docker.io/mongo:4.2.2',
        latestEvent: {
          type: 'sh.keptn.event.production.delivery-direct.finished',
          id: '3d06b4f7-a5e7-43b2-83b3-b3ec68ad6aea',
          time: '2021-07-13T06:56:40.267Z',
        },
      },
    ],
  },
  {
    name: 'delivery',
    service: 'carts',
    project: 'sockshop',
    time: '2021-07-12T11:40:09.871Z',
    shkeptncontext: '0073b17b-c614-4d0a-9757-b4b93019e1a2',
    state: 'finished',
    stages: [
      {
        name: 'dev',
        image: 'docker.io/keptnexamples/carts:0.12.1',
        latestEvaluation: {
          result: 'pass',
          score: 0,
        },
        latestEvent: {
          type: 'sh.keptn.event.dev.delivery.finished',
          id: 'a42319cc-c1ff-48c1-8582-d2978a3048f3',
          time: '2021-07-12T11:43:53.484Z',
        },
      },
      {
        name: 'staging',
        image: 'docker.io/keptnexamples/carts:0.12.1',
        latestEvaluation: {
          result: 'fail',
          score: 33.33333333333333,
        },
        latestEvent: {
          type: 'sh.keptn.event.staging.rollback.finished',
          id: 'e8b00713-8e2e-439e-9150-bdcbf1dc6131',
          time: '2021-07-12T11:48:46.208Z',
        },
        latestFailedEvent: {
          type: 'sh.keptn.event.staging.delivery.finished',
          id: 'a8ee31d1-629e-4f3c-8211-0d42e32aa6e2',
          time: '2021-07-12T11:48:35.575Z',
        },
      },
    ],
  },
  {
    name: 'delivery',
    service: 'carts',
    project: 'sockshop',
    time: '2021-07-12T11:17:04.414Z',
    shkeptncontext: '0fca7c48-106f-47ff-a23c-caa9724645e4',
    state: 'finished',
    stages: [
      {
        name: 'dev',
        image: 'docker.io/keptnexamples/carts:0.12.2',
        latestEvaluation: {
          result: 'pass',
          score: 0,
        },
        latestEvent: {
          type: 'sh.keptn.event.dev.delivery.finished',
          id: 'fb1fec26-0b0b-4ddd-83e7-0c6a74838f13',
          time: '2021-07-12T11:20:50.481Z',
        },
      },
      {
        name: 'staging',
        image: 'docker.io/keptnexamples/carts:0.12.2',
        latestEvaluation: {
          result: 'fail',
          score: 0,
        },
        latestEvent: {
          type: 'sh.keptn.event.staging.rollback.finished',
          id: 'e6d53765-4e00-42e0-b428-dc39b8bff019',
          time: '2021-07-12T11:31:35.387Z',
        },
        latestFailedEvent: {
          type: 'sh.keptn.event.staging.delivery.finished',
          id: '85581bab-2f86-4d8f-b1bf-3bcefa9be5c6',
          time: '2021-07-12T11:31:24.658Z',
        },
      },
    ],
  },
  {
    name: 'delivery',
    service: 'carts',
    project: 'sockshop',
    time: '2021-07-12T11:04:59.465Z',
    shkeptncontext: '67af8b62-6cea-48ca-9a42-81dceab3fa69',
    state: 'finished',
    stages: [
      {
        name: 'dev',
        image: 'docker.io/keptnexamples/carts:0.12.1',
        latestEvaluation: {
          result: 'pass',
          score: 0,
        },
        latestEvent: {
          type: 'sh.keptn.event.dev.delivery.finished',
          id: '76b3c992-3c67-493e-883b-95041a7e7dd8',
          time: '2021-07-12T11:08:51.261Z',
        },
      },
      {
        name: 'staging',
        image: 'docker.io/keptnexamples/carts:0.12.1',
        latestEvaluation: {
          result: 'fail',
          score: 33.33333333333333,
        },
        latestEvent: {
          type: 'sh.keptn.event.staging.rollback.finished',
          id: '838c152d-8d3c-4726-a1f8-3227e7e8a724',
          time: '2021-07-12T11:13:44.653Z',
        },
        latestFailedEvent: {
          type: 'sh.keptn.event.staging.delivery.finished',
          id: '1872fa32-0295-45cd-b79e-f3f7e17427bb',
          time: '2021-07-12T11:13:34.253Z',
        },
      },
    ],
  },
  {
    name: 'delivery',
    service: 'carts',
    project: 'sockshop',
    time: '2021-07-12T10:52:41.635Z',
    shkeptncontext: '27f03034-ab6c-4679-ba4f-6592a86fb823',
    state: 'finished',
    stages: [
      {
        name: 'dev',
        image: 'docker.io/keptnexamples/carts:0.12.3',
        latestEvaluation: {
          result: 'pass',
          score: 0,
        },
        latestEvent: {
          type: 'sh.keptn.event.dev.delivery.finished',
          id: '07f71278-1d98-4abe-aafa-026819d65d69',
          time: '2021-07-12T10:56:31.864Z',
        },
      },
      {
        name: 'staging',
        image: 'docker.io/keptnexamples/carts:0.12.3',
        latestEvaluation: {
          result: 'fail',
          score: 33.33333333333333,
        },
        latestEvent: {
          type: 'sh.keptn.event.staging.rollback.finished',
          id: '1e445c95-f6c6-4713-9ac8-0389e4abd0c6',
          time: '2021-07-12T11:02:13.051Z',
        },
        latestFailedEvent: {
          type: 'sh.keptn.event.staging.delivery.finished',
          id: '3d2d3f1b-752e-4513-a185-f01afe3355d3',
          time: '2021-07-12T11:02:00.354Z',
        },
      },
    ],
  },
  {
    name: 'evaluation',
    service: 'carts',
    project: 'sockshop',
    time: '2021-07-06T09:25:10.147Z',
    shkeptncontext: '5e8d9270-0701-42a2-b7c3-6d5b9aaef9e8',
    state: 'finished',
    stages: [
      {
        name: 'staging',
        latestEvaluation: {
          result: 'pass',
          score: 100,
        },
        latestEvent: {
          type: 'sh.keptn.event.staging.evaluation.finished',
          id: '4953c64b-314d-45a0-ab5e-a74d9b47b24e',
          time: '2021-07-06T09:25:20.242Z',
        },
      },
    ],
  },
  {
    name: 'evaluation',
    service: 'carts',
    project: 'sockshop',
    time: '2021-07-06T09:24:42.020Z',
    shkeptncontext: '2f8fce31-3513-4931-bb3d-fa8fc470334c',
    state: 'finished',
    stages: [
      {
        name: 'staging',
        latestEvaluation: {
          result: 'fail',
          score: 66.66666666666666,
        },
        latestEvent: {
          type: 'sh.keptn.event.staging.evaluation.finished',
          id: '99b80add-8f24-4259-b52e-e3a0eac74fd3',
          time: '2021-07-06T09:24:51.353Z',
        },
        latestFailedEvent: {
          type: 'sh.keptn.event.staging.evaluation.finished',
          id: '99b80add-8f24-4259-b52e-e3a0eac74fd3',
          time: '2021-07-06T09:24:51.353Z',
        },
      },
    ],
  },
  {
    name: 'evaluation',
    service: 'carts',
    project: 'sockshop',
    time: '2021-07-06T09:22:56.433Z',
    shkeptncontext: '43580c5f-258f-4aae-9883-5a5c33fe4516',
    state: 'finished',
    stages: [
      {
        name: 'staging',
        latestEvaluation: {
          result: 'fail',
          score: 66.66666666666666,
        },
        latestEvent: {
          type: 'sh.keptn.event.staging.evaluation.finished',
          id: '118d424d-33ac-4d15-88c4-2ae77fcc8ad1',
          time: '2021-07-06T09:23:09.238Z',
        },
        latestFailedEvent: {
          type: 'sh.keptn.event.staging.evaluation.finished',
          id: '118d424d-33ac-4d15-88c4-2ae77fcc8ad1',
          time: '2021-07-06T09:23:09.238Z',
        },
      },
    ],
  },
  {
    name: 'evaluation',
    service: 'carts',
    project: 'sockshop',
    time: '2021-07-06T09:21:43.659Z',
    shkeptncontext: '88d5f9c0-1e2c-474b-af27-13fe12e038a7',
    state: 'finished',
    stages: [
      {
        name: 'staging',
        latestEvaluation: {
          result: 'fail',
          score: 66.66666666666666,
        },
        latestEvent: {
          type: 'sh.keptn.event.staging.evaluation.finished',
          id: '16614663-aea7-4420-bafe-41f480b572d8',
          time: '2021-07-06T09:21:53.145Z',
        },
        latestFailedEvent: {
          type: 'sh.keptn.event.staging.evaluation.finished',
          id: '16614663-aea7-4420-bafe-41f480b572d8',
          time: '2021-07-06T09:21:53.145Z',
        },
      },
    ],
  },
  {
    name: 'evaluation',
    service: 'carts',
    project: 'sockshop',
    time: '2021-07-06T09:13:00.939Z',
    shkeptncontext: 'f83f1f12-4601-4fbd-a6f7-304458e5442c',
    state: 'finished',
    stages: [
      {
        name: 'staging',
        latestEvaluation: {
          result: 'fail',
          score: 66.66666666666666,
        },
        latestEvent: {
          type: 'sh.keptn.event.staging.evaluation.finished',
          id: '10643917-00cc-46da-873b-5254626663cd',
          time: '2021-07-06T09:13:13.942Z',
        },
        latestFailedEvent: {
          type: 'sh.keptn.event.staging.evaluation.finished',
          id: '10643917-00cc-46da-873b-5254626663cd',
          time: '2021-07-06T09:13:13.942Z',
        },
      },
    ],
  },
  {
    name: 'evaluation',
    service: 'carts',
    project: 'sockshop',
    time: '2021-07-06T09:10:45.712Z',
    shkeptncontext: '714397cc-e648-474f-bd6e-fa79463691d3',
    state: 'finished',
    stages: [
      {
        name: 'staging',
        latestEvaluation: {
          result: 'fail',
          score: 66.66666666666666,
        },
        latestEvent: {
          type: 'sh.keptn.event.staging.evaluation.finished',
          id: 'e5cf45e4-b8f0-48a9-8cdb-ec0ee3ea0219',
          time: '2021-07-06T09:10:55.248Z',
        },
        latestFailedEvent: {
          type: 'sh.keptn.event.staging.evaluation.finished',
          id: 'e5cf45e4-b8f0-48a9-8cdb-ec0ee3ea0219',
          time: '2021-07-06T09:10:55.248Z',
        },
      },
    ],
  },
  {
    name: 'evaluation',
    service: 'carts',
    project: 'sockshop',
    time: '2021-07-06T09:02:09.612Z',
    shkeptncontext: '5b8d8f9c-faba-43f2-9c87-fc0b69d6fc3e',
    state: 'finished',
    stages: [
      {
        name: 'staging',
        latestEvaluation: {
          result: 'pass',
          score: 100,
        },
        latestEvent: {
          type: 'sh.keptn.event.staging.evaluation.finished',
          id: 'ec195981-f095-4f7b-9568-be9ca0fe18c5',
          time: '2021-07-06T09:02:26.842Z',
        },
      },
    ],
  },
  {
    name: 'evaluation',
    service: 'carts',
    project: 'sockshop',
    time: '2021-07-06T09:02:06.635Z',
    shkeptncontext: 'd229e6f9-573b-4d07-8d08-b2502cb8c320',
    state: 'finished',
    stages: [
      {
        name: 'staging',
        latestEvaluation: {
          result: 'pass',
          score: 100,
        },
        latestEvent: {
          type: 'sh.keptn.event.staging.evaluation.finished',
          id: '4a3bd7d7-a786-4e36-a781-a1b22e06b556',
          time: '2021-07-06T09:02:23.043Z',
        },
      },
    ],
  },
  {
    name: 'evaluation',
    service: 'carts',
    project: 'sockshop',
    time: '2021-07-06T09:01:40.989Z',
    shkeptncontext: '450e8156-ac3e-4e76-89f9-364a92faef3a',
    state: 'finished',
    stages: [
      {
        name: 'staging',
        latestEvaluation: {
          result: 'fail',
          score: 66.66666666666666,
        },
        latestEvent: {
          type: 'sh.keptn.event.staging.evaluation.finished',
          id: '46fc0a0d-739e-4101-b043-bad09399bf70',
          time: '2021-07-06T09:01:50.662Z',
        },
        latestFailedEvent: {
          type: 'sh.keptn.event.staging.evaluation.finished',
          id: '46fc0a0d-739e-4101-b043-bad09399bf70',
          time: '2021-07-06T09:01:50.662Z',
        },
      },
    ],
  },
  {
    name: 'evaluation',
    service: 'carts',
    project: 'sockshop',
    time: '2021-07-06T08:59:04.719Z',
    shkeptncontext: 'bd41e3fe-edbd-4245-8d18-93433143ac7f',
    state: 'finished',
    stages: [
      {
        name: 'staging',
        latestEvaluation: {
          result: 'pass',
          score: 100,
        },
        latestEvent: {
          type: 'sh.keptn.event.staging.evaluation.finished',
          id: '31d7a31f-acb7-4f2b-820f-ee387d4f187d',
          time: '2021-07-06T08:59:14.340Z',
        },
      },
    ],
  },
  {
    name: 'evaluation',
    service: 'carts',
    project: 'sockshop',
    time: '2021-07-06T08:58:36.662Z',
    shkeptncontext: '9a7eaa00-1019-41e5-be3a-441738338cf4',
    state: 'finished',
    stages: [
      {
        name: 'staging',
        latestEvaluation: {
          result: 'fail',
          score: 66.66666666666666,
        },
        latestEvent: {
          type: 'sh.keptn.event.staging.evaluation.finished',
          id: '993ec5e5-c3e2-4515-8bf0-f9be93f77125',
          time: '2021-07-06T08:58:45.965Z',
        },
        latestFailedEvent: {
          type: 'sh.keptn.event.staging.evaluation.finished',
          id: '993ec5e5-c3e2-4515-8bf0-f9be93f77125',
          time: '2021-07-06T08:58:45.965Z',
        },
      },
    ],
  },
  {
    name: 'evaluation',
    service: 'carts',
    project: 'sockshop',
    time: '2021-07-06T08:57:38.264Z',
    shkeptncontext: '238e1c43-501a-4a7f-a725-802cf36e16f5',
    state: 'finished',
    stages: [
      {
        name: 'staging',
        latestEvaluation: {
          result: 'fail',
          score: 66.66666666666666,
        },
        latestEvent: {
          type: 'sh.keptn.event.staging.evaluation.finished',
          id: 'a6222bca-2619-48ea-933e-d4f2fbf97937',
          time: '2021-07-06T08:57:48.047Z',
        },
        latestFailedEvent: {
          type: 'sh.keptn.event.staging.evaluation.finished',
          id: 'a6222bca-2619-48ea-933e-d4f2fbf97937',
          time: '2021-07-06T08:57:48.047Z',
        },
      },
    ],
  },
  {
    name: 'delivery',
    service: 'carts',
    project: 'sockshop',
    time: '2021-07-06T08:13:53.766Z',
    shkeptncontext: '1663de8a-a414-47ba-9566-10a9730f406f',
    state: 'finished',
    stages: [
      {
        name: 'dev',
        image: 'docker.io/keptnexamples/carts:0.12.3',
        latestEvaluation: {
          result: 'pass',
          score: 0,
        },
        latestEvent: {
          type: 'sh.keptn.event.dev.delivery.finished',
          id: '92207803-7d49-4081-ad0a-897e3f4229d2',
          time: '2021-07-06T08:17:42.285Z',
        },
      },
      {
        name: 'staging',
        image: 'docker.io/keptnexamples/carts:0.12.3',
        latestEvaluation: {
          result: 'fail',
          score: 33.33333333333333,
        },
        latestEvent: {
          type: 'sh.keptn.event.staging.rollback.finished',
          id: 'd23e5ed6-0025-4273-bc7f-09ff94eab160',
          time: '2021-07-06T08:22:47.310Z',
        },
        latestFailedEvent: {
          type: 'sh.keptn.event.staging.delivery.finished',
          id: '26d927cd-274e-406d-bb40-d7a47a553632',
          time: '2021-07-06T08:22:36.842Z',
        },
      },
    ],
  },
  {
    name: 'delivery',
    service: 'carts',
    project: 'sockshop',
    time: '2021-07-06T08:04:16.338Z',
    shkeptncontext: '26d6fbd3-42e5-4ec7-9f22-4fd102a2e414',
    state: 'finished',
    stages: [
      {
        name: 'dev',
        image: 'docker.io/keptnexamples/carts:0.12.1',
        latestEvaluation: {
          result: 'pass',
          score: 0,
        },
        latestEvent: {
          type: 'sh.keptn.event.dev.delivery.finished',
          id: 'a3ba36e1-2edb-419f-8a24-d624e2100d43',
          time: '2021-07-06T08:07:55.662Z',
        },
      },
      {
        name: 'staging',
        image: 'docker.io/keptnexamples/carts:0.12.1',
        latestEvaluation: {
          result: 'fail',
          score: 33.33333333333333,
        },
        latestEvent: {
          type: 'sh.keptn.event.staging.rollback.finished',
          id: '9892abf9-0528-4108-834c-47f98dc6caa2',
          time: '2021-07-06T08:12:57.479Z',
        },
        latestFailedEvent: {
          type: 'sh.keptn.event.staging.delivery.finished',
          id: '9693825e-b6d7-43a3-97ee-d136ed3f03c3',
          time: '2021-07-06T08:12:47.145Z',
        },
      },
    ],
  },
  {
    name: 'delivery',
    service: 'carts',
    project: 'sockshop',
    time: '2021-07-06T07:50:28.004Z',
    shkeptncontext: '29761b1a-374f-4bb1-b5ab-3e9d809b6108',
    state: 'finished',
    stages: [
      {
        name: 'dev',
        image: 'docker.io/keptnexamples/carts:0.12.3',
        latestEvaluation: {
          result: 'pass',
          score: 0,
        },
        latestEvent: {
          type: 'sh.keptn.event.dev.delivery.finished',
          id: 'fbb56094-9ef7-45df-bc4e-525202e4c302',
          time: '2021-07-06T07:54:21.353Z',
        },
      },
      {
        name: 'staging',
        image: 'docker.io/keptnexamples/carts:0.12.3',
        latestEvaluation: {
          result: 'fail',
          score: 33.33333333333333,
        },
        latestEvent: {
          type: 'sh.keptn.event.staging.rollback.finished',
          id: 'c5f88cd1-7efb-4f24-b0b3-eb5998a9d00d',
          time: '2021-07-06T07:59:12.673Z',
        },
        latestFailedEvent: {
          type: 'sh.keptn.event.staging.delivery.finished',
          id: '40dc9d08-f87a-4785-a037-7d76a9745a0c',
          time: '2021-07-06T07:58:58.241Z',
        },
      },
    ],
  },
  {
    name: 'evaluation',
    service: 'carts',
    project: 'sockshop',
    time: '2021-07-05T14:12:15.566Z',
    shkeptncontext: 'dda949d7-81e8-43db-a002-b736c2d27843',
    state: 'finished',
    stages: [
      {
        name: 'staging',
        latestEvaluation: {
          result: 'pass',
          score: 100,
        },
        latestEvent: {
          type: 'sh.keptn.event.staging.evaluation.finished',
          id: '4b7e9733-46d2-44a5-aadc-709e881d6403',
          time: '2021-07-05T14:12:25.169Z',
        },
      },
    ],
  },
  {
    name: 'evaluation',
    service: 'carts',
    project: 'sockshop',
    time: '2021-07-05T13:47:23.145Z',
    shkeptncontext: '23e404ee-e9ab-4e06-a97f-6408854e6553',
    state: 'finished',
    stages: [
      {
        name: 'staging',
        latestEvaluation: {
          result: 'pass',
          score: 100,
        },
        latestEvent: {
          type: 'sh.keptn.event.staging.evaluation.finished',
          id: 'bbed850c-8c43-431e-86b0-0cffbc0f938f',
          time: '2021-07-05T13:47:32.749Z',
        },
      },
    ],
  },
  {
    name: 'evaluation',
    service: 'carts',
    project: 'sockshop',
    time: '2021-07-05T13:35:36.482Z',
    shkeptncontext: '7a5127e2-8dc0-4d85-9046-450ef054484b',
    state: 'finished',
    stages: [
      {
        name: 'staging',
        latestEvaluation: {
          result: 'fail',
          score: 66.66666666666666,
        },
        latestEvent: {
          type: 'sh.keptn.event.staging.evaluation.finished',
          id: '074fe4e5-af9d-41fb-ac7e-92e7e8c4e0a3',
          time: '2021-07-05T13:35:45.947Z',
        },
        latestFailedEvent: {
          type: 'sh.keptn.event.staging.evaluation.finished',
          id: '074fe4e5-af9d-41fb-ac7e-92e7e8c4e0a3',
          time: '2021-07-05T13:35:45.947Z',
        },
      },
    ],
  },
  {
    name: 'delivery',
    service: 'carts',
    project: 'sockshop',
    time: '2021-07-05T13:10:34.235Z',
    shkeptncontext: 'db6a190d-eea6-49b3-8c7c-aaea24c0015c',
    state: 'finished',
    stages: [
      {
        name: 'dev',
        image: 'docker.io/keptnexamples/carts:0.12.3',
        latestEvaluation: {
          result: 'pass',
          score: 0,
        },
        latestEvent: {
          type: 'sh.keptn.event.dev.delivery.finished',
          id: 'cac88040-0b0a-478f-a75e-57e17f6a1171',
          time: '2021-07-05T13:14:26.464Z',
        },
      },
      {
        name: 'staging',
        image: 'docker.io/keptnexamples/carts:0.12.3',
        latestEvaluation: {
          result: 'fail',
          score: 33.33333333333333,
        },
        latestEvent: {
          type: 'sh.keptn.event.staging.rollback.finished',
          id: '88cd8f8c-def5-4868-aefd-8174ed94738a',
          time: '2021-07-05T13:19:59.459Z',
        },
        latestFailedEvent: {
          type: 'sh.keptn.event.staging.delivery.finished',
          id: '4374b25b-03f6-449b-b9aa-a488d7ed7c4e',
          time: '2021-07-05T13:19:48.958Z',
        },
      },
    ],
  },
] as Sequence[];
sequencesData = sequencesData.map((sequence) => Sequence.fromJSON(sequence));
export { sequencesData as SequencesData };
