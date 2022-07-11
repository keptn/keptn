export interface ServiceDeploymentInformation {
  version?: string;
  name: string;
  image?: string;
  stages: {
    name: string;
    time: string; // ISO-string
  }[];
  keptnContext: string;
}

export class ServiceState {
  name: string;
  deploymentInformation: ServiceDeploymentInformation[] = [];

  constructor(name: string) {
    this.name = name;
  }
}
