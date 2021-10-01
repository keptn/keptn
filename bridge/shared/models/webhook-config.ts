import { PreviousWebhookConfig, WebhookConfig as wc, WebhookConfigMethod } from '../interfaces/webhook-config';
import { UniformSubscriptionFilter } from '../interfaces/uniform-subscription';

export type WebhookConfigFilter = { projects: string[], stages: string[] | [undefined], services: string[] | [undefined] };

export class WebhookConfig implements wc {

  public type: string;
  public filter!: UniformSubscriptionFilter;
  public prevConfiguration?: PreviousWebhookConfig;
  public method: WebhookConfigMethod;
  public url: string;
  public payload: string;
  public header: { name: string, value: string }[];
  public proxy?: string;

  constructor() {
    this.type = '';
    this.method = 'POST';
    this.url = '';
    this.payload = '';
    this.header = [];
  }
}
