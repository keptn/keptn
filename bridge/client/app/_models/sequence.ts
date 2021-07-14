import {ResultTypes} from './result-types';

export class Sequence {
  name: string;
  project: string;
  service: string;
  shkeptncontext: string;
  stages: [
    {
      image: string,
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
  state: 'triggered' | 'finished' | 'waiting';
  time: string;
  problemTitle?: string;
  traces: Trace[] = [];

  static fromJSON(data: any): Sequence {
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
      !!this.getStage(stageName).latestFailedEvent
      : this.stages.some(stage => stage.latestFailedEvent);
  }

  public isFinished(stageName?: string): boolean {
    return stageName ? this.getStage(stageName)?.latestEvent.type.endsWith('finished') : this.state === 'finished';
  }

  public getEvaluation(stage: string): EvaluationResult {
    return this.getStage(stage).latestEvaluation;
  }

  public hasPendingApproval(stageName?: string): boolean {
    return stageName ?
        this.getStage(stageName)?.latestEvent.type === EventTypes.APPROVAL_TRIGGERED
      : this.stages.some(stage => stage.latestEvent.type === EventTypes.APPROVAL_TRIGGERED);
  }

  public getStatus(): string {
    let status: any = this.state;
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
      icon = EVENT_ICONS[Sequence.getShortType(stage?.latestEvent.type)] || DEFAULT_ICON;
    }
    return icon;
  }

  public getLabel(): string {
    return this.name;
  }

  public getStatusLabel(): string {
    return this.state;
  }

  public getShortImageName(): string | null {
    return this.stages[0]?.image?.split('/').pop();
  }

  public findTrace(comp: (args: Trace) => any): Trace {
    return this.traces.reduce((result, trace) => result || trace.findTrace(comp), null);
  }

  public findLastTrace(comp: (args: Trace) => any): Trace {
    return this.traces.reduce((result, trace) => trace.findTrace(comp) || result, null);
  }

  public getLabels(): Map<string, string> {
    return this.getLastTrace()?.getFinishedEvent()?.data.labels || this.getFirstTrace()?.data.labels;
  }

  private getLastTrace(): Trace {
    return this.traces[this.traces.length - 1];
  }

  private getFirstTrace(): Trace | null {
    return this.traces[0];
  }

  public getStage(stageName: string) {
    return this.stages.find(stage => stage.name === stageName);
  }

  static fromJSON(data: any): Sequence {
    return Object.assign(new this(), data);
  }
}
