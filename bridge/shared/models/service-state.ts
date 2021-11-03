export interface DeploymentInformation {
  version?: string;
  name: string;
  image?: string;
  stages: {
    name: string;
    hasOpenRemediations: boolean;
    time: string; // ISO-string
  }[];
  keptnContext: string;
}

export class ServiceState {
  name: string;
  deployments: DeploymentInformation[] = [];

  constructor(name: string) {
    this.name = name;
  }
}
