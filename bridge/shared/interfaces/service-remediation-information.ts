import { ISequenceState } from './sequence';

export interface ServiceRemediationInformation {
  stages: {
    name: string;
    remediations: ISequenceState[];
  }[];
}
