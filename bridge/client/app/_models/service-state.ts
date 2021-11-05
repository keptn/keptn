import { DeploymentInformation as di, ServiceState as svs } from '../../../shared/models/service-state';
import { Deployment } from './deployment';

export interface DeploymentInformation extends di {
  deployment?: Deployment;
}

export class ServiceState extends svs {
  name!: string;
  deploymentInformation: DeploymentInformation[] = [];

  public static fromJSON(states: svs): ServiceState {
    return Object.assign(new this(svs.name), states);
  }

  public static update(serviceStates: ServiceState[], newServiceStates: ServiceState[]): void {
    // TODO: check if update works
    for (const newServiceState of newServiceStates) {
      // deployments.length === 0 means that there aren't any updates for a service
      if (newServiceState.deploymentInformation.length) {
        const serviceStateOriginal = serviceStates.find((serviceStateO) => serviceStateO.name === newServiceState.name);
        serviceStateOriginal?.update(newServiceState);
        if (!serviceStateOriginal) {
          // new service with deployments
          serviceStates.push(newServiceState);
        }
      } else if (!serviceStates.some((s) => s.name === newServiceState.name)) {
        // new service
        serviceStates.push(newServiceState);
      }
      // remove deleted services
      for (let i = 0; i < serviceStates.length; ++i) {
        const serviceStateOriginal = serviceStates[i];
        if (!serviceStates.some((state) => state.name === serviceStateOriginal.name)) {
          serviceStates.splice(i, 1);
        }
      }
    }
  }

  public update(serviceState: ServiceState): void {
    for (const deploymentNew of serviceState.deploymentInformation) {
      const deploymentOriginal = this.deploymentInformation.find(
        (deployment) => deployment.keptnContext === deploymentNew.keptnContext
      );
      if (deploymentOriginal) {
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
      } else {
        // add new deployment
        this.deploymentInformation.push(deploymentNew);
      }
    }
  }
}
