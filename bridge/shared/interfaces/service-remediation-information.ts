import { ISequence } from './sequence';

export interface ServiceRemediationInformation {
  stages: {
    name: string;
    remediations: ISequence[];
    config?: string;
  }[];
}
