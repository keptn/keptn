import { DeploymentInformation } from '../../../shared/models/service-state';
import { Deployment } from '../_models/deployment';

export interface DeploymentInformationSelection {
  deployment: DeploymentInformation;
  stage: string;
}

export interface DeploymentSelection {
  deployment: Deployment;
  stage: string;
}
