import { Sequence } from '../models/sequence';

export interface ServiceRemediationInformation {
  stages: {
    name: string;
    remediations: Sequence[];
  }[];
}
