import request from 'supertest';
import MockAdapter from 'axios-mock-adapter';
import { ProjectQualityGatesResponse, ProjectResponse } from '../fixtures/project-response.mock';
import { OpenRemediationsResponse } from '../fixtures/open-remediations-response.mock';
import {
  ServiceStateQualityGatesOnlyResponse,
  ServiceStateResponse,
} from '../../shared/fixtures/service-state-response.mock';

let axiosMock: MockAdapter;

describe('Test project/:projectName/serviceStates', () => {
  beforeAll(() => {
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

    const response = await request(global.app).get(`/api/project/${projectName}/serviceStates`);
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

    const response = await request(global.app).get(`/api/project/${projectName}/serviceStates`);
    expect(response.body).toEqual(ServiceStateQualityGatesOnlyResponse);
    expect(response.statusCode).toBe(200);
  });
});
