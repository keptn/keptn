import { EventTypes } from '../interfaces/event-types';
import { ResultTypes } from './result-types';
import { IndicatorResult } from '../interfaces/indicator-result';
import { ProblemStates } from '../../client/app/_models/problem-states';
import { ApprovalStates } from '../../client/app/_models/approval-states';
import { KeptnService } from './keptn-service';
import { DateUtil } from '../utils/date.utils';

export interface TraceData {
  project?: string;
  service?: string;
  stage?: string;

  image?: string;
  tag?: string;

  deployment?: {
    deploymentNames: string[];
    deploymentURIsLocal: string[];
    deploymentURIsPublic: string[];
    deploymentstrategy: string;
    gitCommit: string;
  };

  deploymentURILocal?: string;
  deploymentURIPublic?: string;

  message?: string;

  labels?: { [key: string]: string };
  result?: ResultTypes;
  teststrategy?: string;

  start?: Date;
  end?: Date;

  canary?: {
    action: string;
    value: number;
  };
  eventContext?: {
    shkeptncontext: string;
    token: string;
  };
  configurationChange?: {
    values: {
      image: unknown;
    };
  };

  evaluation?: {
    comparedEvents?: string[];
    indicatorResults: IndicatorResult[];
    result: ResultTypes;
    score: number;
    sloFileContent: string;
    timeEnd: Date;
    timeStart: Date;
    // eslint-disable-next-line @typescript-eslint/naming-convention
    score_pass: string;
    // eslint-disable-next-line @typescript-eslint/naming-convention
    score_warning: string;
    // eslint-disable-next-line @typescript-eslint/naming-convention
    compare_with: string;
    // eslint-disable-next-line @typescript-eslint/naming-convention
    include_result_with_score: string;
    // eslint-disable-next-line @typescript-eslint/naming-convention
    number_of_comparison_results: number;
    // eslint-disable-next-line @typescript-eslint/naming-convention
    number_of_missing_comparison_results: number;
    sloFileContentParsed: string;
  };

  evaluationHistory?: Trace[];

  problem?: {
    // eslint-disable-next-line @typescript-eslint/naming-convention
    ProblemTitle: string;
    // eslint-disable-next-line @typescript-eslint/naming-convention
    ImpactedEntity: string;

    ProblemDetails: {
      tagsOfAffectedEntities: {
        key: string;
        value: string;
      }[];
    };
  };

  approval?: {
    result: string;
    status: string;
  };

  action?: {
    action: string;
    description: string;
    name: string;
  };
  // eslint-disable-next-line @typescript-eslint/naming-convention
  Tags?: string;
  // eslint-disable-next-line @typescript-eslint/naming-convention
  State?: string;
}

export class Trace {
  id!: string;
  shkeptncontext!: string;
  triggeredid?: string;
  type!: EventTypes | string;
  time?: string; // 2021-10-29T08:43:11.702Z
  data!: TraceData;
  traces: Trace[] = [];
  finished?: boolean;
  source?: string;
  label?: string;

  public getShortImageName(): string | undefined {
    let image;
    if (this.data.image && this.data.tag) {
      image = [this.data.image.split('/').pop(), this.data.tag].join(':');
    } else if (this.data.image) {
      image = this.data.image.split('/').pop();
    } else if (this.data.configurationChange?.values) {
      image = this.getConfigurationChangeImage();
    }
    return image;
  }

  public getConfigurationChangeImage(): string | undefined {
    return typeof this.data.configurationChange?.values.image === 'string'
      ? this.data.configurationChange.values.image.split('/').pop()
      : undefined;
  }

  public getDeploymentUrl(): string | undefined {
    return this.data.deployment?.deploymentURIsPublic?.find(() => true);
  }

  protected static traceMapperGlobal<T extends Trace>(traces: Trace[]): T[] {
    traces.sort(DateUtil.compareTraceTimesDesc);
    return traces.reduce((seq: Trace[], trace: Trace) => {
      let trigger: Trace | undefined;
      if (trace.triggeredid) {
        trigger = traces.reduce(
          (acc: Trace | undefined, r: Trace) => acc || r.findTrace((t) => t.id === trace.triggeredid),
          undefined
        );
      } else if (trace.isProblem() && trace.isProblemResolvedOrClosed()) {
        trigger = traces.reduce(
          (acc: Trace | undefined, r: Trace) =>
            acc || r.findTrace((t) => t.isProblem() && !t.isProblemResolvedOrClosed()),
          undefined
        );
      } else if (trace.isFinished()) {
        trigger = traces.reduce(
          (acc: Trace | undefined, r: Trace) =>
            acc || r.findTrace((t) => !t.triggeredid && t.type.slice(0, -8) === trace.type.slice(0, -9)),
          undefined
        );
      }

      if (trigger) {
        trigger.traces.push(trace);
      } else if (trace.isSequence()) {
        seq.push(trace);
      } else if (seq.length > 0) {
        seq
          .reduce(
            (lastSeq: Trace | undefined, s: Trace) => (s.data.stage === trace.data.stage ? s : lastSeq),
            undefined
          )
          ?.traces.push(trace);
      } else {
        seq.push(trace);
      }

      return seq;
    }, []) as T[];
  }

