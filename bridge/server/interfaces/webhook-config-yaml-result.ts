export type Webhook = {
  subscriptionId: string,
  type: string, // type === event
  requests: string[]
};

export interface WebhookConfigYamlResult {
  apiVersion: 'webhookconfig.keptn.sh/v1alpha1';
  kind: 'WebhookConfig';
  metadata: {
    name: 'webhook-configuration'
  };
  spec: {
    webhooks: Webhook[]
  };
}
