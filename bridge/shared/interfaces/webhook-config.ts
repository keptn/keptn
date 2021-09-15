import { UniformSubscriptionFilter } from './uniform-subscription';

export type WebhookConfigMethod = 'POST' | 'PUT';

export interface PreviousWebhookConfig {
  filter: UniformSubscriptionFilter;
  type: string;
}

export interface WebhookConfig {
  type: string;
  filter: UniformSubscriptionFilter;
  prevConfiguration?: PreviousWebhookConfig;
  method: WebhookConfigMethod;
  url: string;
  payload: string;
  header: { name: string, value: string }[];
  proxy?: string;
}
