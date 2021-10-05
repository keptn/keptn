import { Deployment } from './deployment';
import { EventTypes } from '../../../shared/interfaces/event-types';
import { Sequence } from './sequence';
import { Trace } from './trace';
import { Service as sv } from '../../../shared/models/service';
import { Approval } from '../_interfaces/approval';
import { ResultTypes } from '../../../shared/models/result-types';

export type DeploymentInformation = { deploymentUrl?: string; image?: string };

export class Service extends sv {
  allDeploymentsLoaded = false;
  deployments: Deployment[] = [];
  sequences: Sequence[] = [];
  openApprovals: Approval[] = [];
  openRemediations: Sequence[] = [];
  latestSequence?: Sequence;

  static fromJSON(data: unknown): Service {
    const service = Object.assign(new this(), data);
    if (service.latestSequence) {
      service.latestSequence = Sequence.fromJSON(service.latestSequence);
    }

    // Support old (deprecated) format from API - openRemediation should be a Sequence but old format just provides an event
    // If openRemediations do not have stages, it is in the old format and should not be processed as Sequence
    const hasStages = service.openRemediations?.some((remediation) => remediation.stages);
    if (hasStages) {
      service.openRemediations = service.openRemediations?.map((remediation) => Sequence.fromJSON(remediation)) ?? [];
    } else {
      service.openRemediations = [];
    }
    service.openRemediations = service.openRemediations?.map((remediation) => Sequence.fromJSON(remediation)) ?? [];
    service.openApprovals = service.openApprovals.map((approval) => {
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
    return (
      this.lastEventTypes?.[EventTypes.DEPLOYMENT_FINISHED]?.time ||
      this.lastEventTypes?.[EventTypes.EVALUATION_FINISHED]?.time
    );
  }

  get evaluationContext(): string | undefined {
    return this.lastEventTypes?.[EventTypes.EVALUATION_FINISHED]?.keptnContext;
  }

  public getShortImageName(): string | undefined {
    return this.deployedImage
      ?.split('/')
      .pop()
      ?.split(':')
      .find(() => true);
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
    return this.openRemediations.length > 0;
  }

  public hasFailedEvaluation(): boolean {
    return this.latestSequence?.getEvaluation(this.stage)?.result === ResultTypes.FAILED;
  }

  public getFailedEvaluationSequence(): Sequence | undefined {
    return this.latestSequence?.getEvaluation(this.stage)?.result === ResultTypes.FAILED
      ? this.latestSequence
      : undefined;
  }
}
