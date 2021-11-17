import { Sequence } from './sequence';
import { Approval } from '../../shared/interfaces/approval';
import { Service as sv } from '../../shared/models/service';
import { Remediation } from '../../shared/models/remediation';
import { EventTypes } from '../../shared/interfaces/event-types';

type ServiceEvent = { eventId: string; keptnContext: string; time: string };
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

  public get latestDeploymentEvent(): ServiceEvent | undefined {
    return this.deploymentEvent ?? this.evaluationEvent;
  }

  public get deploymentEvent(): ServiceEvent | undefined {
    return this.lastEventTypes?.[EventTypes.DEPLOYMENT_FINISHED];
  }

  public get evaluationEvent(): ServiceEvent | undefined {
    return this.lastEventTypes?.[EventTypes.EVALUATION_FINISHED];
  }

  public hasEvaluationUpdate(fromTime: Date): boolean {
    return this.hasUpdate(EventTypes.EVALUATION_FINISHED, fromTime);
  }

  public hasUpdate(type: EventTypes, fromTime: Date): boolean {
    const serviceEvent = this.lastEventTypes?.[type];
    return !!serviceEvent && new Date(+serviceEvent.time / 1_000_000) > fromTime;
  }

  public hasRemediationUpdate(fromTime: Date): boolean {
    return (
      this.hasUpdate(EventTypes.REMEDIATION_TRIGGERED, fromTime) ||
      this.hasUpdate(EventTypes.REMEDIATION_FINISHED, fromTime) ||
      this.hasUpdate(EventTypes.ACTION_TRIGGERED, fromTime)
    );
  }

  public hasApprovalUpdate(fromTime: Date): boolean {
    return (
      this.hasUpdate(EventTypes.APPROVAL_STARTED, fromTime) || this.hasUpdate(EventTypes.APPROVAL_FINISHED, fromTime)
    );
  }
}
