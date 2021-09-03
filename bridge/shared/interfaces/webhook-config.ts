export interface WebhookConfig {
  type: string;
  filter: { projects: string[] | null; stages: string[] | null; services: string[] | null };
  prevFilter: { projects: string[] | null; stages: string[] | null; services: string[] | null } | undefined;
  method: string;
  url: string;
  payload: string;
  header?: {name: string, value: string}[];
  proxy?: string;
}
