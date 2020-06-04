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
        .reduce((traces: Trace[], root: Root) => {
          return [...traces, ...root.traces];
        }, [])
        .find(trace => trace.type == 'sh.keptn.events.deployment-finished' && (!stage || (trace.data.stage == stage.stageName && currentService.roots.find(r => r.shkeptncontext == trace.shkeptncontext).isFaulty() != stage.stageName)));
    else
      return null;
  }

  getDeploymentsFromPrevStage(trace: Trace): Trace[] {
    let prevStage = this.getStages().find((s, i, stages) => {
      if(stages[i+1] && stages[i+1].stageName == trace.data.stage)
        return true;
    });
    if(prevStage) {
      let currentService = this.getServices().find(s => s.serviceName == trace.data.service);
      if(currentService.roots) {
        let traces = currentService.roots.reduce((traces: Trace[], root: Root) => { return [...traces, ...root.traces]; }, []).filter(t => t.type == 'sh.keptn.events.deployment-finished' && t.data.stage == prevStage.stageName && currentService.roots.find(r => r.shkeptncontext == t.shkeptncontext).isFaulty() != prevStage.stageName);
        let currentDeployment = traces.find(t => t.shkeptncontext == trace.shkeptncontext);
        return traces.slice(0, traces.indexOf(currentDeployment));
      }
    }
    return [];
  }

  getDeploymentEvaluation(trace: Trace): Trace {
    let currentService = this.getServices().find(s => s.serviceName == trace.data.service);
    if(currentService.roots) {
      return currentService.roots
        .reduce((traces: Trace[], root: Root) => {
          return [...traces, ...root.traces];
        }, [])
        .find(t => t.type == 'sh.keptn.events.evaluation-done' && t.shkeptncontext == trace.shkeptncontext);
    }
    return null;
  }

  getDeployedServices(stage: Stage) {
    return stage.services.filter(service => !!this.getLatestDeployment(service, stage));
  }

  getServicesOutOfSync(stage: Stage) {
    let services = this.getDeployedServices(stage);
    let servicesOutOfSync: Service[] = [];
    services.forEach(service => {
      let deployment = this.getLatestDeployment(service, stage);
      if(this.getDeploymentsFromPrevStage(deployment).length > 0)
        servicesOutOfSync.push(service);
    });
    return servicesOutOfSync;
  }

  static fromJSON(data: any) {
    return Object.assign(new this, data);
  }
}
