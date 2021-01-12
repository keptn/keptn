import {Stage} from "./stage";
import {Service} from "./service";
import {Trace} from "./trace";
import {EventTypes} from "./event-types";
import {Root} from "./root";

export class Project {
  projectName: string;
  gitUser: string;
  gitRemoteURI: string;
  gitToken: string;

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
    return this.getServices().map(service => service.roots?.find(root => root.traces.some(trace => trace.data.stage === stage.stageName))).filter(root => root && root?.isProblem() && !root?.isProblemResolvedOrClosed());
  }

  getRootEvent(service: Service, event: Trace): Root {
    return service.roots?.find(root => root.shkeptncontext == event.shkeptncontext);
  }

  getDeploymentEvaluation(trace: Trace): Trace {
    let service = this.getServices().find(s => s.serviceName == trace.data.service);
    let root = this.getRootEvent(service, trace);
    return root?.traces.slice().reverse().find(t => t.type == EventTypes.EVALUATION_FINISHED);
  }

  static fromJSON(data: any) {
    return Object.assign(new this, data);
  }
}
