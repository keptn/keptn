import { init } from '../app';
import { Express } from 'express';
import request from 'supertest';
// eslint-disable-next-line import/no-extraneous-dependencies
// import { jest } from '@jest/globals';
//
// jest.mock(axios);
// import { ShipyardResponse } from '../fixtures/shipyard-response';
let app: Express;

describe('Test the root path', () => {
  beforeAll(async () => {
    app = await init();
    app.set('port', 80);
    app.listen(80, '0.0.0.0');
  });
  test('It should return bridgeInfo', async () => {
    const response = await request(app).get('/api/bridgeInfo');
    expect(response.body).toEqual({
      apiUrl: 'http://localhost',
      apiToken: 'apiToken',
      cliDownloadLink: 'https://github.com/keptn/keptn/releases',
      enableVersionCheckFeature: true,
      showApiToken: true,
      authType: 'NONE',
    });
    expect(response.statusCode).toBe(200);
  });
});
