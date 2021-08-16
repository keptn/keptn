export interface UniformSubscription {
  topic: string;
  filter: {
    projects: string[] | null,
    stages: string[] | null,
    services: string[] | null
  };
  parameters?: {key: string, value: string, visible: boolean}[];
}
