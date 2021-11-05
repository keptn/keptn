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
      // deployments.length === 0 means that there aren't any updates for a service
      if (newServiceState.deploymentInformation.length) {
        const serviceStateOriginal = serviceStates.find((serviceStateO) => serviceStateO.name === newServiceState.name);
        if (serviceStateOriginal) {
          serviceStateOriginal.update(newServiceState);
        } else {
          // new service with deployments
          serviceStates.push(newServiceState);
        }
      } else if (!serviceStates.some((s) => s.name === newServiceState.name)) {
        // new service
        serviceStates.push(newServiceState);
      }
    }

    this.deleteOldServices(serviceStates, newServiceStates);
  }

  private static deleteOldServices(serviceStates: ServiceState[], newServiceStates: ServiceState[]): void {
    for (let i = 0; i < serviceStates.length; ) {
      const serviceStateOriginal = serviceStates[i];
      if (!newServiceStates.some((state) => state.name === serviceStateOriginal.name)) {
        serviceStates.splice(i, 1);
        i--;
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
        this.updateDeploymentInformation(deploymentOriginal, deploymentNew);
      } else {
        // add new deployment
        this.deploymentInformation.push(deploymentNew);
      }
    }

    this.deleteOldDeployments(serviceState.deploymentInformation);
  }

  private updateDeploymentInformation(
    deploymentOriginal: ServiceDeploymentInformation,
    deploymentNew: ServiceDeploymentInformation
  ): void {
    // update existing deployment
    deploymentOriginal.stages = [...deploymentOriginal.stages, ...deploymentNew.stages];

    // update other deployments (remove the stages)
    for (let i = 0; i < this.deploymentInformation.length; ++i) {
      const deployment = this.deploymentInformation[i];
      if (deployment !== deploymentOriginal) {
        deployment.stages = deployment.stages.filter((stage) =>
          deploymentNew.stages.some((st) => st.name === stage.name)
        );
        // delete deployment if it does not exist anymore
        if (deployment.stages.length === 0) {
          this.deploymentInformation.splice(i, 1);
        }
      }
    }
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

  public hasRemediations(): boolean {
    return this.deploymentInformation.some((deployment) =>
      deployment.stages.some((stage) => stage.hasOpenRemediations)
    );
  }

  public getLatestDeploymentTime(): Date | undefined {
    let latestTime: Date | undefined;
    for (const deployment of this.deploymentInformation) {
      for (const stage of deployment.stages) {
        const date = new Date(stage.time);
        if (!latestTime || date > latestTime) {
          latestTime = date;
        }
      }
    }
    return latestTime;
  }
}
