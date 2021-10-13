import moment from 'moment';
import { EventTypes } from '../../../shared/interfaces/event-types';
import { ResultTypes } from '../../../shared/models/result-types';
import { ApprovalStates } from './approval-states';
import { EVENT_ICONS } from './event-icons';
import { ProblemStates } from './problem-states';
import { DateUtil } from '../_utils/date.utils';
import { Trace as tc, TraceData } from '../../../shared/models/trace';
import { DtIconType } from '@dynatrace/barista-icons';
import { KeptnService } from '../../../shared/models/keptn-service';

class Trace extends tc {
  traces: Trace[] = [];
  triggeredid?: string;
  type!: EventTypes | string;
  data!: Omit<TraceData, 'evaluationHistory'> & {
    evaluationHistory?: Trace[];
  };
  started?: boolean;
  finished?: boolean;
  source?: string;
  label?: string;
  heatmapLabel?: string;
  icon?: DtIconType;
  image?: string;
  plainEvent?: string;
  time?: Date;
  labelMap?: Map<string, string>;

  static fromJSON(data: unknown): Trace {
    if (data instanceof Trace) {
      return data;
    }

    const plainEvent = JSON.parse(JSON.stringify(data));
    const trace: Trace = Object.assign(new this(), data, { plainEvent });

    if (trace.data?.evaluationHistory?.length) {
      trace.data.evaluationHistory = trace.data.evaluationHistory.map((t) => Trace.fromJSON(t));
    }

    return trace;
  }

  static traceMapper(traces: Trace[]): Trace[] {
    traces = traces.map((trace) => Trace.fromJSON(trace)).sort(DateUtil.compareTraceTimesDesc);

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
          .reduce((lastSeq: Trace | undefined, s: Trace) => (s.stage === trace.stage ? s : lastSeq), undefined)
          ?.traces.push(trace);
      } else {
        seq.push(trace);
      }

