import request from 'supertest';
import MockAdapter from 'axios-mock-adapter';
import { ProjectQualityGatesResponse, ProjectResponse } from '../../shared/fixtures/project-response.mock';
import { OpenRemediationsResponse } from '../../shared/fixtures/open-remediations-response.mock';
import {
  ServiceStateQualityGatesOnlyResponse,
  ServiceStateResponse,
} from '../../shared/fixtures/service-state-response.mock';
import { setupServer } from '../.jest/setupServer';
import { Express } from 'express';

let axiosMock: MockAdapter;

describe('Test project/:projectName/serviceStates', () => {
  let app: Express;

  beforeAll(async () => {
    app = await setupServer();
    axiosMock = new MockAdapter(global.axiosInstance);
  });

  afterEach(() => {
    axiosMock.reset();
  });

  it('should retrieve service states', async () => {
    const projectName = 'sockshop';
    axiosMock.onGet(global.baseUrl + `/controlPlane/v1/project/${projectName}`).reply(200, ProjectResponse);

    axiosMock
      .onGet(`${global.baseUrl}/controlPlane/v1/sequence/${projectName}`, {
        params: {
          pageSize: '100',
          name: 'remediation',
          state: 'started',
        },
      })
      .reply(200, OpenRemediationsResponse);

    const response = await request(app).get(`/api/project/${projectName}/serviceStates`);
    expect(response.body).toEqual(ServiceStateResponse);
    expect(response.statusCode).toBe(200);
  });

  it('should retrieve service states for quality gates only', async () => {
    const projectName = 'sockshop';
    axiosMock.onGet(global.baseUrl + `/controlPlane/v1/project/${projectName}`).reply(200, ProjectQualityGatesResponse);

    axiosMock
      .onGet(`${global.baseUrl}/controlPlane/v1/sequence/${projectName}`, {
        params: {
          pageSize: '100',
          name: 'remediation',
          state: 'started',
        },
      })
      .reply(200, { states: [] });

    const response = await request(app).get(`/api/project/${projectName}/serviceStates`);
    expect(response.body).toEqual(ServiceStateQualityGatesOnlyResponse);
    expect(response.statusCode).toBe(200);
  });
});
