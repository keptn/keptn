import { IServiceRemediationInformation } from '../_interfaces/service-remediation-information';
import { Sequence } from './sequence';

export class ServiceRemediationInformation implements IServiceRemediationInformation {
  stages!: { name: string; remediations: Sequence[] }[];

  public static fromJSON(data: IServiceRemediationInformation): ServiceRemediationInformation {
    const serviceRemediationInformation = Object.assign(new this(), data);
    serviceRemediationInformation.stages = serviceRemediationInformation.stages.map((st) => {
      st.remediations = st.remediations.map((rem) => Sequence.fromJSON(rem));
      return st;
    });
    return serviceRemediationInformation;
  }
}
