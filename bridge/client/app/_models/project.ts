import semver from 'semver';
import { Stage } from './stage';
import { DeploymentInformation, Service } from './service';
import { Trace } from './trace';
import { Deployment } from './deployment';
import { EventTypes } from '../../../shared/interfaces/event-types';
import moment from 'moment';
import { DeploymentStage } from './deployment-stage';
import { Sequence } from './sequence';
import { Project as pj } from '../../../shared/models/project';
import { Approval } from '../_interfaces/approval';

export class Project extends pj {
  allSequencesLoaded = false;
  projectDetailsLoaded = false; // true if project was fetched via project endpoint of bridge server
  stages: Stage[] = [];
  services?: Service[];
  sequences: Sequence[] = [];

  static fromJSON(data: unknown): Project {
    const project: Project = Object.assign(new this(), data);
    project.stages = project.stages.map((stage) => Stage.fromJSON(stage));
    project.setDeployments();
    return project;
  }

  // returns a project without default values
  get reduced(): Partial<Project> {
    const { sequences, allSequencesLoaded, projectDetailsLoaded, ...copyProject } = this;
    return copyProject;
  }

  getServices(stageName?: string): Service[] {
    if (!stageName) {
      if (!this.services) {
        let services: Service[] = [];
        for (const currentStage of this.stages) {
          services = services.concat(
            // eslint-disable-next-line @typescript-eslint/no-loop-func
            currentStage.services.filter((s: Service) => !services.some((ss) => ss.serviceName === s.serviceName))
          );
        }
        this.services = services;
      }
      return this.services;
    } else {
      return this.stages.find((s) => s.stageName === stageName)?.services ?? [];
    }
  }

  getShipyardVersion(): string {
    return this.shipyardVersion?.split('/').pop() ?? '';
  }

  isShipyardNotSupported(supportedVersion: string | undefined): boolean {
    const version = this.getShipyardVersion();
    return !version || !supportedVersion || semver.lt(version, supportedVersion);
  }

  getService(serviceName: string): Service | undefined {
    return this.getServices().find((s) => s.serviceName === serviceName);
  }

  getStage(stageName: string): Stage | undefined {
    return this.stages.find((s) => s.stageName === stageName);
  }

  getLatestDeployment(service?: Service, stage?: Stage): DeploymentInformation | undefined {
    let currentService: Service | undefined;
    if (service) {
      if (stage) {
        currentService = this.getServices(stage.stageName)?.find((s) => s.serviceName === service.serviceName);
      } else {
        currentService = this.getService(service.serviceName);
      }
    }
    return currentService?.deploymentInformation;
  }

  getLatestDeploymentTraceOfSequence(service: Service | undefined, stage?: Stage): Trace | undefined {
    const currentService = service ? this.getService(service.serviceName) : undefined;

    return currentService?.sequences
      ?.find((r) => r.shkeptncontext === currentService.lastEventTypes?.[EventTypes.DEPLOYMENT_FINISHED]?.keptnContext)
      ?.findTrace((trace) => (stage ? trace.isDeployment() === stage.stageName : !!trace.isDeployment()));
  }

  getSequence(service: Service, event: Trace): Sequence | undefined {
    return service.sequences.find((sequence) => sequence.shkeptncontext === event.shkeptncontext);
  }

  getApprovalEvaluation(trace: Trace): Trace | undefined {
    let evaluation: Approval | undefined;
    if (trace.stage) {
      const stage = this.getStage(trace.stage);
      if (stage) {
        evaluation = stage.services.reduce(
          (foundApproval: Approval | undefined, service: Service) =>
            foundApproval || service.openApprovals.find((a) => a.trace.shkeptncontext === trace.shkeptncontext),
          undefined
        );
      }
    }
    return evaluation?.evaluationTrace;
  }

  private setDeployments(): void {
    for (const service of this.getServices()) {
      service.deployments = this.getDeploymentsOfService(service.serviceName);
    }
  }

  private getDeploymentsOfService(serviceName: string): Deployment[] {
    const deployments: Deployment[] = [];
    this.stages.forEach((stage) => {
      const service = stage.services.find((s) => s.serviceName === serviceName);
      if (service?.deploymentContext) {
        const image = service.getImageVersion();
        const deployment = deployments.find(
          (dp) => dp.version === image && dp.shkeptncontext === service.deploymentContext
        );
        const stageDetails = new DeploymentStage(stage.stageName, service.evaluationContext);
        if (deployment) {
          deployment.stages.push(stageDetails);
        } else {
          const newDeployment = new Deployment(image, service.serviceName, stageDetails, service.deploymentContext);
          deployments.push(newDeployment);
        }
      }
    });
    return deployments.sort((a, b) =>
      a.version &&
      b.version &&
      semver.valid(a.version) != null &&
      semver.valid(b.version) != null &&
      semver.gt(a.version, b.version, true)
        ? -1
        : 1
    );
  }

  public getServiceWithLatestDeployment(serviceName: string): Service | undefined {
    let lastService: Service | undefined;
    this.stages.forEach((stage: Stage) => {
      const service = stage.services.find((s) => s.serviceName === serviceName);
      if (
        service?.deploymentContext &&
        (!lastService ||
          (service.deploymentTime &&
            lastService.deploymentTime &&
            moment.unix(service.deploymentTime).isAfter(moment.unix(lastService.deploymentTime))))
      ) {
        lastService = service;
      }
    });
    return lastService;
  }

  public getStages(parent: string[] | null): Stage[] {
    return this.stages.filter(
      (s) => (parent && s.parentStages?.every((element, i) => element === parent[i])) || (!parent && !s.parentStages)
    );
  }

  public getParentStages(): [null, ...string[][]] {
    return this.stages.reduce(
      (stages: [null, ...string[][]], stage) => {
        if (
          stage.parentStages &&
          !stages.find((parent) => parent?.every((element, i) => element === stage.parentStages?.[i]))
        ) {
          stages.push(stage.parentStages);
        }
        return stages;
      },
      [null]
    );
  }
}
