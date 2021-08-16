import { UniformRegistration as ur } from '../../../server/interfaces/uniform-registration';
import { UniformSubscription } from './uniform-subscription';


export class UniformRegistration extends ur {
  public subscriptions: UniformSubscription[] = [];
  public unreadEventsCount!: number;

  public static fromJSON(data: unknown): UniformRegistration {
    const uniformRegistration = Object.assign(new this(), data);
    uniformRegistration.subscriptions = uniformRegistration.subscriptions?.map(subscription => UniformSubscription.fromJSON(subscription)) ?? [];
    return uniformRegistration;
  }

  public getSubscriptions(projectName: string): UniformSubscription[] {
    const subscriptions = this.subscriptions.filter(subscription => subscription.project === projectName || !subscription.project);
    subscriptions.sort((a, b) => {
      let status;
      if (!a.project) {
        status = -1;
      }
      else if (!b.project) {
        status = 1;
      }
      else {
        status = a.topic.localeCompare(b.topic);
      }
      return status;
    });
    return subscriptions;
  }

  public hasSubscriptions(projectName: string): boolean {
    return this.subscriptions.some(subscription => subscription.project === projectName || !subscription.project);
  }

  public formatSubscriptions(projectName: string): string | undefined {
    const subscriptions = this.subscriptions.reduce((accSubscriptions: string[], subscription: UniformSubscription) => {
      if (subscription.project === projectName || !subscription.project) {
        accSubscriptions.push(subscription.topic);
      }
      return accSubscriptions;
    }, []);
    return subscriptions.length !== 0 ? subscriptions.join('<br/>') : undefined;
  }
}
