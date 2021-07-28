import { Service } from '../models/service.js';

export interface StageResult {
  services: Service[];
  stageName: string;
  parentStages?: string[];
}
