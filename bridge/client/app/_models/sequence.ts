import {ResultTypes} from './result-types';
import {Trace} from './trace';
import {EventTypes} from './event-types';
import {EvaluationResult} from './evaluation-result';
import {EVENT_ICONS} from './event-icons';

const DEFAULT_ICON = 'information';

export class Sequence {
  name!: string;
  project!: string;
  service!: string;
  shkeptncontext!: string;
  stages!: [
    {
      image?: string,
      latestEvaluation?: EvaluationResult,
      latestEvent?: {
        id: string,
        time: string,
        type: string
      },
      latestFailedEvent?: {
        id: string,
        time: string,
        type: string
      },
      name: string
    }
  ];
  state!: 'triggered' | 'finished' | 'waiting';
  time!: string;
  problemTitle?: string;
  traces: Trace[] = [];

  public static fromJSON(data: unknown): Sequence {
    return Object.assign(new this(), data);
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
    return stageName ? (this.getStage(stageName)?.latestEvent?.type.endsWith('finished') ?? false) : this.state === 'finished';
  }

  public getEvaluation(stage: string): EvaluationResult | undefined {
    return this.getStage(stage)?.latestEvaluation;
  }

  public hasPendingApproval(stageName?: string): boolean {
    return stageName ?
        this.getStage(stageName)?.latestEvent?.type === EventTypes.APPROVAL_TRIGGERED
      : this.stages.some(stage => stage.latestEvent?.type === EventTypes.APPROVAL_TRIGGERED);
  }

  public getStatus(): string {
    let status: string = this.state;
    if (this.state === 'finished') {
      if (this.stages.some(stage => stage.latestFailedEvent)) {
        status = 'failed';
      }
      else {
        status = 'succeeded';
      }
    }
    else if (this.state === 'triggered') {
      status = 'started';
    }
    return status;
  }

  public isLoading(stageName?: string): boolean {
    return stageName ? this.state === 'triggered' && !this.isFinished(stageName) : this.state === 'triggered';
  }

  public isSuccessful(stageName?: string): boolean {
    return stageName ? !this.isFaulty(stageName) && this.isFinished(stageName) : this.state === 'finished' && !this.isFaulty();
  }

  public isWarning(stageName: string): boolean {
    return this.getStage(stageName)?.latestEvaluation?.result === ResultTypes.WARNING;
  }

  public isWaiting(stageName?: string): boolean {
    return stageName ? !this.isFinished(stageName) && this.state === 'waiting' : this.state === 'waiting';
  }

  public isRemediation(): boolean {
    return this.name === 'remediation';
  }

  public getLatestEvent() {
    return this.stages[this.stages.length - 1]?.latestEvent;
  }

  public getIcon(stageName?: string): string {
    let icon;
    if (this.state === 'waiting') {
      icon = EVENT_ICONS.waiting;
    }
    else {
      const stage = stageName ? this.getStage(stageName) : this.stages[this.stages.length - 1];
      icon = stage?.latestEvent?.type ? EVENT_ICONS[Sequence.getShortType(stage?.latestEvent?.type)] || DEFAULT_ICON : DEFAULT_ICON;
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

  private getFirstTrace(): Trace | null {
    return this.traces[0];
  }
}
