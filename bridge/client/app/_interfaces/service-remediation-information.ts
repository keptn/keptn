import { ServiceRemediationInformation as sri } from '../../../shared/interfaces/service-remediation-information';
import { Sequence } from '../_models/sequence';

export interface IServiceRemediationInformation extends sri {
  stages: {
    name: string;
    remediations: Sequence[];
  }[];
}
