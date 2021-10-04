import { Sequence } from './sequence';
import { Approval } from '../../shared/interfaces/approval';
import { Service as sv } from '../../shared/models/service';
import { Remediation } from './remediation';

type ServiceEvent = { eventId: string; keptnContext: string; time: number };
export type DeploymentInformation = { deploymentUrl?: string, image?: string };

export class Service extends sv {
  lastEventTypes: { [p: string]: ServiceEvent } = {};
  latestSequence?: Sequence;
  openRemediations: Remediation[] = [];
  openApprovals: Approval[] = [];
  deploymentInformation?: DeploymentInformation;

  public static fromJSON(data: unknown): Service {
    return Object.assign(new this(), data);
  }

  public getLatestSequence(): string | undefined {
    let latestSequence: ServiceEvent | undefined;
    for (const key of Object.keys(this.lastEventTypes)) {
      if(!latestSequence || this.lastEventTypes[key].time > latestSequence.time) {
        latestSequence = this.lastEventTypes[key];
      }
    }
    return latestSequence?.keptnContext;
  }
}
