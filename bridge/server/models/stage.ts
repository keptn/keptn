import { Service } from './service';
import { Stage as sg } from '../../shared/models/stage';

export class Stage extends sg {
  services: Service[] = [];

  public static fromJSON(data: unknown) {
    const stage = Object.assign(new this(), data);
    stage.services = stage.services.map((s) => {
      return Service.fromJSON(s);
    });
    return stage;
  }
}
