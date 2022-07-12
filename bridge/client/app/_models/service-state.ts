import { ServiceDeploymentInformation as sdi, ServiceState as svs } from '../../../shared/models/service-state';
import { Deployment } from './deployment';

export interface ServiceDeploymentInformation extends sdi {
  deployment?: Deployment;
}

export class ServiceState extends svs {
  name!: string;
  deploymentInformation: ServiceDeploymentInformation[] = [];

  public static fromJSON(states: svs): ServiceState {
    return Object.assign(new this(svs.name), states);
  }

  public static update(serviceStates: ServiceState[], newServiceStates: ServiceState[]): void {
    for (const newServiceState of newServiceStates) {
      const serviceStateOriginal = serviceStates.find((serviceStateO) => serviceStateO.name === newServiceState.name);
      if (serviceStateOriginal) {
        serviceStateOriginal.update(newServiceState);
      } else {
        // new service with deployments
        serviceStates.push(newServiceState);
      }
    }

    this.deleteOldServices(serviceStates, newServiceStates);
    serviceStates.sort((a, b) => a.name.localeCompare(b.name));
  }

  private static deleteOldServices(serviceStates: ServiceState[], newServiceStates: ServiceState[]): void {
    for (let i = 0; i < serviceStates.length; ) {
      const serviceStateOriginal = serviceStates[i];
      if (!newServiceStates.some((state) => state.name === serviceStateOriginal.name)) {
        serviceStates.splice(i, 1);
        --i;
      }
      ++i;
    }
  }

  public update(serviceState: ServiceState): void {
    for (const deploymentNew of serviceState.deploymentInformation) {
      const deploymentOriginal = this.deploymentInformation.find(
        (deployment) => deployment.keptnContext === deploymentNew.keptnContext
      );
      if (deploymentOriginal) {
        deploymentOriginal.stages = deploymentNew.stages;
      } else {
        // add new deployment
        this.deploymentInformation.push(deploymentNew);
      }
    }

    this.deleteOldDeployments(serviceState.deploymentInformation);
  }

  private deleteOldDeployments(deploymentInformation: ServiceDeploymentInformation[]): void {
    for (let i = 0; i < this.deploymentInformation.length; ) {
      const deploymentOld = this.deploymentInformation[i];
      if (!deploymentInformation.some((st) => st.keptnContext === deploymentOld.keptnContext)) {
        this.deploymentInformation.splice(i, 1);
        --i;
      }
      ++i;
    }
  }

  public getLatestImage(): string {
    const unknownImage = 'unknown';
    let latestTime: Date | undefined;
    let image = unknownImage;
    for (const deployment of this.deploymentInformation) {
      const latestStageTime = deployment.stages.reduce((max: undefined | Date, stage) => {
        const date = new Date(stage.time);
        return max && max > date ? max : date;
      }, undefined);
      if (latestStageTime && (!latestTime || latestStageTime > latestTime)) {
        image = deployment.image ? `${deployment.image}:${deployment.version}` : unknownImage;
        latestTime = latestStageTime;
      }
    }
    return image;
  }
}
