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

  public getLatestSequence(stageName: string): string | undefined {
    const sequenceEvents = this.getSequenceEvents(stageName);
    return sequenceEvents.reduce((latestSequence: ServiceEvent | undefined, currentSequence: ServiceEvent) => {
      return latestSequence && latestSequence.time > currentSequence.time ? latestSequence : currentSequence;
    }, undefined)?.keptnContext;
  }

  private getSequenceEvents(stageName: string): ServiceEvent[] {
    const sequenceEvents: ServiceEvent[] = [];
    for (const key of Object.keys(this.lastEventTypes)) {
      if (this.isSequence(key, stageName)) {
        sequenceEvents.push(this.lastEventTypes[key]);
      }
    }
    return sequenceEvents;
  }

  private isSequence(eventType: string, stageName: string): boolean {
    return eventType.split('.').length === 6 && eventType.includes(stageName);
  }
}
