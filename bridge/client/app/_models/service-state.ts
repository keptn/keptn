import { DeploymentInformation as di, ServiceState as svs } from '../../../shared/models/service-state';
import { Deployment } from './deployment';

export interface DeploymentInformation extends di {
  deployment?: Deployment;
}

export class ServiceState extends svs {
  name!: string;
  deploymentInformation: DeploymentInformation[] = [];
}
