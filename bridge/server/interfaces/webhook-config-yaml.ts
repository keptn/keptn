import Yaml from 'yaml';

const order: { [key: string]: number } = {
  apiVersion: 0,
  kind: 1,
  metadata: 2,
  spec: 3,
};

export class WebhookConfigYaml {
  apiVersion: 'webhookconfig.keptn.sh/v1alpha1';
  kind: 'WebhookConfig';
  metadata: {
    name: 'webhook-configuration'
  };
  spec: {
    webhooks: {
      type: string, // type === event
      requests: string[]
    } []
  };

  constructor(type?: string, request?: string) {
    this.spec = {
      webhooks: [],
    };
    this.metadata = {
      name: 'webhook-configuration',
    };
    this.apiVersion = 'webhookconfig.keptn.sh/v1alpha1';
    this.kind = 'WebhookConfig';
    if (type && request) {
      this.spec.webhooks.push({
        requests: [request],
        type,
      });
    }
  }

  public static fromJSON(data: unknown): WebhookConfigYaml {
    return Object.assign(new this(), data);
  }

  /**
   * @params eventType
   * @returns true if the webhooks have been changed
   */
  public removeWebhook(eventType: string): boolean {
    const index = this.spec.webhooks.findIndex(webhook => webhook.type === eventType);
    const changed = index !== -1;
    if (changed) {
      this.spec.webhooks[index].requests.splice(0, 1);
      if (this.spec.webhooks[index].requests.length === 0) {
        this.spec.webhooks.splice(index, 1);
      }
    }
    return changed;
  }

  public hasWebhooks(): boolean {
    return this.spec.webhooks.length !== 0;
  }

  /**
   * Either adds a webhook or overwrites it if there is already one for the given eventType
   * @params eventType
   * @params curl
   */
  public addWebhook(eventType: string, curl: string): void {
    const webhook = this.spec.webhooks.find(w => w.type === eventType);
    if (!webhook) {
      this.spec.webhooks.push({type: eventType, requests: [curl]});
    } else { // overwrite
      webhook.type = eventType;
      webhook.requests = [curl];
    }

  }

  public toYAML(): string {
    return Yaml.stringify(this, {
      sortMapEntries: (a, b) => {
        return order[a.key] - order[b.key];
      },
    });
  }
}
