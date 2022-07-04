export interface IUniformSubscriptionFilter {
  projects: string[] | null;
  stages: string[] | null;
  services: string[] | null;
}

export interface IUniformSubscription {
  event: string;
  filter: IUniformSubscriptionFilter;
  parameters?: { key: string; value: string; visible: boolean }[];
  id?: string;
}
