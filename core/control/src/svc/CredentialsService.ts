import { KeptnGithubCredentials } from '../lib/types/KeptnGithubCredentials';
import { K8sClientFactory } from '../lib/k8s/K8sClientFactory';
import * as K8sApi from 'kubernetes-client';

import { Logger } from '../lib/log/Logger';

import { base64encode, base64decode } from 'nodejs-base64';

export class CredentialsService {

  private static instance: CredentialsService;

  private k8sClient: K8sApi.ApiRoot;
  private constructor() {
    this.k8sClient = new K8sClientFactory().createK8sClient();
  }

  static getInstance() {
    if (CredentialsService.instance === undefined) {
      CredentialsService.instance = new CredentialsService();
    }
    return CredentialsService.instance;
  }

  async getKeptnApiToken(): Promise<string> {
    let token;
    try {
      const secret = await this.k8sClient.api.v1
        .namespaces('keptn').secrets
        .get({ name: 'keptn-api-token', pretty: true, exact: true, export: true });

      if (secret.body.items && secret.body.items.length > 0) {
        const apiToken = secret.body.items.find(item => item.metadata.name === 'keptn-api-token');
        if (apiToken && apiToken.data !== undefined) {
          token = base64decode(apiToken.data['keptn-api-token']);
        } else {
          console.log('[keptn] The secret does not contain the proper information.');
        }
      }
    } catch (e) {
      Logger.error('', `Error retrieving API Token: ${e}`);
      token = '';
    }

    return token;
  }

}
