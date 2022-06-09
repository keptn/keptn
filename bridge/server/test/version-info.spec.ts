import Axios from 'axios';
import { Express } from 'express';
// eslint-disable-next-line import/no-extraneous-dependencies
import { jest } from '@jest/globals';
import request from 'supertest';
import MockAdapter from 'axios-mock-adapter';
import { KeptnVersions } from '../../shared/interfaces/keptn-versions';

const mockInstance = Axios.create();

jest.unstable_mockModule('../services/axios-instance', () => {
  return { axios: mockInstance };
});

const { init } = await import('../app');

describe('Test /version.json', () => {
  let app: Express;
  let axiosMock: MockAdapter;

  beforeAll(async () => {
    app = await init();
    axiosMock = new MockAdapter(mockInstance);
  });

  afterEach(() => {
    axiosMock.reset();
  });

  it('should return successful calls unmodified', async () => {
    const anyPayload = { prop: Math.random() };
    axiosMock.onGet(`https://get.keptn.sh/version.json`).reply(200, anyPayload);
    const response = await request(app).get(`/api/version.json`);
    expect(response.body).toEqual(anyPayload);
    expect(response.statusCode).toBe(200);
  });

  it('should return a default version if the call to the Keptn page fails', async () => {
    axiosMock.onGet(`https://get.keptn.sh/version.json`).reply(500);
    const response = await request(app).get(`/api/version.json`);
    const keptnVersions = response.body as KeptnVersions;
    expect(keptnVersions).toBeDefined();
    expect(keptnVersions).toEqual({
      cli: { stable: [], prerelease: [] },
      bridge: { stable: [], prerelease: [] },
      keptn: { stable: [] },
    });
  });
});
