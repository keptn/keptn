import { UniformSubscription } from './uniformSubscription';

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
  name!: string;
  subscriptions: UniformSubscription[] = [];

  public static fromJSON(data: unknown): UniformRegistration {
    const uniformRegistration = Object.assign(new this(), data);
    uniformRegistration.subscriptions = uniformRegistration.subscriptions?.map(subscription => UniformSubscription.fromJSON(subscription)) ?? [];
    return uniformRegistration;
  }

  public getSubscriptions(projectName: string): UniformSubscription[] {
    return this.subscriptions.filter(subscription => subscription.project === projectName || subscription.project === undefined);
  }

  public hasSubscriptions(projectName: string): boolean {
    return this.subscriptions.some(subscription => subscription.project === projectName || subscription.project === undefined);
  }

  public formatSubscriptions(projectName: string): string | undefined {
    const subscriptions = this.subscriptions.reduce((accSubscriptions: string[], subscription: UniformSubscription) => {
      if (subscription.project === projectName) {
        accSubscriptions.push(subscription.topics[0]);
      }
      return accSubscriptions;
    }, []);
    return subscriptions.length !== 0 ? subscriptions.join('<br/>') : undefined;
  }
}

