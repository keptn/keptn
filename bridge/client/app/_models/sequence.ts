import { Trace } from './trace';
import { EventTypes } from '../../../shared/interfaces/event-types';
import { EvaluationResult } from '../../../shared/interfaces/evaluation-result';
import { EVENT_ICONS } from './event-icons';
import { RemediationAction } from '../../../shared/models/remediation-action';
import { Sequence as sq, SequenceEvent, SequenceStage, SequenceState } from '../../../shared/models/sequence';
import { DtIconType } from '@dynatrace/barista-icons';
import { ResultTypes } from '../../../shared/models/result-types';

type SeqStage = SequenceStage & {
  latestEvaluationTrace?: Trace;
  actions?: RemediationAction[];
};

export class Sequence extends sq {
  stages!: SeqStage[];
  problemTitle?: string;
  traces: Trace[] = [];

  public static fromJSON(data: unknown): Sequence {
    const sequence = Object.assign(new this(), data);
    for (const stage of sequence.stages) {
      stage.actions = stage.actions?.map((s) => RemediationAction.fromJSON(s)) ?? [];
      if (stage.latestEvaluationTrace) {
        stage.latestEvaluationTrace = Trace.fromJSON(stage.latestEvaluationTrace);
      }
    }
    return sequence;
  }

  public static isFinished(state: SequenceState): boolean {
    return (
      state === SequenceState.FINISHED ||
      state === SequenceState.TIMEDOUT ||
      state === SequenceState.ABORTED ||
      state === SequenceState.SUCCEEDED
    );
  }

  public static getShortType(type: string): string {
    const parts = type?.split('.');
    if (parts.length === 6) {
      return parts[4];
    } else if (parts.length === 5) {
      return parts[3];
    } else {
      return type;
    }
  }

  public getStage(stageName: string): SeqStage | undefined {
    return this.stages.find((stage) => stage.name === stageName);
  }

  public getStages(): string[] {
    return this.stages.map((stage) => stage.name);
  }

  public getLastStage(): string | undefined {
    return this.stages[this.stages.length - 1]?.name;
  }

  public isFaulty(stageName?: string): boolean {
    return (
      (stageName
        ? !!this.getStage(stageName)?.latestFailedEvent
        : this.stages.some((stage) => stage.latestFailedEvent)) || this.isTimedOut(stageName)
    );
  }

  public isFinished(stageName?: string): boolean {
    const state = stageName ? this.getStage(stageName)?.state : this.state;
    return !!state && Sequence.isFinished(state);
  }

  public getEvaluation(stage: string): EvaluationResult | undefined {
    return this.getStage(stage)?.latestEvaluation;
  }

  public getEvaluationTrace(stage: string): Trace | undefined {
    return this.getStage(stage)?.latestEvaluationTrace;
  }

  public hasPendingApproval(stageName?: string): boolean {
    return stageName
      ? this.getStage(stageName)?.latestEvent?.type === EventTypes.APPROVAL_TRIGGERED ||
          this.getStage(stageName)?.latestEvent?.type === EventTypes.APPROVAL_STARTED
      : this.stages.some(
          (stage) =>
            stage.latestEvent?.type === EventTypes.APPROVAL_TRIGGERED ||
            stage.latestEvent?.type === EventTypes.APPROVAL_STARTED
        );
  }

  public getStatus(): string {
    let status: string = this.state;
    if (this.state === SequenceState.FINISHED) {
      if (this.stages.some((stage) => stage.latestFailedEvent)) {
        status = 'failed';
      } else {
        status = 'succeeded';
      }
    } else if (this.isAborted()) {
      status = 'aborted';
    } else if (this.isWaiting()) {
      status = 'waiting';
    } else if (this.isTimedOut()) {
      status = 'timed out';
    }
    return status;
  }

  public isLoading(stageName?: string): boolean {
    const isStarted = this.state === SequenceState.TRIGGERED || this.state === SequenceState.STARTED;
    return isStarted && (!stageName || !this.isFinished(stageName));
  }

  public isSuccessful(stageName?: string): boolean {
    return (
      this.isSucceeded(stageName) ||
      (!this.isFaulty(stageName) &&
        !this.isWarning(stageName) &&
        !this.isAborted(stageName) &&
        this.isFinished(stageName))
    );
  }

  private isSucceeded(stageName?: string): boolean {
    return (stageName ? this.getStage(stageName)?.state : this.state) === SequenceState.SUCCEEDED;
  }

  public isWarning(stageName?: string): boolean {
    return (
      !this.isFaulty(stageName) &&
      (stageName
        ? this.getStage(stageName)?.latestEvaluation?.result === ResultTypes.WARNING
        : this.stages.some((st) => st.latestEvaluation?.result === ResultTypes.WARNING))
    );
  }

  public isWaiting(): boolean {
    if (this.isRemediation()) {
      return false;
    }

    const lastStageName = this.getLastStage();
    if (lastStageName && this.state === SequenceState.STARTED) {
      const lastStage = this.getStage(lastStageName);
      return (
        lastStage?.state === SequenceState.FINISHED || // last stages is finished, but sequence is still started, means it is waiting for next stage to be triggered
        (lastStage?.state === SequenceState.TRIGGERED && !!lastStage?.latestEvent?.type.endsWith('.triggered'))
      ); // last stage is triggered, but has no running tasks
    } else {
      // no stages yet, sequence is triggered, so waiting
      return this.state === SequenceState.TRIGGERED;
    }
  }

  public isRemediation(): boolean {
    return this.name === 'remediation';
  }

  public isPaused(): boolean {
    return this.state === SequenceState.PAUSED;
  }

  public isUnknownState(): boolean {
    return this.state === SequenceState.UNKNOWN;
  }

  public isAborted(stageName?: string): boolean {
    return (stageName ? this.getStage(stageName)?.state : this.state) === SequenceState.ABORTED;
  }

  public isTimedOut(stageName?: string): boolean {
    return (stageName ? this.getStage(stageName)?.state : this.state) === SequenceState.TIMEDOUT;
  }

  public getLatestEvent(): SequenceEvent | undefined {
    return this.stages[this.stages.length - 1]?.latestEvent;
  }

  public getIcon(stageName?: string): DtIconType {
    const stage = stageName ? this.getStage(stageName) : this.stages[this.stages.length - 1];
    return stage?.latestEvent?.type
      ? EVENT_ICONS[Sequence.getShortType(stage?.latestEvent?.type)] || EVENT_ICONS.default
      : EVENT_ICONS.default;
  }

  public getShortImageName(): string | undefined {
    return this.stages[0]?.image?.split('/').pop();
  }

  public getTraces(stageName: string): Trace[] {
    return this.traces.filter((trace) => trace.stage === stageName);
  }

  public findTrace(comp: (args: Trace) => boolean): Trace | undefined {
    return this.traces.reduce((result: Trace | undefined, trace: Trace) => result || trace.findTrace(comp), undefined);
  }

  public findLastTrace(comp: (args: Trace) => boolean): Trace | undefined {
    return this.traces.reduce((result: Trace | undefined, trace: Trace) => trace.findTrace(comp) || result, undefined);
  }

  public getLabels(): Map<string, string> | undefined {
    return this.getLastTrace()?.getFinishedEvent()?.labels || this.getFirstTrace()?.labels;
  }

  private getLastTrace(): Trace {
    return this.traces[this.traces.length - 1];
  }

  private getFirstTrace(): Trace | undefined {
    return this.traces[0];
  }

  public getRemediationActions(): RemediationAction[] {
    return this.stages[0]?.actions ?? [];
  }

  public setState(state: SequenceState): void {
    this.state = state;
  }
}
