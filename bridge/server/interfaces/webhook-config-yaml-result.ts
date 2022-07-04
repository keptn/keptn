import { IWebhookHeader, WebhookConfigMethod } from '../../shared/interfaces/webhook-config';

export enum WebhookApiVersions {
  V1ALPHA1 = 'webhookconfig.keptn.sh/v1alpha1',
  V1BETA1 = 'webhookconfig.keptn.sh/v1beta1',
}

export interface IWebhookSecret {
  name: string;
  secretRef: {
    name: string;
    key: string;
  };
}

export interface IWebhookV1Alpha1 {
  subscriptionID: string;
  sendFinished?: boolean;
  sendStarted?: boolean;
  type: string; // type === event
  requests: string[];
  envFrom?: IWebhookSecret[];
}

export interface IWebhookRequestV1Beta1 {
  url: string;
  method: WebhookConfigMethod;
  headers?: IWebhookHeader[];
  payload?: string;
  options?: string;
}

export interface IWebhookV1Beta1 {
  subscriptionID: string;
  sendFinished?: boolean;
  sendStarted?: boolean;
  type: string; // type === event
  requests: IWebhookRequestV1Beta1[];
  envFrom?: IWebhookSecret[];
}

export interface IWebhookConfigYamlResultV1Alpha1 {
  apiVersion: WebhookApiVersions.V1ALPHA1;
  kind: 'WebhookConfig';
  metadata: {
    name: 'webhook-configuration';
  };
  spec: {
    webhooks: IWebhookV1Alpha1[];
  };
}

export interface IWebhookConfigYamlResultV1Beta1 {
  apiVersion: WebhookApiVersions.V1BETA1;
  kind: 'WebhookConfig';
  metadata: {
    name: 'webhook-configuration';
  };
  spec: {
    webhooks: IWebhookV1Beta1[];
  };
}

export type IWebhookConfigYamlResult = IWebhookConfigYamlResultV1Alpha1 | IWebhookConfigYamlResultV1Beta1;
