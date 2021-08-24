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
      } else if (!b.project) {
        status = 1;
      } else {
        status = a.event.localeCompare(b.event);
      }
      return status;
    });
    return subscriptions;
  }

  public hasSubscriptions(projectName: string): boolean {
    return this.subscriptions.some(subscription => subscription.project === projectName || !subscription.project);
  }
}
