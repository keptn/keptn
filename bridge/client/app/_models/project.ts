import semver from 'semver';
import {Stage} from "./stage";
import {Service} from "./service";
import {Trace} from "./trace";
import {Root} from "./root";
import { Deployment } from './deployment';
import {EventTypes} from "./event-types";
import * as moment from 'moment';

export class Project {
  projectName: string;
  gitUser: string;
  gitRemoteURI: string;
  gitToken: string;
  shipyardVersion: string;
  allSequencesLoaded: boolean;

  stages: Stage[];
  services: Service[];
  sequences: Root[];

  getServices(stage?: Stage): Service[] {
    if(this.services && !stage) {
      return this.services;
    } else if(!this.services && !stage) {
      this.services = [];
      this.stages.forEach((stage: Stage) => {
        this.services = this.services.concat(stage.services.filter(s => !this.services.some(ss => ss.serviceName == s.serviceName)));
      });
      return this.services;
    } else {
      return this.stages.find(s => s.stageName == stage.stageName).services;
    }
  }

  getShipyardVersion(): string {
    return this.shipyardVersion?.split('/').pop();
  }

  isShipyardNotSupported(supportedVersion: string): boolean {
    const version = this.getShipyardVersion();
    return !version || !supportedVersion || semver.lt(version, supportedVersion);
  }

  getService(serviceName: string): Service {
    return this.getServices().find(s => s.serviceName == serviceName);
  }

  getStage(stageName: string): Stage {
    return this.stages.find(s => s.stageName == stageName);
  }

  getLatestDeploymentTrace(service: Service, stage?: Stage): Trace {
    let currentService = this.getService(service.serviceName);

    return currentService.roots
      ?.find(r => r.shkeptncontext == currentService.lastEventTypes[EventTypes.DEPLOYMENT_FINISHED]?.keptnContext)
      ?.findTrace(trace => stage ? trace.isDeployment() == stage.stageName : !!trace.isDeployment());
  }

  getLatestFailedRootEvents(stage: Stage): Root[] {
    return this.getServices(stage).map(service => service.getRecentSequence()).filter(seq => seq?.isFailedEvaluation() === stage.stageName);
  }

  getLatestProblemEvents(stage: Stage): Root[] {
    return stage.getOpenProblems();
  }

  getRootEvent(service: Service, event: Trace): Root {
    return service.roots?.find(root => root.shkeptncontext == event.shkeptncontext);
  }

  getDeploymentEvaluation(trace: Trace): Trace {
    let service = this.getServices().find(s => s.serviceName == trace.data.service);
    let root = this.getRootEvent(service, trace);
    return root?.findLastTrace(t => t.isEvaluation() && t.isFinished())?.getFinishedEvent();
  }

  getDeploymentsOfService(serviceName: string): Deployment[] {
    const deployments: Deployment[] = [];
    this.stages.forEach(stage => {
      const service = stage.services.find(service => service.serviceName === serviceName);
      if (service?.deploymentContext) {
        const image = service.getImageVersion();
        const deployment = deployments.find(deployment => deployment.version === image && deployment.shkeptncontext === service.deploymentContext);
        if (deployment) {
          deployment.stages.push(stage.stageName);
        } else {
          const deployment = Deployment.fromJSON({
            version: image,
            service: service.serviceName,
            stages: [stage.stageName],
            shkeptncontext: service.deploymentContext
          } as Deployment);

          deployments.push(deployment);
        }
      }
    });
    return deployments.sort((a, b) => a.version && b.version && semver.gt(a.version, b.version) ? -1 : 1);
  }

  public getLatestDeployment(serviceName: string): Service {
    let lastService: Service;
    this.stages.forEach((stage: Stage) => {
      const service = stage.services.find(s => s.serviceName === serviceName);
      if(service?.deploymentContext && (!lastService || moment.unix(service.deploymentTime).isAfter(moment.unix(lastService.deploymentTime)))) {
        lastService = service;
      }
    });
    return lastService;
  }

  public hasDeployment(serviceName: string): boolean {
    return this.getDeploymentsOfService(serviceName).length !== 0;
  }

  public getStages(parentStage: string): Stage[] {
    return this.stages.filter(s => (parentStage && s.parentStages && s.parentStages.includes(parentStage)) || (!parentStage && !s.parentStages));
  }

  static fromJSON(data: any) {
    return Object.assign(new this, data);
  }
}
