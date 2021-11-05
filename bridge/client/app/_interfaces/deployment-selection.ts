import { ServiceDeploymentInformation } from '../_models/service-state';
import { Deployment } from '../_models/deployment';

export interface DeploymentInformationSelection {
  deploymentInformation: ServiceDeploymentInformation;
  stage: string;
}

export interface DeploymentSelection {
  deployment: Deployment;
  stage: string;
}
