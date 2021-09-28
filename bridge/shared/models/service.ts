import { Sequence } from './sequence';
import { Approval } from '../interfaces/approval';

type ServiceEvent = { eventId: string; keptnContext: string; time: number };
export type DeploymentInformation = { deploymentUrl?: string; image?: string };

export class Service {
  serviceName!: string;
  creationDate!: number;
  stage!: string;
  deployedImage?: string;
  lastEventTypes: { [p: string]: ServiceEvent } = {};
  latestSequence?: Sequence;
  openRemediations: Sequence[] = [];
  openApprovals?: Approval[] = [];
  deploymentInformation?: DeploymentInformation;

  public static fromJSON(data: unknown): Service {
    return Object.assign(new this(), data);
  }
}
