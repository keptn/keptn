import { DeploymentInformation } from '../_models/service-state';
import { Deployment } from '../_models/deployment';

export interface DeploymentInformationSelection {
  deploymentInformation: DeploymentInformation;
  stage: string;
}

export interface DeploymentSelection {
  deployment: Deployment;
  stage: string;
}
