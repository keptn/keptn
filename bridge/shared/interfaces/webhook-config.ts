import { IUniformSubscriptionFilter } from './uniform-subscription';

export type WebhookConfigMethod = 'POST' | 'PUT' | 'GET';

export interface PreviousWebhookConfig {
  filter: IUniformSubscriptionFilter;
  type: string;
}

export interface IWebhookHeader {
  key: string;
  value: string;
}

export interface IWebhookConfigClient {
  type: string;
  prevConfiguration?: PreviousWebhookConfig; // send from client to server
  method: WebhookConfigMethod;
  url: string;
  payload: string;
  header: IWebhookHeader[];
  proxy?: string;
  sendFinished: boolean;
  sendStarted: boolean;
}
