import { Document, Scalar, YAMLMap, YAMLSeq } from 'yaml';
import { WebhookConfigMethod } from '../../shared/interfaces/webhook-config';
import { WebhookConfig, WebhookSecret } from '../../shared/models/webhook-config';
import { Webhook, WebhookConfigYamlResult } from '../interfaces/webhook-config-yaml-result';
import { parseCurl } from '../utils/curl.utils';

const order: { [key in keyof WebhookConfigYamlResult]: number } = {
  apiVersion: 0,
  kind: 1,
  metadata: 2,
  spec: 3,
};

export class WebhookConfigYaml implements WebhookConfigYamlResult {
  apiVersion: 'webhookconfig.keptn.sh/v1alpha1';
  kind: 'WebhookConfig';
  metadata: {
    name: 'webhook-configuration';
  };
  spec: {
    webhooks: Webhook[];
  };

  constructor() {
    this.spec = {
      webhooks: [],
    };
    this.metadata = {
      name: 'webhook-configuration',
    };
    this.apiVersion = 'webhookconfig.keptn.sh/v1alpha1';
    this.kind = 'WebhookConfig';
  }

  public static fromJSON(data: WebhookConfigYamlResult): WebhookConfigYaml {
    return Object.assign(new this(), data);
  }

  /**
   * @params subscriptionId
   * @returns true if the webhooks have been changed
   */
  public removeWebhook(subscriptionId: string): boolean {
    const index = this.getWebhookIndex(subscriptionId);
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
   * Either adds a webhook or updates it if there is already one for the given subscriptionId
   * @params eventType
   * @params curl
   */
  public addWebhook(
    eventType: string,
    curl: string,
    subscriptionId: string,
    secrets: WebhookSecret[],
    sendFinished: boolean,
    sendStarted: boolean
  ): void {
    const webhook = this.getWebhook(subscriptionId);
    if (!webhook) {
      this.spec.webhooks.push({
        type: eventType,
        requests: [curl],
        ...(secrets.length && { envFrom: secrets }),
        subscriptionID: subscriptionId,
        sendFinished: sendFinished,
        sendStarted: sendStarted,
      });
    } else {
      // overwrite
      webhook.type = eventType;
      webhook.requests[0] = curl;
      webhook.sendFinished = sendFinished;
      webhook.sendStarted = sendStarted;
      if (secrets.length) {
        webhook.envFrom = secrets;
      } else {
        delete webhook.envFrom;
      }
    }
  }

  private getWebhook(subscriptionId: string): Webhook | undefined {
    return this.spec.webhooks.find(this.findWebhook(subscriptionId));
  }

  private getWebhookIndex(subscriptionId: string): number {
    return this.spec.webhooks.findIndex(this.findWebhook(subscriptionId));
  }

  private findWebhook(subscriptionId: string): (webhook: Webhook) => boolean {
    return (webhook: Webhook): boolean => webhook.subscriptionID === subscriptionId;
  }

  public parsedRequest(subscriptionId: string): WebhookConfig | undefined {
    const webhook = this.getWebhook(subscriptionId);
    if (webhook) {
      const curl = webhook.requests[0];
      if (curl) {
        const parsedConfig = this.parseConfig(curl);
        parsedConfig.secrets = webhook.envFrom;
        parsedConfig.sendFinished = webhook.sendFinished ?? false;
        parsedConfig.sendStarted = webhook.sendStarted ?? true;
        parsedConfig.type = webhook.type;
        return parsedConfig;
      }
    }
    return undefined;
  }

  private parseConfig(curl: string): WebhookConfig {
    const config = new WebhookConfig();
    const result = parseCurl(curl);
    config.url = result._?.join(' ') ?? '';
    config.payload = this.formatJSON(result.data?.[0] ?? '');
    config.proxy = result.proxy?.[0] ?? '';
    config.method = (result.request?.[0] ?? '') as WebhookConfigMethod;
    const headers: { name: string; value: string }[] = [];
    if (result.header) {
      for (const header of result.header) {
        const headerInfo = header.split(':');

        headers.push({
          name: headerInfo[0]?.trim(),
          value: headerInfo[1]?.trim(),
        });
      }
    }

    config.header = headers;
    return config;
  }

  private formatJSON(data: string): string {
    try {
      data = JSON.stringify(JSON.parse(data), null, 2);
    } catch {}
    return data;
  }

  public toYAML(): string {
    const yamlDoc = new Document(this, {
      sortMapEntries: (a, b): number =>
        order[a.key as keyof WebhookConfigYamlResult] - order[b.key as keyof WebhookConfigYamlResult],
      toStringDefaults: {
        lineWidth: 0,
      },
    });
    this.setCurlToBlockFolded(yamlDoc);
    return yamlDoc.toString();
  }

  private setCurlToBlockFolded(yamlDoc: Document): void {
    const yamlSeq = yamlDoc.getIn(['spec', 'webhooks'], true) as YAMLSeq;
    for (const webhookYaml of yamlSeq.items) {
      const requests = (webhookYaml as YAMLMap).get('requests', true) as unknown as YAMLSeq;
      for (const curl of requests.items) {
        (curl as Scalar).type = 'BLOCK_FOLDED';
      }
    }
  }
}
