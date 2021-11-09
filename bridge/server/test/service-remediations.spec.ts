import MockAdapter from 'axios-mock-adapter';
import request from 'supertest';
import { TestUtils } from '../.jest/test.utils';
import {
  ServiceRemediationWithConfigResponse,
  ServiceRemediationWithoutConfigResponse,
} from '../fixtures/open-remediations-response.mock';

let axiosMock: MockAdapter;

describe('Test /project/:projectName/service/:serviceName/openRemediations', () => {
  beforeAll(() => {
    axiosMock = new MockAdapter(global.axiosInstance);
  });

  afterEach(() => {
    axiosMock.reset();
  });

  it('should retrieve open remediations with config', async () => {
    const projectName = 'sockshop';
    const serviceName = 'carts';
    TestUtils.mockOpenRemediations(axiosMock, projectName);
    const response = await request(global.app).get(
      `/api/project/${projectName}/service/${serviceName}/openRemediations?config=true`
    );
    expect(response.body).toEqual(ServiceRemediationWithConfigResponse);
    expect(response.statusCode).toBe(200);
  });

  it('should retrieve open remediations without config', async () => {
    const projectName = 'sockshop';
    const serviceName = 'carts';
    TestUtils.mockOpenRemediations(axiosMock, projectName);
    const response = await request(global.app).get(
      `/api/project/${projectName}/service/${serviceName}/openRemediations`
    );
    expect(response.body).toEqual(ServiceRemediationWithoutConfigResponse);
    expect(response.statusCode).toBe(200);
  });
});
