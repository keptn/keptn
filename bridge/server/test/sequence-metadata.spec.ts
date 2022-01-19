import request from 'supertest';
import { StagesResponse } from '../fixtures/stages.mock';
import MockAdapter from 'axios-mock-adapter';
import { setupServer } from '../.jest/setupServer';
import { Express } from 'express';

let axiosMock: MockAdapter;

describe('Test /project/:projectName/sequences/metadata', () => {
  let app: Express;

  beforeAll(async () => {
    app = await setupServer();
    axiosMock = new MockAdapter(global.axiosInstance);
  });

  afterEach(() => {
    axiosMock.reset();
  });

  it('should return sequence metadata', async () => {
    const projectName = 'sockshop';
    axiosMock.onGet(`${global.baseUrl}/controlPlane/v1/project/${projectName}/stage`).reply(200, StagesResponse);
    const response = await request(app).get(`/api/project/${projectName}/sequences/metadata`);
    expect(response.body).toEqual({
      deployments: [
        {
          stage: {
            name: 'dev',
            services: [
              {
                name: 'carts',
                image: 'carts:0.12.1',
              },
              {
                name: 'carts-db',
                image: 'mongo:4.2.2',
              },
            ],
          },
        },
        {
          stage: {
            name: 'production',
            services: [
              {
                name: 'carts',
                image: 'carts:0.12.1',
              },
              {
                name: 'carts-db',
                image: 'mongo:4.2.2',
              },
            ],
          },
        },
        {
          stage: {
            name: 'staging',
            services: [
              {
                name: 'carts',
                image: 'carts:0.12.1',
              },
              {
                name: 'carts-db',
                image: 'mongo:4.2.2',
              },
            ],
          },
        },
      ],
      filter: {
        stages: ['dev', 'production', 'staging'],
        services: ['carts', 'carts-db'],
      },
    });
    expect(response.statusCode).toBe(200);
  });
});
