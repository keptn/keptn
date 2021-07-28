import { Deployment } from './deployment';
import { EventTypes } from './event-types';
import { Sequence } from './sequence';
import { Trace } from './trace';
import { Approval } from './approval';

export type DeploymentInformation = { deploymentUrl?: string, image?: string };

export class Service {
  serviceName!: string;
  deployedImage?: string;
  stage!: string;
  allDeploymentsLoaded = false;
  deployments: Deployment[] = [];
  lastEventTypes?: {[key: string]: {eventId: string, keptnContext: string, time: number}};
  sequences: Sequence[] = [];
  openApprovals: Approval[] = [];
  openRemediations: Sequence[] = [];
  latestSequence?: Sequence;
  deploymentTrace?: Trace;
  deploymentInformation?: DeploymentInformation;

  static fromJSON(data: unknown): Service {
    const service = Object.assign(new this(), data);
    if (service.latestSequence) {
      service.latestSequence = Sequence.fromJSON(service.latestSequence);
    }
    if (service.deploymentTrace) {
      service.deploymentTrace = Trace.fromJSON(service.deploymentTrace);
    }
    service.openRemediations = service.openRemediations?.map(remediation => Sequence.fromJSON(remediation)) ?? [];
    service.openApprovals = service.openApprovals.map(approval => {
      approval.trace = Trace.fromJSON(approval.trace);
      if (approval.evaluationTrace) {
        approval.evaluationTrace = Trace.fromJSON(approval.evaluationTrace);
      }
      return approval;
    });
    return service;
  }

  get deploymentContext(): string | undefined {
    return this.lastEventTypes?.[EventTypes.DEPLOYMENT_FINISHED]?.keptnContext ?? this.evaluationContext;
  }

  get deploymentTime(): number | undefined {
    return this.lastEventTypes?.[EventTypes.DEPLOYMENT_FINISHED]?.time || this.lastEventTypes?.[EventTypes.EVALUATION_FINISHED]?.time;
  }

  get evaluationContext(): string | undefined {
    return this.lastEventTypes?.[EventTypes.EVALUATION_FINISHED]?.keptnContext;
  }
  public getShortImageName(): string | undefined {
    return this.deployedImage?.split('/').pop()?.split(':').find(() => true);
  }

  getImageName(): string | undefined {
    return this.deployedImage?.split('/').pop();
  }

  getImageVersion(): string | undefined {
    return this.deployedImage?.split(':').pop();
  }

  getOpenApprovals(): Approval[] {
    return this.openApprovals;
  }

  public hasRemediations(): boolean {
    return this.deployments.some(d => d.stages.some(s => s.remediations.length !== 0));
  }
}
