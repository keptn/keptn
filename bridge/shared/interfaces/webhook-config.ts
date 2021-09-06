import { UniformSubscriptionFilter } from './uniform-subscription';

export interface WebhookConfig {
  type: string;
  filter: UniformSubscriptionFilter;
  prevFilter?: UniformSubscriptionFilter;
  method: string;
  url: string;
  payload: string;
  header?: { name: string, value: string }[];
  proxy?: string;
}
