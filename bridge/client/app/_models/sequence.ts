import { ResultTypes } from '../../../shared/models/result-types';
import { Trace } from './trace';
import { EventTypes } from '../../../shared/interfaces/event-types';
import { EvaluationResult } from '../../../shared/interfaces/evaluation-result';
import { EVENT_ICONS } from './event-icons';
import { RemediationAction } from '../../../shared/models/remediation-action';
import { Sequence as sq, SequenceStage, SequenceState } from '../../../shared/models/sequence';
import { DtIconType } from '@dynatrace/barista-icons';

export class Sequence extends sq {
  stages!: (SequenceStage &
     {
      latestEvaluationTrace?: Trace,
      actions?: RemediationAction[]
    })[]
  ;
  problemTitle?: string;
  traces: Trace[] = [];

  public static fromJSON(data: unknown): Sequence {
    const sequence = Object.assign(new this(), data);
    for (const stage of sequence.stages) {
      stage.actions = stage.actions?.map(s => RemediationAction.fromJSON(s)) ?? [];
      if (stage.latestEvaluationTrace) {
        stage.latestEvaluationTrace = Trace.fromJSON(stage.latestEvaluationTrace);
      }
    }
    return sequence;
  }

  public static getShortType(type: string): string {
    const parts = type.split('.');
    if (parts.length === 6) {
      return parts[4];
    }
    else if (parts.length === 5) {
      return parts[3];
    }
    else {
      return type;
    }
  }

  public getStage(stageName: string) {
    return this.stages.find(stage => stage.name === stageName);
  }

  public getStages(): string[] {
    return  this.stages.map(stage => stage.name);
  }

  public getLastStage(): string | undefined {
    return this.stages[this.stages.length - 1]?.name;
  }

  public isFaulty(stageName?: string): boolean {
    return stageName ?
      !!this.getStage(stageName)?.latestFailedEvent
      : this.stages.some(stage => stage.latestFailedEvent);
  }

  public isFinished(stageName?: string): boolean {
    return stageName ? (this.getStage(stageName)?.latestEvent?.type.endsWith(SequenceState.FINISHED) ?? false) : this.state === SequenceState.FINISHED;
  }

  public getEvaluation(stage: string): EvaluationResult | undefined {
    return this.getStage(stage)?.latestEvaluation;
  }

  public getEvaluationTrace(stage: string): Trace | undefined {
    return this.getStage(stage)?.latestEvaluationTrace;
  }

  public hasPendingApproval(stageName?: string): boolean {
    return stageName ?
        this.getStage(stageName)?.latestEvent?.type === EventTypes.APPROVAL_TRIGGERED || this.getStage(stageName)?.latestEvent?.type === EventTypes.APPROVAL_STARTED
      : this.stages.some(stage => stage.latestEvent?.type === EventTypes.APPROVAL_TRIGGERED || stage.latestEvent?.type === EventTypes.APPROVAL_STARTED);
  }

  public getStatus(): string {
    let status: string = this.state;
    if (this.state === SequenceState.FINISHED) {
      if (this.stages.some(stage => stage.latestFailedEvent)) {
        status = 'failed';
      }
      else {
        status = 'succeeded';
      }
    }
    else if (this.state === SequenceState.TRIGGERED) {
      status = 'started';
    }
    return status;
  }

  public isLoading(stageName?: string): boolean {
    return stageName ? this.state === SequenceState.TRIGGERED && !this.isFinished(stageName) : this.state === SequenceState.TRIGGERED;
  }

  public isSuccessful(stageName?: string): boolean {
    return stageName ? !this.isFaulty(stageName) && this.isFinished(stageName) : this.state === SequenceState.FINISHED && !this.isFaulty();
  }

  public isWarning(stageName: string): boolean {
    return this.getStage(stageName)?.latestEvaluation?.result === ResultTypes.WARNING;
  }

  public isWaiting(stageName?: string): boolean {
    return stageName ? !this.isFinished(stageName) && this.state === SequenceState.WAITING : this.state === SequenceState.WAITING;
  }

  public isRemediation(): boolean {
    return this.name === 'remediation';
  }

  public getLatestEvent(): {id: string, time: string, type: string} | undefined {
    return this.stages[this.stages.length - 1]?.latestEvent;
  }

  public getIcon(stageName?: string): DtIconType {
    let icon;
    if (this.state === SequenceState.WAITING) {
      icon = EVENT_ICONS.waiting;
    }
    else {
      const stage = stageName ? this.getStage(stageName) : this.stages[this.stages.length - 1];
      icon = stage?.latestEvent?.type ? EVENT_ICONS[Sequence.getShortType(stage?.latestEvent?.type)] || EVENT_ICONS.default : EVENT_ICONS.default;
    }
    return icon;
  }

  public getShortImageName(): string | undefined {
    return this.stages[0]?.image?.split('/').pop();
  }

  public findTrace(comp: (args: Trace) => boolean): Trace | undefined {
    return this.traces.reduce((result: Trace | undefined, trace: Trace) => result || trace.findTrace(comp), undefined);
  }

  public findLastTrace(comp: (args: Trace) => boolean): Trace | undefined {
    return this.traces.reduce((result: Trace | undefined, trace: Trace) => trace.findTrace(comp) || result, undefined);
  }

  public getLabels(): Map<string, string> | undefined {
    return (this.getLastTrace()?.getFinishedEvent()?.labels || this.getFirstTrace()?.labels);
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
}
