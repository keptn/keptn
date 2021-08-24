import { UniformRegistration as ur } from '../../../server/interfaces/uniform-registration';
import { UniformSubscription } from './uniform-subscription';
import semver from 'semver/preload';


const preventSubscriptionUpdate = ['approval-service', 'remediation-service', 'lighthouse-service'];

export class UniformRegistration extends ur {
  public subscriptions: UniformSubscription[] = [];
  public unreadEventsCount!: number;

  public static fromJSON(data: unknown): UniformRegistration {
    const uniformRegistration = Object.assign(new this(), data);
    uniformRegistration.subscriptions = uniformRegistration.subscriptions?.map(subscription => UniformSubscription.fromJSON(subscription)) ?? [];
    return uniformRegistration;
  }

  public getSubscriptions(projectName: string): UniformSubscription[] {
    return this.subscriptions.filter(subscription => subscription.hasProject(projectName, true));
  }

  public hasSubscriptions(projectName: string): boolean {
    return this.subscriptions.some(subscription => subscription.hasProject(projectName, true));
  }

  public formatSubscriptions(projectName: string): string | undefined {
    const subscriptions = this.subscriptions.reduce((accSubscriptions: string[], subscription: UniformSubscription) => {
      if (subscription.hasProject(projectName, true)) {
        accSubscriptions.push(subscription.formattedEvent);
      }
      return accSubscriptions;
    }, []);
    return subscriptions.length !== 0 ? subscriptions.join('<br/>') : undefined;
  }

  public canEditSubscriptions(): boolean {
    return !!(semver.valid(this.metadata.distributorversion) && semver.gte(this.metadata.distributorversion, '0.9.0'));
  }

  public isChangeable(): boolean {
    return !preventSubscriptionUpdate.includes(this.name);
  }
}
