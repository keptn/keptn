import { Service } from './service';
import { Stage as st } from '../../../shared/models/stage';
import { ResultTypes } from '../../../shared/models/result-types';

export class Stage extends st {
  services: Service[] = [];

  static fromJSON(data: unknown): Stage {
    const stage = Object.assign(new this(), data);
    stage.services = stage.services.map((s) => {
      s.stage = stage.stageName;
      return Service.fromJSON(s);
    });
    return stage;
  }

  public getServicesWithOpenApprovals(): Service[] {
    return this.services.filter((s) => s.getOpenApprovals().length > 0);
  }

  public getServicesWithFailedEvaluation(): Service[] {
    return this.services.filter(
      (service) => service.latestSequence?.getEvaluation(this.stageName)?.result === ResultTypes.FAILED
    );
  }

  public getServicesWithRemediations(): Service[] {
    return this.services.filter((service) => service.openRemediations.length > 0);
  }
}
