import request from 'supertest';
import MockAdapter from 'axios-mock-adapter';
import { ProjectResponse } from '../fixtures/project-response.mock';
import { RemediationConfigResponse } from '../fixtures/remediation-config-response.mock';
import { DeploymentTracesResponseMock } from '../fixtures/deployment-traces-response.mock';
import { KeptnService } from '../../shared/models/keptn-service';
import { EventTypes } from '../../shared/interfaces/event-types';
import { ServiceDeploymentMock } from '../fixtures/service-deployment.mock';
import { OpenRemediationsResponse, RemediationTraceResponse } from '../fixtures/open-remediations-response.mock';
import { SequenceDeliveryResponseMock } from '../fixtures/sequence-delivery-response.mock';

let axiosMock: MockAdapter;

describe('Test /project/:projectName/deployment/:keptnContext', () => {
  beforeAll(() => {
    axiosMock = new MockAdapter(global.axiosInstance);
  });

  afterEach(() => {
    axiosMock.reset();
  });

  it('should retrieve deployment of service', async () => {
    const projectName = 'sockshop';
    const keptnContext = '2c0e568b-8bd3-4726-a188-e528423813ed';
    axiosMock.onGet(global.baseUrl + `/controlPlane/v1/project/${projectName}`).reply(200, ProjectResponse);
    axiosMock
      .onGet(`${global.baseUrl}/mongodb-datastore/event/type/${EventTypes.EVALUATION_FINISHED}`, {
        params: {
          filter: `data.project:${projectName} AND data.service:carts AND data.stage:production AND source:${KeptnService.LIGHTHOUSE_SERVICE}`,
          excludeInvalidated: 'true',
          limit: '1',
        },
      })
      .reply(200, {
        events: [
          {
            data: {
              evaluation: {
                gitCommit: '',
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
      });

    axiosMock.onGet(`${global.baseUrl}/mongodb-datastore/event/type/${EventTypes.EVALUATION_FINISHED}`).reply(200, {
      pageSize: 1,
      events: [],
      totalCount: 0,
    });

    axiosMock
      .onGet(
        `${global.baseUrl}/configuration-service/v1/project/${projectName}/stage/production/service/carts/resource/remediation.yaml`
      )
      .reply(200, RemediationConfigResponse);
    axiosMock
      .onGet(`${global.baseUrl}/mongodb-datastore/event`, {
        params: {
          keptnContext,
          project: projectName,
          pageSize: '100',
        },
      })
      .reply(200, DeploymentTracesResponseMock);

    axiosMock
      .onGet(`${global.baseUrl}/mongodb-datastore/event`, {
        params: {
          project: 'sockshop',
          service: 'carts',
          stage: 'production',
          keptnContext: '35383737-3630-4639-b037-353138323631',
          pageSize: '1',
        },
      })
      .reply(200, RemediationTraceResponse);

    axiosMock
      .onGet(`${global.baseUrl}/controlPlane/v1/sequence/${projectName}`, {
        params: {
          pageSize: '1',
          keptnContext,
        },
      })
      .reply(200, SequenceDeliveryResponseMock);

    axiosMock
      .onGet(`${global.baseUrl}/controlPlane/v1/sequence/${projectName}`, {
        params: {
          pageSize: '100',
          name: 'remediation',
          state: 'started',
        },
      })
      .reply(200, OpenRemediationsResponse);

    const response = await request(global.app).get(`/api/project/${projectName}/deployment/${keptnContext}`);
    expect(response.body).toEqual(ServiceDeploymentMock);
    expect(response.statusCode).toBe(200);
  });
});
