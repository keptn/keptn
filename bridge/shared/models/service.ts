import { Sequence } from './sequence';
import { Approval } from '../interfaces/approval';
import { IService, IServiceEvent } from '../interfaces/service';

export type DeploymentInformation = { deploymentUrl?: string; image?: string };

export class Service implements IService {
  serviceName!: string;
  creationDate!: string;
  deployedImage?: string;
  lastEventTypes: { [event: string]: IServiceEvent | undefined } = {};
  latestSequence?: Sequence;
  openRemediations: Sequence[] = [];
  openApprovals?: Approval[] = [];
  deploymentInformation?: DeploymentInformation;

  public static fromJSON(data: IService): Service {
    return Object.assign(new this(), data);
  }

  public getLatestEvent(): IServiceEvent | undefined {
    let latestSequence: IServiceEvent | undefined;
    for (const key of Object.keys(this.lastEventTypes)) {
      const event = this.lastEventTypes[key];
      if (!latestSequence || (event && +event.time > +latestSequence.time)) {
        latestSequence = event;
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
