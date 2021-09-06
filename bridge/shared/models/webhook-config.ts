import { WebhookConfig as wc } from '../interfaces/webhook-config';
import { UniformSubscriptionFilter } from '../interfaces/uniform-subscription';

export class WebhookConfig implements wc {

  public type: string;
  public filter!: UniformSubscriptionFilter;
  public prevFilter?: UniformSubscriptionFilter;
  public method: string;
  public url: string;
  public payload: string;
  public header?: { name: string, value: string }[];
  public proxy?: string;

  constructor() {
    this.type = '';
    this.method = '';
    this.url = '';
    this.payload = '';
    this.header = [];
  }
}
