import request from 'supertest';
import MockAdapter from 'axios-mock-adapter';
import { ProjectResponse } from '../../shared/fixtures/project-response.mock';
import {
  DeploymentTracesResponseMock,
  DeploymentTracesWithPendingApprovalResponseMock,
  EvaluationFinishedProductionResponse,
  EvaluationFinishedStagingResponse,
} from '../../shared/fixtures/traces-response.mock';
import { KeptnService } from '../../shared/models/keptn-service';
import { EventTypes } from '../../shared/interfaces/event-types';
import {
  ServiceDeploymentMock,
  ServiceDeploymentWithApprovalMock,
  ServiceDeploymentWithFromTimeMock,
} from '../../shared/fixtures/service-deployment-response.mock';
import {
  SequenceDeliveryResponseMock,
  SequenceDeliveryTillStagingResponseMock,
} from '../fixtures/sequence-response.mock';
import { TestUtils } from '../.jest/test.utils';
import { setupServer } from '../.jest/setupServer';
import { Express } from 'express';

let axiosMock: MockAdapter;

describe('Test /project/:projectName/deployment/:keptnContext', () => {
  let app: Express;

  beforeAll(async () => {
    app = await setupServer();
    axiosMock = new MockAdapter(global.axiosInstance);
  });

  afterEach(() => {
    axiosMock.reset();
  });

  it('should retrieve deployment of service', async () => {
    const projectName = 'sockshop';
    const keptnContext = '2c0e568b-8bd3-4726-a188-e528423813ed';
    axiosMock
      .onGet(`${global.baseUrl}/mongodb-datastore/event`, {
        params: {
          keptnContext,
          project: projectName,
          pageSize: '100',
        },
      })
      .reply(200, DeploymentTracesResponseMock);
    init(ProjectResponse, SequenceDeliveryResponseMock, projectName, keptnContext);

    const response = await request(app).get(`/api/project/${projectName}/deployment/${keptnContext}`);
    expect(response.body).toEqual(ServiceDeploymentMock);
    expect(response.statusCode).toBe(200);
  });

  it('should retrieve deployment of service with fromTime', async () => {
    const projectName = 'sockshop';
    const keptnContext = '2c0e568b-8bd3-4726-a188-e528423813ed';
    axiosMock
      .onGet(`${global.baseUrl}/mongodb-datastore/event`, {
        params: {
          keptnContext,
          project: projectName,
          pageSize: '100',
        },
      })
      .reply(200, DeploymentTracesResponseMock);
    init(ProjectResponse, SequenceDeliveryResponseMock, projectName, keptnContext);

    const response = await request(app).get(
      `/api/project/${projectName}/deployment/${keptnContext}?fromTime=2021-10-13T10:59:45.104Z`
    );
    expect(response.body).toEqual(ServiceDeploymentWithFromTimeMock);
    expect(response.statusCode).toBe(200);
  });

  it('should return service deployment with pending approval', async () => {
    // service.latestDeployment === keptnContext in getTraces()
    const projectName = 'sockshop';
    const keptnContext = '2c0e568b-8bd3-4726-a188-e528423813ed';
    axiosMock
      .onGet(`${global.baseUrl}/mongodb-datastore/event`, {
        params: {
          keptnContext,
          project: projectName,
          pageSize: '100',
        },
      })
      .reply(200, DeploymentTracesWithPendingApprovalResponseMock);
    axiosMock
      .onGet(`${global.baseUrl}/mongodb-datastore/event/type/${EventTypes.EVALUATION_FINISHED}`, {
        params: {
          filter: `data.project:${projectName} AND data.service:carts AND data.stage:staging AND source:${KeptnService.LIGHTHOUSE_SERVICE}`,
          excludeInvalidated: 'true',
          limit: '1',
        },
      })
      .reply(200, EvaluationFinishedStagingResponse);
    init(ProjectResponse, SequenceDeliveryTillStagingResponseMock, projectName, keptnContext);

    const response = await request(app).get(`/api/project/${projectName}/deployment/${keptnContext}`);
    expect(response.body).toEqual(ServiceDeploymentWithApprovalMock);
    expect(response.statusCode).toBe(200);
  });

  function init(projectMock: unknown, sequenceMock: unknown, projectName: string, keptnContext: string): void {
    axiosMock.onGet(global.baseUrl + `/controlPlane/v1/project/${projectName}`).reply(200, projectMock);
    axiosMock
      .onGet(`${global.baseUrl}/mongodb-datastore/event/type/${EventTypes.EVALUATION_FINISHED}`, {
        params: {
          filter: `data.project:${projectName} AND data.service:carts AND data.stage:production AND source:${KeptnService.LIGHTHOUSE_SERVICE}`,
          excludeInvalidated: 'true',
          limit: '1',
        },
      })
      .reply(200, EvaluationFinishedProductionResponse);

    axiosMock.onGet(`${global.baseUrl}/mongodb-datastore/event/type/${EventTypes.EVALUATION_FINISHED}`).reply(200, {
      pageSize: 1,
      events: [],
      totalCount: 0,
    });

    TestUtils.mockOpenRemediations(axiosMock, projectName);

    axiosMock
      .onGet(`${global.baseUrl}/controlPlane/v1/sequence/${projectName}`, {
        params: {
          pageSize: '1',
          keptnContext,
        },
      })
      .reply(200, sequenceMock);
  }
});