  public findTrace<T extends Trace>(this: T, comp: (args: T) => boolean): T | undefined {
    if (comp(this)) {
      return this;
    } else {
      return (this.traces as T[]).reduce(
        (result: T | undefined, trace: T) => result || trace.findTrace(comp),
        undefined
      );
    }
  }

  public isProblem(): boolean {
    return this.type === EventTypes.PROBLEM_DETECTED || this.type === EventTypes.PROBLEM_OPEN;
  }

  public isProblemResolvedOrClosed(): boolean {
    if (!this.traces || this.traces.length === 0) {
      return this.data.State === ProblemStates.RESOLVED || this.data.State === ProblemStates.CLOSED;
    } else {
      return this.traces.some((t) => t.isProblem() && t.isProblemResolvedOrClosed());
    }
  }

  public isSequence(): boolean {
    return this.type.split('.').length === 6 && !!this.data.stage && this.type.includes(this.data.stage);
  }

  public isFinished(): boolean {
    if (!this.finished) {
      if (!this.traces || this.traces.length === 0) {
        this.finished = this.isFinishedEvent();
      } else if (this.isProblem()) {
        this.finished = this.isProblemResolvedOrClosed();
      } else {
        const countStarted = this.traces.filter((t) => t.isStartedEvent()).length;
        const countFinished = this.traces.filter((t) => t.isFinishedEvent()).length;
        this.finished = countFinished >= countStarted && countFinished !== 0;
      }
    }

    return this.finished;
  }

  public isFaulty(stageName?: string): boolean {
    let result = false;
    if (this.data) {
      if (
        this.isFailed() ||
        (this.isProblem() && !this.isProblemResolvedOrClosed()) ||
        (this.isRemediation() && !this.isSuccessfulRemediation()) ||
        this.traces.some((t) => t.isFaulty())
      ) {
        result = stageName ? this.data.stage === stageName : true;
      }
    }
    return result;
  }

  public isSuccessfulRemediation(): boolean {
    if (!this.traces || this.traces.length === 0) {
      return this.type.endsWith(EventTypes.REMEDIATION_FINISHED_SUFFIX) && this.data.result !== ResultTypes.FAILED;
    } else {
      return this.traces.some((t) => t.isSuccessfulRemediation());
    }
  }

  protected isApprovalFinished(): boolean {
    return this.type === EventTypes.APPROVAL_FINISHED;
  }

  public isFailed(): boolean {
    return (
      this.getFinishedEvent()?.data.result === ResultTypes.FAILED || (this.isApprovalFinished() && this.isDeclined())
    );
  }

  public getFinishedEvent<T extends Trace>(this: T): T | undefined {
    return (this.isFinishedEvent() ? this : this.traces.find((t) => t.type.endsWith('.finished'))) as T;
  }

  private isDeclined(): boolean {
    return this.data.approval?.result === ApprovalStates.DECLINED;
  }

  public isRemediation(): boolean {
    return this.type.endsWith(EventTypes.REMEDIATION_TRIGGERED_SUFFIX);
  }

  public isWarning(stageName?: string): boolean {
    let result = false;
    if (this.getFinishedEvent()?.data.result === ResultTypes.WARNING) {
      result = stageName ? this.data.stage === stageName : true;
    }
    return result;
  }

  public getEvaluation<T extends Trace>(this: T, stageName: string): T | undefined {
    return this.findTrace((t) => !!t.isEvaluation() && t.data.stage === stageName);
  }

  public isEvaluation(): string | undefined {
    return this.type.endsWith(EventTypes.EVALUATION_TRIGGERED_SUFFIX) && !this.isSequence()
      ? this.data.stage
      : undefined;
  }

  public getEvaluationFinishedEvent<T extends Trace>(this: T, stage?: string): T | undefined {
    return this.findTrace(
      (trace) =>
        trace.source === KeptnService.LIGHTHOUSE_SERVICE &&
        trace.type.endsWith(EventTypes.EVALUATION_FINISHED) &&
        (!stage || trace.data.stage === stage)
    );
  }

  public getLabel(): string {
    if (!this.label) {
      this.label = this.getShortType();
    }
    return this.label;
  }

  public getShortType(): string {
    const parts = this.type.split('.');
    if (parts.length === 6) {
      return parts[4];
    } else if (parts.length === 5) {
      return parts[3];
    } else {
      return this.type;
    }
  }

  public isApproval(): string | undefined {
    return this.type === EventTypes.APPROVAL_TRIGGERED || this.type === EventTypes.APPROVAL_STARTED
      ? this.data.stage
      : undefined;
  }

  public isApprovalPending(): boolean {
    let pending = true;
    for (let i = 0; i < this.traces.length && pending; ++i) {
      if (this.traces[i].isApprovalFinished()) {
        pending = false;
      }
    }
    return pending;
  }

  public isChangedEvent(): boolean {
    return this.type.endsWith('.changed');
  }

  public isFinishedEvent(): boolean {
    return this.type.endsWith('.finished');
  }

  public isStartedEvent(): boolean {
    return this.type.endsWith('.started');
  }
}
