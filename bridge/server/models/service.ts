import { Sequence } from './sequence';
import { Approval } from '../../shared/interfaces/approval';
import { Service as sv } from '../../shared/models/service';
import { Remediation } from './remediation';

type ServiceEvent = { eventId: string; keptnContext: string; time: number };
export type DeploymentInformation = { deploymentUrl?: string; image?: string };

export class Service extends sv {
  lastEventTypes: { [p: string]: ServiceEvent } = {};
  latestSequence?: Sequence;
  openRemediations: Remediation[] = [];
  openApprovals: Approval[] = [];
  deploymentInformation?: DeploymentInformation;

  public static fromJSON(data: unknown): Service {
    return Object.assign(new this(), data);
  }
}
