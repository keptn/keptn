import {
  DevCartsDbDeploymentFinished,
  DevCartsDeploymentFinished,
  ProductionCartsDbDeploymentFinished,
  ProductionCartsDeploymentFinished,
  StagingCartsDbDeploymentFinished,
  StagingCartsDeploymentFinished,
  StagingCartsEvaluationFinished,
} from './project-response.mock';

const keptnContext = '2c0e568b-8bd3-4726-a188-e528423813ed';

const defaultEvaluationData = {
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
        value: 304.2952915485157,
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
  timeEnd: '2021-10-13T10:53:29Z',
  timeStart: '2021-10-13T10:51:08Z',
};

const defaultEvaluationFinishedTrace = {
  data: {
    evaluation: defaultEvaluationData,
    labels: {
      DtCreds: 'dynatrace',
    },
    project: 'sockshop',
    result: 'pass',
    service: 'carts',
    stage: 'dev',
    status: 'succeeded',
    temporaryData: {
      distributor: {
        subscriptionID: '',
      },
    },
  },
  id: '93c2eba9-b77c-4976-b079-29a0188d86ef',
  source: 'lighthouse-service',
  specversion: '1.0',
  time: '2021-10-13T10:54:43.112Z',
  type: 'sh.keptn.event.evaluation.finished',
  shkeptncontext: keptnContext,
  shkeptnspecversion: '0.2.3',
  triggeredid: '1cc9c272-721a-43de-98f6-9eceae484cf5',
};

const defaultDeploymentData = {
  deploymentNames: ['direct'],
  deploymentURIsLocal: ['http://carts.sockshop-dev:80'],
  deploymentURIsPublic: ['http://carts.sockshop-dev.35.192.209.116.nip.io:80'],
  deploymentstrategy: 'direct',
};

const defaultDeploymentFinishedTrace = {
  data: {
    deployment: defaultDeploymentData,
    message: 'Successfully deployed',
    project: 'sockshop',
    result: 'pass',
    service: 'carts',
    stage: 'dev',
    status: 'succeeded',
    temporaryData: {
      distributor: {
        subscriptionID: '',
      },
    },
  },
  id: '8b65ab52-f1bb-4ecf-8559-48854acfa60d',
  source: 'helm-service',
  specversion: '1.0',
  time: '2021-10-13T10:46:41.861Z',
  type: 'sh.keptn.event.deployment.finished',
  shkeptncontext: keptnContext,
  shkeptnspecversion: '0.2.3',
  triggeredid: '89265987-90b9-4554-8d13-98deb7b44f9d',
};

