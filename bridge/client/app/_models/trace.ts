import moment from 'moment';
import {EventTypes} from './event-types';
import {ResultTypes} from './result-types';
import {ApprovalStates} from './approval-states';
import {EVENT_ICONS} from './event-icons';
import {ProblemStates} from './problem-states';
import {DateUtil} from '../_utils/date.utils';
import { IndicatorResult } from './indicator-result';

const DEFAULT_ICON = 'information';

class Trace {
  traces: Trace[] = [];
  id!: string;
  shkeptncontext!: string;
  triggeredid?: string;
  type!: EventTypes | string;
  data!: {
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
        image: unknown
      }
    };

    evaluation?: {
      comparedEvents?: string[];
      indicatorResults: IndicatorResult[];
      result: string;
      score: number;
      sloFileContent: string;
      timeEnd: Date;
      timeStart: Date;
      score_pass: string;
      score_warning: string;
      compare_with: string;
      include_result_with_score: string;
      number_of_comparison_results: number;
      number_of_missing_comparison_results: number;
      sloFileContentParsed: string;
    };

    evaluationHistory?: Trace[];

    problem?: {
      ProblemTitle: string;
      ImpactedEntity: string;
      ProblemDetails: {
        tagsOfAffectedEntities: {
          key: string;
          value: string;
        }[]
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
    }

    Tags?: string;
    State?: string;
  };
  started?: boolean;
  finished?: boolean;
  source?: string;
  label?: string;
  heatmapLabel?: string;
  icon?: string;
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

    // if (trace?.evaluationHistory?.length > 0) {
    //   trace.evaluationHistory = trace.evaluationHistory.map(t => Trace.fromJSON(t));
    // }// TODO: Where does this come from? It comes from the mock for 10 SLIs (ktb-evaluation-details.component.spec.ts)

    if (trace.data?.evaluationHistory?.length) {
      trace.data.evaluationHistory = trace.data.evaluationHistory.map(t => Trace.fromJSON(t));
    }

    return trace;
  }

  static traceMapper(traces: Trace[]) {
    traces = traces
      .map(trace => Trace.fromJSON(trace))
      .sort(DateUtil.compareTraceTimesDesc);

    return traces.reduce((seq: Trace[], trace: Trace) => {
      let trigger: Trace | undefined;
      if (trace.triggeredid) {
        trigger = traces.reduce((acc: Trace | undefined, r: Trace) =>
          acc
          || r.findTrace((t) => t.id === trace.triggeredid), undefined);
      } else if (trace.isProblem() && trace.isProblemResolvedOrClosed()) {
        trigger = traces.reduce((acc: Trace | undefined, r: Trace) =>
          acc
          || r.findTrace((t) => t.isProblem() && !t.isProblemResolvedOrClosed()), undefined);
      } else if (trace.isFinished()) {
        trigger = traces.reduce((acc: Trace | undefined, r: Trace) =>
          acc
          || r.findTrace((t) => !t.triggeredid && t.type.slice(0, -8) === trace.type.slice(0, -9)), undefined);
      }

      if (trigger) {
        trigger.traces.push(trace);
      } else if (trace.isSequence()) {
        seq.push(trace);
      } else if (seq.length > 0) {
        seq.reduce((lastSeq: Trace | undefined, s: Trace) => {
          return s.stage === trace.stage ? s : lastSeq;
        }, undefined)?.traces.push(trace);
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
      if (this.isFailed() ||
        (this.isProblem() && !this.isProblemResolvedOrClosed()) ||
        (this.isRemediation() && !this.isSuccessfulRemediation()) ||
        this.traces.some(t => t.isFaulty())) {
        result = stageName ? this.data.stage === stageName : true;
      }
    }
    return result;
  }

  isFailedEvaluation(): string | undefined {
    let result: string | undefined;
    if (this.data) {
      if (this.getFinishedEvent()?.type === EventTypes.EVALUATION_FINISHED && this.isFailed()) {
        result = this.data.stage;
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
    if (this.isFinished()
        && this.getFinishedEvent()?.data.result === ResultTypes.PASSED
      || this.isApprovalFinished()
        && this.isApproved()
      || this.isProblem()
        && this.isProblemResolvedOrClosed()
      || this.isSuccessfulRemediation()) {
      result = stageName ? this.data.stage === stageName : true;
    }
    return !this.isFaulty() && result;
  }

  public isFailed(): boolean {
    return this.getFinishedEvent()?.data.result === ResultTypes.FAILED || this.isApprovalFinished() && this.isDeclined();
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

  public getRemediationActionDetails(): string | undefined {
    return this.data.action?.description || this.data.action?.name;
  }

  public isProblemResolvedOrClosed(): boolean {
    if (!this.traces || this.traces.length === 0) {
      return this.data.State === ProblemStates.RESOLVED || this.data.State === ProblemStates.CLOSED;
    }
    else {
      return this.traces.some(t => t.isProblem() && t.isProblemResolvedOrClosed());
    }
  }

  public isSuccessfulRemediation(): boolean {
    if (!this.traces || this.traces.length === 0) {
      return this.type.endsWith(EventTypes.REMEDIATION_FINISHED_SUFFIX) && this.data.result !== ResultTypes.FAILED;
    }
    else {
      return this.traces.some(t => t.isSuccessfulRemediation());
    }
  }

  public isApproval(): string | undefined {
    return this.type === EventTypes.APPROVAL_TRIGGERED ? this.data.stage : undefined;
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
    return this.type.endsWith(EventTypes.EVALUATION_TRIGGERED_SUFFIX) && !this.isSequence() ? this.data.stage : undefined;
  }

  public isEvaluationInvalidation(): boolean {
    return this.type === EventTypes.EVALUATION_INVALIDATED;
  }

  hasLabels(): boolean {
    return Object.keys(this.data.labels || {}).length > 0;
  }

  getLabel(): string {
    // TODO: use translation file; see also EVENT_LABELS in event-labels.ts
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
    }
    else if (parts.length === 5) {
      return parts[3];
    }
    else {
      return this.type;
    }
  }

  getIcon(): string {
    if (!this.icon) {
      this.icon = EVENT_ICONS[this.getShortType()] || DEFAULT_ICON;
    }
    return this.icon;
  }

  getShortImageName(): string | undefined {
    if (!this.image) {
      if (this.data.image && this.data.tag) {
        this.image = [this.data.image.split('/').pop(), this.data.tag].join(':');
      }
      else if (this.data.image) {
        this.image = this.data.image.split('/').pop();
      }
      else if (this.data.configurationChange?.values) {
        this.image = this.getConfigurationChangeImage();
      }
    }

    return this.image;
  }

  public getConfigurationChangeImage(): string | undefined {
    return typeof this.data.configurationChange?.values.image === 'string'
      ? this.data.configurationChange.values.image.split('/').pop()
      : undefined;
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

  setHeatmapLabel(label: string) {
    this.heatmapLabel = label;
  }

  isStarted() {
    if (!this.started && this.traces) {
      this.started = this.traces.some(t => t.type.endsWith('.started') || t.isStarted());
    }

    return this.started;
  }

  isChanged() {
    return this.type.endsWith('.changed');
  }

  isFinished() {
    if (!this.finished) {
      if (!this.traces || this.traces.length === 0) {
        this.finished = this.type.endsWith('.finished');
      } else if (this.isProblem()) {
        this.finished = this.isProblemResolvedOrClosed();
      } else {
        const countStarted = this.traces.filter(t => t.type.endsWith('.started')).length;
        const countFinished = this.traces.filter(t => t.type.endsWith('.finished')).length;
        this.finished = countFinished >= countStarted && countFinished !== 0;
      }
    }

    return this.finished;
  }

  isTriggered() {
    return this.type.endsWith('.triggered');
  }

  isLoading() {
    return this.isStarted() && !this.isFinished();
  }

  isInvalidated() {
    return !!this.traces.find(e => e.isEvaluationInvalidation() && e.triggeredid === this.id);
  }

  getFinishedEvent() {
    return this.type.endsWith('.finished') ? this : this.traces.find(t => t.type.endsWith('.finished'));
  }

  getRemediationAction() {
    return this.findTrace(t => t.isRemediationAction());
  }

  getEvaluation(stageName: string): Trace | undefined {
    return this.findTrace(t => !!t.isEvaluation() && t.stage === stageName);
  }

  getDeploymentUrl() {
    return this.data.deployment?.deploymentURIsPublic?.find(() => true);
  }

  findTrace(comp: (args: Trace) => boolean): Trace | undefined {
    if (comp(this)) {
      return this;
    }
    else {
      return this.traces.reduce((result: Trace | undefined, trace: Trace) => result || trace.findTrace(comp), undefined);
    }
  }

  findLastTrace(comp: (args: Trace) => boolean) {
    if (comp(this)) {
      return this;
    }
    else {
      return this.traces.reduce((result: Trace | undefined, trace) => trace.findTrace(comp) || result, undefined);
    }
  }

  getLastTrace(): Trace {
    return this.traces.length ? this.traces[this.traces.length - 1].getLastTrace() : this;
  }

  isSequence(): boolean {
    return this.type.split('.').length === 6 && !!this.stage && this.type.includes(this.stage);
  }

  getProblemDetails() {
    return this.data.problem?.ImpactedEntity || this.data.problem?.ProblemTitle;
  }
}

export {Trace};
