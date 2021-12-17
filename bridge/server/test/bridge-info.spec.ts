import request from 'supertest';
import { setupServer } from '../.jest/setupServer';

describe('Test /bridgeInfo', () => {
  beforeAll(async () => {
    await setupServer();
  });

  it('should return bridgeInfo', async () => {
    const response = await request(global.app).get('/api/bridgeInfo');
    expect(response.body).toEqual({
      bridgeVersion: 'develop',
      apiUrl: global.baseUrl,
      apiToken: 'apiToken',
      cliDownloadLink: 'https://github.com/keptn/keptn/releases',
      enableVersionCheckFeature: true,
      showApiToken: true,
      featureFlags: {},
      authType: 'NONE',
    });
    expect(response.statusCode).toBe(200);
  });
});
