import { UniformSubscription } from '../../shared/interfaces/uniform-subscription';
import { KeptnService } from '../../shared/models/keptn-service';


export class UniformRegistration {
  id!: string;
  metadata!: {
    deplyomentname: string,
    distributorversion: string,
    hostname: string,
    integrationversion: string,
    kubernetesmetadata: {
      deploymentname: string,
      namespace: string,
      podname: string
    },
    location: string,
    status: string,
    lastseen: string
  };
  unreadEventsCount?: number;
  name!: string;
  subscriptions: UniformSubscription[] = [];

  public static fromJSON(data: unknown): UniformRegistration {
    return Object.assign(new this(), data);
  }

  public get isWebhookService(): boolean {
    return this.name === KeptnService.WEBHOOK_SERVICE;
  }
}

