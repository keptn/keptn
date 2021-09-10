import Yaml from 'yaml';
import { WebhookConfigMethod } from '../../shared/interfaces/webhook-config';
import { WebhookConfig } from '../../shared/models/webhook-config';

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

  public parsedRequest(eventType: string): WebhookConfig | undefined {
    const curl = this.spec.webhooks.find(w => w.type === eventType)?.requests[0];
    return curl ? this.parseCurl(curl) : undefined;
  }

  private parseCurl(curl: string): WebhookConfig {
    const config = new WebhookConfig();
    config.url = curl.match(/.* (.*)$/)?.[1] ?? '';
    config.payload = this.formatJSON(this.getCommandData('--data', curl).data);
    config.proxy = this.getCommandData('--proxy', curl).data;
    config.method = this.getCommandData('--request', curl).data as WebhookConfigMethod;
    config.header = this.getHeaders('--header', curl);
    return config;
  }

  private getHeaders(arg: string, command: string): { name: string, value: string }[] {
    let index = 0;
    const headers: { name: string, value: string }[] = [];
    while (index !== -1) {
      const result = this.getCommandData(arg, command, index);
      index = result.index;
      if (result.data) {
        const headerInfo = result.data.split(':');

        headers.push({
          name: headerInfo[0]?.trim(),
          value: headerInfo[1]?.trim(),
        });
      }
    }
    return headers;
  }

  private getCommandData(arg: string, command: string, fromIndex = 0): { data: string, index: number } {
    arg = `${arg} `;
    const dataIndex = command.indexOf(arg, fromIndex);
    let data = '';
    let startIndex = -1;
    if (dataIndex > -1) {
      startIndex = dataIndex + arg.length;
      const chars = [...command.substring(startIndex)];
      const startsWith = chars[0];
      if (startsWith === '\'' || startsWith === '\"') {
        let i = 1;
        for (; i < chars.length && (chars[i] !== startsWith || chars[i] === startsWith && chars[i - 1] === '\\'); ++i) {
        }
        data = command.substring(startIndex + 1, i + startIndex);
      } else {
        data = command.substring(startIndex, chars.findIndex(c => c === ' ') + startIndex);
      }
    }
    const formattedData = data.replace(/\\"/g, '"');
    return {data: formattedData, index: startIndex};
  }

  private formatJSON(data: string): string {
    try {
      data = JSON.stringify(JSON.parse(data), null, 2);
    } catch {
    }
    return data;
  }

  public toYAML(): string {
    return Yaml.stringify(this, {
      sortMapEntries: (a, b) => {
        return order[a.key] - order[b.key];
      },
    });
  }
}
