import { Trace } from './trace';
import { EventTypes } from '../../../shared/interfaces/event-types';
import { EvaluationResult } from '../../../shared/interfaces/evaluation-result';
import { EVENT_ICONS } from './event-icons';
import { RemediationAction } from '../../../shared/models/remediation-action';
import { ISequence, SequenceEvent, SequenceStage, SequenceState } from '../../../shared/interfaces/sequence';
import { DtIconType } from '@dynatrace/barista-icons';
import { ResultTypes } from '../../../shared/models/result-types';

const pendingApprovalTypes = [EventTypes.APPROVAL_TRIGGERED, EventTypes.APPROVAL_STARTED];
const finishedStates = [SequenceState.FINISHED, SequenceState.TIMEDOUT, SequenceState.ABORTED, SequenceState.SUCCEEDED];
const startedStates = [SequenceState.TRIGGERED, SequenceState.STARTED];

type SeqStage = SequenceStage & {
  latestEvaluationTrace?: Trace;
  actions?: RemediationAction[];
};

export interface SequenceStateInfo {
  finished: boolean;
  loading: boolean;
  waiting: boolean;
  pendingApproval: boolean;
  successful: boolean;
  warning: boolean;
  faulty: boolean;
  aborted: boolean;
  timedOut: boolean;
  steady: boolean;
  icon: DtIconType;
  statusText: string;
  evaluation: EvaluationResult | undefined;
}

function getStatusText(
  status: string,
  state: { successful: boolean; faulty: boolean; aborted: boolean; waiting: boolean; timedOut: boolean }
): string {
  if (state.successful) {
    status = 'succeeded';
  } else if (state.timedOut) {
    status = 'timed out';
  } else if (state.faulty) {
    status = 'failed';
  } else if (state.aborted) {
    status = 'aborted';
  } else if (state.waiting) {
    status = 'waiting';
  }
  return status;
}

export function isSequenceStarted(state: SequenceState): boolean {
  return startedStates.includes(state);
}

export function isSequenceLoading(sequenceState: SequenceState, stageState?: SequenceState): boolean {
  return isSequenceStarted(sequenceState) && (!stageState || !isFinished(stageState));
}

export function isSequenceAborted(state: SequenceState): boolean {
  return state === SequenceState.ABORTED;
}

export function createSequenceStateInfo(sequence: ISequence, stageName?: string): SequenceStateInfo {
  const stage = getStage(sequence, stageName);
  const latestStage = stage ? stage : sequence.stages[sequence.stages.length - 1];
  const stages = stage ? [stage] : sequence.stages;
  const state = stageName ? stage?.state : sequence.state;
  const finished = !!state && isFinished(state);
  const loading = isSequenceLoading(sequence.state, stage?.state);
  const waiting = state === SequenceState.WAITING;
  const aborted = !!state && isSequenceAborted(state);
  const timedOut = state === SequenceState.TIMEDOUT;
  const faulty = stages.some((s) => s.latestFailedEvent) || timedOut;
  const warning = !faulty && stages.some((s) => s.latestEvaluation?.result === ResultTypes.WARNING);
  const successful = finished && !faulty && !warning && !aborted;
  const pendingApproval = stages.some((s) => pendingApprovalTypes.includes(s.latestEvent?.type as EventTypes));
  const icon = latestStage?.latestEvent?.type
    ? EVENT_ICONS[getShortType(latestStage?.latestEvent?.type)] || EVENT_ICONS.default
    : EVENT_ICONS.default;
  const steady = (!waiting && !loading) || pendingApproval;
  const statusText = getStatusText(sequence.state, { successful, faulty, waiting, aborted, timedOut });
  const evaluation = stage?.latestEvaluation;
  return {
    finished,
    loading,
    waiting,
    pendingApproval,
    successful,
    warning,
    timedOut,
    faulty,
    aborted,
    steady,
    icon,
    statusText,
    evaluation,
  };
}

export function getStage(sequence: ISequence, stageName?: string): SeqStage | undefined {
  return sequence.stages.find((stage) => stage.name === stageName);
}

export function getStageNames(sequence: ISequence): string[] {
  return sequence.stages.map((stage) => stage.name);
}

export function getLastStageName(sequence: ISequence): string | undefined {
  return sequence.stages[sequence.stages.length - 1]?.name;
}

export function isFinished(state: SequenceState): boolean {
  return finishedStates.includes(state);
}

export function getShortType(type: string): string {
  const parts = type?.split('.');
  if (parts.length === 6) {
    return parts[4];
  } else if (parts.length === 5) {
    return parts[3];
  } else {
    return type;
  }
}

export class Sequence implements ISequence {
  name!: string;
  project!: string;
  service!: string;
  shkeptncontext!: string;
  state!: SequenceState;
  time!: string;
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
    return isFinished(state);
  }

  public static getShortType(type: string): string {
    return getShortType(type);
  }

  public getStage(stageName: string): SeqStage | undefined {
    return this.stages.find((stage) => stage.name === stageName);
  }

  public getStageTime(stageName: string): string | undefined {
    return this.findTrace((trace) => trace.stage === stageName)?.time;
  }

  public getStages(): string[] {
    return getStageNames(this);
  }

  public getLastStage(): string | undefined {
    return getLastStageName(this);
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
      this.isFinished(stageName) &&
      !this.isFaulty(stageName) &&
      !this.isWarning(stageName) &&
      !this.isAborted(stageName)
    );
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
    return this.state === SequenceState.WAITING;
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
