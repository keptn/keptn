import {Stage} from "./stage";
import {Service} from "./service";

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

  static fromJSON(data: any) {
    return Object.assign(new this, data);
  }
}
