import { K8sClientFactory } from '../k8s/K8sClientFactory';

export class ChannelReconciler {

  private k8sClient: KubernetesClient.ApiRoot;

  constructor() {
    const clientFactory = new K8sClientFactory();
    this.k8sClient = clientFactory.createK8sClient();
  }

  public async resolveChannel(channelName: string): Promise<string> {
    const services = await this.k8sClient.api.v1.namespace('keptn').services.get();
    if (services.body === undefined) {
      return '';
    }
    if (services.body.items === undefined || services.body.items.length === 0) {
      return '';
    }
    const channelData = services.body.items.find((svc) => {
      if (svc.metadata === undefined) {
        return false;
      }
      if (svc.metadata.labels === undefined) {
        return false;
      }
      if (svc.metadata.labels.channel === undefined) {
        return false;
      }
      if (svc.metadata.labels.channel === channelName) {
        return true;
      }
      return false;
    });
    if (channelData === undefined) {
      return '';
    }
    return `${channelData.metadata.name}.keptn.svc.cluster.local`;
  }
}
