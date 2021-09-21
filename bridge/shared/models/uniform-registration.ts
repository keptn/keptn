import { UniformRegistrationResult } from '../interfaces/uniform-registration-result';
import { UniformSubscription } from '../interfaces/uniform-subscription';
import { KeptnService } from './keptn-service';

export class UniformRegistration implements UniformRegistrationResult {
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

  public get isWebhookService(): boolean {
    return this.name === KeptnService.WEBHOOK_SERVICE;
  }
}

