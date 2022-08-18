import { EventTypes } from '../../../shared/interfaces/event-types';
import { SequenceState } from './sequenceState';
import { Trace } from './trace';
import { Service as sv } from '../../../shared/models/service';
import { ResultTypes } from '../../../shared/models/result-types';

export type DeploymentInformation = { deploymentUrl?: string; image?: string };

export class Service extends sv {
  stage!: string;
  openApprovals: Trace[] = [];
  openRemediations: SequenceState[] = [];
  latestSequence?: SequenceState;

  static fromJSON(data: unknown): Service {
    const service = Object.assign(new this(), data);
    if (service.latestSequence) {
      service.latestSequence = SequenceState.fromJSON(service.latestSequence);
    }

    // Support old (deprecated) format from API - openRemediation should be a Sequence but old format just provides an event
    // If openRemediations do not have stages, it is in the old format and should not be processed as Sequence
    const hasStages = service.openRemediations?.some((remediation) => remediation.stages);
    if (hasStages) {
      service.openRemediations =
        service.openRemediations?.map((remediation) => SequenceState.fromJSON(remediation)) ?? [];
    } else {
      service.openRemediations = [];
    }
    service.openRemediations =
      service.openRemediations?.map((remediation) => SequenceState.fromJSON(remediation)) ?? [];
    service.openApprovals = service.openApprovals.map((approval) => Trace.fromJSON(approval));
    return service;
  }

  get deploymentContext(): string | undefined {
    return this.lastEventTypes?.[EventTypes.DEPLOYMENT_FINISHED]?.keptnContext ?? this.evaluationContext;
  }

  get evaluationContext(): string | undefined {
    return this.lastEventTypes?.[EventTypes.EVALUATION_FINISHED]?.keptnContext;
  }

  getImageName(): string | undefined {
    return this.deployedImage?.split('/').pop();
  }

  getOpenApprovals(): Trace[] {
    return this.openApprovals;
  }

  public hasRemediations(): boolean {
    return this.openRemediations.length > 0;
  }

  public hasFailedEvaluation(): boolean {
    return this.latestSequence?.getEvaluation(this.stage)?.result === ResultTypes.FAILED;
  }

  public getFailedEvaluationSequence(): SequenceState | undefined {
    return this.latestSequence?.getEvaluation(this.stage)?.result === ResultTypes.FAILED
      ? this.latestSequence
      : undefined;
  }
}
