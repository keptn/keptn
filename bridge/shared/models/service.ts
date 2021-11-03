import { Sequence } from './sequence';
import { Approval } from '../interfaces/approval';

interface ServiceEvent {
  eventId: string;
  keptnContext: string;
  time: string; // nanoseconds
}
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

  public getLatestSequence(): string | undefined {
    let latestSequence: ServiceEvent | undefined;
    for (const key of Object.keys(this.lastEventTypes)) {
      if (!latestSequence || this.lastEventTypes[key].time > latestSequence.time) {
        latestSequence = this.lastEventTypes[key];
      }
    }
    return latestSequence?.keptnContext;
  }

  public getImageVersion(): string | undefined {
    return this.deployedImage?.split(':').pop();
  }

  public getShortImageName(): string | undefined {
    return this.getShortImage()
      ?.split(':')
      .find(() => true);
  }

  public getShortImage(): string | undefined {
    return this.deployedImage?.split('/').pop();
  }
}
