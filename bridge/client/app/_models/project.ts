import semver from 'semver';
import {Stage} from "./stage";
import {Service} from "./service";
import {Trace} from "./trace";
import {Root} from "./root";
import { Deployment } from './Deployment';

export class Project {
  projectName: string;
  gitUser: string;
  gitRemoteURI: string;
  gitToken: string;
  shipyardVersion: string;

  stages: Stage[];
  services: Service[];
  sequences: Root[];

  getServices(): Service[] {
    if(!this.services) {
      this.services = [];
      this.stages.forEach((stage: Stage) => {
        this.services = this.services.concat(stage.services.filter(s => !this.services.some(ss => ss.serviceName == s.serviceName)));
      });
    }
    return this.services;
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

  getLatestDeployment(service: Service, stage?: Stage): Trace {
    let currentService = this.getService(service.serviceName);

    if(currentService.roots)
      return currentService.roots
        .filter(root => !stage || root.isFaulty() != stage.stageName || root.getDeploymentDetails(stage)?.isDirectDeployment())
        .reduce((traces: Trace[], root) => [...traces, ...root.traces], [])
        .find(trace => stage ? trace.isDeployment() == stage.stageName : !!trace.isDeployment());
    else
      return null;
  }

  getLatestSuccessfulArtifact(service: Service, stage?: Stage): Trace {
    let currentService = this.getService(service.serviceName);

    if(currentService.roots)
      return currentService.roots
        .filter(root => (root.isEvaluation() || root.isDeployment()) && (!stage || root.isFaulty() != stage.stageName || root.isDeployment() && root.getDeploymentDetails(stage)?.isDirectDeployment()))
        .reduce((traces: Trace[], root) => [...traces, ...root.traces], [])
        .find(trace => stage ? trace.isDeployment() == stage.stageName || trace.isEvaluation() == stage.stageName : !!trace.isDeployment() || !!trace.isEvaluation());
    else
      return null;
  }

  getLatestArtifact(service: Service, stage?: Stage): Trace {
    let currentService = this.getService(service.serviceName);

    if(currentService.roots)
      return currentService.roots
        .filter(root => root.isEvaluation() || root.isDeployment())
        .reduce((traces: Trace[], root) => [...traces, ...root.traces], [])
        .find(trace => stage ? trace.isDeployment() == stage.stageName || trace.isEvaluation() == stage.stageName : !!trace.isDeployment() || !!trace.isEvaluation());
    else
      return null;
  }

  getLatestFailedRootEvents(stage: Stage): Root[] {
    return this.getServices().map(service => service.roots?.find(root => (root?.isDeployment() || root?.isEvaluation()) && root.traces.some(trace => trace.data.stage === stage.stageName))).filter(root => root?.isFailedEvaluation() === stage.stageName);
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
    return root?.traces.slice().reverse().find(t => t.isEvaluation() && t.isFinished())?.getFinishedEvent();
  }

  getDeploymentsOfService(serviceName: string): Deployment[] {
    const deployments: Deployment[] = [];
    this.stages.forEach(stage => {
      const service = stage.services.find(service => service.serviceName === serviceName);
      if (service) {
        const image = service.getImageVersion();
        const deployment = deployments.find(deployment => deployment.version === image);
        if (deployment) {
          deployment.stages.push(stage.stageName);
        } else {
          const deployment = Object.assign(new Deployment(), {
            version: image,
            service: service.serviceName,
            stages: [stage.stageName]
          });

          deployments.push(deployment);
        }
      }
    });
    return deployments.sort((a,b) => a.version && b.version && semver.gt(a.version, b.version) ? -1 : 1);
  }

  static fromJSON(data: any) {
    return Object.assign(new this, data);
  }
}
