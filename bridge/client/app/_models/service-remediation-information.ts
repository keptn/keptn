import { IServiceRemediationInformation } from '../_interfaces/service-remediation-information';
import { SequenceState } from './sequenceState';

export class ServiceRemediationInformation implements IServiceRemediationInformation {
  stages!: { name: string; remediations: SequenceState[] }[];

  public static fromJSON(data: IServiceRemediationInformation): ServiceRemediationInformation {
    const serviceRemediationInformation = Object.assign(new this(), data);
    serviceRemediationInformation.stages = serviceRemediationInformation.stages.map((st) => {
      st.remediations = st.remediations.map((rem) => SequenceState.fromJSON(rem));
      return st;
    });
    return serviceRemediationInformation;
  }
}
