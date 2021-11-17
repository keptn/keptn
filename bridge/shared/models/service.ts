import { Sequence } from './sequence';
import { Approval } from '../interfaces/approval';
import { IService, IServiceEvent } from '../interfaces/service';

export type DeploymentInformation = { deploymentUrl?: string; image?: string };

export class Service implements IService {
  serviceName!: string;
  creationDate!: number;
  stage!: string;
  deployedImage?: string;
  lastEventTypes: { [p: string]: IServiceEvent } = {};
  latestSequence?: Sequence;
  openRemediations: Sequence[] = [];
  openApprovals?: Approval[] = [];
  deploymentInformation?: DeploymentInformation;

  public static fromJSON(data: unknown): Service {
    return Object.assign(new this(), data);
  }

  public getLatestEvent(): IServiceEvent | undefined {
    let latestSequence: IServiceEvent | undefined;
    for (const key of Object.keys(this.lastEventTypes)) {
      if (!latestSequence || this.lastEventTypes[key].time > latestSequence.time) {
        latestSequence = this.lastEventTypes[key];
      }
    }
    return latestSequence;
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
