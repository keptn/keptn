import { UniformSubscriptionFilter } from './uniform-subscription';

export type WebhookConfigMethod = 'POST' | 'PUT' | 'GET';

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