      return seq;
    }, []);
  }

  get project(): string | undefined {
    return this.data.project;
  }

  get service(): string | undefined {
    return this.data.service;
  }

  get stage(): string | undefined {
    return this.data.stage;
  }

  get labels(): Map<string, string> | undefined {
    if (!this.labelMap) {
      let map: Map<string, string> | undefined;
      if (this.data.labels) {
        map = new Map<string, string>();
        for (const key of Object.keys(this.data.labels)) {
          map.set(key, this.data.labels[key]);
        }
      }
      this.labelMap = map;
    }
    return this.labelMap;
  }

  isFaulty(stageName?: string): boolean {
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

  isWarning(stageName?: string): boolean {
    let result = false;
    if (this.getFinishedEvent()?.data.result === ResultTypes.WARNING) {
      result = stageName ? this.data.stage === stageName : true;
    }
    return result;
  }

  isSuccessful(stageName?: string): boolean {
    let result = false;
    if (
      (this.isFinished() && this.getFinishedEvent()?.data.result === ResultTypes.PASSED) ||
      (this.isApprovalFinished() && this.isApproved()) ||
      (this.isProblem() && this.isProblemResolvedOrClosed()) ||
      this.isSuccessfulRemediation()
    ) {
      result = stageName ? this.data.stage === stageName : true;
    }
    return !this.isFaulty() && result;
  }

  public isFailed(): boolean {
    return (
      this.getFinishedEvent()?.data.result === ResultTypes.FAILED || (this.isApprovalFinished() && this.isDeclined())
    );
  }

  public isProblem(): boolean {
    return this.type === EventTypes.PROBLEM_DETECTED || this.type === EventTypes.PROBLEM_OPEN;
  }

  public isRemediation(): boolean {
    return this.type.endsWith(EventTypes.REMEDIATION_TRIGGERED_SUFFIX);
  }

  public isRemediationAction(): boolean {
    return this.type === EventTypes.ACTION_TRIGGERED;
  }

  public getRemediationActionDescription(): string | undefined {
    return this.data.action?.description;
  }

  public getRemediationActionName(): string | undefined {
    return this.data.action?.name;
  }

  public isProblemResolvedOrClosed(): boolean {
    if (!this.traces || this.traces.length === 0) {
      return this.data.State === ProblemStates.RESOLVED || this.data.State === ProblemStates.CLOSED;
    } else {
      return this.traces.some((t) => t.isProblem() && t.isProblemResolvedOrClosed());
    }
  }

  public isSuccessfulRemediation(): boolean {
    if (!this.traces || this.traces.length === 0) {
      return this.type.endsWith(EventTypes.REMEDIATION_FINISHED_SUFFIX) && this.data.result !== ResultTypes.FAILED;
    } else {
      return this.traces.some((t) => t.isSuccessfulRemediation());
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

  private isApprovalFinished(): boolean {
    return this.type === EventTypes.APPROVAL_FINISHED;
  }

  public isDirectDeployment(): boolean {
    return this.type === EventTypes.DEPLOYMENT_FINISHED && this.data?.deployment?.deploymentstrategy === 'direct';
  }

  private isApproved(): boolean {
    return this.data.approval?.result === ApprovalStates.APPROVED;
  }

  private isDeclined(): boolean {
    return this.data.approval?.result === ApprovalStates.DECLINED;
  }

  public isDeployment(): string | undefined {
    return this.type === EventTypes.DEPLOYMENT_TRIGGERED ? this.data.stage : undefined;
  }

  public isEvaluation(): string | undefined {
    return this.type.endsWith(EventTypes.EVALUATION_TRIGGERED_SUFFIX) && !this.isSequence()
      ? this.data.stage
      : undefined;
  }

  public isEvaluationInvalidation(): boolean {
    return this.type === EventTypes.EVALUATION_INVALIDATED;
  }

  public getEvaluationFinishedEvent(): Trace | undefined {
    return this.findTrace(
      (trace) => trace.source === KeptnService.LIGHTHOUSE_SERVICE && trace.type.endsWith(EventTypes.EVALUATION_FINISHED)
    );
  }

  hasLabels(): boolean {
    return Object.keys(this.data.labels || {}).length > 0;
  }

  getLabel(): string {
    if (!this.label) {
      this.label = this.getShortType();
    }

    return this.label;
  }

  getProblemTitle(): string | undefined {
    return this.data.problem?.ProblemTitle;
  }

  getShortType(): string {
    const parts = this.type.split('.');
    if (parts.length === 6) {
      return parts[4];
    } else if (parts.length === 5) {
      return parts[3];
    } else {
      return this.type;
    }
  }

  getIcon(): DtIconType {
    if (!this.icon) {
      this.icon = EVENT_ICONS[this.getShortType()] || EVENT_ICONS.default;
    }
    return this.icon;
  }

  getChartLabel(): string {
    return this.data.labels?.buildId ?? moment(this.time).format('YYYY-MM-DD HH:mm');
  }

  getHeatmapLabel(): string {
    if (!this.heatmapLabel) {
      this.heatmapLabel = this.getChartLabel();
    }
    return this.heatmapLabel;
  }

  setHeatmapLabel(label: string): void {
    this.heatmapLabel = label;
  }

  isStarted(): boolean {
    if (!this.started && this.traces) {
      this.started = this.traces.some((t) => t.type.endsWith('.started') || t.isStarted());
    }

    return !!this.started;
  }

  isChanged(): boolean {
    return this.type.endsWith('.changed');
  }

  isFinished(): boolean {
    if (!this.finished) {
      if (!this.traces || this.traces.length === 0) {
        this.finished = this.type.endsWith('.finished');
      } else if (this.isProblem()) {
        this.finished = this.isProblemResolvedOrClosed();
      } else {
        const countStarted = this.traces.filter((t) => t.type.endsWith('.started')).length;
        const countFinished = this.traces.filter((t) => t.type.endsWith('.finished')).length;
        this.finished = countFinished >= countStarted && countFinished !== 0;
      }
    }

    return this.finished;
  }

  isTriggered(): boolean {
    return this.type.endsWith('.triggered');
  }

  isLoading(): boolean {
    return this.isStarted() && !this.isFinished();
  }

  isInvalidated(): boolean {
    return !!this.traces.find((e) => e.isEvaluationInvalidation() && e.triggeredid === this.id);
  }

  getFinishedEvent(): Trace | undefined {
    return this.type.endsWith('.finished') ? this : this.traces.find((t) => t.type.endsWith('.finished'));
  }

  getRemediationAction(): Trace | undefined {
    return this.findTrace((t) => t.isRemediationAction());
  }

  getEvaluation(stageName: string): Trace | undefined {
    return this.findTrace((t) => !!t.isEvaluation() && t.stage === stageName);
  }

  getDeploymentUrl(): string | undefined {
    return this.data.deployment?.deploymentURIsPublic?.find(() => true);
  }

  findTrace(comp: (args: Trace) => boolean): Trace | undefined {
    if (comp(this)) {
      return this;
    } else {
      return this.traces.reduce(
        (result: Trace | undefined, trace: Trace) => result || trace.findTrace(comp),
        undefined
      );
    }
  }

  findLastTrace(comp: (args: Trace) => boolean): Trace | undefined {
    if (comp(this)) {
      return this;
    } else {
      return this.traces.reduce((result: Trace | undefined, trace) => trace.findTrace(comp) || result, undefined);
    }
  }

  getLastTrace(): Trace {
    return this.traces.length ? this.traces[this.traces.length - 1].getLastTrace() : this;
  }

  isSequence(): boolean {
    return this.type.split('.').length === 6 && !!this.stage && this.type.includes(this.stage);
  }

  getProblemDetails(): string | undefined {
    return this.data.problem?.ImpactedEntity || this.data.problem?.ProblemTitle;
  }
}

export { Trace };
