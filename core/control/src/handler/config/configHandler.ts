import { KeptnConfig } from '../../lib/types/KeptnConfig';
import { K8sClientFactory } from '../../lib/k8s/K8sClientFactory';
import * as K8sApi from 'kubernetes-client';

import { KeptnConfigSecretFactory } from '../../lib/types/KeptnConfigSecretFactory';
import { KeptnConfigSecret } from '../../lib/types/KeptnConfigSecret';

export class ConfigHandler {

  private k8sClient: K8sApi.ApiRoot;
  constructor() {

  }

  async init() {
    this.k8sClient = new K8sClientFactory().createK8sClient();
  }

  async updateKeptnConfig(keptnConfig: KeptnConfig) {
    const secret = new KeptnConfigSecretFactory().createKeptnConfigSecret(keptnConfig);

    const created = await this.updateGithubCredentials(secret);
    console.log(created);
  }

  private async updateGithubCredentials(secret: KeptnConfigSecret) {
    try {
      const deleteResult = await this.k8sClient.api.v1
        .namespaces('keptn').secrets('github-credentials').delete();
      console.log(deleteResult);
    }
    catch (e) { }
    const created = await this.k8sClient.api.v1.namespaces('keptn').secrets.post({
      body: secret,
    });

    return created;
  }
}
