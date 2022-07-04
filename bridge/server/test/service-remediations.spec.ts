import MockAdapter from 'axios-mock-adapter';
import request from 'supertest';
import { TestUtils } from '../.jest/test.utils';
import { setupServer } from '../.jest/setupServer';
import { Express } from 'express';
import { ServiceRemediationResponse } from '../../shared/fixtures/open-remediations-response.mock';

let axiosMock: MockAdapter;

describe('Test /project/:projectName/service/:serviceName/openRemediations', () => {
  let app: Express;

  beforeAll(async () => {
    app = await setupServer();
    axiosMock = new MockAdapter(global.axiosInstance);
  });

  afterEach(() => {
    axiosMock.reset();
  });

  it('should retrieve open remediations', async () => {
    const projectName = 'sockshop';
    const serviceName = 'carts';
    TestUtils.mockOpenRemediations(axiosMock, projectName);
    const response = await request(app).get(`/api/project/${projectName}/service/${serviceName}/openRemediations`);
    expect(response.body).toEqual(ServiceRemediationResponse);
    expect(response.statusCode).toBe(200);
  });
});
