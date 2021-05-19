import {Service} from "./service";

export class Stage {
  stageName: string;
  parentStages: string[];
  services: Service[];

  public getOpenApprovals() {
    return this.services.reduce((openApprovals, service) => [...openApprovals, ...service.getOpenApprovals()], []);
  }

  public servicesWithOpenApprovals() {
    return this.services.filter(s => s.getOpenApprovals().length > 0);
  }

  public getOpenProblems() {
    return this.services.reduce((openProblems, service) => [...openProblems, ...service.getOpenProblems()], []);
  }

  static fromJSON(data: any) {
    return Object.assign(new this, data);
  }
}
