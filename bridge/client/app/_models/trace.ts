import * as moment from 'moment';

import {EventTypes} from "./event-types";
import {ResultTypes} from "./result-types";
import {ApprovalStates} from "./approval-states";
import {EVENT_LABELS} from "./event-labels";
import {EVENT_ICONS} from "./event-icons";
import {ProblemStates} from "./problem-states";
import {DateUtil} from '../_utils/date.utils';

const DEFAULT_ICON = "information";

class Trace {
  traces: Trace[] = [];

  id: string;
  shkeptncontext: string;
  triggeredid: string;
  started: boolean;
  finished: boolean;
  source: string;
  time: Date;
  type: string;
  label: string;
  heatmapLabel: string;
  icon: string;
  image: string;
  plainEvent: string;
  data: {
    project: string;
    service: string;
    stage: string;

    image: string;
    tag: string;

    deployment: {
      deploymentNames: string[];
      deploymentURIsLocal: string[];
      deploymentURIsPublic: string[];
      deploymentstrategy: string;
      gitCommit: string;
    };

    deploymentURILocal: string;
    deploymentURIPublic: string;

    message: string;

    labels: Map<string, string>;
    result: string;
    teststrategy: string;

    start: Date;
    end: Date;

    canary: {
      action: string;
      value: number;
    };
    eventContext: {
      shkeptncontext: string;
      token: string;
    };
    configurationChange: {
      values: {
        image: string
      }
    };

    evaluation: {
      comparedEvents: string[];
      indicatorResults: any;
      result: string;
      score: number;
      sloFileContent: string;
      timeEnd: Date;
      timeStart: Date;

      score_pass: any;
      score_warning: any;
      compare_with: string;
      include_result_with_score: string;
      number_of_comparison_results: number;
      number_of_missing_comparison_results: number;
      sloFileContentParsed: string;
    };

    evaluationHistory: Trace[];

    ProblemTitle: string;
    ImpactedEntity: string;
    ProblemDetails: {
      tagsOfAffectedEntities: {
        key: string;
        value: string;
      }
    };

    approval: {
      result: string;
      status: string;
    };

    action: {
      action: string;
      description: string;
      name: string;
    }

    Tags: string;
    State: string;
  };

  isFaulty(): string {
    let result: string = null;
    if(this.data) {
      if(this.isFailed() ||
        (this.isProblem() && !this.isProblemResolvedOrClosed()) ||
        this.traces.some(t => t.isFailed() || (t.isProblem() && !t.isProblemResolvedOrClosed()))) {
        result = this.data.stage;
      }
    }
    return result;
  }

  isFailedEvaluation() {
    let result: string = null;
    if(this.data) {
      if(this.getFinishedEvent()?.type === EventTypes.EVALUATION_FINISHED && this.isFailed()) {
        result = this.data.stage;
      }
    }
    return result;
  }

  isWarning(): string {
    let result: string = null;
    if(this.getFinishedEvent()?.data.result == ResultTypes.WARNING) {
      result = this.data.stage;
    }
    return result;
  }

  isSuccessful(): string {
    let result: string = null;
    if ( this.isFinished() && this.getFinishedEvent()?.data.result === ResultTypes.PASSED || this.isApprovalFinished() && this.isApproved() || this.isProblem() && this.isProblemResolvedOrClosed() || this.isSuccessfulRemediation()) {
      result = this.data.stage;
    }
    return !this.isFaulty() && result ? result : null;
  }

  public isFailed(): boolean {
    return this.getFinishedEvent()?.data.result == ResultTypes.FAILED || this.isApprovalFinished() && this.isDeclined();
  }

  public isProblem(): boolean {
    return this.type === EventTypes.PROBLEM_DETECTED || this.type === EventTypes.PROBLEM_OPEN;
  }

  public isRemediation(): boolean {
    return this.type.endsWith(EventTypes.REMEDIATION_TRIGGERED_SUFFIX);
  }

  public isProblemResolvedOrClosed(): boolean {
    if (!this.traces || this.traces.length === 0)
      return this.data.State === ProblemStates.RESOLVED || this.data.State === ProblemStates.CLOSED;
    else
      return this.traces.some(t => t.isProblem() && t.isProblemResolvedOrClosed());
  }

  public isSuccessfulRemediation(): boolean {
    return this.type === EventTypes.REMEDIATION_FINISHED && this.data.result == ResultTypes.PASSED;
  }

  public isApproval(): string {
    return this.type === EventTypes.APPROVAL_TRIGGERED ? this.data.stage : null;
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
    return this.type === EventTypes.DEPLOYMENT_FINISHED && this.data?.deployment?.deploymentstrategy == "direct";
  }

  private isApproved(): boolean {
    return this.data.approval?.result == ApprovalStates.APPROVED;
  }

