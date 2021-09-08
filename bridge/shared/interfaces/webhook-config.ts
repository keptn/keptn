import { UniformSubscriptionFilter } from './uniform-subscription';

export type WebhookConfigMethod = 'POST' | 'PUT';

export interface WebhookConfig {
  type: string;
  filter: UniformSubscriptionFilter;
  prevFilter?: UniformSubscriptionFilter;
  method: WebhookConfigMethod;
  url: string;
  payload: string;
  header?: { name: string, value: string }[];
  proxy?: string;
}
