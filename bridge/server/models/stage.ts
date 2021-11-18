import { Service } from './service';
import { Stage as sg } from '../../shared/models/stage';
import { IStage } from '../../shared/interfaces/stage';

export class Stage extends sg {
  services: Service[] = [];

  public static fromJSON(data: IStage): Stage {
    const stage = Object.assign(new this(), data);
    stage.services = stage.services.map((s) => Service.fromJSON(s));
    return stage;
  }

  public static getAllServices(stages: Stage[]): Service[] {
    return stages.reduce((services: Service[], stage) => [...services, ...stage.services], []);
  }
}
