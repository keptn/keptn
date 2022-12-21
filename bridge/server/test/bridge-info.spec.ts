import request from 'supertest';
import { Express } from 'express';
import MockAdapter from 'axios-mock-adapter';
import { IMetadata } from '../../shared/interfaces/metadata';
import { BridgeInfo } from '../../shared/interfaces/bridge-info';
import { KeptnVersions } from '../../shared/interfaces/keptn-versions';
import { KeptnInfoResult } from '../../shared/interfaces/keptn-info-result';
// eslint-disable-next-line import/no-extraneous-dependencies
import { jest } from '@jest/globals';

const apiToken = 'abcdefg';
const getOAuthSecretsSpy = jest.fn();
const getBasicSecretsSpy = jest.fn();
const getMongoDbSecretsSpy = jest.fn();
getOAuthSecretsSpy.mockReturnValue({
  sessionSecret: 'session_secret',
  databaseEncryptSecret: 'database_secret_'.repeat(2),
  clientSecret: '',
});
getBasicSecretsSpy.mockReturnValue({
  apiToken: apiToken,
  user: '',
  password: '',
});
getMongoDbSecretsSpy.mockReturnValue({
  user: 'user',
  password: 'pwd',
});

jest.unstable_mockModule('../user/secrets', () => {
  return {
    getOAuthSecrets: getOAuthSecretsSpy,
    getBasicSecrets: getBasicSecretsSpy,
    getMongoDbSecrets: getMongoDbSecretsSpy,
    getOAuthMongoExternalConnectionString: (): string => '',
    getMongodbFolder: (): string => '',
    mongodbPasswordFileName: 'pwdName',
    mongodbUserFileName: 'userName',
  };
});

const { getConfiguration } = await import('../utils/configuration');
const { getBaseOptions, setupServer } = await import('../.jest/setupServer');

describe('Test /bridgeInfo', () => {
  let app: Express;

  const apiUrl = 'http://localhost/api/';
  const provMsg = '  message   ';
  const authMsg = 'a string';
  const version = 'testVersion';
  let axiosMock: MockAdapter;

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
    axiosMock = new MockAdapter(global.axiosInstance);
  });

  afterEach(() => {
    axiosMock.reset();
  });

  it('should fetch versions if enabled', async () => {
    mockMetadata();
    const getVersionCallCount = mockVersions();
    const response = await request(app).get('/api/bridgeInfo?isVersionCheckEnabled=true');
    expect(response.body).toEqual<BridgeInfo>({
      info: getInfoResponse(),
      metadata: getMetadataResponse(),
      versions: getVersionsResponse(),
    });
    expect(response.statusCode).toBe(200);
    expect(getVersionCallCount()).toBe(1);
  });

  it('should not fetch versions if disabled', async () => {
    mockMetadata();
    const getVersionCallCount = mockVersions();
    const response = await request(app).get('/api/bridgeInfo');
    expect(response.body).toEqual<BridgeInfo>({
      info: getInfoResponse(),
      metadata: getMetadataResponse(),
    });
    expect(response.statusCode).toBe(200);
    expect(getVersionCallCount()).toBe(0);
  });

  it('should fail if metadata call fails', async () => {
    mockMetadata(401);
    const response = await request(app).get('/api/bridgeInfo');
    expect(response.statusCode).toBe(401);
  });

  it('should not fail if version call fails', async () => {
    mockMetadata();
    const getVersionCallCount = mockVersions(401);
    const response = await request(app).get('/api/bridgeInfo?isVersionCheckEnabled=true');
    expect(response.body).toEqual<BridgeInfo>({
      info: getInfoResponse(),
      metadata: getMetadataResponse(),
      versions: {
        cli: {
          stable: [],
          prerelease: [],
        },
        bridge: {
          stable: [],
          prerelease: [],
        },
        keptn: {
          stable: [],
        },
      },
    });
    expect(response.statusCode).toBe(200);
    expect(getVersionCallCount()).toBe(1);
  });

  function mockMetadata(statusCode = 200): void {
    axiosMock.onGet(`${global.baseUrl}/v1/metadata`).reply<IMetadata>(statusCode, getMetadataResponse());
  }

  function mockVersions(statusCode = 200): () => number {
    const url = 'https://get.keptn.sh/version.json';
    axiosMock.onGet(url).reply<KeptnVersions>(statusCode, getVersionsResponse());
    return () => requestCount('get', url);
  }

  function requestCount(method: 'get' | 'post' | 'put', url: string): number {
    return axiosMock.history[method].filter((hr) => hr.url === url).length;
  }

  function getInfoResponse(): KeptnInfoResult {
    return {
      bridgeVersion: version,
      apiUrl: apiUrl,
      apiToken: apiToken,
      cliDownloadLink: 'https://github.com/keptn/keptn/releases',
      enableVersionCheckFeature: true,
      showApiToken: true,
      featureFlags: {},
      authType: 'NONE',
      automaticProvisioningMsg: provMsg.trim(),
      authMsg: authMsg,
      keptnInstallationType: 'QUALITY_GATES,CONTINUOUS_OPERATIONS,CONTINUOUS_DELIVERY',
      projectsPageSize: 50,
      servicesPageSize: 50,
    };
  }

  function getMetadataResponse(): IMetadata {
    return {
      bridgeversion: 'develop',
      keptnlabel: 'keptn',
      keptnversion: 'develop',
      namespace: '',
      shipyardversion: 'develop',
    };
  }

  function getVersionsResponse(): KeptnVersions {
    return {
      cli: {
        stable: ['0.0.1'],
        prerelease: ['0.0.2'],
      },
      bridge: {
        stable: ['0.0.2'],
        prerelease: ['0.0.3'],
      },
      keptn: {
        stable: [
          {
            version: '0.0.3',
            upgradableVersions: ['0.0.1', '0.0.2'],
          },
        ],
      },
    };
  }
});
