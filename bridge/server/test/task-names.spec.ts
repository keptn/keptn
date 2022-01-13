import request from 'supertest';
import MockAdapter from 'axios-mock-adapter';
import { ShipyardResponse } from '../fixtures/shipyard-response.mock';
import { setupServer } from '../.jest/setupServer';
import { Express } from 'express';

let axiosMock: MockAdapter;

describe('Test /project/:projectName/tasks', () => {
  let app: Express;

  beforeAll(async () => {
    app = await setupServer();
    axiosMock = new MockAdapter(global.axiosInstance);
  });

  afterEach(() => {
    axiosMock.reset();
  });

  it('should retrieve task names', async () => {
    const projectName = 'sockshop';
    axiosMock
      .onGet(`${global.baseUrl}/configuration-service/v1/project/${projectName}/resource/shipyard.yaml`)
      .reply(200, ShipyardResponse);
    const response = await request(app).get(`/api/project/${projectName}/tasks`);
    expect(response.body).toEqual(['evaluation', 'deployment', 'test', 'release', 'rollback', 'get-action', 'action']);
    expect(response.statusCode).toBe(200);
  });
});
