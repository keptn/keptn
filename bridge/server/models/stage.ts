import { StageResult } from '../interfaces/stage-result.js';
import { Service } from './service.js';

export class Stage implements StageResult {
  services: Service[] = [];
  stageName!: string;
  parentStages?: string[];

  public static fromJSON(data: unknown) {
    const stage = Object.assign(new this(), data);
    stage.services = stage.services.map(s => {
      return Service.fromJSON(s);
    });
    return stage;
  }
}
