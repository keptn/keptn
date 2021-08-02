import semver from 'semver';
import {Stage} from './stage';
import {Service} from './service';
import {Trace} from './trace';
import {Root} from './root';
import { Deployment } from './deployment';
import {EventTypes} from './event-types';
import moment from 'moment';
import {DeploymentStage} from './deployment-stage';
import {Sequence} from './sequence';

export class Project {
  projectName!: string;
  gitUser?: string;
  gitRemoteURI?: string;
  shipyardVersion?: string;
  allSequencesLoaded = false;
  stages: Stage[] = [];
  services?: Service[];
  roots: Root[] = [];
  sequences: Sequence[] = [];

  static fromJSON(data: unknown) {
    const project: Project = Object.assign(new this(), data);
    project.stages = project.stages.map(stage => {
      stage.services = stage.services.map(service => {
        service.stage = stage.stageName;
        return Service.fromJSON(service);
      });
      return Stage.fromJSON(stage);
    });
    project.setDeployments();
    return project;
  }

  getServices(stageName?: string): Service[] {
    if (!stageName) {
      if (!this.services) {
        let services: Service[] = [];
        for (const currentStage of this.stages){
          services = services.concat(
            currentStage.services.filter(s => !services.some(ss => ss.serviceName === s.serviceName))
          );
        }
        this.services = services;
      }
      return this.services;
    }
    else {
      return this.stages.find(s => s.stageName === stageName)?.services ?? [];
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
    return this.getServices().find(s => s.serviceName === serviceName);
  }

  getStage(stageName: string): Stage | undefined {
    return this.stages.find(s => s.stageName === stageName);
  }

  getLatestDeploymentTrace(service: Service, stage?: Stage): Trace | undefined {
    const currentService = this.getService(service.serviceName);

    return currentService?.roots
      ?.find(r => r.shkeptncontext === currentService.lastEventTypes?.[EventTypes.DEPLOYMENT_FINISHED]?.keptnContext)
      ?.findTrace(trace => stage ? trace.isDeployment() === stage.stageName : !!trace.isDeployment());
  }

  getLatestDeploymentTraceOfSequence(service: Service, stage?: Stage): Trace | undefined{
    const currentService = this.getService(service.serviceName);

    return currentService?.sequences
      ?.find(r => r.shkeptncontext === currentService.lastEventTypes?.[EventTypes.DEPLOYMENT_FINISHED]?.keptnContext)
      ?.findTrace(trace => stage ? trace.isDeployment() === stage.stageName : !!trace.isDeployment());
  }

  getLatestFailedRootEvents(stage: Stage): Root[] {
    return this.getServices(stage.stageName).map(service => service.getRecentRoot()).filter(seq => seq?.hasFailedEvaluation() === stage.stageName);
  }

  getLatestProblemEvents(stage: Stage): Root[] {
    return stage.getOpenProblems();
  }

  getSequence(service: Service, event: Trace): Sequence | undefined {
    return service.sequences.find(sequence => sequence.shkeptncontext === event.shkeptncontext);
  }

  getRootEvent(service: Service, event: Trace): Root | undefined {
    return service.roots.find(root => root.shkeptncontext === event.shkeptncontext);
  }

  getDeploymentEvaluation(trace: Trace, isSequence: boolean): Trace | undefined {
    const service = this.getServices().find(s => s.serviceName === trace.data.service);
    let root: Sequence | Root | undefined;
    if (service) {
      root = (isSequence ? this.getSequence : this.getRootEvent)(service, trace);
    }
    return root?.findLastTrace(t => !!(t.isEvaluation() && t.isFinished()))?.getFinishedEvent();
  }

  private setDeployments() {
    for (const service of this.getServices()) {
      service.deployments = this.getDeploymentsOfService(service.serviceName);
    }
  }

  private getDeploymentsOfService(serviceName: string): Deployment[] {
    const deployments: Deployment[] = [];
    this.stages.forEach(stage => {
      const service = stage.services.find(s => s.serviceName === serviceName);
      if (service?.deploymentContext) {
        const image = service.getImageVersion();
        const deployment = deployments.find(dp => dp.version === image && dp.shkeptncontext === service.deploymentContext);
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
      a.version && b.version && semver.valid(a.version) != null &&
      semver.valid(b.version) != null && semver.gt(a.version, b.version, true) ? -1 : 1);
  }

  public getLatestDeployment(serviceName: string): Service | undefined {
    let lastService: Service | undefined;
    this.stages.forEach((stage: Stage) => {
      const service = stage.services.find(s => s.serviceName === serviceName);
      if (service?.deploymentContext &&
        (!lastService
          || service.deploymentTime && lastService.deploymentTime
            && moment.unix(service.deploymentTime).isAfter(moment.unix(lastService.deploymentTime)))) {
        lastService = service;
      }
    });
    return lastService;
  }

  public getStages(parent: (string | null)[]): Stage[] {
    return this.stages.filter(s => (
      parent && s.parentStages?.every((element, i) => element === parent[i]))
      || (!parent && !s.parentStages));
  }

  public getParentStages(): [null, ...string[][]] {
    return this.stages.reduce((stages: ([null, ...string[][]]), stage) => {
      if (stage.parentStages &&
        !stages.find((parent) =>
          parent?.every((element, i) => element === stage.parentStages?.[i]))) {
        stages.push(stage.parentStages);
      }
      return stages;
    }, [null]);
  }
}
