import { Service } from './service';
import { IStage } from '../interfaces/stage';

export class Stage implements IStage {
  stageName!: string;
  parentStages?: string[];
  services: Service[] = [];
}