  private isDeclined(): boolean {
    return this.data.approval?.result == ApprovalStates.DECLINED;
  }

  public isDeployment(): string {
    return this.type === EventTypes.DEPLOYMENT_TRIGGERED ? this.data.stage : null;
  }

  public isEvaluation(): string {
    return this.type.endsWith(EventTypes.EVALUATION_TRIGGERED_SUFFIX) && !this.isSequence() ? this.data.stage : null;
  }

  public isEvaluationInvalidation(): boolean {
    return this.type === EventTypes.EVALUATION_INVALIDATED;
  }

  hasLabels(): boolean {
    return Object.keys(this.data.labels||{}).length > 0;
  }

  getLabel(): string {
    // TODO: use translation file; see also EVENT_LABELS in event-labels.ts
    if(!this.label) {
      this.label = this.getShortType();
    }

    return this.label;
  }

  getStage(): string {
    return this.data?.stage;
  }

  getShortType(): string {
    let parts = this.type.split(".");
    if(parts.length == 6)
      return parts[4];
    else if(parts.length == 5)
      return parts[3];
    else
      return this.type;
  }

  getIcon(): string {
    if(!this.icon) {
      this.icon = EVENT_ICONS[this.getShortType()] || DEFAULT_ICON;
    }
    return this.icon;
  }

  getShortImageName() {
    if(!this.image) {
      if(this.data.image && this.data.tag)
        this.image = [this.data.image.split("/").pop(), this.data.tag].join(":");
      else if(this.data.image)
        this.image = this.data.image.split("/").pop();
      else if(this.data.configurationChange?.values)
        this.image = this.data.configurationChange.values.image?.split("/").pop();
    }

    return this.image;
  }

  getProject(): string {
    return this.data.project;
  }

  getService(): string {
    return this.data.service;
  }

  getChartLabel(): string {
    return this.data.labels?.["buildId"] ?? moment(this.time).format("YYYY-MM-DD HH:mm");
  }

  getHeatmapLabel(): string {
    if(!this.heatmapLabel) {
      this.heatmapLabel = this.getChartLabel();
    }
    return this.heatmapLabel;
  }

  setHeatmapLabel(label: string) {
    this.heatmapLabel = label;
  }

  isStarted() {
    if(!this.started && this.traces) {
      this.started = this.traces.some(t => t.type.endsWith('.started') || t.isStarted());
    }

    return this.started;
  }

  isChanged() {
    return this.type.endsWith('.changed')
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
    return !!this.traces.find(e => e.isEvaluationInvalidation() && e.triggeredid == this.id);
  }

  getFinishedEvent() {
    return this.type.endsWith('.finished') ? this : this.traces.find(t => t.type.endsWith('.finished'));
  }

  getDeploymentUrl() {
    return this.data.deployment?.deploymentURIsPublic?.find(e => true);
  }

  findTrace(comp: <T = Trace>(args: Trace) => any): Trace {
    if (comp(this))
      return this;
    else
      return this.traces.reduce((result, trace) => result || trace.findTrace(comp), null);
  }

  findLastTrace(comp) {
    if(comp(this))
      return this;
    else
      return this.traces.reduce((result, trace) => trace.findTrace(comp) || result, null);
  }

  getLastTrace(): Trace {
    return this.traces.length ? this.traces[this.traces.length - 1].getLastTrace() : this;
  }

  isSequence() {
    return this.type.split(".").length == 6 && this.type.includes(this.getStage());
  }

  static fromJSON(data: any) {
    if(data instanceof Trace)
      return data;

    const plainEvent = JSON.parse(JSON.stringify(data));
    return Object.assign(new this, data, { plainEvent });
  }

  static traceMapper(traces: Trace[]) {
    traces = traces
      .map(trace => Trace.fromJSON(trace))
      .sort(DateUtil.compareTraceTimesDesc);

    return traces.reduce((seq: Trace[], trace: Trace) => {
      let trigger: Trace = null;
      if(trace.triggeredid) {
        trigger = traces.reduce((trigger, r) => trigger || r.findTrace((t) => t.id == trace.triggeredid), null);
      } else if(trace.isProblem() && trace.isProblemResolvedOrClosed()) {
        trigger = traces.reduce((trigger, r) => trigger || r.findTrace((t) => t.isProblem() && !t.isProblemResolvedOrClosed()), null);
      } else if(trace.isFinished()) {
        trigger = traces.reduce((trigger, r) => trigger || r.findTrace((t) => !t.triggeredid && t.type.slice(0, -8) === trace.type.slice(0, -9)), null);
      }

      if (trigger) {
        trigger.traces.push(trace);
      } else if (trace.isSequence()) {
        seq.push(trace);
      } else if(seq.length > 0) {
        seq[seq.length-1].traces.push(trace);
      } else {
        seq.push(trace);
      }

      return seq;
    }, []);
  }
}

export {Trace}
