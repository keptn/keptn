import {Subscription} from './subscription';

export class KeptnService {
  public name: string;
  public status: string;
  public host: string;
  public location: string;
  public version: string;
  public namespace: string;
  public subscriptions: Subscription[] = [];

  static fromJSON(data: any) {
    const service = Object.assign(new this, data);
    service.subscriptions = service.subscriptions.map(subscription => Subscription.fromJSON(subscription));
    return service;
  }

  public deleteSubscription(index: number) {
    this.subscriptions.splice(index, 1);
  }

  public addSubscription(subscription: Subscription) {
    this.subscriptions.push(subscription);
  }
}
