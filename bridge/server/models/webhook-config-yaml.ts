import { Document, Pair, Scalar, YAMLMap, YAMLSeq } from 'yaml';
import { IWebhookConfigClient } from '../../shared/interfaces/webhook-config';
import { IWebhookSecret } from '../interfaces/webhook-config';
import {
  IWebhookConfigYamlResult,
  IWebhookConfigYamlResultV1Beta1,
  IWebhookRequestV1Beta1,
  IWebhookV1Beta1,
  WebhookApiVersions,
} from '../interfaces/webhook-config-yaml-result';
import { stringifyPayload } from './webhook-config.utils';

const order: { [key in keyof IWebhookConfigYamlResultV1Beta1]: number } = {
  apiVersion: 0,
  kind: 1,
  metadata: 2,
  spec: 3,
};

export class WebhookConfigYaml implements IWebhookConfigYamlResultV1Beta1 {
  apiVersion: WebhookApiVersions.V1BETA1;
  kind: 'WebhookConfig';
  metadata: {
    name: 'webhook-configuration';
  };
  spec: {
    webhooks: IWebhookV1Beta1[];
  };

  constructor() {
    this.spec = {
      webhooks: [],
    };
    this.metadata = {
      name: 'webhook-configuration',
    };
    this.apiVersion = WebhookApiVersions.V1BETA1;
    this.kind = 'WebhookConfig';
  }

  public static fromJSON(data: IWebhookConfigYamlResultV1Beta1): WebhookConfigYaml {
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
  public addWebhook(config: IWebhookConfigClient, subscriptionId: string, secrets: IWebhookSecret[]): void {
    const webhook = this.getWebhook(subscriptionId);
    const request: IWebhookRequestV1Beta1 = {
      ...(config.proxy && { options: `--proxy ${config.proxy}` }),
      url: config.url,
      payload: stringifyPayload(config.payload),
      headers: config.header,
      method: config.method,
    };
    if (!webhook) {
      this.spec.webhooks.push({
        type: config.type,
        requests: [request],
        ...(secrets.length && { envFrom: secrets }),
        subscriptionID: subscriptionId,
        sendFinished: config.sendFinished,
        sendStarted: config.sendStarted,
      });
    } else {
      // overwrite
      webhook.type = config.type;
      webhook.requests[0] = request;
      webhook.sendFinished = config.sendFinished;
      webhook.sendStarted = config.sendStarted;
      if (secrets.length) {
        webhook.envFrom = secrets;
      } else {
        delete webhook.envFrom;
      }
    }
  }

  public getWebhook(subscriptionId: string): IWebhookV1Beta1 | undefined {
    return this.spec.webhooks.find(this.findWebhook(subscriptionId));
  }

  private getWebhookIndex(subscriptionId: string): number {
    return this.spec.webhooks.findIndex(this.findWebhook(subscriptionId));
  }

  private findWebhook(subscriptionId: string): (webhook: IWebhookV1Beta1) => boolean {
    return (webhook: IWebhookV1Beta1): boolean => webhook.subscriptionID === subscriptionId;
  }

  public toYAML(): string {
    const yamlDoc = new Document(this, {
      sortMapEntries: (a: Pair, b: Pair): number =>
        order[a.key as keyof IWebhookConfigYamlResult] - order[b.key as keyof IWebhookConfigYamlResult],
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
      for (const requestConfig of requests.items) {
        const payload = (requestConfig as YAMLMap).get('payload', true);
        (payload as Scalar).type = 'BLOCK_FOLDED';
      }
    }
  }
}
