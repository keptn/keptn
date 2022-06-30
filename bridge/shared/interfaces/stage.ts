import { IService } from './service';

export interface IStage {
  services: IService[];
  stageName: string;
  parentStages?: string[];
}
