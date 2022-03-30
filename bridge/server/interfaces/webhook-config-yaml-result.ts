import { WebhookSecret } from '../../shared/models/webhook-config';

export type Webhook = {
  subscriptionID: string;
  sendFinished?: boolean;
  sendStarted?: boolean;
  type: string; // type === event
  requests: string[];
  envFrom?: WebhookSecret[];
};

export interface WebhookConfigYamlResult {
  apiVersion: 'webhookconfig.keptn.sh/v1alpha1';
  kind: 'WebhookConfig';
  metadata: {
    name: 'webhook-configuration';
  };
  spec: {
    webhooks: Webhook[];
  };
}
