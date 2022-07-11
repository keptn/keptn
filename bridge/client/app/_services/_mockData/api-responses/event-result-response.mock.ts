import { TestUtils } from '../../../_utils/test.utils';

const eventData = [
  {
    data: {
      deployment: {
        deploymentNames: ['canary'],
        deploymentURIsLocal: ['http://carts.sockshop-production:80'],
        deploymentURIsPublic: ['http://carts.sockshop-production.35.188.183.151.nip.io:80'],
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
    },
    id: '0ef7698a-284c-4ce6-b349-82abf8f41f60',
    source: 'helm-service',
    specversion: '1.0',
    time: '2021-11-11T13:21:21.355Z',
    type: 'sh.keptn.event.deployment.finished',
    shkeptncontext: '76f0b0af-0290-458e-82da-56bec6ec5868',
    shkeptnspecversion: '0.2.3',
    triggeredid: '2884b495-b467-4ab2-b8d0-873c4cfc9781',
  },
];

const data = TestUtils.mapTraces(eventData);
export { data as EventResultResponseMock };
