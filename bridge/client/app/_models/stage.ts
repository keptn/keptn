import {Service} from './service';
import { Root } from './root';

export class Stage {
  stageName!: string;
  parentStages?: string[];
  services: Service[] = [];

  static fromJSON(data: unknown) {
    return Object.assign(new this(), data);
  }

  public servicesWithOpenApprovals(): Service[] {
    return this.services.filter(s => s.getOpenApprovals().length > 0);
  }

  public getOpenProblems(): Root[] {
    return this.services.reduce((openProblems: Root[], service: Service) => [...openProblems, ...service.getOpenProblems()], []);
  }
}
