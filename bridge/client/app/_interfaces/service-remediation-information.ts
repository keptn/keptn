import { ServiceRemediationInformation as sri } from '../../../shared/interfaces/service-remediation-information';
import { SequenceState } from '../_models/sequenceState';

export interface IServiceRemediationInformation extends sri {
  stages: {
    name: string;
    remediations: SequenceState[];
  }[];
}