const deliveryFinishedProduction = {
  data: {
    labels: {
      DtCreds: 'dynatrace',
    },
    message: 'Finished release',
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
  id: '46ec4f2d-21c4-49bd-94c5-3976c1808abc',
  source: 'shipyard-controller',
  specversion: '1.0',
  time: '2021-10-13T11:03:12.215Z',
  type: 'sh.keptn.event.production.delivery.finished',
  shkeptncontext: keptnContext,
  shkeptnspecversion: '0.2.3',
  triggeredid: '61a2d8f9-7368-4097-b469-2fb81af50eb3',
};

const releaseTracesProduction = [
  {
    data: {
      labels: {
        DtCreds: 'dynatrace',
      },
      message: 'Finished release',
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
    id: '934c5141-a130-4aa7-9cba-00c162309524',
    source: 'helm-service',
    specversion: '1.0',
    time: '2021-10-13T11:03:12.129Z',
    type: 'sh.keptn.event.release.finished',
    shkeptncontext: keptnContext,
    shkeptnspecversion: '0.2.3',
    triggeredid: '3926781b-9cb0-4538-938d-f3f6dfdd46d4',
  },
  {
    data: {
      labels: {
        DtCreds: 'dynatrace',
      },
      project: 'sockshop',
      service: 'carts',
      stage: 'production',
      status: 'succeeded',
      temporaryData: {
        distributor: {
          subscriptionID: '',
        },
      },
    },
    id: '0049ede6-11a7-4146-8b0e-8fd718595a64',
    source: 'helm-service',
    specversion: '1.0',
    time: '2021-10-13T11:01:18.704Z',
    type: 'sh.keptn.event.release.started',
    shkeptncontext: keptnContext,
    shkeptnspecversion: '0.2.3',
    triggeredid: '3926781b-9cb0-4538-938d-f3f6dfdd46d4',
  },
  {
    data: {
      configurationChange: {
        values: {
          image: 'docker.io/keptnexamples/carts:0.12.1',
        },
      },
      deployment: {
        deploymentNames: ['canary', 'direct', 'canary', 'direct'],
        deploymentURIsLocal: ['http://carts.sockshop-production:80'],
        deploymentURIsPublic: [
          'http://carts.sockshop-production.35.192.209.116.nip.io:80',
          'http://carts.sockshop-dev.35.192.209.116.nip.io:80',
          'http://carts.sockshop-staging.35.192.209.116.nip.io:80',
          'http://carts.sockshop-dev.35.192.209.116.nip.io:80',
        ],
        deploymentstrategy: 'duplicate',
      },
      evaluation: {
        comparedEvents: ['3344487d-e384-4cd9-a0e0-fcf157a33ad6'],
        indicatorResults: null,
        responses: ['Redirecting to https://keptn.sh/\n', 'Redirecting to https://keptn.sh/\n'],
        result: 'pass',
        score: 0,
        sloFileContent: '',
        timeEnd: '2021-10-13T10:47:11Z',
        timeStart: '2021-10-13T10:46:42Z',
      },
      labels: {
        DtCreds: 'dynatrace',
      },
      message: '',
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
      test: {
        end: '2021-10-13T10:47:11Z',
        start: '2021-10-13T10:46:42Z',
      },
    },
    id: '3926781b-9cb0-4538-938d-f3f6dfdd46d4',
    source: 'shipyard-controller',
    specversion: '1.0',
    time: '2021-10-13T11:01:18.701Z',
    type: 'sh.keptn.event.release.triggered',
    shkeptncontext: keptnContext,
    shkeptnspecversion: '0.2.3',
  },
];

const deliveryTriggeredProduction = {
  data: {
    configurationChange: {
      values: {
        image: 'docker.io/keptnexamples/carts:0.12.1',
      },
    },
    deployment: {
      deploymentNames: ['direct', 'canary', 'direct'],
      deploymentURIsLocal: null,
      deploymentURIsPublic: [
        'http://carts.sockshop-dev.35.192.209.116.nip.io:80',
        'http://carts.sockshop-staging.35.192.209.116.nip.io:80',
        'http://carts.sockshop-dev.35.192.209.116.nip.io:80',
      ],
      deploymentstrategy: '',
    },
    evaluation: {
      comparedEvents: ['3344487d-e384-4cd9-a0e0-fcf157a33ad6'],
      indicatorResults: null,
      responses: ['Redirecting to https://keptn.sh/\n', 'Redirecting to https://keptn.sh/\n'],
      result: 'pass',
      score: 0,
      sloFileContent: '',
      timeEnd: '2021-10-13T10:47:11Z',
      timeStart: '2021-10-13T10:46:42Z',
    },
    labels: {
      DtCreds: 'dynatrace',
    },
    message: 'Finished release',
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
    test: {
      end: '2021-10-13T10:47:11Z',
      start: '2021-10-13T10:46:42Z',
    },
  },
  id: '61a2d8f9-7368-4097-b469-2fb81af50eb3',
  source: 'shipyard-controller',
  specversion: '1.0',
  time: '2021-10-13T10:59:45.304Z',
  type: 'sh.keptn.event.production.delivery.triggered',
  shkeptncontext: keptnContext,
  shkeptnspecversion: '0.2.3',
};

const deliveryFinishedStaging = {
  data: {
    labels: {
      DtCreds: 'dynatrace',
    },
    message: 'Finished release',
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
  id: 'f622d9ac-e673-4e72-b115-7926cc4a137f',
  source: 'shipyard-controller',
  specversion: '1.0',
  time: '2021-10-13T10:59:45.104Z',
  type: 'sh.keptn.event.staging.delivery.finished',
  shkeptncontext: keptnContext,
  shkeptnspecversion: '0.2.3',
  triggeredid: '3cbb5949-3852-4073-a70d-27ec52e04b93',
};

const deliveryTriggeredStaging = {
  data: {
    configurationChange: {
      values: {
        image: 'docker.io/keptnexamples/carts:0.12.1',
      },
    },
    deployment: {
      deploymentNames: ['direct'],
      deploymentURIsLocal: null,
      deploymentURIsPublic: ['http://carts.sockshop-dev.35.192.209.116.nip.io:80'],
      deploymentstrategy: '',
    },
    evaluation: {
      indicatorResults: null,
      responses: ['Redirecting to https://keptn.sh/\n'],
      result: 'pass',
      score: 0,
      sloFileContent: '',
      timeEnd: '2021-10-13T10:47:11Z',
      timeStart: '2021-10-13T10:46:42Z',
    },
    labels: {
      DtCreds: 'dynatrace',
    },
    message: 'Finished release',
    project: 'sockshop',
    result: '',
    service: 'carts',
    stage: 'staging',
    status: '',
    temporaryData: {
      distributor: {
        subscriptionID: '',
      },
    },
    test: {
      end: '2021-10-13T10:47:11Z',
      start: '2021-10-13T10:46:42Z',
    },
  },
  id: '3cbb5949-3852-4073-a70d-27ec52e04b93',
  source: 'shipyard-controller',
  specversion: '1.0',
  time: '2021-10-13T10:49:30.202Z',
  type: 'sh.keptn.event.staging.delivery.triggered',
  shkeptncontext: keptnContext,
  shkeptnspecversion: '0.2.3',
};

const deploymentTracesProduction = [
  {
    data: {
      deployment: {
        deploymentNames: ['canary'],
        deploymentURIsLocal: ['http://carts.sockshop-production:80'],
        deploymentURIsPublic: ['http://carts.sockshop-production.35.192.209.116.nip.io:80'],
        deploymentstrategy: 'duplicate',
      },
      labels: {
        DtCreds: 'dynatrace',
      },
      message: 'Successfully deployed',
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
    id: '9998714f-64f1-431f-8423-7bbcef957604',
    source: 'helm-service',
    specversion: '1.0',
    time: '2021-10-13T11:01:18.562Z',
    type: 'sh.keptn.event.deployment.finished',
    shkeptncontext: keptnContext,
    shkeptnspecversion: '0.2.3',
    triggeredid: '983e3fb5-c868-4fcc-90d2-5724b499826a',
  },
  {
    data: {
      labels: {
        DtCreds: 'dynatrace',
      },
      project: 'sockshop',
      service: 'carts',
      stage: 'production',
      status: 'succeeded',
      temporaryData: {
        distributor: {
          subscriptionID: '',
        },
      },
    },
    id: 'bef96141-19c6-473d-8eda-468f8d4f7dce',
    source: 'helm-service',
    specversion: '1.0',
    time: '2021-10-13T10:59:46.822Z',
    type: 'sh.keptn.event.deployment.started',
    shkeptncontext: keptnContext,
    shkeptnspecversion: '0.2.3',
    triggeredid: '983e3fb5-c868-4fcc-90d2-5724b499826a',
  },
  {
    data: {
      configurationChange: {
        values: {
          image: 'docker.io/keptnexamples/carts:0.12.1',
        },
      },
      deployment: {
        deploymentNames: ['direct', 'canary', 'direct'],
        deploymentURIsLocal: null,
        deploymentURIsPublic: [
          'http://carts.sockshop-dev.35.192.209.116.nip.io:80',
          'http://carts.sockshop-staging.35.192.209.116.nip.io:80',
          'http://carts.sockshop-dev.35.192.209.116.nip.io:80',
        ],
        deploymentstrategy: 'blue_green_service',
      },
      evaluation: {
        comparedEvents: ['3344487d-e384-4cd9-a0e0-fcf157a33ad6'],
        indicatorResults: null,
        responses: ['Redirecting to https://keptn.sh/\n', 'Redirecting to https://keptn.sh/\n'],
        result: 'pass',
        score: 0,
        sloFileContent: '',
        timeEnd: '2021-10-13T10:47:11Z',
        timeStart: '2021-10-13T10:46:42Z',
      },
      labels: {
        DtCreds: 'dynatrace',
      },
      message: '',
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
      test: {
        end: '2021-10-13T10:47:11Z',
        start: '2021-10-13T10:46:42Z',
      },
    },
    id: '983e3fb5-c868-4fcc-90d2-5724b499826a',
    source: 'shipyard-controller',
    specversion: '1.0',
    time: '2021-10-13T10:59:46.817Z',
    type: 'sh.keptn.event.deployment.triggered',
    shkeptncontext: keptnContext,
    shkeptnspecversion: '0.2.3',
  },
];

const releaseTracesStaging = [
  {
    data: {
      labels: {
        DtCreds: 'dynatrace',
      },
      message: 'Finished release',
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
    id: '6430a35c-84d1-4d56-9d66-05ea30d5c9f0',
    source: 'helm-service',
    specversion: '1.0',
    time: '2021-10-13T10:59:44.569Z',
    type: 'sh.keptn.event.release.finished',
    shkeptncontext: keptnContext,
    shkeptnspecversion: '0.2.3',
    triggeredid: 'ddbb13be-918e-498d-9c8a-148693237c23',
  },
  {
    data: {
      labels: {
        DtCreds: 'dynatrace',
      },
      project: 'sockshop',
      service: 'carts',
      stage: 'staging',
      status: 'succeeded',
      temporaryData: {
        distributor: {
          subscriptionID: '',
        },
      },
    },
    id: '43278c61-5cb3-408b-bf0c-54b360ff834f',
    source: 'helm-service',
    specversion: '1.0',
    time: '2021-10-13T10:57:44.406Z',
    type: 'sh.keptn.event.release.started',
    shkeptncontext: keptnContext,
    shkeptnspecversion: '0.2.3',
    triggeredid: 'ddbb13be-918e-498d-9c8a-148693237c23',
  },
  {
    data: {
      configurationChange: {
        values: {
          image: 'docker.io/keptnexamples/carts:0.12.1',
        },
      },
      deployment: {
        deploymentNames: ['canary', 'direct'],
        deploymentURIsLocal: ['http://carts.sockshop-staging:80'],
        deploymentURIsPublic: [
          'http://carts.sockshop-staging.35.192.209.116.nip.io:80',
          'http://carts.sockshop-dev.35.192.209.116.nip.io:80',
        ],
        deploymentstrategy: 'duplicate',
      },
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
              value: 304.2952915485157,
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
        responses: ['Redirecting to https://keptn.sh/\n'],
        result: 'pass',
        score: 100,
        sloFileContent:
          'LS0tCnNwZWNfdmVyc2lvbjogIjEuMCIKY29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciCiAgY29tcGFyZV93aXRoOiAic2luZ2xlX3Jlc3VsdCIKICBpbmNsdWRlX3Jlc3VsdF93aXRoX3Njb3JlOiAicGFzcyIKICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxCmZpbHRlcjoKb2JqZWN0aXZlczoKICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1IgogICAgZGlzcGxheU5hbWU6ICJSZXNwb25zZSB0aW1lIFA5NSIKICAgIGtleV9zbGk6IGZhbHNlCiAgICBwYXNzOiAgICAgICAgICAgICAjIHBhc3MgaWYgKHJlbGF0aXZlIGNoYW5nZSA8PSAxMCUgQU5EIGFic29sdXRlIHZhbHVlIGlzIDwgNjAwbXMpCiAgICAgIC0gY3JpdGVyaWE6CiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykKICAgICAgICAgIC0gIjw2MDAiICAgICMgYWJzb2x1dGUgdmFsdWVzIG9ubHkgcmVxdWlyZSBhIGxvZ2ljYWwgb3BlcmF0b3IKICAgIHdhcm5pbmc6ICAgICAgICAgICMgaWYgdGhlIHJlc3BvbnNlIHRpbWUgaXMgYmVsb3cgODAwbXMsIHRoZSByZXN1bHQgc2hvdWxkIGJlIGEgd2FybmluZwogICAgICAtIGNyaXRlcmlhOgogICAgICAgICAgLSAiPD04MDAiCiAgICB3ZWlnaHQ6IDEKdG90YWxfc2NvcmU6CiAgcGFzczogIjkwJSIKICB3YXJuaW5nOiAiNzUlIg==',
        timeEnd: '2021-10-13T10:53:29Z',
        timeStart: '2021-10-13T10:51:08Z',
      },
      labels: {
        DtCreds: 'dynatrace',
      },
      message: '',
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
      test: {
        end: '2021-10-13T10:53:29Z',
        start: '2021-10-13T10:51:08Z',
      },
    },
    id: 'ddbb13be-918e-498d-9c8a-148693237c23',
    source: 'shipyard-controller',
    specversion: '1.0',
    time: '2021-10-13T10:57:44.402Z',
    type: 'sh.keptn.event.release.triggered',
    shkeptncontext: keptnContext,
    shkeptnspecversion: '0.2.3',
  },
];

const deploymentTracesStaging = [
  {
    data: {
      deployment: {
        deploymentNames: ['canary'],
        deploymentURIsLocal: ['http://carts.sockshop-staging:80'],
        deploymentURIsPublic: ['http://carts.sockshop-staging.35.192.209.116.nip.io:80'],
        deploymentstrategy: 'duplicate',
      },
      labels: {
        DtCreds: 'dynatrace',
      },
      message: 'Successfully deployed',
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
    id: '4515cce9-bf02-4263-ba0d-2d46e7877530',
    source: 'helm-service',
    specversion: '1.0',
    time: '2021-10-13T10:51:07.960Z',
    type: 'sh.keptn.event.deployment.finished',
    shkeptncontext: keptnContext,
    shkeptnspecversion: '0.2.3',
    triggeredid: '3e82e52c-a129-43c3-a9b4-e9dce6aeaf4e',
  },
  {
    data: {
      labels: {
        DtCreds: 'dynatrace',
      },
      project: 'sockshop',
      service: 'carts',
      stage: 'staging',
      status: 'succeeded',
      temporaryData: {
        distributor: {
          subscriptionID: '',
        },
      },
    },
    id: '19356816-2df9-46e4-8ef1-aac7cfbf43a6',
    source: 'helm-service',
    specversion: '1.0',
    time: '2021-10-13T10:49:33.714Z',
    type: 'sh.keptn.event.deployment.started',
    shkeptncontext: keptnContext,
    shkeptnspecversion: '0.2.3',
    triggeredid: '3e82e52c-a129-43c3-a9b4-e9dce6aeaf4e',
  },
  {
    data: {
      configurationChange: {
        values: {
          image: 'docker.io/keptnexamples/carts:0.12.1',
        },
      },
      deployment: {
        deploymentNames: ['direct'],
        deploymentURIsLocal: null,
        deploymentURIsPublic: ['http://carts.sockshop-dev.35.192.209.116.nip.io:80'],
        deploymentstrategy: 'blue_green_service',
      },
      evaluation: {
        indicatorResults: null,
        responses: ['Redirecting to https://keptn.sh/\n'],
        result: 'pass',
        score: 0,
        sloFileContent: '',
        timeEnd: '2021-10-13T10:47:11Z',
        timeStart: '2021-10-13T10:46:42Z',
      },
      labels: {
        DtCreds: 'dynatrace',
      },
      message: '',
      project: 'sockshop',
      result: '',
      service: 'carts',
      stage: 'staging',
      status: '',
      temporaryData: {
        distributor: {
          subscriptionID: '',
        },
      },
      test: {
        end: '2021-10-13T10:47:11Z',
        start: '2021-10-13T10:46:42Z',
      },
    },
    id: '3e82e52c-a129-43c3-a9b4-e9dce6aeaf4e',
    source: 'shipyard-controller',
    specversion: '1.0',
    time: '2021-10-13T10:49:33.711Z',
    type: 'sh.keptn.event.deployment.triggered',
    shkeptncontext: keptnContext,
    shkeptnspecversion: '0.2.3',
  },
];

const approvalFinishedTraceStaging = {
  data: {
    labels: {
      DtCreds: 'dynatrace',
    },
    message: '',
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
  id: '0b31d407-db29-4f80-82eb-e30f45a74c1f',
  source: 'https://github.com/keptn/keptn/bridge#approval.finished',
  specversion: '1.0',
  time: '2021-10-13T10:57:44.251Z',
  type: 'sh.keptn.event.approval.finished',
  shkeptncontext: keptnContext,
  shkeptnspecversion: '0.2.3',
  triggeredid: '7d93986c-37bb-41b0-8e97-2e8ee7eee6c0',
};

const approvalTriggeredTraceStaging = {
  data: {
    approval: {
      pass: 'manual',
      warning: 'manual',
    },
    configurationChange: {
      values: {
        image: 'docker.io/keptnexamples/carts:0.12.1',
      },
    },
    deployment: {
      deploymentNames: ['canary', 'direct'],
      deploymentURIsLocal: ['http://carts.sockshop-staging:80'],
      deploymentURIsPublic: [
        'http://carts.sockshop-staging.35.192.209.116.nip.io:80',
        'http://carts.sockshop-dev.35.192.209.116.nip.io:80',
      ],
      deploymentstrategy: 'duplicate',
    },
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
            value: 304.2952915485157,
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
      responses: ['Redirecting to https://keptn.sh/\n'],
      result: 'pass',
      score: 100,
      sloFileContent:
        'LS0tCnNwZWNfdmVyc2lvbjogIjEuMCIKY29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciCiAgY29tcGFyZV93aXRoOiAic2luZ2xlX3Jlc3VsdCIKICBpbmNsdWRlX3Jlc3VsdF93aXRoX3Njb3JlOiAicGFzcyIKICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxCmZpbHRlcjoKb2JqZWN0aXZlczoKICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1IgogICAgZGlzcGxheU5hbWU6ICJSZXNwb25zZSB0aW1lIFA5NSIKICAgIGtleV9zbGk6IGZhbHNlCiAgICBwYXNzOiAgICAgICAgICAgICAjIHBhc3MgaWYgKHJlbGF0aXZlIGNoYW5nZSA8PSAxMCUgQU5EIGFic29sdXRlIHZhbHVlIGlzIDwgNjAwbXMpCiAgICAgIC0gY3JpdGVyaWE6CiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykKICAgICAgICAgIC0gIjw2MDAiICAgICMgYWJzb2x1dGUgdmFsdWVzIG9ubHkgcmVxdWlyZSBhIGxvZ2ljYWwgb3BlcmF0b3IKICAgIHdhcm5pbmc6ICAgICAgICAgICMgaWYgdGhlIHJlc3BvbnNlIHRpbWUgaXMgYmVsb3cgODAwbXMsIHRoZSByZXN1bHQgc2hvdWxkIGJlIGEgd2FybmluZwogICAgICAtIGNyaXRlcmlhOgogICAgICAgICAgLSAiPD04MDAiCiAgICB3ZWlnaHQ6IDEKdG90YWxfc2NvcmU6CiAgcGFzczogIjkwJSIKICB3YXJuaW5nOiAiNzUlIg==',
      timeEnd: '2021-10-13T10:53:29Z',
      timeStart: '2021-10-13T10:51:08Z',
    },
    labels: {
      DtCreds: 'dynatrace',
    },
    message: '',
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
    test: {
      end: '2021-10-13T10:53:29Z',
      start: '2021-10-13T10:51:08Z',
    },
  },
  id: '7d93986c-37bb-41b0-8e97-2e8ee7eee6c0',
  source: 'shipyard-controller',
  specversion: '1.0',
  time: '2021-10-13T10:54:43.311Z',
  type: 'sh.keptn.event.approval.triggered',
  shkeptncontext: keptnContext,
  shkeptnspecversion: '0.2.3',
};

const approvalStartedTracesStaging = [
  {
    data: {
      labels: {
        DtCreds: 'dynatrace',
      },
      project: 'sockshop',
      service: 'carts',
      stage: 'staging',
      status: 'succeeded',
      temporaryData: {
        distributor: {
          subscriptionID: '',
        },
      },
    },
    id: '851442f9-7201-44da-8577-2ab1bf08c476',
    source: 'approval-service',
    specversion: '1.0',
    time: '2021-10-13T10:54:43.315Z',
    type: 'sh.keptn.event.approval.started',
    shkeptncontext: keptnContext,
    shkeptnspecversion: '0.2.3',
    triggeredid: '7d93986c-37bb-41b0-8e97-2e8ee7eee6c0',
  },
  approvalTriggeredTraceStaging,
];

const approvalTracesStaging = [...[approvalFinishedTraceStaging], ...approvalStartedTracesStaging];

const evaluationAndSliTracesStaging = [
  {
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
              value: 304.2952915485157,
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
        timeEnd: '2021-10-13T10:53:29Z',
        timeStart: '2021-10-13T10:51:08Z',
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
    id: '93c2eba9-b77c-4976-b079-29a0188d86ef',
    source: 'lighthouse-service',
    specversion: '1.0',
    time: '2021-10-13T10:54:43.112Z',
    type: 'sh.keptn.event.evaluation.finished',
    shkeptncontext: keptnContext,
    shkeptnspecversion: '0.2.3',
    triggeredid: '1cc9c272-721a-43de-98f6-9eceae484cf5',
  },
  {
    data: {
      'get-sli': {
        end: '2021-10-13T10:53:29Z',
        indicatorValues: [
          {
            metric: 'response_time_p95',
            success: true,
            value: 304.2952915485157,
          },
        ],
        start: '2021-10-13T10:51:08Z',
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
    id: '6a890b21-49c0-45bd-8eaf-72672053aa08',
    source: 'dynatrace-service',
    specversion: '1.0',
    time: '2021-10-13T10:54:39.911Z',
    type: 'sh.keptn.event.get-sli.finished',
    shkeptncontext: keptnContext,
    shkeptnspecversion: '0.2.1',
    triggeredid: '5b8b3a7b-b4da-4c1c-addd-459e2ad3efe4',
  },
  {
    data: {
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
    id: '5085265f-174f-4535-aab3-bf4233b7ce46',
    source: 'dynatrace-service',
    specversion: '1.0',
    time: '2021-10-13T10:53:32.926Z',
    type: 'sh.keptn.event.get-sli.started',
    shkeptncontext: keptnContext,
    shkeptnspecversion: '0.2.1',
    triggeredid: '5b8b3a7b-b4da-4c1c-addd-459e2ad3efe4',
  },
  {
    data: {
      deployment: 'canary',
      'get-sli': {
        end: '2021-10-13T10:53:29Z',
        indicators: ['response_time_p95'],
        sliProvider: 'dynatrace',
        start: '2021-10-13T10:51:08Z',
      },
      labels: {
        DtCreds: 'dynatrace',
      },
      project: 'sockshop',
      service: 'carts',
      stage: 'staging',
      temporaryData: {
        distributor: {
          subscriptionID: '',
        },
      },
    },
    id: '5b8b3a7b-b4da-4c1c-addd-459e2ad3efe4',
    source: 'lighthouse-service',
    specversion: '1.0',
    time: '2021-10-13T10:53:32.923Z',
    type: 'sh.keptn.event.get-sli.triggered',
    shkeptncontext: keptnContext,
    shkeptnspecversion: '0.2.3',
  },
  {
    data: {
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
    id: 'b8d1914c-6bc1-4870-8da7-c04b0681197b',
    source: 'lighthouse-service',
    specversion: '1.0',
    time: '2021-10-13T10:53:30.109Z',
    type: 'sh.keptn.event.evaluation.started',
    shkeptncontext: keptnContext,
    shkeptnspecversion: '0.2.3',
    triggeredid: '1cc9c272-721a-43de-98f6-9eceae484cf5',
  },
  {
    data: {
      configurationChange: {
        values: {
          image: 'docker.io/keptnexamples/carts:0.12.1',
        },
      },
      deployment: {
        deploymentNames: ['canary', 'direct'],
        deploymentURIsLocal: ['http://carts.sockshop-staging:80'],
        deploymentURIsPublic: [
          'http://carts.sockshop-staging.35.192.209.116.nip.io:80',
          'http://carts.sockshop-dev.35.192.209.116.nip.io:80',
        ],
        deploymentstrategy: 'duplicate',
      },
      evaluation: {
        indicatorResults: null,
        responses: ['Redirecting to https://keptn.sh/\n'],
        result: 'pass',
        score: 0,
        sloFileContent: '',
        timeEnd: '2021-10-13T10:47:11Z',
        timeStart: '2021-10-13T10:46:42Z',
      },
      labels: {
        DtCreds: 'dynatrace',
      },
      message: '',
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
      test: {
        end: '2021-10-13T10:53:29Z',
        start: '2021-10-13T10:51:08Z',
      },
    },
    id: '1cc9c272-721a-43de-98f6-9eceae484cf5',
    source: 'shipyard-controller',
    specversion: '1.0',
    time: '2021-10-13T10:53:30.106Z',
    type: 'sh.keptn.event.evaluation.triggered',
    shkeptncontext: keptnContext,
    shkeptnspecversion: '0.2.3',
  },
];

const testTracesStaging = [
  {
    data: {
      labels: {
        DtCreds: 'dynatrace',
      },
      message:
        'Tests for performance with status = true.Project: sockshop, Service: carts, Stage: staging, TestStrategy: performance, Context: 2c0e568b-8bd3-4726-a188-e528423813ed',
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
      test: {
        end: '2021-10-13T10:53:29Z',
        start: '2021-10-13T10:51:08Z',
      },
    },
    id: '108a2efa-c474-49d5-a132-99962fb0a39b',
    source: 'jmeter-service',
    specversion: '1.0',
    time: '2021-10-13T10:53:29.967Z',
    type: 'sh.keptn.event.test.finished',
    shkeptncontext: keptnContext,
    shkeptnspecversion: '0.2.3',
    triggeredid: '4ac95fe7-bc8f-407d-9665-c32a3eb74791',
  },
  {
    data: {
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
    id: '9a7ceb5d-6816-4a18-a22f-3924e85504f9',
    source: 'jmeter-service',
    specversion: '1.0',
    time: '2021-10-13T10:51:08.213Z',
    type: 'sh.keptn.event.test.started',
    shkeptncontext: keptnContext,
    shkeptnspecversion: '0.2.3',
    triggeredid: '4ac95fe7-bc8f-407d-9665-c32a3eb74791',
  },
  {
    data: {
      configurationChange: {
        values: {
          image: 'docker.io/keptnexamples/carts:0.12.1',
        },
      },
      deployment: {
        deploymentNames: ['canary', 'direct'],
        deploymentURIsLocal: ['http://carts.sockshop-staging:80'],
        deploymentURIsPublic: [
          'http://carts.sockshop-staging.35.192.209.116.nip.io:80',
          'http://carts.sockshop-dev.35.192.209.116.nip.io:80',
        ],
        deploymentstrategy: 'duplicate',
      },
      evaluation: {
        indicatorResults: null,
        responses: ['Redirecting to https://keptn.sh/\n'],
        result: 'pass',
        score: 0,
        sloFileContent: '',
        timeEnd: '2021-10-13T10:47:11Z',
        timeStart: '2021-10-13T10:46:42Z',
      },
      labels: {
        DtCreds: 'dynatrace',
      },
      message: '',
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
      test: {
        end: '2021-10-13T10:47:11Z',
        start: '2021-10-13T10:46:42Z',
        teststrategy: 'performance',
      },
    },
    id: '4ac95fe7-bc8f-407d-9665-c32a3eb74791',
    source: 'shipyard-controller',
    specversion: '1.0',
    time: '2021-10-13T10:51:08.208Z',
    type: 'sh.keptn.event.test.triggered',
    shkeptncontext: keptnContext,
    shkeptnspecversion: '0.2.3',
  },
];

const deliveryFinishedDev = {
  data: {
    labels: {
      DtCreds: 'dynatrace',
    },
    message: 'Finished release',
    project: 'sockshop',
    result: 'pass',
    service: 'carts',
    stage: 'dev',
    status: 'succeeded',
    temporaryData: {
      distributor: {
        subscriptionID: '',
      },
    },
  },
  id: '861c0880-0555-44e5-b7c3-478e4a7edae4',
  source: 'shipyard-controller',
  specversion: '1.0',
  time: '2021-10-13T10:49:30.005Z',
  type: 'sh.keptn.event.dev.delivery.finished',
  shkeptncontext: keptnContext,
  shkeptnspecversion: '0.2.3',
  triggeredid: '08e89fdb-02db-4fc7-a5fd-386d03e1c4a9',
};

const releaseTracesDev = [
  {
    data: {
      labels: {
        DtCreds: 'dynatrace',
      },
      message: 'Finished release',
      project: 'sockshop',
      result: 'pass',
      service: 'carts',
      stage: 'dev',
      status: 'succeeded',
      temporaryData: {
        distributor: {
          subscriptionID: '',
        },
      },
    },
    id: '277ad08f-b94b-4853-a50f-706852bbde2f',
    source: 'helm-service',
    specversion: '1.0',
    time: '2021-10-13T10:49:27.724Z',
    type: 'sh.keptn.event.release.finished',
    shkeptncontext: keptnContext,
    shkeptnspecversion: '0.2.3',
    triggeredid: '26dbcd4e-baa7-4d57-9dfe-48b03b0eb891',
  },
  {
    data: {
      labels: {
        DtCreds: 'dynatrace',
      },
      project: 'sockshop',
      service: 'carts',
      stage: 'dev',
      status: 'succeeded',
      temporaryData: {
        distributor: {
          subscriptionID: '',
        },
      },
    },
    id: '40c8756b-c44a-4338-ae5e-cf7638bfaa39',
    source: 'helm-service',
    specversion: '1.0',
    time: '2021-10-13T10:49:27.719Z',
    type: 'sh.keptn.event.release.started',
    shkeptncontext: keptnContext,
    shkeptnspecversion: '0.2.3',
    triggeredid: '26dbcd4e-baa7-4d57-9dfe-48b03b0eb891',
  },
  {
    data: {
      configurationChange: {
        values: {
          image: 'docker.io/keptnexamples/carts:0.12.1',
        },
      },
      deployment: {
        deploymentNames: ['direct'],
        deploymentURIsLocal: ['http://carts.sockshop-dev:80'],
        deploymentURIsPublic: ['http://carts.sockshop-dev.35.192.209.116.nip.io:80'],
        deploymentstrategy: 'direct',
      },
      evaluation: {
        indicatorResults: null,
        responses: ['Redirecting to https://keptn.sh/\n'],
        result: 'pass',
        score: 0,
        sloFileContent: '',
        timeEnd: '2021-10-13T10:47:11Z',
        timeStart: '2021-10-13T10:46:42Z',
      },
      labels: {
        DtCreds: 'dynatrace',
      },
      message: '',
      project: 'sockshop',
      release: null,
      result: 'pass',
      service: 'carts',
      stage: 'dev',
      status: 'succeeded',
      temporaryData: {
        distributor: {
          subscriptionID: '',
        },
      },
      test: {
        end: '2021-10-13T10:47:11Z',
        start: '2021-10-13T10:46:42Z',
      },
    },
    id: '26dbcd4e-baa7-4d57-9dfe-48b03b0eb891',
    source: 'shipyard-controller',
    specversion: '1.0',
    time: '2021-10-13T10:49:27.714Z',
    type: 'sh.keptn.event.release.triggered',
    shkeptncontext: keptnContext,
    shkeptnspecversion: '0.2.3',
  },
];

const evaluationAndSliTracesDev = [
  {
    data: {
      evaluation: {
        indicatorResults: null,
        result: 'pass',
        score: 0,
        sloFileContent: '',
        timeEnd: '2021-10-13T10:47:11Z',
        timeStart: '2021-10-13T10:46:42Z',
      },
      labels: {
        DtCreds: 'dynatrace',
      },
      message: 'no evaluation performed by lighthouse because no SLO file configured for project sockshop',
      project: 'sockshop',
      result: 'pass',
      service: 'carts',
      stage: 'dev',
      status: 'succeeded',
      temporaryData: {
        distributor: {
          subscriptionID: '',
        },
      },
    },
    id: 'b286a8ea-1b11-4106-b892-fade06245bcf',
    source: 'lighthouse-service',
    specversion: '1.0',
    time: '2021-10-13T10:49:27.606Z',
    type: 'sh.keptn.event.evaluation.finished',
    shkeptncontext: keptnContext,
    shkeptnspecversion: '0.2.3',
    triggeredid: '5a548dc5-b29f-4f2d-a2be-4d4534de523f',
  },
  {
    data: {
      labels: {
        DtCreds: 'dynatrace',
      },
      project: 'sockshop',
      result: 'pass',
      service: 'carts',
      stage: 'dev',
      status: 'succeeded',
      temporaryData: {
        distributor: {
          subscriptionID: '',
        },
      },
    },
    id: 'x286a8ea-1b11-4106-b892-fade06245bcf',
    source: 'webhook-service',
    specversion: '1.0',
    time: '2021-10-13T10:49:27.607Z',
    type: 'sh.keptn.event.evaluation.started',
    shkeptncontext: keptnContext,
    shkeptnspecversion: '0.2.3',
    triggeredid: 'b286a8ea-1b11-4106-b892-fade06245bcf',
  },
  {
    data: {
      'get-sli': {
        end: '2021-10-13T10:47:11Z',
        indicatorValues: [
          {
            message: "Couldn't retrieve any SLI Results",
            metric: 'no metric',
            success: false,
            value: 0,
          },
        ],
        start: '2021-10-13T10:46:42Z',
      },
      labels: {
        DtCreds: 'dynatrace',
      },
      project: 'sockshop',
      result: 'pass',
      service: 'carts',
      stage: 'dev',
      status: 'succeeded',
      temporaryData: {
        distributor: {
          subscriptionID: '',
        },
      },
    },
    id: '892b6b90-0681-45ac-9213-988457c2b9a2',
    source: 'dynatrace-service',
    specversion: '1.0',
    time: '2021-10-13T10:49:25.806Z',
    type: 'sh.keptn.event.get-sli.finished',
    shkeptncontext: keptnContext,
    shkeptnspecversion: '0.2.1',
    triggeredid: 'd746bba3-3077-42fa-9b55-6adae66b2d80',
  },
  {
    data: {
      project: 'sockshop',
      result: 'pass',
      service: 'carts',
      stage: 'dev',
      status: 'succeeded',
      temporaryData: {
        distributor: {
          subscriptionID: '',
        },
      },
    },
    id: '73f6a0ba-a953-48fc-8633-08fc73ab1ccd',
    source: 'dynatrace-service',
    specversion: '1.0',
    time: '2021-10-13T10:47:14.732Z',
    type: 'sh.keptn.event.get-sli.started',
    shkeptncontext: keptnContext,
    shkeptnspecversion: '0.2.1',
    triggeredid: 'd746bba3-3077-42fa-9b55-6adae66b2d80',
  },
  {
    data: {
      deployment: 'direct',
      'get-sli': {
        end: '2021-10-13T10:47:11Z',
        sliProvider: 'dynatrace',
        start: '2021-10-13T10:46:42Z',
      },
      project: 'sockshop',
      service: 'carts',
      stage: 'dev',
      temporaryData: {
        distributor: {
          subscriptionID: '',
        },
      },
    },
    id: 'd746bba3-3077-42fa-9b55-6adae66b2d80',
    source: 'lighthouse-service',
    specversion: '1.0',
    time: '2021-10-13T10:47:14.727Z',
    type: 'sh.keptn.event.get-sli.triggered',
    shkeptncontext: keptnContext,
    shkeptnspecversion: '0.2.3',
  },
  {
    data: {
      project: 'sockshop',
      result: 'pass',
      service: 'carts',
      stage: 'dev',
      status: 'succeeded',
      temporaryData: {
        distributor: {
          subscriptionID: '',
        },
      },
    },
    id: '25929ca4-8c87-4a09-969a-454df32a9a1a',
    source: 'lighthouse-service',
    specversion: '1.0',
    time: '2021-10-13T10:47:11.520Z',
    type: 'sh.keptn.event.evaluation.started',
    shkeptncontext: keptnContext,
    shkeptnspecversion: '0.2.3',
    triggeredid: '5a548dc5-b29f-4f2d-a2be-4d4534de523f',
  },
  {
    data: {
      configurationChange: {
        values: {
          image: 'docker.io/keptnexamples/carts:0.12.1',
        },
      },
      deployment: {
        deploymentNames: ['direct'],
        deploymentURIsLocal: ['http://carts.sockshop-dev:80'],
        deploymentURIsPublic: ['http://carts.sockshop-dev.35.192.209.116.nip.io:80'],
        deploymentstrategy: 'direct',
      },
      evaluation: null,
      message: '',
      project: 'sockshop',
      result: 'pass',
      service: 'carts',
      stage: 'dev',
      status: 'succeeded',
      temporaryData: {
        distributor: {
          subscriptionID: '',
        },
      },
      test: {
        end: '2021-10-13T10:47:11Z',
        start: '2021-10-13T10:46:42Z',
      },
    },
    id: '5a548dc5-b29f-4f2d-a2be-4d4534de523f',
    source: 'shipyard-controller',
    specversion: '1.0',
    time: '2021-10-13T10:47:11.509Z',
    type: 'sh.keptn.event.evaluation.triggered',
    shkeptncontext: keptnContext,
    shkeptnspecversion: '0.2.3',
  },
];

const testTracesDev = [
  {
    data: {
      message:
        'Tests for functional with status = true.Project: sockshop, Service: carts, Stage: dev, TestStrategy: functional, Context: 2c0e568b-8bd3-4726-a188-e528423813ed',
      project: 'sockshop',
      result: 'pass',
      service: 'carts',
      stage: 'dev',
      status: 'succeeded',
      temporaryData: {
        distributor: {
          subscriptionID: '',
        },
      },
      test: {
        end: '2021-10-13T10:47:11Z',
        start: '2021-10-13T10:46:42Z',
      },
    },
    id: 'f7830cc9-1889-4aca-9194-a10f7ce79938',
    source: 'jmeter-service',
    specversion: '1.0',
    time: '2021-10-13T10:47:11.404Z',
    type: 'sh.keptn.event.test.finished',
    shkeptncontext: keptnContext,
    shkeptnspecversion: '0.2.3',
    triggeredid: '29edc18b-538d-4184-8b13-25fb0aaf5ad4',
  },
  {
    data: {
      project: 'sockshop',
      result: 'pass',
      service: 'carts',
      stage: 'dev',
      status: 'succeeded',
      temporaryData: {
        distributor: {
          subscriptionID: '',
        },
      },
    },
    id: 'c2889bef-37dd-466f-b539-460559cd434b',
    source: 'jmeter-service',
    specversion: '1.0',
    time: '2021-10-13T10:46:42.218Z',
    type: 'sh.keptn.event.test.started',
    shkeptncontext: keptnContext,
    shkeptnspecversion: '0.2.3',
    triggeredid: '29edc18b-538d-4184-8b13-25fb0aaf5ad4',
  },
  {
    data: {
      configurationChange: {
        values: {
          image: 'docker.io/keptnexamples/carts:0.12.1',
        },
      },
      deployment: {
        deploymentNames: ['direct'],
        deploymentURIsLocal: ['http://carts.sockshop-dev:80'],
        deploymentURIsPublic: ['http://carts.sockshop-dev.35.192.209.116.nip.io:80'],
        deploymentstrategy: 'direct',
      },
      message: '',
      project: 'sockshop',
      result: 'pass',
      service: 'carts',
      stage: 'dev',
      status: 'succeeded',
      temporaryData: {
        distributor: {
          subscriptionID: '',
        },
      },
      test: {
        teststrategy: 'functional',
      },
    },
    id: '29edc18b-538d-4184-8b13-25fb0aaf5ad4',
    source: 'shipyard-controller',
    specversion: '1.0',
    time: '2021-10-13T10:46:42.214Z',
    type: 'sh.keptn.event.test.triggered',
    shkeptncontext: keptnContext,
    shkeptnspecversion: '0.2.3',
  },
];

const deploymentTracesDev = [
  {
    data: {
      deployment: {
        deploymentNames: ['direct'],
        deploymentURIsLocal: ['http://carts.sockshop-dev:80'],
        deploymentURIsPublic: ['http://carts.sockshop-dev.35.192.209.116.nip.io:80'],
        deploymentstrategy: 'direct',
      },
      message: 'Successfully deployed',
      project: 'sockshop',
      result: 'pass',
      service: 'carts',
      stage: 'dev',
      status: 'succeeded',
      temporaryData: {
        distributor: {
          subscriptionID: '',
        },
      },
    },
    id: '8b65ab52-f1bb-4ecf-8559-48854acfa60d',
    source: 'helm-service',
    specversion: '1.0',
    time: '2021-10-13T10:46:41.861Z',
    type: 'sh.keptn.event.deployment.finished',
    shkeptncontext: keptnContext,
    shkeptnspecversion: '0.2.3',
    triggeredid: '89265987-90b9-4554-8d13-98deb7b44f9d',
  },
  {
    data: {
      project: 'sockshop',
      service: 'carts',
      stage: 'dev',
      status: 'succeeded',
      temporaryData: {
        distributor: {
          subscriptionID: '',
        },
      },
    },
    id: '40b7572f-2d63-4f32-8db7-45d6a99bf106',
    source: 'helm-service',
    specversion: '1.0',
    time: '2021-10-13T10:45:06.960Z',
    type: 'sh.keptn.event.deployment.started',
    shkeptncontext: keptnContext,
    shkeptnspecversion: '0.2.3',
    triggeredid: '89265987-90b9-4554-8d13-98deb7b44f9d',
  },
  {
    data: {
      configurationChange: {
        values: {
          image: 'docker.io/keptnexamples/carts:0.12.1',
        },
      },
      deployment: {
        deploymentURIsLocal: null,
        deploymentstrategy: 'direct',
      },
      message: '',
      project: 'sockshop',
      result: '',
      service: 'carts',
      stage: 'dev',
      status: '',
      temporaryData: {
        distributor: {
          subscriptionID: '',
        },
      },
    },
    id: '89265987-90b9-4554-8d13-98deb7b44f9d',
    source: 'shipyard-controller',
    specversion: '1.0',
    time: '2021-10-13T10:45:06.901Z',
    type: 'sh.keptn.event.deployment.triggered',
    shkeptncontext: keptnContext,
    shkeptnspecversion: '0.2.3',
  },
];

const deliveryTriggeredDev = {
  data: {
    configurationChange: {
      values: {
        image: 'docker.io/keptnexamples/carts:0.12.1',
      },
    },
    deployment: {
      deploymentURIsLocal: null,
      deploymentstrategy: '',
    },
    project: 'sockshop',
    service: 'carts',
    stage: 'dev',
    temporaryData: {
      distributor: {
        subscriptionID: '',
      },
    },
  },
  id: '08e89fdb-02db-4fc7-a5fd-386d03e1c4a9',
  source: 'https://github.com/keptn/keptn/cli#configuration-change',
  specversion: '1.0',
  time: '2021-10-13T10:45:03.780Z',
  type: 'sh.keptn.event.dev.delivery.triggered',
  shkeptncontext: keptnContext,
  shkeptnspecversion: '0.2.3',
};

const deploymentTracesResponseMock = {
  events: [
    ...[deliveryFinishedProduction],
    ...releaseTracesProduction,
    ...deploymentTracesProduction,
    ...[deliveryTriggeredProduction],
    ...[deliveryFinishedStaging],
    ...releaseTracesStaging,
    ...approvalTracesStaging,
    ...evaluationAndSliTracesStaging,
    ...testTracesStaging,
    ...deploymentTracesStaging,
    ...[deliveryTriggeredStaging],
    ...[deliveryFinishedDev],
    ...releaseTracesDev,
    ...evaluationAndSliTracesDev,
    ...testTracesDev,
    ...deploymentTracesDev,
    ...[deliveryTriggeredDev],
  ],
  pageSize: 100,
  totalCount: 47,
};

const deploymentTracesWithPendingApprovalResponseMock = {
  events: [
    ...approvalStartedTracesStaging,
    ...deploymentTracesStaging,
    ...[deliveryTriggeredStaging],
    ...[deliveryFinishedDev],
    ...releaseTracesDev,
    ...evaluationAndSliTracesDev,
    ...testTracesDev,
    ...deploymentTracesDev,
    ...[deliveryTriggeredDev],
  ],
  pageSize: 100,
  totalCount: 20,
};

const evaluationFinishedProductionResponse = {
  events: [
    {
      data: {
        evaluation: {
          indicatorResults: null,
          result: 'pass',
          score: 0,
          sloFileContent: '',
          timeEnd: '2021-11-08T14:02:14.147Z',
          timeStart: '2021-11-08T13:57:14.147Z',
        },
        labels: {
          DtCreds: 'dynatrace',
        },
        message: 'no evaluation performed by lighthouse because no SLO file configured for project sockshop',
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
      id: 'e2800526-c0f9-4817-9aec-9a3037d42837',
      source: 'lighthouse-service',
      specversion: '1.0',
      time: '2021-11-08T14:02:30.265Z',
      type: 'sh.keptn.event.evaluation.finished',
      shkeptncontext: '01b41825-899b-48d0-bfc2-4e92d7c9bf29',
      shkeptnspecversion: '0.2.3',
      triggeredid: '4f01fc62-06b5-4651-b090-2ff672c364b1',
    },
  ],
};

const evaluationFinishedStagingResponse = {
  events: [
    {
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
                value: 304.2952915485157,
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
          timeEnd: '2021-10-13T10:53:29Z',
          timeStart: '2021-10-13T10:51:08Z',
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
      id: '93c2eba9-b77c-4976-b079-29a0188d86ef',
      source: 'lighthouse-service',
      specversion: '1.0',
      time: '2021-10-13T10:54:43.112Z',
      type: 'sh.keptn.event.evaluation.finished',
      shkeptncontext: '2c0e568b-8bd3-4726-a188-e528423813ee',
      shkeptnspecversion: '0.2.3',
      triggeredid: '1cc9c272-721a-43de-98f6-9eceae484cf5',
    },
  ],
};

const latestFinishedDeployments = {
  events: [
    {
      ...defaultDeploymentFinishedTrace,
      data: {
        service: 'carts',
        stage: 'dev',
        project: 'sockshop',
        deployment: {
          ...defaultDeploymentData,
          deploymentURIsPublic: ['http://dev-carts.com'],
        },
        result: 'pass',
      },
      id: DevCartsDeploymentFinished.eventId,
      shkeptncontext: DevCartsDeploymentFinished.keptnContext,
    },
    {
      ...defaultDeploymentFinishedTrace,
      data: {
        service: 'carts-db',
        stage: 'dev',
        project: 'sockshop',
        deployment: {
          ...defaultDeploymentData,
          deploymentURIsPublic: ['http://dev-carts-db.com'],
        },
        result: 'pass',
      },
      id: DevCartsDbDeploymentFinished.eventId,
      shkeptncontext: DevCartsDbDeploymentFinished.keptnContext,
    },
    {
      ...defaultDeploymentFinishedTrace,
      data: {
        service: 'carts',
        stage: 'staging',
        project: 'sockshop',
        deployment: {
          ...defaultDeploymentData,
          deploymentURIsPublic: ['http://staging-carts.com'],
        },
        result: 'pass',
      },
      id: StagingCartsDeploymentFinished.eventId,
      shkeptncontext: StagingCartsDeploymentFinished.keptnContext,
    },
    {
      ...defaultDeploymentFinishedTrace,
      data: {
        service: 'carts-db',
        stage: 'staging',
        project: 'sockshop',
        deployment: {
          ...defaultDeploymentData,
          deploymentURIsPublic: ['http://staging-carts-db.com'],
        },
        result: 'pass',
      },
      id: StagingCartsDbDeploymentFinished.eventId,
      shkeptncontext: StagingCartsDbDeploymentFinished.keptnContext,
    },
    {
      ...defaultDeploymentFinishedTrace,
      data: {
        service: 'carts',
        stage: 'production',
        project: 'sockshop',
        deployment: {
          ...defaultDeploymentData,
          deploymentURIsPublic: ['http://production-carts.com'],
        },
        result: 'pass',
      },
      id: ProductionCartsDeploymentFinished.eventId,
      shkeptncontext: ProductionCartsDeploymentFinished.keptnContext,
    },
    {
      ...defaultDeploymentFinishedTrace,
      data: {
        service: 'carts-db',
        stage: 'production',
        project: 'sockshop',
        deployment: {
          ...defaultDeploymentData,
          deploymentURIsPublic: ['http://production-carts-db.com'],
        },
        result: 'pass',
      },
      id: ProductionCartsDbDeploymentFinished.eventId,
      shkeptncontext: ProductionCartsDbDeploymentFinished.keptnContext,
    },
  ],
  pageSize: 100,
  totalCount: 6,
};

const latestFinishedEvaluations = {
  events: [
    // {
    //   ...defaultEvaluationFinishedTrace,
    //   data: {
    //     service: 'carts',
    //     stage: 'dev',
    //     project: 'sockshop',
    //     evaluation: defaultEvaluationData,
    //   },
    //   id: DevCartsEvaluationFinished.eventId,
    //   shkeptncontext: DevCartsEvaluationFinished.keptnContext,
    // }, //evaluation is not part of latest Sequence
    {
      ...defaultEvaluationFinishedTrace,
      data: {
        service: 'carts',
        stage: 'staging',
        project: 'sockshop',
        evaluation: defaultEvaluationData,
        result: 'pass',
      },
      id: StagingCartsEvaluationFinished.eventId,
      shkeptncontext: StagingCartsEvaluationFinished.keptnContext,
    },
  ],
  pageSize: 100,
  totalCount: 1,
};

const openApprovalsResponse = {
  events: [approvalTriggeredTraceStaging],
  pageSize: 100,
  totalCount: 1,
};

const defaultDeploymentTriggeredEvent = {
  data: {
    configurationChange: {
      values: {
        image: 'docker.io/keptnexamples/carts:0.12.3',
      },
    },
    project: 'sockshop',
    service: 'carts',
    stage: 'dev',
  },
  id: '8b65ab52-f1bb-4ecf-8559-48854acfa60d',
  source: 'helm-service',
  specversion: '1.0',
  time: '2021-10-13T10:46:41.861Z',
  type: 'sh.keptn.event.deployment.triggered',
  shkeptncontext: keptnContext,
  shkeptnspecversion: '0.2.3',
  triggeredid: '89265987-90b9-4554-8d13-98deb7b44f9d',
};

const defaultDeploymentStartedEvent = {
  ...defaultDeploymentTriggeredEvent,
  type: 'sh.keptn.event.deployment.started',
};

const intersectDeploymentTriggeredResponse = {
  events: [
    {
      ...defaultDeploymentTriggeredEvent,
      data: {
        ...defaultDeploymentTriggeredEvent.data,
        project: 'sockshop',
        service: 'carts',
        stage: 'dev',
      },
    },
    {
      ...defaultDeploymentTriggeredEvent,
      data: {
        ...defaultDeploymentTriggeredEvent.data,
        project: 'sockshop',
        service: 'carts-db',
        stage: 'dev',
      },
    },
    {
      ...defaultDeploymentTriggeredEvent,
      data: {
        ...defaultDeploymentTriggeredEvent.data,
        project: 'sockshop',
        service: 'carts',
        stage: 'staging',
      },
    },
    {
      ...defaultDeploymentTriggeredEvent,
      data: {
        ...defaultDeploymentTriggeredEvent.data,
        project: 'sockshop',
        service: 'carts-db',
        stage: 'staging',
      },
    },
  ],
};

const intersectDeploymentStartedResponse = {
  events: [
    {
      ...defaultDeploymentStartedEvent,
      data: {
        ...defaultDeploymentStartedEvent.data,
        project: 'sockshop',
        service: 'carts',
        stage: 'dev',
      },
    },
    {
      ...defaultDeploymentStartedEvent,
      data: {
        ...defaultDeploymentStartedEvent.data,
        project: 'sockshop',
        service: 'carts-db',
        stage: 'dev',
      },
    },
    {
      ...defaultDeploymentStartedEvent,
      data: {
        ...defaultDeploymentStartedEvent.data,
        project: 'sockshop',
        service: 'carts',
        stage: 'staging',
      },
    },
    {
      ...defaultDeploymentStartedEvent,
      data: {
        ...defaultDeploymentStartedEvent.data,
        project: 'sockshop',
        service: 'carts-db',
        stage: 'staging',
      },
    },
  ],
};

const intersectDeploymentFinishedResponse = {
  events: [
    {
      ...defaultDeploymentFinishedTrace,
      data: {
        ...defaultDeploymentFinishedTrace.data,
        project: 'sockshop',
        service: 'carts',
        stage: 'dev',
      },
    },
    {
      ...defaultDeploymentFinishedTrace,
      data: {
        ...defaultDeploymentFinishedTrace.data,
        project: 'sockshop',
        service: 'carts-db',
        stage: 'dev',
      },
    },
    {
      ...defaultDeploymentFinishedTrace,
      data: {
        ...defaultDeploymentFinishedTrace.data,
        project: 'sockshop',
        service: 'carts',
        stage: 'staging',
      },
    },
  ],
};

export { deploymentTracesResponseMock as DeploymentTracesResponseMock };
export { deploymentTracesWithPendingApprovalResponseMock as DeploymentTracesWithPendingApprovalResponseMock };
export { evaluationFinishedProductionResponse as EvaluationFinishedProductionResponse };
export { evaluationFinishedStagingResponse as EvaluationFinishedStagingResponse };
export { latestFinishedDeployments as LatestFinishedDeployments };
export { latestFinishedEvaluations as LatestFinishedEvaluations };
export { openApprovalsResponse as OpenApprovalsResponse };
export { defaultDeploymentFinishedTrace as DefaultDeploymentFinishedTrace };
export { defaultDeploymentData as DefaultDeploymentData };
export { defaultEvaluationFinishedTrace as DefaultEvaluationFinishedTrace };
export { intersectDeploymentTriggeredResponse as IntersectDeploymentTriggeredResponse };
export { intersectDeploymentStartedResponse as IntersectDeploymentStartedResponse };
export { intersectDeploymentFinishedResponse as IntersectDeploymentFinishedResponse };
