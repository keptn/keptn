import { UniformSubscription } from './uniform-subscription';

export interface UniformRegistrationResult {
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
  subscriptions: UniformSubscription[];
}
