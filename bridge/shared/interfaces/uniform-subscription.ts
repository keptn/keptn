interface IUniformSubscriptionFilterParameter {
  key: string;
  value: string;
  visible: boolean;
}

export interface IUniformSubscriptionFilter {
  projects: string[] | null;
  stages: string[] | null;
  services: string[] | null;
}

export interface IUniformSubscription {
  event: string;
  filter: IUniformSubscriptionFilter;
  parameters?: IUniformSubscriptionFilterParameter[];
  id?: string;
}
