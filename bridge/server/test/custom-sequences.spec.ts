import request from 'supertest';
import MockAdapter from 'axios-mock-adapter';
import { ShipyardResponse } from '../fixtures/shipyard-response.mock';
import { setupServer } from '../.jest/setupServer';
import { Express } from 'express';

let axiosMock: MockAdapter;

describe('Test /project/:projectName/customSequences', () => {
  let app: Express;

  beforeAll(async () => {
    app = await setupServer();
    axiosMock = new MockAdapter(global.axiosInstance);
  });

  afterEach(() => {
    axiosMock.reset();
  });

  it('should retrieve custom sequences', async () => {
    const projectName = 'sockshop';
    axiosMock
      .onGet(`${global.baseUrl}/configuration-service/v1/project/${projectName}/resource/shipyard.yaml`)
      .reply(200, ShipyardResponse);
    const response = await request(app).get(`/api/project/${projectName}/customSequences`);
    expect(response.body).toEqual(['delivery-direct', 'rollback', 'remediation']);
    expect(response.statusCode).toBe(200);
  });
});
