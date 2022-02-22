import request from 'supertest';
import MockAdapter from 'axios-mock-adapter';
import { setupServer } from '../.jest/setupServer';
import { Express } from 'express';
import { StagesResponse } from '../fixtures/stages.mock';

let axiosMock: MockAdapter;

describe('Test /project/:projectName/stages', () => {
  let app: Express;

  beforeAll(async () => {
    app = await setupServer();
    axiosMock = new MockAdapter(global.axiosInstance);
  });

  afterEach(() => {
    axiosMock.reset();
  });

  it('should retrieve stage names', async () => {
    const projectName = 'sockshop';
    axiosMock.onGet(`${global.baseUrl}/controlPlane/v1/project/${projectName}/stage`).reply(200, StagesResponse);
    const response = await request(app).get(`/api/project/${projectName}/stages`);
    expect(response.body).toEqual(['dev', 'production', 'staging']);
    expect(response.statusCode).toBe(200);
  });
});
