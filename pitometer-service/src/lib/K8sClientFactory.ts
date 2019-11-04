import * as Api from 'kubernetes-client';

export class K8sClientFactory {

  constructor() { }

  createK8sClient(): Api.ApiRoot {
    // tslint:disable-next-line: variable-name
    const Client = Api.Client1_13;
    const config = Api.config;
    let k8sClient;

    if (process.env.NODE_ENV === 'production') {
      k8sClient = new Client({ config: config.getInCluster() });
    } else {
      k8sClient = new Client({ config: config.fromKubeconfig() });
    }

    return k8sClient;
  }
}
