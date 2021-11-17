import { Service } from './service';
import { Stage as sg } from '../../shared/models/stage';
import { IStage } from '../../shared/interfaces/stage';
import { EventTypes } from '../../shared/interfaces/event-types';

export class Stage extends sg {
  services: Service[] = [];

  public static fromJSON(data: IStage): Stage {
    const stage = Object.assign(new this(), data);
    stage.services = stage.services.map((s) => Service.fromJSON(s));
    return stage;
  }

  public static hasApprovalUpdate(stages: Stage[], fromTime: Date): boolean {
    return (
      this.hasUpdate(stages, EventTypes.APPROVAL_STARTED, fromTime) ||
      this.hasUpdate(stages, EventTypes.APPROVAL_FINISHED, fromTime)
    );
  }

  public static hasRemediationUpdate(stages: Stage[], fromTime: Date): boolean {
    return (
      this.hasUpdate(stages, EventTypes.REMEDIATION_TRIGGERED, fromTime) ||
      this.hasUpdate(stages, EventTypes.REMEDIATION_FINISHED, fromTime)
    );
  }

  public static hasEvaluationsUpdate(stages: Stage[], fromTime: Date): boolean {
    return this.hasUpdate(stages, EventTypes.EVALUATION_FINISHED, fromTime);
  }

  private static hasUpdate(stages: Stage[], eventType: EventTypes, fromTime: Date): boolean {
    for (const stage of stages) {
      for (const service of stage.services) {
        if (service.hasUpdate(eventType, fromTime)) {
          return true;
        }
      }
    }
    return false;
  }

  public static getAllServices(stages: Stage[]): Service[] {
    return stages.reduce((services: Service[], stage) => [...services, ...stage.services], []);
  }
}
