import {EventTypes} from "./event-types";
import {ResultTypes} from "./result-types";
import {ApprovalStates} from "./approval-states";
import {EVENT_LABELS} from "./event-labels";
import {EVENT_ICONS} from "./event-icons";

enum ProblemStates {
  RESOLVED = 'RESOLVED'
};

const DEFAULT_ICON = "information";

class Trace {
  id: string;
  shkeptncontext: string;
  source: string;
  time: Date;
  type: string;
  label: string;
  icon: string;
  image: string;
  plainEvent: string;
  data: {
    project: string;
    service: string;
    stage: string;

    image: string;
    tag: string;

    deploymentURILocal: string;
    deploymentURIPublic: string;

    deploymentstrategy: string;
    labels: Map<any, any>;
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
    valuesCanary: {
      image: string;
    };

    evaluationdetails: {
      indicatorResults: any;
      result: string;
      score: number;
      sloFileContent: string;
      timeEnd: Date;
      timeStart: Date;
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

    Tags: string;
    State: string;
  };

  isFaulty(): string {
    let result: string = null;
    if(this.data) {
      if(this.isFailed() || this.isProblem()) {
        result = this.data.stage;
      }
    }
    return result;
  }

  isWarning(): string {
    let result: string = null;
    if(this.data) {
      if(this.data.result == ResultTypes.WARNING) {
        result = this.data.stage;
      }
    }
    return result;
  }

  isFailed(): boolean {
    return this.data.result == ResultTypes.FAILED || this.type === EventTypes.APPROVAL_FINISHED && this.data.approval.result == ApprovalStates.DECLINED;
  }

  isProblem(): boolean {
    return this.type.indexOf('problem') != -1;
  }

  isSuccessful(): boolean {
    let result: boolean = false;
    if(this.data) {
      if(this.data.result == ResultTypes.PASSED || this.type === EventTypes.APPROVAL_FINISHED && this.data.approval.result == ApprovalStates.APPROVED) {
        result = true;
      }
    }
    return !this.isFaulty() && result;
  }

  getLabel(): string {
    // TODO: use translation file
    if(!this.label) {
      if(this.type === EventTypes.PROBLEM_DETECTED && this.data.State === ProblemStates.RESOLVED) {
        this.label = EVENT_LABELS[EventTypes.PROBLEM_RESOLVED];
      } else {
        this.label = EVENT_LABELS[this.type] || this.type;
      }
    }

    return this.label;
  }

  getIcon() {
    if(!this.icon) {
      this.icon = EVENT_ICONS[this.type] || DEFAULT_ICON;
    }
    return this.icon;
  }

  getShortImageName() {
    if(!this.image) {
      if(this.data.image && this.data.tag)
        this.image = [this.data.image.split("/").pop(), this.data.tag].join(":");
      else if(this.data.image)
        this.image = this.data.image.split("/").pop();
      else if(this.data.valuesCanary)
        this.image = this.data.valuesCanary.image.split("/").pop();
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
    return this.data.labels?.get("buildId") ?? this.time;
  }

  static fromJSON(data: any) {
    return Object.assign(new this, data, { plainEvent: JSON.parse(JSON.stringify(data)) });
  }
}

export {Trace}
