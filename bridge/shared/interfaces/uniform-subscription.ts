export type UniformSubscriptionFilter = {
  projects: string[] | null;
  stages: string[] | null;
  services: string[] | null;
};

export interface UniformSubscription {
  event: string;
  filter: UniformSubscriptionFilter;
  parameters?: { key: string; value: string; visible: boolean }[];
  id?: string;
}
