import { UniformSubscription } from '../../shared/interfaces/uniform-subscription';


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
    status: string
  };
  unreadEventsCount?: number;
  name!: string;
  subscriptions: UniformSubscription[] = [];
}

