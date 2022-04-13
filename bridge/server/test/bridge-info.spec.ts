import request from 'supertest';
import { setupServer } from '../.jest/setupServer';
import { Express } from 'express';

describe('Test /bridgeInfo', () => {
  let app: Express;

  beforeAll(async () => {
    app = await setupServer();
  });

  it('should return bridgeInfo', async () => {
    const response = await request(app).get('/api/bridgeInfo');
    expect(response.body).toEqual({
      bridgeVersion: 'develop',
      apiUrl: global.baseUrl,
      apiToken: 'apiToken',
      cliDownloadLink: 'https://github.com/keptn/keptn/releases',
      enableVersionCheckFeature: true,
      showApiToken: true,
      featureFlags: {
        RESOURCE_SERVICE_ENABLED: false,
      },
      authType: 'NONE',
    });
    expect(response.statusCode).toBe(200);
  });
});
