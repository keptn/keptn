import * as moment from 'moment';

import {EventTypes} from "./event-types";
import {ResultTypes} from "./result-types";
import {ApprovalStates} from "./approval-states";
import {EVENT_LABELS} from "./event-labels";
import {EVENT_ICONS} from "./event-icons";
import {ProblemStates} from "./problem-states";


const DEFAULT_ICON = "information";

class Trace {
  traces: Trace[] = [];

  id: string;
  shkeptncontext: string;
  triggeredid: string;
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
      if(this.type === EventTypes.EVALUATION_FINISHED && this.isFailed()) {
        result = this.data.stage;
      }
    }
    return result;
  }

  isWarning(): string {
    let result: string = null;
    if(this.data) {
      if(this.type === EventTypes.EVALUATION_FINISHED && this.data.result == ResultTypes.WARNING) {
        result = this.data.stage;
      }
    }
    return result;
  }

  isSuccessful(): string {
    let result: string = null;
    if(this.data) {
      if(this.data.result == ResultTypes.PASSED || this.isApprovalFinished() && this.isApproved() || this.isProblem() && this.isProblemResolvedOrClosed() || this.isSuccessfulRemediation()) {
        result = this.data.stage;
      }
    }
    return !this.isFaulty() && result ? result : null;
  }

  public isFailed(): boolean {
    return this.data.result == ResultTypes.FAILED || this.isApprovalFinished() && this.isDeclined();
  }

  public isProblem(): boolean {
    return this.type === EventTypes.PROBLEM_DETECTED || this.type === EventTypes.PROBLEM_OPEN;
  }

  public isProblemResolvedOrClosed(): boolean {
    return this.data.State === ProblemStates.RESOLVED || this.data.State === ProblemStates.CLOSED;
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
    return this.type === EventTypes.EVALUATION_TRIGGERED ? this.data.stage : null;
  }

  public isEvaluationInvalidation(): boolean {
    return this.type === EventTypes.EVALUATION_INVALIDATED;
  }

  hasLabels(): boolean {
    return Object.keys(this.data.labels||{}).length > 0;
  }

  getLabel(): string {
    // TODO: use translation file
    if(!this.label) {
      if(this.isProblem() && this.isProblemResolvedOrClosed()) {
        this.label = EVENT_LABELS[EventTypes.PROBLEM_RESOLVED];
      } else if(this.isApprovalFinished()) {
        this.label = EVENT_LABELS[EventTypes.APPROVAL_FINISHED][this.data.approval?.result] || this.getShortType();
      } else {
        this.label = EVENT_LABELS[this.type] || this.getShortType();
      }
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

  isFinished() {
    if(!this.finished) {
      if(!this.traces || this.traces.length == 0)
        this.finished = this.type.includes(".finished");
      else
        this.finished = this.traces.some(t => t.type.includes(".finished"));
    }

    return this.finished;
  }

  getFinishedEvent() {
    return this.traces.find(t => t.type.includes(".finished"));
  }

  static fromJSON(data: any) {
    if(data instanceof Trace)
      return data;

    const plainEvent = JSON.parse(JSON.stringify(data));
    return Object.assign(new this, data, { plainEvent });
  }
}

export {Trace}
