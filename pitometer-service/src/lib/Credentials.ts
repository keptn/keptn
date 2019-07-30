import { DynatraceCredentialsModel } from './DynatraceCredentialsModel';
import { K8sClientFactory } from './K8sClientFactory';

import * as K8sApi from 'kubernetes-client';

import { base64decode } from 'nodejs-base64';
import { Logger } from './Logger';

export class Credentials {

  private static instance: Credentials;

  private k8sClient: K8sApi.ApiRoot;
  private constructor() {
    this.k8sClient = new K8sClientFactory().createK8sClient();
  }

  static getInstance() {
    if (Credentials.instance === undefined) {
      Credentials.instance = new Credentials();
    }
    return Credentials.instance;
  }

  async getDynatraceCredentials(): Promise<DynatraceCredentialsModel> {
    const dynatraceCredentials: DynatraceCredentialsModel = {} as DynatraceCredentialsModel;

    const secret = await this.k8sClient.api.v1
      .namespaces('keptn').secrets
      .get({ name: 'dynatrace', pretty: true, exact: true, export: true });

    if (secret.body.items && secret.body.items.length > 0) {
      const dtItem = secret.body.items.find(item => item.metadata.name === 'dynatrace');
      if (dtItem && dtItem.data !== undefined) {
        dynatraceCredentials.tenant = base64decode(dtItem.data.DT_TENANT);
        dynatraceCredentials.apiToken = base64decode(dtItem.data.DT_API_TOKEN);
      }
    }

    return dynatraceCredentials;
  }
}
