export interface IWebhookConfigFilter {
  projects: string[];
  stages: string[] | [undefined];
  services: string[] | [undefined];
}
