import request from 'supertest';
import { StagesResponse } from '../fixtures/stages.mock';
import MockAdapter from 'axios-mock-adapter';
import { setupServer } from '../.jest/setupServer';
import { Express } from 'express';

let axiosMock: MockAdapter;

describe('Test /project/:projectName/sequences/filter', () => {
  let app: Express;

  beforeAll(async () => {
    app = await setupServer();
    axiosMock = new MockAdapter(global.axiosInstance);
  });

  afterEach(() => {
    axiosMock.reset();
  });

  it('should return sequence filter', async () => {
    const projectName = 'sockshop';
    axiosMock.onGet(`${global.baseUrl}/controlPlane/v1/project/${projectName}/stage`).reply(200, StagesResponse);
    const response = await request(app).get(`/api/project/${projectName}/sequences/filter`);
    expect(response.body).toEqual({
      services: ['carts', 'carts-db'],
      stages: ['dev', 'production', 'staging'],
    });
    expect(response.statusCode).toBe(200);
  });
});
