import { inject, injectable } from 'inversify';
import { K8sClientFactory } from '../k8s/K8sClientFactory';
import { KeptnOrgToRepoConfigMapModel } from './KeptnOrgToRepoConfigMapModel';

@injectable()
export class OrgToRepoMapper {

  private k8sClient: KubernetesClient.ApiRoot;

  constructor() {
    const clientFactory = new K8sClientFactory();
    this.k8sClient = clientFactory.createK8sClient();
  }

  public async getRepoForOrg(org: string): Promise<string> {
    let result = '';

    const configMap = await this.k8sClient.api.v1.namespace('keptn').configmap('keptn-orgs').get();

    if (configMap === undefined || configMap.body === undefined) {
      return result;
    }

    const cm = configMap.body as KeptnOrgToRepoConfigMapModel;

    const orgsToRepos: any = JSON.parse(cm.data.orgsToRepos);

    if (orgsToRepos[org] !== undefined) {
      result = orgsToRepos[org];
    }

    return result;
  }
}
