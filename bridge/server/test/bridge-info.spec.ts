import request from 'supertest';
import { getBaseOptions, setupServer } from '../.jest/setupServer';
import { Express } from 'express';
import { getConfiguration } from '../utils/configuration';

describe('Test /bridgeInfo', () => {
  let app: Express;

  const apiUrl = 'http://localhost:8090/api';
  const apiToken = 'abcdefg';
  const provMsg = '  message   ';
  const authMsg = 'a string';
  const version = 'testVersion';

  beforeAll(async () => {
    const conf = getConfiguration({
      ...getBaseOptions(),
      api: {
        url: apiUrl,
        token: apiToken,
        showToken: true,
      },
      auth: {
        authMessage: authMsg,
      },
      features: {
        automaticProvisioningMessage: provMsg,
      },
      version: version,
    });
    app = await setupServer(conf);
  });

  it('should return bridgeInfo', async () => {
    const response = await request(app).get('/api/bridgeInfo');
    expect(response.body).toEqual({
      bridgeVersion: version,
      apiUrl: apiUrl,
      apiToken: apiToken,
      cliDownloadLink: 'https://github.com/keptn/keptn/releases',
      enableVersionCheckFeature: true,
      showApiToken: true,
      featureFlags: {
        RESOURCE_SERVICE_ENABLED: false,
        D3_HEATMAP_ENABLED: false,
      },
      authType: 'NONE',
      automaticProvisioningMsg: provMsg.trim(),
      authMsg: authMsg,
      keptnInstallationType: 'QUALITY_GATES,CONTINUOUS_OPERATIONS,CONTINUOUS_DELIVERY',
      projectsPageSize: 50,
      servicesPageSize: 50,
    });
    expect(response.statusCode).toBe(200);
  });
});
