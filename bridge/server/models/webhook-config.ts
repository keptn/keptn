import { WebhookConfig as wc } from '../../shared/interfaces/webhook-config';

export class WebhookConfig implements wc {

  public type: string;
  public filter!: { projects: string[] | null; stages: string[] | null; services: string[] | null };
  public prevFilter: { projects: string[] | null; stages: string[] | null; services: string[] | null } | undefined;
  public method: string;
  public url: string;
  public payload: string;
  public header?: {name: string, value: string}[];
  public proxy?: string;

  constructor() {
    this.type = '';
    this.method = '';
    this.url = '';
    this.payload = '';
    this.header = [];
  }

  public static fromJSON(data: unknown): WebhookConfig {
    return Object.assign(new this(), data);
  }
}
