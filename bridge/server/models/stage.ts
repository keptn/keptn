import { Service } from './service';
import { Stage as sg } from '../../shared/models/stage';

export class Stage extends sg {
  services: Service[] = [];

  public static fromJSON(data: unknown): Stage {
    const stage = Object.assign(new this(), data);
    stage.services = stage.services.map((s) => Service.fromJSON(s));
    return stage;
  }
}
