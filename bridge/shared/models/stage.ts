import { Service } from './service';

export class Stage {
  stageName!: string;
  parentStages?: string[];
  services: Service[] = [];
}
