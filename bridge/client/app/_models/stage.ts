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

  public update(newStage: Stage): void {
    this.parentStages = newStage.parentStages;
    for (const newService of newStage.services) {
      const existingService = this.services.find((service) => service.serviceName === newService.serviceName);
      if (existingService) {
        // update existing service
        Object.assign(existingService, newService);
      } else {
        // add new service
        this.services.push(newService);
      }
    }
    // remove deleted services
    for (let i = 0; i < this.services.length; ) {
      if (!newStage.services.some((service) => service.serviceName === this.services[i].serviceName)) {
        this.services.splice(i, 1);
        --i;
      }
      ++i;
    }
    this.services.sort(this.compareServices());
  }

  private compareServices() {
    return (a: Service, b: Service): number => {
      if (!a.latestSequence && !b.latestSequence) {
        return new Date(b.creationDate).getTime() - new Date(a.creationDate).getTime();
      } else if (!a.latestSequence) {
        return 1;
      } else if (!b.latestSequence) {
        return -1;
      } else {
        return new Date(b.latestSequence.time).getTime() - new Date(a.latestSequence.time).getTime();
      }
    };
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
