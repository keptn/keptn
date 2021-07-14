import {Service} from './service';

export class Stage {
  stageName: string;
  parentStages: string[];
  services: Service[];

  static fromJSON(data: any) {
    return Object.assign(new this(), data);
  }

  public servicesWithOpenApprovals(): Service[] {
    return this.services.filter(s => s.getOpenApprovals().length > 0);
  }

  public getOpenProblems() {
    return this.services.reduce((openProblems, service) => [...openProblems, ...service.getOpenProblems()], []);
  }
}
