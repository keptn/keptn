import request from 'supertest';
import { Express } from 'express';
// import {expect, jest, test, } from '@jest/globals';
// eslint-disable-next-line import/no-extraneous-dependencies
//
// jest.mock(axios);
// import { ShipyardResponse } from '../fixtures/shipyard-response';
let app: Express;

describe('Test the root path', () => {
  beforeAll(async () => {
    app = global.app;
  });
  test('It should return bridgeInfo', async () => {
    const response = await request(app).get('/api/bridgeInfo');
    expect(response.body).toEqual({
      apiUrl: global.baseUrl,
      apiToken: 'apiToken',
      cliDownloadLink: 'https://github.com/keptn/keptn/releases',
      enableVersionCheckFeature: true,
      showApiToken: true,
      authType: 'NONE',
    });
    expect(response.statusCode).toBe(200);
  });
  it('should retrieve task names', async () => {});
});
