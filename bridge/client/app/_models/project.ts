import {Stage} from "./stage";
import {Service} from "./service";
import {Trace} from "./trace";
import {Root} from "./root";

export class Project {
  projectName: string;
  gitUser: string;
  gitRemoteURI: string;
  gitToken: string;

  stages: Stage[];
  services: Service[];

  getServices(): Service[] {
    if(!this.services) {
      this.services = [];
      this.stages.forEach((stage: Stage) => {
        this.services = this.services.concat(stage.services.filter(s => !this.services.some(ss => ss.serviceName == s.serviceName)));
      });
    }
    return this.services;
  }

  getStages(): Stage[] {
    return this.stages;
  }

  getLatestDeployment(service: Service, stage?: Stage): Trace {
    let currentService = this.getServices().find(s => s.serviceName == service.serviceName);

    if(currentService.roots)
      return currentService.roots
        .reduce((traces: Trace[], root) => [...traces, ...root.traces], [])
        .find(trace => trace.type == 'sh.keptn.events.deployment-finished' && (!stage || (trace.data.stage == stage.stageName && currentService.roots.find(r => r.shkeptncontext == trace.shkeptncontext).isFaulty() != stage.stageName)));
    else
      return null;
  }

  getDeploymentEvaluation(trace: Trace): Trace {
    let currentService = this.getServices().find(s => s.serviceName == trace.data.service);
    if(currentService.roots) {
      return currentService.roots
        .reduce((traces: Trace[], root) => [...traces, ...root.traces], [])
        .find(t => t.type == 'sh.keptn.events.evaluation-done' && t.shkeptncontext == trace.shkeptncontext);
    }
    return null;
  }

  static fromJSON(data: any) {
    return Object.assign(new this, data);
  }
}
