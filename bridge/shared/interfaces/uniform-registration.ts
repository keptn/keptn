import { IUniformSubscription } from './uniform-subscription';
import { KeptnService } from '../models/keptn-service';

export function isWebhookService(ur: IUniformRegistration): boolean {
  return ur.name === KeptnService.WEBHOOK_SERVICE;
}

export interface IUniformRegistration {
  id: string;
  metadata: {
    deplyomentname: string;
    distributorversion: string;
    hostname: string;
    integrationversion: string;
    kubernetesmetadata: {
      deploymentname: string;
      namespace: string;
      podname: string;
    };
    location: string;
    status: string;
    lastseen: string;
  };
  unreadEventsCount?: number;
  name: string;
  subscriptions: IUniformSubscription[];
}
