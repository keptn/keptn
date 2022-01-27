import MockAdapter from 'axios-mock-adapter';
import {
  OpenRemediationsResponse,
  RemediationTraceResponse,
} from '../../shared/fixtures/open-remediations-response.mock';
import { RemediationConfigResponse } from '../fixtures/remediation-config-response.mock';
import { init } from '../app';
import { Express } from 'express';

export class TestUtils {
  public static mockOpenRemediations(axiosMock: MockAdapter, projectName: string): void {
    axiosMock
      .onGet(`${global.baseUrl}/controlPlane/v1/sequence/${projectName}`, {
        params: {
          pageSize: '100',
          name: 'remediation',
          state: 'started',
        },
      })
      .reply(200, OpenRemediationsResponse);

    axiosMock
      .onGet(
        `${global.baseUrl}/configuration-service/v1/project/${projectName}/stage/production/service/carts/resource/remediation.yaml`
      )
      .reply(200, RemediationConfigResponse);

    axiosMock
      .onGet(`${global.baseUrl}/mongodb-datastore/event`, {
        params: {
          project: 'sockshop',
          service: 'carts',
          stage: 'production',
          keptnContext: '35383737-3630-4639-b037-353138323631',
          pageSize: '50',
          type: `sh.keptn.event.production.remediation.triggered`,
        },
      })
      .reply(200, RemediationTraceResponse);
  }

  public static async setupOAuthTest(): Promise<Express> {
    process.env.OAUTH_ENABLED = 'true';
    process.env.OAUTH_CLIENT_ID = 'myClientID';
    process.env.OAUTH_BASE_URL = 'http://localhost';
    process.env.OAUTH_DISCOVERY = 'http://localhost/.well-known/openid-configuration';
    return init();
  }
}
