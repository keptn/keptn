import { Service } from './service';
import { Sequence } from './sequence';
import { Approval } from './approval';

export class Stage {
  stageName!: string;
  parentStages?: string[];
  services: Service[] = [];

  static fromJSON(data: unknown): Stage {
    const stage = Object.assign(new this(), data);
    stage.services = stage.services.map(s => {
      s.stage = stage.stageName;
      return Service.fromJSON(s);
    });
    return stage;
  }

  public servicesWithOpenApprovals(): Service[] {
    return this.services.filter(s => s.getOpenApprovals().length > 0);
  }

  public getOpenProblems(): Sequence[] {
    return this.services.reduce((remediations: Sequence[], service: Service) => [...remediations, ...service.openRemediations], []);
  }

  public getOpenApprovals(): Approval[] {
    return this.services.reduce((openApprovals: Approval[], service: Service) => [...openApprovals, ...service.openApprovals], []);
  }
}
